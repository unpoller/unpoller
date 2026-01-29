package promunifi

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/unpoller/unifi/v5"
)

// These are shared by all four device types: UDM, UAP, USG, USW.
type unifiDevice struct {
	Info                     *prometheus.Desc
	Uptime                   *prometheus.Desc
	Temperature              *prometheus.Desc
	Storage                  *prometheus.Desc
	TotalMaxPower            *prometheus.Desc // sw only
	OutletACPowerConsumption *prometheus.Desc // pdu only
	PowerSource              *prometheus.Desc // pdu only
	FanLevel                 *prometheus.Desc // sw only
	TotalTxBytes             *prometheus.Desc
	TotalRxBytes             *prometheus.Desc
	TotalBytes               *prometheus.Desc
	BytesR                   *prometheus.Desc // ap only
	BytesD                   *prometheus.Desc // ap only
	TxBytesD                 *prometheus.Desc // ap only
	RxBytesD                 *prometheus.Desc // ap only
	Counter                  *prometheus.Desc
	Loadavg1                 *prometheus.Desc
	Loadavg5                 *prometheus.Desc
	Loadavg15                *prometheus.Desc
	MemBuffer                *prometheus.Desc
	MemTotal                 *prometheus.Desc
	MemUsed                  *prometheus.Desc
	CPU                      *prometheus.Desc
	Mem                      *prometheus.Desc
	Upgradeable              *prometheus.Desc
}

func descDevice(ns string) *unifiDevice {
	labels := []string{"type", "site_name", "name", "source"}
	infoLabels := []string{"version", "model", "serial", "mac", "ip", "id", "tag"}

	return &unifiDevice{
		Info:   prometheus.NewDesc(ns+"info", "Device Information", append(labels, infoLabels...), nil),
		Uptime: prometheus.NewDesc(ns+"uptime_seconds", "Device Uptime", append(labels, "tag"), nil),
		Temperature: prometheus.NewDesc(ns+"temperature_celsius", "Temperature",
			append(labels, "temp_area", "temp_type", "tag"), nil),
		Storage: prometheus.NewDesc(ns+"storage", "Storage",
			append(labels, "mountpoint", "storage_name", "storage_reading", "tag"), nil),
		TotalMaxPower:            prometheus.NewDesc(ns+"max_power_total", "Total Max Power", append(labels, "tag"), nil),
		OutletACPowerConsumption: prometheus.NewDesc(ns+"outlet_ac_power_consumption", "Outlet AC Power Consumption", append(labels, "tag"), nil),
		PowerSource:              prometheus.NewDesc(ns+"power_source", "Power Source", append(labels, "tag"), nil),
		FanLevel:                 prometheus.NewDesc(ns+"fan_level", "Fan Level", append(labels, "tag"), nil),
		TotalTxBytes:             prometheus.NewDesc(ns+"transmit_bytes_total", "Total Transmitted Bytes", append(labels, "tag"), nil),
		TotalRxBytes:             prometheus.NewDesc(ns+"receive_bytes_total", "Total Received Bytes", append(labels, "tag"), nil),
		TotalBytes:               prometheus.NewDesc(ns+"bytes_total", "Total Bytes Transferred", append(labels, "tag"), nil),
		BytesR:                   prometheus.NewDesc(ns+"rate_bytes", "Transfer Rate", append(labels, "tag"), nil),
		BytesD:                   prometheus.NewDesc(ns+"d_bytes", "Total Bytes D???", append(labels, "tag"), nil),
		TxBytesD:                 prometheus.NewDesc(ns+"d_tranmsit_bytes", "Transmit Bytes D???", append(labels, "tag"), nil),
		RxBytesD:                 prometheus.NewDesc(ns+"d_receive_bytes", "Receive Bytes D???", append(labels, "tag"), nil),
		Counter:                  prometheus.NewDesc(ns+"stations", "Number of Stations", append(labels, "station_type", "tag"), nil),
		Loadavg1:                 prometheus.NewDesc(ns+"load_average_1", "System Load Average 1 Minute", append(labels, "tag"), nil),
		Loadavg5:                 prometheus.NewDesc(ns+"load_average_5", "System Load Average 5 Minutes", append(labels, "tag"), nil),
		Loadavg15:                prometheus.NewDesc(ns+"load_average_15", "System Load Average 15 Minutes", append(labels, "tag"), nil),
		MemUsed:                  prometheus.NewDesc(ns+"memory_used_bytes", "System Memory Used", append(labels, "tag"), nil),
		MemTotal:                 prometheus.NewDesc(ns+"memory_installed_bytes", "System Installed Memory", append(labels, "tag"), nil),
		MemBuffer:                prometheus.NewDesc(ns+"memory_buffer_bytes", "System Memory Buffer", append(labels, "tag"), nil),
		CPU:                      prometheus.NewDesc(ns+"cpu_utilization_ratio", "System CPU % Utilized", append(labels, "tag"), nil),
		Mem:                      prometheus.NewDesc(ns+"memory_utilization_ratio", "System Memory % Utilized", append(labels, "tag"), nil),
		Upgradeable:              prometheus.NewDesc(ns+"upgradable", "Upgrade-able", append(labels, "tag"), nil),
	}
}

// UDM is a collection of stats from USG, USW and UAP. It has no unique stats.
func (u *promUnifi) exportUDM(r report, d *unifi.UDM) {
	if !d.Adopted.Val || d.Locating.Val {
		return
	}

	baseLabels := []string{d.Type, d.SiteName, d.Name, d.SourceName}
	baseInfoLabels := []string{d.Version, d.Model, d.Serial, d.Mac, d.IP, d.ID}

	// Export metrics with tags - create separate series for each tag
	u.exportWithTags(r, d.Tags, func(tagLabels []string) {
		tag := tagLabels[0]
		labels := baseLabels

		infoLabels := append(baseInfoLabels, tag)

		// Shared data (all devices do this).
		u.exportBYTstats(r, append(labels, tag), d.TxBytes, d.RxBytes)
		u.exportSYSstats(r, append(labels, tag), d.SysStats, d.SystemStats)
		u.exportSTAcount(r, append(labels, tag), d.UserNumSta, d.GuestNumSta, d.NumDesktop, d.NumMobile, d.NumHandheld)
		// Switch Data
		u.exportUSWstats(r, append(labels, tag), d.Stat.Sw)
		u.exportPRTtable(r, append(labels, tag), d.PortTable)
		// Gateway Data
		u.exportWANPorts(r, append(labels, tag), d.Wan1, d.Wan2)
		u.exportUSGstats(r, append(labels, tag), d.Stat.Gw, d.SpeedtestStatus, d.Uplink)
		// Dream Machine System Data.
		r.send([]*metric{
			{u.Device.Info, gauge, 1.0, append(labels, infoLabels...)},
			{u.Device.Uptime, gauge, d.Uptime, append(labels, tag)},
			{u.Device.Upgradeable, gauge, d.Upgradeable.Val, append(labels, tag)},
		})

		// UDM pro has special temp sensors. UDM non-pro may not have temp; not sure.
		for _, t := range d.Temperatures {
			r.send([]*metric{{u.Device.Temperature, gauge, t.Value, append(labels, t.Name, t.Type, tag)}})
		}

		// UDM pro and UXG have hard drives.
		for _, t := range d.Storage {
			r.send([]*metric{
				{u.Device.Storage, gauge, t.Size.Val, append(labels, t.MountPoint, t.Name, "size", tag)},
				{u.Device.Storage, gauge, t.Used.Val, append(labels, t.MountPoint, t.Name, "used", tag)},
			})
		}
	})

	// Wireless Data - UDM (non-pro) only
	if d.Stat.Ap != nil && d.VapTable != nil {
		u.exportWithTags(r, d.Tags, func(tagLabels []string) {
			tag := tagLabels[0]
			labels := append(baseLabels, tag)
			u.exportUAPstats(r, labels, d.Stat.Ap, d.BytesD, d.TxBytesD, d.RxBytesD, d.BytesR)
			u.exportVAPtable(r, labels, *d.VapTable)
			u.exportRADtable(r, labels, *d.RadioTable, *d.RadioTableStats)
		})
	}
}

// exportWithTags exports metrics with tag support. If device has multiple tags,
// each tag creates a separate metric series. If no tags, exports with tag="".
func (u *promUnifi) exportWithTags(_ report, tags []string, fn func([]string)) {
	if len(tags) == 0 {
		// No tags - export once with empty tag
		fn([]string{""})

		return
	}
	// Multiple tags - export once per tag
	for _, tag := range tags {
		fn([]string{tag})
	}
}

// Shared by all.
func (u *promUnifi) exportBYTstats(r report, labels []string, tx, rx unifi.FlexInt) {
	r.send([]*metric{
		{u.Device.TotalTxBytes, counter, tx, labels},
		{u.Device.TotalRxBytes, counter, rx, labels},
		{u.Device.TotalBytes, counter, tx.Val + rx.Val, labels},
	})
}

// Shared by all, pass 2 or 5 stats.
func (u *promUnifi) exportSTAcount(r report, labels []string, stas ...unifi.FlexInt) {
	r.send([]*metric{
		{u.Device.Counter, gauge, stas[0], append(labels, "user")},
		{u.Device.Counter, gauge, stas[1], append(labels, "guest")},
	})

	if len(stas) > 2 { // nolint: gomnd
		r.send([]*metric{
			{u.Device.Counter, gauge, stas[2], append(labels, "desktop")},
			{u.Device.Counter, gauge, stas[3], append(labels, "mobile")},
			{u.Device.Counter, gauge, stas[4], append(labels, "handheld")},
		})
	}
}

// Shared by all.
func (u *promUnifi) exportSYSstats(r report, labels []string, s unifi.SysStats, ss unifi.SystemStats) {
	r.send([]*metric{
		{u.Device.Loadavg1, gauge, s.Loadavg1, labels},
		{u.Device.Loadavg5, gauge, s.Loadavg5, labels},
		{u.Device.Loadavg15, gauge, s.Loadavg15, labels},
		{u.Device.MemUsed, gauge, s.MemUsed, labels},
		{u.Device.MemTotal, gauge, s.MemTotal, labels},
		{u.Device.MemBuffer, gauge, s.MemBuffer, labels},
		{u.Device.CPU, gauge, ss.CPU.Val / 100.0, labels},
		{u.Device.Mem, gauge, ss.Mem.Val / 100.0, labels},
	})
}
