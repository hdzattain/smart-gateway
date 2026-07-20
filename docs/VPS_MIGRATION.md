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
git clone https://github.com/hdzattain/smart-gateway.git /root/smart-gateway
cd /root/smart-gateway
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
cd /root/smart-gateway/web
bun install

cd /root/smart-gateway/web/default
export NODE_OPTIONS="--max-old-space-size=3072"
DISABLE_ESLINT_PLUGIN=true VITE_REACT_APP_VERSION=dev bun run build

cd /root/smart-gateway/web/classic
export NODE_OPTIONS="--max-old-space-size=3072"
VITE_REACT_APP_VERSION=dev bun run build

cd /root/smart-gateway
go mod download
go build -ldflags "-s -w" -o /root/smart-gateway/smart-gateway-server main.go
```

## 5. Restore runtime data

From the old VPS:

```bash
scp root@OLD_VPS_IP:/root/smart-gateway/smart-gateway.db /root/smart-gateway/smart-gateway.db
scp root@OLD_VPS_IP:/root/smart-gateway/.admin_password /root/smart-gateway/.admin_password
chmod 600 /root/smart-gateway/.admin_password
```

If starting fresh, run the setup wizard from the web UI and create a new admin.

## 6. Install systemd service

```bash
cat >/etc/systemd/system/smart-gateway.service <<'EOF'
[Unit]
Description=Smart Gateway (Smart Gateway)
After=network.target

[Service]
Type=simple
WorkingDirectory=/root/smart-gateway
Environment=PORT=3000
Environment=SQLITE_PATH=/root/smart-gateway/smart-gateway.db
ExecStart=/root/smart-gateway/smart-gateway-server
Restart=always
RestartSec=5
LimitNOFILE=65535
StandardOutput=append:/root/smart-gateway/logs/smart-gateway.log
StandardError=append:/root/smart-gateway/logs/smart-gateway.err.log

[Install]
WantedBy=multi-user.target
EOF

mkdir -p /root/smart-gateway/logs
systemctl daemon-reload
systemctl enable smart-gateway
systemctl restart smart-gateway
```

## 7. Configure Nginx

```bash
cp /root/smart-gateway/deploy/nginx/smart-gateway.shop.conf /etc/nginx/sites-available/smart-gateway.shop
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
