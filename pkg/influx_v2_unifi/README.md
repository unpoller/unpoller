# influx_v2_unifi

## UnPoller InfluxDB v2 Plugin

Collects UniFi data from a UniFi controller using the API.

This is meant for InfluxDB users 1.8+ and 2.x series.

## Configuration

```yaml
influxdb2:
  # to enable this
  enable: true
  # How often to poll UniFi and report to Datadog.
  interval: "2m"
  # the influxdb url to post data
  url: http://somehost:1234
  # the secret auth token
  auth_token: somesecret
  # the influxdb org
  org: my-org
  # the influxdb bucket
  bucket: my-bucket
  # how many points to batch write per flush.
  batch_size: 20
```
