import { act } from 'react';
import { createRoot, type Root } from 'react-dom/client';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { CommandEnvironmentNotice } from './command-environment-notice';

const mockPost = vi.fn();
vi.mock('#/lib/api-client', () => ({
  apiClient: {
    post: (...args: unknown[]) => mockPost(...args),
  },
}));

describe('CommandEnvironmentNotice', () => {
  let root: Root;
  let container: HTMLDivElement;

  beforeEach(() => {
    vi.stubGlobal('IS_REACT_ACT_ENVIRONMENT', true);
    mockPost.mockReset();
    container = document.createElement('div');
    document.body.append(container);
    root = createRoot(container);
  });

  afterEach(async () => {
    await act(async () => root.unmount());
    container.remove();
    vi.unstubAllGlobals();
  });

  it('renders variables without exposing values and invokes migrateDraft', async () => {
    const migrate = vi.fn();
    await act(async () => {
      root.render(
        <CommandEnvironmentNotice
          commands={{ buildCommand: "TOKEN='secret-value' bun run build" }}
          migrateDraft={migrate}
        />
      );
    });

    expect(container.textContent).toContain('TOKEN');
    expect(container.textContent).not.toContain('secret-value');

    const button = container.querySelector('button');
    expect(button).not.toBeNull();
    await act(async () => button!.click());

    expect(migrate).toHaveBeenCalledWith(
      { buildCommand: 'bun run build' },
      { TOKEN: 'secret-value' }
    );
  });

  it('does not allow migration if shell expansion is used', async () => {
    const migrate = vi.fn();
    await act(async () => {
      root.render(
        <CommandEnvironmentNotice
          commands={{ buildCommand: 'TOKEN=$SECRET build' }}
          migrateDraft={migrate}
        />
      );
    });
    expect(container.querySelector('button')).toBeNull();
    expect(container.textContent).toContain('Review this command manually');
  });

  it('calls api client when projectId is supplied', async () => {
    const onMigrated = vi.fn();
    mockPost.mockResolvedValueOnce({
      data: { commands: { startCommand: 'serve' } },
    });

    await act(async () => {
      root.render(
        <CommandEnvironmentNotice
          projectId="proj-123"
          serviceId="svc-1"
          commands={{ startCommand: 'PORT=8080 serve' }}
          onMigrated={onMigrated}
        />
      );
    });

    const button = container.querySelector('button');
    expect(button).not.toBeNull();
    await act(async () => button!.click());

    expect(mockPost).toHaveBeenCalledWith('/projects/proj-123/env/migrate-commands', {
      environment: 'production',
      serviceId: 'svc-1',
      commands: { startCommand: 'PORT=8080 serve' },
    });
    expect(onMigrated).toHaveBeenCalledWith({ startCommand: 'serve' });
  });
});
