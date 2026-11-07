import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { render, screen, fireEvent } from '@testing-library/react';
import type * as React from 'react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { TokensLayout } from './tokens-page';

const mockNavigate = vi.fn();
let currentSearch = { tab: 'private-keys' };

vi.mock('@tanstack/react-router', () => ({
  useNavigate: () => mockNavigate,
  useRouterState: ({
    select,
  }: {
    select: (s: { location: { search: typeof currentSearch } }) => unknown;
  }) => select({ location: { search: currentSearch } }),
  Link: ({ children, to, ...props }: { children: React.ReactNode; to: string }) => (
    <a href={to} {...props}>
      {children}
    </a>
  ),
}));

vi.mock('#/hooks/use-managed', () => ({
  useManagedCredentials: () => ({ data: [], isLoading: false }),
  useCreateManagedCredential: () => ({ mutate: vi.fn(), isPending: false }),
  useDeleteManagedCredential: () => ({ mutate: vi.fn(), isPending: false }),
}));

vi.mock('#/hooks/use-profile', () => ({
  useListTokens: () => ({ data: { data: [] }, isLoading: false }),
  useCreateToken: () => ({ mutate: vi.fn(), isPending: false }),
  useDeleteToken: () => ({ mutate: vi.fn(), isPending: false }),
}));

vi.mock('#/hooks/use-canvas', () => ({
  useListCanvasSummaries: () => ({ data: { data: [] } }),
}));

vi.mock('#/hooks/use-organizations', () => ({
  useListOrganizations: () => ({
    data: [{ id: 'org-test-1', name: 'Solomon Team' }],
    isLoading: false,
  }),
}));

vi.mock('#/stores/organization-store', () => ({
  useOrganizationStore: (selector: (s: { activeOrganizationId: string }) => unknown) =>
    selector({ activeOrganizationId: 'org-test-1' }),
}));

describe('TokensLayout (Right Sidebar Navigation)', () => {
  let queryClient: QueryClient;

  beforeEach(() => {
    queryClient = new QueryClient({
      defaultOptions: { queries: { retry: false } },
    });
    vi.clearAllMocks();
    currentSearch = { tab: 'private-keys' };
  });

  const renderComponent = () =>
    render(
      <QueryClientProvider client={queryClient}>
        <TokensLayout />
      </QueryClientProvider>
    );

  it('renders Keys & Tokens header and all five credential tabs', () => {
    renderComponent();

    expect(screen.getByText('Keys & Tokens')).toBeDefined();
    expect(screen.getAllByText('Private Keys').length).toBeGreaterThan(0);
    expect(screen.getAllByText('Cloud Tokens').length).toBeGreaterThan(0);
    expect(screen.getAllByText('Integration Tokens').length).toBeGreaterThan(0);
    expect(screen.getAllByText('Cloud-Init Scripts').length).toBeGreaterThan(0);
    expect(screen.getAllByText('API Tokens').length).toBeGreaterThan(0);
  });

  it('renders private keys tab content by default', () => {
    renderComponent();

    expect(screen.getByText('New private key')).toBeDefined();
    expect(screen.getByText('Default Deploy Key')).toBeDefined();
  });

  it('navigates when clicking a right rail tab', () => {
    renderComponent();

    const cloudTokensBtn = screen.getAllByRole('button', { name: /Cloud Tokens/i })[0];
    fireEvent.click(cloudTokensBtn);

    expect(mockNavigate).toHaveBeenCalledWith({
      to: '/tokens',
      search: { tab: 'cloud-tokens' },
      replace: true,
    });
  });
});
