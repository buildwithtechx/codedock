import { cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { afterEach, describe, expect, it, vi } from 'vitest';
import { installCommand } from '../lib/product-links';
import { InstallCommand } from './install-command';

afterEach(() => {
  cleanup();
  vi.unstubAllGlobals();
});
describe('install command clipboard', () => {
  it('copies the exact installer command and announces completion', async () => {
    const writeText = vi.fn().mockResolvedValue(undefined);
    vi.stubGlobal('navigator', { clipboard: { writeText } });
    render(<InstallCommand />);
    fireEvent.click(screen.getByRole('button', { name: 'Copy install command' }));
    await waitFor(() => expect(screen.getByText('Install command copied.')).toBeTruthy());
    expect(writeText).toHaveBeenCalledWith(installCommand);
  });
  it('keeps the command available when clipboard access fails', async () => {
    vi.stubGlobal('navigator', {
      clipboard: { writeText: vi.fn().mockRejectedValue(new Error('denied')) },
    });
    render(<InstallCommand />);
    fireEvent.click(screen.getByRole('button', { name: 'Copy install command' }));
    await waitFor(() =>
      expect(
        screen.getByText('Clipboard unavailable. Select and copy the command above.')
      ).toBeTruthy()
    );
    expect(screen.getByText(installCommand)).toBeTruthy();
  });
});
