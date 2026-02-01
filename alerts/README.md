# UniFi Infrastructure Alerts

Example Prometheus and Loki alerting rules for monitoring UniFi infrastructure with unPoller.

## Overview

- **Prometheus** – Metrics from devices, clients, UPS/PDU, controller, sites, WAN, and more
- **Loki** – Logs for events, alarms, IDS, anomalies, and system logs

These examples assume the default Prometheus namespace `unpoller`. Adjust metric names if you use a custom `prometheus.namespace`.

## Prometheus Alerts

Place `prometheus/unifi-alerts.yaml` in your Prometheus `rule_files` or Grafana Alerting.

**Configuration example (prometheus.yml):**

```yaml
rule_files:
  - /etc/prometheus/rules/unifi-alerts.yaml
```

## Loki Alerts

Place `loki/unifi-alerts.yaml` in your Loki Ruler config. Loki must be run with the `-ruler.enable=true` flag and `-ruler.storage.path` configured.

**Configuration example (loki-config.yaml):**

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
