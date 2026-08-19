# Plan: opt-in UNAS Pro support (issue #785)

Reference implementation: https://github.com/alexgreenbank/unaspoller (MIT, David Newhall II
copyright header — same lineage as unpoller, so the code is cleanly upstreamable). Credit
alexgreenbank in the package header and README.

## Design decision: separate input plugin, not a `save_unas` toggle

UNAS Pro is a standalone UniFi OS console with **no Network application**: no sites, no
`/status`, its own credentials, its own host. `inputunifi`'s sites → devices → per-site flow
has nothing to offer it, and folding it in would mean bolting a second credential set onto
`Controller`. So: a new `pkg/inputunas` plugin with its own `[unas]` config block, off by
default and a no-op when unconfigured.

### The one real blocker

`unifi.NewUnifi()` is unusable against a UNAS target. It ends in `GetServerData()` → GET
`APIStatusPath` (`/status`), which `path()` rewrites to `/proxy/network/status` when
`u.new` is set. A storage-only console has no Network app to answer that, and unaspoller's
endpoint list never touches `/status`, so we have no evidence it responds. API-key auth is
out for the same reason — `Login()` uses `/status` as its login path when `APIKey != ""`.

Fix by construction, not by testing against hardware we don't have: a login-only
constructor that stops after `Login()`.

What already works in our favour: `path()` (types.go:234) deliberately passes anything
starting with `/proxy/` through untouched, so `GetData`/`GetJSON` work on `/proxy/drive/…`
with zero changes. UniFi OS login, the cookie jar, csrf capture
(`x-csrf-token`/`x-updated-csrf-token`), and 429 `Retry-After` handling all already exist.

## Scope for v1

**In** — the four endpoints unaspoller actually polls:

| Endpoint | Yields |
|---|---|
| `/proxy/drive/api/v2/systems/device-info` | name, model, version, firmware, uptime, cpu load/temp, memory free/total/avail |
| `/proxy/drive/api/v2/storage` | pools (capacity/usage/status/raid groups) and disks (health score, temp, power-on hours, rpm, size, bad/uncorrectable sectors, read/write KBPS) |
| `/proxy/users/drive/api/v2/drives` | per-share id/name/type/status/quota/usage/member count |
| `/proxy/drive/api/v2/systems/network-io` | receive/transmit KBPS |

**Explicitly cut from v1**, to be revisited once community data exists:

- `/proxy/drive/api/v2/systems/disk-stats` — needs `?start=&end=&interval=`, returns
  time-series arrays that don't map onto gauges, and its per-disk `readKBPS`/`writeKBPS`
  are already in `/storage`.
- Cache slots, expansions (unaspoller has them as empty placeholder structs).
- `/proxy/users/drive/api/v1/systems/{identity,info}`, `/proxy/users/drive/api/v2/{groups,storage}`,
  `/proxy/drive/api/v1/systems/performance/file-operations`.

## Part 1 — `../unifi` (github.com/unpoller/unifi/v5) — **done**

Landed as `unas.go` + `unas_test.go`, four path constants in `types.go`, two sentinels
(`ErrNilConfig`, `ErrAPIKeyUnsupported`) in `unifi.go`, synthetic fixtures under
`endpoints_data/unas-*.json`, and a UNAS section in `endpoints_data/ENDPOINTS.md`.
`go test ./...` and `golangci-lint run ./...` are clean.

Refinements over the plan below:

1. `NewUNASClient` sets `u.new = true` directly instead of calling `checkNewStyleAPI()`. A
   UNAS Pro is always a UniFi OS console, so probing a GET of `/` only adds a way for login
   to resolve to the wrong path. This removes the "does the console answer GET `/`?" open
   item entirely.
2. **`GetUNASDevice` does not return an error on partial failure.** Failing endpoints are
   logged to `ErrorLog` and their fields left nil; only a total failure returns an error
   (and then a nil device). Returning `(device, err)` together was the original design, but
   it made the reflexive `if err != nil { return }` in `pkg/inputunas` silently drop every
   metric from a console whenever one endpoint 404s — exactly the case the partial path
   exists to serve. The footgun is removed at the source rather than documented around.
3. `NewUnifi` now returns `ErrNilConfig` too, instead of its own `fmt.Errorf("config is
   nil")` for the identical condition. Callers gain `errors.Is`; the lib stops having two
   errors for one condition in one file.

### 1a. `unas.go`

Path constants alongside the existing `/proxy/protect/…` block in `types.go` (or local to
`unas.go`, matching however `protect_log.go` does it):

```go
APIUNASDeviceInfoPath = "/proxy/drive/api/v2/systems/device-info"
APIUNASStoragePath    = "/proxy/drive/api/v2/storage"
APIUNASDrivesPath     = "/proxy/users/drive/api/v2/drives"
APIUNASNetworkIOPath  = "/proxy/drive/api/v2/systems/network-io"
```

Login-only client:

```go
type UNASClient struct{ *Unifi }

func NewUNASClient(config *Config) (*UNASClient, error) // newUnifi → checkNewStyleAPI → Login, stop
```

`newUnifi` is unexported but we're in-package; this is ~10 lines lifted from `NewUnifi`
plus the SSL-fingerprint loop, minus `GetServerData()`. Reject a non-empty `APIKey` with a
clear error (`ErrAPIKeyUnsupported` or reuse an existing sentinel) since key auth routes
through `/status`.

Getters on `*UNASClient` — `GetUNASDeviceInfo`, `GetUNASStorage`, `GetUNASDrives`,
`GetUNASNetworkIO`, plus one aggregate `GetUNASDevice()` returning the composed
`*UNASDevice`. **Do not add these to the `UnifiClient` interface**: CONVENTIONS.md requires
`mocks.MockUnifi` implement every method, so widening it breaks the mocks build for no
gain — `inputunifi` holds `*unifi.Unifi` directly and nothing consumes the interface here.

### 1b. Structs — port with two bug fixes

Port unaspoller's `drivetypes.go`, renaming `DriveApiV2*` → `UNAS*` to match lib style.
Two dead fields to fix rather than inherit:

1. `DriveApiV2SystemsDeviceInfo.NetworkInterfaces` has an **empty** tag (`json:""`) — it
   never unmarshals today. Confirm the real key (likely `networkInterfaces`) from a probe
   dump before relying on it; leave it out of v1 metrics either way.
2. `DriveApiV2Pool.activeRaidGroupId` is **unexported** with a json tag — also dead.
   Export as `ActiveRaidGroupID`.

Use `FlexInt`/`FlexBool` for every numeric/boolean field. unaspoller's own comments flag
several as "float64? only ever seen ints" (temperatures, sector counts) — the flex types
are the lib convention and pre-empt exactly the unmarshal breakage its author predicted.

### 1c. Tests

JSON fixtures under `endpoints_data/`, table tests in the existing lib style.

## Part 2 — `unpoller` — **done**

Landed as `pkg/inputunas/{config.go,input.go,input_test.go,README.md}`, the `UNASDevices`
pair in `pkg/poller`, a blank import in `main.go`, `unas.go` in promunifi/influxunifi/
datadogunifi with their dispatch wiring, `pkg/promunifi/unas_test.go`, and an `[unas]` section
in all three example configs. `go test ./...` and `golangci-lint run` are clean.

Refinements over the plan below:

1. **Opt-in is `enable`, defaulting to false — not `disable`.** A field named `disable`
   zero-values to false, so it cannot make a plugin default-off and reads as a double
   negative; `enable bool` defaults to off and says what it does. Two further guards:
   `Initialize` returns silently when the device list is empty, and unlike `inputunifi` it
   does **not** synthesize a default device from `[unas.defaults]` — doing so would poll a
   host the operator never named. Devices configured while `enable` is false log one error,
   since that combination is always a mistake.
2. **`Metrics` returns `(metrics, nil)` whenever any console was collected.** Not politeness:
   `poller.collectMetrics` uses `if result.err != nil { ... } else if result.metric != nil`,
   so returning both discards every metric in that cycle. One dead console out of three would
   have thrown away the two that worked. (`inputunifi.Metrics` hits this at
   `interface.go:303`; left alone as out of scope.)
3. **Re-auth fires on any total failure, not on a 401.** A mid-session 401 from `GetData`
   surfaces as `ErrInvalidStatusCode`, not `ErrAuthenticationFailed` (`unifi.go:603`), so
   there is no sentinel to match. Session expiry fails all four endpoints at once, which is
   exactly the total-failure case, and one wasted re-login against an unreachable console is
   cheaper than never recovering from an expired session.
4. **Env var names come from the `xml` tags, not `json`** — verified by binding a real config:
   `UP_UNAS_DEFAULT_TIMEOUT` (singular) and `UP_UNAS_DEVICE_0_URL`.

### Original Part 2 plan

### 2a. `pkg/inputunas/`

`config.go` + `input.go` + `README.md`. Config shape:

```toml
[unas]
  enable = true
  [unas.defaults]
    user       = "unpoller"
    pass       = ""
    verify_ssl = false
    timeout    = "60s"
  [[unas.device]]
    url = "https://192.168.1.10"
    user = "..."
    pass = "..."
```

`json`/`toml`/`xml`/`yaml` tags on every field; env vars map as `UP_UNAS_*` via cnfg.
Implement all five `poller.Input` methods — `Events` returns empty, `RawMetrics` and
`DebugInput` as stubs; `Metrics` returns `&poller.Metrics{UNASDevices: …}`.

**Do not re-log endpoint failures.** `GetUNASDevice` already reports each failing endpoint
through `ErrorLog` once per poll — accepted deliberately: a 404ing endpoint is something an
operator should see, and `inputunifi` logs per-cycle problems the same way. The plugin must
not add a second layer of the same message.

**Own re-auth-on-401 loop.** We do not inherit `inputunifi`'s. UNAS session tokens expire
around 2h, and re-login must re-attach fresh cookies — this is the bug alexgreenbank burned
days on. One re-login attempt per request, then fail.

### 2b. Plumbing (the easy-to-miss half)

3. `main.go` — blank import `_ "github.com/unpoller/unpoller/pkg/inputunas"`.
4. `pkg/poller/config.go` — add `UNASDevices []any` to `Metrics` **and** the matching
   append in `AppendMetrics` (`pkg/poller/inputs.go` ~336). Omitting the second is silent
   metric loss when two inputs run, and nothing errors.

**One aggregate family, not five.** `UNASDevices []any` carrying `*unifi.UNASDevice`
(device-info + pools + disks + drives + net-io) keeps `AppendMetrics` to one line and lets
one export function per output emit every gauge family. Nothing needs pools or disks
correlated independently — no site name, no PII redaction — so flat-per-entity families buy
nothing and cost 5× the plumbing.

### 2c. Outputs — follow the `ups.go` pattern

- `pkg/promunifi/unas.go` + four edits in `collector.go`: descriptor struct field (~71),
  `descUNASDevice()` call (~321), the `Describe` reflect list (~607), the collect loop (~836).
- `pkg/influxunifi/unas.go` + the type-switch loop in `influxdb.go`.
- `pkg/datadogunifi/unas.go` + `datadog.go` loop and type-switch (~476).
- `pkg/lokiunifi` — skip, events only.
- `pkg/otelunifi` — **skip in v1.** `report.go:178` walks only `m.Devices`; it's per-family
  and UPSDevices is already absent there. Adding UNAS is a separate follow-up.

**Namespace divergence, stated deliberately:** promunifi prefixes with `u.Namespace + "_"`,
so these land as `unifi_unas_disk_temperature`, not unaspoller's `unas_disk_temperature`.
The Grafana dashboard in issue #785 is built on the latter and will need its queries
adjusted. Consistency with every other unpoller metric wins.

### 2d. Docs

`examples/up.conf.example` (new `[unas]` section after `[unifi]` at line 133),
`up.json.example`, `up.yaml.example`, `pkg/inputunas/README.md`, and the `UP_UNAS_*` env
mapping.

## Sequencing across the two repos

`go.mod` pins `github.com/unpoller/unifi/v5 v5.30.0` with **no** `replace` directive
(CLAUDE.md's `/Users/briangates/unifi` reference is stale). So:

1. Land the `../unifi` change, tag a release.
2. Bump `unifi/v5` in unpoller, land `pkg/inputunas`.

Both steps are done: Part 1 released as **v5.31.0**, and unpoller's `go.mod` now pins that
version. Because the tag is published and proxy-resolvable, the `replace` directive this plan
originally called for was never needed — which removes the "strip the replace before merging"
footgun entirely. The pin bump also pulled `testify` to v1.12.0 and `golang.org/x/net` to
v0.58.0, both required by the new lib version.

## Open items needing real hardware or community data

- Real json key and shape for `networkInterfaces`, `usbs`, `sfpAggregation`. The port assumes
  `networkInterfaces` (the reference implementation's empty tag means this was never
  verified); if a live console disagrees, that one tag is the fix.
- Shapes for `cacheSlots`, `riskReasons`, `incompatibleReasons`, `expansions` — unaspoller
  left these as empty placeholders. Out of v1 scope.
