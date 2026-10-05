import { ApiClient } from '../api-client.js';
import { printError, printJson, printSuccess, printTable } from '../output-format.js';
import type { CliContext, ServerRecord } from '../types.js';

interface ServerCreateOptions {
  name?: string;
  ip?: string;
  port?: number | string;
  user?: string;
  key?: string;
  isLocal?: boolean;
}

export async function serversCommand(
  ctx: CliContext,
  args: string[],
  options: ServerCreateOptions
): Promise<void> {
  const subAction = args[0] || 'list';

  if (!ctx.token) {
    printError('Not logged in. Please run "codedock login" first.');
    process.exit(1);
  }

  const client = new ApiClient(ctx);

  if (subAction === 'list') {
    try {
      const res = await client.get<ServerRecord[]>('/api/servers');
      const servers = res.data || [];
      if (ctx.json) {
        printJson(servers);
        return;
      }
      if (servers.length === 0) {
        console.log('No servers found.');
        return;
      }
      const rows = servers.map((s) => [
        s.id.slice(0, 8),
        s.name,
        s.ipAddress || '127.0.0.1',
        s.isLocal ? 'Local Docker' : 'Remote SSH',
        s.status || 'active',
        s.dockerVersion || 'Docker Engine',
      ]);
      printTable(['ID', 'NAME', 'HOST', 'TYPE', 'STATUS', 'DOCKER'], rows);
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : String(err);
      printError(`Failed to list servers: ${msg}`);
      process.exit(1);
    }
    return;
  }

  if (subAction === 'create') {
    const name = options.name || args[1];
    if (!name) {
      printError(
        'Usage: codedock servers create --name <name> [--ip <host>] [--port <port>] [--user <user>] [--is-local]'
      );
      process.exit(1);
    }
    try {
      const payload = {
        name,
        ipAddress: options.ip || '127.0.0.1',
        sshPort: options.port ? Number(options.port) : 22,
        sshUser: options.user || 'root',
        sshKey: options.key || '',
        isLocal: options.isLocal ?? false,
      };
      const res = await client.post<ServerRecord>('/api/servers', payload);
      printSuccess(`Server "${res.data?.name || name}" created successfully (ID: ${res.data?.id})`);
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : String(err);
      printError(`Failed to create server: ${msg}`);
      process.exit(1);
    }
    return;
  }

  if (subAction === 'delete') {
    const serverId = args[1];
    if (!serverId) {
      printError('Usage: codedock servers delete <server-id>');
      process.exit(1);
    }
    try {
      await client.del(`/api/servers/${serverId}`);
      printSuccess(`Server ${serverId} deleted successfully`);
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : String(err);
      printError(`Failed to delete server: ${msg}`);
      process.exit(1);
    }
    return;
  }

  printError(`Unknown servers action: "${subAction}". Available: list, create, delete`);
  process.exit(1);
}
