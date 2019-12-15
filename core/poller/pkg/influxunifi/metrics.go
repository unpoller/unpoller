// Package influxunifi provides the methods to turn UniFi measurements into influx
// data-points with appropriate tags and fields.
package influxunifi

import (
	"crypto/tls"
	"fmt"
	"log"
	"time"

	"github.com/davidnewhall/unifi-poller/pkg/poller"
	influx "github.com/influxdata/influxdb1-client/v2"
	conf "golift.io/config"
)

const (
	defaultInterval   = 30 * time.Second
	defaultInfluxDB   = "unifi"
	defaultInfluxUser = "unifi"
	defaultInfluxURL  = "http://127.0.0.1:8086"
)

// Config defines the data needed to store metrics in InfluxDB
type Config struct {
	Interval  conf.Duration `json:"interval,omitempty" toml:"interval,omitempty" xml:"interval" yaml:"interval"`
	Disable   bool          `json:"disable" toml:"disable" xml:"disable" yaml:"disable"`
	VerifySSL bool          `json:"verify_ssl" toml:"verify_ssl" xml:"verify_ssl" yaml:"verify_ssl"`
	URL       string        `json:"url,omitempty" toml:"url,omitempty" xml:"url" yaml:"url"`
	User      string        `json:"user,omitempty" toml:"user,omitempty" xml:"user" yaml:"user"`
	Pass      string        `json:"pass,omitempty" toml:"pass,omitempty" xml:"pass" yaml:"pass"`
	DB        string        `json:"db,omitempty" toml:"db,omitempty" xml:"db" yaml:"db"`
}

// InfluxDB allows the data to be nested in the config file.
type InfluxDB struct {
	Config Config `json:"influxdb" toml:"influxdb" xml:"influxdb" yaml:"influxdb"`
}

// InfluxUnifi is returned by New() after you provide a Config.
type InfluxUnifi struct {
	Collector poller.Collect
	influx    influx.Client
	LastCheck time.Time
	*InfluxDB
}

type metric struct {
	Table  string
	Tags   map[string]string
	Fields map[string]interface{}
}

func init() {
	u := &InfluxUnifi{InfluxDB: &InfluxDB{}, LastCheck: time.Now()}
	poller.NewOutput(&poller.Output{
		Name:   "influxdb",
		Config: u.InfluxDB,
		Method: u.Run,
	})
}

// PollController runs forever, polling UniFi and pushing to InfluxDB
// This is started by Run() or RunBoth() after everything checks out.
func (u *InfluxUnifi) PollController() {
	interval := u.Config.Interval.Round(time.Second)
	log.Printf("[INFO] Everything checks out! Poller started, InfluxDB interval: %v", interval)

	ticker := time.NewTicker(interval)
	for u.LastCheck = range ticker.C {
		metrics, err := u.Collector.Metrics()
		if err != nil {
			u.Collector.LogErrorf("%v", err)
			continue
		}

		report, err := u.ReportMetrics(metrics)
		if err != nil {
			// XXX: reset and re-auth? not sure..
			u.Collector.LogErrorf("%v", err)
			continue
		}

		u.LogInfluxReport(report)
	}
}

// Run runs a ticker to poll the unifi server and update influxdb.
func (u *InfluxUnifi) Run(c poller.Collect) error {
	var err error

	if u.Config.Disable {
		return nil
	}

	u.Collector = c
	u.setConfigDefaults()

	u.influx, err = influx.NewHTTPClient(influx.HTTPConfig{
		Addr:      u.Config.URL,
		Username:  u.Config.User,
		Password:  u.Config.Pass,
		TLSConfig: &tls.Config{InsecureSkipVerify: !u.Config.VerifySSL},
	})
	if err != nil {
		return err
	}

	u.PollController()

	return nil
}

func (u *InfluxUnifi) setConfigDefaults() {
	if u.Config.URL == "" {
		u.Config.URL = defaultInfluxURL
	}

	if u.Config.User == "" {
		u.Config.User = defaultInfluxUser
	}

	if u.Config.Pass == "" {
		u.Config.Pass = defaultInfluxUser
	}

	if u.Config.DB == "" {
		u.Config.DB = defaultInfluxDB
	}

	if u.Config.Interval.Duration == 0 {
		u.Config.Interval = conf.Duration{Duration: defaultInterval}
	} else if u.Config.Interval.Duration < defaultInterval/2 {
		u.Config.Interval = conf.Duration{Duration: defaultInterval / 2}
	}

	u.Config.Interval = conf.Duration{Duration: u.Config.Interval.Duration.Round(time.Second)}
}

// ReportMetrics batches all device and client data into influxdb data points.
// Call this after you've collected all the data you care about.
// Returns an error if influxdb calls fail, otherwise returns a report.
func (u *InfluxUnifi) ReportMetrics(m *poller.Metrics) (*Report, error) {
	r := &Report{Metrics: m, ch: make(chan *metric), Start: time.Now()}
	defer close(r.ch)
	// Make a new Influx Points Batcher.
	var err error
	r.bp, err = influx.NewBatchPoints(influx.BatchPointsConfig{Database: u.Config.DB})
	if err != nil {
		return nil, fmt.Errorf("influx.NewBatchPoints: %v", err)
	}

	go u.collect(r, r.ch)
	// Batch all the points.
	u.loopPoints(r)
	r.wg.Wait() // wait for all points to finish batching!

	// Send all the points.
	if err = u.influx.Write(r.bp); err != nil {
		return nil, fmt.Errorf("influxdb.Write(points): %v", err)
	}
	r.Elapsed = time.Since(r.Start)
	return r, nil
}

// collect runs in a go routine and batches all the points.
func (u *InfluxUnifi) collect(r report, ch chan *metric) {
	for m := range ch {
		pt, err := influx.NewPoint(m.Table, m.Tags, m.Fields, r.metrics().TS)
		if err != nil {
			r.error(err)
		} else {
			r.batch(m, pt)
		}
		r.done()
	}
}

// loopPoints kicks off 3 or 7 go routines to process metrics and send them
// to the collect routine through the metric channel.
func (u *InfluxUnifi) loopPoints(r report) {
	m := r.metrics()
	r.add()
	go func() {
		defer r.done()
		for _, s := range m.Sites {
			u.batchSite(r, s)
		}
	}()
	r.add()
	go func() {
		defer r.done()
		for _, s := range m.Clients {
			u.batchClient(r, s)
		}
	}()
	r.add()
	go func() {
		defer r.done()
		for _, s := range m.IDSList {
			u.batchIDS(r, s)
		}
	}()
	if m.Devices == nil {
		return
	}

	r.add()
	go func() {
		defer r.done()
		for _, s := range m.UAPs {
			u.batchUAP(r, s)
		}
	}()
	r.add()
	go func() {
		defer r.done()
		for _, s := range m.USGs {
			u.batchUSG(r, s)
		}
	}()
	r.add()
	go func() {
		defer r.done()
		for _, s := range m.USWs {
			u.batchUSW(r, s)
		}
	}()
	r.add()
	go func() {
		defer r.done()
		for _, s := range m.UDMs {
			u.batchUDM(r, s)
		}
	}()
}

// LogInfluxReport writes a log message after exporting to influxdb.
func (u *InfluxUnifi) LogInfluxReport(r *Report) {
	idsMsg := fmt.Sprintf("IDS Events: %d, ", len(r.Metrics.IDSList))
	u.Collector.Logf("UniFi Metrics Recorded. Sites: %d, Clients: %d, "+
		"UAP: %d, USG/UDM: %d, USW: %d, %sPoints: %d, Fields: %d, Errs: %d, Elapsed: %v",
		len(r.Metrics.Sites), len(r.Metrics.Clients), len(r.Metrics.UAPs),
		len(r.Metrics.UDMs)+len(r.Metrics.USGs), len(r.Metrics.USWs), idsMsg, r.Total,
		r.Fields, len(r.Errors), r.Elapsed.Round(time.Millisecond))
}
