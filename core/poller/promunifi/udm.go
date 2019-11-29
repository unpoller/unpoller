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
func (u *unifiCollector) exportUDM(r *Report, s *unifi.UDM) {
	labels := []string{s.IP, s.Type, s.Version, s.SiteName, s.Mac, s.Model, s.Name, s.Serial}
	// Gateway System Data.
	r.send([]*metricExports{
		{u.USG.Uptime, prometheus.GaugeValue, s.Uptime, labels},
		{u.USG.TotalTxBytes, prometheus.CounterValue, s.TxBytes, labels},
		{u.USG.TotalRxBytes, prometheus.CounterValue, s.RxBytes, labels},
		{u.USG.TotalBytes, prometheus.CounterValue, s.Bytes, labels},
		{u.USG.NumSta, prometheus.GaugeValue, s.NumSta, labels},
		{u.USG.UserNumSta, prometheus.GaugeValue, s.UserNumSta, labels},
		{u.USG.GuestNumSta, prometheus.GaugeValue, s.GuestNumSta, labels},
		{u.USG.NumDesktop, prometheus.GaugeValue, s.NumDesktop, labels},
		{u.USG.NumMobile, prometheus.GaugeValue, s.NumMobile, labels},
		{u.USG.NumHandheld, prometheus.GaugeValue, s.NumHandheld, labels},
		{u.USG.Loadavg1, prometheus.GaugeValue, s.SysStats.Loadavg1, labels},
		{u.USG.Loadavg5, prometheus.GaugeValue, s.SysStats.Loadavg5, labels},
		{u.USG.Loadavg15, prometheus.GaugeValue, s.SysStats.Loadavg15, labels},
		{u.USG.MemUsed, prometheus.GaugeValue, s.SysStats.MemUsed, labels},
		{u.USG.MemTotal, prometheus.GaugeValue, s.SysStats.MemTotal, labels},
		{u.USG.MemBuffer, prometheus.GaugeValue, s.SysStats.MemBuffer, labels},
		{u.USG.CPU, prometheus.GaugeValue, s.SystemStats.CPU, labels},
		{u.USG.Mem, prometheus.GaugeValue, s.SystemStats.Mem, labels},
	})
	u.exportUSWstats(r, s.Stat.Sw, labels)
	u.exportUSGstats(r, s.Stat.Gw, s.SpeedtestStatus, labels)
	u.exportWANPorts(r, labels, s.Wan1, s.Wan2)
	u.exportPortTable(r, s.PortTable, labels[4:])
	if s.Stat.Ap != nil && s.VapTable != nil {
		u.exportUAPstats(r, labels[2:], s.Stat.Ap)
		u.exportVAPtable(r, labels[2:], *s.VapTable, *s.RadioTable, *s.RadioTableStats)
	}
}
