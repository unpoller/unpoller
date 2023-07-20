package lokiunifi

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/unpoller/unpoller/pkg/poller"
	"github.com/unpoller/unpoller/pkg/webserver"
	"golift.io/cnfg"
)

const (
	maxInterval     = 10 * time.Minute
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
	Disable   bool          `json:"disable"    toml:"disable"    xml:"disable"    yaml:"disable"`
	VerifySSL bool          `json:"verify_ssl" toml:"verify_ssl" xml:"verify_ssl" yaml:"verify_ssl"`
	URL       string        `json:"url"        toml:"url"        xml:"url"        yaml:"url"`
	Username  string        `json:"user"       toml:"user"       xml:"user"       yaml:"user"`
	Password  string        `json:"pass"       toml:"pass"       xml:"pass"       yaml:"pass"`
	TenantID  string        `json:"tenant_id"  toml:"tenant_id"  xml:"tenant_id"  yaml:"tenant_id"`
	Interval  cnfg.Duration `json:"interval"   toml:"interval"   xml:"interval"   yaml:"interval"`
	Timeout   cnfg.Duration `json:"timeout"    toml:"timeout"    xml:"timeout"    yaml:"timeout"`
}

// Loki is the main library struct. This satisfies the poller.Output interface.
type Loki struct {
	Collect poller.Collect
	*Config `json:"loki" toml:"loki" xml:"loki" yaml:"loki"`
	client  *Client
	last    time.Time
}

var _ poller.OutputPlugin = &Loki{}

// init is how this modular code is initialized by the main app.
// This module adds itself as an output module to the poller core.
func init() { // nolint: gochecknoinits
	l := &Loki{Config: &Config{
		Interval: cnfg.Duration{Duration: defaultInterval},
		Timeout:  cnfg.Duration{Duration: defaultTimeout},
	}}

	poller.NewOutput(&poller.Output{
		Name:         PluginName,
		Config:       l,
		OutputPlugin: l,
	})
}

func (l *Loki) Enabled() bool {
	if l == nil {
		return false
	}

	if l.Config == nil {
		return false
	}

	if l.URL == "" {
		return false
	}

	return !l.Disable
}

func (l *Loki) DebugOutput() (bool, error) {
	if l == nil {
		return true, nil
	}

	if !l.Enabled() {
		return true, nil
	}

	if err := l.ValidateConfig(); err != nil {
		return false, err
	}

	return true, nil
}

// Run is fired from the poller library after the Config is unmarshalled.
func (l *Loki) Run(collect poller.Collect) error {
	l.Collect = collect
	if !l.Enabled() {
		l.LogDebugf("Loki config missing (or disabled), Loki output disabled!")

		return nil
	}

	l.Logf("Loki enabled")

	if err := l.ValidateConfig(); err != nil {
		l.LogErrorf("invalid loki config")

		return err
	}

	fake := *l.Config
	fake.Password = strconv.FormatBool(fake.Password != "")

	webserver.UpdateOutput(&webserver.Output{Name: PluginName, Config: fake})
	l.PollController()
	l.LogErrorf("Loki Output Plugin Stopped!")

	return nil
}

// ValidateConfig sets initial "last" update time. Also creates an http client,
// makes sure URL is sane, and sets interval within min/max limits.
func (l *Loki) ValidateConfig() error {
	if l.Interval.Duration > maxInterval {
		l.Interval.Duration = maxInterval
	} else if l.Interval.Duration < minInterval {
		l.Interval.Duration = minInterval
	}

	if strings.HasPrefix(l.Password, "file://") {
		pass, err := os.ReadFile(strings.TrimPrefix(l.Password, "file://"))
		if err != nil {
			l.LogErrorf("Reading Loki Password File: %v", err)
			
			return fmt.Errorf("error reading password file")
		}

		l.Password = strings.TrimSpace(string(pass))
	}

	l.last = time.Now().Add(-l.Interval.Duration)
	l.client = l.httpClient()
	l.URL = strings.TrimRight(l.URL, "/") // gets a path appended to it later.

	return nil
}

// PollController runs forever, polling UniFi for events and pushing them to Loki.
// This is started by Run().
func (l *Loki) PollController() {
	interval := l.Interval.Round(time.Second)

	l.Logf("Loki Event collection started, interval: %v, URL: %s", interval, l.URL)

	ticker := time.NewTicker(interval)
	for start := range ticker.C {
		events, err := l.Collect.Events(&poller.Filter{Name: InputName})
		if err != nil {
			l.LogErrorf("event fetch for Loki failed: %v", err)

			continue
		}

		err = l.ProcessEvents(l.NewReport(start), events)
		if err != nil {
			l.LogErrorf("%v", err)
		}
	}
}

// ProcessEvents offloads some of the loop from PollController.
func (l *Loki) ProcessEvents(report *Report, events *poller.Events) error {
	// Sometimes it gets stuck on old messages. This gets it past that.
	if time.Since(l.last) > 4*l.Interval.Duration {
		l.last = time.Now().Add(-4 * l.Interval.Duration)
	}

	logs := report.ProcessEventLogs(events)
	if err := l.client.Post(logs); err != nil {
		return fmt.Errorf("sending to Loki failed: %w", err)
	}

	l.last = report.Start

	l.Logf("Events sent to Loki. %v", report)

	return nil
}
