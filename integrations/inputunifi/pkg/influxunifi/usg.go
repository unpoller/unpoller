package influxunifi

import (
	"strings"

	"golift.io/unifi"
)

// batchUSG generates Unifi Gateway datapoints for InfluxDB.
// These points can be passed directly to influx.
func (u *InfluxUnifi) batchUSG(r report, s *unifi.USG) {
	if s.Stat.Gw == nil {
		s.Stat.Gw = &unifi.Gw{}
	}
	tags := map[string]string{
		"mac":       s.Mac,
		"site_name": s.SiteName,
		"name":      s.Name,
		"version":   s.Version,
		"model":     s.Model,
		"serial":    s.Serial,
		"type":      s.Type,
	}
	fields := map[string]interface{}{
		"ip":                             s.IP,
		"bytes":                          s.Bytes.Val,
		"last_seen":                      s.LastSeen.Val,
		"license_state":                  s.LicenseState,
		"guest-num_sta":                  s.GuestNumSta.Val,
		"rx_bytes":                       s.RxBytes.Val,
		"tx_bytes":                       s.TxBytes.Val,
		"uptime":                         s.Uptime.Val,
		"state":                          s.State.Val,
		"user-num_sta":                   s.UserNumSta.Val,
		"version":                        s.Version,
		"num_desktop":                    s.NumDesktop.Val,
		"num_handheld":                   s.NumHandheld.Val,
		"num_mobile":                     s.NumMobile.Val,
		"speedtest-status_latency":       s.SpeedtestStatus.Latency.Val,
		"speedtest-status_runtime":       s.SpeedtestStatus.Runtime.Val,
		"speedtest-status_ping":          s.SpeedtestStatus.StatusPing.Val,
		"speedtest-status_xput_download": s.SpeedtestStatus.XputDownload.Val,
		"speedtest-status_xput_upload":   s.SpeedtestStatus.XputUpload.Val,
		"wan1_bytes-r":                   s.Wan1.BytesR.Val,
		"wan1_enable":                    s.Wan1.Enable.Val,
		"wan1_full_duplex":               s.Wan1.FullDuplex.Val,
		"wan1_gateway":                   s.Wan1.Gateway,
		"wan1_ifname":                    s.Wan1.Ifname,
		"wan1_ip":                        s.Wan1.IP,
		"wan1_mac":                       s.Wan1.Mac,
		"wan1_max_speed":                 s.Wan1.MaxSpeed.Val,
		"wan1_name":                      s.Wan1.Name,
		"wan1_rx_bytes":                  s.Wan1.RxBytes.Val,
		"wan1_rx_bytes-r":                s.Wan1.RxBytesR.Val,
		"wan1_rx_dropped":                s.Wan1.RxDropped.Val,
		"wan1_rx_errors":                 s.Wan1.RxErrors.Val,
		"wan1_rx_multicast":              s.Wan1.RxMulticast.Val,
		"wan1_rx_packets":                s.Wan1.RxPackets.Val,
		"wan1_type":                      s.Wan1.Type,
		"wan1_speed":                     s.Wan1.Speed.Val,
		"wan1_up":                        s.Wan1.Up.Val,
		"wan1_tx_bytes":                  s.Wan1.TxBytes.Val,
		"wan1_tx_bytes-r":                s.Wan1.TxBytesR.Val,
		"wan1_tx_dropped":                s.Wan1.TxDropped.Val,
		"wan1_tx_errors":                 s.Wan1.TxErrors.Val,
		"wan1_tx_packets":                s.Wan1.TxPackets.Val,
		"wan2_bytes-r":                   s.Wan2.BytesR.Val,
		"wan2_enable":                    s.Wan2.Enable.Val,
		"wan2_full_duplex":               s.Wan2.FullDuplex.Val,
		"wan2_gateway":                   s.Wan2.Gateway,
		"wan2_ifname":                    s.Wan2.Ifname,
		"wan2_ip":                        s.Wan2.IP,
		"wan2_mac":                       s.Wan2.Mac,
		"wan2_max_speed":                 s.Wan2.MaxSpeed.Val,
		"wan2_name":                      s.Wan2.Name,
		"wan2_rx_bytes":                  s.Wan2.RxBytes.Val,
		"wan2_rx_bytes-r":                s.Wan2.RxBytesR.Val,
		"wan2_rx_dropped":                s.Wan2.RxDropped.Val,
		"wan2_rx_errors":                 s.Wan2.RxErrors.Val,
		"wan2_rx_multicast":              s.Wan2.RxMulticast.Val,
		"wan2_rx_packets":                s.Wan2.RxPackets.Val,
		"wan2_type":                      s.Wan2.Type,
		"wan2_speed":                     s.Wan2.Speed.Val,
		"wan2_up":                        s.Wan2.Up.Val,
		"wan2_tx_bytes":                  s.Wan2.TxBytes.Val,
		"wan2_tx_bytes-r":                s.Wan2.TxBytesR.Val,
		"wan2_tx_dropped":                s.Wan2.TxDropped.Val,
		"wan2_tx_errors":                 s.Wan2.TxErrors.Val,
		"wan2_tx_packets":                s.Wan2.TxPackets.Val,
		"loadavg_1":                      s.SysStats.Loadavg1.Val,
		"loadavg_5":                      s.SysStats.Loadavg5.Val,
		"loadavg_15":                     s.SysStats.Loadavg15.Val,
		"mem_used":                       s.SysStats.MemUsed.Val,
		"mem_buffer":                     s.SysStats.MemBuffer.Val,
		"mem_total":                      s.SysStats.MemTotal.Val,
		"cpu":                            s.SystemStats.CPU.Val,
		"mem":                            s.SystemStats.Mem.Val,
		"system_uptime":                  s.SystemStats.Uptime.Val,
		"lan-rx_bytes":                   s.Stat.LanRxBytes.Val,
		"lan-rx_packets":                 s.Stat.LanRxPackets.Val,
		"lan-tx_bytes":                   s.Stat.LanTxBytes.Val,
		"lan-tx_packets":                 s.Stat.LanTxPackets.Val,
		"wan-rx_bytes":                   s.Stat.WanRxBytes.Val,
		"wan-rx_dropped":                 s.Stat.WanRxDropped.Val,
		"wan-rx_packets":                 s.Stat.WanRxPackets.Val,
		"wan-tx_bytes":                   s.Stat.WanTxBytes.Val,
		"wan-tx_packets":                 s.Stat.WanTxPackets.Val,
	}
	r.send(&metric{Table: "usg", Tags: tags, Fields: fields})

	for _, p := range s.NetworkTable {
		tags := map[string]string{
			"device_name": s.Name,
			"device_id":   s.ID,
			"device_mac":  s.Mac,
			"site_name":   s.SiteName,
			"up":          p.Up.Txt,
			"enabled":     p.Enabled.Txt,
			"site_id":     p.SiteID,
			"ip":          p.IP,
			"ip_subnet":   p.IPSubnet,
			"mac":         p.Mac,
			"name":        p.Name,
			"domain_name": p.DomainName,
			"purpose":     p.Purpose,
		}
		fields := map[string]interface{}{
			"num_sta":    p.NumSta.Val,
			"rx_bytes":   p.RxBytes.Val,
			"rx_packets": p.RxPackets.Val,
			"tx_bytes":   p.TxBytes.Val,
			"tx_packets": p.TxPackets.Val,
		}
		r.send(&metric{Table: "usg_networks", Tags: tags, Fields: fields})
	}
	for _, p := range s.PortTable {
		tags := map[string]string{
			"device_name": s.Name,
			"device_id":   s.ID,
			"device_mac":  s.Mac,
			"site_name":   s.SiteName,
			"name":        p.Name,
			"ifname":      p.Ifname,
			"ip":          p.IP,
			"mac":         p.Mac,
			"up":          p.Up.Txt,
			"speed":       p.Speed.Txt,
			"full_duplex": p.FullDuplex.Txt,
			"enable":      p.Enable.Txt,
		}
		fields := map[string]interface{}{
			"rx_bytes":     p.RxBytes.Val,
			"rx_dropped":   p.RxDropped.Val,
			"rx_errors":    p.RxErrors.Val,
			"rx_packets":   p.RxBytes.Val,
			"tx_bytes":     p.TxBytes.Val,
			"tx_dropped":   p.TxDropped.Val,
			"tx_errors":    p.TxErrors.Val,
			"tx_packets":   p.TxPackets.Val,
			"rx_multicast": p.RxMulticast.Val,
			"dns_servers":  strings.Join(p.DNS, ","),
		}
		r.send(&metric{Table: "usg_ports", Tags: tags, Fields: fields})

	}
}
