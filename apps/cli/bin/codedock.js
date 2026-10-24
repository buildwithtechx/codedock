#!/usr/bin/env node
import { spawnSync } from 'node:child_process';
import { existsSync } from 'node:fs';
import { dirname, join } from 'node:path';
import { fileURLToPath } from 'node:url';
import { mapPlatform } from '../lib.js';

const root = dirname(dirname(fileURLToPath(import.meta.url)));
const binDir = join(root, 'bin');

function fail(message) {
  console.error(`codedock: ${message}`);
  process.exit(1);
}

let binary;
try {
  binary = join(binDir, mapPlatform(process.platform, process.arch).binary);
} catch {
  fail(`unsupported platform ${process.platform}/${process.arch}`);
}
if (!existsSync(binary)) {
  fail(
    'binary is missing (postinstall download did not run or failed). ' +
      'Reinstall with network access, or download a release from https://github.com/buildwithtechx/codedock/releases'
  );
}
const result = spawnSync(binary, process.argv.slice(2), { stdio: 'inherit' });
process.exit(result.status ?? 1);
