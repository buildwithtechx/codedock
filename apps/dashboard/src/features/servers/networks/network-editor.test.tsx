import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { networksService } from '#/services/networks';
import { NetworkEditor } from './network-editor';

describe('NetworkEditor', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders initial form fields for WireGuard mesh', () => {
    render(<NetworkEditor />);

    expect(screen.getByText('Private Network Mesh')).toBeDefined();
    expect(screen.getByLabelText(/network name/i)).toBeDefined();
    expect(screen.getByLabelText(/subnet cidr/i)).toBeDefined();
    expect(screen.getByRole('button', { name: /managed wireguard/i })).toBeDefined();
    expect(screen.getByRole('button', { name: /apply & arm watchdog/i })).toBeDefined();
    expect(screen.getByText(/member nodes \(2\)/i)).toBeDefined();
  });

  it('adds a member server when Add Server is clicked', () => {
    render(<NetworkEditor />);

    const addBtn = screen.getByRole('button', { name: /add server/i });
    fireEvent.click(addBtn);

    expect(screen.getByText(/member nodes \(3\)/i)).toBeDefined();
  });

  it('submits create request and arms 120s rollback watchdog', async () => {
    const createSpy = vi.spyOn(networksService, 'create').mockResolvedValue({
      id: 'net-mesh-1',
      name: 'beta-mesh',
      mode: 'wireguard',
      subnetCidr: '10.42.0.0/16',
      mtu: 1420,
      verificationPort: 51821,
      status: 'applying',
      revision: 1,
      members: [
        { serverId: 'srv-1', privateIp: '10.42.0.1', interfaceName: 'wg0', status: 'ready', latencyMs: 1.2 },
        { serverId: 'srv-2', privateIp: '10.42.0.2', interfaceName: 'wg0', status: 'ready', latencyMs: 1.9 },
      ],
      accessPolicy: { mode: 'full_mesh' },
      watchdog: {
        operationId: 'op-999',
        armed: true,
        timeoutSeconds: 120,
        remainingSeconds: 120,
        expiresAt: '2026-10-16T12:02:00Z',
        canCommit: true,
        canRollback: true,
        status: 'armed',
      },
      createdAt: '2026-10-16T12:00:00Z',
      updatedAt: '2026-10-16T12:00:00Z',
    });

    const onSaved = vi.fn();
    render(<NetworkEditor onSaved={onSaved} />);

    const nameInput = screen.getByLabelText(/network name/i);
    fireEvent.change(nameInput, { target: { value: 'beta-mesh' } });

    const submitBtn = screen.getByRole('button', { name: /apply & arm watchdog/i });
    fireEvent.click(submitBtn);

    await waitFor(() => {
      expect(createSpy).toHaveBeenCalledWith(
        expect.objectContaining({
          name: 'beta-mesh',
          mode: 'wireguard',
        })
      );
      expect(onSaved).toHaveBeenCalled();
      expect(screen.getByText(/Rollback Watchdog Armed \(120s\)/i)).toBeDefined();
    });
  });
});
