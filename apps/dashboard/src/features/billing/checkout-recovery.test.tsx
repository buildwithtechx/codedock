import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, fireEvent, render, screen } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import type { PendingCheckout } from '#/interfaces/billing';
import { CheckoutRecovery } from './checkout-recovery';

const createTestQueryClient = () =>
  new QueryClient({
    defaultOptions: {
      queries: { retry: false },
      mutations: { retry: false },
    },
  });

function renderWithClient(ui: React.ReactElement) {
  const queryClient = createTestQueryClient();
  return render(<QueryClientProvider client={queryClient}>{ui}</QueryClientProvider>);
}

describe('CheckoutRecovery', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  const mockPending: PendingCheckout = {
    checkoutId: 'cs_test_123',
    idempotencyKey: 'idem_key_456',
    priceId: 'price_pro_month',
    url: 'https://checkout.stripe.com/pay/cs_test_123',
    status: 'open',
    canCancel: true,
    createdAt: '2026-10-16T14:00:00Z',
  };

  it('renders nothing when pending is null or status is complete', () => {
    const { container: empty } = renderWithClient(<CheckoutRecovery pending={null} />);
    expect(empty.firstChild).toBeNull();

    const { container: complete } = renderWithClient(
      <CheckoutRecovery pending={{ ...mockPending, status: 'complete' }} />
    );
    expect(complete.firstChild).toBeNull();
  });

  it('renders unsettled payment banner with resume payment link', () => {
    renderWithClient(<CheckoutRecovery pending={mockPending} />);

    expect(screen.getByText('Unsettled Payment in Progress')).toBeDefined();
    expect(screen.getByText(/open/i)).toBeDefined();

    const resumeLink = screen.getByRole('link', { name: /resume payment/i });
    expect(resumeLink.getAttribute('href')).toBe('https://checkout.stripe.com/pay/cs_test_123');
  });

  it('cancels checkout with id and idempotency key when clicked', async () => {
    const handleCancel = vi.fn().mockResolvedValue(undefined);
    const handleSettled = vi.fn();
    renderWithClient(
      <CheckoutRecovery pending={mockPending} onCancel={handleCancel} onSettled={handleSettled} />
    );

    const cancelBtn = screen.getByRole('button', { name: /cancel checkout/i });
    await act(async () => {
      fireEvent.click(cancelBtn);
    });

    expect(handleCancel).toHaveBeenCalledTimes(1);
    expect(handleCancel).toHaveBeenCalledWith('cs_test_123', 'idem_key_456');
    expect(handleSettled).toHaveBeenCalledTimes(1);
  });

  it('disables cancel and displays reason when canCancel is false', () => {
    renderWithClient(
      <CheckoutRecovery
        pending={{
          ...mockPending,
          canCancel: false,
          cancelUnavailableReason: 'provider_unavailable',
        }}
      />
    );

    const cancelBtn = screen.getByRole('button', { name: /cancel checkout/i });
    expect((cancelBtn as HTMLButtonElement).disabled).toBe(true);
    expect(screen.getByText(/cancellation is temporarily unavailable/i)).toBeDefined();
  });

  it('retains session on uncertain outcome failure and displays alert', async () => {
    const failingCancel = vi.fn().mockRejectedValue(new Error('Network error'));
    renderWithClient(<CheckoutRecovery pending={mockPending} onCancel={failingCancel} />);

    const cancelBtn = screen.getByRole('button', { name: /cancel checkout/i });
    await act(async () => {
      fireEvent.click(cancelBtn);
    });

    expect(screen.getByRole('alert')).toBeDefined();
    expect(screen.getByText(/retaining checkout session for retry/i)).toBeDefined();
  });
});
