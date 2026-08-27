## UnPoller InfluxDB  Plugin

Collects UniFi data from a UniFi controller using the API.

This supports InfluxDB 1.x, 2.x, and 3.x.

## Configuration

### InfluxDB 3.x

Set `version: 3` and provide a token plus database name. For InfluxDB 3 Core/Enterprise,
leave `use_v2_api` unset (defaults to the native v3 write API). For InfluxDB Cloud
Serverless or Clustered, set `use_v2_api: true`.

```yaml
influxdb:
  disable: false
  version: 3
  interval: "30s"
  url: http://influxdb3:8181
  auth_token: somesecret
  database: unifi
  verify_ssl: false
```

See [MIGRATION.md](MIGRATION.md) for schema changes affecting Grafana dashboards.

### InfluxDB 1.8+, 2.x

Note the use of `auth_token` to enable v2 mode when `version` is omitted.

```yaml
influxdb:
  disable: false
  # How often to poll UniFi and report to InfluxDB.
  interval: "2m"
  # the influxdb url to post data
  url: http://somehost:1234
  # the secret auth token, this enables InfluxDB 1.8, 2.x compatibility.
  auth_token: somesecret
  # the influxdb org
  org: my-org
  # the influxdb bucket
  bucket: my-bucket
  # how many points to batch write per flush.
  batch_size: 20
```

### InfluxDB pre 1.8

Note the lack of `auth_token` to enable this mode.

```yaml
influxdb:
  disable: false
  # How often to poll UniFi and report to InfluxDB.
  interval: "2m"
  # the influxdb url to post data
  url: http://somehost:1234
  # the database
  db: mydb
  # the influxdb api user
  user: unifi
  # the influxdb api password 
  pass: supersecret
```

### Global Tags

Tags configured under `tags` are attached to every measurement written to
InfluxDB. Per-metric tags (site, device id, etc.) take precedence on key
collision so they cannot be overwritten by a misconfigured global tag.

```yaml
influxdb:
  tags:
    customer: abc_corp
    env: prod
```

Equivalent TOML:

```toml
[influxdb.tags]
  customer = "abc_corp"
  env      = "prod"
```
