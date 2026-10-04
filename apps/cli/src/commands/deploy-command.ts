import { existsSync, statSync } from 'node:fs';
import { resolve } from 'node:path';
import { ApiClient } from '../api-client.js';
import { printError, printInfo, printSuccess } from '../output-format.js';
import type { CliContext, DeploymentRecord, ProjectRecord } from '../types.js';

interface DeployOptions {
  project?: string;
  branch?: string;
}

export async function deployCommand(
  ctx: CliContext,
  args: string[],
  options: DeployOptions
): Promise<void> {
  const target = args[0] || '.';

  if (!ctx.token) {
    printError('Not logged in. Please run "codedock login" first.');
    process.exit(1);
  }

  const client = new ApiClient(ctx);

  if (existsSync(target) && statSync(target).isDirectory()) {
    const absPath = resolve(target);
    printInfo(`Deploying project from directory: ${absPath}`);

    let projectId = options.project;
    if (!projectId) {
      try {
        const projRes = await client.get<ProjectRecord[]>('/api/projects');
        const projects = projRes.data || [];
        if (projects.length > 0) {
          projectId = projects[0].id;
        } else {
          const createRes = await client.post<ProjectRecord>('/api/projects', {
            name: 'default-app',
            description: 'Created via Codedock CLI deploy',
          });
          projectId = createRes.data?.id;
        }
      } catch (err: unknown) {
        const msg = err instanceof Error ? err.message : String(err);
        printError(`Failed to resolve or create project for deploy: ${msg}`);
        process.exit(1);
      }
    }

    try {
      const deployRes = await client.post<DeploymentRecord>(`/api/projects/${projectId}/deploy`, {
        branch: options.branch || 'main',
      });
      printSuccess(`Deployment triggered successfully (ID: ${deployRes.data?.id || 'queued'})`);
      printInfo(`Monitor status with: codedock status ${projectId}`);
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : String(err);
      printError(`Deployment failed: ${msg}`);
      process.exit(1);
    }
    return;
  }

  printInfo(`Triggering deployment for service ID: ${target}`);
  try {
    const deployRes = await client.post<DeploymentRecord>('/api/deployments', {
      serviceId: target,
      branch: options.branch || 'main',
    });
    printSuccess(
      `Deployment triggered successfully for service ${target} (ID: ${deployRes.data?.id})`
    );
  } catch (err: unknown) {
    const msg = err instanceof Error ? err.message : String(err);
    printError(`Failed to trigger deployment: ${msg}`);
    process.exit(1);
  }
}
