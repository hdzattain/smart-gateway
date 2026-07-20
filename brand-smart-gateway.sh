#!/usr/bin/env bash
# Post-deploy branding for Smart Gateway.
# Applies white-label branding via the SUPPORTED options API (no schema edits,
# no source patching). Upstream "Smart Gateway" attribution is preserved per
# Apache-2.0 license and the project's AGENTS.md Rule 5.
#
# Usage:
#   BASE_URL=https://your-host ADMIN_USER=admin ./brand-smart-gateway.sh
set -euo pipefail

BASE_URL="${BASE_URL:-http://localhost:3000}"
ADMIN_USER="${ADMIN_USER:-admin}"
ADMIN_PASSWORD="${ADMIN_PASSWORD:?set ADMIN_PASSWORD}"
COOKIES="$(mktemp)"

echo ">> Logging in as ${ADMIN_USER} ..."
ADMIN_ID="$(curl -s -c "$COOKIES" -X POST "$BASE_URL/api/user/login" \
  -H 'Content-Type: application/json' \
  -d "{\"username\":\"$ADMIN_USER\",\"password\":\"$ADMIN_PASSWORD\"}" \
  | python3 -c 'import sys,json;print(json.load(sys.stdin)["data"]["id"])')"
echo "   admin id=$ADMIN_ID"

set_opt() {
  curl -s -b "$COOKIES" -H "Smart-Gateway-User: $ADMIN_ID" -H 'Content-Type: application/json' \
    -X PUT "$BASE_URL/api/option/" -d "$1" >/dev/null && echo "   set: $1"
}

echo ">> Applying branding (attribution preserved) ..."
set_opt '{"key":"SystemName","value":"Smart Gateway"}'
set_opt '{"key":"Footer","value":"© 2026 Smart Gateway. Powered by <a href=\"https://github.com/hdzattain/smart-gateway\">Smart Gateway</a> (open-source, Apache-2.0)."}'
set_opt '{"key":"Announcements","value":"[{\"content\":\"Welcome to Smart Gateway.\",\"type\":\"default\",\"publishDate\":\"2026-06-11\"}]"}'

echo ">> Done. Verify:"
curl -s "$BASE_URL/api/status" | python3 -c 'import sys,json;d=json.load(sys.stdin)["data"];print("   system_name:",d["system_name"]);print("   footer:",d["footer_html"])'
rm -f "$COOKIES"
