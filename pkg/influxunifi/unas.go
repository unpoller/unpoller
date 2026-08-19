package influxunifi

import (
	"github.com/unpoller/unifi/v5"
)

// batchUNASDevice generates InfluxDB points for one UNAS Pro storage console.
//
// Four measurements are emitted rather than one: the console, its storage pools, its physical
// disks, and its shares each have their own tag set, and folding them together would give
// every point the union of those tags with most of them empty.
//
// Each section is guarded independently because GetUNASDevice leaves the field for a failing
// endpoint nil rather than failing the whole poll -- a console that serves storage but not
// network-io still reports its disks.
func (u *InfluxUnifi) batchUNASDevice(r report, d *unifi.UNASDevice) { // nolint: funlen
	if d == nil {
		return
	}

	tags := map[string]string{
		"source": d.SourceName,
		"name":   d.Name(),
		"model":  d.Model(),
	}

	fields := map[string]any{"present": 1}

	if info := d.DeviceInfo; info != nil {
		tags["version"] = info.Version
		tags["firmware_version"] = info.FirmwareVersion
		tags["status"] = info.Status
		fields["cpu_load"] = info.CPU.CurrentLoad.Val
		fields["cpu_temperature"] = info.CPU.Temperature.Val
		fields["memory_total"] = info.Memory.Total.Val
		fields["memory_free"] = info.Memory.Free.Val
		fields["memory_available"] = info.Memory.Available.Val
		fields["startup_time"] = info.StartupTime
	}

	if io := d.NetworkIO; io != nil {
		fields["receive_kbps"] = io.ReceiveKBPS.Val
		fields["transmit_kbps"] = io.TransmitKBPS.Val
	}

	r.send(&metric{Table: "unas_device", Tags: tags, Fields: fields})

	if s := d.Storage; s != nil {
		u.batchUNASStorage(r, d, s)
	}

	for _, share := range d.Drives {
		if share == nil {
			continue
		}

		r.send(&metric{
			Table: "unas_share",
			Tags: map[string]string{
				"source":            d.SourceName,
				"name":              d.Name(),
				"share_id":          share.ID,
				"share_name":        share.Name,
				"share_type":        share.Type,
				"status":            share.Status,
				"storage_pool_id":   share.StoragePoolID,
				"role":              share.Role,
				"encryption_status": share.Protections.EncryptionStatus,
			},
			Fields: map[string]any{
				"quota":                 share.Quota.Val,
				"usage":                 share.Usage.Val,
				"member_count":          share.MemberCount.Val,
				"snapshot_enabled":      share.Protections.SnapshotEnabled.Val,
				"remote_backup_enabled": share.Protections.RemoteBackupEnabled.Val,
				"deduplication":         share.Deduplication,
				"compression_level":     share.CompressionLevel,
			},
		})
	}
}

func (u *InfluxUnifi) batchUNASStorage(r report, d *unifi.UNASDevice, s *unifi.UNASStorage) {
	for _, p := range s.Pools {
		fields := map[string]any{
			"capacity":            p.Capacity.Val,
			"usage":               p.Usage.Val,
			"raid_group_count":    len(p.RaidGroups),
			"initializing_status": p.InitializingStatus,
		}

		// A pool has one active RAID group; its redundancy is the number an operator alerts
		// on, so it is flattened onto the pool point rather than given its own measurement.
		for _, rg := range p.RaidGroups {
			if rg.ID != p.ActiveRaidGroupID {
				continue
			}

			fields["current_protection"] = rg.CurrentProtection.Val
			fields["expected_protection"] = rg.ExpectedProtection.Val
			fields["progress"] = rg.Progress.Val
			fields["current_level"] = rg.CurrentLevel
		}

		r.send(&metric{
			Table: "unas_pool",
			Tags: map[string]string{
				"source":       d.SourceName,
				"name":         d.Name(),
				"pool_id":      p.ID,
				"pool_type":    p.Type,
				"status":       p.Status,
				"prefer_level": p.PreferLevel,
			},
			Fields: fields,
		})
	}

	for _, disk := range s.Disks {
		r.send(&metric{
			Table: "unas_disk",
			Tags: map[string]string{
				"source":    d.SourceName,
				"name":      d.Name(),
				"slot_id":   disk.SlotID,
				"location":  disk.Location,
				"pool_id":   disk.PoolID,
				"disk_type": disk.Type,
				"state":     disk.State,
				"model":     disk.Model,
				"serial":    disk.Serial,
				"firmware":  disk.Firmware,
			},
			Fields: map[string]any{
				"size":                       disk.Size.Val,
				"temperature":                disk.Temperature.Val,
				"health_score":               disk.HealthScore.Val,
				"power_on_hours":             disk.PowerOnHours.Val,
				"rpm":                        disk.RPM.Val,
				"bad_sector_count":           disk.BadSectorCount.Val,
				"uncorrectable_sector_count": disk.UncorrectableSectorCount.Val,
				"read_error_rate":            disk.ReadErrorRate.Val,
				"smart_read_error_count":     disk.SmartReadErrorCount.Val,
				"read_kbps":                  disk.ReadKBPS.Val,
				"write_kbps":                 disk.WriteKBPS.Val,
				"is_global_hot_spare":        disk.IsGlobalHotSpare.Val,
				"is_local_hot_spare":         disk.IsLocalHotSpare.Val,
			},
		})
	}
}
