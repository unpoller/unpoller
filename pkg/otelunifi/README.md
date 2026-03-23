# otelunifi — OpenTelemetry Output Plugin

Exports UniFi metrics to any [OpenTelemetry Protocol (OTLP)](https://opentelemetry.io/docs/specs/otel/protocol/) compatible backend via push, using the Go OTel SDK.

Compatible backends include Grafana Alloy/Mimir, Honeycomb, Datadog (via OTel collector), Grafana Tempo, New Relic, Lightstep, and any vendor that accepts OTLP.

## Configuration

The plugin is **disabled by default**. Set `enable = true` (or `UP_OTEL_ENABLE=true`) to enable it.

### TOML

```toml
[otel]
  url      = "http://localhost:4318"   # OTLP HTTP endpoint (default)
  protocol = "http"                    # "http" (default) or "grpc"
  interval = "30s"
  timeout  = "10s"
  enable   = true
  dead_ports = false

  # Optional bearer token for authenticated collectors (e.g. Grafana Cloud)
  api_key  = ""
```

### YAML

```yaml
otel:
  url: "http://localhost:4318"
  protocol: http
  interval: 30s
  timeout: 10s
  enable: true
  dead_ports: false
  api_key: ""
```

### Environment Variables

All config keys use the `UP_OTEL_` prefix:

| Variable | Default | Description |
|---|---|---|
| `UP_OTEL_URL` | `http://localhost:4318` | OTLP endpoint URL |
| `UP_OTEL_PROTOCOL` | `http` | Transport: `http` or `grpc` |
| `UP_OTEL_INTERVAL` | `30s` | Push interval |
| `UP_OTEL_TIMEOUT` | `10s` | Per-export timeout |
| `UP_OTEL_ENABLE` | `false` | Set to `true` to enable |
| `UP_OTEL_API_KEY` | `` | Bearer token for auth |
| `UP_OTEL_DEAD_PORTS` | `false` | Include down/disabled switch ports |

## Protocol Notes

- **HTTP** (`protocol = "http"`): Sends to `<url>/v1/metrics`. Default port `4318`.
- **gRPC** (`protocol = "grpc"`): Sends to `<host>:<port>`. Default `localhost:4317`. The URL for gRPC should be `host:port` (no scheme).

## Exported Metrics

All metrics use the `unifi_` prefix and carry identifying attributes (labels).

### Site metrics (`unifi_site_*`)

Attributes: `site_name`, `source`, `subsystem`, `status`

| Metric | Description |
|---|---|
| `unifi_site_users` | Connected user count |
| `unifi_site_guests` | Connected guest count |
| `unifi_site_iot` | IoT device count |
| `unifi_site_aps` | Access point count |
| `unifi_site_gateways` | Gateway count |
| `unifi_site_switches` | Switch count |
| `unifi_site_adopted` | Adopted device count |
| `unifi_site_disconnected` | Disconnected device count |
| `unifi_site_latency_seconds` | WAN latency |
| `unifi_site_uptime_seconds` | Site uptime |
| `unifi_site_tx_bytes_rate` | Transmit bytes rate |
| `unifi_site_rx_bytes_rate` | Receive bytes rate |

### Client metrics (`unifi_client_*`)

Attributes: `mac`, `name`, `ip`, `site_name`, `source`, `oui`, `network`, `ap_name`, `sw_name`, `wired`

Wireless-only additional attributes: `essid`, `radio`, `radio_proto`

| Metric | Description |
|---|---|
| `unifi_client_uptime_seconds` | Client uptime |
| `unifi_client_rx_bytes` | Total bytes received |
| `unifi_client_tx_bytes` | Total bytes transmitted |
| `unifi_client_rx_bytes_rate` | Receive rate |
| `unifi_client_tx_bytes_rate` | Transmit rate |
| `unifi_client_signal_db` | Signal strength (wireless) |
| `unifi_client_noise_db` | Noise floor (wireless) |
| `unifi_client_rssi_db` | RSSI (wireless) |
| `unifi_client_tx_rate_bps` | TX rate (wireless) |
| `unifi_client_rx_rate_bps` | RX rate (wireless) |

### Device metrics

#### UAP (`unifi_device_uap_*`)

Attributes: `mac`, `name`, `model`, `version`, `type`, `ip`, `site_name`, `source`

Includes: `up`, `uptime_seconds`, `cpu_utilization`, `mem_utilization`, `load_avg_{1,5,15}`, per-radio `channel`/`tx_power_dbm`, per-VAP `num_stations`/`satisfaction`/`rx_bytes`/`tx_bytes`.

#### USW (`unifi_device_usw_*`)

Attributes: `mac`, `name`, `model`, `version`, `type`, `ip`, `site_name`, `source`

Includes: `up`, `uptime_seconds`, `cpu_utilization`, `mem_utilization`, `load_avg_1`, `rx_bytes`, `tx_bytes`, and per-port metrics (`port_up`, `port_speed_mbps`, `port_rx_bytes`, `port_tx_bytes`, `port_poe_*`, etc.).

#### USG (`unifi_device_usg_*`)

Includes: `up`, `uptime_seconds`, `cpu_utilization`, `mem_utilization`, and per-WAN interface (`wan_rx_bytes`, `wan_tx_bytes`, `wan_rx_packets`, `wan_tx_packets`, `wan_rx_errors`, `wan_tx_errors`, `wan_speed_mbps`).

#### UDM (`unifi_device_udm_*`)

Includes: `up`, `uptime_seconds`, `cpu_utilization`, `mem_utilization`, `load_avg_{1,5,15}`.

#### UXG (`unifi_device_uxg_*`)

Includes: `up`, `uptime_seconds`, `cpu_utilization`, `mem_utilization`, `load_avg_1`.

## Example: Grafana Alloy

```alloy
otelcol.receiver.otlp "default" {
  grpc { endpoint = "0.0.0.0:4317" }
  http { endpoint = "0.0.0.0:4318" }

  output {
    metrics = [otelcol.exporter.prometheus.default.input]
  }
}
```

Set `UP_OTEL_URL=http://alloy-host:4318` and `UP_OTEL_ENABLE=true` in unpoller's environment.

## Example: Grafana Cloud (OTLP with auth)

```toml
[otel]
  url      = "https://otlp-gateway-prod-us-central-0.grafana.net/otlp"
  protocol = "http"
  api_key  = "instanceID:grafana_cloud_api_token"
  interval = "60s"
  enable   = true
```
