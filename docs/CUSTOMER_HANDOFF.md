# Customer Handoff Guide

## Production URL

- Website: `http://smart-gateway.shop`
- Health check: `http://smart-gateway.shop/healthz`

## Administrator account

- Username: `admin`
- Password: stored on the production VPS at `/root/smart-gateway/.admin_password`

For security, the real administrator password is not stored in this repository. Retrieve it on the production server:

```bash
ssh root@45.78.79.80
cat /root/smart-gateway/.admin_password
```

Change the password after handoff from the web console.

## Demo customer account

- Username: `reseller1`
- Password: `Cust#Pass2026`
- Role: standard customer/reseller

This demo account is for verification only. Delete it or rotate its password before onboarding real customers.

## Basic administrator workflow

1. Open `http://smart-gateway.shop`.
2. Click **Sign in**.
3. Log in as `admin` using the password from `/root/smart-gateway/.admin_password`.
4. Configure upstream channels in the admin console.
5. Create customer/reseller users.
6. Assign model groups and quotas to customers.
7. Customers create and manage their own API tokens from their console.

## Verified permission boundary

The customer/reseller account cannot access sensitive admin-only upstream configuration:

- `/api/channel/` returns `Unauthorized, insufficient privileges` for standard customer users.
- `/api/user/` returns `Unauthorized, insufficient privileges` for standard customer users.
- `/api/token/` works for standard customer users and returns only their own token data.

## Current deployment notes

- Nginx listens on port `80` and proxies to `127.0.0.1:3000`.
- Port `443` is intentionally reserved for an existing `/usr/local/x-ui/bin/xray-linux-amd64` service on the current VPS.
- Streaming/SSE-compatible Nginx settings are enabled:
  - `proxy_buffering off`
  - `proxy_request_buffering off`
  - `X-Accel-Buffering: no`
  - long read/send timeouts

## Operational commands

```bash
systemctl status smart-gateway
systemctl status nginx
journalctl -u smart-gateway -n 100 --no-pager
curl http://smart-gateway.shop/healthz
```

## Important security tasks before real customer traffic

1. Rotate the `admin` password.
2. Delete or rotate the demo `reseller1` account.
3. Configure upstream providers and billing carefully.
4. Add backups for `/root/smart-gateway/smart-gateway.db` or migrate to MySQL/PostgreSQL.
5. If HTTPS is required on this service, move the existing x-ui/xray process off port `443` or use a separate origin/IP.
