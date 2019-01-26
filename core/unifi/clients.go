package unifi

import (
	"strconv"
	"time"

	influx "github.com/influxdata/influxdb/client/v2"
)

// Points generates Unifi Client datapoints for InfluxDB.
// These points can be passed directly to influx.
func (c *UCL) Points() ([]*influx.Point, error) {
	var points []*influx.Point
	// Fix name and hostname fields. Sometimes one or the other is blank.
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
		"network_id":         c.NetworkID,
		"usergroup_id":       c.UserGroupID,
		"ap_mac":             c.ApMac,
		"gw_mac":             c.GwMac,
		"sw_mac":             c.SwMac,
		"oui":                c.Oui,
		"radio_name":         c.RadioName,
		"radio":              c.Radio,
		"radio_proto":        c.RadioProto,
		"name":               c.Name,
		"fixed_ip":           c.FixedIP,
		"sw_port":            strconv.Itoa(c.SwPort),
		"os_class":           strconv.Itoa(c.OsClass),
		"os_name":            strconv.Itoa(c.OsName),
		"dev_cat":            strconv.Itoa(c.DevCat),
		"dev_id":             strconv.Itoa(c.DevID),
		"dev_family":         strconv.Itoa(c.DevFamily),
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
		"ip":                     c.IP,
		"essid":                  c.Essid,
		"bssid":                  c.Bssid,
		"hostname":               c.Hostname,
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
		"network":                c.Network,
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
		"wired-rx_bytes":         c.WiredRxBytes,
		"wired-rx_bytes-r":       c.WiredRxBytesR,
		"wired-rx_packets":       c.WiredRxPackets,
		"wired-tx_bytes":         c.WiredTxBytes,
		"wired-tx_bytes-r":       c.WiredTxBytesR,
		"wired-tx_packets":       c.WiredTxPackets,
	}
	pt, err := influx.NewPoint("clients", tags, fields, time.Now())
	if err == nil {
		points = append(points, pt)
	}
	return points, err
}
