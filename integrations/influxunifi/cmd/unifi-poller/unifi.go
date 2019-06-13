package main

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
func (u *UnifiPoller) PollController() {
	log.Println("[INFO] Everything checks out! Poller started, interval:", u.Interval.value)
	ticker := time.NewTicker(u.Interval.value)

	for range ticker.C {
		// Get the sites we care about.
		sites, err := u.filterSites(u.Sites)
		if err != nil {
			logErrors([]error{err}, "uni.GetSites()")
		}
		// Get all the points.
		clients, err := u.GetClients(sites)
		if err != nil {
			logErrors([]error{err}, "uni.GetClients()")
		}
		devices, err := u.GetDevices(sites)
		if err != nil {
			logErrors([]error{err}, "uni.GetDevices()")
		}
		// Make a new Points Batcher.
		bp, err := influx.NewBatchPoints(influx.BatchPointsConfig{Database: u.InfluxDB})
		if err != nil {
			logErrors([]error{err}, "influx.NewBatchPoints")
			continue
		}
		// Batch (and send) all the points.
		if errs := batchPoints(devices, clients, bp); errs != nil && hasErr(errs) {
			logErrors(errs, "asset.Points()")
		}
		if err := u.Write(bp); err != nil {
			logErrors([]error{err}, "infdb.Write(bp)")
		}

		// Talk about the data.
		var fieldcount, pointcount int
		for _, p := range bp.Points() {
			pointcount++
			i, _ := p.Fields()
			fieldcount += len(i)
		}
		u.Logf("Unifi Measurements Recorded. Sites: %d Clients: %d, "+
			"Wireless APs: %d, Gateways: %d, Switches: %d, Points: %d, Fields: %d",
			len(sites), len(clients.UCLs),
			len(devices.UAPs), len(devices.USGs), len(devices.USWs), pointcount, fieldcount)
	}
}

// batchPoints combines all device and client data into influxdb data points.
func batchPoints(devices *unifi.Devices, clients *unifi.Clients, bp influx.BatchPoints) (errs []error) {
	process := func(asset Asset) error {
		if asset == nil {
			return nil
		}
		influxPoints, err := asset.Points()
		if err != nil {
			return err
		}
		bp.AddPoints(influxPoints)
		return nil
	}
	if devices != nil {
		for _, asset := range devices.UAPs {
			errs = append(errs, process(asset))
		}
		for _, asset := range devices.USGs {
			errs = append(errs, process(asset))
		}
		for _, asset := range devices.USWs {
			errs = append(errs, process(asset))
		}
	}
	if clients != nil {
		for _, asset := range clients.UCLs {
			errs = append(errs, process(asset))
		}
	}
	return
}

// filterSites returns a list of sites to fetch data for.
// Omits requested but unconfigured sites.
func (u *UnifiPoller) filterSites(filter []string) ([]unifi.Site, error) {
	sites, err := u.GetSites()
	if err != nil {
		return nil, err
	} else if len(filter) < 1 || StringInSlice("all", filter) {
		return sites, nil
	}
	var i int
	for _, s := range sites {
		// Only include valid sites in the request filter.
		if StringInSlice(s.Name, filter) {
			sites[i] = s
			i++
		}
	}
	return sites[:i], nil
}
