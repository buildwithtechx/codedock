---
title: Worker Servers
description: Connect and manage Linux worker hosts over direct SSH.
---

Codedock connects to remote workers using SSH. The current worker path does not implement a Yamux tunnel, Cloudflare SSH transport or a jump-host connection.

## Prepare a worker

Use a compatible Linux host with Docker, adequate storage and network access to the images and repositories your workloads need. Allow the control plane to reach its SSH port. Choose an SSH user with the required Docker and host-management permissions.

Keep private keys secure and restrict access to server management. Docker access is powerful host access; choose credentials accordingly.

## Add a server

Open server management and enter a name, host, SSH port, SSH user and authentication method. Use a private key or password according to your installation's supported settings. Test the connection before saving or assigning workloads.

The API provides `POST /api/servers/test-ssh`, `POST /api/servers`, and `PATCH /api/servers/:id`. Direct SSH is the supported transport. A private worker still needs a routable connection from the control plane; there is no automatic reverse tunnel.

## Workload placement

Assign the server where the project or service configuration supports it. Verify container startup, image access, DNS and application connectivity on that host. Remote placement does not automatically make local-only database query paths reachable.

## Inspect and maintain

Use server metrics and service logs to investigate failures. Maintain host security updates, Docker, disk capacity and independent backups. Before deleting a server record, plan how its workloads and data will be moved or removed.

Cloud plan limits apply according to the hosted account. Self-hosted worker setup does not require a Stripe key or a separate paid tunnel.

See [Canvas](/fleet-management/canvas/) and [observability](/operations/observability/).
