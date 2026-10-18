import * as vscode from 'vscode';
import { DatabasesApi } from '../api';
import type { ConnectionManager } from '../auth/connection-manager';
import type { CodedockDatabase } from '../interfaces';
import type { DatabaseNode } from '../providers/workloads-tree-provider';

async function resolveDatabase(
  connectionManager: ConnectionManager,
  node?: DatabaseNode
): Promise<{ db: CodedockDatabase; url: string; token: string } | null> {
  const conn = await connectionManager.getActiveConnection();
  if (!conn) {
    vscode.window.showWarningMessage('Please connect to a Codedock instance first.');
    return null;
  }

  if (node?.database) {
    return { db: node.database, url: conn.url, token: conn.token };
  }

  const dbs = await DatabasesApi.listDatabases(conn.url, conn.token);
  if (dbs.length === 0) {
    vscode.window.showInformationMessage('No databases found on this instance.');
    return null;
  }

  const selected = await vscode.window.showQuickPick(
    dbs.map((d) => ({
      label: d.name,
      description: `${d.type} · ${d.status}`,
      detail: `Internal Port: ${d.port || 'default'}`,
      db: d,
    })),
    { placeHolder: 'Select a database' }
  );

  if (!selected) return null;
  return { db: selected.db, url: conn.url, token: conn.token };
}

export function registerDatabaseCommands(
  context: vscode.ExtensionContext,
  connectionManager: ConnectionManager,
  onRefresh: () => void
): void {
  // Restart Database
  context.subscriptions.push(
    vscode.commands.registerCommand('codedock.restartDatabase', async (node?: DatabaseNode) => {
      const target = await resolveDatabase(connectionManager, node);
      if (!target) return;

      await vscode.window.withProgress(
        {
          location: vscode.ProgressLocation.Notification,
          title: `Restarting database ${target.db.name}...`,
          cancellable: false,
        },
        async () => {
          try {
            await DatabasesApi.restartDatabase(target.url, target.token, target.db.id);
            vscode.window.showInformationMessage(`Restarted database ${target.db.name}`);
            onRefresh();
          } catch (err: unknown) {
            const msg = err instanceof Error ? err.message : String(err);
            vscode.window.showErrorMessage(`Database restart failed: ${msg}`);
          }
        }
      );
    })
  );

  // Stop Database
  context.subscriptions.push(
    vscode.commands.registerCommand('codedock.stopDatabase', async (node?: DatabaseNode) => {
      const target = await resolveDatabase(connectionManager, node);
      if (!target) return;

      const confirm = await vscode.window.showWarningMessage(
        `Stop database "${target.db.name}"? Connected services may fail.`,
        { modal: true },
        'Stop Database'
      );
      if (confirm !== 'Stop Database') return;

      try {
        await DatabasesApi.stopDatabase(target.url, target.token, target.db.id);
        vscode.window.showInformationMessage(`Stopped database ${target.db.name}`);
        onRefresh();
      } catch (err: unknown) {
        const msg = err instanceof Error ? err.message : String(err);
        vscode.window.showErrorMessage(`Database stop failed: ${msg}`);
      }
    })
  );

  // Start Database
  context.subscriptions.push(
    vscode.commands.registerCommand('codedock.startDatabase', async (node?: DatabaseNode) => {
      const target = await resolveDatabase(connectionManager, node);
      if (!target) return;

      try {
        await DatabasesApi.startDatabase(target.url, target.token, target.db.id);
        vscode.window.showInformationMessage(`Started database ${target.db.name}`);
        onRefresh();
      } catch (err: unknown) {
        const msg = err instanceof Error ? err.message : String(err);
        vscode.window.showErrorMessage(`Database start failed: ${msg}`);
      }
    })
  );

  // Reveal Credentials
  context.subscriptions.push(
    vscode.commands.registerCommand('codedock.revealCredentials', async (node?: DatabaseNode) => {
      const target = await resolveDatabase(connectionManager, node);
      if (!target) return;

      try {
        const creds = await DatabasesApi.revealCredentials(target.url, target.token, target.db.id);
        const connectionStr =
          creds.connectionString ||
          `postgres://${creds.username || 'user'}:${creds.password || ''}@localhost:${target.db.port || 5432}/db`;

        const action = await vscode.window.showInformationMessage(
          `Connection credentials retrieved for ${target.db.name}`,
          'Copy Connection String',
          'Show Password'
        );

        if (action === 'Copy Connection String') {
          await vscode.env.clipboard.writeText(connectionStr);
          vscode.window.showInformationMessage('Connection string copied to clipboard!');
        } else if (action === 'Show Password' && creds.password) {
          vscode.window.showInformationMessage(`Password: ${creds.password}`);
        }
      } catch (err: unknown) {
        const msg = err instanceof Error ? err.message : String(err);
        vscode.window.showErrorMessage(`Failed to reveal credentials: ${msg}`);
      }
    })
  );
}
