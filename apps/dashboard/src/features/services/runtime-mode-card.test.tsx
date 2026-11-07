import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { fireEvent, render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { RuntimeModeCard } from './runtime-mode-card';

vi.mock('#/lib/api-client', () => ({
  apiClient: {
    get: vi.fn().mockResolvedValue({ data: { data: {} } }),
    put: vi.fn().mockResolvedValue({ data: { success: true } }),
  },
}));

describe('RuntimeModeCard', () => {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  });

  const renderComponent = (initialData?: {
    runtimeMode?: 'web' | 'worker' | 'static';
    internalPort?: number;
    spaFallback?: boolean;
  }) =>
    render(
      <QueryClientProvider client={queryClient}>
        <RuntimeModeCard serviceId="test-svc" initialData={initialData} />
      </QueryClientProvider>
    );

  it('renders web mode with internal port input by default', () => {
    renderComponent({ runtimeMode: 'web', internalPort: 3000 });
    expect(screen.getByText('Web Service')).toBeDefined();
    expect(screen.getByLabelText(/internal port/i)).toBeDefined();
    expect(screen.queryByText('SPA fallback')).toBeNull();
  });

  it('renders SPA fallback toggle when static mode is active', () => {
    renderComponent({ runtimeMode: 'static', spaFallback: true });
    expect(screen.getByText('Static Site')).toBeDefined();
    expect(screen.getByText('SPA fallback')).toBeDefined();
    expect(screen.getByText(/Static sites: serve index\.html for missing paths/i)).toBeDefined();
    expect(screen.queryByLabelText(/internal port/i)).toBeNull();
  });

  it('allows switching to static mode to reveal SPA fallback control', () => {
    renderComponent({ runtimeMode: 'web' });
    expect(screen.queryByText('SPA fallback')).toBeNull();

    const staticButton = screen.getByRole('button', { name: /static site/i });
    fireEvent.click(staticButton);

    expect(screen.getByText('SPA fallback')).toBeDefined();
  });
});
