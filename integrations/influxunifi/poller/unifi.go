package poller

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/davidnewhall/unifi-poller/influx"
	"github.com/davidnewhall/unifi-poller/metrics"
	"github.com/davidnewhall/unifi-poller/prometheus"
	client "github.com/influxdata/influxdb1-client/v2"
	"golift.io/unifi"
)

// CheckSites makes sure the list of provided sites exists on the controller.
// This does not run in Lambda (run-once) mode.
func (u *UnifiPoller) CheckSites() error {
	if strings.Contains(strings.ToLower(u.Config.Mode), "lambda") {
		return nil // Skip this in lambda mode.
	}
	u.LogDebugf("Checking Controller Sites List")
	sites, err := u.Unifi.GetSites()
	if err != nil {
		return err
	}
	msg := []string{}
	for _, site := range sites {
		msg = append(msg, site.Name+" ("+site.Desc+")")
	}
	u.Logf("Found %d site(s) on controller: %v", len(msg), strings.Join(msg, ", "))
	if StringInSlice("all", u.Config.Sites) {
		u.Config.Sites = []string{"all"}
		return nil
	}
FIRST:
	for _, s := range u.Config.Sites {
		for _, site := range sites {
			if s == site.Name {
				continue FIRST
			}
		}
		// This is fine, it may get added later.
		u.LogErrorf("configured site not found on controller: %v", s)
	}
	return nil
}

// PollController runs forever, polling UniFi
// and pushing to influx OR exporting for prometheus.
// This is started by Run() after everything checks out.
func (u *UnifiPoller) PollController(process func(*metrics.Metrics) error) error {
	interval := u.Config.Interval.Round(time.Second)
	log.Printf("[INFO] Everything checks out! Poller started in %v mode, interval: %v", u.Config.Mode, interval)
	ticker := time.NewTicker(interval)
	for u.LastCheck = range ticker.C {
		var err error
		if u.Config.ReAuth {
			u.LogDebugf("Re-authenticating to UniFi Controller")
			// Some users need to re-auth every interval because the cookie times out.
			if err = u.Unifi.Login(); err != nil {
				u.LogError(err, "re-authenticating")
			}
		}
		if err == nil {
			// Only run this if the authentication procedure didn't return error.
			_ = u.CollectAndProcess(process)
		}
		if u.errorCount > 0 {
			return fmt.Errorf("too many errors, stopping poller")
		}
	}
	return nil
}

// CollectAndProcess collects measurements and then passese them into the
// provided method. The method is either an http exporter or an influxdb update.
// Can be called once or in a ticker loop. This function and all the ones below
// handle their own logging. An error is returned so the calling function may
// determine if there was a read or write error and act on it. This is currently
// called in two places in this library. One returns an error, one does not.
func (u *UnifiPoller) CollectAndProcess(process func(*metrics.Metrics) error) error {
	metrics, err := u.CollectMetrics()
	if err != nil {
		return err
	}
	if err := u.AugmentMetrics(metrics); err != nil {
		return err
	}
	err = process(metrics)
	u.LogError(err, "processing metrics")
	return err
}

// CollectMetrics grabs all the measurements from a UniFi controller and returns them.
func (u *UnifiPoller) CollectMetrics() (*metrics.Metrics, error) {
	m := &metrics.Metrics{TS: u.LastCheck} // At this point, it's the Current Check.
	var err error
	// Get the sites we care about.
	m.Sites, err = u.GetFilteredSites()
	u.LogError(err, "unifi.GetSites()")
	if u.Config.CollectIDS {
		// Check back in time since twice the interval. Dups are discarded by InfluxDB.
		m.IDSList, err = u.Unifi.GetIDS(m.Sites, time.Now().Add(2*u.Config.Interval.Duration), time.Now())
		u.LogError(err, "unifi.GetIDS()")
	}
	// Get all the points.
	m.Clients, err = u.Unifi.GetClients(m.Sites)
	u.LogError(err, "unifi.GetClients()")
	m.Devices, err = u.Unifi.GetDevices(m.Sites)
	u.LogError(err, "unifi.GetDevices()")
	return m, err
}

// AugmentMetrics is our middleware layer between collecting metrics and writing them.
// This is where we can manipuate the returned data or make arbitrary decisions.
// This function currently adds parent device names to client metrics.
func (u *UnifiPoller) AugmentMetrics(metrics *metrics.Metrics) error {
	if metrics == nil || metrics.Devices == nil || metrics.Clients == nil {
		return fmt.Errorf("nil metrics, augment impossible")
	}
	devices := make(map[string]string)
	bssdIDs := make(map[string]string)
	for _, r := range metrics.UAPs {
		devices[r.Mac] = r.Name
		for _, v := range r.VapTable {
			bssdIDs[v.Bssid] = fmt.Sprintf("%s %s %s:", r.Name, v.Radio, v.RadioName)
		}
	}
	for _, r := range metrics.USGs {
		devices[r.Mac] = r.Name
	}
	for _, r := range metrics.USWs {
		devices[r.Mac] = r.Name
	}
	for _, r := range metrics.UDMs {
		devices[r.Mac] = r.Name
	}
	// These come blank, so set them here.
	for i, c := range metrics.Clients {
		metrics.Clients[i].SwName = devices[c.SwMac]
		metrics.Clients[i].ApName = devices[c.ApMac]
		metrics.Clients[i].GwName = devices[c.GwMac]
		metrics.Clients[i].RadioDescription = bssdIDs[metrics.Clients[i].Bssid] + metrics.Clients[i].RadioProto
	}
	return nil
}

// ExportMetrics updates the internal metrics provided via
// HTTP at /metrics for prometheus collection.
func (u *UnifiPoller) ExportMetrics(metrics *metrics.Metrics) error {
	m := &prometheus.Metrics{Metrics: metrics}
	for _, err := range m.ProcessExports() {
		u.LogError(err, "prometheus.ProcessExports")
	}
	u.LogExportReport(m)
	return nil
}

// LogExportReport writes a log line after exporting metrics via HTTP.
func (u *UnifiPoller) LogExportReport(m *prometheus.Metrics) {
	idsMsg := ""
	if u.Config.CollectIDS {
		idsMsg = fmt.Sprintf(", IDS Events: %d, ", len(m.IDSList))
	}
	u.Logf("UniFi Measurements Exported. Sites: %d, Clients: %d, "+
		"Wireless APs: %d, Gateways: %d, Switches: %d%s",
		len(m.Sites), len(m.Clients), len(m.UAPs),
		len(m.UDMs)+len(m.USGs), len(m.USWs), idsMsg)
}

// ReportMetrics batches all the metrics and writes them to InfluxDB.
// This creates an InfluxDB writer, and returns an error if the write fails.
func (u *UnifiPoller) ReportMetrics(metrics *metrics.Metrics) error {
	// Batch (and send) all the points.
	m := &influx.Metrics{Metrics: metrics}
	// Make a new Influx Points Batcher.
	var err error
	m.BatchPoints, err = client.NewBatchPoints(client.BatchPointsConfig{Database: u.Config.InfluxDB})
	if err != nil {
		return fmt.Errorf("influx.NewBatchPoints: %v", err)
	}
	for _, err := range m.ProcessPoints() {
		u.LogError(err, "influx.ProcessPoints")
	}
	if err = u.Influx.Write(m.BatchPoints); err != nil {
		return fmt.Errorf("influxdb.Write(points): %v", err)
	}
	u.LogInfluxReport(m)
	return nil
}

// LogInfluxReport writes a log message after exporting to influxdb.
func (u *UnifiPoller) LogInfluxReport(m *influx.Metrics) {
	var fields, points int
	for _, p := range m.Points() {
		points++
		i, _ := p.Fields()
		fields += len(i)
	}
	idsMsg := ""
	if u.Config.CollectIDS {
		idsMsg = fmt.Sprintf("IDS Events: %d, ", len(m.IDSList))
	}
	u.Logf("UniFi Measurements Recorded. Sites: %d, Clients: %d, "+
		"Wireless APs: %d, Gateways: %d, Switches: %d, %sPoints: %d, Fields: %d",
		len(m.Sites), len(m.Clients), len(m.UAPs),
		len(m.UDMs)+len(m.USGs), len(m.USWs), idsMsg, points, fields)
}

// GetFilteredSites returns a list of sites to fetch data for.
// Omits requested but unconfigured sites. Grabs the full list from the
// controller and returns the sites provided in the config file.
func (u *UnifiPoller) GetFilteredSites() (unifi.Sites, error) {
	sites, err := u.Unifi.GetSites()
	if err != nil {
		return nil, err
	} else if len(u.Config.Sites) < 1 || StringInSlice("all", u.Config.Sites) {
		return sites, nil
	}
	var i int
	for _, s := range sites {
		// Only include valid sites in the request filter.
		if StringInSlice(s.Name, u.Config.Sites) {
			sites[i] = s
			i++
		}
	}
	return sites[:i], nil
}
