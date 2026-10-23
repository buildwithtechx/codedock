import { fireEvent, render, screen } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { BottomDock } from './bottom-dock';
import { useDockStore } from '#/stores/dock-store';

vi.mock('#/features/services/live-logs-viewer', () => ({
  LiveLogsViewer: () => <div data-testid="live-logs-viewer">Live Logs Mock</div>,
}));

vi.mock('#/features/services/web-terminal', () => ({
  WebTerminal: () => <div data-testid="web-terminal">Web Terminal Mock</div>,
}));

describe('BottomDock', () => {
  beforeEach(() => {
    useDockStore.setState({
      isOpen: false,
      activeTab: 'logs',
      serviceId: undefined,
      height: 280,
    });
  });

  it('renders collapsed pill when closed and opens on click', () => {
    render(<BottomDock />);

    const openBtn = screen.getByRole('button', { name: /Logs & Terminal/i });
    expect(openBtn).toBeTruthy();

    fireEvent.click(openBtn);
    expect(useDockStore.getState().isOpen).toBe(true);
  });

  it('toggles open/close on Cmd+J key press', () => {
    render(<BottomDock />);
    expect(useDockStore.getState().isOpen).toBe(false);

    fireEvent.keyDown(window, { key: 'j', metaKey: true });
    expect(useDockStore.getState().isOpen).toBe(true);

    fireEvent.keyDown(window, { key: 'j', metaKey: true });
    expect(useDockStore.getState().isOpen).toBe(false);
  });

  it('renders tabs and empty states when open without serviceId', () => {
    useDockStore.setState({ isOpen: true, activeTab: 'logs' });
    render(<BottomDock />);

    expect(screen.getByRole('button', { name: /Logs/i })).toBeTruthy();
    expect(screen.getByRole('button', { name: /Terminal/i })).toBeTruthy();
    expect(screen.getByText(/No active service selected for live logs/i)).toBeTruthy();

    const terminalTab = screen.getByRole('button', { name: /Terminal/i });
    fireEvent.click(terminalTab);
    expect(useDockStore.getState().activeTab).toBe('terminal');
  });

  it('renders LiveLogsViewer when serviceId is set', async () => {
    useDockStore.setState({ isOpen: true, activeTab: 'logs', serviceId: 'srv-123' });
    render(<BottomDock />);

    const viewer = await screen.findByTestId('live-logs-viewer');
    expect(viewer).toBeTruthy();
  });
});
