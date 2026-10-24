# Changelog

All notable changes to the Codedock VS Code extension will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.0] - 2026-10-18

### Added
- **Multi-Profile Authentication**:
  - Secure credential storage using VS Code `SecretStorage`.
  - Profile switcher supporting local and remote Codedock control planes.
- **Three Dedicated Activity Bar Views**:
  - `Workloads & Databases`: Hierarchical status monitoring for deployed applications and managed databases.
  - `Servers & Clusters`: Overview of registered hosts, control planes, and SSH endpoints.
  - `Recent Deployments`: Deployment history with real-time status and duration.
- **Service & Database Management**:
  - Commands for redeploying, restarting, and stopping services.
  - Database restart and secure credential revelation.
- **Streaming Observability**:
  - Live Server-Sent Events (SSE) container log streaming into VS Code Output Channels.
- **Tooling & Packaging**:
  - esbuild bundling for minimal footprint (`~43 KB`).
  - Strict VS Code engine target `^1.96.0`.
  - Automated `.vsix` packaging and Vitest test suite.