import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { ComposeDeployForm } from './compose-deploy-form';

const api = vi.hoisted(() => ({ post: vi.fn() }));
vi.mock('#/hooks/use-environments', () => ({
  useListByProject: () => ({
    data: { data: [{ id: 'environment', name: 'Production', isDefault: true }] },
  }),
}));
vi.mock('#/lib/api-client', () => ({
  apiClient: { get: async () => ({ data: [] }), post: api.post },
}));

describe('Compose stack review', () => {
  it('rejects a review response after its source changed', async () => {
    let resolve: (value: unknown) => void = () => {};
    api.post.mockReturnValueOnce(
      new Promise((complete) => {
        resolve = complete;
      })
    );
    const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
    render(
      <QueryClientProvider client={client}>
        <ComposeDeployForm projectId="project" />
      </QueryClientProvider>
    );
    fireEvent.click(screen.getByRole('button', { name: 'Validate and review' }));
    await waitFor(() => expect(api.post).toHaveBeenCalledTimes(1));
    fireEvent.change(screen.getByLabelText('Compose configuration'), {
      target: { value: 'services: {}' },
    });
    await act(async () => {
      resolve({
        data: { config: 'stale configuration', digest: 'old', services: ['old-web'], effects: [] },
      });
    });
    expect(screen.queryByText('Deployment review')).toBeNull();
    expect(screen.queryByRole('button', { name: 'Save and deploy' })).toBeNull();
    expect(api.post).toHaveBeenCalledTimes(1);
  });
});
