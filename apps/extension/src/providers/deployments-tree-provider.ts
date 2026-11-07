import * as vscode from 'vscode';
import { DeploymentsApi } from '../api';
import type { ConnectionManager } from '../auth/connection-manager';
import type { CodedockDeployment } from '../interfaces';

export type DeploymentNode = DeploymentItemNode | DeploymentActionNode;

export class DeploymentItemNode extends vscode.TreeItem {
  constructor(public readonly deployment: CodedockDeployment) {
    const label = deployment.serviceName
      ? `${deployment.serviceName} (${deployment.commitHash?.slice(0, 7) || 'latest'})`
      : `Deployment ${deployment.id.slice(0, 8)}`;

    super(label, vscode.TreeItemCollapsibleState.None);
    this.contextValue = 'codedock-deployment';

    const status = (deployment.status || 'pending').toLowerCase();
    const duration = deployment.durationSeconds ? `${deployment.durationSeconds}s` : '';
    this.description = [status, deployment.branch, duration].filter(Boolean).join(' · ');
    this.tooltip = `Service: ${deployment.serviceName || 'App'}\nStatus: ${status}\nCommit: ${deployment.commitMessage || deployment.commitHash || 'Manual'}\nBranch: ${deployment.branch || 'main'}\nStarted: ${deployment.startedAt || 'N/A'}`;

    if (status === 'succeeded') {
      this.iconPath = new vscode.ThemeIcon('check', new vscode.ThemeColor('testing.iconPassed'));
    } else if (status === 'running' || status === 'pending') {
      this.iconPath = new vscode.ThemeIcon('sync~spin');
    } else if (status === 'failed') {
      this.iconPath = new vscode.ThemeIcon('error', new vscode.ThemeColor('testing.iconErrored'));
    } else {
      this.iconPath = new vscode.ThemeIcon('circle-slash');
    }

    this.command = {
      command: 'codedock.viewDeploymentLogs',
      title: 'View Deployment Logs',
      arguments: [this],
    };
  }
}

export class DeploymentActionNode extends vscode.TreeItem {
  constructor(label: string, icon: string, commandId?: string) {
    super(label, vscode.TreeItemCollapsibleState.None);
    this.iconPath = new vscode.ThemeIcon(icon);
    if (commandId) {
      this.command = { command: commandId, title: label };
    }
  }
}

export class DeploymentsTreeProvider implements vscode.TreeDataProvider<DeploymentNode> {
  private readonly _onDidChangeTreeData = new vscode.EventEmitter<
    DeploymentNode | undefined | null
  >();
  readonly onDidChangeTreeData = this._onDidChangeTreeData.event;

  constructor(private readonly connectionManager: ConnectionManager) {
    this.connectionManager.onDidChangeConnection(() => this.refresh());
  }

  refresh(): void {
    this._onDidChangeTreeData.fire(undefined);
  }

  getTreeItem(element: DeploymentNode): vscode.TreeItem {
    return element;
  }

  async getChildren(element?: DeploymentNode): Promise<DeploymentNode[]> {
    if (element) return [];

    const conn = await this.connectionManager.getActiveConnection();
    if (!conn) {
      return [new DeploymentActionNode('Connect to Codedock...', 'plug', 'codedock.connect')];
    }

    try {
      const deployments = await DeploymentsApi.listDeployments(conn.url, conn.token, 25);
      if (deployments.length === 0) {
        return [new DeploymentActionNode('No recent deployments', 'info')];
      }
      return deployments.map((d) => new DeploymentItemNode(d));
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : String(err);
      return [new DeploymentActionNode(`Failed to load deployments: ${msg}`, 'warning')];
    }
  }
}
