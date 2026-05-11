# prometheus

This package provides the interface to turn UniFi measurements into prometheus
exported metrics. Requires the poller package for actual UniFi data collection.

## Scrape cache

Prometheus scrapes are served from an in-memory cache refreshed by a background
poller on a fixed interval. This decouples the scrape cadence from the UniFi
API call cadence: scrapes always return immediately, and upstream backpressure
(e.g. `429 Too Many Requests`) no longer stalls `/metrics`.

Config (TOML):

```toml
[prometheus]
  http_listen = "0.0.0.0:9130"
  # How often the background poller refreshes the cache served to /metrics.
  # Default: 60s. Values below 15s are clamped to 15s.
  interval = "60s"
```

Environment variable: `UP_PROMETHEUS_INTERVAL=60s`.

On poll error the last successful snapshot is preserved, so a transient 429 no
longer empties `/metrics`. To monitor cache staleness, scrape the
`unpoller_prometheus_cache_age_seconds` gauge — it reports seconds since the
last successful background refresh, or `-1` if no refresh has succeeded yet.

The `/scrape` endpoint (per-target dynamic scrapes) still fetches live but
coalesces concurrent requests for the same target via `singleflight`, so a
noisy scraper cannot multiply upstream load.
