package unifipoller

import (
	"log"
	"strings"
	"time"

	"github.com/golift/unifi"
	influx "github.com/influxdata/influxdb1-client/v2"
	"github.com/pkg/errors"
)

// CheckSites makes sure the list of provided sites exists on the controller.
func (u *UnifiPoller) CheckSites() error {
	sites, err := u.GetSites()
	if err != nil {
		return err
	}
	msg := []string{}
	for _, site := range sites {
		msg = append(msg, site.Name+" ("+site.Desc+")")
	}
	u.Logf("Found %d site(s) on controller: %v", len(msg), strings.Join(msg, ", "))
	if StringInSlice("all", u.Sites) {
		u.Sites = []string{"all"}
		return nil
	}
FIRST:
	for _, s := range u.Sites {
		for _, site := range sites {
			if s == site.Name {
				continue FIRST
			}
		}
		return errors.Errorf("configured site not found on controller: %v", s)
	}
	return nil
}

// PollController runs forever, polling UniFi, and pushing to influx.
// This is started by Run() after everything checks out.
func (u *UnifiPoller) PollController() error {
	log.Println("[INFO] Everything checks out! Poller started, interval:", u.Interval.Round(time.Second))
	ticker := time.NewTicker(u.Interval.Round(time.Second))
	for range ticker.C {
		metrics, err := u.CollectMetrics()
		if err == nil {
			u.LogError(u.ReportMetrics(metrics), "reporting metrics")
		}
		if u.MaxErrors >= 0 && u.errorCount > u.MaxErrors {
			return errors.Errorf("reached maximum error count, stopping poller (%d > %d)", u.errorCount, u.MaxErrors)
		}
	}
	return nil
}

// CollectMetrics grabs all the measurements from a UniFi controller and returns them.
// This also creates an InfluxDB writer, and retuns error if that fails.
func (u *UnifiPoller) CollectMetrics() (*Metrics, error) {
	m := &Metrics{}
	var err error
	// Get the sites we care about.
	m.Sites, err = u.GetFilteredSites()
	u.LogError(err, "unifi.GetSites()")
	// Get all the points.
	m.Clients, err = u.GetClients(m.Sites)
	u.LogError(err, "unifi.GetClients()")
	m.Devices, err = u.GetDevices(m.Sites)
	u.LogError(err, "unifi.GetDevices()")
	// Make a new Influx Points Batcher.
	m.BatchPoints, err = influx.NewBatchPoints(influx.BatchPointsConfig{Database: u.InfluxDB})
	u.LogError(err, "influx.NewBatchPoints")
	return m, err
}

// ReportMetrics batches all the metrics and writes them to InfluxDB.
// Returns an error if the write to influx fails.
func (u *UnifiPoller) ReportMetrics(metrics *Metrics) error {
	// Batch (and send) all the points.
	for _, err := range metrics.ProcessPoints() {
		u.LogError(err, "asset.Points()")
	}
	err := u.Write(metrics.BatchPoints)
	if err != nil {
		return errors.Wrap(err, "infdb.Write(bp)")
	}
	var fields, points int
	for _, p := range metrics.Points() {
		points++
		i, _ := p.Fields()
		fields += len(i)
	}
	u.Logf("UniFi Measurements Recorded. Sites: %d, Clients: %d, "+
		"Wireless APs: %d, Gateways: %d, Switches: %d, Points: %d, Fields: %d",
		len(metrics.Sites), len(metrics.Clients), len(metrics.UAPs),
		len(metrics.USGs), len(metrics.USWs), points, fields)
	return nil
}

// ProcessPoints batches all device and client data into influxdb data points.
// Call this after you've collected all the data you care about.
func (m *Metrics) ProcessPoints() (errs []error) {
	for _, asset := range m.Sites {
		errs = append(errs, m.processPoints(asset))
	}
	for _, asset := range m.Clients {
		errs = append(errs, m.processPoints(asset))
	}
	if m.Devices == nil {
		return
	}
	for _, asset := range m.UAPs {
		errs = append(errs, m.processPoints(asset))
	}
	for _, asset := range m.USGs {
		errs = append(errs, m.processPoints(asset))
	}
	for _, asset := range m.USWs {
		errs = append(errs, m.processPoints(asset))
	}
	return
}

// processPoints is helper function for ProcessPoints.
func (m *Metrics) processPoints(asset Asset) error {
	if asset == nil {
		return nil
	}
	points, err := asset.Points()
	if err != nil || points == nil {
		return err
	}
	m.BatchPoints.AddPoints(points)
	return nil
}

// GetFilteredSites returns a list of sites to fetch data for.
// Omits requested but unconfigured sites. Grabs the full list from the
// controller and filters the sites provided in the config file.
func (u *UnifiPoller) GetFilteredSites() (unifi.Sites, error) {
	sites, err := u.GetSites()
	if err != nil {
		return nil, err
	} else if len(u.Sites) < 1 || StringInSlice("all", u.Sites) {
		return sites, nil
	}
	var i int
	for _, s := range sites {
		// Only include valid sites in the request filter.
		if StringInSlice(s.Name, u.Sites) {
			sites[i] = s
			i++
		}
	}
	return sites[:i], nil
}
