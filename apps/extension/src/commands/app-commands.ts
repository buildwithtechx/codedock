import * as vscode from 'vscode';
import { streamRuntimeLogs, WorkloadsApi } from '../api';
import type { ConnectionManager } from '../auth/connection-manager';
import type { CodedockApp } from '../interfaces';
import type { AppNode } from '../providers/workloads-tree-provider';
import type { OutputChannels } from '../views/output-channels';

let activeLogStream: AbortController | null = null;

async function resolveApp(
  connectionManager: ConnectionManager,
  node?: AppNode
): Promise<{ app: CodedockApp; url: string; token: string } | null> {
  const conn = await connectionManager.getActiveConnection();
  if (!conn) {
    vscode.window.showWarningMessage('Please connect to a Codedock instance first.');
    return null;
  }

  if (node?.app) {
    return { app: node.app, url: conn.url, token: conn.token };
  }

  const apps = await WorkloadsApi.listApps(conn.url, conn.token);
  if (apps.length === 0) {
    vscode.window.showInformationMessage('No apps found on this instance.');
    return null;
  }

  const selected = await vscode.window.showQuickPick(
    apps.map((a) => ({
      label: a.name,
      description: a.status || 'unknown',
      detail: a.domain ? `Domain: ${a.domain}` : `Project: ${a.projectName || 'Default'}`,
      app: a,
    })),
    { placeHolder: 'Select an application' }
  );

  if (!selected) return null;
  return { app: selected.app, url: conn.url, token: conn.token };
}

export function registerAppCommands(
  context: vscode.ExtensionContext,
  connectionManager: ConnectionManager,
  outputChannels: OutputChannels,
  onRefresh: () => void
): void {
  // Redeploy
  context.subscriptions.push(
    vscode.commands.registerCommand('codedock.deployApp', async (node?: AppNode) => {
      const target = await resolveApp(connectionManager, node);
      if (!target) return;

      await vscode.window.withProgress(
        {
          location: vscode.ProgressLocation.Notification,
          title: `Redeploying ${target.app.name}...`,
          cancellable: false,
        },
        async () => {
          try {
            await WorkloadsApi.redeployApp(target.url, target.token, target.app.id);
            vscode.window.showInformationMessage(`Redeploy triggered for ${target.app.name}`);
            onRefresh();
          } catch (err: unknown) {
            const msg = err instanceof Error ? err.message : String(err);
            vscode.window.showErrorMessage(`Deploy failed: ${msg}`);
          }
        }
      );
    })
  );

  // Restart
  context.subscriptions.push(
    vscode.commands.registerCommand('codedock.restartApp', async (node?: AppNode) => {
      const target = await resolveApp(connectionManager, node);
      if (!target) return;

      await vscode.window.withProgress(
        {
          location: vscode.ProgressLocation.Notification,
          title: `Restarting ${target.app.name}...`,
          cancellable: false,
        },
        async () => {
          try {
            await WorkloadsApi.restartApp(target.url, target.token, target.app.id);
            vscode.window.showInformationMessage(`Restart triggered for ${target.app.name}`);
            onRefresh();
          } catch (err: unknown) {
            const msg = err instanceof Error ? err.message : String(err);
            vscode.window.showErrorMessage(`Restart failed: ${msg}`);
          }
        }
      );
    })
  );

  // Stop
  context.subscriptions.push(
    vscode.commands.registerCommand('codedock.stopApp', async (node?: AppNode) => {
      const target = await resolveApp(connectionManager, node);
      if (!target) return;

      const confirm = await vscode.window.showWarningMessage(
        `Are you sure you want to stop service "${target.app.name}"?`,
        { modal: true },
        'Stop Service'
      );
      if (confirm !== 'Stop Service') return;

      try {
        await WorkloadsApi.stopApp(target.url, target.token, target.app.id);
        vscode.window.showInformationMessage(`Service ${target.app.name} stopped`);
        onRefresh();
      } catch (err: unknown) {
        const msg = err instanceof Error ? err.message : String(err);
        vscode.window.showErrorMessage(`Stop failed: ${msg}`);
      }
    })
  );

  // Stream Logs
  context.subscriptions.push(
    vscode.commands.registerCommand('codedock.viewLogs', async (node?: AppNode) => {
      const target = await resolveApp(connectionManager, node);
      if (!target) return;

      if (activeLogStream) {
        activeLogStream.abort();
        activeLogStream = null;
      }

      const channel = outputChannels.runtimeLogs;
      channel.clear();
      channel.show(true);
      channel.appendLine(`=== Streaming Runtime Logs for ${target.app.name} ===\n`);

      activeLogStream = streamRuntimeLogs(
        target.url,
        target.token,
        target.app.id,
        (line) => channel.appendLine(line),
        (err) => channel.appendLine(`[Stream Notice]: ${err.message}`)
      );
    })
  );

  // View Env Vars
  context.subscriptions.push(
    vscode.commands.registerCommand('codedock.viewEnvVars', async (node?: AppNode) => {
      const target = await resolveApp(connectionManager, node);
      if (!target) return;

      try {
        const vars = await WorkloadsApi.getAppVariables(target.url, target.token, target.app.id);
        if (vars.length === 0) {
          vscode.window.showInformationMessage(
            `No environment variables defined for ${target.app.name}`
          );
          return;
        }

        const picked = await vscode.window.showQuickPick(
          vars.map((v) => ({
            label: v.key,
            description: v.isSecret ? '••••••••' : v.value,
            detail: v.isSecret ? 'Encrypted secret variable' : undefined,
            variable: v,
          })),
          { placeHolder: `Environment variables for ${target.app.name}` }
        );

        if (picked) {
          const action = await vscode.window.showQuickPick(
            ['Copy Variable Name', 'Copy Variable Value'],
            { placeHolder: picked.label }
          );
          if (action === 'Copy Variable Name') {
            await vscode.env.clipboard.writeText(picked.variable.key);
          } else if (action === 'Copy Variable Value') {
            await vscode.env.clipboard.writeText(picked.variable.value);
            vscode.window.showInformationMessage(
              `Copied ${picked.variable.key} value to clipboard`
            );
          }
        }
      } catch (err: unknown) {
        const msg = err instanceof Error ? err.message : String(err);
        vscode.window.showErrorMessage(`Failed to fetch variables: ${msg}`);
      }
    })
  );

  // Copy Public URL
  context.subscriptions.push(
    vscode.commands.registerCommand('codedock.copyPublicUrl', async (node?: AppNode) => {
      const target = await resolveApp(connectionManager, node);
      if (!target) return;

      if (!target.app.domain) {
        vscode.window.showInformationMessage(`No domain configured for ${target.app.name}`);
        return;
      }

      const domain = target.app.domain.startsWith('http')
        ? target.app.domain
        : `https://${target.app.domain}`;

      await vscode.env.clipboard.writeText(domain);
      vscode.window.showInformationMessage(`Copied ${domain} to clipboard`);
    })
  );

  // Open Dashboard
  context.subscriptions.push(
    vscode.commands.registerCommand('codedock.openDashboard', async (node?: AppNode) => {
      const conn = await connectionManager.getActiveConnection();
      if (!conn) {
        vscode.window.showWarningMessage('Please connect to a Codedock instance first.');
        return;
      }

      let targetUrl = conn.url;
      if (node?.app?.id) {
        targetUrl = `${conn.url}/apps/${encodeURIComponent(node.app.id)}`;
      }

      vscode.env.openExternal(vscode.Uri.parse(targetUrl));
    })
  );
}

export function stopActiveLogStream(): void {
  if (activeLogStream) {
    activeLogStream.abort();
    activeLogStream = null;
  }
}
