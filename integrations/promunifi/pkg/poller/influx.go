package poller

import (
	"fmt"
	"log"
	"time"

	"github.com/davidnewhall/unifi-poller/pkg/influxunifi"
)

// GetInfluxDB returns an InfluxDB interface.
func (u *UnifiPoller) GetInfluxDB() (err error) {
	if u.Influx != nil {
		return nil
	}

	u.Influx, err = influxunifi.New(&influxunifi.Config{
		Database: u.Config.InfluxDB,
		User:     u.Config.InfluxUser,
		Pass:     u.Config.InfluxPass,
		BadSSL:   u.Config.InfxBadSSL,
		URL:      u.Config.InfluxURL,
	})
	if err != nil {
		return fmt.Errorf("influxdb: %v", err)
	}

	u.Logf("Logging Measurements to InfluxDB at %s as user %s", u.Config.InfluxURL, u.Config.InfluxUser)

	return nil
}

// PollController runs forever, polling UniFi and pushing to InfluxDB
// This is started by Run() or RunBoth() after everything checks out.
func (u *UnifiPoller) PollController() {
	interval := u.Config.Interval.Round(time.Second)
	log.Printf("[INFO] Everything checks out! Poller started, InfluxDB interval: %v", interval)

	ticker := time.NewTicker(interval)
	for u.LastCheck = range ticker.C {
		if err := u.CollectAndProcess(); err != nil {
			u.LogErrorf("%v", err)
		}
	}
}

// CollectAndProcess collects measurements and then reports them to InfluxDB
// Can be called once or in a ticker loop. This function and all the ones below
// handle their own logging. An error is returned so the calling function may
// determine if there was a read or write error and act on it. This is currently
// called in two places in this library. One returns an error, one does not.
func (u *UnifiPoller) CollectAndProcess() error {
	if err := u.GetInfluxDB(); err != nil {
		return err
	}

	metrics, err := u.CollectMetrics()
	if err != nil {
		return err
	}

	report, err := u.Influx.ReportMetrics(metrics)
	if err != nil {
		return err
	}

	u.LogInfluxReport(report)
	return nil
}

// LogInfluxReport writes a log message after exporting to influxdb.
func (u *UnifiPoller) LogInfluxReport(r *influxunifi.Report) {
	idsMsg := fmt.Sprintf("IDS Events: %d, ", len(r.Metrics.IDSList))
	u.Logf("UniFi Metrics Recorded. Sites: %d, Clients: %d, "+
		"UAP: %d, USG/UDM: %d, USW: %d, %sPoints: %d, Fields: %d, Errs: %d, Elapsed: %v",
		len(r.Metrics.Sites), len(r.Metrics.Clients), len(r.Metrics.UAPs),
		len(r.Metrics.UDMs)+len(r.Metrics.USGs), len(r.Metrics.USWs), idsMsg, r.Total,
		r.Fields, len(r.Errors), r.Elapsed.Round(time.Millisecond))
}
