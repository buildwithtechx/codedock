import { fireEvent, render, screen } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { SettingsLayout } from './settings-page';

const mockNavigate = vi.fn();
let currentSearch = { tab: 'general' };

vi.mock('@tanstack/react-router', () => ({
  Link: ({ children, to, search }: { children: React.ReactNode; to: string; search?: unknown }) => (
    <a href={`${to}?${JSON.stringify(search)}`}>{children}</a>
  ),
  useNavigate: () => mockNavigate,
  useRouterState: ({ select }: { select: (s: { location: { search: unknown } }) => unknown }) =>
    select({ location: { search: currentSearch } }),
}));

vi.mock('#/stores/auth-store', () => ({
  useAuthStore: (selector: (s: { user: { email: string } }) => unknown) =>
    selector({ user: { email: 'admin@codedock.local' } }),
}));

vi.mock('./profile-settings-tab', () => ({
  ProfileSettingsTab: () => <div data-testid="profile-content">Profile Settings Tab</div>,
}));

vi.mock('./general-settings', () => ({
  GeneralSettings: () => <div data-testid="general-content">General Settings Tab</div>,
}));

vi.mock('./appearance-editor/appearance-editor', () => ({
  AppearanceEditor: () => <div data-testid="appearance-editor-content">Appearance Editor</div>,
}));

vi.mock('./server-settings-tab', () => ({
  ServerSettingsTab: () => <div data-testid="server-content">Server & Routing Tab</div>,
}));

vi.mock('#/features/dns/dns-settings', () => ({
  DnsSettings: () => <div data-testid="dns-content">DNS Providers Tab</div>,
}));

vi.mock('./ai/ai-settings', () => ({
  AISettings: () => <div data-testid="ai-content">AI Hub Tab</div>,
}));

vi.mock('./mcp/mcp-settings', () => ({
  McpSettings: () => <div data-testid="mcp-content">MCP Tab</div>,
}));

vi.mock('./team-settings', () => ({
  TeamSettings: () => <div data-testid="team-content">Team Tab</div>,
}));

vi.mock('#/features/notifications/notifications-settings', () => ({
  NotificationsSettings: () => <div data-testid="notif-content">Notifications Tab</div>,
}));

vi.mock('#/features/sources', () => ({
  GithubIntegration: () => <div>Github</div>,
  GitProviders: () => <div>GitProviders</div>,
}));

vi.mock('#/features/users/api-keys-list', () => ({
  ApiKeysList: () => <div>Tokens</div>,
}));

vi.mock('#/features/users/oauth-providers-list', () => ({
  OAuthProvidersList: () => <div>OAuth</div>,
}));

vi.mock('./maintenance-settings', () => ({
  MaintenancePage: () => <div>Maintenance</div>,
}));

vi.mock('./instance-info', () => ({
  InstanceInfo: () => <div>InstanceInfo</div>,
}));

vi.mock('./update-settings', () => ({
  UpdatesPage: () => <div>Updates</div>,
}));

vi.mock('./migration-settings', () => ({
  MigrationSettings: () => <div>Migrations</div>,
}));

describe('SettingsLayout Modernization', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    currentSearch = { tab: 'general' };
  });

  it('renders categorized sidebar groups', () => {
    render(<SettingsLayout />);

    expect(screen.getAllByText('Workspace & Preferences')[0]).toBeDefined();
    expect(screen.getAllByText('Integrations & Access')[0]).toBeDefined();
    expect(screen.getAllByText('AI & Automation')[0]).toBeDefined();
    expect(screen.getAllByText('Instance & Fleet')[0]).toBeDefined();
  });

  it('renders modernized tabs including DNS Providers and Server & Routing', () => {
    render(<SettingsLayout />);

    expect(screen.getAllByText('DNS Providers')[0]).toBeDefined();
    expect(screen.getAllByText('Server & Routing')[0]).toBeDefined();
    expect(screen.getAllByText('AI Models')[0]).toBeDefined();
    expect(screen.getAllByText('MCP Servers')[0]).toBeDefined();
  });

  it('navigates when a tab button is clicked', () => {
    render(<SettingsLayout />);

    const dnsButtons = screen.getAllByText('DNS Providers');
    fireEvent.click(dnsButtons[0]);

    expect(mockNavigate).toHaveBeenCalledWith({
      to: '/settings',
      search: { tab: 'dns' },
      replace: true,
    });
  });

  it('maps legacy credentials tab query parameter to dns tab seamlessly', () => {
    currentSearch = { tab: 'credentials' };
    render(<SettingsLayout />);

    expect(screen.getByTestId('dns-content')).toBeDefined();
  });

  it('renders AppearanceEditor when view is appearance', () => {
    currentSearch = { tab: 'general', view: 'appearance' } as never;
    render(<SettingsLayout />);

    expect(screen.getByTestId('appearance-editor-content')).toBeDefined();
  });
});
