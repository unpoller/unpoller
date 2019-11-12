package unifipoller

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	influx "github.com/influxdata/influxdb1-client/v2"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/pflag"
	"golift.io/unifi"
)

// Start begins the application from a CLI.
// Parses flags, parses config and executes Run().
func Start() error {
	log.SetFlags(log.LstdFlags)
	up := &UnifiPoller{
		Flag: &Flag{},
		Config: &Config{
			// Preload our defaults.
			InfluxURL:  defaultInfluxURL,
			InfluxUser: defaultInfluxUser,
			InfluxPass: defaultInfluxPass,
			InfluxDB:   defaultInfluxDB,
			UnifiUser:  defaultUnifiUser,
			UnifiPass:  os.Getenv("UNIFI_PASSWORD"), // deprecated name.
			UnifiBase:  defaultUnifiURL,
			Interval:   Duration{defaultInterval},
			Sites:      []string{"all"},
			HTTPListen: defaultHTTPListen,
		}}
	up.Flag.Parse(os.Args[1:])
	if up.Flag.ShowVer {
		fmt.Printf("unifi-poller v%s\n", Version)
		return nil // don't run anything else w/ version request.
	}
	if up.Flag.DumpJSON == "" { // do not print this when dumping JSON.
		up.Logf("Loading Configuration File: %s", up.Flag.ConfigFile)
	}
	// Parse config file.
	if err := up.Config.ParseFile(up.Flag.ConfigFile); err != nil {
		up.Flag.Usage()
		return err
	}
	// Update Config with ENV variable overrides.
	if err := up.Config.ParseENV(); err != nil {
		return err
	}
	return up.Run()
}

// Parse turns CLI arguments into data structures. Called by Start() on startup.
func (f *Flag) Parse(args []string) {
	f.FlagSet = pflag.NewFlagSet("unifi-poller", pflag.ExitOnError)
	f.Usage = func() {
		fmt.Println("Usage: unifi-poller [--config=/path/to/up.conf] [--version]")
		f.PrintDefaults()
	}
	f.StringVarP(&f.DumpJSON, "dumpjson", "j", "",
		"This debug option prints a json payload and exits. See man page for more info.")
	f.StringVarP(&f.ConfigFile, "config", "c", DefaultConfFile, "Poller config file path.")
	f.BoolVarP(&f.ShowVer, "version", "v", false, "Print the version and exit.")
	_ = f.FlagSet.Parse(args) // pflag.ExitOnError means this will never return error.
}

// Run invokes all the application logic and routines.
func (u *UnifiPoller) Run() (err error) {
	if u.Flag.DumpJSON != "" {
		return u.DumpJSONPayload()
	}
	if u.Config.Debug {
		log.SetFlags(log.Lshortfile | log.Lmicroseconds | log.Ldate)
		u.LogDebugf("Debug Logging Enabled")
	}
	log.Printf("[INFO] UniFi Poller v%v Starting Up! PID: %d", Version, os.Getpid())
	if err = u.GetUnifi(); err != nil {
		return err
	}
	u.Logf("Polling UniFi Controller at %s v%s as user %s. Sites: %v",
		u.Config.UnifiBase, u.Unifi.ServerVersion, u.Config.UnifiUser, u.Config.Sites)
	if err = u.GetInfluxDB(); err != nil {
		return err
	}

	switch strings.ToLower(u.Config.Mode) {
	case "influxlambda", "lambdainflux", "lambda_influx", "influx_lambda":
		u.Logf("Logging Measurements to InfluxDB at %s as user %s one time (lambda mode)",
			u.Config.InfluxURL, u.Config.InfluxUser)
		u.LastCheck = time.Now()
		return u.CollectAndProcess(u.ReportMetrics)
	case "prometheus", "exporter":
		u.Logf("Exporting Measurements at https://%s/metrics for Prometheus", u.Config.HTTPListen)
		u.Config.Mode = "http exporter"
		http.Handle("/metrics", promhttp.Handler())
		go func() {
			err = http.ListenAndServe(u.Config.HTTPListen, nil)
			if err != http.ErrServerClosed {
				log.Fatalf("[ERROR] http server: %v", err)
			}
		}()
		return u.PollController(u.ExportMetrics)
	default:
		u.Logf("Logging Measurements to InfluxDB at %s as user %s", u.Config.InfluxURL, u.Config.InfluxUser)
		u.Config.Mode = "influx poller"
		return u.PollController(u.ReportMetrics)
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
