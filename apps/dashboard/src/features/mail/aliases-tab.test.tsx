import type { MailAlias } from '@codedock/contracts';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { mailService } from '#/services/mail';
import { AliasesTab } from './aliases-tab';

vi.mock('#/services/mail', () => ({
  mailService: {
    listAliases: vi.fn(),
    createAlias: vi.fn(),
    deleteAlias: vi.fn(),
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

const mockAliases: MailAlias[] = [
  {
    id: 'alias-1',
    domainId: 'dom-1',
    source: 'support@example.com',
    destination: 'alice@example.com',
    active: true,
    createdAt: '2026-11-07T12:00:00Z',
    updatedAt: '2026-11-07T12:00:00Z',
  },
];

describe('AliasesTab', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders aliases list and strict sender ownership notice', async () => {
    vi.mocked(mailService.listAliases).mockResolvedValue(mockAliases);

    render(<AliasesTab serverId="srv-1" domainId="dom-1" domainName="example.com" />, {
      wrapper: createWrapper(),
    });

    await waitFor(() => {
      expect(screen.getByText('support@example.com')).toBeDefined();
      expect(screen.getByText('alice@example.com')).toBeDefined();
      expect(screen.getByText(/Strict direct alias ownership/i)).toBeDefined();
    });
  });

  it('creates an alias when form is submitted', async () => {
    vi.mocked(mailService.listAliases).mockResolvedValue([]);
    vi.mocked(mailService.createAlias).mockResolvedValue({
      id: 'alias-new',
      domainId: 'dom-1',
      source: 'billing@example.com',
      destination: 'bob@example.com',
      active: true,
      createdAt: '2026-11-07T12:00:00Z',
      updatedAt: '2026-11-07T12:00:00Z',
    });

    render(<AliasesTab serverId="srv-1" domainId="dom-1" domainName="example.com" />, {
      wrapper: createWrapper(),
    });

    await waitFor(() => {
      expect(screen.getByText('Create alias')).toBeDefined();
    });

    fireEvent.click(screen.getByText('Create alias'));

    const sourceInput = screen.getByPlaceholderText('support');
    const destInput = screen.getByPlaceholderText('alice@example.com');

    fireEvent.change(sourceInput, { target: { value: 'billing' } });
    fireEvent.change(destInput, { target: { value: 'bob@example.com' } });

    const submitBtn = screen.getByRole('button', { name: /^Create Alias$/i });
    fireEvent.click(submitBtn);

    await waitFor(() => {
      expect(mailService.createAlias).toHaveBeenCalledWith('srv-1', 'dom-1', {
        source: 'billing@example.com',
        destination: 'bob@example.com',
      });
    });
  });
});
