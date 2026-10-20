---
title: Push Notifications & Mobile Alerts
description: Real-time push notifications via Expo, device tokens, and webhook alerts.
---

Codedock provides a comprehensive, multi-channel alerting system that delivers instant notifications to your mobile devices, team chat apps, and custom webhooks when important platform events occur.

## Supported Notification Channels

- **Mobile Push (iOS & Android)**: Real-time push notifications sent directly to your phone via the Codedock Mobile app powered by Expo.
- **Discord**: Rich embed messages with direct links to failing services or logs.
- **Slack**: Webhook alerts delivered to your team's incident or engineering channels.
- **Email**: Transactional email alerts sent via SMTP or Resend.
- **Generic Webhooks**: Raw JSON event payloads delivered to any HTTP endpoint or incident management tool (PagerDuty, Opsgenie).

## Mobile Device Token Registration

When you sign in to the Codedock Mobile app on your iOS or Android device:
1. The app requests push notification permissions and retrieves an Expo Push Token.
2. The token is registered with your Codedock account (`POST /api/notifications/devices`).
3. You can inspect or revoke active mobile device tokens at any time under **Profile → Mobile Devices**.

## Event Subscriptions

You can customize which events trigger notifications and route different events to specific channels:

| Event Category | Severity | Description |
| --- | --- | --- |
| `deployment_success` | Info | Notifies when a build and rollout finishes successfully. |
| `deployment_failure` | Critical | Alerts when a container build fails or crash-loops on startup. |
| `server_offline` | Critical | Triggered when a worker node stops responding to control plane heartbeats. |
| `resource_attention` | Warning | Triggered on OOMKilled events or abnormal container exit codes. |
| `disk_threshold` | Warning | Alerts when a host server crosses the configured disk usage threshold (e.g. 80%). |
| `backup_completed` | Info | Confirms that scheduled database backups to S3 or R2 completed. |

Configure subscriptions under **Server Settings → Notifications → Subscriptions** or via `PUT /api/notifications/subscriptions`.

## Delivery Logs & Diagnostics

Codedock logs every notification attempt to provide full visibility into delivery reliability:
- View delivery timestamps, target channels, HTTP status codes, and latency under **Notifications → Deliveries** (`GET /api/notifications/deliveries`).
- If an external webhook or push endpoint encounters a temporary network timeout, Codedock automatically retries with exponential backoff.
- Click **Send Test Notification** (`POST /api/settings/notifications/test`) to verify your channel configurations before going live.
