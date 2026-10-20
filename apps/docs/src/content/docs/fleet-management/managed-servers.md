---
title: Managed Cloud Servers (Hetzner)
description: Automated cloud server provisioning and fleet scaling with Hetzner Cloud.
---

Codedock provides native integration with **Hetzner Cloud** to automatically provision, configure, and scale Linux worker servers with zero manual SSH setup.

## Overview

While standard worker hosts require you to manually bring a Linux server with Docker and configure SSH access, Managed Servers automate the entire lifecycle through the cloud provider API.

When you add a Managed Server:
1. Codedock requests a new cloud VPS from Hetzner Cloud.
2. An isolated, high-security SSH key pair is dynamically generated.
3. The server boots with a cloud-init bootstrap script that installs Docker, sets up firewall rules, and configures daemon networking.
4. The server automatically joins your Codedock worker fleet and begins accepting container workloads.

## Prerequisites

1. Create a Hetzner Cloud account and generate an API Token with read & write permissions in the Hetzner Cloud Console.
2. In the Codedock dashboard, navigate to **Admin Settings → Cloud Integrations**.
3. Enter your Hetzner Cloud API Token and save.

## Provisioning a Managed Server

1. Open **Fleet Management → Servers**.
2. Click **Add Server** and choose **Managed Cloud (Hetzner)**.
3. Configure the server parameters:
   - **Server Name**: A distinct identifier for your worker (e.g., `worker-eu-01`).
   - **Location**: Select your preferred datacenter (`nbg1` Nuremberg, `fsn1` Falkenstein, `hel1` Helsinki, or `ash` Ashburn US).
   - **Server Type**: Choose your compute specifications (e.g., `cx22` for standard workloads, `cpx31` for high-concurrency builds).
   - **Image**: Ubuntu LTS is selected and optimized by default.
4. Click **Provision Server**.

## Provisioning Stages

Codedock provides live status updates as the server progresses through its initialization stages:

| Stage | Status | Description |
| --- | --- | --- |
| 1 | `initializing` | Generating cryptographic keys and validating API quotas. |
| 2 | `creating` | Dispatching the VPS creation request to Hetzner Cloud. |
| 3 | `booting` | Waiting for IPv4 assignment and initial OS boot. |
| 4 | `installing` | Running the automated cloud-init bootstrap script to install Docker and system dependencies. |
| 5 | `online` | Control plane connects over secure SSH, performs health checks, and marks the server ready for workloads. |

You can inspect the live installation log at any time during provisioning by clicking **View Boot Logs** in the server details panel.

## Automated Credentials & Security

Managed Servers handle key management automatically:
- **SSH Key Injection**: The public key is registered in Hetzner Cloud and injected directly into the root user's authorized keys during boot.
- **Runtime API Credentials**: Codedock generates dedicated worker API tokens for mutual authentication between the daemon and remote hosts.
- **Credential Rotation**: You can rotate the worker's runtime credentials or SSH access with one click via **Server Settings → Rotate Credentials**.

## Workload Placement & Maintenance

Once a managed server transitions to `online`:
- Select the server as the target host in project environments or individual service configurations.
- View real-time CPU, RAM, and disk utilization via live WebSocket streams.
- If you delete a managed server from Codedock, you are prompted to choose whether to retain or permanently terminate the underlying Hetzner VPS to avoid unexpected cloud billing.
