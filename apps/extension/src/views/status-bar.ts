import * as vscode from 'vscode';
import type { ConnectionManager } from '../auth/connection-manager';

export class CodedockStatusBar implements vscode.Disposable {
  private readonly item: vscode.StatusBarItem;
  private readonly subscription: vscode.Disposable;

  constructor(private readonly connectionManager: ConnectionManager) {
    this.item = vscode.window.createStatusBarItem(vscode.StatusBarAlignment.Left, 100);

    this.subscription = this.connectionManager.onDidChangeConnection(() => {
      this.update();
    });

    this.update();
    this.item.show();
  }

  async update(): Promise<void> {
    const conn = await this.connectionManager.getActiveConnection();

    if (!conn) {
      this.item.text = '$(cloud) Codedock: Disconnected';
      this.item.tooltip = 'Click to connect to a Codedock instance';
      this.item.command = 'codedock.connect';
      return;
    }

    let host = 'Connected';
    try {
      const parsed = new URL(conn.url);
      host = parsed.host;
    } catch {
      host = conn.url;
    }

    const label = conn.instanceName || host;
    this.item.text = `$(cloud) Codedock: ${label}`;
    this.item.tooltip = `Instance: ${conn.url}\nUser: ${conn.userEmail || 'Authenticated'}\nVersion: ${conn.version || 'unknown'}\nClick for options`;
    this.item.command = 'codedock.showQuickMenu';
  }

  dispose(): void {
    this.subscription.dispose();
    this.item.dispose();
  }
}
