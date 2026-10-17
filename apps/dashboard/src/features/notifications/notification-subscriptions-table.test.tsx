import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { render, screen } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { NotificationSubscriptionsTable } from './notification-subscriptions-table';

vi.mock('#/features/settings/api', () => ({
  settingsService: {
    getSubscriptions: vi.fn().mockResolvedValue({ data: [] }),
    getDefaults: vi.fn().mockResolvedValue({ data: [] }),
    upsertSubscription: vi.fn().mockResolvedValue({}),
    upsertDefault: vi.fn().mockResolvedValue({}),
  },
}));

describe('NotificationSubscriptionsTable (Alert preferences)', () => {
  let queryClient: QueryClient;

  beforeEach(() => {
    queryClient = new QueryClient({
      defaultOptions: {
        queries: { retry: false },
      },
    });
  });

  it('renders Alert preferences header with exact 24 events, underline tabs, and card controls', () => {
    render(
      <QueryClientProvider client={queryClient}>
        <NotificationSubscriptionsTable />
      </QueryClientProvider>
    );

    expect(screen.getByText('Alert preferences')).toBeDefined();
    expect(
      screen.getByText(
        'Organization defaults apply when a member has not made an explicit choice for an alert. Members can override them with their own preferences.'
      )
    ).toBeDefined();

    expect(screen.getByText('24')).toBeDefined();
    expect(screen.getAllByText('Deployment').length).toBeGreaterThan(0);
    expect(screen.getAllByText('App health').length).toBeGreaterThan(0);
    expect(screen.getByText('Deploy failed')).toBeDefined();
    expect(screen.getAllByText('Notify me').length).toBeGreaterThan(0);
    expect(screen.getAllByText('Organization defaults').length).toBeGreaterThan(0);
  });
});
