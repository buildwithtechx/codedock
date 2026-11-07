import * as vscode from 'vscode';
import { ConnectionManager } from './auth/connection-manager';
import { registerAllCommands, stopActiveLogStream } from './commands';
import { DeploymentsTreeProvider } from './providers/deployments-tree-provider';
import { ServersTreeProvider } from './providers/servers-tree-provider';
import { WorkloadsTreeProvider } from './providers/workloads-tree-provider';
import { OutputChannels } from './views/output-channels';
import { CodedockStatusBar } from './views/status-bar';

export function activate(context: vscode.ExtensionContext): void {
  const connectionManager = new ConnectionManager(context);
  const workloadsProvider = new WorkloadsTreeProvider(connectionManager);
  const serversProvider = new ServersTreeProvider(connectionManager);
  const deploymentsProvider = new DeploymentsTreeProvider(connectionManager);
  const outputChannels = OutputChannels.getInstance();
  const statusBar = new CodedockStatusBar(connectionManager);

  // Register Tree Views
  context.subscriptions.push(
    vscode.window.registerTreeDataProvider('codedock.views.workloads', workloadsProvider),
    vscode.window.registerTreeDataProvider('codedock.views.servers', serversProvider),
    vscode.window.registerTreeDataProvider('codedock.views.deployments', deploymentsProvider)
  );

  // Register Status Bar & Output Channels
  context.subscriptions.push(statusBar, outputChannels);

  // Register All Commands
  registerAllCommands(
    context,
    connectionManager,
    workloadsProvider,
    serversProvider,
    deploymentsProvider,
    outputChannels
  );
}

export function deactivate(): void {
  stopActiveLogStream();
}
