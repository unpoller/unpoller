# inputunifi

## UnPoller Input Plugin

Polls UniFi controllers and hands their metrics and events to every configured output.
All configuration lives under `[unifi]` — see the commented `[[unifi.controller]]` block in
[`examples/up.conf.example`](../../examples/up.conf.example) for every available option.

## Protect-only consoles (UNVR)

A **UNVR** or **UNVR Pro** runs UniFi Protect with *no Network application installed*.
UnPoller's normal startup probes the Network API to read the controller version, which on
these appliances returns the UniFi OS SPA HTML rather than JSON:

```
[ERROR] Controller 3 of 3 Auth or Connection Error, retrying: unifi controller:
        unable to get server version: invalid character '<' looking for beginning of value
```

Set `disable_network = true` on that controller. UnPoller then skips the Network API
entirely and collects only UniFi Protect:

```toml
[[unifi.controller]]
  url                  = "https://unvr.example.com"
  # Optional: a local read-only account. Only needed for save_protect_logs, which uses the
  # legacy Protect endpoints and authenticates with a session cookie.
  user                 = "unpoller"
  pass                 = "unpoller"
  # Required. Mint this in Protect under Settings -> Control Plane -> Integrations.
  protect_api_key      = "unifiprotectapikey"

  disable_network      = true
  save_protect_devices = true
  save_protect_logs    = false
  verify_ssl           = false
```

As an environment variable this is `UP_UNIFI_CONTROLLER_0_DISABLE_NETWORK=true`.

### What is and isn't collected

| | With `disable_network = true` |
| --- | --- |
| Protect devices — cameras, sensors, lights, bridges, link stations, NVR | ✅ `save_protect_devices` |
| Protect event logs | ✅ `save_protect_logs` |
| Sites, clients, devices, DPI, traffic, rogue APs, speed tests | ❌ never polled |
| Events, syslog, alarms, anomalies, IDs | ❌ never polled |

The Network-only `save_*` options are ignored rather than honoured, so leaving them at their
defaults is fine. UnPoller logs an error at startup if `disable_network` is set with neither
`save_protect_devices` nor `save_protect_logs` — that combination collects nothing at all.

A console that runs *both* applications (a UDM, UCG, or a UniFi OS Server with Protect
installed) should leave `disable_network` at its default of `false` and simply set
`save_protect_devices = true`. This flag is only for consoles with no Network application.

Mixing is fine: a Protect-only console is configured as one more `[[unifi.controller]]`
alongside your normal ones.

See [unpoller/unpoller#1066](https://github.com/unpoller/unpoller/issues/1066).
