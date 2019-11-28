// Package poller provides the CLI interface to setup unifi-poller.
package poller

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/davidnewhall/unifi-poller/promunifi"
	influx "github.com/influxdata/influxdb1-client/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/pflag"
	"golift.io/unifi"
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
			UnifiUser:  defaultUnifiUser,
			UnifiPass:  defaultUnifiUser,
			UnifiBase:  defaultUnifiURL,
			Interval:   Duration{defaultInterval},
			Sites:      []string{"all"},
			HTTPListen: defaultHTTPListen,
			Namespace:  appName,
		}, Flag: &Flag{},
	}
}

// Start begins the application from a CLI.
// Parses cli flags, parses config file, parses env vars, sets up logging, then:
// - dumps a json payload OR - authenticates unifi controller and executes Run().
func (u *UnifiPoller) Start() error {
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
	if err := u.Config.ParseFile(u.Flag.ConfigFile); err != nil {
		u.Flag.Usage()
		return err
	}

	// Update Config with ENV variable overrides.
	if err := u.Config.ParseENV(); err != nil {
		return err
	}

	if u.Flag.DumpJSON != "" {
		return u.DumpJSONPayload()
	}

	if u.Config.Debug {
		log.SetFlags(log.Lshortfile | log.Lmicroseconds | log.Ldate)
		u.LogDebugf("Debug Logging Enabled")
	}

	log.Printf("[INFO] UniFi Poller v%v Starting Up! PID: %d", Version, os.Getpid())
	if err := u.GetUnifi(); err != nil {
		return err
	}

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
	u.Logf("Polling UniFi Controller at %s v%s as user %s. Sites: %v",
		u.Config.UnifiBase, u.Unifi.ServerVersion, u.Config.UnifiUser, u.Config.Sites)

	switch strings.ToLower(u.Config.Mode) {
	default:
		if err := u.GetInfluxDB(); err != nil {
			return err
		}
		u.Logf("Logging Measurements to InfluxDB at %s as user %s", u.Config.InfluxURL, u.Config.InfluxUser)
		u.Config.Mode = "influx poller"
		return u.PollController()

	case "influxlambda", "lambdainflux", "lambda_influx", "influx_lambda":
		if err := u.GetInfluxDB(); err != nil {
			return err
		}
		u.Logf("Logging Measurements to InfluxDB at %s as user %s one time (lambda mode)",
			u.Config.InfluxURL, u.Config.InfluxUser)
		u.LastCheck = time.Now()
		return u.CollectAndProcess()

	case "prometheus", "exporter":
		u.Logf("Exporting Measurements at https://%s/metrics for Prometheus", u.Config.HTTPListen)
		http.Handle("/metrics", promhttp.Handler())
		prometheus.MustRegister(promunifi.NewUnifiCollector(promunifi.UnifiCollectorCnfg{
			Namespace:    strings.Replace(u.Config.Namespace, "-", "", -1),
			CollectFn:    u.ExportMetrics,
			LoggingFn:    u.LogExportReport,
			ReportErrors: true, // XXX: Does this need to be configurable?
		}))
		return http.ListenAndServe(u.Config.HTTPListen, nil)
	}
}

// GetInfluxDB returns an InfluxDB interface.
func (u *UnifiPoller) GetInfluxDB() (err error) {
	u.Influx, err = influx.NewHTTPClient(influx.HTTPConfig{
		Addr:      u.Config.InfluxURL,
		Username:  u.Config.InfluxUser,
		Password:  u.Config.InfluxPass,
		TLSConfig: &tls.Config{InsecureSkipVerify: u.Config.InfxBadSSL},
	})
	if err != nil {
		return fmt.Errorf("influxdb: %v", err)
	}

	return nil
}

// GetUnifi returns a UniFi controller interface.
func (u *UnifiPoller) GetUnifi() (err error) {
	// Create an authenticated session to the Unifi Controller.
	u.Unifi, err = unifi.NewUnifi(&unifi.Config{
		User:      u.Config.UnifiUser,
		Pass:      u.Config.UnifiPass,
		URL:       u.Config.UnifiBase,
		VerifySSL: u.Config.VerifySSL,
		ErrorLog:  u.LogErrorf, // Log all errors.
		DebugLog:  u.LogDebugf, // Log debug messages.
	})
	if err != nil {
		return fmt.Errorf("unifi controller: %v", err)
	}
	u.LogDebugf("Authenticated with controller successfully")

	return u.CheckSites()
}
