---
title: Software-Defined Networks & Watchdog
description: Multi-network topologies, isolated Docker networks, and atomic commit watchdogs.
---

Codedock provides **Software-Defined Networking (SDN)** capabilities that isolate container communications, secure lateral traffic between services, and protect server connectivity using an automated two-phase commit watchdog.

## Container Network Isolation

By default, every Codedock Project Environment receives an isolated, user-defined Docker bridge network:

- **Private DNS Resolution**: Containers within the same project communicate using service names as hostnames (e.g., `http://api:3000` or `postgres:5432`).
- **No Unintended Cross-Project Access**: Containers in Project A cannot reach containers in Project B unless explicitly connected through a shared network.
- **Port Masking**: Database ports (such as PostgreSQL 5432 or Redis 6379) are kept strictly internal to the container network and are never exposed to the public Internet unless Public Access is toggled on.

## Creating Multi-Network Topologies

For complex enterprise setups, you can define custom networks:

1. Navigate to **System Settings → Networks** (or use `POST /api/system/networks`).
2. Click **Create Network**.
3. Choose the driver and subnet parameters:
   - **Bridge**: Standard single-host isolated network.
   - **Overlay**: Multi-host distributed network for clustered environments.
4. Assign individual services to the new network. A single container can be attached to multiple networks simultaneously (e.g., a backend API connected to both an `edge-ingress` network and a private `database-tier` network).

## Network Access Revision Control

Control plane administrators can restrict inter-service communication:

- **Deny by Default**: Lock down sensitive database networks to only allow connections from designated API services.
- **Access Policies**: Define which IP CIDR ranges or subnets are permitted to reach container ports.

## Two-Phase Commit Watchdog

Modifying host network interfaces, firewall rules, or Docker daemon bridges carries the risk of accidentally severing control plane SSH or API access. Codedock prevents permanent lockouts using a **Two-Phase Commit Watchdog**:

```
1. Prepare Changes ──────► 2. Arm Watchdog Timer ──────► 3. Apply Network Changes
                                                                 │
                          ┌──────────────────────────────────────┴──────────────────────────────────────┐
                          ▼ Verified Connection                                 ▼ Severed Connection
                  4. Explicit Commit                                    4. Automated Rollback
            (Changes Persisted Permanently)                          (Revert to Known-Good Network State)
```

### How the Watchdog Works

1. **Watchdog Arming**: When a network topology or interface update is requested, Codedock initiates a timed watchdog operation (`GET /api/system/networks/operations/:operationId/watchdog`).
2. **Temporary Application**: Changes are applied in a temporary, testable state.
3. **Connectivity Verification**: The control plane attempts end-to-end health checks across the modified network.
4. **Automatic Reversion**: If the node fails to check in before the watchdog timer expires (default: 60 seconds), the daemon automatically rolls back the network configuration to its previous state (`POST .../watchdog/rollback`), restoring server access.
5. **Permanent Commit**: If connectivity checks succeed, the administrator or automated test confirms the operation (`POST .../watchdog/commit`), permanently persisting the new configuration.
