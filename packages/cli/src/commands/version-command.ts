import type { CliContext } from '../types.js';

export const CLI_VERSION = '0.1.0';

export async function versionCommand(_ctx: CliContext): Promise<void> {
  console.log(
    `codedock v${CLI_VERSION} (${process.platform}-${process.arch}) Node ${process.version}`
  );
}
