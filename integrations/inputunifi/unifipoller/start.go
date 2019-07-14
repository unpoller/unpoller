package unifipoller

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"code.golift.io/unifi"
	"github.com/BurntSushi/toml"
	influx "github.com/influxdata/influxdb1-client/v2"
	"github.com/pkg/errors"
	flag "github.com/spf13/pflag"
	yaml "gopkg.in/yaml.v2"
)

// Start begins the application from a CLI.
// Parses flags, parses config and executes Run().
func Start() error {
	log.SetFlags(log.LstdFlags)
	up := &UnifiPoller{}
	if up.ParseFlags(os.Args[1:]); up.ShowVer {
		fmt.Printf("unifi-poller v%s\n", Version)
		return nil // don't run anything else w/ version request.
	}
	if err := up.GetConfig(); err != nil {
		up.Flag.Usage()
		return err
	}
	return up.Run()
}

// ParseFlags runs the parser.
func (u *UnifiPoller) ParseFlags(args []string) {
	u.Flag = flag.NewFlagSet("unifi-poller", flag.ExitOnError)
	u.Flag.Usage = func() {
		fmt.Println("Usage: unifi-poller [--config=filepath] [--version]")
		u.Flag.PrintDefaults()
	}
	u.Flag.StringVarP(&u.DumpJSON, "dumpjson", "j", "",
		"This debug option prints a json payload and exits. See man page for more.")
	u.Flag.StringVarP(&u.ConfigFile, "config", "c", DefaultConfFile, "Poller Config File (TOML Format)")
	u.Flag.BoolVarP(&u.ShowVer, "version", "v", false, "Print the version and exit")
	_ = u.Flag.Parse(args)
}

// GetConfig parses and returns our configuration data.
func (u *UnifiPoller) GetConfig() error {
	// Preload our defaults.
	u.Config = &Config{
		InfluxURL:  defaultInfxURL,
		InfluxUser: defaultInfxUser,
		InfluxPass: defaultInfxPass,
		InfluxDB:   defaultInfxDb,
		UnifiUser:  defaultUnifUser,
		UnifiPass:  os.Getenv("UNIFI_PASSWORD"),
		UnifiBase:  defaultUnifURL,
		Interval:   Duration{defaultInterval},
		Sites:      []string{"default"},
		Quiet:      u.DumpJSON != "",
	}
	u.Logf("Loading Configuration File: %s", u.ConfigFile)
	switch buf, err := ioutil.ReadFile(u.ConfigFile); {
	case err != nil:
		return err
	case strings.Contains(u.ConfigFile, ".json"):
		return json.Unmarshal(buf, u.Config)
	case strings.Contains(u.ConfigFile, ".xml"):
		return xml.Unmarshal(buf, u.Config)
	case strings.Contains(u.ConfigFile, ".yaml"):
		return yaml.Unmarshal(buf, u.Config)
	default:
		return toml.Unmarshal(buf, u.Config)
	}
}

// Run invokes all the application logic and routines.
func (u *UnifiPoller) Run() (err error) {
	if u.DumpJSON != "" {
		return u.DumpJSONPayload()
	}
	if u.Debug {
		log.SetFlags(log.Lshortfile | log.Lmicroseconds | log.Ldate)
		u.LogDebugf("Debug Logging Enabled")
	}
	log.Printf("[INFO] UniFi Poller v%v Starting Up! PID: %d", Version, os.Getpid())
	if err = u.GetUnifi(); err != nil {
		return err
	}
	if err = u.GetInfluxDB(); err != nil {
		return err
	}
	switch strings.ToLower(u.Mode) {
	case "influxlambda", "lambdainflux", "lambda_influx", "influx_lambda":
		u.LogDebugf("Lambda Mode Enabled")
		u.LastCheck = time.Now()
		return u.CollectAndReport()
	default:
		return u.PollController()
	}
}

// GetInfluxDB returns an InfluxDB interface.
func (u *UnifiPoller) GetInfluxDB() (err error) {
	u.Client, err = influx.NewHTTPClient(influx.HTTPConfig{
		Addr:     u.InfluxURL,
		Username: u.InfluxUser,
		Password: u.InfluxPass,
	})
	if err != nil {
		return errors.Wrap(err, "influxdb")
	}
	u.Logf("Logging Measurements to InfluxDB at %s as user %s", u.InfluxURL, u.InfluxUser)
	return nil
}

// GetUnifi returns a UniFi controller interface.
func (u *UnifiPoller) GetUnifi() (err error) {
	// Create an authenticated session to the Unifi Controller.
	u.Unifi, err = unifi.NewUnifi(u.UnifiUser, u.UnifiPass, u.UnifiBase, u.VerifySSL)
	if err != nil {
		return errors.Wrap(err, "unifi controller")
	}
	u.Unifi.ErrorLog = u.LogErrorf // Log all errors.
	u.Unifi.DebugLog = u.LogDebugf // Log debug messages.
	u.Logf("Authenticated to UniFi Controller at %s version %s as user %s", u.UnifiBase, u.ServerVersion, u.UnifiUser)
	if err := u.CheckSites(); err != nil {
		return err
	}
	u.Logf("Polling UniFi Controller Sites: %v", u.Sites)
	return nil
}
