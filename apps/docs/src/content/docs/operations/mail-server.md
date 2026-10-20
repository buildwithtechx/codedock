---
title: Codedock Mail Server
description: Complete self-hosted email infrastructure with Stalwart and SnappyMail.
---

Codedock features a dedicated **Mail Server Management** control plane that enables you to host your own secure, modern email infrastructure with automated DKIM signing, DNS verification, and webmail access.

## Architecture

Codedock pairs high-performance mail components with automated PaaS orchestration:

```
[ Inbound / Outbound Email ]
            │
            ▼ (Ports 25, 465, 587, 993)
┌───────────────────────────────────────┐
│        Stalwart Mail Server           │
│   (JMAP, IMAP4, SMTP, Sieve Filters)  │
└───────────────────┬───────────────────┘
                    │
                    ▼
┌───────────────────────────────────────┐
│       Codedock Mail Control Plane     │
│ (Domains, Mailboxes, DKIM, DNS Probes)│
└───────────────────┬───────────────────┘
                    │
                    ▼
┌───────────────────────────────────────┐
│         SnappyMail Webmail UI         │
│   (Modern, Lightweight Web Client)    │
└───────────────────────────────────────┘
```

## Adding a Mail Server

1. Deploy the Codedock Mail template from **Templates → Codedock Mail** or connect an external Stalwart instance.
2. In Codedock, navigate to **Mail → Servers** (or use `POST /api/mail/servers`).
3. Enter your mail server host, management endpoint, and admin credentials.
4. Test the connection to verify administrative API access.

## Configuring Domains & DKIM

1. Open your Mail Server and click **Add Domain** (`POST .../domains`).
2. Enter your email domain (e.g., `company.com`).
3. Click **Generate DKIM Key** (`POST .../dkim`). Codedock automatically generates a 2048-bit RSA key pair, configures Stalwart for outbound signing, and outputs the exact DNS TXT record.

## DNS Records Guidance

Codedock inspects your domain and provides the exact DNS records required for 100% email deliverability:

| Record Type | Host | Value / Target | Purpose |
| --- | --- | --- | --- |
| `MX` | `@` | `mail.example.com` (Priority: 10) | Directs inbound mail to your server. |
| `TXT` | `@` | `v=spf1 mx ~all` | Authorizes your server IP in SPF policy. |
| `TXT` | `codedock._domainkey` | `v=DKIM1; k=rsa; p=MIIBIjANB...` | Cryptographic signature verification. |
| `TXT` | `_dmarc` | `v=DMARC1; p=quarantine; rua=...` | DMARC compliance and reporting. |

Click **Run DNS Health Check** (`GET .../dns` and `GET .../checks`) to test DNS propagation. Codedock queries authoritative nameservers and validates each record in real time.

## Creating Mailboxes

1. In your domain view, click **Create Mailbox** (`POST .../mailboxes`).
2. Enter the email address (e.g., `alex@company.com`) and a secure password.
3. Users can immediately connect desktop and mobile mail clients (such as Apple Mail, Thunderbird, or Outlook) using standard IMAP/SMTP ports:
   - **Incoming (IMAP)**: Port 993 (SSL/TLS)
   - **Outgoing (SMTP)**: Port 465 (SSL/TLS) or 587 (STARTTLS)
4. Or log in directly via the integrated **SnappyMail Webmail** client.

## TLS Certificate Probing

To ensure mail encryption is functioning, click **Probe TLS Certificates** (`POST .../certs/probe`). Codedock validates that STARTTLS and implicit TLS handshakes succeed on ports 25, 465, 587, and 993 with valid Let's Encrypt certificates.
