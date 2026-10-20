---
title: High Availability & Clustering
description: Multi-node control planes, Patroni PostgreSQL clustering, and automated failover.
---

For mission-critical production environments, Codedock can be deployed in a **High Availability (HA)** architecture that eliminates single points of failure across the control plane, database, and ingress layers.

## High Availability Architecture

A standard Codedock instance runs the daemon, Traefik, and metadata database on a single host. In an HA setup, components are decoupled into redundant layers:

```
                  ┌───────────────────────────────┐
                  │    External Load Balancer     │
                  │   (DNS Failover / Anycast)    │
                  └──────────────┬────────────────┘
                                 │
                 ┌───────────────┴───────────────┐
                 │                               │
        ┌────────▼────────┐             ┌────────▼────────┐
        │  Codedock Node  │             │  Codedock Node  │
        │ (Active Leader) │             │ (Standby Node)  │
        └────────┬────────┘             └────────┬────────┘
                 │                               │
                 └───────────────┬───────────────┘
                                 │
                  ┌──────────────▼────────────────┐
                  │    Patroni PostgreSQL Cluster │
                  │  (Synchronous Replication HA) │
                  └───────────────────────────────┘
```

## Patroni PostgreSQL HA Cluster

Codedock stores all system metadata, project configs, audit logs, and deployment states in PostgreSQL. In HA mode, PostgreSQL is managed by **Patroni**:

- **Distributed Consensus**: Patroni uses Raft consensus via etcd or Consul to maintain cluster topology and elect the primary database.
- **Synchronous Replication**: Writes are safely committed to standby replicas before acknowledgement, preventing data loss during host outages.
- **Automated Failover**: If the primary PostgreSQL node becomes unresponsive, a healthy standby is promoted to primary within seconds.
- **HAProxy Routing**: Built-in HAProxy instances automatically detect the elected primary and route incoming database traffic accordingly.

Codedock includes a production-ready HA Compose template in `compose.ha.yml` and pre-built Patroni container configurations in `docker/patroni/`.

## Control Plane Leadership Election

When multiple Codedock daemon instances connect to the shared database:

- An internal engine leadership manager (`internal/engine/leadership`) elects an **active leader** instance using atomic database advisory locks and lease renewals.
- The active leader handles scheduled cron execution, background maintenance sweeps, and system update orchestrations.
- Standby control plane instances serve read-only dashboard traffic and API requests, providing horizontal throughput scaling.
- If the leader instance loses its database lease or shuts down, a standby node seamlessly assumes leadership without service interruption.

## Multi-Node Worker Clusters

In addition to control plane HA, you can cluster your compute fleet:

1. Connect multiple Linux servers across different availability zones or cloud providers.
2. Form runtime server clusters under **Fleet Management → Clusters**.
3. Deploy application containers across distinct worker nodes using round-robin or resource-aware placement policies.
4. If a physical worker server fails, Codedock can reschedule stateless workloads onto surviving nodes in the cluster.

## Health Monitoring & Recovery

- **Node Heartbeats**: The control plane continuously monitors worker connectivity over SSH and Docker probe endpoints.
- **Failover Drills**: Test failure recovery periodically by stopping primary nodes in non-production environments to verify that standby promotions complete smoothly.
- **Backup Redundancy**: Always maintain scheduled S3 or R2 database backups alongside active replication to safeguard against catastrophic multi-datacenter disruptions.
