package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"time"

	influx "github.com/influxdata/influxdb/client/v2"
)

func main() {
	config := GetConfig()
	if err := config.AuthController(); err != nil {
		log.Fatal(err)
	}
	log.Println("Successfully authenticated to Unifi Controller!")

	infdb, err := influx.NewHTTPClient(influx.HTTPConfig{
		Addr:     config.InfluxAddr,
		Username: config.InfluxUser,
		Password: config.InfluxPass,
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Polling Unifi Controller, interval:", config.Interval)
	config.PollUnifiController(infdb)
}

// GetConfig parses and returns our configuration data.
func GetConfig() Config {
	// TODO: A real config file.
	var err error
	config := Config{
		InfluxAddr: os.Getenv("INFLUXDB_ADDR"),
		InfluxUser: os.Getenv("INFLUXDB_USERNAME"),
		InfluxPass: os.Getenv("INFLUXDB_PASSWORD"),
		InfluxDB:   os.Getenv("INFLUXDB_DATABASE"),
		UnifiUser:  os.Getenv("UNIFI_USERNAME"),
		UnifiPass:  os.Getenv("UNIFI_PASSWORD"),
		UnifiBase:  "https://" + os.Getenv("UNIFI_ADDR") + ":" + os.Getenv("UNIFI_PORT"),
	}
	if config.Interval, err = time.ParseDuration(os.Getenv("INTERVAL")); err != nil {
		log.Println("Invalid Interval, defaulting to 15 seconds.")
		config.Interval = time.Duration(time.Second * 15)
	}
	return config
}

// AuthController creates a http.Client with authenticated cookies.
// Used to make additional, authenticated requests to the APIs.
func (c *Config) AuthController() error {
	json := `{"username": "` + c.UnifiUser + `","password": "` + c.UnifiPass + `"}`
	jar, err := cookiejar.New(nil)
	if err != nil {
		return err
	}
	c.uniClient = &http.Client{
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
		Jar:       jar,
	}
	if req, err := c.uniRequest(LoginPath, json); err != nil {
		return err
	} else if resp, err := c.uniClient.Do(req); err != nil {
		return err
	} else if resp.StatusCode != http.StatusOK {
		return errors.New("Error Authenticating with Unifi Controller")
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
			bp.AddPoint(client.Point())
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
