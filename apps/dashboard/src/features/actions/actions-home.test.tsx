import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { actionsService, type ActionRun, type ActionRunner, type ActionWorkflow } from '#/services/actions';
import { ActionsHome } from './actions-home';

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
    expect(screen.getByText(/Recent Runs \(1\)/i)).toBeDefined();
    expect(screen.getByText(/Workflows \(1\)/i)).toBeDefined();
    expect(screen.getByText(/Runners \(1\)/i)).toBeDefined();
  });

  it('navigates to run workspace on run click', async () => {
    render(
      <ActionsHome
        initialWorkflows={mockWorkflows}
        initialRuns={mockRuns}
        initialRunners={mockRunners}
      />
    );

    const runCard = screen.getByText('Build and Verify');
    fireEvent.click(runCard);

    expect(await screen.findByText('Back to Actions Overview')).toBeDefined();
    expect(screen.getByText('#1')).toBeDefined();
  });

  it('opens create workflow modal and creates workflow', async () => {
    const newWorkflow: ActionWorkflow = {
      id: 'wf-new',
      name: 'Deploy Production',
      description: 'Ship to cluster',
      trigger: 'push',
      yaml: 'name: Deploy',
      enabled: true,
      createdAt: '2026-10-16T12:00:00Z',
      updatedAt: '2026-10-16T12:00:00Z',
    };
    vi.spyOn(actionsService, 'createWorkflow').mockResolvedValue(newWorkflow);

    render(
      <ActionsHome
        initialWorkflows={mockWorkflows}
        initialRuns={mockRuns}
        initialRunners={mockRunners}
      />
    );

    const newBtn = screen.getByRole('button', { name: /new workflow/i });
    fireEvent.click(newBtn);

    expect(await screen.findByText('Create Action Workflow')).toBeDefined();

    const nameInput = screen.getByLabelText(/workflow name/i);
    fireEvent.change(nameInput, { target: { value: 'Deploy Production' } });

    const submitBtn = screen.getByRole('button', { name: /save workflow/i });
    fireEvent.click(submitBtn);

    await waitFor(() => {
      expect(actionsService.createWorkflow).toHaveBeenCalled();
    });
  });
});
