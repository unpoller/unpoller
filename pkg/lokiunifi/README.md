# lokiunifi

Loki Output Plugin for UnPoller

This plugin writes UniFi Events, System Logs, IDS, Alarms, and Anomalies to Loki as JSON.

## Log Types

| Application Label | Config Option | API | Description |
|-------------------|---------------|-----|-------------|
| `unifi_system_log` | `save_syslog` | v2 | System log events (UDM recommended) |
| `unifi_event` | `save_events` | v1 | Legacy events (older controllers) |
| `unifi_ids` | `save_ids` | v1 | Intrusion Detection System events |
| `unifi_alarm` | `save_alarms` | v1 | Alarm events |
| `unifi_anomaly` | `save_anomalies` | v1 | Anomaly events |

## Querying in Loki

All logs are stored as JSON. Use Loki's `| json` parser to extract fields:

```logql
{application="unifi_system_log"} | json
```

Filter by severity:
```logql
{application="unifi_system_log", severity="HIGH"} | json
```

Extract specific fields:
```logql
{application="unifi_system_log"} | json | line_format "{{.message}}"
```

## Example Config

```toml
[loki]
  # URL is the only required setting for Loki.
  url = "http://192.168.3.2:3100"

  # How often to poll UniFi and report to Loki.
  interval = "2m"

  # How long to wait for Loki responses.
  timeout = "5s"

  # Set these to use basic auth.
  #user = ""
  #pass = ""

  # Used for auth-less multi-tenant.
  #tenant_id = ""

[unifi.defaults]
  # For UDM/UDM-Pro/UCG devices, use save_syslog (v2 API)
  save_syslog = true

  # For older controllers, use save_events (v1 API)
  save_events = false

  # Other log types
  save_ids = false
  save_alarms = false
  save_anomalies = false
```

## Environment Variables

```bash
UP_LOKI_URL=http://localhost:3100
UP_LOKI_INTERVAL=2m
UP_UNIFI_DEFAULT_SAVE_SYSLOG=true
UP_UNIFI_DEFAULT_SAVE_EVENTS=false
```
