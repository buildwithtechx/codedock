import { execFileSync } from 'node:child_process';
import { copyFileSync, cpSync, existsSync, rmSync } from 'node:fs';
import { dirname, join, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';

const __dirname = dirname(fileURLToPath(import.meta.url));
const root = resolve(__dirname, '../../..');
const pkg = resolve(__dirname, '..');

execFileSync('npm', ['run', 'build', '--workspace=@codedock/sdk'], {
  cwd: root,
  stdio: 'inherit',
});

rmSync(join(pkg, 'dist'), { recursive: true, force: true });
cpSync(join(root, 'packages/sdk/dist'), join(pkg, 'dist/sdk'), { recursive: true });

// Create .d.cts declaration files matching .d.ts for CJS exports
if (existsSync(join(pkg, 'dist/sdk/index.d.ts'))) {
  copyFileSync(join(pkg, 'dist/sdk/index.d.ts'), join(pkg, 'dist/sdk/index.d.cts'));
}

if (existsSync(join(pkg, 'dist/sdk/client.d.ts'))) {
  copyFileSync(join(pkg, 'dist/sdk/client.d.ts'), join(pkg, 'dist/sdk/client.d.cts'));
}

if (existsSync(join(root, 'LICENSE'))) {
  copyFileSync(join(root, 'LICENSE'), join(pkg, 'LICENSE'));
}
