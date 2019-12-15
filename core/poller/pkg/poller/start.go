// Package poller provides the CLI interface to setup unifi-poller.
package poller

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/spf13/pflag"
	"golift.io/config"
)

// New returns a new poller struct preloaded with default values.
// No need to call this if you call Start.c
func New() *UnifiPoller {
	return &UnifiPoller{
		Config: &Config{
			InfluxURL:  defaultInfluxURL,
			InfluxUser: defaultInfluxUser,
			InfluxPass: defaultInfluxPass,
			InfluxDB:   defaultInfluxDB,
			Interval:   config.Duration{Duration: defaultInterval},
			HTTPListen: defaultHTTPListen,
			Namespace:  appName,
		},
		Flag: &Flag{
			ConfigFile: DefaultConfFile,
		},
	}
}

// Start begins the application from a CLI.
// Parses cli flags, parses config file, parses env vars, sets up logging, then:
// - dumps a json payload OR - executes Run().
func (u *UnifiPoller) Start() error {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags)
	u.Flag.Parse(os.Args[1:])

	if u.Flag.ShowVer {
		fmt.Printf("%s v%s\n", appName, Version)
		return nil // don't run anything else w/ version request.
	}

	if u.Flag.DumpJSON == "" { // do not print this when dumping JSON.
		u.Logf("Loading Configuration File: %s", u.Flag.ConfigFile)
	}

	// Parse config file.
	if err := config.ParseFile(u.Config, u.Flag.ConfigFile); err != nil {
		u.Flag.Usage()
		return err
	}

	// Update Config with ENV variable overrides.
	if _, err := config.ParseENV(u.Config, ENVConfigPrefix); err != nil {
		return err
	}

	if len(u.Config.Controllers) < 1 {
		u.Config.Controllers = []Controller{{
			Sites:     []string{"all"},
			User:      defaultUnifiUser,
			Pass:      "",
			URL:       defaultUnifiURL,
			SaveSites: true,
		}}
	}

	if u.Flag.DumpJSON != "" {
		return u.DumpJSONPayload()
	}

	if u.Config.Debug {
		log.SetFlags(log.Lshortfile | log.Lmicroseconds | log.Ldate)
		u.LogDebugf("Debug Logging Enabled")
	}

	log.Printf("[INFO] UniFi Poller v%v Starting Up! PID: %d", Version, os.Getpid())

	return u.Run()
}

// Parse turns CLI arguments into data structures. Called by Start() on startup.
func (f *Flag) Parse(args []string) {
	f.FlagSet = pflag.NewFlagSet(appName, pflag.ExitOnError)
	f.Usage = func() {
		fmt.Printf("Usage: %s [--config=/path/to/up.conf] [--version]", appName)
		f.PrintDefaults()
	}

	f.StringVarP(&f.DumpJSON, "dumpjson", "j", "",
		"This debug option prints a json payload and exits. See man page for more info.")
	f.StringVarP(&f.ConfigFile, "config", "c", DefaultConfFile, "Poller config file path.")
	f.BoolVarP(&f.ShowVer, "version", "v", false, "Print the version and exit.")
	_ = f.FlagSet.Parse(args) // pflag.ExitOnError means this will never return error.
}

// Run picks a mode and executes the associated functions. This will do one of three things:
// 1. Start the collector routine that polls unifi and reports to influx on an interval. (default)
// 2. Run the collector one time and report the metrics to influxdb. (lambda)
// 3. Start a web server and wait for Prometheus to poll the application for metrics.
func (u *UnifiPoller) Run() error {
	for i, c := range u.Config.Controllers {
		if c.Name == "" {
			u.Config.Controllers[i].Name = c.URL
		}

		switch err := u.GetUnifi(c); err {
		case nil:
			u.Logf("Polling UniFi Controller at %s v%s as user %s. Sites: %v",
				c.URL, c.Unifi.ServerVersion, c.User, c.Sites)
		default:
			u.LogErrorf("Controller Auth or Connection failed, but continuing to retry! %s: %v", c.URL, err)
		}
	}

	switch strings.ToLower(u.Config.Mode) {
	default:
		u.PollController()
		return nil
	case "influxlambda", "lambdainflux", "lambda_influx", "influx_lambda":
		u.LastCheck = time.Now()
		return u.CollectAndProcess()
	case "both":
		go u.PollController()
		fallthrough
	case "prometheus", "exporter":
		return u.RunPrometheus()
	}
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
