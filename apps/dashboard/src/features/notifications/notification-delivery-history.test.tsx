import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { NotificationDeliveryHistory } from './notification-delivery-history';

const mocks = vi.hoisted(() => ({ get: vi.fn() }));
vi.mock('#/lib/api-client', () => ({ apiClient: mocks }));

function renderHistory() {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  const view = render(
    <QueryClientProvider client={client}>
      <NotificationDeliveryHistory />
    </QueryClientProvider>
  );
  return { view, client };
}

describe('notification delivery history', () => {
  it('requests the deliveries endpoint with a limit', async () => {
    mocks.get.mockResolvedValueOnce({ data: [] });
    const { client } = renderHistory();
    await screen.findByText('No deliveries yet');
    expect(mocks.get).toHaveBeenCalledWith('/notifications/deliveries?limit=20');
    client.clear();
  });

  it('lists delivered and failed attempts', async () => {
    mocks.get.mockResolvedValueOnce({
      data: [
        {
          id: 'd-1',
          eventType: 'test_global_discord',
          channel: 'discord',
          title: 'hello',
          success: true,
          error: '',
          createdAt: new Date().toISOString(),
        },
        {
          id: 'd-2',
          eventType: 'test_global_slack',
          channel: 'slack',
          title: '',
          success: false,
          error: 'connection refused',
          createdAt: new Date().toISOString(),
        },
      ],
    });
    const { client } = renderHistory();
    await screen.findByText('Delivered');
    await screen.findByText('Failed');
    await screen.findByText(/connection refused/);
    client.clear();
  });
});
