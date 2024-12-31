package promunifi

import (
	"github.com/unpoller/unifi/v5"
)

// exportUXG is a collection of stats from USG and USW. It has no unique stats.
func (u *promUnifi) exportUXG(r report, d *unifi.UXG) {
	if !d.Adopted.Val || d.Locating.Val {
		return
	}

	var gw *unifi.Gw
	if d.Stat != nil {
		gw = d.Stat.Gw
	}

	var sw *unifi.Sw
	if d.Stat != nil {
		sw = d.Stat.Sw
	}

	labels := []string{d.Type, d.SiteName, d.Name, d.SourceName}
	infoLabels := []string{d.Version, d.Model, d.Serial, d.Mac, d.IP, d.ID}
	// Shared data (all devices do this).
	u.exportBYTstats(r, labels, d.TxBytes, d.RxBytes)
	u.exportSYSstats(r, labels, d.SysStats, d.SystemStats)
	u.exportSTAcount(r, labels, d.UserNumSta, d.GuestNumSta, d.NumDesktop, d.NumMobile, d.NumHandheld)
	// Switch Data
	u.exportUSWstats(r, labels, sw)
	u.exportPRTtable(r, labels, d.PortTable)
	// Gateway Data
	u.exportWANPorts(r, labels, d.Wan1, d.Wan2)
	u.exportUSGstats(r, labels, gw, d.SpeedtestStatus, d.Uplink)
	// Dream Machine System Data.
	r.send([]*metric{
		{u.Device.Info, gauge, 1.0, append(labels, infoLabels...)},
		{u.Device.Uptime, gauge, d.Uptime, labels},
	})

	for _, t := range d.Temperatures {
		r.send([]*metric{{u.Device.Temperature, gauge, t.Value, append(labels, t.Name, t.Type)}})
	}

	// UDM pro and UXG have hard drives.
	for _, t := range d.Storage {
		r.send([]*metric{
			{u.Device.Storage, gauge, t.Size.Val, append(labels, t.MountPoint, t.Name, "size")},
			{u.Device.Storage, gauge, t.Used.Val, append(labels, t.MountPoint, t.Name, "used")},
		})
	}
}
