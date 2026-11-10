// Package platform contains technical infrastructure drivers and system-level integrations.
//
// # Architectural Purpose:
//
// Unlike business domain modules, platform packages do not manage REST endpoints
// or API request/response DTOs. Instead, they interact with the host operating system,
// container runtimes, remote networks, and third-party infrastructure.
//
// # Subpackages:
//   - build: Dockerfile synthesizing, Railpack builders, and build cache management.
//   - docker: Docker engine client, GPU access, and Compose stack runtime.
//   - ssh: Remote server SSH connection pools, key management, and script execution.
//   - networking: Traefik dynamic configuration, Let's Encrypt TLS, and port allocation.
//   - storage: S3/SFTP destinations, tar archives, and backup stream producers.
//   - database: Database connection pools, pg_dump/restore drivers, and supervisor routines.
//   - vault: AES-GCM data encryption for secrets and credentials.
//   - scheduler: Background cron job runner and distributed leader election.
//   - observability: Prometheus/TSDB metric workers, Loki log streams, and GeoIP lookups.
//   - notifications: Notification channels (Slack, Discord, Telegram, Pushover, Webhooks, Expo).
//   - bare: Direct bare-metal server script runners and runtime lifecycle.
package platform
