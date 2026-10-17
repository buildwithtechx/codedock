import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { fireEvent, render, screen } from '@testing-library/react';
import type { ReactNode } from 'react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import type { ActionRunner } from '#/services/actions';
import { RunnerCatalog } from './runner-catalog';

vi.mock('@tanstack/react-router', () => ({
  Link: ({ children, to, ...props }: { children?: ReactNode; to?: string }) => (
    <a href={to} {...props}>
      {children}
    </a>
  ),
  useNavigate: () => vi.fn(),
}));

vi.mock('#/services/servers', () => ({
  serverService: {
    list: vi
      .fn()
      .mockResolvedValue([{ id: 'server-1', name: 'Dedicated Node', ipAddress: '1.2.3.4' }]),
  },
}));

function renderWithProviders(ui: ReactNode) {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  });
  return render(<QueryClientProvider client={queryClient}>{ui}</QueryClientProvider>);
}

const mockRunners: ActionRunner[] = [
  {
    id: 'runner-1',
    name: 'Ubuntu 22.04 Fast',
    mode: 'container',
    image: 'catthehacker/ubuntu:act-22.04',
    status: 'online',
    labels: ['ubuntu-latest', 'docker'],
    enabled: true,
    lastHeartbeatAt: '2026-10-16T12:00:00Z',
    createdAt: '2026-10-16T10:00:00Z',
    updatedAt: '2026-10-16T12:00:00Z',
  },
];

describe('RunnerCatalog', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders runner cards and status', () => {
    renderWithProviders(<RunnerCatalog runners={mockRunners} />);

    expect(screen.getByText('Ubuntu 22.04 Fast')).toBeDefined();
    expect(screen.getByText('online')).toBeDefined();
    expect(screen.getByText('catthehacker/ubuntu:act-22.04')).toBeDefined();
    expect(screen.getByText('ubuntu-latest')).toBeDefined();
  });

  it('renders empty state when no runners exist', () => {
    renderWithProviders(<RunnerCatalog runners={[]} />);

    expect(screen.getByText('No runners registered')).toBeDefined();
  });

  it('opens register runner modal on button click', async () => {
    renderWithProviders(<RunnerCatalog runners={mockRunners} />);

    const registerBtn = screen.getByRole('button', { name: /register runner/i });
    fireEvent.click(registerBtn);

    expect(await screen.findByRole('heading', { level: 1, name: /add runner/i })).toBeDefined();
    expect(screen.getByLabelText(/runner name/i)).toBeDefined();
  });
});
