# unpoller

![Version: 2.11.2](https://img.shields.io/badge/Version-2.11.2-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: v2.11.2](https://img.shields.io/badge/AppVersion-v2.11.2-informational?style=flat-square)

A Helm chart for unpoller, a unifi prometheus exporter

**Homepage:** <https://github.com/unpoller/unpoller>

## Maintainers

| Name | Email | Url |
| ---- | ------ | --- |
| XXXX? | <XXXX?> |  |

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

