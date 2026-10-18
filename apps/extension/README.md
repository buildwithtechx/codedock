# Codedock for Visual Studio Code

The official Visual Studio Code extension for [Codedock](https://codedock.run) — deploy, manage, and monitor your containerized applications, databases, and microservices directly from your editor.

---

## Features

- **Multi-Profile Authentication**: Securely manage multiple Codedock instances using Personal Access Tokens (PAT) stored encrypted in VS Code's `SecretStorage`.
- **Workloads & Databases Explorer**: Interactive tree view displaying real-time status of all your services and databases (`running`, `stopped`, `failed`, `building`).
- **Servers & Clusters Explorer**: Monitor connected server nodes, control planes, IP addresses, and SSH hosts.
- **Recent Deployments Explorer**: Track recent build and deployment runs with commit messages, branch names, and durations.
- **One-Click Deploy & Restart**: Trigger redeployments, restarts, or stops with single clicks or keyboard commands.
- **Database Credential Revelation**: Securely reveal and inspect managed database passwords directly inside VS Code.
- **Live SSE Log Streaming**: Stream real-time container runtime logs directly into a dedicated VS Code Output Channel without leaving your editor.
- **Status Bar Integration**: Instantly view your active Codedock instance and switch between connection profiles.

---

## Getting Started

1. Open the Command Palette (`Ctrl+Shift+P` / `Cmd+Shift+P`).
2. Run **`Codedock: Connect to Instance`**.
3. Enter your Codedock daemon URL (e.g. `http://localhost:8080` or `https://codedock.example.com`).
4. Enter your Personal Access Token (PAT) generated from your Codedock profile settings.
5. Access your workloads, databases, and servers in the **Codedock** activity bar icon!

---

## Available Commands

- `Codedock: Connect to Instance` — Configure instance URL and PAT.
- `Codedock: Switch Instance Profile` — Switch between saved Codedock connections.
- `Codedock: Disconnect Instance` — Clear stored credentials.
- `Codedock: Refresh Views` — Refresh all trees (workloads, servers, deployments).
- `Codedock: Redeploy Service` — Trigger a workload redeployment.
- `Codedock: Restart Service` — Restart an active service container.
- `Codedock: Stop Service` — Stop a running service container.
- `Codedock: Restart Database` — Restart an active managed database.
- `Codedock: Reveal Database Password` — Reveal database credentials securely.
- `Codedock: View Runtime Logs` — Stream live container runtime logs via Server-Sent Events.
- `Codedock: Open in Dashboard` — Open the web dashboard in your default browser.
