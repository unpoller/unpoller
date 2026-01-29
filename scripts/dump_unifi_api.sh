#!/usr/bin/env bash
#
# Dump raw JSON from UniFi Controller API endpoints to files.
# Uses unpoller -j "other <path>" for each path and saves to OUTDIR.
#
# Prerequisites: unpoller on PATH, valid config with controller auth.
#
# Usage:
#   ./scripts/dump_unifi_api.sh [-c CONFIG] [-s SITE] [-o OUTDIR]
#
# Options:
#   -c CONFIG   Config file (default: unpoller default locations)
#   -s SITE     Site name for /api/s/<site>/... paths (default: default)
#   -o OUTDIR   Output directory (default: ./api_dump)
#
# Examples:
#   ./scripts/dump_unifi_api.sh -c up.conf -o ./my_dump
#   SITE=my-site ./scripts/dump_unifi_api.sh -o ./api_dump
#

set -euo pipefail

CONFIG=""
SITE="${SITE:-default}"
OUTDIR="${OUTDIR:-./api_dump}"

while getopts "c:s:o:h" opt; do
  case "$opt" in
    c) CONFIG="$OPTARG" ;;
    s) SITE="$OPTARG" ;;
    o) OUTDIR="$OPTARG" ;;
    h) grep -E '^# (Usage|Options|Examples)' "$0" | sed 's/^# //'; exit 0 ;;
    *) exit 1 ;;
  esac
done

# Paths that need site substitution use %s
PATHS=(
  "/api/stat/sites"
  "/api/s/%s/stat/device"
  "/api/s/%s/stat/sta"
  "/api/s/%s/stat/event"
  "/api/s/%s/stat/rogueap"
  "/api/s/%s/stat/sitedpi"
  "/api/s/%s/stat/stadpi"
  "/api/s/%s/stat/alluser"
  "/api/s/%s/rest/networkconf"
  "/api/s/%s/list/alarm"
  "/api/s/%s/stat/ips/event"
  "/api/s/%s/stat/anomalies"
  "/api/s/%s/stat/admins"
  "/api/s/%s/stat/session"
  "/api/s/%s/stat/dashboard"
  "/api/s/%s/stat/health"
  "/v2/api/site/%s/aggregated-dashboard?historySeconds=3600"
)

# Optional: add traffic endpoints with fixed time window (last hour)
NOW_MS=$(($(date +%s) * 1000))
START_MS=$((NOW_MS - 3600000))
PATHS+=(
  "/v2/api/site/%s/traffic?start=${START_MS}&end=${NOW_MS}&includeUnidentified=false"
  "/v2/api/site/%s/country-traffic?start=${START_MS}&end=${NOW_MS}"
)

UNPOLLER="${UNPOLLER:-unpoller}"
if ! command -v "$UNPOLLER" &>/dev/null; then
  echo "error: $UNPOLLER not found (set UNPOLLER to path of unpoller binary)" >&2
  exit 1
fi

mkdir -p "$OUTDIR"
CONF_ARGS=()
if [[ -n "$CONFIG" ]]; then
  CONF_ARGS=(-c "$CONFIG")
fi

dump_one() {
  local path="$1"
  local sub
  sub=$(echo "$path" | sed "s|%s|$SITE|g")
  local fname
  fname=$(echo "$sub" | sed 's|^/||; s|[/?=&]|_|g')
  [[ -z "$fname" ]] && fname="root"
  fname="${fname}.json"
  local out="$OUTDIR/$fname"

  if out_err=$("$UNPOLLER" "${CONF_ARGS[@]}" -j "other $sub" 2>&1); then
    echo "$out_err" > "$out"
    echo "ok  $sub -> $out"
  else
    echo "fail $sub ($out_err)" >&2
  fi
}

echo "Dumping UniFi API responses to $OUTDIR (site=$SITE)"
for p in "${PATHS[@]}"; do
  dump_one "$p"
done
echo "Done. Output in $OUTDIR"
