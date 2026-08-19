package promunifi

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/unpoller/unifi/v5"
)

// unasDevice holds Prometheus descriptors for UNAS Pro storage console metrics.
//
// The struct is deliberately flat even though the metrics come from five different entities
// (console, network IO, pool, RAID group, disk, share) with different label sets: Describe
// reflects over the members of this struct and sends each one to the channel, so a nested
// sub-struct would be sent as something that is not a *prometheus.Desc.
type unasDevice struct {
	// Console.
	Presence       *prometheus.Desc
	CPULoad        *prometheus.Desc
	CPUTemperature *prometheus.Desc
	MemoryTotal    *prometheus.Desc
	MemoryFree     *prometheus.Desc
	MemoryAvail    *prometheus.Desc
	ReceiveKBPS    *prometheus.Desc
	TransmitKBPS   *prometheus.Desc
	// Storage pools and their RAID groups.
	PoolCapacity           *prometheus.Desc
	PoolUsage              *prometheus.Desc
	RaidCurrentProtection  *prometheus.Desc
	RaidExpectedProtection *prometheus.Desc
	RaidProgress           *prometheus.Desc
	// Physical disks.
	DiskSize                     *prometheus.Desc
	DiskTemperature              *prometheus.Desc
	DiskHealthScore              *prometheus.Desc
	DiskPowerOnHours             *prometheus.Desc
	DiskRPM                      *prometheus.Desc
	DiskBadSectorCount           *prometheus.Desc
	DiskUncorrectableSectorCount *prometheus.Desc
	DiskReadErrorRate            *prometheus.Desc
	DiskSmartReadErrorCount      *prometheus.Desc
	DiskReadKBPS                 *prometheus.Desc
	DiskWriteKBPS                *prometheus.Desc
	// Shares ("drives").
	ShareQuota       *prometheus.Desc
	ShareUsage       *prometheus.Desc
	ShareMemberCount *prometheus.Desc
}

func descUNASDevice(ns string) *unasDevice { // nolint: funlen
	ns += "unas_"
	device := []string{"source", "name", "model"}
	console := []string{"source", "name"}
	pool := []string{"source", "name", "pool_id", "pool_type", "status"}
	raid := []string{"source", "name", "pool_id", "raid_group_id", "current_level", "config_level"}
	disk := []string{"source", "name", "slot_id", "pool_id", "disk_type", "state", "model", "serial"}
	share := []string{"source", "name", "share_id", "share_name", "share_type", "status"}

	return &unasDevice{
		Presence: prometheus.NewDesc(ns+"device_present",
			"UNAS console reachable (always 1 when present)", device, nil),
		CPULoad: prometheus.NewDesc(ns+"cpu_load_percent",
			"UNAS console current CPU load", console, nil),
		CPUTemperature: prometheus.NewDesc(ns+"cpu_temperature_celsius",
			"UNAS console CPU temperature", console, nil),
		MemoryTotal: prometheus.NewDesc(ns+"memory_total_bytes",
			"UNAS console total memory", console, nil),
		MemoryFree: prometheus.NewDesc(ns+"memory_free_bytes",
			"UNAS console free memory", console, nil),
		MemoryAvail: prometheus.NewDesc(ns+"memory_available_bytes",
			"UNAS console available memory", console, nil),
		ReceiveKBPS: prometheus.NewDesc(ns+"receive_kbps",
			"UNAS console instantaneous receive throughput in KB/s", console, nil),
		TransmitKBPS: prometheus.NewDesc(ns+"transmit_kbps",
			"UNAS console instantaneous transmit throughput in KB/s", console, nil),
		PoolCapacity: prometheus.NewDesc(ns+"pool_capacity_bytes",
			"UNAS storage pool total capacity", pool, nil),
		PoolUsage: prometheus.NewDesc(ns+"pool_usage_bytes",
			"UNAS storage pool bytes in use", pool, nil),
		RaidCurrentProtection: prometheus.NewDesc(ns+"raid_group_current_protection",
			"UNAS RAID group current redundancy (compare to expected to spot a degraded array)",
			raid, nil),
		RaidExpectedProtection: prometheus.NewDesc(ns+"raid_group_expected_protection",
			"UNAS RAID group expected redundancy", raid, nil),
		RaidProgress: prometheus.NewDesc(ns+"raid_group_progress_percent",
			"UNAS RAID group rebuild or expansion progress", raid, nil),
		DiskSize: prometheus.NewDesc(ns+"disk_size_bytes",
			"UNAS disk size", disk, nil),
		DiskTemperature: prometheus.NewDesc(ns+"disk_temperature_celsius",
			"UNAS disk temperature", disk, nil),
		DiskHealthScore: prometheus.NewDesc(ns+"disk_health_score",
			"UNAS disk health score as reported by the console", disk, nil),
		DiskPowerOnHours: prometheus.NewDesc(ns+"disk_power_on_hours",
			"UNAS disk SMART power-on hours", disk, nil),
		DiskRPM: prometheus.NewDesc(ns+"disk_rpm",
			"UNAS disk rotational speed (0 for solid state)", disk, nil),
		DiskBadSectorCount: prometheus.NewDesc(ns+"disk_bad_sectors",
			"UNAS disk bad sector count", disk, nil),
		DiskUncorrectableSectorCount: prometheus.NewDesc(ns+"disk_uncorrectable_sectors",
			"UNAS disk uncorrectable sector count", disk, nil),
		DiskReadErrorRate: prometheus.NewDesc(ns+"disk_read_error_rate",
			"UNAS disk SMART read error rate", disk, nil),
		DiskSmartReadErrorCount: prometheus.NewDesc(ns+"disk_smart_read_errors",
			"UNAS disk SMART read error count", disk, nil),
		DiskReadKBPS: prometheus.NewDesc(ns+"disk_read_kbps",
			"UNAS disk instantaneous read throughput in KB/s", disk, nil),
		DiskWriteKBPS: prometheus.NewDesc(ns+"disk_write_kbps",
			"UNAS disk instantaneous write throughput in KB/s", disk, nil),
		ShareQuota: prometheus.NewDesc(ns+"share_quota_bytes",
			"UNAS share quota (0 when unlimited)", share, nil),
		ShareUsage: prometheus.NewDesc(ns+"share_usage_bytes",
			"UNAS share bytes in use", share, nil),
		ShareMemberCount: prometheus.NewDesc(ns+"share_members",
			"UNAS share member count", share, nil),
	}
}

// exportUNASDevice emits every metric family for one UNAS console.
//
// Each section is guarded independently because GetUNASDevice leaves the field for a failing
// endpoint nil rather than failing the whole poll -- a console that serves storage but not
// network-io still reports its disks.
func (u *promUnifi) exportUNASDevice(r report, d *unifi.UNASDevice) {
	if d == nil {
		return
	}

	name, model := d.Name(), d.Model()
	console := []string{d.SourceName, name}

	r.send([]*metric{{u.UNASDevice.Presence, gauge, 1.0, []string{d.SourceName, name, model}}})

	if info := d.DeviceInfo; info != nil {
		r.send([]*metric{
			{u.UNASDevice.CPULoad, gauge, info.CPU.CurrentLoad, console},
			{u.UNASDevice.CPUTemperature, gauge, info.CPU.Temperature, console},
			{u.UNASDevice.MemoryTotal, gauge, info.Memory.Total, console},
			{u.UNASDevice.MemoryFree, gauge, info.Memory.Free, console},
			{u.UNASDevice.MemoryAvail, gauge, info.Memory.Available, console},
		})
	}

	if io := d.NetworkIO; io != nil {
		r.send([]*metric{
			{u.UNASDevice.ReceiveKBPS, gauge, io.ReceiveKBPS, console},
			{u.UNASDevice.TransmitKBPS, gauge, io.TransmitKBPS, console},
		})
	}

	if s := d.Storage; s != nil {
		u.exportUNASStorage(r, d, s)
	}

	for _, share := range d.Drives {
		if share == nil {
			continue
		}

		labels := []string{d.SourceName, name, share.ID, share.Name, share.Type, share.Status}

		r.send([]*metric{
			{u.UNASDevice.ShareQuota, gauge, share.Quota, labels},
			{u.UNASDevice.ShareUsage, gauge, share.Usage, labels},
			{u.UNASDevice.ShareMemberCount, gauge, share.MemberCount, labels},
		})
	}
}

func (u *promUnifi) exportUNASStorage(r report, d *unifi.UNASDevice, s *unifi.UNASStorage) {
	name := d.Name()

	for _, p := range s.Pools {
		labels := []string{d.SourceName, name, p.ID, p.Type, p.Status}

		r.send([]*metric{
			{u.UNASDevice.PoolCapacity, gauge, p.Capacity, labels},
			{u.UNASDevice.PoolUsage, gauge, p.Usage, labels},
		})

		for _, rg := range p.RaidGroups {
			raid := []string{d.SourceName, name, p.ID, rg.ID, rg.CurrentLevel, rg.ConfigLevel}

			r.send([]*metric{
				{u.UNASDevice.RaidCurrentProtection, gauge, rg.CurrentProtection, raid},
				{u.UNASDevice.RaidExpectedProtection, gauge, rg.ExpectedProtection, raid},
				{u.UNASDevice.RaidProgress, gauge, rg.Progress, raid},
			})
		}
	}

	for _, disk := range s.Disks {
		labels := []string{
			d.SourceName, name, disk.SlotID, disk.PoolID,
			disk.Type, disk.State, disk.Model, disk.Serial,
		}

		r.send([]*metric{
			{u.UNASDevice.DiskSize, gauge, disk.Size, labels},
			{u.UNASDevice.DiskTemperature, gauge, disk.Temperature, labels},
			{u.UNASDevice.DiskHealthScore, gauge, disk.HealthScore, labels},
			{u.UNASDevice.DiskPowerOnHours, gauge, disk.PowerOnHours, labels},
			{u.UNASDevice.DiskRPM, gauge, disk.RPM, labels},
			{u.UNASDevice.DiskBadSectorCount, gauge, disk.BadSectorCount, labels},
			{u.UNASDevice.DiskUncorrectableSectorCount, gauge, disk.UncorrectableSectorCount, labels},
			{u.UNASDevice.DiskReadErrorRate, gauge, disk.ReadErrorRate, labels},
			{u.UNASDevice.DiskSmartReadErrorCount, gauge, disk.SmartReadErrorCount, labels},
			{u.UNASDevice.DiskReadKBPS, gauge, disk.ReadKBPS, labels},
			{u.UNASDevice.DiskWriteKBPS, gauge, disk.WriteKBPS, labels},
		})
	}
}
