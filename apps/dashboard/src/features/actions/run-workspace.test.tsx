import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { actionsService, type ActionRun } from '#/services/actions';
import { RunWorkspace } from './run-workspace';

const mockRun: ActionRun = {
  id: 'run-101',
  workflowId: 'wf-1',
  name: 'CI Build & Test',
  number: 12,
  status: 'running',
  triggeredBy: 'octocat',
  runnerId: 'runner-ubuntu-1',
  logs: '[Build] Compiling source\n[Test] 139 passed',
  jobs: [
    {
      id: 'job-1',
      name: 'Build and Verify',
      status: 'running',
      runnerId: 'runner-ubuntu-1',
      durationSeconds: 15,
      steps: [
        { id: 's-1', name: 'Checkout repo', status: 'ready', durationSeconds: 2 },
        { id: 's-2', name: 'Run tests', status: 'running', durationSeconds: 13 },
      ],
    },
  ],
  startedAt: '2026-10-16T12:00:00Z',
  durationSeconds: 15,
  createdAt: '2026-10-16T12:00:00Z',
  updatedAt: '2026-10-16T12:00:15Z',
};

describe('RunWorkspace', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders run header, metadata, jobs, and terminal log', () => {
    render(<RunWorkspace run={mockRun} />);

    expect(screen.getByText(/CI Build & Test/i)).toBeDefined();
    expect(screen.getByText('#12')).toBeDefined();
    expect(screen.getByText('octocat')).toBeDefined();
    expect(screen.getByText('15s')).toBeDefined();
    expect(screen.getByText('runner-ubuntu-1')).toBeDefined();
    expect(screen.getByText('Build and Verify')).toBeDefined();
    expect(screen.getByText('Checkout repo')).toBeDefined();
    expect(screen.getByText(/Compiling source/i)).toBeDefined();
  });

  it('allows cancelling an active run', async () => {
    const updatedRun: ActionRun = {
      ...mockRun,
      status: 'cancelled',
      logs: mockRun.logs + '\n[Abort] Run cancelled by operator',
    };
    vi.spyOn(actionsService, 'cancelRun').mockResolvedValue(updatedRun);

    render(<RunWorkspace run={mockRun} />);

    const cancelBtn = screen.getByRole('button', { name: /cancel run/i });
    fireEvent.click(cancelBtn);

    await waitFor(() => {
      expect(actionsService.cancelRun).toHaveBeenCalledWith('run-101');
      expect(screen.getByText('cancelled')).toBeDefined();
    });
  });

  it('allows retrying a run', async () => {
    const retriedRun: ActionRun = {
      ...mockRun,
      id: 'run-102',
      number: 13,
      status: 'running',
      triggeredBy: 'retry',
    };
    vi.spyOn(actionsService, 'retryRun').mockResolvedValue(retriedRun);

    render(<RunWorkspace run={mockRun} />);

    const retryBtn = screen.getByRole('button', { name: /re-run/i });
    fireEvent.click(retryBtn);

    await waitFor(() => {
      expect(actionsService.retryRun).toHaveBeenCalledWith('run-101');
      expect(screen.getByText('#13')).toBeDefined();
    });
  });
});
