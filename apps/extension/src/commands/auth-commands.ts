import * as vscode from 'vscode';
import { InfrastructureApi } from '../api';
import type { ConnectionManager } from '../auth/connection-manager';

export function registerAuthCommands(
  context: vscode.ExtensionContext,
  connectionManager: ConnectionManager
): void {
  // Connect command
  context.subscriptions.push(
    vscode.commands.registerCommand('codedock.connect', async () => {
      const existing = await connectionManager.getActiveConnection();
      const defaultCloudUrl =
        vscode.workspace.getConfiguration('codedock').get<string>('cloudUrl') ||
        process.env.CODEDOCK_CLOUD_URL ||
        'https://app.codedock.run';

      const method = await vscode.window.showQuickPick(
        [
          {
            label: '$(cloud) Authorize with Codedock Cloud',
            description: defaultCloudUrl,
            detail: 'Sign in seamlessly via your browser without entering a token',
            id: 'cloud',
          },
          {
            label: '$(key) Connect with Personal Access Token (PAT)',
            description: 'Enter your instance URL and token manually',
            detail: 'Supports Codedock Cloud and self-hosted instances',
            id: 'token',
          },
          {
            label: '$(server) Connect to Self-Hosted Instance',
            description: existing?.url || 'http://localhost:8080',
            detail: 'Connect directly to your self-hosted Codedock server',
            id: 'self-hosted',
          },
        ],
        { placeHolder: 'Choose how to connect to Codedock' }
      );

      if (!method) return;

      if (method.id === 'cloud') {
        await handleCloudDeviceLogin(defaultCloudUrl, connectionManager);
        return;
      }

      await handleManualTokenLogin(
        method.id === 'self-hosted' ? existing?.url || 'http://localhost:8080' : defaultCloudUrl,
        connectionManager
      );
    })
  );

  // Switch profile command
  context.subscriptions.push(
    vscode.commands.registerCommand('codedock.switchProfile', async () => {
      const profiles = await connectionManager.listProfiles();
      if (profiles.length === 0) {
        vscode.window.showInformationMessage('No saved profiles. Connect to an instance first.');
        return;
      }

      const selected = await vscode.window.showQuickPick(
        profiles.map((p) => ({
          label: p.name,
          description: p.url,
          detail: p.isActive ? '$(check) Currently active' : undefined,
          id: p.id,
        })),
        { placeHolder: 'Select a Codedock profile to activate' }
      );

      if (selected) {
        await connectionManager.switchProfile(selected.id);
        vscode.window.showInformationMessage(`Switched to profile: ${selected.label}`);
      }
    })
  );

  // Disconnect command
  context.subscriptions.push(
    vscode.commands.registerCommand('codedock.disconnect', async () => {
      const confirm = await vscode.window.showWarningMessage(
        'Disconnect from current Codedock instance?',
        { modal: true },
        'Disconnect'
      );
      if (confirm === 'Disconnect') {
        await connectionManager.clearActiveConnection();
        vscode.window.showInformationMessage('Disconnected from Codedock.');
      }
    })
  );

  // Quick Menu command
  context.subscriptions.push(
    vscode.commands.registerCommand('codedock.showQuickMenu', async () => {
      const conn = await connectionManager.getActiveConnection();
      if (!conn) {
        vscode.commands.executeCommand('codedock.connect');
        return;
      }

      const options = [
        { label: '$(sparkle) Ask Codedock AI Copilot', cmd: 'codedock.askAi' },
        { label: '$(plug) Configure MCP Server (.vscode/mcp.json)', cmd: 'codedock.configureMcp' },
        { label: '$(globe) Open Dashboard', cmd: 'codedock.openDashboard' },
        { label: '$(sync) Refresh All Data', cmd: 'codedock.refreshAll' },
        { label: '$(pulse) View System Metrics', cmd: 'codedock.viewSystemStats' },
        { label: '$(settings) Switch Profile', cmd: 'codedock.switchProfile' },
        { label: '$(plug) Add New Instance Connection', cmd: 'codedock.connect' },
        { label: '$(sign-out) Disconnect Current Instance', cmd: 'codedock.disconnect' },
      ];

      const choice = await vscode.window.showQuickPick(options, {
        placeHolder: `Codedock (${conn.instanceName || conn.url})`,
      });

      if (choice?.cmd) {
        vscode.commands.executeCommand(choice.cmd);
      }
    })
  );
}

async function handleCloudDeviceLogin(
  cloudUrl: string,
  connectionManager: ConnectionManager
): Promise<void> {
  try {
    const session = await InfrastructureApi.startDeviceAuth(cloudUrl, 'VS Code Extension');
    const authUrl = `${cloudUrl.replace(/\/+$/, '')}${session.verificationUriComplete}`;

    await vscode.env.openExternal(vscode.Uri.parse(authUrl));

    await vscode.window.withProgress(
      {
        location: vscode.ProgressLocation.Notification,
        title: `Codedock: Waiting for browser authorization (Code: ${session.userCode})...`,
        cancellable: true,
      },
      async (_progress, cancellationToken) => {
        let finished = false;
        cancellationToken.onCancellationRequested(() => {
          finished = true;
        });

        const start = Date.now();
        while (!finished && Date.now() - start < 15 * 60 * 1000) {
          await new Promise((resolve) => setTimeout(resolve, 2000));
          if (finished || cancellationToken.isCancellationRequested) break;

          try {
            const res = await InfrastructureApi.pollDeviceToken(cloudUrl, session.deviceCode);
            if (res.status === 'approved' && res.token) {
              const meta = await InfrastructureApi.getMeta(cloudUrl);
              const me = await InfrastructureApi.getMe(cloudUrl, res.token);
              await connectionManager.saveConnection(
                {
                  url: cloudUrl,
                  token: res.token,
                  userName: me.name,
                  userEmail: me.email,
                  version: meta.version,
                  instanceName: 'Codedock Cloud',
                },
                'cloud'
              );
              vscode.window.showInformationMessage(
                `Connected to Codedock Cloud (${meta.version}) as ${me.email}`
              );
              finished = true;
              return;
            }
            if (res.status === 'denied') {
              vscode.window.showWarningMessage('Codedock authorization request was declined.');
              finished = true;
              return;
            }
          } catch {
            // continue polling
          }
        }
      }
    );
  } catch (err: unknown) {
    const msg = err instanceof Error ? err.message : String(err);
    vscode.window.showErrorMessage(`Failed to start cloud login: ${msg}`);
  }
}

async function handleManualTokenLogin(
  defaultUrl: string,
  connectionManager: ConnectionManager
): Promise<void> {
  const urlInput = await vscode.window.showInputBox({
    title: 'Codedock: Instance URL',
    prompt:
      'Enter your Codedock daemon URL (e.g. https://app.codedock.run or https://codedock.example.com)',
    value: defaultUrl,
    ignoreFocusOut: true,
    validateInput: (val) => {
      if (!val) return 'Instance URL is required';
      if (!val.startsWith('http://') && !val.startsWith('https://')) {
        return 'URL must start with http:// or https://';
      }
      return null;
    },
  });

  if (!urlInput) return;

  const tokenInput = await vscode.window.showInputBox({
    title: 'Codedock: Personal Access Token',
    prompt: 'Enter your Personal Access Token (PAT)',
    password: true,
    ignoreFocusOut: true,
    validateInput: (val) => (!val ? 'Token is required' : null),
  });

  if (!tokenInput) return;

  const profileName = await vscode.window.showInputBox({
    title: 'Codedock: Profile Name',
    prompt: 'Friendly name for this instance profile (e.g. Cloud, Production)',
    value: urlInput.includes('codedock') ? 'Cloud' : 'Default',
    ignoreFocusOut: true,
  });

  await vscode.window.withProgress(
    {
      location: vscode.ProgressLocation.Notification,
      title: 'Connecting to Codedock...',
      cancellable: false,
    },
    async () => {
      try {
        const meta = await InfrastructureApi.getMeta(urlInput);
        const me = await InfrastructureApi.getMe(urlInput, tokenInput);

        const profileId = (profileName || 'default').toLowerCase().replace(/\s+/g, '-');
        await connectionManager.saveConnection(
          {
            url: urlInput,
            token: tokenInput,
            userName: me.name,
            userEmail: me.email,
            version: meta.version,
            instanceName: profileName || meta.name || 'Codedock',
          },
          profileId
        );

        vscode.window.showInformationMessage(
          `Connected to ${profileName || 'Codedock'} (${meta.version}) as ${me.email}`
        );
      } catch (err: unknown) {
        const msg = err instanceof Error ? err.message : String(err);
        vscode.window.showErrorMessage(`Failed to connect: ${msg}`);
      }
    }
  );
}
