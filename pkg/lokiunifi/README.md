# lokiunifi

Loki Output Plugin for UnPoller

This plugin writes UniFi Events and IDS data to Loki. Maybe Alarms too.

Example Config:

```toml
[loki]
  # URL is the only required setting for Loki.
  url = "http://192.168.3.2:3100"

  # How often to poll UniFi and report to Loki.
  interval = "2m"

  # How long to wait for Loki responses.
  timeout = "5s"

  # Set these to use basic auth.
  #user = ""
  #pass = ""

  # Used for auth-less multi-tenant.
  #tenant_id = ""
```
