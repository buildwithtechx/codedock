import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { CloudDeployPlanModal } from './cloud-deploy-plan-modal';

vi.mock('#/hooks/use-billing', () => ({
  useBillingConfig: () => ({
    data: {
      publishableKey: 'pk_test_123',
      plans: [{ id: 'pro', name: 'Pro', price: 39, priceId: 'price_pro' }],
    },
  }),
  useCreateCheckoutSession: () => ({
    mutateAsync: vi.fn(),
    isPending: false,
  }),
}));

vi.mock('#/stores/auth-store', () => ({
  useAuthStore: (selector: any) =>
    selector({
      user: { id: 'u_1', planType: 'free' },
    }),
}));

describe('CloudDeployPlanModal', () => {
  it('renders modal with title and plan options', () => {
    render(<CloudDeployPlanModal open={true} onClose={vi.fn()} />);
    expect(screen.getByText('Choose a Cloud plan')).toBeDefined();
    expect(screen.getByText('Starter')).toBeDefined();
    expect(screen.getByText('Pro')).toBeDefined();
  });

  it('renders capacity limit title when code is CLOUD_CAPACITY_REQUIRED', () => {
    render(
      <CloudDeployPlanModal
        open={true}
        onClose={vi.fn()}
        restriction={{ code: 'CLOUD_CAPACITY_REQUIRED', message: 'Maximum workspaces reached' }}
      />
    );
    expect(screen.getByText('Compute capacity limit reached')).toBeDefined();
    expect(screen.getByText('Maximum workspaces reached')).toBeDefined();
  });
});
