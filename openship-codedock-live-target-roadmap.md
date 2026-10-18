# OpenShip vs Codedock: Live Target & Platform Roadmap

> **Status**: Parity tracks complete 🎯 — owner-side locked drills remain 🔒;
> new own-plan track: notifications gap-fill + clients (items 15–20, locked 21).
> **Reference**: openship `oblien/openship` at `2e42298df` (v0.9.3 reviewed Oct 17)
> **Last Updated**: October 17, 2026
> **Legend**: ✅ done, co-located tests in tree · ⬜ open · 🔒 needs real infra/certs/money (owner-side, not agent-doable)

---

## Remaining Work

### 🔒 Live drills & release CI (real infra, certs, or spend required)

- [ ] **8. Live Hetzner Provisioning & Bootstrap Drill** — disposable project:
  cloud-init <180s, SSH + Docker socket, power-cycle + `cx22`→`cpx21` resize
  keeps data, clean delete with no dangling resources.
- [ ] **9. Live K3s Installation & Operator Recovery** — bootstrap,
  CloudNative-PG & Redis failover under network partition.
- [ ] **10. Compose & Backup Recovery Exercises** — named volumes, Minisign/S3
  target, synthetic wipe & restore.
- [ ] **11. Multi-Host HA Spread & Failover Drills** — WireGuard latency
  matrix, MTU checks, split-brain tolerance.
- [ ] **14. Desktop Sidecar Bundling & Notarization CI** — Tauri `externalBin`,
  Minisign updater, Apple Notarytool (needs Apple certs).
- [ ] **2. Codedock Cloud Offering** — HA control plane, `pg_advisory_lock`
  leadership, reverse agent tunnels. On arrival also adopt: account/org
  identity switching, relay + import isolation, provisioning authority,
  workspace lifecycle drills, Cloud-Spaces wording (v0.9.2, presentation-only
  upstream — nothing to port self-hosted).
- [ ] **21. Client Release Accounts & Signing** 🔒 — VS Code Marketplace
  publisher + `VSCE_PAT`; Apple Developer Program, Google Play Console, EAS
  project/credentials; store listings + privacy policy. Gates items 17/18.

### 📱 Clients & notifications (own plan — not upstream parity)

> Not demand-pulled: upstream's only mobile issue (#1140) is a maintainer-scoped
> Swift/iOS brief (0 reactions); zero editor/extension requests across 105 open
> issues. Sequence by cost. Notification core is already ported (Discord, Slack,
> Telegram, Pushover, Resend, generic webhook, delivery history, subscriptions,
> per-service webhooks) — items below are gaps, not builds.

- [x] **15. Notification Event Coverage + Expo Push Channel** ✅ — `b3480514`:
  standardized event types (`deploy.succeeded`, `deploy.failed`, `server.unreachable`,
  `backup.failed`, `domain.expiring`), dynamic dashboard deep-links, wired missing
  daemon sources (server unreachable on SSH monitor fail, backup failed on engine error),
  device token registration endpoints (`/api/notifications/devices`), and Expo push
  provider (`internal/notifications/expo.go`).
- [x] **16. PWA Completion (service worker only)** ✅ — manifest updated with
  maskable icons and PWA install metadata, `sw.js` app-shell caching with
  network-first navigation, static asset cache, offline fallback page (`offline.html`),
  and Web Push listener scaffold. Registered in `apps/dashboard/src/main.tsx`.
- [x] **17. VS Code Extension (Production Parity)** ✅ — public `GET /api/meta`,
  secure multi-profile authentication with PAT in `SecretStorage`, hierarchical
  workloads & databases explorer, infrastructure servers explorer, recent deployments
  explorer, deploy/redeploy/restart/stop commands, live SSE log streaming into
  dedicated output channels, database credential revelation, modern VS Code engine
  `^1.96.0`, clean modular architecture, bundled with esbuild, packaged via `vsce`.
- [x] **18. Expo Mobile App (iOS + Android)** ✅ — Full production parity matching
  TakeNTrade architecture under unified `src/`: TanStack Query hooks (`useWorkloads`,
  `useDatabases`, `useDeployments`, `useServers`), Axios API client, Zustand store
  (`useAppStore`), SecureStore PAT persistence, Expo push notification registration
  with Item 15 backend, dark zinc theme matching web and dashboard (`#09090B`), 5 native
  tabs (Workloads, Databases, Servers, Deployments, Settings), app.json with brand icons,
  splash screen, and plugins, EAS build and submit scripts, and Jest tests.
- [x] **19. Web Push (VAPID) Channel** (Superseded by Item 18) — Native push is
  handled directly on iOS and Android by the Expo Mobile App using the Item 15
  channel (`/api/notifications/devices`); redundant browser VAPID overhead skipped.
- [x] **20. Client CI & Workspace Coverage** ✅ — Root `package.json` wires
  extension + mobile typechecks and tests into root commands (`typecheck`, `test`,
  `typecheck:mobile`, `typecheck:extension`, `test:mobile`, `test:extension`). Biome
  includes mobile and extension under `src/**`. CI pipeline (`.github/workflows/ci.yml`)
  runs both clients through typechecks, unit tests, compilation, and `vsce package`
  verification.

---

## Historical Notes (Kept for Reference)

- Upstream v0.9.2 (`d2e1a04a9`) deleted its managed tab UI in favor of a
  Resources view; codedock keeps outbound/SSH/network tabs — correct for
  self-hosted hosts. Not a regression.
- `odoo.yaml` pins v18; upstream moved to `odoo:19.0-20260926` (see 21c).
- Template-schema ideas banked for later: `gpus`/`shmSize`, `supportsCloud`,
  `requires[].bundleId`, `generateGroup` shared secrets, `minEngine` gates.
- Upstream Emails-everywhere = connect-external client on cloud, mail-server
  admin stays self-hosted-only (see 26). Cloud Spaces = presentation-only
  (see 2).
- Preview sweep + leases no-op upstream without `OPENSHIP_NATIVE` +
  `OPENSHIP_NATIVE_JOBS` + an org owner; codedock workers are always native —
  single-controller concurrency is simpler than the lease protocol (see 36).
- Roles wave 2 is self-hosted Docker-server scope upstream (configuration
  refused on cloud/cluster/multi-service/static); the handover journal itself
  stays portable (see 37).
- Build caches: upstream `PACKAGE_BUILD_CACHES` dirs are byte-identical to
  v0.9.2 (shared source, 15 managers, yarn omitted); codedock
  `GetPackageCacheDirs` already mirrors — new work is receipts + migration +
  log summary (see 38).
- PR preview rows survive whole-instance transfer; portable export/import
  strips preview settings + cache receipts — adapt to codedock transfer
  semantics when porting (see 36).
- v0.9.3 Cloud-only hardening (build-cache cleanup scoped to owned project
  workers with no daemon-wide prune, monitoring memory safety on
  fuse-overlayfs, guest RAM counters, partial-metric retention, cloud builds
  require the managed executor) — adopt together with item 2.
- v0.9.3 edge docs split self-hosted vs Cloud capabilities (trusted-proxy env
  vars, containerized-edge header caveat, reapply-routing after changes) —
  codedock trusted proxies already done (ledger 23); informational only.
