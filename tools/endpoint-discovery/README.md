# UniFi endpoint discovery (headless)

Runs a headless browser against your UniFi controller, logs in, and records all API requests the UI makes. Use this to see which endpoints your controller exposes (e.g. when debugging 404s like [device-tags #935](https://github.com/unpoller/unpoller/issues/935)).

## What it does

1. Launches Chromium (headless by default).
2. Navigates to your UniFi controller URL.
3. If it sees a login form, fills username/password and submits.
4. Visits common UI paths to trigger API calls.
5. Captures every XHR/fetch request to `/api` or `/proxy/` (same origin).
6. Writes a markdown file: `API_ENDPOINTS_HEADLESS_YYYY-MM-DD.md` in this directory.

## Quick start

```bash
cd tools/endpoint-discovery
python3 -m venv .venv
.venv/bin/pip install playwright
.venv/bin/playwright install chromium
UNIFI_URL=https://YOUR_CONTROLLER UNIFI_USER=admin UNIFI_PASS=yourpassword .venv/bin/python discover.py
```

The script writes `API_ENDPOINTS_HEADLESS_<date>.md` in the same directory. Paste that file (or its contents) in an issue or comment when reporting which endpoints your controller supports.

## Options

- **Headed (see the browser):** `HEADLESS=false` before the command.
- **Output elsewhere:** `OUTPUT_DIR=/path/to/dir` before the command.
- **Use .env:** Copy `.env.example` to `.env`, set your values, then run `python discover.py` (install `python-dotenv` for .env support).

## Requirements

- Python 3.8+
- Direct controller URL (e.g. `https://192.168.1.1`) works best; same-origin requests are captured.
