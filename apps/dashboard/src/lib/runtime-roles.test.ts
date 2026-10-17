import { describe, expect, it } from 'vitest';
import {
  type RuntimeRole,
  roleLogsHref,
  roleLogTargetKey,
  runtimeRoleLogTargets,
  runtimeRoleState,
} from './runtime-roles';

describe('runtime-roles helpers', () => {
  it('computes correct role status based on replica states', () => {
    const runningRole: RuntimeRole = {
      name: 'worker',
      kind: 'worker',
      command: 'run.sh',
      replicas: 2,
      containers: [
        { containerId: 'c1', replica: 1, status: 'running', restartCount: 0 },
        { containerId: 'c2', replica: 2, status: 'running', restartCount: 1 },
      ],
    };
    expect(runtimeRoleState(runningRole)).toBe('running');

    const failingRole: RuntimeRole = {
      ...runningRole,
      containers: [
        { containerId: 'c1', replica: 1, status: 'running', restartCount: 0 },
        { containerId: 'c2', replica: 2, status: 'failed', restartCount: 5 },
      ],
    };
    expect(runtimeRoleState(failingRole)).toBe('failed');
  });

  it('generates role logs href and keys', () => {
    const href = roleLogsHref('proj-1', 'dep-1', 'worker', 1);
    expect(href).toBe('/projects/proj-1/logs?role=worker&replica=1&deploymentId=dep-1');

    const key = roleLogTargetKey('dep-1', 'worker', 2);
    expect(key).toBe('role:["dep-1","worker",2]');
  });

  it('extracts log targets from runtime roles', () => {
    const targets = runtimeRoleLogTargets({
      deploymentId: 'dep-100',
      roles: [
        {
          name: 'crawler',
          kind: 'worker',
          command: 'crawl.sh',
          replicas: 1,
          containers: [{ containerId: 'cid-1', replica: 1, status: 'running', restartCount: 0 }],
        },
      ],
    });
    expect(targets.length).toBe(1);
    expect(targets[0].role).toBe('crawler');
    expect(targets[0].deploymentId).toBe('dep-100');
  });
});
