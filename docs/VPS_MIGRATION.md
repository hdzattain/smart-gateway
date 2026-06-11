# VPS Migration Guide

This guide moves Smart Gateway to a new VPS without committing secrets to GitHub.

## 1. Prepare the new VPS

Use Ubuntu 24.04 or similar.

```bash
apt-get update -y
apt-get install -y git curl ca-certificates nginx
```

## 2. Clone the repository

```bash
git clone https://github.com/hdzattain/smart-gateway.git /root/new-api
cd /root/new-api
```

## 3. Install runtimes

```bash
# Go 1.25.1
curl -fsSL https://go.dev/dl/go1.25.1.linux-amd64.tar.gz -o /tmp/go.tgz
rm -rf /usr/local/go
tar -C /usr/local -xzf /tmp/go.tgz
export PATH=$PATH:/usr/local/go/bin

# Bun
apt-get install -y unzip
curl -fsSL https://bun.sh/install | bash
export PATH=$PATH:$HOME/.bun/bin
```

## 4. Build production assets

Low-memory VPS instances may need swap:

```bash
fallocate -l 3G /swapfile2 || dd if=/dev/zero of=/swapfile2 bs=1M count=3072
chmod 600 /swapfile2
mkswap /swapfile2
swapon /swapfile2
```

Build:

```bash
cd /root/new-api/web
bun install

cd /root/new-api/web/default
export NODE_OPTIONS="--max-old-space-size=3072"
DISABLE_ESLINT_PLUGIN=true VITE_REACT_APP_VERSION=dev bun run build

cd /root/new-api/web/classic
export NODE_OPTIONS="--max-old-space-size=3072"
VITE_REACT_APP_VERSION=dev bun run build

cd /root/new-api
go mod download
go build -ldflags "-s -w" -o /root/new-api/new-api-server main.go
```

## 5. Restore runtime data

From the old VPS:

```bash
scp root@OLD_VPS_IP:/root/new-api/one-api.db /root/new-api/one-api.db
scp root@OLD_VPS_IP:/root/new-api/.admin_password /root/new-api/.admin_password
chmod 600 /root/new-api/.admin_password
```

If starting fresh, run the setup wizard from the web UI and create a new admin.

## 6. Install systemd service

```bash
cat >/etc/systemd/system/smart-gateway.service <<'EOF'
[Unit]
Description=Smart Gateway (New API)
After=network.target

[Service]
Type=simple
WorkingDirectory=/root/new-api
Environment=PORT=3000
Environment=SQLITE_PATH=/root/new-api/one-api.db
ExecStart=/root/new-api/new-api-server
Restart=always
RestartSec=5
LimitNOFILE=65535
StandardOutput=append:/root/new-api/logs/smart-gateway.log
StandardError=append:/root/new-api/logs/smart-gateway.err.log

[Install]
WantedBy=multi-user.target
EOF

mkdir -p /root/new-api/logs
systemctl daemon-reload
systemctl enable smart-gateway
systemctl restart smart-gateway
```

## 7. Configure Nginx

```bash
cp /root/new-api/deploy/nginx/smart-gateway.shop.conf /etc/nginx/sites-available/smart-gateway.shop
ln -sf /etc/nginx/sites-available/smart-gateway.shop /etc/nginx/sites-enabled/smart-gateway.shop
rm -f /etc/nginx/sites-enabled/default
nginx -t
systemctl enable nginx
systemctl restart nginx
```

## 8. Point DNS

Create or update DNS:

- `A smart-gateway.shop -> NEW_VPS_IP`
- `A www.smart-gateway.shop -> NEW_VPS_IP`

## 9. Verify

```bash
curl http://smart-gateway.shop/healthz
curl http://smart-gateway.shop/api/status
systemctl is-active smart-gateway nginx
```
