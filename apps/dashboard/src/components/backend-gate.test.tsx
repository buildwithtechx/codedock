import { render, screen, waitFor } from '@testing-library/react';
import { afterEach, describe, expect, it, vi } from 'vitest';
import { BackendGate } from './backend-gate';

describe('BackendGate', () => {
  afterEach(() => {
    vi.unstubAllGlobals();
  });

  it('renders children once the backend responds', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({ ok: true }));
    render(
      <BackendGate>
        <p>app ready</p>
      </BackendGate>
    );
    expect(await screen.findByText('app ready')).toBeDefined();
  });

  it('holds a waiting state while the backend is unreachable', async () => {
    const fetchMock = vi.fn().mockRejectedValue(new Error('down'));
    vi.stubGlobal('fetch', fetchMock);
    render(
      <BackendGate>
        <p>app ready</p>
      </BackendGate>
    );
    await waitFor(() => expect(fetchMock).toHaveBeenCalled());
    expect(screen.queryByText('app ready')).toBeNull();
    expect(screen.getByText(/Waiting for the daemon/)).toBeDefined();
  });
});
