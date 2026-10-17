import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import type { ReactNode } from 'react';
import { actionsService, type ActionRun, type ActionRunner, type ActionWorkflow } from '#/services/actions';
import { ActionsHome } from './actions-home';

vi.mock('@tanstack/react-router', () => ({
  Link: ({ children, to, ...props }: { children?: ReactNode; to?: string }) => (
    <a href={to} {...props}>{children}</a>
  ),
}));

const mockWorkflows: ActionWorkflow[] = [
  {
    id: 'wf-build',
    name: 'Build and Verify',
    description: 'Unit and end-to-end tests',
    trigger: 'push',
    yaml: 'name: CI\non: [push]\njobs:\n  test:\n    runs-on: ubuntu-latest\n',
    enabled: true,
    createdAt: '2026-10-16T10:00:00Z',
    updatedAt: '2026-10-16T10:00:00Z',
  },
];

const mockRuns: ActionRun[] = [
  {
    id: 'run-1',
    workflowId: 'wf-build',
    name: 'Build and Verify',
    number: 1,
    status: 'ready',
    triggeredBy: 'admin',
    runnerId: 'runner-1',
    logs: 'Done',
    durationSeconds: 12,
    createdAt: '2026-10-16T10:00:00Z',
    updatedAt: '2026-10-16T10:00:12Z',
  },
];

const mockRunners: ActionRunner[] = [
  {
    id: 'runner-1',
    name: 'Ubuntu 22.04 Fast',
    mode: 'container',
    image: 'catthehacker/ubuntu:act-22.04',
    status: 'online',
    labels: ['ubuntu-latest'],
    enabled: true,
    createdAt: '2026-10-16T10:00:00Z',
    updatedAt: '2026-10-16T10:00:00Z',
  },
];

describe('ActionsHome', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders actions dashboard with runs, workflows, and runners count', () => {
    render(
      <ActionsHome
        initialWorkflows={mockWorkflows}
        initialRuns={mockRuns}
        initialRunners={mockRunners}
      />
    );

    expect(screen.getByRole('heading', { name: /^Actions$/i })).toBeDefined();
    expect(screen.getByRole('button', { name: /workflows/i })).toBeDefined();
    expect(screen.getByRole('button', { name: /^runs/i })).toBeDefined();
    expect(screen.getByRole('button', { name: /^runners/i })).toBeDefined();
  });

  it('navigates to run workspace on run click', async () => {
    render(
      <ActionsHome
        initialWorkflows={mockWorkflows}
        initialRuns={mockRuns}
        initialRunners={mockRunners}
      />
    );

    const runsTabButton = screen.getByRole('button', { name: /^runs/i });
    fireEvent.click(runsTabButton);

    const runCard = await screen.findByText('Build and Verify');
    fireEvent.click(runCard);

    expect(await screen.findByText('Back to Actions Overview')).toBeDefined();
    expect(screen.getByText('#1')).toBeDefined();
  });

  it('renders link to create new workflow', () => {
    render(
      <ActionsHome
        initialWorkflows={mockWorkflows}
        initialRuns={mockRuns}
        initialRunners={mockRunners}
      />
    );

    const newLink = screen.getByRole('link', { name: /new workflow/i });
    expect(newLink.getAttribute('href')).toBe('/actions/new');
  });

  it('allows deleting a workflow', async () => {
    vi.spyOn(actionsService, 'deleteWorkflow').mockResolvedValue();
    render(
      <ActionsHome
        initialWorkflows={mockWorkflows}
        initialRuns={mockRuns}
        initialRunners={mockRunners}
      />
    );

    const deleteBtn = screen.getByRole('button', { name: /delete workflow build and verify/i });
    fireEvent.click(deleteBtn);

    await waitFor(() => {
      expect(actionsService.deleteWorkflow).toHaveBeenCalledWith('wf-build');
    });
  });

  it('allows deleting a runner', async () => {
    vi.spyOn(actionsService, 'deleteRunner').mockResolvedValue();
    render(
      <ActionsHome
        initialWorkflows={mockWorkflows}
        initialRuns={mockRuns}
        initialRunners={mockRunners}
      />
    );

    const runnersTabButton = screen.getByRole('button', { name: /^runners/i });
    fireEvent.click(runnersTabButton);

    const deleteBtn = await screen.findByRole('button', { name: /delete runner ubuntu 22\.04 fast/i });
    fireEvent.click(deleteBtn);

    await waitFor(() => {
      expect(actionsService.deleteRunner).toHaveBeenCalledWith('runner-1');
    });
  });
});
