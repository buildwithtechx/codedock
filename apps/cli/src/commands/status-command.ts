import { ApiClient } from '../api-client.js';
import { printError, printJson, printTable } from '../output-format.js';
import type { CliContext, DeploymentRecord, ServiceRecord } from '../types.js';

export async function statusCommand(ctx: CliContext, args: string[]): Promise<void> {
  const serviceId = args[0];
  if (!serviceId) {
    printError('Usage: codedock status <service-id>');
    process.exit(1);
  }

  if (!ctx.token) {
    printError('Not logged in. Please run "codedock login" first.');
    process.exit(1);
  }

  const client = new ApiClient(ctx);
  try {
    const serviceRes = await client.get<ServiceRecord>(`/api/services/${serviceId}`);
    const service = serviceRes.data;
    if (!service) {
      throw new Error(`Service ${serviceId} not found`);
    }

    let deployments: DeploymentRecord[] = [];
    try {
      const depRes = await client.get<DeploymentRecord[]>(`/api/services/${serviceId}/deployments`);
      deployments = depRes.data || [];
    } catch {
      deployments = [];
    }

    if (ctx.json) {
      printJson({ service, deployments });
      return;
    }

    console.log(`\x1b[1mService:\x1b[0m       ${service.name} (${service.id})`);
    console.log(`\x1b[1mStatus:\x1b[0m        ${service.status}`);
    if (service.domain) {
      console.log(`\x1b[1mDomain:\x1b[0m        ${service.domain}`);
    }
    if (service.repositoryUrl) {
      console.log(
        `\x1b[1mRepository:\x1b[0m    ${service.repositoryUrl} (branch: ${service.branch || 'main'})`
      );
    }
    console.log(`\x1b[1mCreated At:\x1b[0m    ${service.createdAt}`);
    console.log('');

    if (deployments.length > 0) {
      console.log('\x1b[1mRecent Deployments:\x1b[0m');
      const rows = deployments
        .slice(0, 5)
        .map((d) => [
          d.id.slice(0, 8),
          d.status,
          d.branch || 'main',
          d.commitHash ? d.commitHash.slice(0, 7) : 'latest',
          d.createdAt,
        ]);
      printTable(['ID', 'STATUS', 'BRANCH', 'COMMIT', 'CREATED'], rows);
    } else {
      console.log('No deployments found for this service.');
    }
  } catch (err: unknown) {
    const msg = err instanceof Error ? err.message : String(err);
    printError(`Failed to fetch service status: ${msg}`);
    process.exit(1);
  }
}
