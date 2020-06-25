package lokiunifi

import (
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/unifi-poller/poller"
	"golift.io/cnfg"
)

const (
	maxInterval     = time.Hour
	minInterval     = 10 * time.Second
	defaultTimeout  = 10 * time.Second
	defaultInterval = 2 * time.Minute
)

const (
	// InputName is the name of plugin that gives us data.
	InputName = "unifi"
	// PluginName is the name of this plugin.
	PluginName = "loki"
)

// Config is the plugin's input data.
type Config struct {
	Disable   bool          `json:"disable" toml:"disable" xml:"disable" yaml:"disable"`
	VerifySSL bool          `json:"verify_ssl" toml:"verify_ssl" xml:"verify_ssl" yaml:"verify_ssl"`
	URL       string        `json:"url" toml:"url" xml:"url" yaml:"url"`
	Username  string        `json:"user" toml:"user" xml:"user" yaml:"user"`
	Password  string        `json:"pass" toml:"pass" xml:"pass" yaml:"pass"`
	TenantID  string        `json:"tenant_id" toml:"tenant_id" xml:"tenant_id" yaml:"tenant_id"`
	Interval  cnfg.Duration `json:"interval" toml:"interval" xml:"interval" yaml:"interval"`
	Timeout   cnfg.Duration `json:"timeout" toml:"timeout" xml:"timeout" yaml:"timeout"`
}

// Loki is the main library struct. This satisfies the poller.Output interface.
type Loki struct {
	poller.Collect
	*Config `json:"loki" toml:"loki" xml:"loki" yaml:"loki"`
	client  *Client
	last    time.Time
}

// init is how this modular code is initialized by the main app.
// This module adds itself as an output module to the poller core.
func init() { // nolint: gochecknoinits
	l := &Loki{Config: &Config{
		Interval: cnfg.Duration{Duration: defaultInterval},
		Timeout:  cnfg.Duration{Duration: defaultTimeout},
	}}

	poller.NewOutput(&poller.Output{
		Name:   PluginName,
		Config: l,
		Method: l.Run,
	})
}

// Run is fired from the poller library after the Config is unmarshalled.
func (l *Loki) Run(collect poller.Collect) error {
	if l.Collect = collect; l.Config == nil || l.URL == "" || l.Disable {
		l.Logf("Loki config missing (or disabled), Loki output disabled!")
		return nil
	}

	l.ValidateConfig()
	l.PollController()
	l.LogErrorf("Loki Output Plugin Stopped!")

	return nil
}

// ValidateConfig sets initial "last" update time. Also creates an http client,
// makes sure URL is sane, and sets interval within min/max limits.
func (l *Loki) ValidateConfig() {
	if l.Interval.Duration > maxInterval {
		l.Interval.Duration = maxInterval
	} else if l.Interval.Duration < minInterval {
		l.Interval.Duration = minInterval
	}

	l.last = time.Now().Add(-l.Interval.Duration)
	l.client = l.httpClient()
	l.URL = strings.TrimRight(l.URL, "/") // gets a path appended to it later.
}

// PollController runs forever, polling UniFi for events and pushing them to Loki.
// This is started by Run().
func (l *Loki) PollController() {
	interval := l.Interval.Round(time.Second)
	l.Logf("Loki Event collection started, interval: %v, URL: %s", interval, l.URL)

	ticker := time.NewTicker(interval)
	for start := range ticker.C {
		if err := l.pollController(start); err != nil {
			l.LogErrorf("%v", err)
		}
	}
}

// pollController offloads the loop from PollController.
func (l *Loki) pollController(start time.Time) error {
	events, err := l.Events(&poller.Filter{Name: InputName})
	if err != nil {
		return errors.Wrap(err, "event fetch for Loki failed")
	}

	report := &Report{
		Events: events,
		Start:  start,
		Logger: l.Collect,
		Client: l.client,
		Last:   &l.last,
	}

	return report.Execute(4 * l.Interval.Duration) // nolint: gomnd
}
