import * as vscode from 'vscode';

export class OutputChannels implements vscode.Disposable {
  private static _instance: OutputChannels;

  readonly runtimeLogs: vscode.OutputChannel;
  readonly deploymentLogs: vscode.OutputChannel;

  private constructor() {
    this.runtimeLogs = vscode.window.createOutputChannel('Codedock: Runtime Logs');
    this.deploymentLogs = vscode.window.createOutputChannel('Codedock: Deployments');
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
  }
}
