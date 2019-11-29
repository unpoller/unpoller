package promunifi

import (
	"github.com/prometheus/client_golang/prometheus"
	"golift.io/unifi"
)

func (u *unifiCollector) exportUDMs(r *Report) {
	if r.Metrics == nil || r.Metrics.Devices == nil || len(r.Metrics.Devices.UDMs) < 1 {
		return
	}
	r.wg.Add(one)
	go func() {
		defer r.wg.Done()
		for _, d := range r.Metrics.Devices.UDMs {
			u.exportUDM(r, d)
		}
	}()
}

// UDM is a collection of stats from USG, USW and UAP. It has no unique stats.
func (u *unifiCollector) exportUDM(r *Report, d *unifi.UDM) {
	labels := []string{d.IP, d.Type, d.Version, d.SiteName, d.Mac, d.Model, d.Name, d.Serial}
	// Gateway System Data.
	r.send([]*metricExports{
		{u.USG.Uptime, prometheus.GaugeValue, d.Uptime, labels},
		{u.USG.TotalTxBytes, prometheus.CounterValue, d.TxBytes, labels},
		{u.USG.TotalRxBytes, prometheus.CounterValue, d.RxBytes, labels},
		{u.USG.TotalBytes, prometheus.CounterValue, d.Bytes, labels},
		{u.USG.NumSta, prometheus.GaugeValue, d.NumSta, labels},
		{u.USG.UserNumSta, prometheus.GaugeValue, d.UserNumSta, labels},
		{u.USG.GuestNumSta, prometheus.GaugeValue, d.GuestNumSta, labels},
		{u.USG.NumDesktop, prometheus.GaugeValue, d.NumDesktop, labels},
		{u.USG.NumMobile, prometheus.GaugeValue, d.NumMobile, labels},
		{u.USG.NumHandheld, prometheus.GaugeValue, d.NumHandheld, labels},
		{u.USG.Loadavg1, prometheus.GaugeValue, d.SysStats.Loadavg1, labels},
		{u.USG.Loadavg5, prometheus.GaugeValue, d.SysStats.Loadavg5, labels},
		{u.USG.Loadavg15, prometheus.GaugeValue, d.SysStats.Loadavg15, labels},
		{u.USG.MemUsed, prometheus.GaugeValue, d.SysStats.MemUsed, labels},
		{u.USG.MemTotal, prometheus.GaugeValue, d.SysStats.MemTotal, labels},
		{u.USG.MemBuffer, prometheus.GaugeValue, d.SysStats.MemBuffer, labels},
		{u.USG.CPU, prometheus.GaugeValue, d.SystemStats.CPU, labels},
		{u.USG.Mem, prometheus.GaugeValue, d.SystemStats.Mem, labels},
	})
	u.exportUSWstats(r, d.Stat.Sw, labels)
	u.exportUSGstats(r, d.Stat.Gw, d.SpeedtestStatus, labels)
	u.exportWANPorts(r, labels, d.Wan1, d.Wan2)
	u.exportPortTable(r, d.PortTable, labels[4:])
	if d.Stat.Ap != nil && d.VapTable != nil {
		u.exportUAPstats(r, labels[2:], d.Stat.Ap)
		u.exportVAPtable(r, labels[2:], *d.VapTable)
		u.exportRadtable(r, labels[2:], *d.RadioTable, *d.RadioTableStats)
	}
}
