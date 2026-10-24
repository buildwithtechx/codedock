import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { render, screen, fireEvent } from '@testing-library/react';
import type * as React from 'react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { TeamLayout } from './team-page';

const mockNavigate = vi.fn();
let currentSearch = { tab: 'general' };

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

vi.mock('#/hooks/use-organizations', () => ({
  useListOrganizations: () => ({
    data: [{ id: 'org-test-1', name: 'Solomon Olatunji Team', createdAt: '2026-01-01T00:00:00Z' }],
    isLoading: false,
  }),
  useUpdateOrganization: () => ({
    mutate: vi.fn(),
    isPending: false,
  }),
  useDeleteOrganization: () => ({
    mutate: vi.fn(),
    isPending: false,
  }),
}));

vi.mock('#/features/settings/team-settings', () => ({
  TeamSettings: () => <div data-testid="team-settings-members">Members Content</div>,
}));

vi.mock('#/features/audit/audit-log-list', () => ({
  AuditLogList: () => <div data-testid="audit-log-list">Audit Logs Content</div>,
}));

vi.mock('#/stores/organization-store', () => ({
  useOrganizationStore: (selector: (s: { activeOrganizationId: string }) => unknown) =>
    selector({ activeOrganizationId: 'org-test-1' }),
}));

describe('TeamLayout (Settings-style with Right Sidebar)', () => {
  let queryClient: QueryClient;

  beforeEach(() => {
    queryClient = new QueryClient({
      defaultOptions: { queries: { retry: false } },
    });
    vi.clearAllMocks();
    currentSearch = { tab: 'general' };
  });

  const renderComponent = () =>
    render(
      <QueryClientProvider client={queryClient}>
        <TeamLayout />
      </QueryClientProvider>
    );

  it('renders Team & Workspace header and right rail tabs', () => {
    renderComponent();

    expect(screen.getByText('Team & Workspace')).toBeDefined();
    expect(screen.getAllByText('Solomon Olatunji Team').length).toBeGreaterThan(0);
    expect(screen.getAllByText('General').length).toBeGreaterThan(0);
    expect(screen.getAllByText('Members').length).toBeGreaterThan(0);
    expect(screen.getAllByText('Audit Logs').length).toBeGreaterThan(0);
    expect(screen.getAllByText('Danger Zone').length).toBeGreaterThan(0);
  });

  it('renders general tab by default with workspace properties', () => {
    renderComponent();

    expect(screen.getByLabelText('Workspace Name')).toBeDefined();
    expect(screen.getByText('Workspace ID')).toBeDefined();
    expect(screen.getByText('org-test-1')).toBeDefined();
  });

  it('navigates when a tab button is clicked', () => {
    renderComponent();

    const membersBtn = screen.getAllByRole('button', { name: /Members/i })[0];
    fireEvent.click(membersBtn);

    expect(mockNavigate).toHaveBeenCalledWith({
      to: '/team',
      search: { tab: 'members' },
      replace: true,
    });
  });
});
