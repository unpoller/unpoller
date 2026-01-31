# Endpoint discovery via `--discover` (Go)

The [unpoller](https://github.com/unpoller/unpoller) repo has a Python/Playwright endpoint-discovery tool in `tools/endpoint-discovery/` (browser-based; discovers unknown endpoints; more support currently). This doc describes the Go-based discovery: it lives in [unpoller/unifi](https://github.com/unpoller/unifi) and unpoller, probes a fixed list of known API paths, and discovers/confirms known endpoints only.

---

## Summary

- **Feature:** `--discover` flag on unpoller that probes known API endpoints on the controller and writes a shareable markdown report.
- **Credentials:** Uses the same config file (and first unifi controller) as normal unpoller runs.
- **Output:** Markdown file (default `api_endpoints_discovery.md`) with method, path, and HTTP status for each endpoint. Users can share this when reporting API/404 issues (e.g. [issue #935](https://github.com/unpoller/unpoller/issues/935)).
- **Dependency:** Requires [unpoller/unifi](https://github.com/unpoller/unifi) with `DiscoverEndpoints` and `Probe` merged.

---

## Changes in unpoller (this PR)

| File | Change |
|------|--------|
| `pkg/poller/config.go` | Add `Discover bool`, `DiscoverOutput string` to `Flags`. |
| `pkg/poller/start.go` | Register `--discover` and `--discover-output`; when `--discover`, call `RunDiscover()` and exit. |
| `pkg/poller/commands.go` | Add `RunDiscover()`: load config, init inputs, find input implementing `Discoverer`, call `Discover(outputPath)`. |
| `pkg/poller/inputs.go` | Add optional interface `Discoverer` with `Discover(outputPath string) error`. |
| `pkg/inputunifi/discover.go` | **New file.** Implement `Discoverer`: first controller, authenticate, get sites, call `c.Unifi.DiscoverEndpoints(site, outputPath)`. |
| `.gitignore` | Add `up.discover-test.json`, `api_endpoints_discovery.md`. |

---

## Usage

```bash
# Same config as normal unpoller (first unifi controller is used)
unpoller --discover --config /path/to/up.conf --discover-output api_endpoints_discovery.md
```

If config is in the default search path:

```bash
unpoller --discover --discover-output api_endpoints_discovery.md
```

---

## PR title

**Add `--discover` to probe API endpoints and write shareable report**

---

## PR description (suggested)

Go-based endpoint discovery: probes known API paths on the controller and writes a shareable report. Uses the [unifi](https://github.com/unpoller/unifi) library; same config as normal polling. The Python tool in `tools/endpoint-discovery/` remains for browser-based discovery (more coverage).

**What it does**
- `unpoller --discover` uses the first unifi controller from your config, authenticates, and probes a set of known API paths.
- Writes a markdown report (default `api_endpoints_discovery.md`) with method, path, and HTTP status for each endpoint.
- Same credentials as normal polling; users can share the file when reporting API/404 issues (e.g. #935).

**Usage**
```bash
unpoller --discover --config /path/to/up.conf --discover-output api_endpoints_discovery.md
```

**Requires** [unpoller/unifi](https://github.com/unpoller/unifi) with `DiscoverEndpoints` (and `Probe`) merged.

**CI:** Merge unifi first, then this PR (or update go.mod to require the new unifi release). For local testing, use `replace github.com/unpoller/unifi/v5 => ../unifi` with the unifi repo checked out.
