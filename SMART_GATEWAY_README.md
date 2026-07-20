# Smart Gateway Deployment Guide

This repository contains the source and non-sensitive deployment assets for the Smart Gateway service.

## What is included

- Upstream application source based on hdzattain/smart-gateway.
- Docker deployment template: `docker-compose.smart-gateway.yml`.
- Runtime branding helper: `brand-smart-gateway.sh`.
- Customer/admin handoff guide: `docs/CUSTOMER_HANDOFF.md`.
- VPS migration guide: `docs/VPS_MIGRATION.md`.

## What is intentionally excluded

The following files are runtime/private and must not be committed:

- `.env`
- `.admin_password`
- `smart-gateway.db` / `*.db`
- `logs/`
- `data/`
- generated binaries such as `smart-gateway-server`
- TLS private keys and certificates
- frontend `node_modules/` and build caches

## Current production endpoint

- HTTP endpoint: `http://smart-gateway.shop`
- Health endpoint: `http://smart-gateway.shop/healthz`

Port `443` is intentionally reserved for an existing x-ui/xray service on the current VPS, so Smart Gateway is exposed over port `80` through Nginx.

## Quick start on a new VPS

```bash
git clone https://github.com/hdzattain/smart-gateway.git /root/smart-gateway
cd /root/smart-gateway
bash scripts/provision-vps.sh
```

Then restore runtime data from the old host if needed:

```bash
scp root@OLD_HOST:/root/smart-gateway/smart-gateway.db /root/smart-gateway/smart-gateway.db
scp root@OLD_HOST:/root/smart-gateway/.admin_password /root/smart-gateway/.admin_password
chmod 600 /root/smart-gateway/.admin_password
systemctl restart smart-gateway nginx
```

For Docker-based deployment, use:

```bash
docker compose -f docker-compose.smart-gateway.yml up -d --build
```

## License and attribution

This deployment preserves upstream open-source attribution and license notices. Do not remove protected upstream identifiers, license files, or notices.
