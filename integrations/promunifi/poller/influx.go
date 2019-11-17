package poller

import (
	"fmt"

	"github.com/davidnewhall/unifi-poller/influxunifi"
	"github.com/davidnewhall/unifi-poller/metrics"
	client "github.com/influxdata/influxdb1-client/v2"
)

// CollectAndProcess collects measurements and then reports them to InfluxDB
// Can be called once or in a ticker loop. This function and all the ones below
// handle their own logging. An error is returned so the calling function may
// determine if there was a read or write error and act on it. This is currently
// called in two places in this library. One returns an error, one does not.
func (u *UnifiPoller) CollectAndProcess() error {
	metrics, err := u.CollectMetrics()
	if err != nil {
		return err
	}
	u.AugmentMetrics(metrics)
	err = u.ReportMetrics(metrics)
	u.LogError(err, "processing metrics")
	return err
}

// ReportMetrics batches all the metrics and writes them to InfluxDB.
// This creates an InfluxDB writer, and returns an error if the write fails.
func (u *UnifiPoller) ReportMetrics(metrics *metrics.Metrics) error {
	// Batch (and send) all the points.
	m := &influxunifi.Metrics{Metrics: metrics}
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
func (u *UnifiPoller) LogInfluxReport(m *influxunifi.Metrics) {
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
