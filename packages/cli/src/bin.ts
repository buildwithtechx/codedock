#!/usr/bin/env node
import { runCli } from './index.js';

runCli().catch((err: unknown) => {
  const msg = err instanceof Error ? err.message : String(err);
  console.error(`Error: ${msg}`);
  process.exit(1);
});
