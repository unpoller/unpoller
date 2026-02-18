## Summary

Fixes crashes and poor behavior when using **Remote API** (UniFi Site Manager / api.ui.com) with an API key—e.g. a “Fabric” or cloud API key—for authentication. We ran into 429 rate limiting, 403s from NVR/non-Network consoles, and **nil pointer panics** in `updateWeb` / `controllerID` during or after 429 handling.

## Issue: Remote API with API key (e.g. Fabric key)

With **remote API mode** enabled and an API key (e.g. from unifi.ui.com / Fabric):

1. **429 Too Many Requests** – Discovery and polling hit the cloud API quickly. The API returns 429 with `Retry-After` (e.g. 1m0s). We either failed immediately or retried once with a short cap, so some controllers were never discovered or we kept hitting 429.
2. **403 on non-Network consoles** – The `/v1/hosts` list includes NVRs, CloudKey+ for Displays, etc. Calling the Network “sites” API on those returns 403. We were calling every console, causing 403 noise and wasting rate limit.
3. **Nil pointer panic in `updateWeb`** – When the API returned 429 during polling, re-auth or internal handling could set `c.Unifi = nil`. We used `defer updateWeb(c, m)`, so when the deferred call ran, `c` or `c.Unifi` could be nil and `controllerID(c)` panicked (`addr=0x28` = dereferencing `c.Unifi` when `c` was nil or `c.Unifi` was nil). The process crashed repeatedly under 429.

## What was required to fix it

### In **unpoller/unifi** (library, separate PR)

- **429**: Retry up to 3 times using the **full** `Retry-After` from the API (no 15s cap) so rate-limited requests can succeed.
- **NVR/Protect filtering**: Added `FilterNetworkConsoles` and `DiscoverNetworkConsoles()` so we only discover consoles that support the Network API; skip names containing `nvr`, `protect`, `cloudkey+ for displays` to avoid 403s and save rate limit.
- **Error response safety**: Cap read size for 4xx/5xx responses and truncate body in error messages to avoid OOM from huge HTML/error bodies.

### In **unpoller** (this PR)

- **Discovery**: Use `DiscoverNetworkConsoles()` and retry each console up to 3 times with a short delay. After 3 failures (429 or 403), log once: `Excluding controller <name>: api key does not have permissions (after 3 attempts)` and stop trying that controller.
- **updateWeb panic**:
  - Removed `defer updateWeb(c, m)`. Call `updateWeb(c, m)` only on the **success path** at the end of `pollController`, guarded by `c != nil && c.Unifi != nil`.
  - In `updateWeb`: guard on `c == nil`, `metrics == nil`, `c.Unifi == nil`; add `defer recover()` to log and swallow any remaining panic so one bad path doesn’t kill the process.
  - In `controllerID`: safe nil checks and single dereference.
  - In `formatSites` / `formatClients` / `formatDevices`: skip nil slice elements.
  - In `pollController`: wrap the `updateWeb` call in a `recover()` so even an old image or race doesn’t crash the poller.
- **Dockerfile.local**: Add a Dockerfile that builds unpoller with a **local** unifi library (build context = parent dir with `unifi/` and `unpoller/`). Uses `go mod edit -replace` and `go build` inside the image so the container is built entirely from local repos for testing.

## Testing

- Ran unpoller in Kubernetes with remote API and an API key; multiple Network and NVR consoles in the account.
- Confirmed: NVRs excluded from discovery, 429 retries with full backoff, no more `controllerID`/`updateWeb` panics when 429 occurs during polling. Excluded controllers get a single log line after 3 attempts.

## Dependencies

- Requires **unpoller/unifi** changes (429 retry, `DiscoverNetworkConsoles`, error body handling). Until that library is released, build unpoller with a local unifi clone and `Dockerfile.local` (or `go mod replace`).

## CI / Workflows

**Workflows are expected to fail on this PR until [unpoller/unifi#200](https://github.com/unpoller/unifi/pull/200) is merged and a new unifi release is available.** This PR calls `DiscoverNetworkConsoles()` and uses 429/error behavior that exist only in the unifi library PR. The current `go.mod` points at the published unifi module, which does not yet provide `DiscoverNetworkConsoles`, so `go build` and `go test` in CI will fail with an undefined method error. Once unifi is updated and unpoller’s `go.mod` is updated to that version (or the unifi PR is merged and a release cut), CI should pass.
