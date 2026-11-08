import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { fireEvent, render, screen } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { RouteProtectionButton } from './route-protection-button';

vi.mock('#/hooks/use-route-protection', () => ({
  useRouteProtection: vi.fn(() => ({
    data: [],
    isLoading: false,
  })),
  useChangeRouteProtection: vi.fn(() => ({
    mutateAsync: vi.fn(),
    isPending: false,
  })),
}));

const createWrapper = () => {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: { retry: false },
      mutations: { retry: false },
    },
  });
  return ({ children }: { children: React.ReactNode }) => (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  );
};

describe('RouteProtectionButton', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders protection dialog and password input when supported', () => {
    render(
      <RouteProtectionButton hostname="api.example.com" projectId="proj-1" unsupported={false} />,
      { wrapper: createWrapper() }
    );

    const trigger = screen.getByLabelText('Route Protection: api.example.com');
    fireEvent.click(trigger);

    expect(screen.getByText('Route Password Protection')).toBeDefined();
    const input = document.querySelector('input[type="password"]');
    expect(input).not.toBeNull();
    expect(screen.getByText('Enable Route Protection')).toBeDefined();
  });

  it('suppresses password input and displays warning notice when route is unsupported (e.g. external ingress / tunnel)', () => {
    render(
      <RouteProtectionButton hostname="tunnel.example.com" projectId="proj-1" unsupported={true} />,
      { wrapper: createWrapper() }
    );

    const trigger = screen.getByLabelText('Route Protection: tunnel.example.com');
    fireEvent.click(trigger);

    expect(screen.getByText('Route Password Protection')).toBeDefined();
    expect(screen.getByText('Protection Not Supported')).toBeDefined();
    expect(
      screen.getByText(
        /Password protection is not supported for hostnames behind external ingress/i
      )
    ).toBeDefined();

    // Password input and enable button must be completely suppressed
    const input = document.querySelector('input[type="password"]');
    expect(input).toBeNull();
    expect(screen.queryByText('Enable Route Protection')).toBeNull();
  });
});
