import { cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { afterEach, describe, expect, it, vi } from 'vitest';
import { flushInstallCopies } from '../lib/analytics';
import { installCommand } from '../lib/product-links';
import { InstallCommand } from './install-command';

afterEach(() => {
  cleanup();
  vi.unstubAllGlobals();
  delete window.posthog;
  delete window.__codedock_install_copies;
});
describe('install command clipboard', () => {
  it('copies the exact installer command and announces completion', async () => {
    const writeText = vi.fn().mockResolvedValue(undefined);
    vi.stubGlobal('navigator', { clipboard: { writeText } });
    render(<InstallCommand source="test" />);
    fireEvent.click(screen.getByRole('button', { name: 'Copy install command' }));
    await waitFor(() => expect(screen.getByText('Install command copied.')).toBeTruthy());
    expect(writeText).toHaveBeenCalledWith(installCommand);
  });
  it('keeps the command available when clipboard access fails', async () => {
    vi.stubGlobal('navigator', {
      clipboard: { writeText: vi.fn().mockRejectedValue(new Error('denied')) },
    });
    render(<InstallCommand source="test" />);
    fireEvent.click(screen.getByRole('button', { name: 'Copy install command' }));
    await waitFor(() =>
      expect(
        screen.getByText('Clipboard unavailable. Select and copy the command above.')
      ).toBeTruthy()
    );
    expect(screen.getByText(installCommand)).toBeTruthy();
  });
  it('delivers early copies once analytics loads without sending them twice', async () => {
    vi.stubGlobal('navigator', { clipboard: { writeText: vi.fn().mockResolvedValue(undefined) } });
    render(<InstallCommand source="bottom_cta" />);
    fireEvent.click(screen.getByRole('button', { name: 'Copy install command' }));
    await waitFor(() => expect(screen.getByText('Install command copied.')).toBeTruthy());
    const capture = vi.fn();
    window.posthog = { init: vi.fn(), capture };
    flushInstallCopies();
    expect(capture).toHaveBeenCalledWith('install_command_copied', {
      source: 'bottom_cta',
      command: installCommand,
    });
    flushInstallCopies();
    expect(capture).toHaveBeenCalledTimes(1);
  });
  it('preserves the hero install conversion event after a successful copy', async () => {
    const capture = vi.fn();
    window.posthog = { init: vi.fn(), capture };
    vi.stubGlobal('navigator', { clipboard: { writeText: vi.fn().mockResolvedValue(undefined) } });
    render(<InstallCommand source="hero_section" />);
    fireEvent.click(screen.getByRole('button', { name: 'Copy install command' }));
    await waitFor(() =>
      expect(capture).toHaveBeenCalledWith('install_command_copied', {
        source: 'hero_section',
        command: installCommand,
      })
    );
  });
});
