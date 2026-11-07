import type { Fail2banStatus } from '@codedock/contracts';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { mailService } from '#/services/mail';
import { Fail2banDrawer } from './fail2ban-drawer';

vi.mock('#/services/mail', () => ({
  mailService: {
    getFail2banStatus: vi.fn(),
    banIP: vi.fn(),
    unbanIP: vi.fn(),
  },
}));

const createWrapper = () => {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: { retry: false },
      mutations: { retry: false },
    },
  });
  return ({ children }: { children: React.ReactNode }) => (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  );
};

const mockFail2banData: Fail2banStatus = {
  jails: ['postfix', 'dovecot'],
  checkedAt: '2026-11-07T12:00:00Z',
  canManage: true,
  details: {
    postfix: {
      jail: 'postfix',
      currentlyFailed: 2,
      totalFailed: 10,
      currentlyBanned: 1,
      totalBanned: 5,
      bannedIps: ['198.51.100.1'],
    },
    dovecot: {
      jail: 'dovecot',
      currentlyFailed: 0,
      totalFailed: 2,
      currentlyBanned: 0,
      totalBanned: 1,
      bannedIps: [],
    },
  },
};

describe('Fail2banDrawer', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders jails and banned IP list when open', async () => {
    vi.mocked(mailService.getFail2banStatus).mockResolvedValue(mockFail2banData);

    render(
      <Fail2banDrawer
        serverId="srv-mail-1"
        serverName="mail.example.com"
        isOpen={true}
        onClose={vi.fn()}
      />,
      { wrapper: createWrapper() }
    );

    await waitFor(() => {
      expect(screen.getByText('Fail2ban Brute-force Protection')).toBeDefined();
      expect(screen.getAllByText('postfix').length).toBeGreaterThanOrEqual(1);
      expect(screen.getAllByText('dovecot').length).toBeGreaterThanOrEqual(1);
      expect(screen.getByText('198.51.100.1')).toBeDefined();
    });
  });

  it('opens confirmation bar and calls unban on confirm', async () => {
    vi.mocked(mailService.getFail2banStatus).mockResolvedValue(mockFail2banData);
    vi.mocked(mailService.unbanIP).mockResolvedValue({ success: true });

    render(
      <Fail2banDrawer
        serverId="srv-mail-1"
        serverName="mail.example.com"
        isOpen={true}
        onClose={vi.fn()}
      />,
      { wrapper: createWrapper() }
    );

    await waitFor(() => {
      expect(screen.getByText('198.51.100.1')).toBeDefined();
    });

    const unbanBtn = screen.getByText('Unban');
    fireEvent.click(unbanBtn);

    expect(screen.getByText('Action executes immediately on the live server.')).toBeDefined();
    expect(screen.getByText('Confirm')).toBeDefined();

    fireEvent.click(screen.getByText('Confirm'));

    await waitFor(() => {
      expect(mailService.unbanIP).toHaveBeenCalledWith('srv-mail-1', '198.51.100.1', 'postfix');
    });
  });

  it('reviews and submits a manual IP ban', async () => {
    vi.mocked(mailService.getFail2banStatus).mockResolvedValue(mockFail2banData);
    vi.mocked(mailService.banIP).mockResolvedValue({ success: true });

    render(
      <Fail2banDrawer
        serverId="srv-mail-1"
        serverName="mail.example.com"
        isOpen={true}
        onClose={vi.fn()}
      />,
      { wrapper: createWrapper() }
    );

    await waitFor(() => {
      expect(screen.getByPlaceholderText(/198.51.100.42/i)).toBeDefined();
    });

    const ipInput = screen.getByPlaceholderText(/198.51.100.42/i);
    fireEvent.change(ipInput, { target: { value: '203.0.113.50' } });

    fireEvent.click(screen.getByText('Review ban'));

    expect(screen.getByText('Confirm')).toBeDefined();
    fireEvent.click(screen.getByText('Confirm'));

    await waitFor(() => {
      expect(mailService.banIP).toHaveBeenCalledWith('srv-mail-1', '203.0.113.50', 'postfix');
    });
  });
});
