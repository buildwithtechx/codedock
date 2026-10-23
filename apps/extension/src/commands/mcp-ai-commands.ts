import * as fs from 'node:fs/promises';
import * as path from 'node:path';
import * as vscode from 'vscode';
import { InfrastructureApi } from '../api';
import type { ConnectionManager } from '../auth/connection-manager';
import type { OutputChannels } from '../views/output-channels';

export function registerMcpAiCommands(
  context: vscode.ExtensionContext,
  connectionManager: ConnectionManager,
  outputChannels: OutputChannels
): void {
  context.subscriptions.push(
    vscode.commands.registerCommand('codedock.configureMcp', async () => {
      const conn = await connectionManager.getActiveConnection();
      if (!conn) {
        vscode.window.showWarningMessage('Connect to a Codedock instance first.');
        vscode.commands.executeCommand('codedock.connect');
        return;
      }

      const mcpEndpoint = `${conn.url.replace(/\/+$/, '')}/api/mcp`;
      const mcpConfig = {
        servers: {
          codedock: {
            type: 'http',
            url: mcpEndpoint,
            headers: {
              Authorization: `Bearer ${conn.token}`,
            },
          },
        },
      };
      const serialized = JSON.stringify(mcpConfig, null, 2);

      const action = await vscode.window.showQuickPick(
        [
          {
            label: '$(file-code) Write to Workspace (.vscode/mcp.json)',
            description: 'Configure MCP for GitHub Copilot & VS Code agents in this workspace',
            id: 'workspace',
          },
          {
            label: '$(clippy) Copy MCP JSON Configuration',
            description: 'Copy to clipboard for Cursor, Windsurf, or Claude Desktop',
            id: 'clipboard',
          },
        ],
        { placeHolder: `Configure Codedock MCP (${mcpEndpoint})` }
      );

      if (!action) return;

      if (action.id === 'clipboard') {
        await vscode.env.clipboard.writeText(serialized);
        vscode.window.showInformationMessage('Copied Codedock MCP configuration to clipboard.');
        return;
      }

      const folder = vscode.workspace.workspaceFolders?.[0];
      if (!folder) {
        await vscode.env.clipboard.writeText(serialized);
        vscode.window.showInformationMessage(
          'No workspace folder open. Copied MCP configuration to clipboard instead.'
        );
        return;
      }

      const vscodeDir = path.join(folder.uri.fsPath, '.vscode');
      const mcpFile = path.join(vscodeDir, 'mcp.json');
      await fs.mkdir(vscodeDir, { recursive: true });
      await fs.writeFile(mcpFile, `${serialized}\n`, 'utf8');

      const doc = await vscode.workspace.openTextDocument(mcpFile);
      await vscode.window.showTextDocument(doc);
      vscode.window.showInformationMessage('Codedock MCP server configured in .vscode/mcp.json.');
    }),

    vscode.commands.registerCommand('codedock.askAi', async () => {
      const conn = await connectionManager.getActiveConnection();
      if (!conn) {
        vscode.window.showWarningMessage('Connect to a Codedock instance first.');
        vscode.commands.executeCommand('codedock.connect');
        return;
      }

      const prompt = await vscode.window.showInputBox({
        title: 'Ask Codedock AI Copilot',
        prompt: 'Ask about your projects, services, deployments, databases, or servers',
        placeHolder: 'e.g. Inspect unhealthy services and recent deployments on this instance',
        ignoreFocusOut: true,
      });

      if (!prompt?.trim()) return;

      await vscode.window.withProgress(
        {
          location: vscode.ProgressLocation.Notification,
          title: 'Codedock Copilot: Inspecting instance...',
          cancellable: false,
        },
        async () => {
          try {
            const res = await InfrastructureApi.askAiAssistant(conn.url, conn.token, prompt.trim());
            const ch = outputChannels.aiCopilot;
            ch.clear();
            ch.appendLine(`=== Codedock AI Copilot (${conn.instanceName || conn.url}) ===`);
            ch.appendLine(`> ${prompt.trim()}\n`);
            if (res.steps && res.steps.length > 0) {
              ch.appendLine('Executed MCP Tools:');
              for (const s of res.steps) {
                ch.appendLine(`  • [${s.status}] ${s.tool} — ${s.label}`);
              }
              ch.appendLine('');
            }
            ch.appendLine(res.reply);
            ch.show(true);
          } catch (err: unknown) {
            const msg = err instanceof Error ? err.message : String(err);
            vscode.window.showErrorMessage(`Codedock AI error: ${msg}`);
          }
        }
      );
    })
  );
}
