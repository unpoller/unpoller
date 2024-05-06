# unpoller

![Version: 2.11.2](https://img.shields.io/badge/Version-2.11.2-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: v2.11.2](https://img.shields.io/badge/AppVersion-v2.11.2-informational?style=flat-square)

A Helm chart for unpoller, a unifi prometheus exporter. This chart helps deploy Unpoller (unifi metrics exporter)
in kubernetes clusters.
It crates a Deployment to run the unpoller container, confiuration is stored in a ConfigMap and mounted in the container.
It supports integration with Prometheus operator, so a PodMonitor is created that will scrape the Deployment for the metrics.
Optionally, it can deploy automatically the dashboards into a Grafana instance through the integration with GrafanaOperator:
* Creates a Grafana CR with the credentials provided (or reuses existing Grafana object)
* Creates a Dashboard instance for all the unpoller provided charts.
See Readme.MD for details, and values.yaml for all the configuration options.

See further documentation in how to install unpoller in Kubernetes in https://unpoller.com/PATH_TBD (will be updated)

**Homepage:** <https://unpoller.com/>

## Source Code

* <https://github.com/unpoller/unpoller>

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| affinity | object | `{}` |  |
| dashboards.create | bool | `true` |  |
| dashboards.grafana.create | bool | `true` |  |
| dashboards.grafana.secret.existingSecretName | string | `""` |  |
| dashboards.grafana.secret.password | string | `"prom-operator"` |  |
| dashboards.grafana.secret.username | string | `"prom"` |  |
| dashboards.grafana.selectorLabels | object | `{}` |  |
| dashboards.grafana.url | string | `""` |  |
| fullnameOverride | string | `""` |  |
| image.pullPolicy | string | `"IfNotPresent"` |  |
| image.repository | string | `"ghcr.io/unpoller/unpoller"` |  |
| image.tag | string | `"v2.11.2"` |  |
| imagePullSecrets | list | `[]` |  |
| livenessProbe.httpGet.path | string | `"/"` |  |
| livenessProbe.httpGet.port | string | `"http"` |  |
| nameOverride | string | `""` |  |
| nodeSelector | object | `{}` |  |
| podAnnotations | object | `{}` |  |
| podLabels | object | `{}` |  |
| podSecurityContext | object | `{}` |  |
| readinessProbe.httpGet.path | string | `"/"` |  |
| readinessProbe.httpGet.port | string | `"http"` |  |
| replicaCount | int | `1` |  |
| resources | object | `{}` |  |
| securityContext | object | `{}` |  |
| serviceAccount.annotations | object | `{}` |  |
| serviceAccount.automount | bool | `true` |  |
| serviceAccount.create | bool | `true` |  |
| serviceAccount.name | string | `""` |  |
| tolerations | list | `[]` |  |
| upConfig | string | `"[poller]\n    debug = false\n    quiet = false\n    plugins = []\n[prometheus]\n  disable = false\n  http_listen = \"0.0.0.0:9130\"\n  report_errors = false\n[influxdb]\n  disable = true\n[unifi]\n    dynamic = false\n[loki]\n    disable = true\n[[unifi.controller]]    \n    url         = \"https://unifi.home:8443\"\n    user        = \"unifi\"\n    pass        = \"unifi\"\n    sites       = [\"all\"]\n    save_ids    = true\n    save_dpi    = true\n    save_sites  = true\n    hash_pii    = false\n    verify_ssl  = false\n"` |  |

