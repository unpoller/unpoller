package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/golift/unifi"
	influx "github.com/influxdata/influxdb1-client/v2"
	"github.com/naoina/toml"
	"github.com/pkg/errors"
	flag "github.com/spf13/pflag"
)

// Asset is used to give all devices and clients a common interface.
type Asset interface {
	Points() ([]*influx.Point, error)
}

func main() {
	configFile := parseFlags()
	log.Println("Unifi-Poller Starting Up! PID:", os.Getpid())
	config, err := GetConfig(configFile)
	if err != nil {
		flag.Usage()
		log.Fatalf("[ERROR] config file '%v': %v", configFile, err)
	}
	if err := config.Run(); err != nil {
		log.Fatalln("[ERROR]", err)
	}
}

// Run invokes all the application logic and routines.
func (c *Config) Run() error {
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

func parseFlags() string {
	flag.Usage = func() {
		fmt.Println("Usage: unifi-poller [--config=filepath] [--version]")
		flag.PrintDefaults()
	}
	configFile := flag.StringP("config", "c", defaultConfFile, "Poller Config File (TOML Format)")
	version := flag.BoolP("version", "v", false, "Print the version and exit")
	if flag.Parse(); *version {
		fmt.Println("unifi-poller version:", Version)
		os.Exit(0) // don't run anything else.
	}
	return *configFile
}

// CheckSites makes sure the list of provided sites exists on the controller.
func (c *Config) CheckSites(controller *unifi.Unifi) error {
	sites, err := controller.GetSites()
	if err != nil {
		return err
	}
	if !c.Quiet {
		msg := make([]string, 0)
		for _, site := range sites {
			msg = append(msg, site.Name+" ("+site.Desc+")")
		}
		log.Printf("Found %d site(s) on controller: %v", len(msg), strings.Join(msg, ", "))
	}
	if StringInSlice("all", c.Sites) {
		return nil
	}
FIRST:
	for _, s := range c.Sites {
		for _, site := range sites {
			if s == site.Name {
				continue FIRST
			}
		}
		return errors.Errorf("configured site not found on controller: %v", s)
	}
	return nil
}

// GetConfig parses and returns our configuration data.
func GetConfig(configFile string) (Config, error) {
	// Preload our defaults.
	config := Config{
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
	if buf, err := ioutil.ReadFile(configFile); err != nil {
		return config, err
		// This is where the defaults in the config variable are overwritten.
	} else if err := toml.Unmarshal(buf, &config); err != nil {
		return config, err
	}
	log.Println("Loaded Configuration:", configFile)
	return config, nil
}

// PollUnifiController runs forever, polling and pushing.
func (c *Config) PollUnifiController(controller *unifi.Unifi, infdb influx.Client) {
	log.Println("[INFO] Everything checks out! Beginning Poller Routine.")
	ticker := time.NewTicker(c.Interval.value)

	for range ticker.C {
		sites, err := filterSites(controller, c.Sites)
		if err != nil {
			logErrors([]error{err}, "uni.GetSites()")
		}
		// Get all the points.
		clients, err := controller.GetClients(sites)
		if err != nil {
			logErrors([]error{err}, "uni.GetClients()")
		}
		devices, err := controller.GetDevices(sites)
		if err != nil {
			logErrors([]error{err}, "uni.GetDevices()")
		}
		bp, err := influx.NewBatchPoints(influx.BatchPointsConfig{Database: c.InfluxDB})
		if err != nil {
			logErrors([]error{err}, "influx.NewBatchPoints")
			continue
		}
		// Batch all the points.
		if errs := batchPoints(devices, clients, bp); errs != nil && hasErr(errs) {
			logErrors(errs, "asset.Points()")
		}
		if err := infdb.Write(bp); err != nil {
			logErrors([]error{err}, "infdb.Write(bp)")
		}
		if !c.Quiet {
			log.Printf("[INFO] Logged Unifi States. Sites: %d Clients: %d, Wireless APs: %d, Gateways: %d, Switches: %d",
				len(sites), len(clients.UCLs), len(devices.UAPs), len(devices.USGs), len(devices.USWs))
		}
	}
}

// filterSites returns a list of sites to fetch data for.
// Omits requested but unconfigured sites.
func filterSites(controller *unifi.Unifi, filter []string) ([]unifi.Site, error) {
	sites, err := controller.GetSites()
	if err != nil {
		return nil, err
	} else if len(filter) < 1 || StringInSlice("all", filter) {
		return sites, nil
	}
	var i int
	for _, s := range sites {
		// Only include valid sites in the request filter.
		if StringInSlice(s.Name, filter) {
			sites[i] = s
			i++
		}
	}
	return sites[:i], nil
}

// batchPoints combines all device and client data into influxdb data points.
func batchPoints(devices *unifi.Devices, clients *unifi.Clients, bp influx.BatchPoints) (errs []error) {
	process := func(asset Asset) error {
		if asset == nil {
			return nil
		}
		influxPoints, err := asset.Points()
		if err != nil {
			return err
		}
		bp.AddPoints(influxPoints)
		return nil
	}
	if devices != nil {
		for _, asset := range devices.UAPs {
			errs = append(errs, process(asset))
		}
		for _, asset := range devices.USGs {
			errs = append(errs, process(asset))
		}
		for _, asset := range devices.USWs {
			errs = append(errs, process(asset))
		}
	}
	if clients != nil {
		for _, asset := range clients.UCLs {
			errs = append(errs, process(asset))
		}
	}
	return
}

// hasErr checks a list of errors for a non-nil.
func hasErr(errs []error) bool {
	for _, err := range errs {
		if err != nil {
			return true
		}
	}
	return false
}

// logErrors writes a slice of errors, with a prefix, to log-out.
func logErrors(errs []error, prefix string) {
	for _, err := range errs {
		if err != nil {
			log.Println("[ERROR]", prefix+":", err.Error())
		}
	}
}

// StringInSlice returns true if a string is in a slice.
func StringInSlice(str string, slc []string) bool {
	for _, s := range slc {
		if strings.EqualFold(s, str) {
			return true
		}
	}
	return false
}
