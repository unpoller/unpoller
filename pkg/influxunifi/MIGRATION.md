# InfluxDB Output Migration Notes

## InfluxDB 3 support

UnPoller now supports InfluxDB 1.x, 2.x, and 3.x from the same `influxdb` output plugin.

Enable v3 explicitly:

```yaml
influxdb:
  version: 3
  url: http://influxdb3:8181
  auth_token: your-token
  database: unifi
```

For InfluxDB Cloud Serverless or Clustered, set `use_v2_api: true` so writes use the v2-compatible endpoint.

See [README.md](README.md) and `init/docker/docker-compose-influxdb3.yml` for examples.

## Schema changes (v1/v2/v3)

InfluxDB 3 rejects line protocol where the same key is used as both a tag and a field on one point. The following keys were adjusted for compatibility:

| Measurement | Change |
|-------------|--------|
| `subsystems` | Removed duplicate field `wan_ip` (remains a tag) |
| `clients` | Renamed tag `channel` to `channel_name` (numeric field `channel` unchanged) |
| `uap_radios` | Renamed field `channel` to `channel_num`; removed duplicate field `radio` (remains a tag) |
| `usg`, `ubb`, `uci`, `udm`, `uxg` | Removed duplicate fields `source` and/or `version` (remain tags) |

### Grafana dashboards

Import updated dashboards from the [unpoller/dashboards](https://github.com/unpoller/dashboards) repository (`v2.0.0` InfluxDB JSON files) if panels reference the old tag/field names.

### Existing InfluxDB 3 databases

If writes previously failed with errors such as `invalid column type for column 'wan_ip'`, drop and recreate the affected database or bucket, then allow UnPoller to recreate the schema with the corrected layout.

## Issue reference

Fixes [unpoller#1061](https://github.com/unpoller/unpoller/issues/1061).
