import * as vscode from 'vscode';
import { InfrastructureApi } from '../api';
import type { ConnectionManager } from '../auth/connection-manager';
import type { CodedockServer } from '../interfaces';

export type ServerNode = ServerItemNode | ServerActionNode;

export class ServerItemNode extends vscode.TreeItem {
  constructor(public readonly server: CodedockServer) {
    super(server.name, vscode.TreeItemCollapsibleState.None);
    this.contextValue = 'codedock-server';

    const status = (server.status || 'offline').toLowerCase();
    this.description = `${server.ipAddress} · ${status}`;
    this.tooltip = `Server: ${server.name}\nIP: ${server.ipAddress}\nStatus: ${status}\nDocker: ${server.dockerVersion || 'N/A'}\nOS: ${server.osInfo || 'Linux'}`;

    if (status === 'online') {
      this.iconPath = new vscode.ThemeIcon('server', new vscode.ThemeColor('testing.iconPassed'));
    } else {
      this.iconPath = new vscode.ThemeIcon('server-environment');
    }
  }
}

export class ServerActionNode extends vscode.TreeItem {
  constructor(label: string, icon: string, commandId?: string) {
    super(label, vscode.TreeItemCollapsibleState.None);
    this.iconPath = new vscode.ThemeIcon(icon);
    if (commandId) {
      this.command = { command: commandId, title: label };
    }
  }
}

export class ServersTreeProvider implements vscode.TreeDataProvider<ServerNode> {
  private readonly _onDidChangeTreeData = new vscode.EventEmitter<ServerNode | undefined | null>();
  readonly onDidChangeTreeData = this._onDidChangeTreeData.event;

  constructor(private readonly connectionManager: ConnectionManager) {
    this.connectionManager.onDidChangeConnection(() => this.refresh());
  }

  refresh(): void {
    this._onDidChangeTreeData.fire(undefined);
  }

  getTreeItem(element: ServerNode): vscode.TreeItem {
    return element;
  }

  async getChildren(element?: ServerNode): Promise<ServerNode[]> {
    if (element) return [];

    const conn = await this.connectionManager.getActiveConnection();
    if (!conn) {
      return [new ServerActionNode('Connect to Codedock...', 'plug', 'codedock.connect')];
    }

    try {
      const servers = await InfrastructureApi.listServers(conn.url, conn.token);
      if (servers.length === 0) {
        return [new ServerActionNode('No servers connected', 'info')];
      }
      return servers.map((s) => new ServerItemNode(s));
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : String(err);
      return [new ServerActionNode(`Failed to load servers: ${msg}`, 'warning')];
    }
  }
}
