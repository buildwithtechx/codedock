import { execFileSync } from 'node:child_process';
import { mkdtempSync, readFileSync, rmSync, statSync } from 'node:fs';
import { tmpdir } from 'node:os';
import { basename, join } from 'node:path';
import type { ApiClient } from './api-client.js';

export async function uploadDirectory(
  client: ApiClient,
  path: string,
  projectId: string
): Promise<{ appId: string; appName: string }> {
  const directory = mkdtempSync(join(tmpdir(), 'codedock-deploy-'));
  try {
    const archive = join(directory, 'source.tar.gz');
    execFileSync(
      'tar',
      [
        '-czf',
        archive,
        '--exclude=.git',
        '--exclude=node_modules',
        '--exclude=.codedock',
        '-C',
        path,
        '.',
      ],
      { stdio: 'pipe' }
    );
    if (statSync(archive).size > 500 * 1024 * 1024)
      throw new Error('Deployment archive exceeds the 500 MB limit');
    const body = new FormData();
    body.set('projectId', projectId);
    body.set('name', basename(path));
    body.set(
      'file',
      new Blob([readFileSync(archive)], { type: 'application/gzip' }),
      'source.tar.gz'
    );
    const result = await client.post<{ appId: string; appName: string }>(
      '/api/deploy/archive',
      body
    );
    if (!result.data?.appId) throw new Error('Server did not return a deployed app');
    return result.data;
  } finally {
    rmSync(directory, { recursive: true, force: true });
  }
}
