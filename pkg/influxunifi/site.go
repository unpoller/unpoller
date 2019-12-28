package influxunifi

import (
	"golift.io/unifi"
)

// batchSite generates Unifi Sites' datapoints for InfluxDB.
// These points can be passed directly to influx.
func (u *InfluxUnifi) batchSite(r report, s *unifi.Site) {
	for _, h := range s.Health {
		tags := map[string]string{
			"name":      s.Name,
			"site_name": s.SiteName,
			"desc":      s.Desc,
			"status":    h.Status,
			"subsystem": h.Subsystem,
			"wan_ip":    h.WanIP,
			"gw_name":   h.GwName,
			"lan_ip":    h.LanIP,
		}
		fields := map[string]interface{}{
			"num_user":                 h.NumUser.Val,
			"num_guest":                h.NumGuest.Val,
			"num_iot":                  h.NumIot.Val,
			"tx_bytes-r":               h.TxBytesR.Val,
			"rx_bytes-r":               h.RxBytesR.Val,
			"num_ap":                   h.NumAp.Val,
			"num_adopted":              h.NumAdopted.Val,
			"num_disabled":             h.NumDisabled.Val,
			"num_disconnected":         h.NumDisconnected.Val,
			"num_pending":              h.NumPending.Val,
			"num_gw":                   h.NumGw.Val,
			"wan_ip":                   h.WanIP,
			"num_sta":                  h.NumSta.Val,
			"gw_cpu":                   h.GwSystemStats.CPU.Val,
			"gw_mem":                   h.GwSystemStats.Mem.Val,
			"gw_uptime":                h.GwSystemStats.Uptime.Val,
			"latency":                  h.Latency.Val,
			"uptime":                   h.Uptime.Val,
			"drops":                    h.Drops.Val,
			"xput_up":                  h.XputUp.Val,
			"xput_down":                h.XputDown.Val,
			"speedtest_ping":           h.SpeedtestPing.Val,
			"speedtest_lastrun":        h.SpeedtestLastrun.Val,
			"num_sw":                   h.NumSw.Val,
			"remote_user_num_active":   h.RemoteUserNumActive.Val,
			"remote_user_num_inactive": h.RemoteUserNumInactive.Val,
			"remote_user_rx_bytes":     h.RemoteUserRxBytes.Val,
			"remote_user_tx_bytes":     h.RemoteUserTxBytes.Val,
			"remote_user_rx_packets":   h.RemoteUserRxPackets.Val,
			"remote_user_tx_packets":   h.RemoteUserTxPackets.Val,
			"num_new_alarms":           s.NumNewAlarms.Val,
		}
		r.send(&metric{Table: "subsystems", Tags: tags, Fields: fields})
	}
}
