package unifipoller

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	influx "github.com/influxdata/influxdb1-client/v2"
	"github.com/spf13/pflag"
	"golift.io/unifi"
	"gopkg.in/yaml.v2"
)

// Start begins the application from a CLI.
// Parses flags, parses config and executes Run().
func Start() error {
	log.SetFlags(log.LstdFlags)
	up := &UnifiPoller{Flag: &Flag{}}
	up.Flag.Parse(os.Args[1:])
	if up.Flag.ShowVer {
		fmt.Printf("unifi-poller v%s\n", Version)
		return nil // don't run anything else w/ version request.
	}
	// Parse config file.
	if err := up.GetConfig(); err != nil {
		up.Flag.Usage()
		return err
	}
	// Update Config with ENV variable overrides.
	if err := up.ENVSetConfig(); err != nil {
		return err
	}
	return up.Run()
}

// Parse turns CLI arguments into data structures. Called by Start() on startup.
func (f *Flag) Parse(args []string) {
	f.FlagSet = pflag.NewFlagSet("unifi-poller", pflag.ExitOnError)
	f.Usage = func() {
		fmt.Println("Usage: unifi-poller [--config=filepath] [--version]")
		f.PrintDefaults()
	}
	f.StringVarP(&f.DumpJSON, "dumpjson", "j", "",
		"This debug option prints a json payload and exits. See man page for more info.")
	f.StringVarP(&f.ConfigFile, "config", "c", DefaultConfFile, "Poller config file path.")
	f.BoolVarP(&f.ShowVer, "version", "v", false, "Print the version and exit.")
	_ = f.FlagSet.Parse(args)
}

// ENVSetConfig copies environment variables into configuration values.
// This is useful for Docker users that find it easier to pass ENV variables
// than a specific configuration file. Uses reflection to find struct tags.
func (u *UnifiPoller) ENVSetConfig() error {
	t := reflect.TypeOf(Config{})
	// Loop each Config struct member; get reflect tag & env var value; update config.
	for i := 0; i < t.NumField(); i++ {
		// Get the ENV variable name from "env" struct tag then pull value from OS.
		tag := t.Field(i).Tag.Get("env")
		env := os.Getenv(ENVConfigPrefix + tag)
		if tag == "" || env == "" {
			continue
		}

		// Reflect and update the u.Config struct member at position i.
		switch c := reflect.ValueOf(u.Config).Elem().Field(i); c.Type().String() {
		// Handle each member type appropriately (differently).
		case "string":
			// This is a reflect package method to update a struct member by index.
			c.SetString(env)
		case "int":
			val, err := strconv.Atoi(env)
			if err != nil {
				return fmt.Errorf("%s: %v", tag, err)
			}
			c.Set(reflect.ValueOf(val))
		case "[]string":
			c.Set(reflect.ValueOf(strings.Split(env, ",")))
		case "unifipoller.Duration":
			val, err := time.ParseDuration(env)
			if err != nil {
				return fmt.Errorf("%s: %v", tag, err)
			}
			c.Set(reflect.ValueOf(Duration{val}))
		case "bool":
			val, err := strconv.ParseBool(env)
			if err != nil {
				return fmt.Errorf("%s: %v", tag, err)
			}
			c.SetBool(val)
		}
	}
	return nil
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
		UnifiPass:  os.Getenv("UNIFI_PASSWORD"), // deprecated name.
		UnifiBase:  defaultUnifURL,
		Interval:   Duration{defaultInterval},
		Sites:      []string{"default"},
		Quiet:      u.Flag.DumpJSON != "", //s uppress the following u.Logf line.
	}
	u.Logf("Loading Configuration File: %s", u.Flag.ConfigFile)
	switch buf, err := ioutil.ReadFile(u.Flag.ConfigFile); {
	case err != nil:
		return err
	case strings.Contains(u.Flag.ConfigFile, ".json"):
		return json.Unmarshal(buf, u.Config)
	case strings.Contains(u.Flag.ConfigFile, ".xml"):
		return xml.Unmarshal(buf, u.Config)
	case strings.Contains(u.Flag.ConfigFile, ".yaml"):
		return yaml.Unmarshal(buf, u.Config)
	default:
		return toml.Unmarshal(buf, u.Config)
	}
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
	u.Logf("Logging Measurements to InfluxDB at %s as user %s", u.Config.InfluxURL, u.Config.InfluxUser)
	switch strings.ToLower(u.Config.Mode) {
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
	u.Influx, err = influx.NewHTTPClient(influx.HTTPConfig{
		Addr:     u.Config.InfluxURL,
		Username: u.Config.InfluxUser,
		Password: u.Config.InfluxPass,
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
	if err := u.CheckSites(); err != nil {
		return err
	}
	return nil
}
