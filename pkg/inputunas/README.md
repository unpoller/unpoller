# UNAS Pro Input Plugin

Polls UniFi UNAS Pro storage consoles and hands their metrics to every configured output.

**This plugin is opt-in and does nothing until you configure at least one device.** With no
`[unas]` section, or a section with no devices, it stays silent and inert — it logs nothing
and polls nothing.

## Why a separate plugin?

A UNAS Pro is a standalone UniFi OS console with **no Network application**: no sites, no
`/status` endpoint, its own credentials, and usually its own host. `inputunifi`'s
sites-then-devices flow has nothing to offer it, so UNAS consoles are configured separately
rather than as another UniFi controller.

## Configuration

```toml
[unas]
  disable = false

# Applied to any device that does not set its own value.
[unas.defaults]
  user       = "unpoller"
  pass       = ""
  verify_ssl = false
  timeout    = "60s"

# Repeat for each console. Use its own local-account credentials.
# Do not add a path after the host.
[[unas.device]]
  url            = "https://192.168.1.10"
  user           = "unpoller"
  pass           = "unpoller"
  verify_ssl     = false
  timeout        = "60s"
  ssl_cert_paths = []
```

JSON and YAML use `devices` (plural) as the list key; TOML uses repeated `[[unas.device]]`
tables. See `examples/up.{conf,json,yaml}.example`.

### Environment variables

| Variable | Meaning |
|---|---|
| `UP_UNAS_DISABLE` | Disable the plugin outright. |
| `UP_UNAS_DEFAULT_USER` | Default username for every device. |
| `UP_UNAS_DEFAULT_PASS` | Default password. |
| `UP_UNAS_DEFAULT_VERIFY_SSL` | Default TLS verification. |
| `UP_UNAS_DEFAULT_TIMEOUT` | Default HTTP timeout, e.g. `60s`. |
| `UP_UNAS_DEVICE_0_URL` | First console's URL. |
| `UP_UNAS_DEVICE_0_USER` | First console's username. |
| `UP_UNAS_DEVICE_0_PASS` | First console's password. |
| `UP_UNAS_DEVICE_0_VERIFY_SSL` | First console's TLS verification. |
| `UP_UNAS_DEVICE_0_TIMEOUT` | First console's HTTP timeout. |

Increment the index for additional consoles. Note the defaults block is `DEFAULT`, singular.

### Authentication

Username and password only. API keys are **not** supported: the UniFi library routes key
auth through `/status`, which a storage-only console does not serve.

A UNAS session expires after roughly two hours. When it does, every endpoint starts failing
at once; the plugin logs back in and retries once per poll, so this recovers without a
restart.

## Metrics

Data comes from the UniFi Drive API. Four endpoints are polled per cycle:

| Endpoint | Yields |
|---|---|
| `/proxy/drive/api/v2/systems/device-info` | name, model, version, firmware, status, CPU load and temperature, memory |
| `/proxy/drive/api/v2/storage` | pools (capacity, usage, status, RAID groups) and disks (health score, temperature, power-on hours, RPM, size, bad and uncorrectable sectors, read/write KB/s) |
| `/proxy/users/drive/api/v2/drives` | per-share id, name, type, status, quota, usage, member count |
| `/proxy/drive/api/v2/systems/network-io` | receive and transmit KB/s |

An endpoint that fails is logged and its metrics omitted for that cycle; the rest are still
reported. A console that fails *every* endpoint is reported as an error, and one dead console
never suppresses the metrics of a healthy one.

### Prometheus

Metric names are prefixed with the configured namespace, so they read as
`unifi_unas_disk_temperature_celsius`, `unifi_unas_pool_usage_bytes`,
`unifi_unas_share_quota_bytes`, and so on. Note this differs from `unaspoller`, which uses a
bare `unas_` prefix — the Grafana dashboard in
[#785](https://github.com/unpoller/unpoller/issues/785) needs its queries adjusted
accordingly. Consistency with every other unpoller metric wins here.

### InfluxDB and DataDog

Four measurements / metric prefixes: `unas_device`, `unas_pool`, `unas_disk`, and
`unas_share` (DataDog also emits `unas_raid_group`). They are kept separate because each has
its own tag set; folding them together would give every point the union of those tags with
most of them empty.

Loki and OpenTelemetry are not wired up: Loki carries events only, and a UNAS console exposes
no event endpoints in this version.

## Credit

Endpoint discovery and the JSON shapes come from
[alexgreenbank/unaspoller](https://github.com/alexgreenbank/unaspoller) (MIT), which worked
out the UniFi Drive API by observing the console's own web UI. See
[#785](https://github.com/unpoller/unpoller/issues/785).

## Known gaps

These need data from real hardware:

- The `networkInterfaces` JSON key is an informed assumption — the reference implementation's
  struct tag for it was empty, so it was never actually verified. Interface data is not
  exported as metrics either way.
- `cacheSlots`, `expansions`, `usbs`, `riskReasons`, and `incompatibleReasons` are kept as
  raw JSON in the library because no populated example exists.
- `/proxy/drive/api/v2/systems/disk-stats` is not polled: it returns time series that do not
  map onto gauges, and its per-disk throughput is already in `/storage`.
