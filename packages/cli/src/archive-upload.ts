import { spawn } from 'node:child_process';
import { randomBytes } from 'node:crypto';
import { basename } from 'node:path';
import { Readable } from 'node:stream';
import type { ApiClient } from './api-client.js';

const archiveLimit = 500 * 1024 * 1024;

export function createArchiveUpload(path: string, projectId: string, limit = archiveLimit) {
  const boundary = `codedock-${randomBytes(16).toString('hex')}`;
  const child = spawn(
    'tar',
    [
      '-czf',
      '-',
      '--exclude=.git',
      '--exclude=node_modules',
      '--exclude=.codedock',
      '-C',
      path,
      '.',
    ],
    { stdio: ['ignore', 'pipe', 'pipe'] }
  );
  let stderr = '';
  child.stderr.on('data', (chunk: Buffer) => {
    stderr = (stderr + chunk.toString()).slice(-4096);
  });
  const completion = new Promise<void>((resolve, reject) => {
    child.once('error', reject);
    child.once('close', (code) => {
      if (code === 0) resolve();
      else
        reject(
          new Error(`Could not package deployment archive: ${stderr || `tar exited with ${code}`}`)
        );
    });
  });
  void completion.catch(() => {});
  const prefix = `--${boundary}\r\nContent-Disposition: form-data; name="projectId"\r\n\r\n${projectId}\r\n--${boundary}\r\nContent-Disposition: form-data; name="name"\r\n\r\n${basename(path)}\r\n--${boundary}\r\nContent-Disposition: form-data; name="file"; filename="source.tar.gz"\r\nContent-Type: application/gzip\r\n\r\n`;
  const stream = Readable.from(
    (async function* () {
      try {
        yield Buffer.from(prefix);
        let size = 0;
        for await (const chunk of child.stdout) {
          const bytes = chunk as Buffer;
          size += bytes.length;
          if (size > limit) throw new Error('Deployment archive exceeds the 500 MB limit');
          yield bytes;
        }
        await completion;
        yield Buffer.from(`\r\n--${boundary}--\r\n`);
      } finally {
        child.kill();
      }
    })(),
    { objectMode: false }
  );
  return {
    stream,
    contentType: `multipart/form-data; boundary=${boundary}`,
    async dispose() {
      stream.destroy();
      child.kill();
      await completion.catch(() => {});
    },
  };
}

export async function uploadDirectory(
  client: ApiClient,
  path: string,
  projectId: string
): Promise<{ appId: string; appName: string }> {
  const upload = createArchiveUpload(path, projectId);
  try {
    const result = await client.postStream<{ appId: string; appName: string }>(
      '/api/deploy/archive',
      Readable.toWeb(upload.stream) as ReadableStream<Uint8Array>,
      upload.contentType
    );
    if (!result.data?.appId) throw new Error('Server did not return a deployed app');
    return result.data;
  } finally {
    await upload.dispose();
  }
}
