---
title: Ingress Controller & L7 Route Rules
description: Path-based routing, Traefik ingress controller, and advanced traffic rules.
---

Codedock features a built-in **Ingress Controller** powered by Traefik, providing automatic TLS certificate issuance, path-based routing, URL rewriting, and custom HTTP middleware for your services.

## Ingress Architecture

Incoming web traffic is routed through the Codedock Ingress Controller before reaching containerized applications:

```
[ Client Request ]
       │
       ▼ (HTTPS :443 / HTTP :80)
┌────────────────────────────────────────────────────────┐
│               Codedock Ingress Controller              │
│       (Traefik Dynamic Routing & Let's Encrypt)        │
└──────────┬───────────────────────────────┬─────────────┘
           │ /api/*                        │ /*
           ▼                               ▼
┌──────────────────────┐       ┌──────────────────────┐
│    Backend Service   │       │   Frontend Service   │
│   (Node / Go / Py)   │       │     (Next.js / SPA)  │
└──────────────────────┘       └──────────────────────┘
```

## Adding a Domain & Ingress

1. In your Project, select the application or database service you wish to expose.
2. Navigate to the **Domains & Ingress** tab.
3. Enter your custom domain (e.g., `app.example.com`).
4. Select your target port (the internal port your container listens on, such as `3000` or `8080`).
5. Codedock automatically discovers ingress entrypoints and provisions Let's Encrypt TLS certificates using HTTP-01 challenges.

## L7 Path-Based Route Rules

You can route different path prefixes on the same domain to distinct backend services:

1. Open your Service and navigate to **Route Rules**.
2. Click **Add Route Rule**.
3. Configure the rule parameters:
   - **Path Prefix**: Specify the path pattern (e.g., `/api`, `/auth`, or `/static`).
   - **Target Service**: Choose which application service receives requests matching this path.
   - **Strip Prefix**: Toggle whether to remove the matched prefix before forwarding to the container (e.g., forwarding `/api/v1/users` as `/v1/users`).
   - **Priority**: Define execution order when multiple rules match a request.

## Custom Middleware & HTTP Headers

Codedock allows you to attach custom HTTP middleware to any ingress route:

### Request & Response Headers

Inject or strip headers for CORS, security policies, or backend identification:

- **Custom Request Headers**: Add headers before forwarding to containers (e.g., `X-Forwarded-Proto: https`).
- **Security Headers**: Automatically inject `HSTS`, `X-Frame-Options`, and `X-Content-Type-Options`.

### Basic Authentication

Protect sensitive endpoints or staging environments with HTTP Basic Authentication:

1. In Route Rules, enable **Basic Auth**.
2. Provide a username and secure password.
3. Codedock generates bcrypt-hashed credentials and attaches Traefik's `basicauth` middleware to the route.

### URL Redirects

Create HTTP 301 (Permanent) or 302 (Temporary) redirects directly at the edge:

- Forward apex domains to `www` (or vice-versa).
- Redirect legacy API paths to modern endpoints without executing application code.

## Automatic SSL & Certificate Verification

- **HTTP-01 Challenge**: Traefik automatically requests Let's Encrypt certificates once your DNS `A` or `CNAME` record points to your server's IP address.
- **Certificate Renewal**: Certificates are renewed automatically 30 days before expiration.
- **TLS Diagnostics**: The Ingress panel displays real-time SSL status, certificate expiry dates, and issuer validation.
