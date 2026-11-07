import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { cleanup, render, screen } from '@testing-library/react';
import type { ReactNode } from 'react';
import { afterEach, describe, expect, it, vi } from 'vitest';
import { PasskeysSection } from '#/features/profile/passkeys-section';

const mockUsePasskeys = vi.fn();
const mockRegisterMutate = vi.fn();
const mockRenameMutate = vi.fn();
const mockDeleteMutate = vi.fn();

vi.mock('#/hooks/use-passkeys', () => ({
  usePasskeys: () => mockUsePasskeys(),
  useRegisterPasskey: () => ({ mutate: mockRegisterMutate, isPending: false }),
  useRenamePasskey: () => ({ mutate: mockRenameMutate, isPending: false }),
  useDeletePasskey: () => ({ mutate: mockDeleteMutate, isPending: false }),
}));

vi.mock('@simplewebauthn/browser', () => ({
  browserSupportsWebAuthn: () => true,
}));

function renderWithClient(ui: ReactNode) {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  });
  return render(<QueryClientProvider client={queryClient}>{ui}</QueryClientProvider>);
}

describe('PasskeysSection', () => {
  afterEach(() => {
    cleanup();
    vi.clearAllMocks();
  });

  it('renders empty state when user has no passkeys', () => {
    mockUsePasskeys.mockReturnValue({
      data: { data: [] },
      isLoading: false,
    });

    renderWithClient(<PasskeysSection />);
    expect(screen.getByText('Passkeys (WebAuthn)')).toBeDefined();
    expect(screen.getByText('No passkeys registered')).toBeDefined();
    expect(screen.getByRole('button', { name: /add passkey/i })).toBeDefined();
  });

  it('renders passkeys list with friendly metadata', () => {
    mockUsePasskeys.mockReturnValue({
      data: {
        data: [
          {
            id: 'pk-1',
            name: 'MacBook Touch ID',
            backedUp: true,
            deviceType: 'multiDevice',
            transports: ['internal'],
            createdAt: '2026-10-19T00:00:00Z',
          },
        ],
      },
      isLoading: false,
    });

    renderWithClient(<PasskeysSection />);
    expect(screen.getByText('MacBook Touch ID')).toBeDefined();
    expect(screen.getByText('Synced')).toBeDefined();
  });
});
