import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { networksService } from '#/services/networks';
import { NetworkRollbackWatchdog } from './network-rollback-watchdog';

const createTestQueryClient = () =>
  new QueryClient({
    defaultOptions: {
      queries: { retry: false },
      mutations: { retry: false },
    },
  });

function renderWithClient(ui: React.ReactElement) {
  const queryClient = createTestQueryClient();
  return render(<QueryClientProvider client={queryClient}>{ui}</QueryClientProvider>);
}

describe('NetworkRollbackWatchdog', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders armed watchdog with 120s timer countdown and action buttons', () => {
    renderWithClient(
      <NetworkRollbackWatchdog
        watchdog={{
          operationId: 'op-123',
          armed: true,
          timeoutSeconds: 120,
          remainingSeconds: 120,
          expiresAt: '2026-10-16T12:02:00Z',
          canCommit: true,
          canRollback: true,
          status: 'armed',
        }}
      />
    );

    expect(screen.getByText(/Rollback Watchdog Armed \(120s\)/i)).toBeDefined();
    expect(screen.getByText('120s remaining')).toBeDefined();
    expect(screen.getByRole('button', { name: /emergency rollback/i })).toBeDefined();
    expect(screen.getByRole('button', { name: /commit changes/i })).toBeDefined();
  });

  it('commits changes when Commit Changes button is clicked', async () => {
    const commitSpy = vi.spyOn(networksService, 'commitWatchdog').mockResolvedValue({
      operationId: 'op-123',
      armed: false,
      timeoutSeconds: 120,
      remainingSeconds: 0,
      expiresAt: '2026-10-16T12:02:00Z',
      canCommit: false,
      canRollback: false,
      status: 'committed',
    });
    const onCommitted = vi.fn();

    renderWithClient(
      <NetworkRollbackWatchdog
        watchdog={{
          operationId: 'op-123',
          armed: true,
          timeoutSeconds: 120,
          remainingSeconds: 110,
          expiresAt: '2026-10-16T12:02:00Z',
          canCommit: true,
          canRollback: true,
          status: 'armed',
        }}
        onCommitted={onCommitted}
      />
    );

    const commitBtn = screen.getByRole('button', { name: /commit changes/i });
    fireEvent.click(commitBtn);

    await waitFor(() => {
      expect(commitSpy).toHaveBeenCalledWith('op-123');
      expect(onCommitted).toHaveBeenCalled();
    });
  });

  it('rolls back changes when Emergency Rollback button is clicked', async () => {
    const rollbackSpy = vi.spyOn(networksService, 'rollbackWatchdog').mockResolvedValue({
      operationId: 'op-123',
      armed: false,
      timeoutSeconds: 120,
      remainingSeconds: 0,
      expiresAt: '2026-10-16T12:02:00Z',
      canCommit: false,
      canRollback: false,
      status: 'rolled_back',
    });
    const onRolledBack = vi.fn();

    renderWithClient(
      <NetworkRollbackWatchdog
        watchdog={{
          operationId: 'op-123',
          armed: true,
          timeoutSeconds: 120,
          remainingSeconds: 95,
          expiresAt: '2026-10-16T12:02:00Z',
          canCommit: true,
          canRollback: true,
          status: 'armed',
        }}
        onRolledBack={onRolledBack}
      />
    );

    const rollbackBtn = screen.getByRole('button', { name: /emergency rollback/i });
    fireEvent.click(rollbackBtn);

    await waitFor(() => {
      expect(rollbackSpy).toHaveBeenCalledWith('op-123');
      expect(onRolledBack).toHaveBeenCalled();
    });
  });
});
