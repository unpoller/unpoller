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

// PollController runs forever, polling unifi, and pushing to influx.
func (u *UnifiPoller) PollController() error {
	log.Println("[INFO] Everything checks out! Poller started, interval:", u.Interval.Round(time.Second))
	ticker := time.NewTicker(u.Interval.Round(time.Second))
	var err error
	for range ticker.C {
		m := &Metrics{}
		// Get the sites we care about.
		if m.Sites, err = u.GetFilteredSites(); err != nil {
			u.LogErrors([]error{err}, "unifi.GetSites()")
		}
		// Get all the points.
		if m.Clients, err = u.GetClients(m.Sites); err != nil {
			u.LogErrors([]error{err}, "unifi.GetClients()")
		}
		if m.Devices, err = u.GetDevices(m.Sites); err != nil {
			u.LogErrors([]error{err}, "unifi.GetDevices()")
		}

		// Make a new Points Batcher.
		m.BatchPoints, err = influx.NewBatchPoints(influx.BatchPointsConfig{Database: u.InfluxDB})
		if err != nil {
			u.LogErrors([]error{err}, "influx.NewBatchPoints")
			continue
		}
		// Batch (and send) all the points.
		if errs := m.SendPoints(); errs != nil && hasErr(errs) {
			u.LogErrors(errs, "asset.Points()")
		}
		if err := u.Write(m.BatchPoints); err != nil {
			u.LogErrors([]error{err}, "infdb.Write(bp)")
		}

		// Talk about the data.
		var fieldcount, pointcount int
		for _, p := range m.Points() {
			pointcount++
			i, _ := p.Fields()
			fieldcount += len(i)
		}
		u.Logf("Unifi Measurements Recorded. Sites: %d, Clients: %d, "+
			"Wireless APs: %d, Gateways: %d, Switches: %d, Points: %d, Fields: %d",
			len(m.Sites), len(m.Clients), len(m.UAPs), len(m.USGs), len(m.USWs), pointcount, fieldcount)

		if u.MaxErrors >= 0 && u.errorCount > u.MaxErrors {
			return errors.Errorf("reached maximum error count, stopping poller (%d > %d)", u.errorCount, u.MaxErrors)
		}
	}
	return nil
}

// SendPoints combines all device and client data into influxdb data points.
// Call this after you've collected all the data you care about.
// This sends all the batched points to InfluxDB.
func (m *Metrics) SendPoints() (errs []error) {
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

// processPoints is helper function for SendPoints.
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
