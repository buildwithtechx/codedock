import assert from 'node:assert/strict';
import test from 'node:test';
import { artifactName, downloadUrl, mapPlatform, releaseTag } from './lib.js';

test('maps node platforms to release artifacts', () => {
  assert.deepEqual(mapPlatform('darwin', 'arm64'), {
    os: 'darwin',
    arch: 'arm64',
    binary: 'codedock',
  });
  assert.deepEqual(mapPlatform('linux', 'x64'), { os: 'linux', arch: 'amd64', binary: 'codedock' });
  assert.deepEqual(mapPlatform('win32', 'x64'), {
    os: 'windows',
    arch: 'amd64',
    binary: 'codedock.exe',
  });
});

test('rejects unsupported platforms', () => {
  assert.throws(() => mapPlatform('freebsd', 'x64'), /Unsupported platform/);
  assert.throws(() => mapPlatform('linux', 'mips'), /Unsupported platform/);
});

test('builds artifact names and download urls', () => {
  assert.equal(artifactName('linux', 'amd64'), 'codedock_linux_amd64.tar.gz');
  assert.equal(releaseTag('0.1.0'), 'v0.1.0');
  assert.equal(
    downloadUrl('0.1.0', 'darwin', 'arm64'),
    'https://github.com/buildwithtechx/codedock/releases/download/v0.1.0/codedock_darwin_arm64.tar.gz'
  );
});
