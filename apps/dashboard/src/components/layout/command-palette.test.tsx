import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { fireEvent, render, screen } from '@testing-library/react';
import type * as React from 'react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { CommandPalette } from './command-palette';

const mockNavigate = vi.fn();

vi.mock('@tanstack/react-router', () => ({
  useNavigate: () => mockNavigate,
}));

vi.mock('#/hooks/use-projects', () => ({
  useListProjects: () => ({
    data: {
      data: [{ id: 'proj-1', name: 'Alpha API', description: 'Main backend' }],
    },
  }),
}));

vi.mock('#/hooks/use-apps', () => ({
  useListAppsByOrganization: () => ({
    data: [{ id: 'app-1', name: 'Web Client', projectId: 'proj-1' }],
  }),
}));

vi.mock('#/hooks/use-servers', () => ({
  useListServers: () => ({
    data: [{ id: 'srv-1', name: 'Production Node 1', ip: '192.168.1.10' }],
  }),
}));

vi.mock('#/hooks/use-databases', () => ({
  useListDatabases: () => ({
    data: [{ id: 'db-1', name: 'Primary Postgres', type: 'postgres', projectId: 'proj-1' }],
  }),
}));

function renderWithClient(ui: React.ReactElement) {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  });
  return render(<QueryClientProvider client={queryClient}>{ui}</QueryClientProvider>);
}

describe('CommandPalette', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    window.HTMLElement.prototype.scrollIntoView = vi.fn();
  });

  it('renders commands, projects, apps, and servers when open', () => {
    renderWithClient(<CommandPalette open={true} onOpenChange={vi.fn()} />);

    expect(screen.getByPlaceholderText(/Type a command or search projects/i)).toBeTruthy();
    expect(screen.getByText('Toggle Theme')).toBeTruthy();
    expect(screen.getByText('New Project')).toBeTruthy();
    expect(screen.getByText('Alpha API')).toBeTruthy();
    expect(screen.getByText('Web Client')).toBeTruthy();
    expect(screen.getByText('Production Node 1')).toBeTruthy();
    expect(screen.getByText('Primary Postgres')).toBeTruthy();
  });

  it('triggers onOpenChange on Cmd+K keyboard shortcut', () => {
    const handleOpenChange = vi.fn();
    renderWithClient(<CommandPalette open={false} onOpenChange={handleOpenChange} />);

    fireEvent.keyDown(window, { key: 'k', metaKey: true });
    expect(handleOpenChange).toHaveBeenCalledWith(true);
  });
});
