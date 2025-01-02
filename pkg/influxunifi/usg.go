package influxunifi

import (
	"github.com/unpoller/unifi/v5"
)

// usgT is used as a name for printed/logged counters.
const usgT = item("USG")

// batchUSG generates Unifi Gateway datapoints for InfluxDB.
// These points can be passed directly to influx.
func (u *InfluxUnifi) batchUSG(r report, s *unifi.USG) {
	if !s.Adopted.Val || s.Locating.Val {
		return
	}

	tags := map[string]string{
		"mac":       s.Mac,
		"site_name": s.SiteName,
		"source":    s.SourceName,
		"name":      s.Name,
		"version":   s.Version,
		"model":     s.Model,
		"serial":    s.Serial,
		"type":      s.Type,
	}
	fields := Combine(
		u.batchUDMtemps(s.Temperatures),
		u.batchSysStats(s.SysStats, s.SystemStats),
		u.batchUSGstats(s.SpeedtestStatus, s.Stat.Gw, s.Uplink),
		map[string]any{
			"ip":            s.IP,
			"bytes":         s.Bytes.Val,
			"last_seen":     s.LastSeen.Val,
			"license_state": s.LicenseState,
			"guest-num_sta": s.GuestNumSta.Val,
			"rx_bytes":      s.RxBytes.Val,
			"tx_bytes":      s.TxBytes.Val,
			"uptime":        s.Uptime.Val,
			"state":         s.State.Val,
			"user-num_sta":  s.UserNumSta.Val,
			"version":       s.Version,
			"num_desktop":   s.NumDesktop.Val,
			"num_handheld":  s.NumHandheld.Val,
			"num_mobile":    s.NumMobile.Val,
			"upgradeable":   s.Upgradable.Val,
			//			"speedtest_rundate": time.Unix(int64(s.SpeedtestStatus.Rundate.Val), 0).String(),
		},
	)

	r.addCount(usgT)
	r.send(&metric{Table: "usg", Tags: tags, Fields: fields})
	u.batchNetTable(r, tags, s.NetworkTable)
	u.batchUSGwans(r, tags, s.Wan1, s.Wan2)
}

func (u *InfluxUnifi) batchUSGstats(ss unifi.SpeedtestStatus, gw *unifi.Gw, ul unifi.Uplink) map[string]any {
	if gw == nil {
		return map[string]any{}
	}

	return map[string]any{
		"uplink_name":                    ul.Name,
		"uplink_latency":                 ul.Latency.Val,
		"uplink_speed":                   ul.Speed.Val,
		"uplink_max_speed":               ul.MaxSpeed.Val,
		"uplink_uptime":                  ul.Uptime.Val,
		"uplink_type":                    ul.Type,
		"speedtest-status_latency":       ss.Latency.Val,
		"speedtest-status_runtime":       ss.Runtime.Val,
		"speedtest-status_rundate":       ss.Rundate.Val,
		"speedtest-status_ping":          ss.StatusPing.Val,
		"speedtest-status_xput_download": ss.XputDownload.Val,
		"speedtest-status_xput_upload":   ss.XputUpload.Val,
		"lan-rx_bytes":                   gw.LanRxBytes.Val,
		"lan-rx_packets":                 gw.LanRxPackets.Val,
		"lan-tx_bytes":                   gw.LanTxBytes.Val,
		"lan-tx_packets":                 gw.LanTxPackets.Val,
		"lan-rx_dropped":                 gw.LanRxDropped.Val,
	}
}

func (u *InfluxUnifi) batchUSGwans(r report, tags map[string]string, wans ...unifi.Wan) {
	for _, wan := range wans {
		if !wan.Up.Val {
			continue
		}

		tags := map[string]string{
			"device_name": tags["name"],
			"site_name":   tags["site_name"],
			"source":      tags["source"],
			"ip":          wan.IP,
			"purpose":     wan.Name,
			"mac":         wan.Mac,
			"ifname":      wan.Ifname,
			"type":        wan.Type,
			"up":          wan.Up.Txt,
			"enabled":     wan.Enable.Txt,
		}
		fields := map[string]any{
			"bytes-r":      wan.BytesR.Val,
			"full_duplex":  wan.FullDuplex.Val,
			"gateway":      wan.Gateway,
			"is_uplink":    wan.IsUplink.Val,
			"max_speed":    wan.MaxSpeed.Val,
			"rx_bytes":     wan.RxBytes.Val,
			"rx_bytes-r":   wan.RxBytesR.Val,
			"rx_dropped":   wan.RxDropped.Val,
			"rx_errors":    wan.RxErrors.Val,
			"rx_broadcast": wan.RxBroadcast.Val,
			"rx_multicast": wan.RxMulticast.Val,
			"rx_packets":   wan.RxPackets.Val,
			"speed":        wan.Speed.Val,
			"tx_bytes":     wan.TxBytes.Val,
			"tx_bytes-r":   wan.TxBytesR.Val,
			"tx_dropped":   wan.TxDropped.Val,
			"tx_errors":    wan.TxErrors.Val,
			"tx_packets":   wan.TxPackets.Val,
			"tx_broadcast": wan.TxBroadcast.Val,
			"tx_multicast": wan.TxMulticast.Val,
		}

		r.send(&metric{Table: "usg_wan_ports", Tags: tags, Fields: fields})
	}
}

func (u *InfluxUnifi) batchNetTable(r report, tags map[string]string, nt unifi.NetworkTable) {
	for _, p := range nt {
		tags := map[string]string{
			"device_name": tags["name"],
			"site_name":   tags["site_name"],
			"source":      tags["source"],
			"up":          p.Up.Txt,
			"enabled":     p.Enabled.Txt,
			"ip":          p.IP,
			"mac":         p.Mac,
			"name":        p.Name,
			"domain_name": p.DomainName,
			"purpose":     p.Purpose,
			"is_guest":    p.IsGuest.Txt,
		}
		fields := map[string]any{
			"num_sta":    p.NumSta.Val,
			"rx_bytes":   p.RxBytes.Val,
			"rx_packets": p.RxPackets.Val,
			"tx_bytes":   p.TxBytes.Val,
			"tx_packets": p.TxPackets.Val,
		}

		r.send(&metric{Table: "usg_networks", Tags: tags, Fields: fields})
	}
}
