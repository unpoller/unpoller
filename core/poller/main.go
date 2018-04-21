package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strconv"
	"time"

	influx "github.com/influxdata/influxdb/client/v2"
)

// Point generates a client's datapoint for InfluxDB.
func (c Client) Point() *influx.Point {
	if c.Name == "" && c.Hostname != "" {
		c.Name = c.Hostname
	} else if c.Hostname == "" && c.Name != "" {
		c.Hostname = c.Name
	} else if c.Hostname == "" && c.Name == "" {
		c.Hostname = "-no-name-"
		c.Name = "-no-name-"
	}
	tags := map[string]string{
		"id":                 c.ID,
		"mac":                c.Mac,
		"user_id":            c.UserID,
		"site_id":            c.SiteID,
		"ip":                 c.IP,
		"fixed_ip":           c.FixedIP,
		"essid":              c.Essid,
		"bssid":              c.Bssid,
		"network":            c.Network,
		"network_id":         c.NetworkID,
		"usergroup_id":       c.UserGroupID,
		"ap_mac":             c.ApMac,
		"gw_mac":             c.GwMac,
		"sw_mac":             c.SwMac,
		"sw_port":            strconv.Itoa(c.SwPort),
		"oui":                c.Oui,
		"name":               c.Name,
		"hostname":           c.Hostname,
		"radio_name":         c.RadioName,
		"radio":              c.Radio,
		"radio_proto":        c.RadioProto,
		"authorized":         strconv.FormatBool(c.Authorized),
		"is_11r":             strconv.FormatBool(c.Is11R),
		"is_wired":           strconv.FormatBool(c.IsWired),
		"is_guest":           strconv.FormatBool(c.IsGuest),
		"is_guest_by_uap":    strconv.FormatBool(c.IsGuestByUAP),
		"is_guest_by_ugw":    strconv.FormatBool(c.IsGuestByUGW),
		"is_guest_by_usw":    strconv.FormatBool(c.IsGuestByUSW),
		"noted":              strconv.FormatBool(c.Noted),
		"powersave_enabled":  strconv.FormatBool(c.PowersaveEnabled),
		"qos_policy_applied": strconv.FormatBool(c.QosPolicyApplied),
		"use_fixedip":        strconv.FormatBool(c.UseFixedIP),
		"channel":            strconv.Itoa(c.Channel),
		"vlan":               strconv.Itoa(c.Vlan),
	}
	fields := map[string]interface{}{
		"dpi_stats_last_updated": c.DpiStatsLastUpdated,
		"last_seen_by_uap":       c.LastSeenByUAP,
		"last_seen_by_ugw":       c.LastSeenByUGW,
		"last_seen_by_usw":       c.LastSeenByUSW,
		"uptime_by_uap":          c.UptimeByUAP,
		"uptime_by_ugw":          c.UptimeByUGW,
		"uptime_by_usw":          c.UptimeByUSW,
		"assoc_time":             c.AssocTime,
		"bytes_r":                c.BytesR,
		"ccq":                    c.Ccq,
		"first_seen":             c.FirstSeen,
		"idle_time":              c.IdleTime,
		"last_seen":              c.LastSeen,
		"latest_assoc_time":      c.LatestAssocTime,
		"noise":                  c.Noise,
		"note":                   c.Note,
		"roam_count":             c.RoamCount,
		"rssi":                   c.Rssi,
		"rx_bytes":               c.RxBytes,
		"rx_bytes_r":             c.RxBytesR,
		"rx_packets":             c.RxPackets,
		"rx_rate":                c.RxRate,
		"signal":                 c.Signal,
		"tx_bytes":               c.TxBytes,
		"tx_bytes_r":             c.TxBytesR,
		"tx_packets":             c.TxPackets,
		"tx_power":               c.TxPower,
		"tx_rate":                c.TxRate,
		"uptime":                 c.Uptime,
	}

	pt, err := influx.NewPoint("clients", tags, fields, time.Now())
	if err != nil {
		log.Println("Error creating point:", err)
		return nil
	}
	return pt
}

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

// GetUnifiClients returns a response full of clients' data from the Unifi Controller.
func (c *Config) GetUnifiClients() ([]Client, error) {
	response := &ClientResponse{}
	if req, err := uniRequest("GET", c.UnifiBase+ClientPath, nil); err != nil {
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

// AuthController creates a http.Client with authenticated cookies.
// Used to make additional, authenticated requests to the APIs.
func (c *Config) AuthController() error {
	json := bytes.NewBufferString(`{"username": "` + c.UnifiUser + `","password": "` + c.UnifiPass + `"}`)
	jar, err := cookiejar.New(nil)
	if err != nil {
		return err
	}
	c.uniClient = &http.Client{
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
		Jar:       jar,
	}
	if req, err := uniRequest("POST", c.UnifiBase+LoginPath, json); err != nil {
		return err
	} else if resp, err := c.uniClient.Do(req); err != nil {
		return err
	} else if resp.StatusCode != http.StatusOK {
		return errors.New("Error Authenticating with Unifi Controller")
	}
	return nil
}

func uniRequest(method string, url string, data io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, data)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	return req, nil
}
