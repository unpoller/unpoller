package influxunifi

import (
	"golift.io/unifi"
)

// batchClient generates Unifi Client datapoints for InfluxDB.
// These points can be passed directly to influx.
func (u *InfluxUnifi) batchClient(r report, s *unifi.Client) {
	tags := map[string]string{
		"mac":         s.Mac,
		"site_name":   s.SiteName,
		"ap_name":     s.ApName,
		"gw_name":     s.GwName,
		"sw_name":     s.SwName,
		"oui":         s.Oui,
		"radio_name":  s.RadioName,
		"radio":       s.Radio,
		"radio_proto": s.RadioProto,
		"name":        s.Name,
		"fixed_ip":    s.FixedIP,
		"sw_port":     s.SwPort.Txt,
		"os_class":    s.OsClass.Txt,
		"os_name":     s.OsName.Txt,
		"dev_cat":     s.DevCat.Txt,
		"dev_id":      s.DevID.Txt,
		"dev_vendor":  s.DevVendor.Txt,
		"dev_family":  s.DevFamily.Txt,
		"is_wired":    s.IsWired.Txt,
		"is_guest":    s.IsGuest.Txt,
		"use_fixedip": s.UseFixedIP.Txt,
		"channel":     s.Channel.Txt,
		"vlan":        s.Vlan.Txt,
	}
	fields := map[string]interface{}{
		"anomalies":        s.Anomalies,
		"ip":               s.IP,
		"essid":            s.Essid,
		"bssid":            s.Bssid,
		"channel":          s.Channel.Val,
		"hostname":         s.Name,
		"radio_desc":       s.RadioDescription,
		"satisfaction":     s.Satisfaction.Val,
		"bytes_r":          s.BytesR,
		"ccq":              s.Ccq,
		"noise":            s.Noise,
		"note":             s.Note,
		"roam_count":       s.RoamCount,
		"rssi":             s.Rssi,
		"rx_bytes":         s.RxBytes,
		"rx_bytes_r":       s.RxBytesR,
		"rx_packets":       s.RxPackets,
		"rx_rate":          s.RxRate,
		"signal":           s.Signal,
		"tx_bytes":         s.TxBytes,
		"tx_bytes_r":       s.TxBytesR,
		"tx_packets":       s.TxPackets,
		"tx_retries":       s.TxRetries,
		"tx_power":         s.TxPower,
		"tx_rate":          s.TxRate,
		"uptime":           s.Uptime,
		"wifi_tx_attempts": s.WifiTxAttempts,
		"wired-rx_bytes":   s.WiredRxBytes,
		"wired-rx_bytes-r": s.WiredRxBytesR,
		"wired-rx_packets": s.WiredRxPackets,
		"wired-tx_bytes":   s.WiredTxBytes,
		"wired-tx_bytes-r": s.WiredTxBytesR,
		"wired-tx_packets": s.WiredTxPackets,
		/*
			"dpi_app":          c.DpiStats.App.Val,
			"dpi_cat":          c.DpiStats.Cat.Val,
			"dpi_rx_bytes":     c.DpiStats.RxBytes.Val,
			"dpi_rx_packets":   c.DpiStats.RxPackets.Val,
			"dpi_tx_bytes":     c.DpiStats.TxBytes.Val,
			"dpi_tx_packets":   c.DpiStats.TxPackets.Val,
		*/
	}
	r.send(&metric{Table: "clients", Tags: tags, Fields: fields})
}
