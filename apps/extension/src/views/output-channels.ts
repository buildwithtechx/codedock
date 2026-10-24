import * as vscode from 'vscode';

export class OutputChannels implements vscode.Disposable {
  private static _instance: OutputChannels;

  readonly runtimeLogs: vscode.OutputChannel;
  readonly deploymentLogs: vscode.OutputChannel;
  readonly aiCopilot: vscode.OutputChannel;

  private constructor() {
    this.runtimeLogs = vscode.window.createOutputChannel('Codedock: Runtime Logs');
    this.deploymentLogs = vscode.window.createOutputChannel('Codedock: Deployments');
    this.aiCopilot = vscode.window.createOutputChannel('Codedock: AI Copilot');
  }

  static getInstance(): OutputChannels {
    if (!OutputChannels._instance) {
      OutputChannels._instance = new OutputChannels();
    }
    return OutputChannels._instance;
  }

  dispose(): void {
    this.runtimeLogs.dispose();
    this.deploymentLogs.dispose();
    this.aiCopilot.dispose();
  }
}
