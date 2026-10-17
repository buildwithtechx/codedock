import { act, fireEvent, render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { CloudConnectionCard } from './cloud-connection-card';

describe('CloudConnectionCard', () => {
  it('renders disconnected state with perks and connect action', () => {
    const handleConnect = vi.fn();
    render(<CloudConnectionCard connected={false} onConnect={handleConnect} />);

    expect(screen.getByText('Connect to CodeDock Cloud')).toBeTruthy();
    expect(screen.getByText(/Free wildcard domains on \*\.codedock\.app/)).toBeTruthy();

    const btn = screen.getByRole('button', { name: /Connect Cloud/ });
    fireEvent.click(btn);
    expect(handleConnect).toHaveBeenCalledOnce();
  });

  it('renders loading state when checking cloud connection', () => {
    render(<CloudConnectionCard loading={true} />);
    expect(screen.getByText('Checking cloud connection...')).toBeTruthy();
  });

  it('renders connected state with user details and disconnect action', async () => {
    const handleDisconnect = vi.fn();
    vi.spyOn(window, 'confirm').mockReturnValue(true);

    render(
      <CloudConnectionCard
        connected={true}
        cloudUser={{ name: 'Jane Dev', email: 'jane@example.com' }}
        onDisconnect={handleDisconnect}
      />
    );

    expect(screen.getByText('Jane Dev')).toBeTruthy();
    expect(screen.getByText('jane@example.com')).toBeTruthy();
    expect(screen.getByText('Connected')).toBeTruthy();

    const disconnectBtn = screen.getByRole('button', { name: /Disconnect/ });
    await act(async () => {
      fireEvent.click(disconnectBtn);
    });
    expect(handleDisconnect).toHaveBeenCalledOnce();
  });

  it('does not disconnect when user cancels confirmation', () => {
    const handleDisconnect = vi.fn();
    vi.spyOn(window, 'confirm').mockReturnValue(false);

    render(
      <CloudConnectionCard
        connected={true}
        cloudUser={{ name: 'Jane Dev' }}
        onDisconnect={handleDisconnect}
      />
    );

    const disconnectBtn = screen.getByRole('button', { name: /Disconnect/ });
    fireEvent.click(disconnectBtn);
    expect(handleDisconnect).not.toHaveBeenCalled();
  });
});
