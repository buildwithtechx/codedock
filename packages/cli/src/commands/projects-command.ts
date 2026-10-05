import { ApiClient } from '../api-client.js';
import { printError, printJson, printSuccess, printTable } from '../output-format.js';
import { listProjects } from '../project-list.js';
import type { CliContext, ProjectRecord } from '../types.js';

interface ProjectOptions {
  name?: string;
  description?: string;
  organization?: string;
}

export async function projectsCommand(
  ctx: CliContext,
  args: string[],
  options: ProjectOptions
): Promise<void> {
  const subAction = args[0] || 'list';

  if (!ctx.token) {
    printError('Not logged in. Please run "codedock login" first.');
    process.exit(1);
  }

  const client = new ApiClient(ctx);

  if (subAction === 'list') {
    try {
      const projects = await listProjects(client, options.organization);
      if (ctx.json) {
        printJson(projects);
        return;
      }
      if (projects.length === 0) {
        console.log('No projects found.');
        return;
      }
      const rows = projects.map((p) => [
        p.id.slice(0, 8),
        p.name,
        p.description || '-',
        p.createdAt || '-',
      ]);
      printTable(['ID', 'NAME', 'DESCRIPTION', 'CREATED'], rows);
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : String(err);
      printError(`Failed to list projects: ${msg}`);
      process.exit(1);
    }
    return;
  }

  if (subAction === 'create') {
    const name = options.name || args[1];
    if (!name) {
      printError('Usage: codedock projects create <name> [--description <desc>]');
      process.exit(1);
    }
    try {
      const res = await client.post<ProjectRecord>('/api/projects', {
        name,
        organizationId: options.organization,
        description: options.description || '',
      });
      printSuccess(
        `Project "${res.data?.name || name}" created successfully (ID: ${res.data?.id})`
      );
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : String(err);
      printError(`Failed to create project: ${msg}`);
      process.exit(1);
    }
    return;
  }

  if (subAction === 'delete') {
    const projectId = args[1];
    if (!projectId) {
      printError('Usage: codedock projects delete <project-id>');
      process.exit(1);
    }
    try {
      await client.del(`/api/projects/${projectId}`);
      printSuccess(`Project ${projectId} deleted successfully`);
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : String(err);
      printError(`Failed to delete project: ${msg}`);
      process.exit(1);
    }
    return;
  }

  printError(`Unknown projects action: "${subAction}". Available: list, create, delete`);
  process.exit(1);
}
