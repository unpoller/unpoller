# Replacement for PR #936: Endpoint discovery via `--discover` (Go)

**Supersedes:** [PR #936](https://github.com/unpoller/unpoller/pull/936) (Python endpoint-discovery tool in `tools/endpoint-discovery/`).

PR #936 added a Python/Playwright tool to discover API endpoints by driving a headless browser. This replacement implements the same goal **inside unpoller** using the unifi library: no Python, no Playwright, same config as normal polling.

---

## Summary

- **Feature:** `--discover` flag on unpoller that probes known API endpoints on the controller and writes a shareable markdown report.
- **Credentials:** Uses the same config file (and first unifi controller) as normal unpoller runs.
- **Output:** Markdown file (default `api_endpoints_discovery.md`) with method, path, and HTTP status for each endpoint. Users can share this when reporting API/404 issues (e.g. [issue #935](https://github.com/unpoller/unpoller/issues/935)).
- **Dependency:** Requires [unpoller/unifi](https://github.com/unpoller/unifi) with `DiscoverEndpoints` and `Probe` (unifi PR/branch with discovery support).

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

**Removed (vs PR #936):** No `tools/endpoint-discovery/` (no Python, no Playwright).

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

## PR title (for replacement PR)

**Add `--discover` to probe API endpoints and write shareable report (replaces #936)**

---

## PR description (suggested)

**Replaces #936** (Python endpoint-discovery tool). Implements endpoint discovery inside unpoller using the unifi library; no Python or Playwright.

**What it does**
- `unpoller --discover` uses the first unifi controller from your config, authenticates, and probes a set of known API paths.
- Writes a markdown report (default `api_endpoints_discovery.md`) with method, path, and HTTP status for each endpoint.
- Same credentials as normal polling; users can share the file when reporting API/404 issues (e.g. #935).

**Usage**
```bash
unpoller --discover --config /path/to/up.conf --discover-output api_endpoints_discovery.md
```

**Requires** unpoller/unifi with `DiscoverEndpoints` (and `Probe`) merged or a compatible unifi version.

**Closes #936** (replaced by this approach).
