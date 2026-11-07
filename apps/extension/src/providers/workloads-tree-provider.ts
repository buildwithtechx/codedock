import * as vscode from 'vscode';
import { DatabasesApi, WorkloadsApi } from '../api';
import type { ConnectionManager } from '../auth/connection-manager';
import type { CodedockApp, CodedockDatabase, CodedockProject } from '../interfaces';

export type WorkloadNode = ProjectNode | AppNode | DatabaseNode | ActionNode;

export class ProjectNode extends vscode.TreeItem {
  constructor(public readonly project: CodedockProject) {
    super(project.name, vscode.TreeItemCollapsibleState.Expanded);
    this.contextValue = 'codedock-project';
    this.iconPath = new vscode.ThemeIcon('project');
    this.tooltip = project.description || `Project: ${project.name}`;
  }
}

export class AppNode extends vscode.TreeItem {
  constructor(public readonly app: CodedockApp) {
    super(app.name, vscode.TreeItemCollapsibleState.None);
    this.contextValue = 'codedock-app';

    const status = (app.status || 'stopped').toLowerCase();
    this.description = app.domain ? `${status} · ${app.domain}` : status;
    this.tooltip = `Service: ${app.name}\nStatus: ${status}\nDomain: ${app.domain || 'None'}\nRepo: ${app.gitRepo || 'None'} (${app.branch || 'main'})`;

    if (status === 'running') {
      this.iconPath = new vscode.ThemeIcon('pass', new vscode.ThemeColor('testing.iconPassed'));
    } else if (status === 'building' || status === 'deploying') {
      this.iconPath = new vscode.ThemeIcon('sync~spin');
    } else if (status === 'failed' || status === 'error') {
      this.iconPath = new vscode.ThemeIcon('error', new vscode.ThemeColor('testing.iconErrored'));
    } else {
      this.iconPath = new vscode.ThemeIcon('circle-slash');
    }
  }
}

export class DatabaseNode extends vscode.TreeItem {
  constructor(public readonly database: CodedockDatabase) {
    super(database.name, vscode.TreeItemCollapsibleState.None);
    this.contextValue = 'codedock-database';

    const status = (database.status || 'stopped').toLowerCase();
    this.description = `${database.type} · ${status}`;
    this.tooltip = `Database: ${database.name}\nType: ${database.type}\nStatus: ${status}\nPort: ${database.port || 'default'}`;

    if (status === 'running') {
      this.iconPath = new vscode.ThemeIcon('database', new vscode.ThemeColor('testing.iconPassed'));
    } else {
      this.iconPath = new vscode.ThemeIcon('database');
    }
  }
}

export class ActionNode extends vscode.TreeItem {
  constructor(label: string, icon: string, commandId?: string) {
    super(label, vscode.TreeItemCollapsibleState.None);
    this.iconPath = new vscode.ThemeIcon(icon);
    if (commandId) {
      this.command = { command: commandId, title: label };
    }
  }
}

export class WorkloadsTreeProvider implements vscode.TreeDataProvider<WorkloadNode> {
  private readonly _onDidChangeTreeData = new vscode.EventEmitter<
    WorkloadNode | undefined | null
  >();
  readonly onDidChangeTreeData = this._onDidChangeTreeData.event;

  private cachedApps: CodedockApp[] = [];
  private cachedDatabases: CodedockDatabase[] = [];

  constructor(private readonly connectionManager: ConnectionManager) {
    this.connectionManager.onDidChangeConnection(() => this.refresh());
  }

  refresh(): void {
    this._onDidChangeTreeData.fire(undefined);
  }

  getTreeItem(element: WorkloadNode): vscode.TreeItem {
    return element;
  }

  async getChildren(element?: WorkloadNode): Promise<WorkloadNode[]> {
    const conn = await this.connectionManager.getActiveConnection();
    if (!conn) {
      return [new ActionNode('Connect to Codedock Instance...', 'plug', 'codedock.connect')];
    }

    try {
      if (!element) {
        // Root level: fetch projects, apps, databases
        const [projects, apps, databases] = await Promise.all([
          WorkloadsApi.listProjects(conn.url, conn.token),
          WorkloadsApi.listApps(conn.url, conn.token),
          DatabasesApi.listDatabases(conn.url, conn.token),
        ]);

        this.cachedApps = apps;
        this.cachedDatabases = databases;

        if (projects.length === 0 && apps.length === 0 && databases.length === 0) {
          return [new ActionNode('No projects or workloads found', 'info')];
        }

        if (projects.length > 0) {
          return projects.map((p) => new ProjectNode(p));
        }

        // Fallback: render apps and databases directly at root if no explicit projects
        const appNodes = apps.map((a) => new AppNode(a));
        const dbNodes = databases.map((d) => new DatabaseNode(d));
        return [...appNodes, ...dbNodes];
      }

      if (element instanceof ProjectNode) {
        const pId = element.project.id;
        const matchingApps = this.cachedApps
          .filter(
            (a) => a.projectId === pId || (!a.projectId && element.project.name === 'Default')
          )
          .map((a) => new AppNode(a));

        const matchingDbs = this.cachedDatabases
          .filter(
            (d) => d.projectId === pId || (!d.projectId && element.project.name === 'Default')
          )
          .map((d) => new DatabaseNode(d));

        const children: WorkloadNode[] = [...matchingApps, ...matchingDbs];
        if (children.length === 0) {
          return [new ActionNode('No services or databases in this project', 'info')];
        }
        return children;
      }

      return [];
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : String(err);
      return [new ActionNode(`Error loading workloads: ${msg}`, 'warning')];
    }
  }
}
