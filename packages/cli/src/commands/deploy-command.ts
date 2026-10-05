import { existsSync, statSync } from 'node:fs';
import { resolve } from 'node:path';
import { ApiClient } from '../api-client.js';
import { uploadDirectory } from '../archive-upload.js';
import { printError, printInfo, printSuccess } from '../output-format.js';
import { listProjects } from '../project-list.js';
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
  if (!ctx.token) {
    printError('Not logged in. Please run "codedock login" first.');
    process.exitCode = 1;
    return;
  }
  const target = args[0] || '.';
  const client = new ApiClient(ctx);
  try {
    if (existsSync(target) && statSync(target).isDirectory()) {
      let projectId = options.project;
      if (!projectId) {
        const projects = await listProjects(client);
        if (projects.length > 1) throw new Error('Choose a target project with --project <id>');
        projectId = projects[0]?.id;
        if (!projectId)
          projectId = (
            await client.post<ProjectRecord>('/api/projects', {
              name: 'default-app',
              description: 'Created via Codedock CLI deploy',
            })
          ).data?.id;
      }
      if (!projectId) throw new Error('Could not resolve a target project');
      printInfo(`Packaging and uploading ${resolve(target)}`);
      const result = await uploadDirectory(client, resolve(target), projectId);
      printSuccess(`Deployed ${result.appName} (ID: ${result.appId})`);
      printInfo(`Monitor status with: codedock status ${result.appId}`);
      return;
    }
    const result = await client.post<DeploymentRecord>(
      `/api/services/${encodeURIComponent(target)}/deploy`,
      { branch: options.branch || 'main' }
    );
    if (!result.data?.id) throw new Error('Server did not return a deployment');
    printSuccess(`Deployment started (ID: ${result.data.id})`);
  } catch (err) {
    printError(err instanceof Error ? err.message : 'Deployment failed');
    process.exitCode = 1;
  }
}
