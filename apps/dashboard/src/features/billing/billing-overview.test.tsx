import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, fireEvent, render, screen } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { BillingOverview } from './billing-overview';

vi.mock('#/services/billing', () => ({
  billingService: {
    createPortalSession: vi.fn(),
  },
}));

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

describe('BillingOverview', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders hobby plan with free badge and upgrade button', () => {
    const handleUpgrade = vi.fn();
    renderWithClient(<BillingOverview planType="hobby" onUpgradeClick={handleUpgrade} />);

    expect(screen.getByText('Hobby Developer')).toBeDefined();
    expect(screen.getByText('Free')).toBeDefined();
    const upgradeBtn = screen.getByRole('button', { name: /upgrade to pro/i });
    fireEvent.click(upgradeBtn);
    expect(handleUpgrade).toHaveBeenCalledTimes(1);
  });

  it('renders pro plan with renewal date and manage billing button', () => {
    renderWithClient(<BillingOverview planType="pro" renewalDate="Nov 20, 2026" />);

    expect(screen.getByText('Pro Platform')).toBeDefined();
    expect(screen.getByText('Active')).toBeDefined();
    expect(screen.getByText('Renews on Nov 20, 2026')).toBeDefined();
    expect(screen.getByRole('button', { name: /manage billing/i })).toBeDefined();
  });

  it('triggers portal session when manage billing is clicked', async () => {
    const portalFn = vi.fn().mockResolvedValue(undefined);
    renderWithClient(<BillingOverview planType="pro" onPortalClick={portalFn} />);

    const manageBtn = screen.getByRole('button', { name: /manage billing/i });
    await act(async () => {
      fireEvent.click(manageBtn);
    });
    expect(portalFn).toHaveBeenCalledTimes(1);
  });

  it('renders canceling state with scheduled end notice', () => {
    renderWithClient(
      <BillingOverview planType="pro" cancelAtPeriodEnd={true} renewalDate="Dec 1, 2026" />
    );

    expect(screen.getByText('Canceling')).toBeDefined();
    expect(screen.getByText('Access remains active until Dec 1, 2026')).toBeDefined();
  });
});
