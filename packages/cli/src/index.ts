import { executeCli } from './cli-router.js';

export async function runCli(): Promise<void> {
  const args = process.argv.slice(2);
  await executeCli(args);
}

export { ApiClient } from './api-client.js';
export { executeCli } from './cli-router.js';
export { clearConfig, loadConfig, saveConfig } from './config-store.js';
export * from './types.js';
