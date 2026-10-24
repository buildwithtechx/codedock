import type * as vscode from 'vscode';
import type { ConnectionManager } from '../auth/connection-manager';
import type { DeploymentsTreeProvider } from '../providers/deployments-tree-provider';
import type { ServersTreeProvider } from '../providers/servers-tree-provider';
import type { WorkloadsTreeProvider } from '../providers/workloads-tree-provider';
import type { OutputChannels } from '../views/output-channels';
import { registerAppCommands } from './app-commands';
import { registerAuthCommands } from './auth-commands';
import { registerDatabaseCommands } from './database-commands';
import { registerMcpAiCommands } from './mcp-ai-commands';
import { registerServerCommands } from './server-commands';

export function registerAllCommands(
  context: vscode.ExtensionContext,
  connectionManager: ConnectionManager,
  workloadsProvider: WorkloadsTreeProvider,
  serversProvider: ServersTreeProvider,
  deploymentsProvider: DeploymentsTreeProvider,
  outputChannels: OutputChannels
): void {
  const refreshAll = () => {
    workloadsProvider.refresh();
    serversProvider.refresh();
    deploymentsProvider.refresh();
  };

  registerAuthCommands(context, connectionManager);
  registerAppCommands(context, connectionManager, outputChannels, refreshAll);
  registerDatabaseCommands(context, connectionManager, refreshAll);
  registerServerCommands(context, connectionManager, outputChannels, refreshAll);
  registerMcpAiCommands(context, connectionManager, outputChannels);
}

export { stopActiveLogStream } from './app-commands';
