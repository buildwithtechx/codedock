import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { VolumeRestoreDialog } from './volume-restore-dialog';

const mocks = vi.hoisted(() => ({ get: vi.fn(), post: vi.fn() }));
vi.mock('#/lib/api-client', () => ({ apiClient: mocks }));
describe('reviewed restore', () => {
  it('requires an expiring confirmation and observes completion', async () => {
    const operation = {
      id: 'operation',
      status: 'REVIEWED',
      phase: '',
      effects: 'Overwrite only reviewed volume',
      error: '',
      logs: '',
      expiresAt: Math.floor(Date.now() / 1000) + 600,
    };
    mocks.get.mockImplementation(() => Promise.resolve({ data: operation }));
    mocks.post.mockImplementation((url) =>
      Promise.resolve(
        url.endsWith('/review') ? { data: { operation, confirmation: 'opaque-confirmation' } } : {}
      )
    );
    const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
    render(
      <QueryClientProvider client={client}>
        <VolumeRestoreDialog recordId="record" onClose={vi.fn()} />
      </QueryClientProvider>
    );
    fireEvent.click(screen.getByRole('button', { name: 'Prepare restore' }));
    await screen.findByText('Overwrite only reviewed volume');
    const apply = screen.getByRole('button', { name: 'Apply restore' }) as HTMLButtonElement;
    expect(apply.disabled).toBe(true);
    fireEvent.click(screen.getByRole('checkbox'));
    fireEvent.click(apply);
    await waitFor(() =>
      expect(mocks.post).toHaveBeenCalledWith('/backup-operations/operation/apply', {
        confirmation: 'opaque-confirmation',
      })
    );
    client.clear();
  });
});
