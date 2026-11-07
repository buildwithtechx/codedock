# Changelog

All notable changes to the Codedock mobile app will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2026-10-18

### Added

- **5 Native Tabs Navigation**:
  - `Workloads`: Service fleet overview with health badges and quick-action detail modal.
  - `Databases`: Managed database inventory (PostgreSQL, Redis, MySQL, MongoDB) with container ports and restart/stop controls.
  - `Servers`: Connected infrastructure hosts and control planes with provider and region indicators.
  - `Deployments`: Chronological deployment history, commit messages, and build durations.
  - `Settings`: Multi-instance pairing with Personal Access Token (PAT) authentication and connection diagnostics.
- **TanStack Query State**:
  - Custom hooks: `useWorkloads`, `useDatabases`, `useDeployments`, `useServers`.
  - Automatic background polling and pull-to-refresh synchronization.
- **Expo Push Notifications**:
  - Integration with Codedock notification engine (`/api/notifications/devices`).
  - Automatic device token registration on connect for deployment and server status alerts.
- **Security & Storage**:
  - Encrypted storage for instance credentials via `expo-secure-store`.
  - Axios HTTP client with Bearer token authentication and timeout handling.
- **Theming & Design**:
  - Codedock dark zinc color palette (`#09090B` background, zinc surfaces and borders).
  - Tree-shakeable static icons via `@react-native-vector-icons/material-design-icons`.
- **EAS Build Automation**:
  - Preconfigured profiles for development, contributor preview (standalone APK), and production releases.
- **Quality & Testing**:
  - Jest unit test suite covering store hydration, utility functions, and API normalization.
  - TypeScript strict mode typechecking across the entire workspace.
