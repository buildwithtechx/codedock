import * as vscode from 'vscode';
import { DeploymentsApi, InfrastructureApi } from '../api';
import type { ConnectionManager } from '../auth/connection-manager';
import type { DeploymentItemNode } from '../providers/deployments-tree-provider';
import type { ServerItemNode } from '../providers/servers-tree-provider';
import type { OutputChannels } from '../views/output-channels';

export function registerServerCommands(
  context: vscode.ExtensionContext,
  connectionManager: ConnectionManager,
  outputChannels: OutputChannels,
  onRefreshAll: () => void
): void {
  // Test Server SSH
  context.subscriptions.push(
    vscode.commands.registerCommand('codedock.testServerSSH', async (node?: ServerItemNode) => {
      const conn = await connectionManager.getActiveConnection();
      if (!conn) {
        vscode.window.showWarningMessage('Please connect to a Codedock instance first.');
        return;
      }

      if (!node?.server) {
        vscode.window.showInformationMessage('Please select a server to test SSH connection.');
        return;
      }

      await vscode.window.withProgress(
        {
          location: vscode.ProgressLocation.Notification,
          title: `Testing SSH connection to ${node.server.name}...`,
          cancellable: false,
        },
        async () => {
          try {
            const res = await InfrastructureApi.testServerSSH(conn.url, conn.token, node.server.id);
            if (res.success) {
              vscode.window.showInformationMessage(
                `SSH Connection to ${node.server.name} successful!`
              );
            } else {
              vscode.window.showErrorMessage(`SSH Test Failed: ${res.message || 'Unreachable'}`);
            }
          } catch (err: unknown) {
            const msg = err instanceof Error ? err.message : String(err);
            vscode.window.showErrorMessage(`SSH Test Error: ${msg}`);
          }
        }
      );
    })
  );

  // View System Stats
  context.subscriptions.push(
    vscode.commands.registerCommand('codedock.viewSystemStats', async () => {
      const conn = await connectionManager.getActiveConnection();
      if (!conn) {
        vscode.window.showWarningMessage('Please connect to a Codedock instance first.');
        return;
      }

      try {
        const stats = await InfrastructureApi.getSystemStats(conn.url, conn.token);
        const cpuPercent = stats.cpu?.percent?.toFixed(1) ?? 'N/A';
        const memPercent = stats.memory?.percent?.toFixed(1) ?? 'N/A';
        const diskPercent = stats.disk?.percent?.toFixed(1) ?? 'N/A';
        const uptimeHours = stats.uptimeSeconds ? (stats.uptimeSeconds / 3600).toFixed(1) : 'N/A';

        vscode.window.showInformationMessage(
          `Codedock Health — CPU: ${cpuPercent}% | RAM: ${memPercent}% | Disk: ${diskPercent}% | Uptime: ${uptimeHours}h`
        );
      } catch (err: unknown) {
        const msg = err instanceof Error ? err.message : String(err);
        vscode.window.showErrorMessage(`Failed to fetch stats: ${msg}`);
      }
    })
  );

  // View Deployment Logs
  context.subscriptions.push(
    vscode.commands.registerCommand(
      'codedock.viewDeploymentLogs',
      async (node?: DeploymentItemNode) => {
        const conn = await connectionManager.getActiveConnection();
        if (!conn || !node?.deployment) return;

        const channel = outputChannels.deploymentLogs;
        channel.clear();
        channel.show(true);
        channel.appendLine(
          `=== Logs for Deployment ${node.deployment.id} (${node.deployment.serviceName || 'Service'}) ===\n`
        );

        try {
          const logs = await DeploymentsApi.getDeploymentLogs(
            conn.url,
            conn.token,
            node.deployment.id
          );
          channel.appendLine(logs || 'No build logs recorded for this deployment.');
        } catch (err: unknown) {
          channel.appendLine(
            `Failed to fetch logs: ${err instanceof Error ? err.message : String(err)}`
          );
        }
      }
    )
  );

  // Refresh All
  context.subscriptions.push(
    vscode.commands.registerCommand('codedock.refreshAll', () => {
      onRefreshAll();
      vscode.window.showInformationMessage(
        'Refreshed Codedock workloads, servers, and deployments.'
      );
    })
  );
}
