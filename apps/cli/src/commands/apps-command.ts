import { ApiClient } from '../api-client.js';
import { printError, printJson, printSuccess, printTable } from '../output-format.js';
import type { CliContext, DeploymentRecord, ServiceRecord } from '../types.js';

interface AppOptions {
  name?: string;
  project?: string;
  repo?: string;
  branch?: string;
}

export async function appsCommand(
  ctx: CliContext,
  args: string[],
  options: AppOptions
): Promise<void> {
  const subAction = args[0] || 'list';

  if (!ctx.token) {
    printError('Not logged in. Please run "codedock login" first.');
    process.exit(1);
  }

  const client = new ApiClient(ctx);

  if (subAction === 'list') {
    try {
      const endpoint = options.project
        ? `/api/projects/${options.project}/services`
        : '/api/services';
      const res = await client.get<ServiceRecord[]>(endpoint);
      const apps = res.data || [];
      if (ctx.json) {
        printJson(apps);
        return;
      }
      if (apps.length === 0) {
        console.log('No apps/services found.');
        return;
      }
      const rows = apps.map((a) => [
        a.id.slice(0, 8),
        a.name,
        a.status || 'stopped',
        a.domain || '-',
        a.repositoryUrl ? `${a.repositoryUrl}#${a.branch || 'main'}` : '-',
      ]);
      printTable(['ID', 'NAME', 'STATUS', 'DOMAIN', 'REPO'], rows);
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : String(err);
      printError(`Failed to list apps: ${msg}`);
      process.exit(1);
    }
    return;
  }

  if (subAction === 'create') {
    const name = options.name || args[1];
    if (!name) {
      printError(
        'Usage: codedock apps create <name> [--project <id>] [--repo <url>] [--branch <branch>]'
      );
      process.exit(1);
    }
    try {
      const res = await client.post<ServiceRecord>('/api/services', {
        name,
        projectId: options.project,
        repositoryUrl: options.repo,
        branch: options.branch || 'main',
      });
      printSuccess(`App "${res.data?.name || name}" created successfully (ID: ${res.data?.id})`);
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : String(err);
      printError(`Failed to create app: ${msg}`);
      process.exit(1);
    }
    return;
  }

  if (subAction === 'delete') {
    const appId = args[1];
    if (!appId) {
      printError('Usage: codedock apps delete <app-id>');
      process.exit(1);
    }
    try {
      await client.del(`/api/services/${appId}`);
      printSuccess(`App ${appId} deleted successfully`);
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : String(err);
      printError(`Failed to delete app: ${msg}`);
      process.exit(1);
    }
    return;
  }

  if (subAction === 'logs') {
    const appId = args[1];
    if (!appId) {
      printError('Usage: codedock apps logs <app-id>');
      process.exit(1);
    }
    try {
      const res = await client.get<string | { logs?: string[] }>(`/api/services/${appId}/logs`);
      if (typeof res.data === 'string') {
        console.log(res.data);
      } else if (Array.isArray(res.data?.logs)) {
        for (const line of res.data.logs) {
          console.log(line);
        }
      } else {
        printJson(res.data);
      }
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : String(err);
      printError(`Failed to fetch logs: ${msg}`);
      process.exit(1);
    }
    return;
  }

  if (subAction === 'deployments') {
    const appId = args[1];
    if (!appId) {
      printError('Usage: codedock apps deployments <app-id>');
      process.exit(1);
    }
    try {
      const res = await client.get<DeploymentRecord[]>(`/api/services/${appId}/deployments`);
      const deployments = res.data || [];
      if (ctx.json) {
        printJson(deployments);
        return;
      }
      if (deployments.length === 0) {
        console.log('No deployments found.');
        return;
      }
      const rows = deployments.map((d) => [
        d.id.slice(0, 8),
        d.status,
        d.branch || 'main',
        d.commitHash ? d.commitHash.slice(0, 7) : '-',
        d.createdAt,
      ]);
      printTable(['ID', 'STATUS', 'BRANCH', 'COMMIT', 'CREATED'], rows);
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : String(err);
      printError(`Failed to list deployments: ${msg}`);
      process.exit(1);
    }
    return;
  }

  printError(
    `Unknown apps action: "${subAction}". Available: list, create, delete, logs, deployments`
  );
  process.exit(1);
}
