#!/usr/bin/env python3
"""
Headless endpoint discovery: navigate the UniFi controller UI and capture
all XHR/fetch requests (method, URL, optional headers), similar to manually
clicking around in Chrome and logging network traffic.

Usage:
  UNIFI_URL=https://192.168.1.1 UNIFI_USER=admin UNIFI_PASS=... python discover.py
  Or copy .env.example to .env and run (with python-dotenv): python discover.py

Output: API_ENDPOINTS_HEADLESS_<date>.md in current dir (or set OUTPUT_DIR)

Requires: pip install playwright && playwright install chromium
Optional: pip install python-dotenv  (to load .env)
"""

import os
from datetime import date
from pathlib import Path
from urllib.parse import urlparse

try:
    from dotenv import load_dotenv
    load_dotenv(Path(__file__).resolve().parent / ".env")
except ImportError:
    pass

from playwright.sync_api import sync_playwright

SCRIPT_DIR = Path(__file__).resolve().parent
# Default: write output in the same directory as the script (easy for users)
OUTPUT_DIR = Path(os.environ.get("OUTPUT_DIR", str(SCRIPT_DIR)))
HEADLESS = os.environ.get("HEADLESS", "true").lower() != "false"

GROUP_ORDER = ["api", "proxy-network-api", "proxy-network-v2", "proxy-users", "proxy-other", "other"]
GROUP_TITLES = {
    "api": "API (legacy)",
    "proxy-network-api": "Proxy /network API (v1)",
    "proxy-network-v2": "Proxy /network v2 API",
    "proxy-users": "Proxy /users API",
    "proxy-other": "Proxy (other)",
    "other": "Other",
}


def is_api_like(pathname: str) -> bool:
    return "/api" in pathname or "/proxy/" in pathname


def normalize_url(url_str: str) -> str:
    try:
        u = urlparse(url_str)
        return f"{u.scheme}://{u.netloc}{u.path or ''}{u.query and '?' + u.query or ''}"
    except Exception:
        return url_str


def run_group(pathname: str) -> str:
    if pathname.startswith("/api/"):
        return "api"
    if pathname.startswith("/proxy/network/api/"):
        return "proxy-network-api"
    if pathname.startswith("/proxy/network/v2/"):
        return "proxy-network-v2"
    if pathname.startswith("/proxy/users/"):
        return "proxy-users"
    if pathname.startswith("/proxy/"):
        return "proxy-other"
    return "other"


def main() -> None:
    base_url = (os.environ.get("UNIFI_URL") or "").rstrip("/")
    user = os.environ.get("UNIFI_USER") or ""
    password = os.environ.get("UNIFI_PASS") or ""

    if not base_url or not user or not password:
        print("Set UNIFI_URL, UNIFI_USER, and UNIFI_PASS (env or .env).", file=__import__("sys").stderr)
        raise SystemExit(1)

    try:
        origin = urlparse(base_url)
        our_origin = f"{origin.scheme}://{origin.netloc}"
    except Exception:
        our_origin = ""

    captured: dict[str, dict] = {}

    def on_request(request):
        if request.resource_type not in ("xhr", "fetch"):
            return
        url = request.url
        try:
            u = urlparse(url)
            req_origin = f"{u.scheme}://{u.netloc}"
            if req_origin != our_origin:
                return
            pathname = u.path or ""
            if not is_api_like(pathname):
                return
            key = f"{request.method} {normalize_url(url)}"
            if key not in captured:
                captured[key] = {
                    "method": request.method,
                    "url": normalize_url(url),
                    "pathname": pathname,
                    "request_headers": request.headers,
                }
        except Exception:
            pass

    with sync_playwright() as p:
        browser = p.chromium.launch(
            headless=HEADLESS,
            args=["--ignore-certificate-errors"],
        )
        context = browser.new_context(
            ignore_https_errors=True,
            user_agent="Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
        )
        page = context.new_page()
        page.on("request", on_request)

        print("Navigating to", base_url)
        try:
            page.goto(base_url, wait_until="networkidle", timeout=30000)
        except Exception:
            pass

        # Login
        if page.locator('input[type="password"]').count() > 0:
            print("Login form detected, submitting credentials...")
            page.locator('input[name="username"], input[name="email"], input[type="text"]').first.fill(user)
            page.locator('input[type="password"]').fill(password)
            page.locator('button[type="submit"], input[type="submit"], button:has-text("Login"), button:has-text("Sign in")').first.click()
            try:
                page.wait_for_load_state("networkidle")
            except Exception:
                pass
            page.wait_for_timeout(3000)

        # Visit common paths to trigger more API calls
        for path in ["/", "/devices", "/clients", "/settings", "/insights", "/topology", "/dashboard"]:
            try:
                page.goto(base_url + path, wait_until="domcontentloaded", timeout=10000)
                page.wait_for_timeout(2000)
            except Exception:
                pass

        browser.close()

    # Build markdown
    today = date.today().isoformat()
    entries = sorted(captured.values(), key=lambda e: e["url"])
    by_group: dict[str, list] = {}
    for e in entries:
        g = run_group(e["pathname"])
        by_group.setdefault(g, []).append(e)

    lines = [
        "# API Endpoints (headless discovery)",
        "",
        f"- **Date**: {today}",
        f"- **Controller**: {base_url}",
        f"- **Total unique requests**: {len(entries)}",
        "",
        "---",
        "",
    ]
    for g in GROUP_ORDER:
        list_ = by_group.get(g)
        if not list_:
            continue
        lines.append(f"## {GROUP_TITLES.get(g, g)}")
        lines.append("")
        for e in list_:
            u = urlparse(e["url"])
            path_only = (u.path or "") + (("?" + u.query) if u.query else "")
            lines.append(f"- `{e['method']} {path_only}`")
        lines.append("")

    lines.extend(["---", "", "## Sample request headers (first request)", ""])
    if entries and entries[0].get("request_headers"):
        lines.append("```")
        for k, v in entries[0]["request_headers"].items():
            kl = k.lower()
            if kl.startswith("x-") or kl in ("accept", "authorization"):
                lines.append(f"{k}: {v}")
        lines.append("```")

    out_path = OUTPUT_DIR / f"API_ENDPOINTS_HEADLESS_{today}.md"
    OUTPUT_DIR.mkdir(parents=True, exist_ok=True)
    out_path.write_text("\n".join(lines) + "\n", encoding="utf-8")
    print("Wrote", out_path, f"({len(entries)} unique endpoints)")


if __name__ == "__main__":
    main()
