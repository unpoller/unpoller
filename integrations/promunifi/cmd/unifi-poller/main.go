package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/golift/unifi"
	"github.com/influxdata/influxdb/client/v2"
	"github.com/naoina/toml"
	flag "github.com/ogier/pflag"
)

func main() {

	configFile := parseFlags()
	log.Println("Unifi-Poller Starting Up! PID:", os.Getpid())
	config, err := GetConfig(configFile)
	if err != nil {
		flag.Usage()
		log.Fatalf("Config Error '%v': %v", configFile, err)
	} else if log.SetFlags(0); config.Debug {
		log.SetFlags(log.Lshortfile | log.Lmicroseconds | log.Ldate)
		log.Println("Debug Logging Enabled")
	}
	log.Println("Loaded Configuration:", configFile)
	// Create an authenticated session to the Unifi Controller.
	device, err := unifi.AuthController(config.UnifiUser, config.UnifiPass, config.UnifiBase, config.VerifySSL)
	if err != nil {
		log.Fatalln("Unifi Controller Error:", err)
	} else if !config.Quiet {
		log.Println("Authenticated to Unifi Controller @", config.UnifiBase, "as user", config.UnifiUser)
	}
	infdb, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     config.InfluxURL,
		Username: config.InfluxUser,
		Password: config.InfluxPass,
	})
	if err != nil {
		log.Fatalln("InfluxDB Error:", err)
	} else if config.Quiet {
		// Doing it this way allows debug error logs (line numbers, etc)
		unifi.Debug = false
	} else {
		log.Println("Logging Unifi Metrics to InfluXDB @", config.InfluxURL, "as user", config.InfluxUser)
		log.Println("Polling Unifi Controller, interval:", config.Interval.value)
	}
	log.Println("Everyting checks out! Beginning Poller Routine.")
	config.PollUnifiController(infdb, device, config.Quiet)
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
	return config, nil
}

// PollUnifiController runs forever, polling and pushing.
func (c *Config) PollUnifiController(infdb client.Client, device *unifi.AuthedReq, quiet bool) {
	ticker := time.NewTicker(c.Interval.value)
	for range ticker.C {
		var clients, devices []unifi.Asset
		var bp client.BatchPoints
		var err error
		if clients, err = device.GetUnifiClientAssets(); err != nil {
			log.Println("ERROR unifi.GetUnifiClientsAssets():", err)
		} else if devices, err = device.GetUnifiDeviceAssets(); err != nil {
			log.Println("ERROR unifi.GetUnifiDeviceAssets():", err)
		} else if bp, err = client.NewBatchPoints(client.BatchPointsConfig{Database: c.InfluxDB}); err != nil {
			log.Println("ERROR influx.NewBatchPoints:", err)
		}
		if err != nil {
			continue
		}
		for _, asset := range append(clients, devices...) {
			if pt, errr := asset.Points(); errr != nil {
				log.Println("ERROR asset.Points():", errr)
			} else {
				bp.AddPoints(pt)
			}
		}
		if err = infdb.Write(bp); err != nil {
			log.Println("ERROR infdb.Write(bp):", err)
			continue
		}
		if !quiet {
			log.Println("Logged Unifi States. Clients:", len(clients), "- Devices:", len(devices))
		}
	}
}
