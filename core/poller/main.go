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

const (
	// LoginPath is Unifi Controller Login API Path
	LoginPath = "/api/login"
	// ClientPath is Unifi Clients API Path
	ClientPath = "/api/s/default/stat/sta"
)

// ClientResponse marshalls the payload from the controller.
type ClientResponse struct {
	Clients []Client `json:"data"`
	Meta    struct {
		Rc string `json:"rc"`
	} `json:"meta"`
}

// Config represents the data needed to poll a controller and report to influxdb.
type Config struct {
	Interval   time.Duration `json:"interval",toml:"interval",yaml:"interval"`
	InfluxAddr string        `json:"influx_addr",toml:"influx_addr",yaml:"influx_addr"`
	InfluxUser string        `json:"influx_user",toml:"influx_user",yaml:"influx_user"`
	InfluxPass string        `json:"influx_pass",toml:"influx_pass",yaml:"influx_pass"`
	InfluxDB   string        `json:"influx_db",toml:"influx_db",yaml:"influx_db"`
	UnifiUser  string        `json:"unifi_user",toml:"unifi_user",yaml:"unifi_user"`
	UnifiPass  string        `json:"unifi_pass"toml:"unifi_pass",yaml:"unifi_pass"`
	UnifiBase  string        `json:"unifi_url",toml:"unifi_url",yaml:"unifi_url"`
	uniClient  *http.Client
}

// DpiStat is for deep packet inspection stats.
// Does not seem to exist in Unifi 5.7.20.
type DpiStat struct {
	App       int64
	Cat       int64
	RxBytes   int64
	RxPackets int64
	TxBytes   int64
	TxPackets int64
}

// Client defines all the data a connected-network client contains.
type Client struct {
	ID                  string    `json:"_id"`
	IsGuestByUAP        bool      `json:"_is_guest_by_uap"`
	IsGuestByUGW        bool      `json:"_is_guest_by_ugw"`
	IsGuestByUSW        bool      `json:"_is_guest_by_usw"`
	LastSeenByUAP       int64     `json:"_last_seen_by_uap"`
	LastSeenByUGW       int64     `json:"_last_seen_by_ugw"`
	LastSeenByUSW       int64     `json:"_last_seen_by_usw"`
	UptimeByUAP         int64     `json:"_uptime_by_uap"`
	UptimeByUGW         int64     `json:"_uptime_by_ugw"`
	UptimeByUSW         int64     `json:"_uptime_by_usw"`
	ApMac               string    `json:"ap_mac"`
	AssocTime           int64     `json:"assoc_time"`
	Authorized          bool      `json:"authorized"`
	Bssid               string    `json:"bssid"`
	BytesR              int64     `json:"bytes-r"`
	Ccq                 int64     `json:"ccq"`
	Channel             int       `json:"channel"`
	DpiStats            []DpiStat `json:"dpi_stats"`
	DpiStatsLastUpdated int64     `json:"dpi_stats_last_updated"`
	Essid               string    `json:"essid"`
	FirstSeen           int64     `json:"first_seen"`
	FixedIP             string    `json:"fixed_ip"`
	Hostname            string    `json:"hostname"`
	GwMac               string    `json:"gw_mac"`
	IdleTime            int64     `json:"idle_time"`
	IP                  string    `json:"ip"`
	Is11R               bool      `json:"is_11r"`
	IsGuest             bool      `json:"is_guest"`
	IsWired             bool      `json:"is_wired"`
	LastSeen            int64     `json:"last_seen"`
	LatestAssocTime     int64     `json:"latest_assoc_time"`
	Mac                 string    `json:"mac"`
	Name                string    `json:"name"`
	Network             string    `json:"network"`
	NetworkID           string    `json:"network_id"`
	Noise               int64     `json:"noise"`
	Note                string    `json:"note"`
	Noted               bool      `json:"noted"`
	Oui                 string    `json:"oui"`
	PowersaveEnabled    bool      `json:"powersave_enabled"`
	QosPolicyApplied    bool      `json:"qos_policy_applied"`
	Radio               string    `json:"radio"`
	RadioName           string    `json:"radio_name"`
	RadioProto          string    `json:"radio_proto"`
	RoamCount           int64     `json:"roam_count"`
	Rssi                int64     `json:"rssi"`
	RxBytes             int64     `json:"rx_bytes"`
	RxBytesR            int64     `json:"rx_bytes-r"`
	RxPackets           int64     `json:"rx_packets"`
	RxRate              int64     `json:"rx_rate"`
	Signal              int64     `json:"signal"`
	SiteID              string    `json:"site_id"`
	SwDepth             int       `json:"sw_depth"`
	SwMac               string    `json:"sw_mac"`
	SwPort              int       `json:"sw_port"`
	TxBytes             int64     `json:"tx_bytes"`
	TxBytesR            int64     `json:"tx_bytes-r"`
	TxPackets           int64     `json:"tx_packets"`
	TxPower             int64     `json:"tx_power"`
	TxRate              int64     `json:"tx_rate"`
	Uptime              int64     `json:"uptime"`
	UserID              string    `json:"user_id"`
	UserGroupID         string    `json:"usergroup_id"`
	UseFixedIP          bool      `json:"use_fixedip"`
	Vlan                int       `json:"vlan"`
}

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

	pt, err := influx.NewPoint("client_state", tags, fields, time.Now())
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
		uniResponse, err := c.GetUnifiClients()
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

		for _, client := range uniResponse.Clients {
			bp.AddPoint(client.Point())
		}
		if err = infdb.Write(bp); err != nil {
			log.Println("infdb.Write(bp):", err)
			continue
		}
		log.Println("Logged client state. Clients:", len(uniResponse.Clients))
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
func (c *Config) GetUnifiClients() (*ClientResponse, error) {
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
	return response, nil
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
