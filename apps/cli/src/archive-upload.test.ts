import { mkdtempSync, rmSync, writeFileSync } from 'node:fs';
import { createServer } from 'node:http';
import { tmpdir } from 'node:os';
import { join } from 'node:path';
import { afterEach, beforeEach, describe, expect, it } from 'vitest';
import { ApiClient } from './api-client.js';
import { createArchiveUpload, uploadDirectory } from './archive-upload.js';

let directory: string;
beforeEach(() => {
  directory = mkdtempSync(join(tmpdir(), 'codedock-stream-test-'));
  writeFileSync(join(directory, 'package.json'), '{"name":"stream-test"}');
});
afterEach(() => rmSync(directory, { recursive: true, force: true }));

describe('archive uploads', () => {
  it('streams multipart content to the server', async () => {
    let received = '';
    let contentType = '';
    let contentLength: string | undefined;
    const server = createServer((request, response) => {
      contentType = request.headers['content-type'] || '';
      contentLength = request.headers['content-length'];
      request.on('data', (chunk: Buffer) => {
        received += chunk.toString('latin1');
      });
      request.on('end', () => {
        response.setHeader('Content-Type', 'application/json');
        response.end(JSON.stringify({ data: { appId: 'service', appName: 'uploaded' } }));
      });
    });
    await new Promise<void>((resolve) => server.listen(0, '127.0.0.1', resolve));
    try {
      const address = server.address();
      if (!address || typeof address === 'string') throw new Error('No server address');
      const client = new ApiClient({
        config: {},
        serverUrl: `http://127.0.0.1:${address.port}`,
        token: 'test',
        json: false,
      });
      expect(await uploadDirectory(client, directory, 'project')).toEqual({
        appId: 'service',
        appName: 'uploaded',
      });
      expect(contentType).toMatch(/^multipart\/form-data; boundary=/);
      expect(contentLength).toBeUndefined();
      expect(received).toContain('filename="source.tar.gz"');
      expect(received).toContain('project');
    } finally {
      await new Promise<void>((resolve, reject) =>
        server.close((error) => (error ? reject(error) : resolve()))
      );
    }
  });

  it('stops packaging when the compressed stream exceeds the limit', async () => {
    const upload = createArchiveUpload(directory, 'project', 1);
    try {
      await expect(
        (async () => {
          for await (const _chunk of upload.stream) {
          }
        })()
      ).rejects.toThrow('exceeds');
    } finally {
      await upload.dispose();
    }
  });
});
