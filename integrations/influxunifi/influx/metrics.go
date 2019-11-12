// Package influx provides the methods to turn UniFi measurements into influx
// data-points with appropriate tags and fields.
package influx

import (
	"time"

	client "github.com/influxdata/influxdb1-client/v2"
	"golift.io/unifi"
)

// Metrics contains all the data from the controller and an influx endpoint to send it to.
type Metrics struct {
	TS time.Time
	unifi.Sites
	unifi.IDSList
	unifi.Clients
	*unifi.Devices
	client.BatchPoints
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
	processPoints := func(m *Metrics, p []*client.Point, err error) {
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
