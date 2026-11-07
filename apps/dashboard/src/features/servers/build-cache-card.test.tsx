import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import type { Server } from '#/interfaces/server';
import { serverService } from '#/services/servers';
import { BuildCacheCard } from './build-cache-card';

const GIB = 1024 ** 3;

const mockServer: Server = {
  id: 'srv-build-test',
  userId: 'usr-1',
  name: 'Build Worker 1',
  ipAddress: '192.0.2.1',
  status: 'online',
  capabilities: {
    buildCacheBudget: true,
  },
  metrics: null,
  createdAt: '2026-10-01T00:00:00Z',
  updatedAt: '2026-10-01T00:00:00Z',
};

const createWrapper = () => {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: { retry: false },
      mutations: { retry: false },
    },
  });
  return ({ children }: { children: React.ReactNode }) => (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  );
};

describe('BuildCacheCard', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders automatic mode by default when budget is null', () => {
    render(<BuildCacheCard server={{ ...mockServer, buildCacheBudgetBytes: null }} />, {
      wrapper: createWrapper(),
    });

    expect(screen.getByText('Build Cache Budget')).toBeDefined();
    expect(screen.getAllByText('Automatic (20% of disk)').length).toBeGreaterThanOrEqual(1);
  });

  it('renders custom mode with size when budget is configured', () => {
    render(<BuildCacheCard server={{ ...mockServer, buildCacheBudgetBytes: 24 * GIB }} />, {
      wrapper: createWrapper(),
    });

    expect(screen.getByText('Custom (24 GiB)')).toBeDefined();
    const input = screen.getByRole('textbox', { name: /budget in gib/i }) as HTMLInputElement;
    expect(input.value).toBe('24');
  });

  it('submits updated custom budget when form is submitted', async () => {
    const updateSpy = vi.spyOn(serverService, 'update').mockResolvedValue({
      ...mockServer,
      buildCacheBudgetBytes: 32 * GIB,
    });
    const onChanged = vi.fn();

    render(
      <BuildCacheCard
        server={{ ...mockServer, buildCacheBudgetBytes: 24 * GIB }}
        onChanged={onChanged}
      />,
      { wrapper: createWrapper() }
    );

    const input = screen.getByRole('textbox', { name: /budget in gib/i });
    fireEvent.change(input, { target: { value: '32' } });

    const saveBtn = screen.getByRole('button', { name: /save budget/i }) as HTMLButtonElement;
    expect(saveBtn.disabled).toBe(false);
    fireEvent.click(saveBtn);

    await waitFor(() => {
      expect(updateSpy).toHaveBeenCalledWith('srv-build-test', {
        buildCacheBudgetBytes: 32 * GIB,
      });
      expect(onChanged).toHaveBeenCalled();
    });
  });

  it('shows validation error for invalid budget input', () => {
    render(<BuildCacheCard server={{ ...mockServer, buildCacheBudgetBytes: 24 * GIB }} />, {
      wrapper: createWrapper(),
    });

    const input = screen.getByRole('textbox', { name: /budget in gib/i });
    fireEvent.change(input, { target: { value: '0.2' } }); // Below 1 GiB

    expect(screen.getByRole('alert')).toBeDefined();
    expect(screen.getByText(/budget must be between 1 gib and 65536 gib/i)).toBeDefined();

    const saveBtn = screen.getByRole('button', { name: /save budget/i }) as HTMLButtonElement;
    expect(saveBtn.disabled).toBe(true);
  });
});
