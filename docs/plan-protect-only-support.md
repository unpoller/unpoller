# Plan: Protect-only console support (issue #1066)

v5 added UniFi Protect device metrics, but they cannot be collected from a **Protect-only
console** — a UNVR or UNVR Pro, which runs UniFi Protect with no Network application
installed. The Protect collectors are complete; only controller initialisation blocks them.

## Design decision: a `disable_network` flag, not a separate input plugin

UNAS Pro got its own plugin (`pkg/inputunas`, see [plan-unas-support.md](plan-unas-support.md))
because it shares nothing with a UniFi controller: its own credentials, its own host, its own
JSON API. A Protect-only console is the opposite case. It is reached at the same URL, with the
same `Controller` config, by the same `*unifi.Unifi` client, and `collectProtect` /
`collectProtectLogs` already live in `inputunifi` and are already not site-scoped. A second
plugin would duplicate the controller config block to gate two calls it already makes.

So: a per-controller `disable_network` flag, defaulting to `false`, as
[@platinummonkey asked for on the issue](https://github.com/unpoller/unpoller/issues/1066).

Detection is by explicit flag only. The issue documents two unauthenticated probes that
identify a UNVR (`/api/system` reporting `hardware.shortname`, and `/proxy/network/status`
answering HTML rather than a 401). Auto-detection was cut: it adds network calls and cached
state to every startup, and the flag has to exist as an override regardless.

### The blockers

1. **`unifi.NewUnifi()` cannot construct a client for the console at all.** It ends in
   `GetServerData()` → GET `APIStatusPath` (`/status`), which `path()` rewrites to
   `/proxy/network/status`. With no Network application to route to, UniFi OS serves its own
   SPA HTML and the call fails on its first byte:
   `invalid character '<' looking for beginning of value`. `getUnifi()` treats any non-429
   error as fatal, so the entry never even prints a config summary.
2. **`pollController` aborts on `getFilteredSites`** long before it reaches `collectProtect`.
   `collectControllerEvents` aborts in the same place, before `collectProtectLogs`.
3. **`Metrics()` counts a poll successful only if it produced devices or clients.** A
   Protect-only console produces neither, so a filtered scrape of one — the Prometheus
   per-target path — falls through to the dynamic-controller branch and reports
   `ErrDynamicLookupsDisabled` despite a successful collection.

What already works in our favour: `path()` passes anything starting with `/proxy/` through
untouched, so the Protect Integration paths need no changes; `collectProtect` and
`collectProtectLogs` take no `sites` argument; and every output plugin already handles
`ProtectDevices`.

## Part 1 — `../unifi` (github.com/unpoller/unifi/v6)

`NewProtectClient(config *Config) (*Unifi, error)` in `protect.go`, modelled on
`NewUNASClient` in `unas.go`, which exists for the same underlying reason.

- Rejects a nil config, and a config with neither `ProtectAPIKey` nor `APIKey`:
  Integration/v1 is `X-API-Key` only and has no cookie fallback.
- Sets `u.new = true` directly rather than probing with `checkNewStyleAPI()`. A Protect
  console is always a UniFi OS console, so `Login()` resolves to `/api/auth/login` without
  depending on how the console answers a GET of `/`.
- Logs in **only** when `APIKey == "" && User != ""`. Integration/v1 needs no session, but the
  legacy Protect endpoints (`GetProtectLogs`, `GetProtectEventThumbnail`) authenticate with a
  session cookie. Skipping login when `APIKey` is set is required for correctness: `Login()`
  routes through `/status` in that case, the one endpoint this console cannot serve.
- Ends by probing `/v1/meta/info`, the way `NewUnifi` ends by probing Network, so a caller
  that gets no error has a console it can really poll.
- Returns a plain `*Unifi` rather than a wrapper type, unlike `NewUNASClient`: the Protect
  getters are already `*Unifi` methods.

**`ServerStatus` must be populated, and that is not cosmetic.** `Unifi` embeds
`*ServerStatus`, so leaving it nil turns any caller's `u.ServerVersion` into a nil
dereference — including `inputunifi`'s config summary. `ServerVersion` holds the *Protect*
application version, since a Protect-only console has no Network version to report.

## Part 2 — `unpoller`

### 2a. Config

`Controller.DisableNetwork *bool`, tagged for json/toml/xml/yaml, defaulting to `false` in
`setDefaults` and inheritable from `[unifi.defaults]` via `setControllerDefaults`, matching
every neighbouring flag. Added to `formatControllers` so the web UI reflects it.

### 2b. Skipping the Network application

| Site | Change |
|---|---|
| `getUnifi` | Calls `unifi.NewProtectClient` instead of `unifi.NewUnifi`, inside the unchanged 429-retry loop. |
| `Initialize`, `DebugInput` | Skip `checkSites`. |
| `pollController` | The Network pass is extracted into a new `pollNetwork(c, sites, m) error` and skipped wholesale, leaving site discovery and `collectProtect` behind. |
| `collectControllerEvents` | Skips site discovery and reduces the collector list to `collectProtectLogs`, the only site-independent one. |
| `Metrics` | Counts `ProtectDevices` toward a successful poll. |
| `RawMetrics` | Answers the raw-path kind and rejects the site-scoped kinds with `ErrNetworkDisabled`, rather than returning a confusing empty result. |
| `logController` | Marks the mode and omits the Network-only lines. |

`extractDevices` also gains a nil guard on `metrics.Devices`, which it dereferenced
unconditionally — a latent panic in its own right, now reachable whenever the Network pass
does not run.

### 2c. Warnings, not failures

`warnProtectOnly` logs an error at startup for the two configurations that can never collect
anything: `disable_network` with neither Protect save flag, and `save_protect_devices` with no
`protect_api_key` or `api_key`. Neither is fatal — silently collecting nothing is the failure
mode hardest to diagnose from a log, so it is called out rather than acted on.

### 2d. Docs and tests

`pkg/inputunifi/README.md` gains a "Protect-only consoles (UNVR)" section; the three
`examples/up.*.example` files gain `disable_network`, defaulting off.

`pkg/inputunifi` had no tests before this change. The new `input_test.go` follows
`pkg/inputunas/input_test.go`: external test package, an `httptest` fake UNVR that serves the
console's SPA HTML for everything but the Protect paths and the login, and a prose comment
above each test. It covers initialisation, metrics, events, the filtered scrape, `RawMetrics`,
both warnings, config binding across toml/json/yaml/env, and that the shipped examples do not
disable the Network application. `TestProtectOnlyControllerFailsWithoutFlag` pins the original
bug against the same fake console.

## Sequencing across the two repos

`unpoller` consumes `unpoller/unifi` as a tagged module, so `NewProtectClient` has to exist
and be released before anything can call it. Part 1 lands and is tagged first; Part 2 stays a
draft until then, developed against a local `go.work` workspace so `go.mod` is never
temporarily rewritten.

## Open items

- Auto-detection of Protect-only consoles, if operators find the flag a papercut.
- The Protect bootstrap API (`/proxy/protect/api/bootstrap`) carries richer per-camera data —
  `isRecording`, `isConnected`, NVR `version` — but it is private and undocumented, so it is
  out of scope here as it was in #1015.
