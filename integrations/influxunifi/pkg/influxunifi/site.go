package influxunifi

import (
	"time"

	influx "github.com/influxdata/influxdb1-client/v2"
	"golift.io/unifi"
)

// SitePoints generates Unifi Sites' datapoints for InfluxDB.
// These points can be passed directly to influx.
func SitePoints(u *unifi.Site, now time.Time) ([]*influx.Point, error) {
	points := []*influx.Point{}
	for _, s := range u.Health {
		tags := map[string]string{
			"name":      u.Name,
			"site_name": u.SiteName,
			"desc":      u.Desc,
			"status":    s.Status,
			"subsystem": s.Subsystem,
			"wan_ip":    s.WanIP,
			"netmask":   s.Netmask,
			"gw_name":   s.GwName,
			"gw_mac":    s.GwMac,
			"lan_ip":    s.LanIP,
		}
		fields := map[string]interface{}{
			"num_user":                 s.NumUser.Val,
			"num_guest":                s.NumGuest.Val,
			"num_iot":                  s.NumIot.Val,
			"tx_bytes-r":               s.TxBytesR.Val,
			"rx_bytes-r":               s.RxBytesR.Val,
			"num_ap":                   s.NumAp.Val,
			"num_adopted":              s.NumAdopted.Val,
			"num_disabled":             s.NumDisabled.Val,
			"num_disconnected":         s.NumDisconnected.Val,
			"num_pending":              s.NumPending.Val,
			"num_gw":                   s.NumGw.Val,
			"wan_ip":                   s.WanIP,
			"num_sta":                  s.NumSta.Val,
			"gw_cpu":                   s.GwSystemStats.CPU.Val,
			"gw_mem":                   s.GwSystemStats.Mem.Val,
			"gw_uptime":                s.GwSystemStats.Uptime.Val,
			"latency":                  s.Latency.Val,
			"uptime":                   s.Uptime.Val,
			"drops":                    s.Drops.Val,
			"xput_up":                  s.XputUp.Val,
			"xput_down":                s.XputDown.Val,
			"speedtest_ping":           s.SpeedtestPing.Val,
			"speedtest_lastrun":        s.SpeedtestLastrun.Val,
			"num_sw":                   s.NumSw.Val,
			"remote_user_num_active":   s.RemoteUserNumActive.Val,
			"remote_user_num_inactive": s.RemoteUserNumInactive.Val,
			"remote_user_rx_bytes":     s.RemoteUserRxBytes.Val,
			"remote_user_tx_bytes":     s.RemoteUserTxBytes.Val,
			"remote_user_rx_packets":   s.RemoteUserRxPackets.Val,
			"remote_user_tx_packets":   s.RemoteUserTxPackets.Val,
			"num_new_alarms":           u.NumNewAlarms.Val,
		}
		pt, err := influx.NewPoint("subsystems", tags, fields, time.Now())
		if err != nil {
			return points, err
		}
		points = append(points, pt)
	}
	return points, nil
}
