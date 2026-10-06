import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, renderHook, waitFor } from '@testing-library/react';
import type { ReactNode } from 'react';
import { describe, expect, it, vi } from 'vitest';
import { useApplicationSetup } from './use-application-setup';

const mocks = vi.hoisted(() => ({ create: vi.fn(), variable: vi.fn(), inspect: vi.fn() }));
vi.mock('#/hooks/use-environments', () => ({
  useListByProject: () => ({ data: { data: [{ id: 'environment', isDefault: true }] } }),
}));
vi.mock('#/stores/organization-store', () => ({ useOrganizationStore: () => 'organization' }));
vi.mock('#/services/canvas', () => ({
  canvasService: { listCanvasSummaries: async () => ({ data: [] }) },
}));
vi.mock('#/lib/api-client', () => ({
  apiClient: { get: async () => ({ data: { serverId: '' } }), post: mocks.inspect },
}));
vi.mock('#/services/apps', () => ({ appsService: { createApp: mocks.create } }));
vi.mock('#/services/service-variables', () => ({
  serviceVariablesService: { create: mocks.variable },
}));
vi.mock('#/services/deployments', () => ({ deploymentsService: { trigger: vi.fn() } }));

function Wrapper({ children }: { children: ReactNode }) {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return <QueryClientProvider client={client}>{children}</QueryClientProvider>;
}

describe('application setup recovery', () => {
  it('reuses the created application and saved variables after a configuration failure', async () => {
    mocks.create.mockResolvedValue({ data: { id: 'created-app' } });
    mocks.variable
      .mockResolvedValueOnce({})
      .mockRejectedValueOnce(new Error('storage unavailable'))
      .mockResolvedValueOnce({});
    const onOpenChange = vi.fn();
    const { result } = renderHook(
      () =>
        useApplicationSetup({
          projectId: 'project',
          isOpen: true,
          onOpenChange,
          initialSource: 'image',
        }),
      { wrapper: Wrapper }
    );
    await waitFor(() => expect(result.current.project.data).toBeDefined());
    act(() => {
      result.current.update('name', 'Application');
      result.current.update('imageRef', 'nginx');
      result.current.setVariables('FIRST=one\nSECOND=two');
    });
    await act(async () => {
      await expect(result.current.apply.mutateAsync(false)).rejects.toThrow('storage unavailable');
    });
    await act(async () => {
      await result.current.apply.mutateAsync(false);
    });
    expect(mocks.create).toHaveBeenCalledTimes(1);
    expect(mocks.variable.mock.calls.map((call) => call[1].key)).toEqual([
      'FIRST',
      'SECOND',
      'SECOND',
    ]);
    expect(onOpenChange).toHaveBeenCalledWith(false);
  });

  it('keeps manual values when repository detection fails', async () => {
    mocks.inspect.mockRejectedValueOnce(new Error('repository unavailable'));
    const { result } = renderHook(
      () => useApplicationSetup({ projectId: 'project', isOpen: true, onOpenChange: vi.fn() }),
      { wrapper: Wrapper }
    );
    act(() => {
      result.current.update('repositoryUrl', 'https://github.com/example/app');
      result.current.update('buildCommand', 'custom build');
    });
    await act(async () => {
      await expect(result.current.inspect.mutateAsync()).rejects.toThrow();
    });
    expect(result.current.draft.buildCommand).toBe('custom build');
  });
});
