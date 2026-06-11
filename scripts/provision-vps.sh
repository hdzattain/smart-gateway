#!/usr/bin/env bash
set -euo pipefail

APP_DIR="${APP_DIR:-/root/new-api}"
DOMAIN="${DOMAIN:-smart-gateway.shop}"
GO_VERSION="${GO_VERSION:-1.25.1}"

if [ "$(id -u)" -ne 0 ]; then
  echo "Run as root" >&2
  exit 1
fi

apt-get update -y
DEBIAN_FRONTEND=noninteractive apt-get install -y curl ca-certificates unzip nginx git

if ! command -v go >/dev/null 2>&1 || ! go version | grep -q "go${GO_VERSION}"; then
  curl -fsSL "https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz" -o /tmp/go.tgz
  rm -rf /usr/local/go
  tar -C /usr/local -xzf /tmp/go.tgz
fi
export PATH="$PATH:/usr/local/go/bin:$HOME/.bun/bin"

if ! command -v bun >/dev/null 2>&1; then
  curl -fsSL https://bun.sh/install | bash
fi
export PATH="$PATH:$HOME/.bun/bin"

cd "$APP_DIR"

if [ ! -d web/node_modules ]; then
  cd "$APP_DIR/web"
  bun install
fi

if [ ! -f "$APP_DIR/web/default/dist/index.html" ]; then
  cd "$APP_DIR/web/default"
  export NODE_OPTIONS="${NODE_OPTIONS:---max-old-space-size=3072}"
  DISABLE_ESLINT_PLUGIN=true VITE_REACT_APP_VERSION=dev bun run build
fi

if [ ! -f "$APP_DIR/web/classic/dist/index.html" ]; then
  cd "$APP_DIR/web/classic"
  export NODE_OPTIONS="${NODE_OPTIONS:---max-old-space-size=3072}"
  VITE_REACT_APP_VERSION=dev bun run build
fi

cd "$APP_DIR"
go mod download
go build -ldflags "-s -w" -o "$APP_DIR/new-api-server" main.go

mkdir -p "$APP_DIR/logs"
cat >/etc/systemd/system/smart-gateway.service <<EOF
[Unit]
Description=Smart Gateway (New API)
After=network.target

[Service]
Type=simple
WorkingDirectory=$APP_DIR
Environment=PORT=3000
Environment=SQLITE_PATH=$APP_DIR/one-api.db
ExecStart=$APP_DIR/new-api-server
Restart=always
RestartSec=5
LimitNOFILE=65535
StandardOutput=append:$APP_DIR/logs/smart-gateway.log
StandardError=append:$APP_DIR/logs/smart-gateway.err.log

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable smart-gateway
systemctl restart smart-gateway

cp "$APP_DIR/deploy/nginx/smart-gateway.shop.conf" /etc/nginx/sites-available/smart-gateway.shop
ln -sf /etc/nginx/sites-available/smart-gateway.shop /etc/nginx/sites-enabled/smart-gateway.shop
rm -f /etc/nginx/sites-enabled/default
nginx -t
systemctl enable nginx
systemctl restart nginx

sleep 3
curl -fsS "http://127.0.0.1:3000/api/status" >/dev/null
curl -fsS "http://127.0.0.1/healthz" >/dev/null

echo "Smart Gateway provisioned. Point DNS for ${DOMAIN} to this VPS and open http://${DOMAIN}."
