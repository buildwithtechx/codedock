import assert from 'node:assert/strict';
import { execFileSync } from 'node:child_process';
import { mkdtempSync, readdirSync, rmSync, writeFileSync } from 'node:fs';
import { tmpdir } from 'node:os';
import { dirname, join, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';

const __dirname = dirname(fileURLToPath(import.meta.url));
const root = resolve(__dirname, '..');
const tmpDir = mkdtempSync(join(tmpdir(), 'codedock-verify-'));

try {
  execFileSync(
    'npm',
    ['pack', '--workspace=codedock', `--pack-destination=${tmpDir}`],
    { cwd: root, stdio: 'inherit' }
  );

  const files = readdirSync(tmpDir);
  const tarballName = files.find((f) => f.endsWith('.tgz'));
  assert(tarballName, 'Failed to find generated tarball in temporary directory');
  const tarballPath = join(tmpDir, tarballName);

  console.log(`Installing tarball ${tarballName} into test environment...`);
  writeFileSync(
    join(tmpDir, 'package.json'),
    JSON.stringify({
      name: 'codedock-verify-sandbox',
      type: 'module',
      dependencies: {
        codedock: `file:${tarballPath}`,
      },
    }, null, 2)
  );

  execFileSync(
    'npm',
    ['install', '--ignore-scripts', '--no-audit', '--no-fund', '--no-package-lock'],
    { cwd: tmpDir, stdio: 'inherit' }
  );

  console.log('Testing ESM imports...');
  const esmTestScript = `
    import assert from 'node:assert/strict';
    const sdk = await import('codedock');
    assert.equal(typeof sdk.createCodedockClient, 'function', 'createCodedockClient should be a function');
    assert.equal(typeof sdk.CodedockClient, 'function', 'CodedockClient should be a function');
    assert.equal(typeof sdk.CodedockApiError, 'function', 'CodedockApiError should be a function');
    assert.equal(typeof sdk.CodedockRateLimitError, 'function', 'CodedockRateLimitError should be a function');
    assert.equal(typeof sdk.CodedockBillingGateError, 'function', 'CodedockBillingGateError should be a function');

    const subpath = await import('codedock/client');
    assert.equal(typeof subpath.CodedockClient, 'function', 'subpath CodedockClient should be a function');
    assert.equal(typeof subpath.createCodedockClient, 'function', 'subpath createCodedockClient should be a function');

    console.log('✓ ESM imports verified successfully');
  `;
  execFileSync('node', ['--input-type=module', '-e', esmTestScript], { cwd: tmpDir, stdio: 'inherit' });

  console.log('Testing CJS imports...');
  const cjsTestScript = `
    const assert = require('node:assert/strict');
    const sdk = require('codedock');
    assert.equal(typeof sdk.createCodedockClient, 'function', 'CJS createCodedockClient should be a function');
    assert.equal(typeof sdk.CodedockClient, 'function', 'CJS CodedockClient should be a function');

    const subpath = require('codedock/client');
    assert.equal(typeof subpath.CodedockClient, 'function', 'CJS subpath CodedockClient should be a function');

    console.log('✓ CJS imports verified successfully');
  `;
  execFileSync('node', ['-e', cjsTestScript], { cwd: tmpDir, stdio: 'inherit' });

  console.log('Testing TypeScript declarations...');
  const tsProbe = `
    import { createCodedockClient, type CodedockClientOptions } from 'codedock';
    import { CodedockClient } from 'codedock/client';

    const opts: CodedockClientOptions = { baseUrl: 'https://example.com' };
    const client = createCodedockClient(opts);
    const direct = new CodedockClient(opts);
    void client;
    void direct;
  `;
  writeFileSync(join(tmpDir, 'probe.ts'), tsProbe);
  const tscBin = resolve(root, 'node_modules/.bin/tsc');
  execFileSync(
    tscBin,
    ['--noEmit', '--skipLibCheck', '--moduleResolution', 'bundler', '--module', 'ESNext', join(tmpDir, 'probe.ts')],
    { cwd: tmpDir, stdio: 'inherit' }
  );
  console.log('✓ TypeScript declarations verified successfully');

  console.log('\n✅ All package verifications passed successfully!');
} finally {
  rmSync(tmpDir, { recursive: true, force: true });
}
