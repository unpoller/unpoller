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
func (u *UnifiPoller) Run() error {
	c := u.Config
	if u.DumpJSON != "" {
		return c.DumpJSON(u.DumpJSON)
	}
	if log.SetFlags(0); c.Debug {
		log.SetFlags(log.Lshortfile | log.Lmicroseconds | log.Ldate)
		log.Println("[DEBUG] Debug Logging Enabled")
	}
	log.Printf("[INFO] Unifi-Poller v%v Starting Up! PID: %d", Version, os.Getpid())
	controller, err := c.GetController()
	if err != nil {
		return err
	}
	infdb, err := c.GetInfluxDB()
	if err != nil {
		return err
	}
	c.PollUnifiController(controller, infdb)
	return nil
}

func (c *Config) GetInfluxDB() (influx.Client, error) {
	infdb, err := influx.NewHTTPClient(influx.HTTPConfig{
		Addr:     c.InfluxURL,
		Username: c.InfluxUser,
		Password: c.InfluxPass,
	})
	if err != nil {
		return nil, errors.Wrap(err, "influxdb")
	}
	c.Logf("Logging Measurements to InfluxDB at %s as user %s", c.InfluxURL, c.InfluxUser)
	return infdb, nil
}

func (c *Config) GetController() (*unifi.Unifi, error) {
	// Create an authenticated session to the Unifi Controller.
	controller, err := unifi.NewUnifi(c.UnifiUser, c.UnifiPass, c.UnifiBase, c.VerifySSL)
	if err != nil {
		return nil, errors.Wrap(err, "unifi controller")
	}
	controller.ErrorLog = log.Printf // Log all errors.
	// Doing it this way allows debug error logs (line numbers, etc)
	if c.Debug && !c.Quiet {
		controller.DebugLog = log.Printf // Log debug messages.
	}
	c.Logf("Authenticated to Unifi Controller at %s as user %s", c.UnifiBase, c.UnifiUser)
	if err := c.CheckSites(controller); err != nil {
		return nil, err
	}
	c.Logf("Polling Unifi Controller Sites: %v", c.Sites)
	return controller, nil
}
