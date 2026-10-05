import { ApiClient } from '../api-client.js';
import { printError, printJson, printSuccess, printTable } from '../output-format.js';
import type { CliContext, DatabaseRecord } from '../types.js';

interface DatabaseOptions {
  name?: string;
  type?: string;
  server?: string;
  project?: string;
}

export async function databasesCommand(
  ctx: CliContext,
  args: string[],
  options: DatabaseOptions
): Promise<void> {
  const subAction = args[0] || 'list';

  if (!ctx.token) {
    printError('Not logged in. Please run "codedock login" first.');
    process.exit(1);
  }

  const client = new ApiClient(ctx);

  if (subAction === 'list') {
    try {
      const res = await client.get<DatabaseRecord[]>('/api/databases');
      const databases = res.data || [];
      if (ctx.json) {
        printJson(databases);
        return;
      }
      if (databases.length === 0) {
        console.log('No managed databases found.');
        return;
      }
      const rows = databases.map((d) => [
        d.id.slice(0, 8),
        d.name,
        d.engine,
        d.status || 'running',
        d.createdAt || '-',
      ]);
      printTable(['ID', 'NAME', 'TYPE', 'STATUS', 'CREATED'], rows);
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : String(err);
      printError(`Failed to list databases: ${msg}`);
      process.exit(1);
    }
    return;
  }

  if (subAction === 'create') {
    const name = options.name || args[1];
    const type = options.type || 'postgres';
    if (!name || !options.project) {
      printError(
        'Usage: codedock db create <name> --project <id> --type <postgres|mysql|redis|mongodb>'
      );
      process.exit(1);
    }
    try {
      const res = await client.post<DatabaseRecord>('/api/databases', {
        name,
        engine: type,
        projectId: options.project,
        serverId: options.server,
      });
      printSuccess(
        `Database "${res.data?.name || name}" (${type}) created successfully (ID: ${res.data?.id})`
      );
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : String(err);
      printError(`Failed to create database: ${msg}`);
      process.exit(1);
    }
    return;
  }

  if (subAction === 'backup') {
    const dbId = args[1];
    if (!dbId) {
      printError('Usage: codedock db backup <db-id>');
      process.exit(1);
    }
    try {
      await client.post(`/api/databases/${dbId}/backups`);
      printSuccess(`Backup triggered successfully for database ${dbId}`);
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : String(err);
      printError(`Failed to trigger database backup: ${msg}`);
      process.exit(1);
    }
    return;
  }

  printError(`Unknown db action: "${subAction}". Available: list, create, backup`);
  process.exit(1);
}
