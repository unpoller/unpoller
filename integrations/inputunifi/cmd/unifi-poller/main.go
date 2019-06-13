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
		os.Exit(0) // don't run anything else.
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
	if !u.Config.Quiet {
		log.Println("Loaded Configuration:", u.ConfigFile)
	}
	return nil
}

// Run invokes all the application logic and routines.
func (u *UnifiPoller) Run() error {
	c := u.Config
	if u.DumpJSON != "" {
		if err := c.DumpJSON(u.DumpJSON); err != nil {
			log.Fatalln("[ERROR] dumping JSON:", err)
		}
		return nil
	}
	log.Println("Unifi-Poller Starting Up! PID:", os.Getpid())
	// Create an authenticated session to the Unifi Controller.
	controller, err := unifi.NewUnifi(c.UnifiUser, c.UnifiPass, c.UnifiBase, c.VerifySSL)
	if err != nil {
		return errors.Wrap(err, "unifi controller")
	}
	if !c.Quiet {
		log.Println("Authenticated to Unifi Controller @", c.UnifiBase, "as user", c.UnifiUser)
	}
	if err := c.CheckSites(controller); err != nil {
		return err
	}
	controller.ErrorLog = log.Printf // Log all errors.
	if log.SetFlags(0); c.Debug {
		log.Println("Debug Logging Enabled")
		log.SetFlags(log.Lshortfile | log.Lmicroseconds | log.Ldate)
		controller.DebugLog = log.Printf // Log debug messages.
	}
	infdb, err := influx.NewHTTPClient(influx.HTTPConfig{
		Addr:     c.InfluxURL,
		Username: c.InfluxUser,
		Password: c.InfluxPass,
	})
	if err != nil {
		return errors.Wrap(err, "influxdb")
	}
	if c.Quiet {
		// Doing it this way allows debug error logs (line numbers, etc)
		controller.DebugLog = nil
	} else {
		log.Println("Logging Unifi Metrics to InfluXDB @", c.InfluxURL, "as user", c.InfluxUser)
		log.Printf("Polling Unifi Controller (sites %v), interval: %v", c.Sites, c.Interval.value)
	}
	c.PollUnifiController(controller, infdb)
	return nil
}
