## UnPoller InfluxDB  Plugin

Collects UniFi data from a UniFi controller using the API.

This is meant for InfluxDB users 1.8+ and 2.x series.

## Configuration

### InfluxDB 1.8+, 2.x

Note the use of `auth_token` to enable this mode.

```yaml
influxdb:
  disable: false
  # How often to poll UniFi and report to Datadog.
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
  # How often to poll UniFi and report to Datadog.
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
