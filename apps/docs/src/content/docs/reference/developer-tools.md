---
title: Developer Tools & Mobile Apps
description: Official VS Code extension, desktop client, mobile app, and device authorization.
---

Codedock provides a comprehensive suite of developer tooling designed to integrate seamlessly into your day-to-day workflow, including an official **VS Code Extension**, a **Mobile App** for iOS & Android, a cross-platform **Desktop Client**, and secure **Device Authorization**.

## Official VS Code Extension

The **Codedock VS Code Extension** (`vscode-codedock`) brings the full power of your PaaS directly into your code editor.

### Key Capabilities
- **One-Click Deploy**: Right-click any file or open the command palette (`Ctrl+Shift+P` / `Cmd+Shift+P`) and run `Codedock: Deploy Current Project`.
- **Integrated Log Streaming**: Stream live container build and runtime logs directly into dedicated VS Code output channels with syntax highlighting.
- **Embedded Web Terminal**: Open an interactive shell terminal into any running container without leaving your editor.
- **Environment Management**: View, add, and modify service environment variables with immediate push synchronization.
- **Server Health Status**: Monitor CPU, RAM, and deployment statuses directly from your VS Code status bar.

### Installation
Search for `Codedock` in the VS Code Extensions Marketplace or install from `.vsix` in `apps/extension/`.

## Codedock Mobile App (iOS & Android)

Available for iPhone, iPad, and Android devices, the Codedock Mobile App (built with React Native and Expo) keeps you connected to your infrastructure on the go.

### Features
- **Real-Time Push Alerts**: Receive instant push notifications for deployment failures, server downtime, and disk saturation.
- **Emergency Operations**: Trigger emergency zero-downtime rollbacks, restart failing containers, or cancel runaway builds with a single tap.
- **Live Metrics**: Monitor live CPU and RAM usage charts across your entire worker fleet.
- **Biometric Security**: Protect mobile access with Face ID, Touch ID, or fingerprint authentication.

## Codedock Desktop App

For developers who prefer a dedicated desktop experience without keeping browser tabs open, Codedock provides a lightweight desktop client built with Tauri and Rust.
- **Minimal Resource Footprint**: Uses native operating system webviews for ultra-low memory consumption compared to Electron apps.
- **System Tray Quick-Actions**: Check cluster health at a glance from your macOS menu bar or Windows system tray.

## Device Authorization Flow (OAuth 2.0 Device Flow)

Authenticating CLI tools, IDE extensions, or mobile devices without entering your master password or copying sensitive API tokens is enabled via the **Device Authorization Flow**:

1. Run `codedock login` in your terminal or click **Sign In** in the VS Code extension.
2. The client requests a user code and presents an authorization URL (`/auth/device/verify`).
3. Open the link in your browser, verify that the 8-character code matches, and click **Authorize**.
4. The client automatically receives a cryptographically scoped Personal Access Token and securely stores it in your operating system keychain.
