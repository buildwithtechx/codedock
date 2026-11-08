import type { ChatAppConnection } from '@codedock/contracts';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { fireEvent, render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { ChatAppsSettings } from './chat-apps-settings';

vi.mock('@tanstack/react-router', () => ({
  useSearch: vi.fn(() => ({})),
}));

const mockConnections: ChatAppConnection[] = [
  {
    id: 'chat_slack_1',
    userId: 'user_1',
    provider: 'slack',
    transport: 'outbound',
    purpose: 'both',
    name: 'Production Slack Bot',
    botUsername: '@codedock',
    config: {},
    enabled: true,
    status: 'connected',
    createdAt: '2026-11-08T00:00:00Z',
    updatedAt: '2026-11-08T00:00:00Z',
  },
];

vi.mock('#/hooks/use-chat-apps', () => ({
  useChatConnections: vi.fn(() => ({
    data: mockConnections,
    isLoading: false,
  })),
  useChatIdentities: vi.fn(() => ({
    data: [],
    isLoading: false,
  })),
  usePairingChallenge: vi.fn(() => ({
    data: null,
    isLoading: false,
  })),
  useCreateChatConnection: vi.fn(() => ({
    mutateAsync: vi.fn(),
    isPending: false,
  })),
  useUpdateChatConnection: vi.fn(() => ({
    mutate: vi.fn(),
    mutateAsync: vi.fn(),
    isPending: false,
  })),
  useDeleteChatConnection: vi.fn(() => ({
    mutate: vi.fn(),
    isPending: false,
  })),
  useTestChatConnection: vi.fn(() => ({
    mutate: vi.fn(),
    isPending: false,
  })),
  useDeleteChatIdentity: vi.fn(() => ({
    mutate: vi.fn(),
    isPending: false,
  })),
  useVerifyPairing: vi.fn(() => ({
    mutateAsync: vi.fn(),
    isPending: false,
  })),
}));

function renderWithClient(ui: React.ReactElement) {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: { retry: false },
    },
  });
  return render(<QueryClientProvider client={queryClient}>{ui}</QueryClientProvider>);
}

describe('ChatAppsSettings', () => {
  it('renders all 4 provider cards and status indicators', () => {
    renderWithClient(<ChatAppsSettings />);

    expect(screen.getByText('Inbound Chat Apps')).toBeDefined();
    expect(screen.getByText('Slack')).toBeDefined();
    expect(screen.getByText('Discord')).toBeDefined();
    expect(screen.getByText('Telegram')).toBeDefined();
    expect(screen.getByText('Microsoft Teams')).toBeDefined();

    // Slack is connected
    expect(screen.getByText('Connected')).toBeDefined();
    expect(screen.getByText('Production Slack Bot')).toBeDefined();

    // Other providers are not configured
    const notConfiguredList = screen.getAllByText('Not Configured');
    expect(notConfiguredList.length).toBe(3);
  });

  it('renders the chat command reference cheat sheet', () => {
    renderWithClient(<ChatAppsSettings />);

    expect(screen.getByText('Chat Command Reference')).toBeDefined();
    expect(screen.getByText('/codedock link')).toBeDefined();
    expect(screen.getByText('/codedock projects')).toBeDefined();
    expect(screen.getByText('/codedock redeploy <id>')).toBeDefined();
    expect(screen.getByText('/codedock confirm <token>')).toBeDefined();
  });

  it('opens configuration dialog when clicking Connect on an unconfigured provider', () => {
    renderWithClient(<ChatAppsSettings />);

    const connectDiscordBtn = screen.getByText('Connect Discord');
    fireEvent.click(connectDiscordBtn);

    expect(screen.getByText('Connect discord')).toBeDefined();
    expect(screen.getByText('Public Key (Ed25519)')).toBeDefined();
  });
});
