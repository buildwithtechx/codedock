import { fireEvent, render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { PayAsYouGoPlan } from './pay-as-you-go-plan';

describe('PayAsYouGoPlan', () => {
  it('renders sizing sliders and pool limits', () => {
    render(<PayAsYouGoPlan />);
    expect(screen.getByText('Host resource sizing & pool limits')).toBeDefined();
    expect(screen.getAllByText('Dedicated vCPU').length).toBeGreaterThan(0);
    expect(screen.getByText('Estimated Burn Rate')).toBeDefined();
  });

  it('updates selected deposit package when clicked', () => {
    const handleAddCredits = vi.fn();
    render(<PayAsYouGoPlan onAddCredits={handleAddCredits} />);

    // Click on $50 package
    const packageBtn = screen.getByRole('button', { name: /\$50/i });
    fireEvent.click(packageBtn);

    const submitBtn = screen.getByRole('button', { name: /Add Credits & Deploy/i });
    fireEvent.click(submitBtn);

    expect(handleAddCredits).toHaveBeenCalledWith(5000);
  });

  it('allows switching tiers', () => {
    render(<PayAsYouGoPlan />);
    const tier2Tab = screen.getByRole('tab', { name: /Tier 2/i });
    fireEvent.click(tier2Tab);
    expect(tier2Tab.getAttribute('aria-selected')).toBe('true');
  });
});
