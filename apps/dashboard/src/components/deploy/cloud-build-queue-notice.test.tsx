import { describe, expect, it, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import {
  CloudBuildQueueNotice,
  type CloudBuildQueueStatus,
} from './cloud-build-queue-notice';

describe('CloudBuildQueueNotice', () => {
  const baseQueue: CloudBuildQueueStatus = {
    position: 2,
    reason: 'concurrency',
    simultaneousBuilds: 1,
    activeBuilds: 1,
    workspaceId: 'ws-123',
  };

  it('renders queue position, simultaneous allowance, and compare plans button', () => {
    const handleCompare = vi.fn();
    render(<CloudBuildQueueNotice queue={baseQueue} onComparePlans={handleCompare} />);

    expect(screen.getByText('Position 2 in queue')).toBeTruthy();
    expect(
      screen.getByText(/Simultaneous builds allowed by your plan: 1./)
    ).toBeTruthy();
    expect(screen.getByText(/Active builds: 1./)).toBeTruthy();
    expect(screen.getByText('Upgrading is optional')).toBeTruthy();

    const button = screen.getByRole('button', { name: 'Compare plans' });
    expect(button).toBeTruthy();
    fireEvent.click(button);
    expect(handleCompare).toHaveBeenCalledWith('ws-123', 'concurrency');
  });

  it('offers optional plan upgrade for resource waits even with unlimited concurrency', () => {
    const handleCompare = vi.fn();
    render(
      <CloudBuildQueueNotice
        queue={{ ...baseQueue, reason: 'resources', simultaneousBuilds: null }}
        onComparePlans={handleCompare}
      />
    );

    expect(
      screen.getByText(/Waiting for available build resources/)
    ).toBeTruthy();
    expect(screen.getByText(/No simultaneous build limit./)).toBeTruthy();

    const button = screen.getByRole('button', { name: 'Compare plans' });
    expect(button).toBeTruthy();
    fireEvent.click(button);
    expect(handleCompare).toHaveBeenCalledWith('ws-123', 'resources');
  });

  it('hides plan upgrade action for server waits', () => {
    render(
      <CloudBuildQueueNotice
        queue={{ ...baseQueue, reason: 'server', position: null, simultaneousBuilds: 1 }}
      />
    );

    expect(screen.getByText('Build queued')).toBeTruthy();
    expect(screen.getByText(/Waiting for build server/)).toBeTruthy();
    expect(screen.queryByRole('button', { name: 'Compare plans' })).toBeNull();
  });

  it('hides plan upgrade action for concurrency waits when allowance is already unlimited', () => {
    render(
      <CloudBuildQueueNotice
        queue={{ ...baseQueue, reason: 'concurrency', simultaneousBuilds: null }}
      />
    );

    expect(screen.getByText(/No simultaneous build limit./)).toBeTruthy();
    expect(screen.queryByRole('button', { name: 'Compare plans' })).toBeNull();
  });
});
