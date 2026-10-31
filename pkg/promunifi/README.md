# prometheus

This package provides the interface to turn UniFi measurements into prometheus
exported metrics. Requires the poller package for actual UniFi data collection.

## UniFi OS updates

For supported local session-authenticated consoles, two additional gauges expose
UniFi OS updates independently of the Network application:

| Metric (default namespace) | Labels | Meaning |
| --- | --- | --- |
| `unpoller_console_update_available` | `source`, `name` | `1` when an OS update is offered; `0` when explicitly up to date |
| `unpoller_console_info` | `source`, `name`, `version`, `available_version`, `firmware_status` | Firmware metadata; value always `1` |

There is one set of console readings per controller, not per Network site. Version
labels live on the info metric so version changes do not reset an alert's `for`
timer on the update gauge. Existing `controller_update_available` and
`device_upgradable` metrics retain their Network API meaning.

An offered version or `firmwareStatus: "updateAvailable"` produces `1`.
`firmwareStatus: "upToDate"` without an offered version produces `0`. Other states
without an offered version omit the update gauge. A missing/null offered version
alone does not prove the console is up to date. A valid partial response without
firmware entries produces neither metric. If the endpoint is unavailable,
forbidden, or cannot be decoded, both metrics are also absent for that poll;
other controller metrics continue to be collected.

For an update alert, use `unpoller_console_update_available == 1`. Treat an absent
series separately from `0`, and check `unpoller_prometheus_cache_age_seconds`
when using the scrape cache. API-key, remote, Protect-only, and application-update
collection are outside this initial scope.

## Scrape cache

By default, Prometheus scrapes are served from an in-memory cache refreshed by a
background poller on a fixed interval. This decouples the scrape cadence from the
UniFi API call cadence: scrapes always return immediately, and upstream backpressure
(e.g. `429 Too Many Requests`) no longer stalls `/metrics`. Set `interval = 0` to
disable the cache and restore on-demand `/metrics` fetches.

Config (TOML):

```toml
[prometheus]
  http_listen = "0.0.0.0:9130"
  # How often the background poller refreshes the cache served to /metrics.
  # Omitted: cache on, 60s. Set to 0 / "0" / "0s" to disable the cache so
  # /metrics fetches live. Values below 15s log a warning but are not clamped.
  interval = "60s"
```

Environment variable: `UP_PROMETHEUS_INTERVAL` (same contract; `0` / `0s` disables the cache).

On poll error the last successful snapshot is preserved, so a transient 429 no
longer empties `/metrics`. To monitor cache staleness, scrape the
`unpoller_prometheus_cache_age_seconds` gauge — it reports seconds since the
last successful background refresh, or `-1` if no refresh has succeeded yet.

The `/scrape` endpoint (per-target dynamic scrapes) still fetches live but
coalesces concurrent requests for the same target via `singleflight`, so a
noisy scraper cannot multiply upstream load.
