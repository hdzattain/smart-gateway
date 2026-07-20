# Payment Integration API for Customer Engineering Teams

This document describes the payment integration surface for Smart Gateway.

Base URL:

```text
http://smart-gateway.shop
```

Authentication:

- User-facing payment initiation endpoints require a logged-in web session.
- Webhook endpoints are public but verified by the configured provider signature/secret.
- Admin configuration endpoints require the root/admin account.

## Supported payment rails

Smart Gateway already includes these payment providers:

- Stripe: recommended for credit cards.
- Waffo / Waffo Pancake: card and wallet checkout rails, depending on provider account configuration.
- Creem: checkout/product-based payment.
- EPay: callback-style payment integration.
- Balance payment for subscriptions.

## Required admin setup

Log in as `admin`, then configure payment options in the system settings.

Critical setup items:

1. Confirm payment compliance from the admin options endpoint/UI.
2. Configure provider credentials.
3. Configure subscription plans or top-up amounts.
4. Register the webhook URL at the provider dashboard.

Admin-only compliance endpoint:

```http
POST /api/option/payment_compliance
Cookie: session cookie
Smart-Gateway-User: <admin_user_id>
```

## Stripe credit-card integration

### Webhook URL

Configure this endpoint in Stripe Dashboard:

```text
http://smart-gateway.shop/api/stripe/webhook
```

If/when HTTPS is enabled on the Smart Gateway service, switch to:

```text
https://smart-gateway.shop/api/stripe/webhook
```

### User top-up amount calculation

```http
POST /api/user/self/stripe/amount
Content-Type: application/json
Cookie: user session
Smart-Gateway-User: <user_id>
```

Typical request body:

```json
{
  "amount": 10
}
```

### User top-up checkout session

```http
POST /api/user/self/stripe/pay
Content-Type: application/json
Cookie: user session
Smart-Gateway-User: <user_id>
```

Typical request body:

```json
{
  "amount": 10,
  "top_up_count": 1
}
```

The response returns checkout/session data used by the frontend to redirect the user to the payment page.

### Subscription checkout

```http
POST /api/subscription/stripe/pay
Content-Type: application/json
Cookie: user session
Smart-Gateway-User: <user_id>
```

Typical request body:

```json
{
  "plan_id": 1
}
```

Plans are created and managed from admin APIs/UI:

```http
GET  /api/subscription/admin/plans
POST /api/subscription/admin/plans
PUT  /api/subscription/admin/plans/:id
PATCH /api/subscription/admin/plans/:id
```

## Waffo / Waffo Pancake integration

### Webhooks

```text
POST http://smart-gateway.shop/api/waffo/webhook
POST http://smart-gateway.shop/api/waffo-pancake/webhook/prod
POST http://smart-gateway.shop/api/waffo-pancake/webhook/test
```

### Top-up amount and checkout

```http
POST /api/user/self/waffo/amount
POST /api/user/self/waffo/pay
POST /api/user/self/waffo-pancake/amount
POST /api/user/self/waffo-pancake/pay
```

### Subscription checkout

```http
POST /api/subscription/waffo-pancake/pay
```

### Admin helper endpoints

```http
POST /api/option/waffo-pancake/pair
POST /api/option/waffo-pancake/save
POST /api/option/waffo-pancake/subscription-product
POST /api/option/waffo-pancake/subscription-product-options
```

These endpoints require root/admin authentication.

## Creem integration

### Webhook

```text
POST http://smart-gateway.shop/api/creem/webhook
```

### User top-up checkout

```http
POST /api/user/self/creem/pay
```

### Subscription checkout

```http
POST /api/subscription/creem/pay
```

## EPay integration

### Top-up notify/return endpoints

```text
POST http://smart-gateway.shop/api/user/epay/notify
GET  http://smart-gateway.shop/api/user/epay/notify
```

### Subscription notify/return endpoints

```text
POST http://smart-gateway.shop/api/subscription/epay/notify
GET  http://smart-gateway.shop/api/subscription/epay/notify
GET  http://smart-gateway.shop/api/subscription/epay/return
POST http://smart-gateway.shop/api/subscription/epay/return
```

### User top-up checkout

```http
POST /api/user/self/pay
POST /api/user/self/amount
```

### Subscription checkout

```http
POST /api/subscription/epay/pay
```

## Recommended customer engineering handoff

Ask the customer payment team for:

- Stripe publishable key and secret key.
- Stripe webhook signing secret.
- Allowed success/cancel URLs.
- Card settlement currency.
- Tax/invoice requirements.
- Production vs sandbox account separation.

Operational test plan:

1. Configure provider credentials in admin UI.
2. Create a low-value test top-up plan.
3. Start checkout from a non-admin customer account.
4. Complete sandbox card payment.
5. Confirm webhook is received.
6. Confirm quota/subscription is credited.
7. Confirm idempotent handling by replaying a webhook from the provider dashboard.

## Security notes

- Do not expose provider secret keys in frontend code or GitHub.
- Restrict admin configuration APIs to admin users only.
- Use provider webhook signature verification.
- Move to HTTPS before processing production credit-card traffic.
- Rotate demo accounts before production onboarding.
