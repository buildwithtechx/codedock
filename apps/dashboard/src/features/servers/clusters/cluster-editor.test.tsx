import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { fireEvent, render, screen } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { clustersService } from '#/services/clusters';
import { ClusterEditor } from './cluster-editor';

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

describe('ClusterEditor', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.spyOn(clustersService, 'getScaling').mockResolvedValue({
      status: 'ready',
      k3sVersion: 'v1.31.1+k3s1',
      controlNodes: 1,
      workerNodes: 2,
      updatedAt: '2026-10-16T12:00:00Z',
    });
  });

  it('renders initial form fields and nodes', () => {
    render(<ClusterEditor />, { wrapper: createWrapper() });

    expect(screen.getByText('K3s Server Cluster Pool')).toBeDefined();
    expect(screen.getByLabelText(/cluster name/i)).toBeDefined();
    expect(screen.getByRole('button', { name: /1 node \(standalone\)/i })).toBeDefined();
    expect(screen.getByRole('button', { name: /3 nodes \(ha quorum\)/i })).toBeDefined();
    expect(screen.getByText(/assigned server nodes \(3\)/i)).toBeDefined();
  });

  it('adds a node when Add Node is clicked', () => {
    render(<ClusterEditor />, { wrapper: createWrapper() });

    const addBtn = screen.getByRole('button', { name: /add node/i });
    fireEvent.click(addBtn);

    expect(screen.getByText(/assigned server nodes \(4\)/i)).toBeDefined();
  });

  it('submits saved cluster pool', () => {
    const onSaved = vi.fn();
    render(<ClusterEditor onSaved={onSaved} />, { wrapper: createWrapper() });

    const nameInput = screen.getByLabelText(/cluster name/i);
    fireEvent.change(nameInput, { target: { value: 'staging-k3s' } });

    const saveBtn = screen.getByRole('button', { name: /save cluster pool/i });
    fireEvent.click(saveBtn);

    expect(onSaved).toHaveBeenCalledWith(
      expect.objectContaining({
        name: 'staging-k3s',
        version: 'v1.31.1+k3s1',
      })
    );
  });
});
