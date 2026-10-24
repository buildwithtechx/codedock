import { describe, expect, it } from 'vitest';
import { parseCatalogEntry, parseScalar } from './generate-app-catalog.mjs';

const fixture = `x-codedock:
  is_one_click: true
  verified: true
  name: Supabase
  description: Full self-hosted Supabase — Postgres, Auth and Studio behind one gateway.
  category: backend
  icon: supabase
  connection_string: postgresql://postgres:\${secret:PASSWORD}@db:5432/postgres
  secrets:
    - var: POSTGRES_PASSWORD
      label: Database password
      length: 24
  inputs:
    - var: DASHBOARD_USERNAME
      label: Studio dashboard username
services:
  kong:
    image: kong:3.9.1
  db:
    image: supabase/postgres:17.6.1.136
`;

describe('catalog parser', () => {
  it('parses scalar metadata', () => {
    expect(parseScalar('  name: Supabase')).toEqual(['name', 'Supabase']);
    expect(parseScalar('  verified: true')).toEqual(['verified', true]);
    expect(parseScalar('    - var: X')).toBeNull();
  });

  it('extracts entry, services and secrets without inputs leakage', () => {
    const entry = parseCatalogEntry('supabase', fixture);
    expect(entry?.id).toBe('supabase');
    expect(entry?.name).toBe('Supabase');
    expect(entry?.category).toBe('backend');
    expect(entry?.icon).toBe('supabase');
    expect(entry?.verified).toBe(true);
    expect(entry?.services).toEqual(['kong', 'db']);
    expect(entry?.secrets).toEqual([
      { var: 'POSTGRES_PASSWORD', label: 'Database password', length: 24 },
    ]);
  });

  it('skips non-one-click entries', () => {
    expect(parseCatalogEntry('plain', 'services:\n  web:\n    image: x\n')).toBeNull();
  });
});
