package unifipoller

import (
	"fmt"
	"log"
	"strings"
	"time"

	influx "github.com/influxdata/influxdb1-client/v2"
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

// PollController runs forever, polling UniFi, and pushing to influx.
// This is started by Run() after everything checks out.
func (u *UnifiPoller) PollController() error {
	interval := u.Config.Interval.Round(time.Second)
	log.Println("[INFO] Everything checks out! Poller started, interval:", interval)
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
			_ = u.CollectAndReport()
		}
		if u.Config.MaxErrors >= 0 && u.errorCount > u.Config.MaxErrors {
			return fmt.Errorf("reached maximum error count, stopping poller (%d > %d)",
				u.errorCount, u.Config.MaxErrors)
		}
	}
	return nil
}

// CollectAndReport collects measurements and reports them to influxdb.
// Can be called once or in a ticker loop. This function and all the ones below
// handle their own logging. An error is returned so the calling function may
// determine if there was a read or write error and act on it. This is currently
// called in two places in this library. One returns an error, one does not.
func (u *UnifiPoller) CollectAndReport() error {
	metrics, err := u.CollectMetrics()
	if err != nil {
		return err
	}
	if err := u.AugmentMetrics(metrics); err != nil {
		return err
	}
	err = u.ReportMetrics(metrics)
	u.LogError(err, "reporting metrics")
	return err
}

// CollectMetrics grabs all the measurements from a UniFi controller and returns them.
// This also creates an InfluxDB writer, and returns an error if that fails.
func (u *UnifiPoller) CollectMetrics() (*Metrics, error) {
	m := &Metrics{TS: u.LastCheck} // At this point, it's the Current Check.
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
	// Make a new Influx Points Batcher.
	m.BatchPoints, err = influx.NewBatchPoints(influx.BatchPointsConfig{Database: u.Config.InfluxDB})
	u.LogError(err, "influx.NewBatchPoints")
	return m, err
}

// AugmentMetrics is our middleware layer between collecting metrics and writing them.
// This is where we can manipuate the returned data or make arbitrary decisions.
// This function currently adds parent device names to client metrics.
func (u *UnifiPoller) AugmentMetrics(metrics *Metrics) error {
	devices := make(map[string]string)
	for _, r := range metrics.UAPs {
		devices[r.Mac] = r.Name
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
	}
	return nil
}

// ReportMetrics batches all the metrics and writes them to InfluxDB.
// Returns an error if the write to influx fails.
func (u *UnifiPoller) ReportMetrics(metrics *Metrics) error {
	// Batch (and send) all the points.
	for _, err := range metrics.ProcessPoints() {
		u.LogError(err, "asset.Points()")
	}
	err := u.Influx.Write(metrics.BatchPoints)
	if err != nil {
		return fmt.Errorf("influxdb.Write(points): %v", err)
	}
	var fields, points int
	for _, p := range metrics.Points() {
		points++
		i, _ := p.Fields()
		fields += len(i)
	}
	idsMsg := ""
	if u.Config.CollectIDS {
		idsMsg = fmt.Sprintf("IDS Events: %d, ", len(metrics.IDSList))
	}
	u.Logf("UniFi Measurements Recorded. Sites: %d, Clients: %d, "+
		"Wireless APs: %d, Gateways: %d, Switches: %d, %sPoints: %d, Fields: %d",
		len(metrics.Sites), len(metrics.Clients), len(metrics.UAPs),
		len(metrics.UDMs)+len(metrics.USGs), len(metrics.USWs), idsMsg, points, fields)
	return nil
}

// ProcessPoints batches all device and client data into influxdb data points.
// Call this after you've collected all the data you care about.
// This function is sorta weird and returns a slice of errors. The reasoning is
// that some points may process while others fail, so we attempt to process them
// all. This is (usually) run in a loop, so we can't really exit on error,
// we just log the errors and tally them on a counter. In reality, this never
// returns any errors because we control the data going in; cool right? But we
// still check&log it in case the data going is skewed up and causes errors!
func (m *Metrics) ProcessPoints() []error {
	errs := []error{}
	processPoints := func(m *Metrics, p []*influx.Point, err error) {
		switch {
		case err != nil:
			errs = append(errs, err)
		case p == nil:
		default:
			m.BatchPoints.AddPoints(p)
		}
	}

	for _, asset := range m.Sites {
		pts, err := SitePoints(asset, m.TS)
		processPoints(m, pts, err)
	}
	for _, asset := range m.Clients {
		pts, err := ClientPoints(asset, m.TS)
		processPoints(m, pts, err)
	}
	for _, asset := range m.IDSList {
		pts, err := IDSPoints(asset) // no m.TS.
		processPoints(m, pts, err)
	}

	if m.Devices == nil {
		return errs
	}
	for _, asset := range m.Devices.UAPs {
		pts, err := UAPPoints(asset, m.TS)
		processPoints(m, pts, err)
	}
	for _, asset := range m.Devices.USGs {
		pts, err := USGPoints(asset, m.TS)
		processPoints(m, pts, err)
	}
	for _, asset := range m.Devices.USWs {
		pts, err := USWPoints(asset, m.TS)
		processPoints(m, pts, err)
	}
	for _, asset := range m.Devices.UDMs {
		pts, err := UDMPoints(asset, m.TS)
		processPoints(m, pts, err)
	}
	return errs
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
