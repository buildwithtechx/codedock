import { act } from 'react';
import { createRoot, type Root } from 'react-dom/client';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import type { RuntimeRole } from '#/interfaces/roles';
import { RuntimeRolesCard } from './runtime-roles';

describe('RuntimeRolesCard', () => {
  let root: Root;
  let container: HTMLDivElement;

  beforeEach(() => {
    vi.stubGlobal('IS_REACT_ACT_ENVIRONMENT', true);
    container = document.createElement('div');
    document.body.append(container);
    root = createRoot(container);
  });

  afterEach(async () => {
    await act(async () => root.unmount());
    container.remove();
    vi.unstubAllGlobals();
  });

  it('renders nothing when roles list is empty', async () => {
    await act(async () => {
      root.render(<RuntimeRolesCard projectId="p1" deploymentId="d1" roles={[]} />);
    });
    expect(container.innerHTML).toBe('');
  });

  it('renders role containers with replica status and logs button', async () => {
    const roles: RuntimeRole[] = [
      {
        name: 'worker-bg',
        kind: 'worker',
        command: 'python run.py',
        replicas: 1,
        containers: [
          {
            containerId: 'cont-1',
            replica: 1,
            status: 'running',
            restartCount: 2,
          },
        ],
      },
    ];

    await act(async () => {
      root.render(<RuntimeRolesCard projectId="p1" deploymentId="d1" roles={roles} />);
    });

    expect(container.textContent).toContain('worker-bg');
    expect(container.textContent).toContain('Replica #1');
    expect(container.textContent).toContain('running');
    expect(container.textContent).toContain('2 restarts');
    expect(container.querySelector('button')).not.toBeNull();
  });
});
