import { existsSync, mkdirSync, readFileSync, rmSync, writeFileSync } from 'node:fs';
import { homedir } from 'node:os';
import { dirname, join } from 'node:path';
import type { CliConfig } from './types.js';

export function getConfigPath(): string {
  const home = process.env.CODEDOCK_CONFIG_DIR || homedir();
  return join(home, '.codedock', 'config.json');
}

export function loadConfig(): CliConfig {
  const configPath = getConfigPath();
  if (!existsSync(configPath)) {
    return {};
  }
  try {
    const raw = readFileSync(configPath, 'utf8');
    return JSON.parse(raw) as CliConfig;
  } catch {
    return {};
  }
}

export function saveConfig(cfg: CliConfig): void {
  const configPath = getConfigPath();
  const dir = dirname(configPath);
  if (!existsSync(dir)) {
    mkdirSync(dir, { recursive: true });
  }
  writeFileSync(configPath, JSON.stringify(cfg, null, 2), 'utf8');
}

export function clearConfig(): void {
  const configPath = getConfigPath();
  if (existsSync(configPath)) {
    try {
      rmSync(configPath);
    } catch {
      writeFileSync(configPath, JSON.stringify({}, null, 2), 'utf8');
    }
  }
}
