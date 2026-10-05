import { ApiClient } from '../api-client.js';
import { printError, printJson, printSuccess, printTable } from '../output-format.js';
import type { CliContext } from '../types.js';

export async function envCommand(ctx: CliContext, args: string[]): Promise<void> {
  const subAction = args[0] || 'list';

  if (!ctx.token) {
    printError('Not logged in. Please run "codedock login" first.');
    process.exit(1);
  }

  const client = new ApiClient(ctx);

  if (subAction === 'list') {
    const projectId = args[1];
    if (!projectId) {
      printError('Usage: codedock env list <project-id>');
      process.exit(1);
    }
    try {
      const res = await client.get<Record<string, string>>(`/api/projects/${projectId}/env`);
      const vars = res.data || {};
      if (ctx.json) {
        printJson(vars);
        return;
      }
      const keys = Object.keys(vars);
      if (keys.length === 0) {
        console.log('No environment variables set.');
        return;
      }
      const rows = keys.map((k) => [k, vars[k]]);
      printTable(['KEY', 'VALUE'], rows);
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : String(err);
      printError(`Failed to list env vars: ${msg}`);
      process.exit(1);
    }
    return;
  }

  if (subAction === 'get') {
    const projectId = args[1];
    const key = args[2];
    if (!projectId || !key) {
      printError('Usage: codedock env get <project-id> <key>');
      process.exit(1);
    }
    try {
      const res = await client.get<Record<string, string>>(`/api/projects/${projectId}/env`);
      const vars = res.data || {};
      const val = vars[key];
      if (val === undefined) {
        printError(`Variable "${key}" not found in project`);
        process.exit(1);
      }
      console.log(val);
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : String(err);
      printError(`Failed to get env var: ${msg}`);
      process.exit(1);
    }
    return;
  }

  if (subAction === 'set') {
    const projectId = args[1];
    const pairs = args.slice(2);
    if (!projectId || pairs.length === 0) {
      printError('Usage: codedock env set <project-id> KEY1=VAL1 KEY2=VAL2 ...');
      process.exit(1);
    }
    const variables: Record<string, string> = {};
    for (const pair of pairs) {
      const eqIdx = pair.indexOf('=');
      if (eqIdx === -1) {
        printError(`Invalid variable pair: "${pair}". Must be KEY=VALUE`);
        process.exit(1);
      }
      const k = pair.slice(0, eqIdx).trim();
      const v = pair.slice(eqIdx + 1);
      variables[k] = v;
    }
    try {
      await client.put(`/api/projects/${projectId}/env`, variables);
      printSuccess(
        `Updated ${Object.keys(variables).length} environment variable(s) for project ${projectId}`
      );
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : String(err);
      printError(`Failed to set env vars: ${msg}`);
      process.exit(1);
    }
    return;
  }

  printError(`Unknown env action: "${subAction}". Available: list, get, set`);
  process.exit(1);
}
