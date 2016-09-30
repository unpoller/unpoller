package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"time"

	influx "github.com/influxdata/influxdb/client/v2"
)

var (
	API   *http.Client
	stats influx.Client
)

type Response struct {
	Data []Client
	Meta struct {
		Rc string
	}
}

type DpiStat struct {
	App       int64
	Cat       int
	RxBytes   int64
	RxPackets int64
	TxBytes   int64
	TxPackets int64
}

type Client struct {
	ID                  string `json:"_id"`
	IsGuestByUAP        bool   `json:"_is_guest_by_uap"`
	IsGuestByUGW        bool   `json:"_is_guest_by_ugw"`
	LastSeenByUAP       int64  `json:"_last_seen_by_uap"`
	LastSeenByUGW       int64  `json:"_last_seen_by_ugw"`
	UptimeByUAP         int64  `json:"_uptime_by_uap"`
	UptimeByUGW         int64  `json:"_uptime_by_ugw"`
	ApMac               string `json:"ap_mac"`
	AssocTime           int64  `json:"assoc_time"`
	Authorized          bool
	Bssid               string
	BytesR              int `json:"bytes-r"`
	Ccq                 int
	Channel             int
	DpiStats            []DpiStat `json:"dpi_stats"`
	DpiStatsLastUpdated int64     `json:"dpi_stats_last_updated"`
	Essid               string
	FirstSeen           int64  `json:"first_seen"`
	FixedIP             string `json:"fixed_ip"`
	Hostname            string
	GwMac               string `json:"gw_mac"`
	IdleTime            int64  `json:"idle_time"`
	Ip                  string
	IsGuest             bool  `json:"is_guest"`
	IsWired             bool  `json:"is_wired"`
	LastSeen            int64 `json:"last_seen"`
	LatestAssocTime     int64 `json:"latest_assoc_time"`
	Mac                 string
	Name                string
	Network             string
	NetworkID           string `json:"network_id"`
	Noise               int
	Oui                 string
	PowersaveEnabled    bool `json:"powersave_enabled"`
	QosPolicyApplied    bool `json:"qos_policy_applied"`
	Radio               string
	RadioProto          string `json:"radio_proto"`
	RoamCount           int    `json:"roam_count"`
	Rssi                int
	RxBytes             int64 `json:"rx_bytes"`
	RxBytesR            int   `json:"rx_bytes-r"`
	RxPackets           int64 `json:"rx_packets"`
	RxRate              int64 `json:"rx_rate"`
	Signal              int
	SiteID              string `json:"site_id"`
	TxBytes             int64  `json:"tx_bytes"`
	TxBytesR            int    `json:"tx_bytes-r"`
	TxPackets           int64  `json:"tx_packets"`
	TxPower             int    `json:"tx_power"`
	TxRate              int64  `json:"tx_rate"`
	Uptime              int64
	UserID              string `json:"user_id"`
	Vlan                int
}

func (c Client) Point() *influx.Point {
	tags := map[string]string{
		"mac":     c.Mac,
		"user_id": c.UserID,
		"site_id": c.SiteID,
		"ip":      c.Ip,
		"essid":   c.Essid,
		"network": c.Network,
		"ap_mac":  c.ApMac,
		"name":    c.Name,
	}
	fields := map[string]interface{}{
		"is_guest_by_uap":        c.IsGuestByUAP,
		"is_guest_by_ugw":        c.IsGuestByUGW,
		"authorized":             c.Authorized,
		"last_seen_by_uap":       c.LastSeenByUAP,
		"last_seen_by_ugw":       c.LastSeenByUGW,
		"uptime_by_uap":          c.UptimeByUAP,
		"uptime_by_ugw":          c.UptimeByUGW,
		"assoc_time":             c.AssocTime,
		"bytes_r":                c.BytesR,
		"ccq":                    c.Ccq,
		"channel":                c.Channel,
		"dpi_stats_last_updated": c.DpiStatsLastUpdated,
		"first_seen":             c.FirstSeen,
		"idle_time":              c.IdleTime,
		"last_seen":              c.LastSeen,
		"latest_assoc_time":      c.LatestAssocTime,
		"noise":                  c.Noise,
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
		"vlan":                   c.Vlan,
	}

	pt, err := influx.NewPoint("client_state", tags, fields, time.Now())
	if err != nil {
		return nil
	}

	return pt
}

func main() {
	var err error
	tickRate := os.Getenv("TICK_RATE")
	interval, err := time.ParseDuration(tickRate)
	if err != nil {
		panic(err)
	}

	API, err = login()
	if err != nil {
		panic(err)
	}

	database := os.Getenv("INFLUXDB_DATABASE")
	stats, err = influx.NewHTTPClient(influx.HTTPConfig{
		Addr:     os.Getenv("INFLUXDB_ADDR"),
		Username: os.Getenv("INFLUXDB_USERNAME"),
		Password: os.Getenv("INFLUXDB_PASSWORD"),
	})
	if err != nil {
		panic(err)
	}

	log.Printf("Starting to poll Unifi every %+v\n", interval)
	for {
		devices, err := fetch()
		if err != nil {
			log.Println(err)
		} else {
			bp, _ := influx.NewBatchPoints(influx.BatchPointsConfig{
				Database: database,
			})

			for _, device := range devices {
				bp.AddPoint(device.Point())
			}

			err = stats.Write(bp)
			if err != nil {
				log.Println(err)
			}

			log.Println("Logged client state...")
		}

		time.Sleep(interval)
	}
}

func fetch() ([]Client, error) {
	format := "https://%s:%s/api/s/default/stat/sta"
	url := fmt.Sprintf(format, os.Getenv("UNIFI_ADDR"), os.Getenv("UNIFI_PORT"))
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/json")
	resp, err := API.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	response := &Response{}
	err = json.Unmarshal(body, response)
	if err != nil {
		return nil, err
	}

	return response.Data, nil
}

func login() (*http.Client, error) {
	url := fmt.Sprintf("https://%s:%s/api/login", os.Getenv("UNIFI_ADDR"), os.Getenv("UNIFI_PORT"))
	auth := map[string]string{
		"username": os.Getenv("UNIFI_USERNAME"),
		"password": os.Getenv("UNIFI_PASSWORD"),
	}
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Transport: transport,
		Jar:       jar,
	}
	json, _ := json.Marshal(auth)
	params := bytes.NewReader(json)
	req, err := http.NewRequest("POST", url, params)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Not a successful request")
	}

	return client, nil
}
