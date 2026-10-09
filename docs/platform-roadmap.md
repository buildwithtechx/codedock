# Platform implementation roadmap

Baseline: Codedock main 1c03759.
Managed cloud provider: Hetzner. Existing projects and self-hosted installations remain supported.

Phases 1 through 7 plus the migration lifecycle, traffic analytics, desktop lifecycle, app catalogue and REST schema scope are implemented; the completed checklist is preserved in git history. This document now tracks only what is left.

## Delivery workflow

Later phases can ship together in focused commits and one pull request, with one combined CI and Cubic review cycle. Run Biome and Go formatting and relevant regression checks before committing. The current batch uses incremental commits and pushes without creating a PR; CI and Cubic review remain pending until a PR is authorized. Merge only when authorized. No feature is complete merely because its settings form or API route exists. Runtime changes require target validation, permission checks, durable progress, cancellation and recovery.

## Remaining product scope

- [ ] Localization and RTL support across shared dashboard components.
- [ ] Mail server and webmail as an optional capability with separate operational requirements.
- [ ] Managed-server setup wizard in the dashboard for Hetzner provisioning.
- [ ] Codedock Cloud hosted offering on the HA control plane.
- [ ] Cloud connection and cloud deploy flow in settings, including plan picker, checkout feedback and credit alerts.
- [ ] Support center with ticket list, conversations and cloud-gated availability.
- [ ] Billing and usage experience: monthly host plans plus custom configurator, prepaid pay-as-you-go with top-ups, enterprise path, usage meters and Stripe portal/checkout recovery.
- [ ] Server cluster and network provisioning editors with preparation and operation tracking.
- [ ] Full managed data apps: Supabase, Convex and Neon equivalents with correct backing services.

## Live-target validation backlog

Source-level and mocked checks do not establish live infrastructure reliability. Each item below needs disposable targets and recovery exercises.

- [ ] Live Hetzner provisioning, cloud-init Docker bootstrap, resize power-cycling and SSH reachability on a disposable account.
- [ ] Live K3s installation, PostgreSQL/Redis operator recovery and native-server install/recovery exercises.
- [ ] Compose, SSH Docker and backup recovery exercises against disposable targets.
- [ ] Multi-host HA spread and failover drills for production certification.
- [ ] Maintainer install pass over the app catalogue to set verified badges.
- [ ] GeoIP database for IP-level traffic geography beyond private ranges.
- [ ] Desktop sidecar release bundling, updater signing keys and notarization CI.

## Related backlogs

Tunnel integration work is tracked separately: [TODO](TODO.md).
Implementation boundaries for deployment and recovery: [Application, Compose and topology workflows](platform-deployment.md), [Backup recovery and cluster setup](platform-recovery.md).

## Verification boundaries

Source-level and mocked provider checks do not establish live infrastructure reliability. Compose, backups, SSH, Hetzner, K3s and operators need disposable integration targets and recovery exercises. Provider secrets stay out of source control. Schema changes follow the repository initial-schema convention and numbered upgrades when required.
