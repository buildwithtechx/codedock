# Codedock for Visual Studio Code

The official Visual Studio Code extension for [Codedock](https://codedock.io) — deploy, manage, and monitor your containerized applications, databases, and microservices directly from your editor.

## Features

- **Direct Instance Connection**: Connect securely to any self-hosted or cloud Codedock daemon using Personal Access Tokens (PAT) stored safely in VS Code's `SecretStorage`.
- **Apps & Services Explorer**: Interactive sidebar view displaying real-time status of all your services (`running`, `stopped`, `failed`, `building`).
- **One-Click Deploy & Restart**: Trigger rebuilds, redeployments, or quick service restarts with single clicks or keyboard commands.
- **Live SSE Log Streaming**: Stream real-time container runtime logs directly into a dedicated VS Code Output Channel without leaving your code editor.
- **Status Bar Integration**: Always see your connected Codedock instance and jump straight to the web dashboard.

## Getting Started

1. Open the Command Palette (`Ctrl+Shift+P` / `Cmd+Shift+P`).
2. Run **`Codedock: Connect to Instance`**.
3. Enter your Codedock daemon URL (e.g. `http://localhost:8080` or `https://codedock.example.com`).
4. Enter your Personal Access Token (PAT) generated from your Codedock user settings.
5. Explore and manage your applications in the **Codedock** activity bar icon!

## Commands

- `Codedock: Connect to Instance` — Configure instance URL and PAT.
- `Codedock: Disconnect Instance` — Clear stored credentials.
- `Codedock: Refresh Apps` — Refresh the applications tree view.
- `Codedock: Redeploy App` — Trigger an application deployment.
- `Codedock: Restart App` — Restart an active application container.
- `Codedock: View Runtime Logs` — Stream live container runtime logs.
- `Codedock: Open in Dashboard` — Open the web dashboard in your default browser.
