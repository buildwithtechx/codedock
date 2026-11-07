import { describe, expect, it } from 'vitest';
import {
  CAPABILITIES,
  capabilitiesForRoot,
  capabilitiesImpliedBy,
  capabilityAction,
  capabilityFor,
  isCapability,
} from '../src/index';

describe('access-capabilities', () => {
  it('validates capabilities', () => {
    expect(isCapability('domains:read')).toBe(true);
    expect(isCapability('domains:write')).toBe(true);
    expect(isCapability('domains:delete')).toBe(true);
    expect(isCapability('deployments:read')).toBe(true);
    expect(isCapability('deployments:trigger')).toBe(true);
    expect(isCapability('invalid:action')).toBe(false);
  });

  it('maps capability to coarse actions', () => {
    expect(capabilityAction('domains:read')).toBe('read');
    expect(capabilityAction('domains:write')).toBe('write');
    expect(capabilityAction('domains:delete')).toBe('admin');
    expect(capabilityAction('deployments:trigger')).toBe('write');
  });

  it('resolves capability for resource type and action', () => {
    expect(capabilityFor('domain', 'read')).toBe('domains:read');
    expect(capabilityFor('domain', 'admin')).toBe('domains:delete');
    expect(capabilityFor('deployment', 'write')).toBe('deployments:trigger');
    expect(capabilityFor('unknown', 'read')).toBeNull();
  });

  it('expands coarse permissions into implied capabilities', () => {
    const readImplied = capabilitiesImpliedBy(['read']);
    expect(readImplied).toContain('domains:read');
    expect(readImplied).toContain('deployments:read');
    expect(readImplied).not.toContain('domains:write');

    const writeImplied = capabilitiesImpliedBy(['write']);
    expect(writeImplied).toContain('domains:read');
    expect(writeImplied).toContain('domains:write');
    expect(writeImplied).toContain('deployments:trigger');
    expect(writeImplied).not.toContain('domains:delete');

    const adminImplied = capabilitiesImpliedBy(['admin']);
    expect(adminImplied).toContain('domains:delete');
    expect(adminImplied.length).toBe(CAPABILITIES.length);
  });

  it('lists capabilities for project root', () => {
    const projectCaps = capabilitiesForRoot('project');
    expect(projectCaps).toEqual([...CAPABILITIES]);
  });
});
