package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"time"

	influx "github.com/influxdata/influxdb/client/v2"
	"github.com/pkg/errors"
)

func main() {
	config := GetConfig()
	if err := config.AuthController(); err != nil {
		log.Fatal(err)
	}
	log.Println("Authenticated to Unifi Controller", config.UnifiBase, "as user", config.UnifiUser)

	infdb, err := influx.NewHTTPClient(influx.HTTPConfig{
		Addr:     config.InfluxURL,
		Username: config.InfluxUser,
		Password: config.InfluxPass,
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Logging Unifi Metrics to InfluXDB @", config.InfluxURL, "as user", config.InfluxUser)
	log.Println("Polling Unifi Controller, interval:", config.Interval)
	config.PollUnifiController(infdb)
}

// GetConfig parses and returns our configuration data.
func GetConfig() Config {
	// TODO: A real config file.
	interval, err := time.ParseDuration(os.Getenv("INTERVAL"))
	if err != nil {
		log.Println("Invalid Interval, defaulting to", defaultInterval)
		interval = time.Duration(defaultInterval)
	}
	return Config{
		InfluxURL:  os.Getenv("INFLUXDB_URL"),
		InfluxUser: os.Getenv("INFLUXDB_USERNAME"),
		InfluxPass: os.Getenv("INFLUXDB_PASSWORD"),
		InfluxDB:   os.Getenv("INFLUXDB_DATABASE"),
		UnifiUser:  os.Getenv("UNIFI_USERNAME"),
		UnifiPass:  os.Getenv("UNIFI_PASSWORD"),
		UnifiBase:  "https://" + os.Getenv("UNIFI_ADDR") + ":" + os.Getenv("UNIFI_PORT"),
		Interval:   interval,
	}
}

// AuthController creates a http.Client with authenticated cookies.
// Used to make additional, authenticated requests to the APIs.
func (c *Config) AuthController() error {
	json := `{"username": "` + c.UnifiUser + `","password": "` + c.UnifiPass + `"}`
	jar, err := cookiejar.New(nil)
	if err != nil {
		return errors.Wrap(err, "cookiejar.New(nil)")
	}
	c.uniClient = &http.Client{
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
		Jar:       jar,
	}
	if req, err := c.uniRequest(LoginPath, json); err != nil {
		return errors.Wrap(err, "c.uniRequest(LoginPath, json)")
	} else if resp, err := c.uniClient.Do(req); err != nil {
		return errors.Wrap(err, "c.uniClient.Do(req)")
	} else if resp.StatusCode != http.StatusOK {
		return errors.Errorf("authentication failed (%v): %v (status: %v/%v)",
			c.UnifiUser, c.UnifiBase+LoginPath, resp.StatusCode, resp.Status)
	}
	return nil
}

// PollUnifiController runs forever, polling and pushing.
func (c *Config) PollUnifiController(infdb influx.Client) {
	ticker := time.NewTicker(c.Interval)
	for range ticker.C {
		clients, err := c.GetUnifiClients()
		if err != nil {
			log.Println("GetUnifiClients(unifi):", err)
			continue
		}
		bp, err := influx.NewBatchPoints(influx.BatchPointsConfig{
			Database: c.InfluxDB,
		})
		if err != nil {
			log.Println("influx.NewBatchPoints:", err)
			continue
		}

		for _, client := range clients {
			if pt, errr := client.Point(); errr != nil {
				log.Println("client.Point():", errr)
			} else {
				bp.AddPoint(pt)
			}
		}
		if err = infdb.Write(bp); err != nil {
			log.Println("infdb.Write(bp):", err)
			continue
		}
		log.Println("Logged client state. Clients:", len(clients))
	}
}

// GetUnifiClients returns a response full of clients' data from the Unifi Controller.
func (c *Config) GetUnifiClients() ([]Client, error) {
	response := &ClientResponse{}
	if req, err := c.uniRequest(ClientPath, ""); err != nil {
		return nil, err
	} else if resp, err := c.uniClient.Do(req); err != nil {
		return nil, err
	} else if body, err := ioutil.ReadAll(resp.Body); err != nil {
		return nil, err
	} else if err = json.Unmarshal(body, response); err != nil {
		return nil, err
	} else if err = resp.Body.Close(); err != nil {
		log.Println("resp.Body.Close():", err) // Not fatal? Just log it.
	}
	return response.Clients, nil
}

// uniRequest is a small helper function that adds an Accept header.
func (c *Config) uniRequest(url string, params string) (req *http.Request, err error) {
	if params != "" {
		req, err = http.NewRequest("POST", c.UnifiBase+url, bytes.NewBufferString(params))
	} else {
		req, err = http.NewRequest("GET", c.UnifiBase+url, nil)
	}
	if err == nil {
		req.Header.Add("Accept", "application/json")
	}
	return
}
