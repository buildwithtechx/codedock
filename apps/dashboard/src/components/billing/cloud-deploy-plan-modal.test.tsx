import { describe, expect, it, vi } from 'vitest';
import { render, screen, fireEvent, act } from '@testing-library/react';
import { CloudDeployPlanModal } from './cloud-deploy-plan-modal';

describe('CloudDeployPlanModal', () => {
  it('renders modal with restriction message and default plans', () => {
    const handleClose = vi.fn();
    render(
      <CloudDeployPlanModal
        restriction={{
          code: 'PLAN_UPGRADE_REQUIRED',
          message: 'Your concurrent build limit was reached.',
        }}
        onClose={handleClose}
      />
    );

    expect(screen.getByText('Choose your deployment capacity')).toBeTruthy();
    expect(screen.getByText('Your concurrent build limit was reached.')).toBeTruthy();
    expect(screen.getByText('Starter')).toBeTruthy();
    expect(screen.getByText('Pro')).toBeTruthy();
    expect(screen.getByText('Scale')).toBeTruthy();

    const closeBtn = screen.getByRole('button', { name: 'Close dialog' });
    fireEvent.click(closeBtn);
    expect(handleClose).toHaveBeenCalledOnce();
  });

  it('selects plan and triggers upgrade callback', async () => {
    const handleClose = vi.fn();
    const handleSelect = vi.fn().mockResolvedValue(undefined);
    const handleRetry = vi.fn().mockResolvedValue(undefined);

    render(
      <CloudDeployPlanModal
        onClose={handleClose}
        onSelectPlan={handleSelect}
        onRetry={handleRetry}
      />
    );

    const scaleCard = screen.getByText('Scale');
    await act(async () => {
      fireEvent.click(scaleCard);
    });

    const submitBtn = screen.getByRole('button', { name: 'Select plan' });
    await act(async () => {
      fireEvent.click(submitBtn);
    });

    expect(handleSelect).toHaveBeenCalledWith('scale');
  });

  it('closes dialog on cancel button click', () => {
    const handleClose = vi.fn();
    render(<CloudDeployPlanModal onClose={handleClose} />);

    const cancelBtn = screen.getByRole('button', { name: 'Cancel' });
    fireEvent.click(cancelBtn);
    expect(handleClose).toHaveBeenCalledOnce();
  });
});
