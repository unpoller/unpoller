package datadogunifi

import (
	"github.com/unpoller/unifi/v6"
)

// batchUNASDevice generates UNAS Pro storage console datapoints for Datadog.
//
// Each section is guarded independently because GetUNASDevice leaves the field for a failing
// endpoint nil rather than failing the whole poll -- a console that serves storage but not
// network-io still reports its disks.
func (u *DatadogUnifi) batchUNASDevice(r report, d *unifi.UNASDevice) {
	if d == nil {
		return
	}

	metricName := metricNamespace("unas_device")

	tags := cleanTags(map[string]string{
		"source": d.SourceName,
		"name":   d.Name(),
		"model":  d.Model(),
	})

	deviceTags := tagMapToTags(tags)

	_ = r.reportGauge(metricName("present"), 1.0, deviceTags)

	if info := d.DeviceInfo; info != nil {
		_ = r.reportGauge(metricName("cpu_load"), info.CPU.CurrentLoad.Val, deviceTags)
		_ = r.reportGauge(metricName("cpu_temperature"), info.CPU.Temperature.Val, deviceTags)
		_ = r.reportGauge(metricName("memory_total"), info.Memory.Total.Val, deviceTags)
		_ = r.reportGauge(metricName("memory_free"), info.Memory.Free.Val, deviceTags)
		_ = r.reportGauge(metricName("memory_available"), info.Memory.Available.Val, deviceTags)
	}

	if io := d.NetworkIO; io != nil {
		_ = r.reportGauge(metricName("receive_kbps"), io.ReceiveKBPS.Val, deviceTags)
		_ = r.reportGauge(metricName("transmit_kbps"), io.TransmitKBPS.Val, deviceTags)
	}

	if s := d.Storage; s != nil {
		u.batchUNASPools(r, d, s)
		u.batchUNASDisks(r, d, s)
	}

	u.batchUNASShares(r, d)
}

func (u *DatadogUnifi) batchUNASPools(r report, d *unifi.UNASDevice, s *unifi.UNASStorage) {
	metricName := metricNamespace("unas_pool")
	raidName := metricNamespace("unas_raid_group")

	for _, p := range s.Pools {
		tags := tagMapToTags(cleanTags(map[string]string{
			"source":    d.SourceName,
			"name":      d.Name(),
			"pool_id":   p.ID,
			"pool_type": p.Type,
			"status":    p.Status,
		}))

		_ = r.reportGauge(metricName("capacity"), p.Capacity.Val, tags)
		_ = r.reportGauge(metricName("usage"), p.Usage.Val, tags)

		for _, rg := range p.RaidGroups {
			raidTags := tagMapToTags(cleanTags(map[string]string{
				"source":        d.SourceName,
				"name":          d.Name(),
				"pool_id":       p.ID,
				"raid_group_id": rg.ID,
				"current_level": rg.CurrentLevel,
				"config_level":  rg.ConfigLevel,
			}))

			_ = r.reportGauge(raidName("current_protection"), rg.CurrentProtection.Val, raidTags)
			_ = r.reportGauge(raidName("expected_protection"), rg.ExpectedProtection.Val, raidTags)
			_ = r.reportGauge(raidName("progress"), rg.Progress.Val, raidTags)
		}

		if c := p.Cache; c != nil {
			u.batchUNASPoolCache(r, d, p.ID, c)
		}
	}
}

func (u *DatadogUnifi) batchUNASPoolCache(r report, d *unifi.UNASDevice, poolID string, c *unifi.UNASPoolCache) {
	metricName := metricNamespace("unas_pool_cache")

	tags := tagMapToTags(cleanTags(map[string]string{
		"source":  d.SourceName,
		"name":    d.Name(),
		"pool_id": poolID,
		"mode":    c.Mode,
		"status":  c.Status,
	}))

	_ = r.reportGauge(metricName("size"), c.Size.Val, tags)
	_ = r.reportGauge(metricName("used"), c.Composition.Used.Val, tags)
	_ = r.reportGauge(metricName("dirty"), c.Composition.Dirty.Val, tags)
	_ = r.reportGauge(metricName("read_hit_rate"), c.Composition.ReadHitRate.Val, tags)
	_ = r.reportGauge(metricName("write_hit_rate"), c.Composition.WriteHitRate.Val, tags)
}

func (u *DatadogUnifi) batchUNASDisks(r report, d *unifi.UNASDevice, s *unifi.UNASStorage) {
	metricName := metricNamespace("unas_disk")

	for _, disk := range s.Disks {
		u.batchUNASDisk(r, d, disk)
	}

	// Cache bays hold SSDs with the same fields as pool disks; disk_type tells them apart.
	for _, disk := range s.CacheSlots {
		if disk.State == "empty" {
			continue
		}

		tags := u.batchUNASDisk(r, d, disk)
		_ = r.reportGauge(metricName("life_span"), disk.LifeSpan.Val, tags)
	}
}

// batchUNASDisk reports one disk's gauges and returns its tags.
func (u *DatadogUnifi) batchUNASDisk(r report, d *unifi.UNASDevice, disk unifi.UNASDisk) []string {
	metricName := metricNamespace("unas_disk")

	tags := tagMapToTags(cleanTags(map[string]string{
		"source":    d.SourceName,
		"name":      d.Name(),
		"slot_id":   disk.SlotID,
		"location":  disk.Location,
		"pool_id":   disk.PoolID,
		"disk_type": disk.Type,
		"state":     disk.State,
		"model":     disk.Model,
		"serial":    disk.Serial,
	}))

	_ = r.reportGauge(metricName("size"), disk.Size.Val, tags)
	_ = r.reportGauge(metricName("temperature"), disk.Temperature.Val, tags)
	_ = r.reportGauge(metricName("health_score"), disk.HealthScore.Val, tags)
	_ = r.reportGauge(metricName("power_on_hours"), disk.PowerOnHours.Val, tags)
	_ = r.reportGauge(metricName("rpm"), disk.RPM.Val, tags)
	_ = r.reportGauge(metricName("bad_sector_count"), disk.BadSectorCount.Val, tags)
	_ = r.reportGauge(metricName("uncorrectable_sector_count"), disk.UncorrectableSectorCount.Val, tags)
	_ = r.reportGauge(metricName("read_error_rate"), disk.ReadErrorRate.Val, tags)
	_ = r.reportGauge(metricName("smart_read_error_count"), disk.SmartReadErrorCount.Val, tags)
	_ = r.reportGauge(metricName("read_kbps"), disk.ReadKBPS.Val, tags)
	_ = r.reportGauge(metricName("write_kbps"), disk.WriteKBPS.Val, tags)

	return tags
}

func (u *DatadogUnifi) batchUNASShares(r report, d *unifi.UNASDevice) {
	metricName := metricNamespace("unas_share")

	for _, share := range d.Drives {
		if share == nil {
			continue
		}

		tags := tagMapToTags(cleanTags(map[string]string{
			"source":     d.SourceName,
			"name":       d.Name(),
			"share_id":   share.ID,
			"share_name": share.Name,
			"share_type": share.Type,
			"status":     share.Status,
			"role":       share.Role,
		}))

		_ = r.reportGauge(metricName("quota"), share.Quota.Val, tags)
		_ = r.reportGauge(metricName("usage"), share.Usage.Val, tags)
		_ = r.reportGauge(metricName("member_count"), share.MemberCount.Val, tags)
	}
}
