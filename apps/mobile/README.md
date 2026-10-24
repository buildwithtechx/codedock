# Codedock Mobile App

The official iOS and Android client for [Codedock](https://github.com/buildwithtechx/codedock) — monitor your infrastructure, trigger service redeploys, inspect recent deployments, manage databases, and receive real-time push alerts on the go.

---

## Highlights & Features

- **5 Native Tabs**:
  - 📦 **Workloads**: Monitor running services, health status, domains, and trigger redeploy/restart/stop commands.
  - 🗄️ **Databases**: Inspect managed databases (PostgreSQL, Redis, MySQL, MongoDB), connection endpoints, container ports, and restart/stop controls.
  - 🖥️ **Servers**: Fleet overview of connected servers, control planes, IP addresses, providers, and regions.
  - 🚀 **Deployments**: Chronological history of recent deployments, commit messages, branches, and durations.
  - ⚙️ **Settings**: Instant instance connection with URL and Personal Access Token (PAT), stored securely in device hardware via `expo-secure-store`.
- **Native Push Notifications**: Automatic device token registration with Codedock's notification pipeline for deployment success/failure, backup alerts, and server unreachable events.
- **Dark Zinc Theme**: Matches Codedock's web dashboard and documentation (`#09090B` background, zinc cards and borders).
- **Offline & Cache-Optimized**: Powered by TanStack Query with background polling and instant pull-to-refresh.

---

## Tech Stack

- **Framework**: [Expo SDK 57](https://expo.dev) + [React Native 0.86](https://reactnative.dev)
- **Routing**: [Expo Router](https://docs.expo.dev/router/introduction/) (file-based navigation under `src/app`)
- **State & Caching**: [@tanstack/react-query](https://tanstack.com/query) + [Zustand](https://zustand.docs.pmnd.rs)
- **Network**: [Axios](https://axios-http.com) client with authorization interceptors
- **Storage**: [expo-secure-store](https://docs.expo.dev/versions/latest/sdk/securestore/)
- **Icons**: [@react-native-vector-icons/material-design-icons](https://www.npmjs.com/package/@react-native-vector-icons/material-design-icons) (static, tree-shakeable)
- **Testing**: [Jest](https://jestjs.io) with `jest-expo`

---

## Directory Structure

```text
apps/mobile/
├── assets/images/          # Official brand icons and splash screens
├── eas.json                # EAS Build configuration for dev, preview & prod
├── test/                   # Jest unit test suites
│   ├── store.test.ts
│   └── utils.test.ts
├── src/
│   ├── app/                # Expo Router native routes & tabs
│   │   ├── (tabs)/         # 5 Tabs: index, databases, servers, deploys, settings
│   │   ├── _layout.tsx     # Root navigation & QueryClientProvider
│   │   ├── modal.tsx       # Workload action modal
│   │   └── +not-found.tsx  # Missing route fallback
│   ├── components/         # Reusable native UI components (Icon)
│   ├── hooks/              # TanStack Query hooks (workloads, databases, deploys, servers)
│   ├── interfaces/         # TypeScript type definitions
│   ├── services/           # API client, SecureStore storage, Push notifications
│   ├── store/              # Zustand app state
│   ├── theme/              # Codedock dark zinc color palette & spacing
│   └── utils/              # Formatting and helper utilities
├── babel.config.js
├── jest.config.js
├── metro.config.js
└── package.json
```

---

## Getting Started

### 1. Install Dependencies

From the repository root:

```bash
npm install
```

### 2. Start the Development Server

```bash
npm run dev:mobile
# or inside apps/mobile:
npm start
```

### 3. Run on Emulator / Physical Device

```bash
# Android emulator or connected device
npm run android

# iOS simulator (macOS only)
npm run ios
```

---

## Contributor Builds (EAS)

The app includes predefined build profiles in [`eas.json`](./eas.json):

### Contributor Preview Builds (Standalone Installers)

To build a standalone preview APK for direct installation on physical Android devices without requiring store access:

```bash
npm run build:preview:android
```

For iOS simulator preview builds:

```bash
npm run build:preview:ios
```

### Development Client Builds

To build custom development clients with all native plugins compiled:

```bash
npm run build:dev:android
npm run build:dev:ios
```

### Production Release Builds

```bash
npm run build:prod
```

---

## Testing & Quality

```bash
# Typecheck
npm run typecheck

# Unit tests
npm test

# Biome code formatting
npm run format
```
