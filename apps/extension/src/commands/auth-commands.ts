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
      const defaultUrl = existing?.url || 'http://localhost:8080';

      const urlInput = await vscode.window.showInputBox({
        title: 'Codedock: Instance URL',
        prompt:
          'Enter your Codedock daemon URL (e.g. http://localhost:8080 or https://codedock.example.com)',
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
        prompt: 'Friendly name for this instance profile (e.g. Local Dev, Production)',
        value: 'Default',
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
