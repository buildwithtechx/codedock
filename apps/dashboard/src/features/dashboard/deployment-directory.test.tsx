import { render, screen } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { DeploymentDirectory } from './deployment-directory';

vi.mock('#/hooks/use-deployments', () => ({
  useListByOrganization: vi.fn(),
}));

vi.mock('#/features/projects', () => ({
  useListProjects: vi.fn(),
}));

vi.mock('@tanstack/react-router', () => ({
  Link: ({ children, to, ...props }: any) => <a href={to} {...props}>{children}</a>,
}));

import { useListByOrganization } from '#/hooks/use-deployments';
import { useListProjects } from '#/features/projects';

describe('DeploymentDirectory', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(useListProjects).mockReturnValue({
      data: { records: [{ id: 'p1', name: 'Project 1' }], total: 1 },
    } as any);
  });

  it('renders loading state with Activity icon without crashing', () => {
    vi.mocked(useListByOrganization).mockReturnValue({
      data: undefined,
      isLoading: true,
      isError: false,
      refetch: vi.fn(),
    } as any);

    const { container } = render(<DeploymentDirectory hideHeader />);
    const activityIcon = container.querySelector('.animate-pulse');
    expect(activityIcon).not.toBeNull();
  });

  it('renders full overview, filters, and empty state matching OpenShip design', () => {
    vi.mocked(useListByOrganization).mockReturnValue({
      data: { records: [], total: 0 },
      isLoading: false,
      isError: false,
      refetch: vi.fn(),
    } as any);

    render(<DeploymentDirectory hideHeader />);

    // Verify Filters
    expect(screen.getByPlaceholderText('Search deployments...')).toBeDefined();
    expect(screen.getByText('All Projects')).toBeDefined();
    expect(screen.getByRole('button', { name: 'All' })).toBeDefined();
    expect(screen.getByRole('button', { name: 'Success' })).toBeDefined();
    expect(screen.getByRole('button', { name: 'Failed' })).toBeDefined();
    expect(screen.getByRole('button', { name: 'Building' })).toBeDefined();
    expect(screen.getByRole('button', { name: 'In Progress' })).toBeDefined();
    expect(screen.getByRole('button', { name: 'Canceled' })).toBeDefined();

    // Verify Deployment Overview & Metrics
    expect(screen.getByText('Deployment overview')).toBeDefined();
    expect(screen.getByText('Deployment attempts')).toBeDefined();
    expect(screen.getByText('Success rate')).toBeDefined();
    expect(screen.getByText('Across the selected window')).toBeDefined();

    // Verify Daily outcomes & Project breakdown cards
    expect(screen.getByText('Daily outcomes')).toBeDefined();
    expect(screen.getByText('Project breakdown')).toBeDefined();

    // Verify History heading and Empty State
    expect(screen.getByText('Deployment history')).toBeDefined();
    expect(screen.getByText('No deployments yet')).toBeDefined();
    expect(screen.getByText('Deploy Project')).toBeDefined();
    expect(screen.getByText('Browse Templates')).toBeDefined();
  });
});
