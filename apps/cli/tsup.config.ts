import { defineConfig } from 'tsup';

export default defineConfig({
  clean: true,
  dts: false,
  entry: ['src/index.ts', 'src/bin.ts'],
  format: ['esm'],
  sourcemap: true,
  shims: true,
});
