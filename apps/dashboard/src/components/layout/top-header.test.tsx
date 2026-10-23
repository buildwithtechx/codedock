import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { fireEvent, render, screen } from '@testing-library/react';
import type * as React from 'react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { TopHeader } from './top-header';

let mockPathname = '/';

vi.mock('@tanstack/react-router', () => ({
  Link: ({ children, to }: { children: React.ReactNode; to: string }) => (
    <a href={to}>{children}</a>
  ),
  useRouterState: () => mockPathname,
}));

vi.mock('#/hooks/use-projects', () => ({
  useGetProject: () => ({
    data: { data: { id: 'p1', name: 'Alpha Project' } },
  }),
}));

vi.mock('#/hooks/use-servers', () => ({
  useServer: () => ({
    data: { id: 's1', name: 'Fleet Master Node' },
  }),
}));

vi.mock('#/features/projects/environment-switcher', () => ({
  EnvironmentSwitcher: () => <div data-testid="env-switcher">Environment Switcher</div>,
}));

const { mockUser } = vi.hoisted(() => ({
  mockUser: { id: 'u1', name: 'Solomon Olatunji', email: 'admin@codedock.local' },
}));

vi.mock('#/stores/auth-store', () => {
  const store = (selector: (s: { user: typeof mockUser }) => unknown) =>
    selector({ user: mockUser });
  store.getState = () => ({ user: mockUser });
  store.subscribe = vi.fn(() => () => {});
  return { useAuthStore: store };
});

function renderHeader(ui: React.ReactElement) {
  const qc = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return render(<QueryClientProvider client={qc}>{ui}</QueryClientProvider>);
}

describe('TopHeader', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockPathname = '/';
  });

  it('renders team breadcrumb and search trigger', () => {
    const handleCommand = vi.fn();
    renderHeader(
      <TopHeader mobileOpen={false} onToggleMobile={vi.fn()} onOpenCommand={handleCommand} />
    );

    expect(screen.getByText('Solomon Olatunji')).toBeTruthy();
    expect(screen.getByText('Dashboard')).toBeTruthy();
    expect(screen.getByText(/Search/i)).toBeTruthy();

    const searchBtn = screen.getByRole('button', { name: /Search/i });
    fireEvent.click(searchBtn);
    expect(handleCommand).toHaveBeenCalledOnce();
  });

  it('renders static section breadcrumb for /servers', () => {
    mockPathname = '/servers';
    renderHeader(<TopHeader mobileOpen={false} onToggleMobile={vi.fn()} onOpenCommand={vi.fn()} />);

    expect(screen.getByText('Servers')).toBeTruthy();
  });

  it('renders project breadcrumbs and inline environment switcher for /projects/:id', () => {
    mockPathname = '/projects/p1';
    renderHeader(<TopHeader mobileOpen={false} onToggleMobile={vi.fn()} onOpenCommand={vi.fn()} />);

    expect(screen.getByText('Projects')).toBeTruthy();
    expect(screen.getByText('Alpha Project')).toBeTruthy();
    expect(screen.getByTestId('env-switcher')).toBeTruthy();
  });

  it('toggles mobile menu on mobile button click', () => {
    const handleToggle = vi.fn();
    renderHeader(
      <TopHeader mobileOpen={false} onToggleMobile={handleToggle} onOpenCommand={vi.fn()} />
    );

    const toggleBtn = screen.getByRole('button', { name: /Open menu/i });
    fireEvent.click(toggleBtn);
    expect(handleToggle).toHaveBeenCalledOnce();
  });
});
