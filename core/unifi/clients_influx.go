package unifi

import (
	"strconv"
	"time"

	influx "github.com/influxdata/influxdb1-client/v2"
)

// Points generates Unifi Client datapoints for InfluxDB.
// These points can be passed directly to influx.
func (c UCL) Points() ([]*influx.Point, error) {
	// Fix name and hostname fields. Sometimes one or the other is blank.
	switch {
	case c.Hostname == "" && c.Name == "":
		c.Hostname = "-no-name-"
		c.Name = "-no-name-"
	case c.Hostname == "" && c.Name != "":
		c.Hostname = c.Name
	case c.Name == "" && c.Hostname != "":
		c.Name = c.Hostname
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
		"authorized":         c.Authorized.Txt,
		"is_11r":             c.Is11R.Txt,
		"is_wired":           c.IsWired.Txt,
		"is_guest":           c.IsGuest.Txt,
		"is_guest_by_uap":    c.IsGuestByUAP.Txt,
		"is_guest_by_ugw":    c.IsGuestByUGW.Txt,
		"is_guest_by_usw":    c.IsGuestByUSW.Txt,
		"noted":              c.Noted.Txt,
		"powersave_enabled":  c.PowersaveEnabled.Txt,
		"qos_policy_applied": c.QosPolicyApplied.Txt,
		"use_fixedip":        c.UseFixedIP.Txt,
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
	if err != nil {
		return nil, err
	}
	return []*influx.Point{pt}, nil
}
