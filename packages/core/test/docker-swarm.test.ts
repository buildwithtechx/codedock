import { describe, expect, it } from 'vitest';
import { assertStandaloneContainer, swarmTaskOwner } from '../src/index';

describe('docker-swarm utilities', () => {
  it('detects swarm service and task labels', () => {
    expect(swarmTaskOwner(null)).toBeNull();
    expect(swarmTaskOwner({})).toBeNull();
    expect(swarmTaskOwner({ 'com.docker.swarm.service.name': 'production_web' })).toBe(
      'production_web'
    );
    expect(swarmTaskOwner({ 'com.docker.swarm.task.id': 'task-12345' })).toBe('task-12345');
    expect(swarmTaskOwner({ 'com.docker.stack.namespace': 'my-stack' })).toBe('my-stack');
  });

  it('permits standalone container assertion', () => {
    expect(() => assertStandaloneContainer({})).not.toThrow();
    expect(() => assertStandaloneContainer({ 'codedock.service_id': 'srv-1' })).not.toThrow();
  });

  it('rejects swarm task container adoption as standalone', () => {
    expect(() => assertStandaloneContainer({ 'com.docker.swarm.service.name': 'db' })).toThrowError(
      /belongs to Docker Swarm \(db\)/
    );
  });
});
