import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { actionsService, type ActionRunner } from '#/services/actions';
import { RunnerCatalog } from './runner-catalog';

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
    render(<RunnerCatalog runners={mockRunners} />);

    expect(screen.getByText('Ubuntu 22.04 Fast')).toBeDefined();
    expect(screen.getByText('online')).toBeDefined();
    expect(screen.getByText('catthehacker/ubuntu:act-22.04')).toBeDefined();
    expect(screen.getByText('ubuntu-latest')).toBeDefined();
  });

  it('renders empty state when no runners exist', () => {
    render(<RunnerCatalog runners={[]} />);

    expect(screen.getByText('No runners registered')).toBeDefined();
  });

  it('opens register runner modal on button click', async () => {
    render(<RunnerCatalog runners={mockRunners} />);

    const registerBtn = screen.getByRole('button', { name: /register runner/i });
    fireEvent.click(registerBtn);

    expect(await screen.findByText('Register Action Runner')).toBeDefined();
    expect(screen.getByLabelText(/runner name/i)).toBeDefined();
  });
});
