import { describe, expect, it, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { DeploymentHistoryCard } from './deployment-history-card';
import type { OrganizationDeployment } from '#/interfaces/services';

vi.mock('@tanstack/react-router', () => ({
  Link: ({ children, ...props }: any) => <a {...props}>{children}</a>,
}));

const mockDeployment: OrganizationDeployment = {
  id: 'dep-fail-1',
  serviceId: 'svc-1',
  serviceName: 'web',
  environmentId: 'env-1',
  projectId: 'proj-1',
  projectName: 'My Project',
  status: 'FAILED',
  branch: 'main',
  commitHash: 'commit1234567',
  commitMessage: 'Failed cutover',
  trigger: 'manual',
  containerId: '',
  resourceIssue: '',
  artifactRetainedAt: null,
  pinned: false,
  createdAt: new Date().toISOString(),
  updatedAt: new Date().toISOString(),
};

describe('DeploymentHistoryCard Retry', () => {
  it('displays Retry deployment in action menu for failed deployments', async () => {
    const queryClient = new QueryClient({
      defaultOptions: { queries: { retry: false } },
    });

    render(
      <QueryClientProvider client={queryClient}>
        <DeploymentHistoryCard deployment={mockDeployment} />
      </QueryClientProvider>
    );

    const menuButton = screen.getByRole('button', { name: /actions for web/i });
    expect(menuButton).not.toBeNull();

    fireEvent.pointerDown(menuButton);

    const retryOption = await screen.findByRole('menuitem', { name: /retry deployment/i });
    expect(retryOption).not.toBeNull();
  });
});
