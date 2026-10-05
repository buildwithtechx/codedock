import { clearConfig } from '../config-store.js';
import { printSuccess } from '../output-format.js';
import type { CliContext } from '../types.js';

export async function logoutCommand(_ctx: CliContext): Promise<void> {
  clearConfig();
  printSuccess('Logged out successfully and cleared local credentials');
}
