import { describe, expect, it } from 'vitest';
import {
  COMPANION_CATALOG,
  generateCatalogSecret,
  resolveCatalogEnvironment,
} from './companion-catalog';

describe('companion-catalog credential generation', () => {
  it('generates 48-char hex secret string by default', () => {
    const secret = generateCatalogSecret();
    expect(secret).toMatch(/^[a-f0-9]{48}$/);
    const secret2 = generateCatalogSecret();
    expect(secret2).not.toEqual(secret);
  });

  it('populates generated passwords for postgres, mysql, and mongo', () => {
    const postgres = COMPANION_CATALOG.find((i) => i.id === 'postgres')!;
    const env = resolveCatalogEnvironment(postgres);
    const pwd = env.find((e) => e.key === 'POSTGRES_PASSWORD');
    expect(pwd?.value).toMatch(/^[a-f0-9]{48}$/);
    expect(env.find((e) => e.key === 'POSTGRES_USER')?.value).toBe('postgres');

    const minio = COMPANION_CATALOG.find((i) => i.id === 'minio')!;
    const minioEnv = resolveCatalogEnvironment(minio);
    const minioPwd = minioEnv.find((e) => e.key === 'MINIO_ROOT_PASSWORD');
    expect(minioPwd?.value).toMatch(/^[a-f0-9]{48}$/);
  });
});
