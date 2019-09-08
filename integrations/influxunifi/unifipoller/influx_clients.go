package unifipoller

import (
	"time"

	influx "github.com/influxdata/influxdb1-client/v2"
	"golift.io/unifi"
)

// ClientPoints generates Unifi Client datapoints for InfluxDB.
// These points can be passed directly to influx.
func ClientPoints(c *unifi.Client, now time.Time) ([]*influx.Point, error) {
	tags := map[string]string{
		"id":                 c.ID,
		"mac":                c.Mac,
		"user_id":            c.UserID,
		"site_id":            c.SiteID,
		"site_name":          c.SiteName,
		"network_id":         c.NetworkID,
		"usergroup_id":       c.UserGroupID,
		"ap_mac":             c.ApMac,
		"gw_mac":             c.GwMac,
		"sw_mac":             c.SwMac,
		"ap_name":            c.ApName,
		"gw_name":            c.GwName,
		"sw_name":            c.SwName,
		"oui":                c.Oui,
		"radio_name":         c.RadioName,
		"radio":              c.Radio,
		"radio_proto":        c.RadioProto,
		"name":               c.Name,
		"fixed_ip":           c.FixedIP,
		"sw_port":            c.SwPort.Txt,
		"os_class":           c.OsClass.Txt,
		"os_name":            c.OsName.Txt,
		"dev_cat":            c.DevCat.Txt,
		"dev_id":             c.DevID.Txt,
		"dev_vendor":         c.DevVendor.Txt,
		"dev_family":         c.DevFamily.Txt,
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
		"channel":            c.Channel.Txt,
		"vlan":               c.Vlan.Txt,
	}
	fields := map[string]interface{}{
		"anomalies":              c.Anomalies,
		"ip":                     c.IP,
		"essid":                  c.Essid,
		"bssid":                  c.Bssid,
		"radio_desc":             c.RadioDescription,
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
		"wifi_tx_attempts":       c.WifiTxAttempts,
		"wired-rx_bytes":         c.WiredRxBytes,
		"wired-rx_bytes-r":       c.WiredRxBytesR,
		"wired-rx_packets":       c.WiredRxPackets,
		"wired-tx_bytes":         c.WiredTxBytes,
		"wired-tx_bytes-r":       c.WiredTxBytesR,
		"wired-tx_packets":       c.WiredTxPackets,
		"dpi_app":                c.DpiStats.App.Val,
		"dpi_cat":                c.DpiStats.Cat.Val,
		"dpi_rx_bytes":           c.DpiStats.RxBytes.Val,
		"dpi_rx_packets":         c.DpiStats.RxPackets.Val,
		"dpi_tx_bytes":           c.DpiStats.TxBytes.Val,
		"dpi_tx_packets":         c.DpiStats.TxPackets.Val,
	}
	pt, err := influx.NewPoint("clients", tags, fields, now)
	if err != nil {
		return nil, err
	}
	return []*influx.Point{pt}, nil
}
