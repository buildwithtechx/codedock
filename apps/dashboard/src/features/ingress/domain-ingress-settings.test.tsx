import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { ingressService } from '#/services/ingress';
import { DomainIngressSettings } from './domain-ingress-settings';

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

describe('DomainIngressSettings', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders manage tunnel button and status badge when active', async () => {
    vi.spyOn(ingressService, 'getDomain').mockResolvedValue({
      id: 'igh-1',
      domainId: 'dom-1',
      connectionId: 'ing-1',
      hostname: 'api.example.com',
      status: 'ready',
      phase: 'dns',
      error: null,
      updatedAt: '2026-10-16T12:00:00Z',
    });
    vi.spyOn(ingressService, 'list').mockResolvedValue([]);

    render(<DomainIngressSettings domainId="dom-1" hostname="api.example.com" />, {
      wrapper: createWrapper(),
    });

    expect(await screen.findByText('Tunnel: ready')).toBeDefined();
    expect(screen.getByRole('button', { name: /tunnel routing/i })).toBeDefined();
  });

  it('opens modal and toggles ingress routing', async () => {
    vi.spyOn(ingressService, 'getDomain').mockResolvedValue(null);
    vi.spyOn(ingressService, 'list').mockResolvedValue([
      {
        id: 'ing-1',
        serverId: 'srv-1',
        credentialId: 'cf-cred-1',
        accountId: 'cf-acc-1',
        name: 'Main Edge',
        connector: 'managed',
        edgePort: 80,
        managedTunnel: true,
        defaultForDomains: true,
        status: 'ready',
        createdAt: '2026-10-16T12:00:00Z',
        updatedAt: '2026-10-16T12:00:00Z',
      },
    ]);

    const setDomainSpy = vi.spyOn(ingressService, 'setDomain').mockResolvedValue({
      id: 'igh-1',
      domainId: 'dom-1',
      connectionId: 'ing-1',
      hostname: 'api.example.com',
      status: 'ready',
      phase: 'dns',
      error: null,
      updatedAt: '2026-10-16T12:00:00Z',
    });

    render(<DomainIngressSettings domainId="dom-1" serverId="srv-1" hostname="api.example.com" />, {
      wrapper: createWrapper(),
    });

    const button = screen.getByRole('button', { name: /tunnel routing/i });
    fireEvent.click(button);

    expect(await screen.findByText('Tunnel Ingress')).toBeDefined();
    expect(screen.getByText('api.example.com')).toBeDefined();
    expect(await screen.findByText('Main Edge')).toBeDefined();

    const switchBtn = await screen.findByRole('switch');
    await waitFor(() => expect(switchBtn.hasAttribute('disabled')).toBe(false));
    fireEvent.click(switchBtn);

    await waitFor(() => {
      expect(setDomainSpy).toHaveBeenCalledWith({
        domainId: 'dom-1',
        connectionId: 'ing-1',
      });
    });
  });
});
