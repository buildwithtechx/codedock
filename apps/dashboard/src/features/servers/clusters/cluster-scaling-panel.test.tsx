import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { clustersService } from '#/services/clusters';
import { ClusterScalingPanel } from './cluster-scaling-panel';

const createWrapper = () => {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: { retry: false },
    },
  });
  return ({ children }: { children: React.ReactNode }) => (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  );
};

describe('ClusterScalingPanel', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders not_enabled state with Enable scaling button', async () => {
    vi.spyOn(clustersService, 'getScaling').mockResolvedValue({
      status: 'not_enabled',
      k3sVersion: 'v1.31.1+k3s1',
      controlNodes: 1,
      workerNodes: 0,
      updatedAt: '2026-10-16T12:00:00Z',
    });

    render(<ClusterScalingPanel clusterId="cls-1" />, { wrapper: createWrapper() });

    expect(await screen.findByText('Cluster Scaling Runtime')).toBeDefined();
    expect(screen.getByText('not enabled')).toBeDefined();
    expect(screen.getByRole('button', { name: /enable scaling/i })).toBeDefined();
  });

  it('triggers setup when Enable scaling button is clicked', async () => {
    vi.spyOn(clustersService, 'getScaling').mockResolvedValue({
      status: 'not_enabled',
      k3sVersion: 'v1.31.1+k3s1',
      controlNodes: 1,
      workerNodes: 0,
      updatedAt: '2026-10-16T12:00:00Z',
    });

    const setupSpy = vi.spyOn(clustersService, 'setupScaling').mockResolvedValue({
      status: 'ready',
      k3sVersion: 'v1.31.1+k3s1',
      controlNodes: 1,
      workerNodes: 2,
      updatedAt: '2026-10-16T12:00:00Z',
    });

    const onReadyChange = vi.fn();
    render(<ClusterScalingPanel clusterId="cls-1" onReadyChange={onReadyChange} />, {
      wrapper: createWrapper(),
    });

    const enableBtn = await screen.findByRole('button', { name: /enable scaling/i });
    await waitFor(() => expect(enableBtn.hasAttribute('disabled')).toBe(false));
    fireEvent.click(enableBtn);

    await waitFor(() => {
      expect(setupSpy).toHaveBeenCalledWith('cls-1');
      expect(onReadyChange).toHaveBeenCalledWith(true);
    });
  });

  it('toggles technical details section', async () => {
    vi.spyOn(clustersService, 'getScaling').mockResolvedValue({
      status: 'ready',
      k3sVersion: 'v1.31.1+k3s1',
      controlNodes: 3,
      workerNodes: 4,
      podCidr: '10.42.0.0/16',
      serviceCidr: '10.43.0.0/16',
      clusterDomain: 'cluster.local',
      updatedAt: '2026-10-16T12:00:00Z',
    });

    render(<ClusterScalingPanel clusterId="cls-1" />, { wrapper: createWrapper() });

    const toggleBtn = await screen.findByRole('button', { name: /show technical details/i });
    fireEvent.click(toggleBtn);

    expect(screen.getByText(/Pod CIDR:/i)).toBeDefined();
    expect(screen.getByText('10.42.0.0/16')).toBeDefined();
    expect(screen.getByText('cluster.local')).toBeDefined();
  });
});
