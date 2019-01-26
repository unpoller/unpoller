package unifi

import (
	"strconv"
	"time"

	influx "github.com/influxdata/influxdb1-client/v2"
)

// Points generates Unifi Switch datapoints for InfluxDB.
// These points can be passed directly to influx.
func (u USW) Points() ([]*influx.Point, error) {
	var points []*influx.Point
	tags := map[string]string{
		"id":                     u.ID,
		"mac":                    u.Mac,
		"device_type":            u.Stat.O,
		"device_oid":             u.Stat.Oid,
		"site_id":                u.SiteID,
		"name":                   u.Name,
		"addopted":               strconv.FormatBool(u.Adopted),
		"adopt_ip":               u.AdoptIP,
		"adopt_url":              u.AdoptURL,
		"cfgversion":             u.Cfgversion,
		"config_network_ip":      u.ConfigNetwork.IP,
		"config_network_type":    u.ConfigNetwork.Type,
		"connect_request_ip":     u.ConnectRequestIP,
		"connect_request_port":   u.ConnectRequestPort,
		"default":                strconv.FormatBool(u.Default),
		"device_id":              u.DeviceID,
		"discovered_via":         u.DiscoveredVia,
		"inform_ip":              u.InformIP,
		"last_uplink_mac":        u.LastUplink.UplinkMac,
		"known_cfgversion":       u.KnownCfgversion,
		"led_override":           u.LedOverride,
		"locating":               strconv.FormatBool(u.Locating),
		"model":                  u.Model,
		"outdoor_mode_override":  u.OutdoorModeOverride,
		"serial":                 u.Serial,
		"type":                   u.Type,
		"version_incompatible":   strconv.FormatBool(u.VersionIncompatible),
		"dot1x_portctrl_enabled": strconv.FormatBool(u.Dot1XPortctrlEnabled),
		"flowctrl_enabled":       strconv.FormatBool(u.FlowctrlEnabled),
		"has_fan":                strconv.FormatBool(u.HasFan),
		"has_temperature":        strconv.FormatBool(u.HasTemperature),
		"jumboframe_enabled":     strconv.FormatBool(u.JumboframeEnabled),
		"stp_priority":           u.StpPriority,
		"stp_version":            u.StpVersion,
	}
	fields := map[string]interface{}{
		"fw_caps":                  u.FwCaps,
		"guest-num_sta":            u.GuestNumSta,
		"ip":                       u.IP,
		"bytes":                    u.Bytes,
		"fan_level":                u.FanLevel,
		"general_temperature":      u.GeneralTemperature,
		"last_seen":                u.LastSeen,
		"license_state":            u.LicenseState,
		"overheating":              u.Overheating,
		"rx_bytes":                 u.RxBytes,
		"tx_bytes":                 u.TxBytes,
		"uptime":                   u.Uptime,
		"considered_lost_at":       u.ConsideredLostAt,
		"next_heartbeat_at":        u.NextHeartbeatAt,
		"roll_upgrade":             u.Rollupgrade,
		"state":                    u.State,
		"upgradable":               u.Upgradable,
		"user-num_sta":             u.UserNumSta,
		"version":                  u.Version,
		"loadavg_1":                u.SysStats.Loadavg1,
		"loadavg_5":                u.SysStats.Loadavg5,
		"loadavg_15":               u.SysStats.Loadavg15,
		"mem_buffer":               u.SysStats.MemBuffer,
		"mem_used":                 u.SysStats.MemUsed,
		"mem_total":                u.SysStats.MemTotal,
		"cpu":                      u.SystemStats.CPU,
		"mem":                      u.SystemStats.Mem,
		"system_uptime":            u.SystemStats.Uptime,
		"stat_bytes":               u.Stat.Bytes,
		"stat_duration":            u.Stat.Duration,
		"stat_guest-rx_bytes":      u.Stat.RxBytes,
		"stat_guest-rx_crypts":     u.Stat.RxCrypts,
		"stat_guest-rx_dropped":    u.Stat.RxDropped,
		"stat_guest-rx_errors":     u.Stat.RxErrors,
		"stat_guest-rx_frags":      u.Stat.RxFrags,
		"stat_guest-rx_packets":    u.Stat.RxPackets,
		"stat_guest-tx_bytes":      u.Stat.TxBytes,
		"stat_guest-tx_dropped":    u.Stat.TxDropped,
		"stat_guest-tx_errors":     u.Stat.TxErrors,
		"stat_guest-tx_packets":    u.Stat.TxPackets,
		"stat_guest-tx_retries":    u.Stat.TxRetries,
		"stat_port_1-rx_broadcast": u.Stat.Port1RxBroadcast,
		"stat_port_1-rx_bytes":     u.Stat.Port1RxBytes,
		"stat_port_1-rx_multicast": u.Stat.Port1RxMulticast,
		"stat_port_1-rx_packets":   u.Stat.Port1RxPackets,
		"stat_port_1-tx_broadcast": u.Stat.Port1TxBroadcast,
		"stat_port_1-tx_bytes":     u.Stat.Port1TxBytes,
		"stat_port_1-tx_multicast": u.Stat.Port1TxMulticast,
		"stat_port_1-tx_packets":   u.Stat.Port1TxPackets,
		"stat_rx_bytes":            u.Stat.RxBytes,
		"stat_rx_crypts":           u.Stat.RxCrypts,
		"stat_rx_dropped":          u.Stat.RxDropped,
		"stat_rx_errors":           u.Stat.RxErrors,
		"stat_rx_frags":            u.Stat.RxFrags,
		"stat_rx_packets":          u.Stat.TxPackets,
		"stat_tx_bytes":            u.Stat.TxBytes,
		"stat_tx_dropped":          u.Stat.TxDropped,
		"stat_tx_errors":           u.Stat.TxErrors,
		"stat_tx_packets":          u.Stat.TxPackets,
		"stat_tx_retries":          u.Stat.TxRetries,
		"uplink_depth":             strconv.FormatFloat(u.UplinkDepth, 'f', 6, 64),
		// Add the port stats too.
	}
	pt, err := influx.NewPoint("usw", tags, fields, time.Now())
	if err == nil {
		points = append(points, pt)
	}
	return points, err
}
