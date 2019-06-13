package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/golift/unifi"
	influx "github.com/influxdata/influxdb1-client/v2"
	"github.com/naoina/toml"
	"github.com/pkg/errors"
	flag "github.com/spf13/pflag"
)

func main() {
	u := &UnifiPoller{}
	if u.ParseFlags(os.Args[1:]); u.ShowVer {
		fmt.Printf("unifi-poller v%s\n", Version)
		return // don't run anything else.
	}
	if err := u.GetConfig(); err != nil {
		u.Flag.Usage()
		log.Fatalf("[ERROR] config file '%v': %v", u.ConfigFile, err)
	}
	if err := u.Run(); err != nil {
		log.Fatalln("[ERROR]", err)
	}
}

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
		Interval:   Dur{value: defaultInterval},
		Sites:      []string{"default"},
	}
	if buf, err := ioutil.ReadFile(u.ConfigFile); err != nil {
		return err
		// This is where the defaults in the config variable are overwritten.
	} else if err := toml.Unmarshal(buf, u.Config); err != nil {
		return err
	}
	if u.DumpJSON != "" {
		u.Quiet = true
	}
	u.Config.Logf("Loaded Configuration: %s", u.ConfigFile)
	return nil
}

// Run invokes all the application logic and routines.
func (u *UnifiPoller) Run() (err error) {
	if u.DumpJSON != "" {
		return u.DumpJSONPayload()
	}
	if log.SetFlags(0); u.Debug {
		log.SetFlags(log.Lshortfile | log.Lmicroseconds | log.Ldate)
		log.Println("[DEBUG] Debug Logging Enabled")
	}
	log.Printf("[INFO] Unifi-Poller v%v Starting Up! PID: %d", Version, os.Getpid())

	if err = u.GetUnifi(); err != nil {
		return err
	}
	if err = u.GetInfluxDB(); err != nil {
		return err
	}
	u.PollController()
	return nil
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
	u.Unifi.ErrorLog = log.Printf // Log all errors.
	// Doing it this way allows debug error logs (line numbers, etc)
	if u.Debug && !u.Quiet {
		u.Unifi.DebugLog = log.Printf // Log debug messages.
	}
	u.Logf("Authenticated to Unifi Controller at %s as user %s", u.UnifiBase, u.UnifiUser)
	if err = u.CheckSites(); err != nil {
		return err
	}
	u.Logf("Polling Unifi Controller Sites: %v", u.Sites)
	return nil
}
