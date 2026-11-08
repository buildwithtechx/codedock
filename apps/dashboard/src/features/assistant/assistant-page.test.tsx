import { fireEvent, render, screen } from '@testing-library/react';
import type * as React from 'react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { AssistantPage } from './assistant-page';

const mockOpenSidebar = vi.fn();

vi.mock('@tanstack/react-router', () => ({
  Link: ({ children, to }: { children: React.ReactNode; to: string }) => (
    <a href={to}>{children}</a>
  ),
}));

vi.mock('#/stores/ai-sidebar-store', () => ({
  useAiSidebarStore: (selector: (s: { openSidebar: () => void }) => unknown) =>
    selector({ openSidebar: mockOpenSidebar }),
}));

describe('AssistantPage (Hero AI Agents Workspace)', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders hero title, badges, and zero-lockin action section', () => {
    render(<AssistantPage />);

    expect(screen.getByRole('heading', { name: /Put your agents/i })).toBeDefined();
    expect(screen.getByText('Codex')).toBeDefined();
    expect(screen.getByText('Claude Code')).toBeDefined();
    expect(screen.getByText('or your own API')).toBeDefined();
    expect(screen.getByText('Run agents with zero lock-in')).toBeDefined();
    expect(screen.getByText('Open Assistant Co-pilot')).toBeDefined();
    expect(screen.getByText('Manage servers')).toBeDefined();
    expect(screen.getByText('Configure AI models ↗')).toBeDefined();
  });

  it('triggers assistant co-pilot drawer on button click', () => {
    render(<AssistantPage />);

    const openBtn = screen.getByRole('button', { name: /Open Assistant Co-pilot/i });
    fireEvent.click(openBtn);

    expect(mockOpenSidebar).toHaveBeenCalledTimes(1);
  });

  it('renders interactive preview examples with tab switching', () => {
    render(<AssistantPage />);

    expect(screen.getByText('Explore an example')).toBeDefined();
    expect(screen.getByRole('tab', { name: /Deploy/i })).toBeDefined();
    expect(screen.getByRole('tab', { name: /Investigate/i })).toBeDefined();
    expect(screen.getByRole('tab', { name: /Automate/i })).toBeDefined();

    // Default tab is Deploy
    expect(screen.getByText('Deploy my API from GitHub and connect its domain.')).toBeDefined();
    expect(screen.getByText('Your API is live.')).toBeDefined();

    // Click Investigate tab
    fireEvent.click(screen.getByRole('tab', { name: /Investigate/i }));
    expect(screen.getByText('The latest deploy failed. Find the cause and fix it.')).toBeDefined();
    expect(screen.getByText('Found the cause. Fix ready to review.')).toBeDefined();

    // Click Automate tab
    fireEvent.click(screen.getByRole('tab', { name: /Automate/i }));
    expect(
      screen.getByText('When my app goes down, investigate and send me a report.')
    ).toBeDefined();
    expect(screen.getByText('The next incident has an agent on it.')).toBeDefined();
  });
});
