package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/davidnewhall/unifi-poller/unidev"

	influx "github.com/influxdata/influxdb/client/v2"
	"github.com/naoina/toml"
	flg "github.com/ogier/pflag"
	"github.com/pkg/errors"
)

func main() {
	flg.Usage = func() {
		fmt.Println("Usage: unifi-poller [--config=filepath] [--debug] [--version]")
		flg.PrintDefaults()
	}
	configFile := flg.StringP("config", "c", defaultConfFile, "Poller Config File (TOML Format)")
	flg.BoolVarP(&Debug, "debug", "D", false, "Turn on the Spam (default false)")
	version := flg.BoolP("version", "v", false, "Print the version and exit.")
	flg.Parse()
	if *version {
		fmt.Println("unifi-poller version:", Version)
		os.Exit(0) // don't run anything else.
	}
	if log.SetFlags(0); Debug {
		log.SetFlags(log.Lshortfile | log.Lmicroseconds | log.Ldate)
		unidev.Debug = true
	}
	config, err := GetConfig(*configFile)
	if err != nil {
		flg.Usage()
		log.Fatalln("Config Error:", err)
	}
	// Create an authenticated session to the Unifi Controller.
	unifi, err := unidev.AuthController(config.UnifiUser, config.UnifiPass, config.UnifiBase)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Authenticated to Unifi Controller @", config.UnifiBase, "as user", config.UnifiUser)

	infdb, err := influx.NewHTTPClient(influx.HTTPConfig{
		Addr:     config.InfluxURL,
		Username: config.InfluxUser,
		Password: config.InfluxPass,
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Logging Unifi Metrics to InfluXDB @", config.InfluxURL, "as user", config.InfluxUser)
	log.Println("Polling Unifi Controller, interval:", config.Interval.value)
	config.PollUnifiController(infdb, unifi)
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
	}
	if buf, err := ioutil.ReadFile(configFile); err != nil {
		return config, err
		// This is where the defaults in the config variable are overwritten.
	} else if err := toml.Unmarshal(buf, &config); err != nil {
		return config, errors.Wrap(err, "invalid config")
	}
	return config, nil
}

// PollUnifiController runs forever, polling and pushing.
func (c *Config) PollUnifiController(infdb influx.Client, unifi *unidev.AuthedReq) {
	ticker := time.NewTicker(c.Interval.value)
	for range ticker.C {
		clients, err := unifi.GetUnifiClients()
		if err != nil {
			log.Println("unifi.GetUnifiClients():", err)
			continue
		}
		devices, err := unifi.GetUnifiDevices()
		if err != nil {
			log.Println("unifi.GetUnifiDevices():", err)
			continue
		}
		bp, err := influx.NewBatchPoints(influx.BatchPointsConfig{
			Database: c.InfluxDB,
		})
		if err != nil {
			log.Println("influx.NewBatchPoints:", err)
			continue
		}

		for _, asset := range append(clients, devices...) {
			if pt, errr := asset.Points(); errr != nil {
				log.Println("asset.Points():", errr)
			} else {
				bp.AddPoints(pt)
			}
		}

		if err = infdb.Write(bp); err != nil {
			log.Println("infdb.Write(bp):", err)
			continue
		}
		log.Println("Logged client state. Clients:", len(clients), "- Devices:", len(devices))
	}
}
