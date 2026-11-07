import { fireEvent, render, screen } from '@testing-library/react';
import type * as React from 'react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { AppSidebar } from './app-sidebar';

vi.mock('@tanstack/react-router', () => ({
  Link: ({ children, to }: { children: React.ReactNode; to: string }) => (
    <a href={to}>{children}</a>
  ),
  useRouterState: () => ({
    location: { pathname: '/', search: {} },
  }),
}));

vi.mock('#/stores/auth-store', () => ({
  useAuthStore: (selector: (s: { user: { email: string; name: string } }) => unknown) =>
    selector({ user: { email: 'daniel@takentrade.com', name: 'Daniel Uko' } }),
}));

vi.mock('#/stores/dock-store', () => ({
  useDockStore: () => vi.fn(),
}));

vi.mock('./use-sidebar-counts', () => ({
  useSidebarCounts: () => ({
    projects: 5,
    servers: 2,
    runningApps: 3,
    serverStatus: 'healthy',
    hasIssues: false,
  }),
}));

describe('AppSidebar (Coolify & Dokploy Developer Design)', () => {
  const onToggle = vi.fn();
  const onMobileClose = vi.fn();
  const onOpenCommand = vi.fn();

  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders brand header, version badge, and quick search trigger', () => {
    render(
      <AppSidebar
        collapsed={false}
        onToggle={onToggle}
        mobileOpen={false}
        onMobileClose={onMobileClose}
        onOpenCommand={onOpenCommand}
      />
    );

    expect(screen.getByText('Codedock')).toBeDefined();
    expect(screen.getByText('v0.1.0')).toBeDefined();
    expect(screen.getByText('Search')).toBeDefined();
    expect(screen.getByText('Ctrl+K')).toBeDefined();
  });

  it('renders Workspace, Infrastructure, and Manage navigation groups', () => {
    render(
      <AppSidebar
        collapsed={false}
        onToggle={onToggle}
        mobileOpen={false}
        onMobileClose={onMobileClose}
        onOpenCommand={onOpenCommand}
      />
    );

    // Workspace
    expect(screen.getByText('Workspace')).toBeDefined();
    expect(screen.getByText('Dashboard')).toBeDefined();
    expect(screen.getByText('Projects')).toBeDefined();
    expect(screen.getByText('Apps')).toBeDefined();
    expect(screen.getByText('Deployments')).toBeDefined();
    expect(screen.getByText('Actions')).toBeDefined();
    expect(screen.getByText('Analytics')).toBeDefined();
    expect(screen.queryByText('Terminal')).toBeNull();

    // Infrastructure
    expect(screen.getByText('Infrastructure')).toBeDefined();
    expect(screen.getByText('Servers')).toBeDefined();
    expect(screen.getByText('Domains')).toBeDefined();
    expect(screen.getByText('Emails')).toBeDefined();
    expect(screen.getByText('Backups')).toBeDefined();
    expect(screen.getByText('Cron Jobs')).toBeDefined();

    // Manage
    expect(screen.getByText('Manage')).toBeDefined();
    expect(screen.getByText('Team')).toBeDefined();
    expect(screen.getByText('Keys & Tokens')).toBeDefined();
    expect(screen.queryByText('Audit Logs')).toBeNull();
    expect(screen.queryByText('Terminal')).toBeNull();
    expect(screen.getByText('Settings')).toBeDefined();
  });

  it('renders expand button when collapsed', () => {
    render(
      <AppSidebar
        collapsed={true}
        onToggle={onToggle}
        mobileOpen={false}
        onMobileClose={onMobileClose}
        onOpenCommand={onOpenCommand}
      />
    );

    const expandBtn = screen.getByRole('button', { name: /Expand sidebar/i });
    expect(expandBtn).toBeDefined();
    fireEvent.click(expandBtn);
    expect(onToggle).toHaveBeenCalledOnce();
  });

  it('renders user profile row with name and email', () => {
    render(
      <AppSidebar
        collapsed={false}
        onToggle={onToggle}
        mobileOpen={false}
        onMobileClose={onMobileClose}
        onOpenCommand={onOpenCommand}
      />
    );

    expect(screen.getByText('Daniel Uko')).toBeDefined();
    expect(screen.getByText('daniel@takentrade.com')).toBeDefined();
  });

  it('clicking search trigger triggers onOpenCommand', () => {
    render(
      <AppSidebar
        collapsed={false}
        onToggle={onToggle}
        mobileOpen={false}
        onMobileClose={onMobileClose}
        onOpenCommand={onOpenCommand}
      />
    );

    const searchBtn = screen.getByText('Search');
    fireEvent.click(searchBtn);
    expect(onOpenCommand).toHaveBeenCalled();
  });
});
