package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/golift/unifi"
	influx "github.com/influxdata/influxdb1-client/v2"
	"github.com/naoina/toml"
	flag "github.com/ogier/pflag"
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
		log.Fatalf("Config Error '%v': %v", configFile, err)
	}
	// Create an authenticated session to the Unifi Controller.
	controller, err := unifi.GetController(config.UnifiUser, config.UnifiPass, config.UnifiBase, config.VerifySSL)
	if err != nil {
		log.Fatalln("Unifi Controller Error:", err)
	} else if !config.Quiet {
		log.Println("Authenticated to Unifi Controller @", config.UnifiBase, "as user", config.UnifiUser)
	}
	controller.ErrorLog = log.Printf
	if log.SetFlags(0); config.Debug {
		controller.DebugLog = log.Printf
		log.SetFlags(log.Lshortfile | log.Lmicroseconds | log.Ldate)
		log.Println("Debug Logging Enabled")
	}
	infdb, err := influx.NewHTTPClient(influx.HTTPConfig{
		Addr:     config.InfluxURL,
		Username: config.InfluxUser,
		Password: config.InfluxPass,
	})
	if err != nil {
		log.Fatalln("InfluxDB Error:", err)
	}
	if config.Quiet {
		// Doing it this way allows debug error logs (line numbers, etc)
		controller.DebugLog = nil
	} else {
		log.Println("Logging Unifi Metrics to InfluXDB @", config.InfluxURL, "as user", config.InfluxUser)
		log.Println("Polling Unifi Controller, interval:", config.Interval.value)
	}
	config.PollUnifiController(controller, infdb)
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
		VerifySSL:  defaultVerifySSL,
		Debug:      defaultDebug,
		Quiet:      defaultQuiet,
		Interval:   Dur{value: defaultInterval},
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
	log.Println("[INFO] Everyting checks out! Beginning Poller Routine.")
	ticker := time.NewTicker(c.Interval.value)
	for range ticker.C {
		if clients, err := controller.GetClients(); err != nil {
			logErrors([]error{err}, "uni.GetClients()")
		} else if devices, err := controller.GetDevices(); err != nil {
			logErrors([]error{err}, "uni.GetDevices()")
		} else if bp, err := influx.NewBatchPoints(influx.BatchPointsConfig{Database: c.InfluxDB}); err != nil {
			logErrors([]error{err}, "influx.NewBatchPoints")
		} else if errs := batchPoints(devices, clients, bp); errs != nil && hasErr(errs) {
			logErrors(errs, "asset.Points()")
		} else if err := infdb.Write(bp); err != nil {
			logErrors([]error{err}, "infdb.Write(bp)")
		} else if !c.Quiet {
			log.Println("[INFO] Logged Unifi States. Clients:", len(clients.UCLs), "- Wireless APs:",
				len(devices.UAPs), "Gateways:", len(devices.USGs), "Switches:", len(devices.USWs))
		}
	}
}

// batchPoints combines all device and client data into influxdb data points.
func batchPoints(devices *unifi.Devices, clients *unifi.Clients, batchPoints influx.BatchPoints) (errs []error) {
	process := func(asset Asset) error {
		influxPoints, err := asset.Points()
		if err != nil {
			return err
		}
		batchPoints.AddPoints(influxPoints)
		return nil
	}
	for _, asset := range devices.UAPs {
		errs = append(errs, process(asset))
	}
	for _, asset := range devices.USGs {
		errs = append(errs, process(asset))
	}
	for _, asset := range devices.USWs {
		errs = append(errs, process(asset))
	}
	for _, asset := range clients.UCLs {
		errs = append(errs, process(asset))
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
