import { execFileSync } from 'node:child_process';
import { createWriteStream, promises as fs } from 'node:fs';
import { dirname, join } from 'node:path';
import { pipeline } from 'node:stream/promises';
import { fileURLToPath } from 'node:url';
import { downloadUrl, mapPlatform } from './lib.js';

const root = dirname(fileURLToPath(import.meta.url));
const binDir = join(root, 'bin');

async function main() {
  if (process.env.CODEDOCK_SKIP_POSTINSTALL === '1') {
    console.log('codedock: skipping binary download (CODEDOCK_SKIP_POSTINSTALL=1)');
    return;
  }
  const { version } = JSON.parse(await fs.readFile(join(root, 'package.json'), 'utf8'));
  const { os, arch, binary } = mapPlatform(process.platform, process.arch);
  const url = downloadUrl(version, os, arch);
  const archive = join(binDir, `codedock-${os}-${arch}.tar.gz`);
  console.log(`codedock: downloading ${url}`);
  await fs.mkdir(binDir, { recursive: true });
  const response = await fetch(url, { redirect: 'follow' });
  if (!response.ok || !response.body) {
    throw new Error(`download failed with status ${response.status}`);
  }
  await pipeline(response.body, createWriteStream(archive));
  execFileSync('tar', ['-xzf', archive, '-C', binDir], { stdio: 'inherit' });
  await fs.unlink(archive);
  if (os !== 'windows') {
    await fs.chmod(join(binDir, binary), 0o755);
  }
  console.log(`codedock: installed ${binary} ${version} (${os}/${arch})`);
}

try {
  await main();
} catch (error) {
  console.warn(`codedock: WARNING: binary download failed (${error.message}).`);
  console.warn('codedock: the CLI will not run until a binary is installed.');
  console.warn(
    'codedock: check your network, then reinstall, or install from https://github.com/buildwithtechx/codedock/releases'
  );
}
