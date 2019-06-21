package unifipoller

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/golift/unifi"
	influx "github.com/influxdata/influxdb1-client/v2"
	"github.com/pkg/errors"
	flag "github.com/spf13/pflag"
	yaml "gopkg.in/yaml.v2"
)

// ParseFlags runs the parser.
func (u *UnifiPoller) ParseFlags(args []string) {
	u.Flag = flag.NewFlagSet("unifi-poller", flag.ExitOnError)
	u.Flag.Usage = func() {
		fmt.Println("Usage: unifi-poller [--config=filepath] [--version]")
		u.Flag.PrintDefaults()
	}
	u.Flag.StringVarP(&u.DumpJSON, "dumpjson", "j", "",
		"This debug option prints the json payload for a device and exits.")
	u.Flag.StringVarP(&u.ConfigFile, "config", "c", defaultConfFile, "Poller Config File (TOML Format)")
	u.Flag.BoolVarP(&u.ShowVer, "version", "v", false, "Print the version and exit")
	_ = u.Flag.Parse(args)
}

// GetConfig parses and returns our configuration data.
func (u *UnifiPoller) GetConfig() (err error) {
	// Preload our defaults.
	u.Config = &Config{
		InfluxURL:  defaultInfxURL,
		InfluxUser: defaultInfxUser,
		InfluxPass: defaultInfxPass,
		InfluxDB:   defaultInfxDb,
		UnifiUser:  defaultUnifUser,
		UnifiPass:  os.Getenv("UNIFI_PASSWORD"),
		UnifiBase:  defaultUnifURL,
		Interval:   Dur{defaultInterval},
		Sites:      []string{"default"},
	}
	var buf []byte
	switch buf, err = ioutil.ReadFile(u.ConfigFile); {
	case err != nil:
		return err
	case strings.HasSuffix(u.ConfigFile, ".json"):
		err = json.Unmarshal(buf, u.Config)
	case strings.HasSuffix(u.ConfigFile, ".xml"):
		err = xml.Unmarshal(buf, u.Config)
	case strings.HasSuffix(u.ConfigFile, ".yaml"):
		err = yaml.Unmarshal(buf, u.Config)
	default:
		err = toml.Unmarshal(buf, u.Config)
	}
	if err != nil {
		return err
	}

	if u.DumpJSON != "" {
		u.Quiet = true
	}
	u.Logf("Loaded Configuration: %s", u.ConfigFile)
	return nil
}

// Run invokes all the application logic and routines.
func (u *UnifiPoller) Run() (err error) {
	if u.DumpJSON != "" {
		return u.DumpJSONPayload()
	}
	if log.SetFlags(0); u.Debug {
		log.SetFlags(log.Lshortfile | log.Lmicroseconds | log.Ldate)
		u.LogDebugf("Debug Logging Enabled")
	}
	log.Printf("[INFO] Unifi-Poller v%v Starting Up! PID: %d", Version, os.Getpid())

	if err = u.GetUnifi(); err != nil {
		return err
	}
	if err = u.GetInfluxDB(); err != nil {
		return err
	}
	return u.PollController()
}

// GetInfluxDB returns an influxdb interface.
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

// GetUnifi returns a Unifi controller interface.
func (u *UnifiPoller) GetUnifi() (err error) {
	// Create an authenticated session to the Unifi Controller.
	u.Unifi, err = unifi.NewUnifi(u.UnifiUser, u.UnifiPass, u.UnifiBase, u.VerifySSL)
	if err != nil {
		return errors.Wrap(err, "unifi controller")
	}
	u.Unifi.ErrorLog = u.LogErrorf // Log all errors.
	u.Unifi.DebugLog = u.LogDebugf // Log debug messages.
	v, err := u.GetServer()
	if err != nil {
		v.ServerVersion = "unknown"
	}
	u.Logf("Authenticated to Unifi Controller at %s version %s as user %s", u.UnifiBase, v.ServerVersion, u.UnifiUser)
	if err = u.CheckSites(); err != nil {
		return err
	}
	u.Logf("Polling Unifi Controller Sites: %v", u.Sites)
	return nil
}
