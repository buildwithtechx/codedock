import { fireEvent, render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { CreditAlertBanner } from './credit-alert-banner';

describe('CreditAlertBanner', () => {
  it('renders nothing when alert is null or state is ok', () => {
    const { container: c1 } = render(<CreditAlertBanner alert={null} />);
    expect(c1.firstChild).toBeNull();

    const { container: c2 } = render(<CreditAlertBanner alert={{ state: 'ok' }} />);
    expect(c2.firstChild).toBeNull();

    const { container: c3 } = render(
      <CreditAlertBanner alert={{ state: 'depleted' }} billingMode="monthly" />
    );
    expect(c3.firstChild).toBeNull();
  });

  it('renders low credit warning and fires top-up action', () => {
    const onTopUp = vi.fn();
    const onViewPlans = vi.fn();

    render(
      <CreditAlertBanner
        alert={{ state: 'low', percent: 12 }}
        onTopUp={onTopUp}
        onViewPlans={onViewPlans}
      />
    );

    expect(screen.getByText('Build credits running low')).toBeTruthy();
    expect(screen.getByText(/Your balance is below 15% \(12% remaining\)/)).toBeTruthy();

    const topUpBtn = screen.getByRole('button', { name: 'Add credits' });
    fireEvent.click(topUpBtn);
    expect(onTopUp).toHaveBeenCalledOnce();

    const plansBtn = screen.getByRole('button', { name: /View plans/ });
    fireEvent.click(plansBtn);
    expect(onViewPlans).toHaveBeenCalledOnce();
  });

  it('renders grace period active warning', () => {
    render(<CreditAlertBanner alert={{ state: 'grace' }} />);

    expect(screen.getByText('Grace period active')).toBeTruthy();
    expect(screen.getByText(/Your balance is in grace credit usage/)).toBeTruthy();
  });

  it('renders critical depleted gate when credits are exhausted', () => {
    render(<CreditAlertBanner alert={{ state: 'depleted' }} />);

    expect(screen.getByText('Build credits exhausted')).toBeTruthy();
    expect(
      screen.getByText(/New builds and deployments are paused until credits are added/)
    ).toBeTruthy();
  });
});
