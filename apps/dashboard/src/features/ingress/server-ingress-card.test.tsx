import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { ingressService } from '#/services/ingress';
import { ServerIngressCard } from './server-ingress-card';

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

describe('ServerIngressCard', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders empty state when no ingress tunnel exists', async () => {
    vi.spyOn(ingressService, 'list').mockResolvedValue([]);

    render(<ServerIngressCard serverId="srv-1" />, { wrapper: createWrapper() });

    expect(await screen.findByText('No ingress tunnel connected')).toBeDefined();
    expect(screen.getByText(/Connect Cloudflare Tunnel to expose services/i)).toBeDefined();
    expect(screen.getAllByRole('button', { name: /connect tunnel/i })).toHaveLength(2);
  });

  it('renders active tunnel information and preferences button', async () => {
    vi.spyOn(ingressService, 'list').mockResolvedValue([
      {
        id: 'ing-1',
        serverId: 'srv-1',
        credentialId: 'cf-cred-1',
        accountId: 'cf-acc-1',
        tunnelId: 'tun-1',
        name: 'Production Ingress',
        connector: 'managed',
        edgePort: 80,
        managedTunnel: true,
        defaultForDomains: true,
        defaultAccessPolicyId: null,
        status: 'ready',
        error: null,
        createdAt: '2026-10-16T12:00:00Z',
        updatedAt: '2026-10-16T12:00:00Z',
      },
    ]);

    render(<ServerIngressCard serverId="srv-1" />, { wrapper: createWrapper() });

    expect(await screen.findByText('Production Ingress')).toBeDefined();
    expect(screen.getByText('ready')).toBeDefined();
    expect(screen.getByText(/Automatic enrollment active/i)).toBeDefined();
    expect(screen.getByRole('button', { name: /preferences/i })).toBeDefined();
  });

  it('renders retry button when status is error', async () => {
    vi.spyOn(ingressService, 'list').mockResolvedValue([
      {
        id: 'ing-1',
        serverId: 'srv-1',
        credentialId: 'cf-cred-1',
        accountId: 'cf-acc-1',
        tunnelId: 'tun-1',
        name: 'Stalled Tunnel',
        connector: 'managed',
        edgePort: 80,
        managedTunnel: true,
        defaultForDomains: false,
        status: 'error',
        error: 'Connector handshake timeout',
        createdAt: '2026-10-16T12:00:00Z',
        updatedAt: '2026-10-16T12:00:00Z',
      },
    ]);
    const retrySpy = vi.spyOn(ingressService, 'retry').mockResolvedValue({
      id: 'ing-1',
      serverId: 'srv-1',
      credentialId: 'cf-cred-1',
      accountId: 'cf-acc-1',
      name: 'Stalled Tunnel',
      connector: 'managed',
      edgePort: 80,
      managedTunnel: true,
      defaultForDomains: false,
      status: 'ready',
      createdAt: '2026-10-16T12:00:00Z',
      updatedAt: '2026-10-16T12:00:00Z',
    });

    render(<ServerIngressCard serverId="srv-1" />, { wrapper: createWrapper() });

    expect(await screen.findByText('Stalled Tunnel')).toBeDefined();
    expect(screen.getByText('error')).toBeDefined();

    const retryBtn = screen.getByRole('button', { name: /retry/i });
    expect(retryBtn).toBeDefined();
    fireEvent.click(retryBtn);

    await waitFor(() => {
      expect(retrySpy).toHaveBeenCalledWith('ing-1');
    });
  });
});
