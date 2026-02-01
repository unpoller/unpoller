# UniFi Infrastructure Alerts

Example Prometheus and Loki alerting rules for monitoring UniFi infrastructure with unPoller.

## Overview

- **Prometheus** – Metrics from devices, clients, UPS/PDU, controller, sites, WAN, DHCP, rogue APs, and more
- **Loki** – Logs for events, alarms, IDS, anomalies, and system logs

These examples assume the default Prometheus namespace `unpoller`. Adjust metric names if you use a custom `prometheus.namespace`.

---

## Prometheus Alerts

Place `prometheus/unifi-alerts.yaml` in your Prometheus `rule_files` or Grafana Alerting.

### UPS (unifi-ups)

| Alert | Trigger | Severity | Description |
|-------|---------|----------|-------------|
| **UnifiUPSLowBattery** | Battery level < 20% for 5m | warning | UPS needs attention; plan for charging or replacement |
| **UnifiUPSCriticalBattery** | Battery level < 10% for 2m | critical | UPS near depletion; prepare for shutdown |
| **UnifiUPSOnBattery** | Running on battery for 1m | warning | Power outage or AC loss; UPS sustaining load |
| **UnifiUPSLowRuntime** | Runtime < 5 min (and known) for 5m | warning | Little runtime left; prioritize critical loads |
| **UnifiUPSHighLoad** | Load > 80% of capacity for 10m | warning | UPS near capacity; consider load shedding |
| **UnifiUPSBMSAnomaly** | BMS anomaly count > 0 for 5m | warning | Battery management system issue; check UPS health |
| **UnifiUPSNotCharging** | Not charging and battery < 100% for 30m | warning | Battery not charging; check power or battery |

*Requires: PDU/UPS devices with vbms_table (e.g. USW-DA-23-POE-UPS)*

### Controller (unifi-controller)

| Alert | Trigger | Severity | Description |
|-------|---------|----------|-------------|
| **UnifiControllerUpdateAvailable** | Update available for 1h | info | Controller firmware update available |
| **UnifiControllerUnsupportedDevices** | Unsupported device count > 0 for 1h | warning | Devices no longer supported; plan upgrades |

### Controller Health (unifi-controller-health)

| Alert | Trigger | Severity | Description |
|-------|---------|----------|-------------|
| **UnifiControllerRecentlyRestarted** | Uptime < 1h for 5m | info | Controller recently restarted; may indicate maintenance or crash |
| **UnifiControllerBackupDisabled** | Auto backup disabled for 24h | info | Backups disabled; enable for disaster recovery |

### Devices (unifi-devices)

| Alert | Trigger | Severity | Description |
|-------|---------|----------|-------------|
| **UnifiDeviceHighCPU** | CPU > 90% for 10m | warning | Device under heavy load; investigate |
| **UnifiDeviceHighMemory** | Memory > 90% for 10m | warning | Device memory pressure; may impact performance |
| **UnifiDeviceUpgradeAvailable** | Firmware upgrade available for 1h | info | Device has firmware update available |

### Site (unifi-site)

| Alert | Trigger | Severity | Description |
|-------|---------|----------|-------------|
| **UnifiSiteHighDisconnectedDevices** | Disconnected devices > 0 (WLAN/WAN/LAN) for 15m | warning | Devices offline; check power, connectivity, adoption |
| **UnifiSitePendingAdoptions** | Pending adoptions > 0 for 1h | info | Devices awaiting adoption |
| **UnifiSiteWANDrops** | WAN disconnections in last 1h > 0 | warning | Internet connectivity issues |
| **UnifiSiteHighLatency** | Internet latency > 500ms for 10m | warning | Poor internet performance |

*Requires: save_sites=true*

### WAN (unifi-wan)

| Alert | Trigger | Severity | Description |
|-------|---------|----------|-------------|
| **UnifiWANLowUptime** | WAN uptime < 95% for 15m | warning | WAN link unstable; check ISP or cabling |
| **UnifiWANPeakDownloadUtilization** | Peak download > 90% of capacity for 10m | info | Download near capacity; consider upgrade |
| **UnifiWANPeakUploadUtilization** | Peak upload > 90% of capacity for 10m | info | Upload near capacity; consider upgrade |

*Requires: WAN metrics (UDM/UDM-Pro/UCG)*

### DHCP (unifi-dhcp)

| Alert | Trigger | Severity | Description |
|-------|---------|----------|-------------|
| **UnifiDHCPPoolExhaustion** | Pool utilization > 90% for 15m | warning | DHCP pool nearly full; expand range or reduce lease time |
| **UnifiDHCPPoolCritical** | Pool utilization > 98% for 5m | critical | Pool almost exhausted; new devices may not get IPs |

*Requires: save_dhcp or DHCP lease collection*

### Rogue AP (unifi-rogue)

| Alert | Trigger | Severity | Description |
|-------|---------|----------|-------------|
| **UnifiRogueAPDetected** | Any rogue AP detected for 5m | warning | Unauthorized access point; investigate and remediate |

*Requires: save_rogue=true*

---

## Prometheus Recording Rules

Place `prometheus/unifi-recording-rules.yaml` in your Prometheus `rule_files` to pre-compute aggregates for faster dashboards and simpler alerting.

### UPS Recording Rules (interval: 1m)

| Recorded Metric | Expression | Description |
|-----------------|------------|-------------|
| `unpoller:ups_on_battery:count` | Count of UPSes with battery_mode=1 by site | UPS devices running on battery per site |
| `unpoller:ups_min_battery_level_percent:min` | Min battery level by site | Worst battery level per site |
| `unpoller:ups_min_runtime_seconds:min` | Min runtime (≥0) by site | Worst runtime remaining per site |
| `unpoller:ups_total_power_output_watts:sum` | Sum of power output by site | Total UPS load per site |
| `unpoller:ups_total_power_budget_watts:sum` | Sum of power budget by site | Total UPS capacity per site |
| `unpoller:ups_bms_anomaly_count:sum` | Sum of devices with BMS anomaly by site | UPSes with BMS issues per site |

### Device Recording Rules (interval: 1m)

| Recorded Metric | Expression | Description |
|-----------------|------------|-------------|
| `unpoller:device_count:by_type` | Count of devices by type (uap, usw, pdu, etc.) per site | Device inventory by type |
| `unpoller:device_count:total` | Total device count per site | Total devices per site |
| `unpoller:device_high_cpu_count:count` | Count of devices with CPU > 90% per site | Overloaded devices per site |
| `unpoller:device_high_memory_count:count` | Count of devices with memory > 90% per site | Memory-pressure devices per site |

### Controller Recording Rules (interval: 5m)

| Recorded Metric | Expression | Description |
|-----------------|------------|-------------|
| `unpoller:controller_update_available:count` | Count of controllers with update available | Controllers needing updates |
| `unpoller:controller_unsupported_devices_total:sum` | Sum of unsupported devices | Total unsupported devices across controllers |

---

## Loki Alerts

Place `loki/unifi-alerts.yaml` in your Loki Ruler config. Loki must be run with `-ruler.enable=true` and `-ruler.storage.path` configured.

### Alarms (unifi-alarms)

| Alert | Trigger | Severity | Description |
|-------|---------|----------|-------------|
| **UnifiHighAlarmRate** | > 20 alarms in 15m for 5m | warning | Elevated alarm volume; review controller |

*Requires: save_alarms=true*

### IDS (unifi-ids)

| Alert | Trigger | Severity | Description |
|-------|---------|----------|-------------|
| **UnifiIDSEvent** | Any IDS event in 5m for 1m | warning | Intrusion detection triggered; review logs |
| **UnifiIDSHighVolume** | > 50 IDS events in 1h for 5m | critical | High IDS volume; possible attack |

*Requires: save_ids=true*

### Anomalies (unifi-anomalies)

| Alert | Trigger | Severity | Description |
|-------|---------|----------|-------------|
| **UnifiAnomalyDetected** | > 5 anomalies in 10m for 5m | warning | Multiple anomalies; check network health |

*Requires: save_anomalies=true*

### System Log (unifi-system-log)

| Alert | Trigger | Severity | Description |
|-------|---------|----------|-------------|
| **UnifiSystemLogCritical** | Any CRITICAL log in 5m for 1m | critical | Critical system log; immediate attention |
| **UnifiSystemLogHighSeverity** | > 10 CRITICAL/HIGH/ERROR logs in 15m for 5m | warning | High volume of severe logs |
| **UnifiSystemLogAuthFailure** | > 5 auth failure matches in 1h for 5m | warning | Authentication failures; possible brute force |

*Requires: save_syslog=true (UDM/UDM-Pro) or save_events=true (older controllers)*

### Events (unifi-events)

| Alert | Trigger | Severity | Description |
|-------|---------|----------|-------------|
| **UnifiEventSpike** | > 100 events in 5m for 5m | info | Event spike; may indicate churn or issue |

*Requires: save_events=true*

---

## Configuration

**Prometheus (prometheus.yml):**

```yaml
rule_files:
  - /etc/prometheus/rules/unifi-alerts.yaml
  - /etc/prometheus/rules/unifi-recording-rules.yaml
```

**Loki (loki-config.yaml):**

```yaml
ruler:
  enable_api: true
  storage:
    type: local
    local:
      directory: /loki/rules
  rule_path: /loki/rules-temp
  alertmanager_url: http://alertmanager:9093
```

Mount the `loki/` directory into your Loki container at `/loki/rules/`.

## AlertManager Integration

Both Prometheus and Loki can forward alerts to Alertmanager. Configure Alertmanager receivers (Slack, PagerDuty, email, etc.) as needed.

## Customization

- Tune thresholds (battery %, runtime seconds, CPU %, etc.) for your environment
- Add or remove labels in `annotations` for your notification channels
- Adjust `for` durations to reduce noise or catch issues sooner
- Disable alert groups that don't apply (e.g. remove UPS alerts if you have no UPS devices)
