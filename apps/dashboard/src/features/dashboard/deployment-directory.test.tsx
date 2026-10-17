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

vi.mock('./deployment-directory-summary', () => ({
  DeploymentSummary: () => <div data-testid="deployment-summary">Summary</div>,
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
    // Verify the animated pulse icon for Activity is rendered
    const activityIcon = container.querySelector('.animate-pulse');
    expect(activityIcon).not.toBeNull();
  });

  it('renders empty state when there are no deployments', () => {
    vi.mocked(useListByOrganization).mockReturnValue({
      data: { records: [], total: 0 },
      isLoading: false,
      isError: false,
      refetch: vi.fn(),
    } as any);

    render(<DeploymentDirectory hideHeader />);
    expect(screen.getByText('No deployments yet')).toBeDefined();
  });
});
