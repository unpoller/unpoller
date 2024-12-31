package influxunifi

import (
	"github.com/unpoller/unifi/v5"
)

// batchSite generates Unifi Sites' datapoints for InfluxDB.
// These points can be passed directly to influx.
func (u *InfluxUnifi) batchSite(r report, s *unifi.Site) {
	for _, h := range s.Health {
		tags := map[string]string{
			"name":      s.Name,
			"site_name": s.SiteName,
			"source":    s.SourceName,
			"desc":      s.Desc,
			"status":    h.Status,
			"subsystem": h.Subsystem,
			"wan_ip":    h.WanIP,
			"gw_name":   h.GwName,
			"lan_ip":    h.LanIP,
		}
		fields := map[string]any{
			"num_user":                 h.NumUser.Val,
			"num_guest":                h.NumGuest.Val,
			"num_iot":                  h.NumIot.Val,
			"tx_bytes-r":               h.TxBytesR.Int64(),
			"rx_bytes-r":               h.RxBytesR.Int64(),
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
			"uptime":                   h.Uptime.Int64(),
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

func (u *InfluxUnifi) batchSiteDPI(r report, v any) {
	s, ok := v.(*unifi.DPITable)
	if !ok {
		u.LogErrorf("invalid type given to batchSiteDPI: %T", v)

		return
	}

	for _, dpi := range s.ByApp {
		r.send(&metric{
			Table: "sitedpi",
			Tags: map[string]string{
				"category":    unifi.DPICats.Get(dpi.Cat.Int()),
				"application": unifi.DPIApps.GetApp(dpi.Cat.Int(), dpi.App.Int()),
				"site_name":   s.SiteName,
				"source":      s.SourceName,
			},
			Fields: map[string]any{
				"tx_packets": dpi.TxPackets.Int64(),
				"rx_packets": dpi.RxPackets.Int64(),
				"tx_bytes":   dpi.TxBytes.Int64(),
				"rx_bytes":   dpi.RxBytes.Int64(),
			},
		})
	}
}
