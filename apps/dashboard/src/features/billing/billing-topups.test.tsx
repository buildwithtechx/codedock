import { act, fireEvent, render, screen } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { BillingTopups } from './billing-topups';

describe('BillingTopups', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders credit balance and standard credit packs', () => {
    render(<BillingTopups balanceCredits={15400} />);

    expect(screen.getByText('Prepaid Compute Credits')).toBeDefined();
    expect(screen.getByText(/15,400/)).toBeDefined();
    expect(screen.getByText('10,000 credits')).toBeDefined();
    expect(screen.getByText('55,000 credits')).toBeDefined();
    expect(screen.getByText('+10%')).toBeDefined();
    expect(screen.getByText('+15%')).toBeDefined();
    expect(screen.getByText('+20%')).toBeDefined();
  });

  it('triggers onBuy with pack and idempotency key when clicked', async () => {
    const handleBuy = vi.fn().mockResolvedValue(undefined);
    render(<BillingTopups onBuy={handleBuy} />);

    const buyFiftyBtn = screen.getByRole('button', { name: /^buy \$50$/i });
    await act(async () => {
      fireEvent.click(buyFiftyBtn);
    });

    expect(handleBuy).toHaveBeenCalledTimes(1);
    expect(handleBuy).toHaveBeenCalledWith(
      expect.objectContaining({ id: 'topup_50', priceCents: 5000 }),
      expect.stringMatching(/^topup_topup_50_/)
    );
  });

  it('displays error message when purchase fails', async () => {
    const failingBuy = vi.fn().mockRejectedValue(new Error('Gateway failure'));
    render(<BillingTopups onBuy={failingBuy} />);

    const buyTwentyFiveBtn = screen.getByRole('button', { name: /^buy \$25$/i });
    await act(async () => {
      fireEvent.click(buyTwentyFiveBtn);
    });

    expect(screen.getByRole('alert')).toBeDefined();
    expect(screen.getByText(/failed to initiate checkout/i)).toBeDefined();
  });
});
