package promunifi

import (
	"github.com/prometheus/client_golang/prometheus"
	"golift.io/unifi"
)

// These are shared by all four device types: UDM, UAP, USG, USW
type unifiDevice struct {
	Info          *prometheus.Desc // uptime
	Temperature   *prometheus.Desc // sw only
	TotalMaxPower *prometheus.Desc // sw only
	FanLevel      *prometheus.Desc // sw only
	TotalTxBytes  *prometheus.Desc
	TotalRxBytes  *prometheus.Desc
	TotalBytes    *prometheus.Desc
	BytesR        *prometheus.Desc // ap only
	BytesD        *prometheus.Desc // ap only
	Bytes         *prometheus.Desc // ap only
	TxBytesD      *prometheus.Desc // ap only
	RxBytesD      *prometheus.Desc // ap only
	NumSta        *prometheus.Desc
	NumDesktop    *prometheus.Desc // gw only
	NumMobile     *prometheus.Desc // gw only
	NumHandheld   *prometheus.Desc // gw only
	Loadavg1      *prometheus.Desc
	Loadavg5      *prometheus.Desc
	Loadavg15     *prometheus.Desc
	MemBuffer     *prometheus.Desc
	MemTotal      *prometheus.Desc
	MemUsed       *prometheus.Desc
	CPU           *prometheus.Desc
	Mem           *prometheus.Desc
}

func descDevice(ns string) *unifiDevice {
	labels := []string{"type", "site_name", "name"}
	infoLabels := []string{"version", "model", "serial", "mac", "ip", "id", "bytes"}
	return &unifiDevice{
		Info:          prometheus.NewDesc(ns+"info", "Device Information", append(labels, infoLabels...), nil),
		Temperature:   prometheus.NewDesc(ns+"temperature_celsius", "Temperature", labels, nil),
		TotalMaxPower: prometheus.NewDesc(ns+"max_power_total", "Total Max Power", labels, nil),
		FanLevel:      prometheus.NewDesc(ns+"fan_level", "Fan Level", labels, nil),
		TotalTxBytes:  prometheus.NewDesc(ns+"transmit_bytes_total", "Total Transmitted Bytes", labels, nil),
		TotalRxBytes:  prometheus.NewDesc(ns+"receive_bytes_total", "Total Received Bytes", labels, nil),
		TotalBytes:    prometheus.NewDesc(ns+"bytes_total", "Total Bytes Transferred", labels, nil),
		BytesR:        prometheus.NewDesc(ns+"rate_bytes", "Transfer Rate", labels, nil),
		BytesD:        prometheus.NewDesc(ns+"d_bytes", "Total Bytes D???", labels, nil),
		Bytes:         prometheus.NewDesc(ns+"transferred_bytes_total", "Bytes Transferred", labels, nil),
		TxBytesD:      prometheus.NewDesc(ns+"d_tranmsit_bytes", "Transmit Bytes D???", labels, nil),
		RxBytesD:      prometheus.NewDesc(ns+"d_receive_bytes", "Receive Bytes D???", labels, nil),
		NumSta:        prometheus.NewDesc(ns+"stations", "Number of Stations", append(labels, "station_type"), nil),
		NumDesktop:    prometheus.NewDesc(ns+"desktops", "Number of Desktops", labels, nil),
		NumMobile:     prometheus.NewDesc(ns+"mobile", "Number of Mobiles", labels, nil),
		NumHandheld:   prometheus.NewDesc(ns+"handheld", "Number of Handhelds", labels, nil),
		Loadavg1:      prometheus.NewDesc(ns+"load_average_1", "System Load Average 1 Minute", labels, nil),
		Loadavg5:      prometheus.NewDesc(ns+"load_average_5", "System Load Average 5 Minutes", labels, nil),
		Loadavg15:     prometheus.NewDesc(ns+"load_average_15", "System Load Average 15 Minutes", labels, nil),
		MemUsed:       prometheus.NewDesc(ns+"memory_used_bytes", "System Memory Used", labels, nil),
		MemTotal:      prometheus.NewDesc(ns+"memory_installed_bytes", "System Installed Memory", labels, nil),
		MemBuffer:     prometheus.NewDesc(ns+"memory_buffer_bytes", "System Memory Buffer", labels, nil),
		CPU:           prometheus.NewDesc(ns+"cpu_utilization_ratio", "System CPU % Utilized", labels, nil),
		Mem:           prometheus.NewDesc(ns+"memory_utilization_ratio", "System Memory % Utilized", labels, nil),
	}
}

// UDM is a collection of stats from USG, USW and UAP. It has no unique stats.
func (u *promUnifi) exportUDM(r report, d *unifi.UDM) {
	labels := []string{d.Type, d.SiteName, d.Name}
	infoLabels := []string{d.Version, d.Model, d.Serial, d.Mac, d.IP, d.ID, d.Bytes.Txt}
	labelsGuest := append(labels, "guest")
	labelsUser := append(labels, "user")
	// Dream Machine System Data.
	r.send([]*metric{
		{u.Device.Info, prometheus.GaugeValue, d.Uptime, append(labels, infoLabels...)},
		{u.Device.TotalTxBytes, prometheus.CounterValue, d.TxBytes, labels},
		{u.Device.TotalRxBytes, prometheus.CounterValue, d.RxBytes, labels},
		{u.Device.TotalBytes, prometheus.CounterValue, d.Bytes, labels},
		{u.Device.NumSta, prometheus.GaugeValue, d.UserNumSta, labelsUser},
		{u.Device.NumSta, prometheus.GaugeValue, d.GuestNumSta, labelsGuest},
		{u.Device.NumDesktop, prometheus.GaugeValue, d.NumDesktop, labels},
		{u.Device.NumMobile, prometheus.GaugeValue, d.NumMobile, labels},
		{u.Device.NumHandheld, prometheus.GaugeValue, d.NumHandheld, labels},
		{u.Device.Loadavg1, prometheus.GaugeValue, d.SysStats.Loadavg1, labels},
		{u.Device.Loadavg5, prometheus.GaugeValue, d.SysStats.Loadavg5, labels},
		{u.Device.Loadavg15, prometheus.GaugeValue, d.SysStats.Loadavg15, labels},
		{u.Device.MemUsed, prometheus.GaugeValue, d.SysStats.MemUsed, labels},
		{u.Device.MemTotal, prometheus.GaugeValue, d.SysStats.MemTotal, labels},
		{u.Device.MemBuffer, prometheus.GaugeValue, d.SysStats.MemBuffer, labels},
		{u.Device.CPU, prometheus.GaugeValue, d.SystemStats.CPU.Val / 100.0, labels},
		{u.Device.Mem, prometheus.GaugeValue, d.SystemStats.Mem.Val / 100.0, labels},
	})

	// Switch Data
	u.exportUSWstats(r, labels, d.Stat.Sw)
	u.exportPortTable(r, labels, d.PortTable)
	// Gateway Data
	u.exportUSGstats(r, labels, d.Stat.Gw, d.SpeedtestStatus, d.Uplink)
	u.exportWANPorts(r, labels, d.Wan1, d.Wan2)
	// Wireless Data - UDM (non-pro) only
	if d.Stat.Ap != nil && d.VapTable != nil {
		u.exportUAPstats(r, labels, d.Stat.Ap)
		u.exportVAPtable(r, labels, *d.VapTable)
		u.exportRadtable(r, labels, *d.RadioTable, *d.RadioTableStats)
	}
}
