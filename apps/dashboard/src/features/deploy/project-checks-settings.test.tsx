import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { checksService, type ProjectChecksConfig } from '#/services/checks';
import { ProjectChecksSettings } from './project-checks-settings';

const mockConfig: ProjectChecksConfig = {
  projectId: 'proj-1',
  enabled: true,
  reportCommitStatus: true,
  requirePassingChecks: false,
  targetBranches: ['main', 'master'],
  contextName: 'codedock/deploy',
  updatedAt: '2026-10-16T12:00:00Z',
};

describe('ProjectChecksSettings', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders form fields with initial checks config', () => {
    render(<ProjectChecksSettings projectId="proj-1" initialConfig={mockConfig} />);

    expect(screen.getByText('GitHub Checks & Deployment Admission')).toBeDefined();
    expect(screen.getByLabelText(/report commit status to github/i)).toBeDefined();
    expect(screen.getByLabelText(/require passing checks before deploy/i)).toBeDefined();
    expect(screen.getByDisplayValue('main, master')).toBeDefined();
    expect(screen.getByDisplayValue('codedock/deploy')).toBeDefined();
  });

  it('submits updated checks configuration', async () => {
    const updatedConfig: ProjectChecksConfig = {
      ...mockConfig,
      requirePassingChecks: true,
      contextName: 'ci/codedock-check',
    };
    vi.spyOn(checksService, 'updateChecks').mockResolvedValue(updatedConfig);

    render(<ProjectChecksSettings projectId="proj-1" initialConfig={mockConfig} />);

    const requireSwitch = screen.getByLabelText(/require passing checks before deploy/i);
    fireEvent.click(requireSwitch);

    const contextInput = screen.getByLabelText(/github status context/i);
    fireEvent.change(contextInput, { target: { value: 'ci/codedock-check' } });

    const submitBtn = screen.getByRole('button', { name: /save checks configuration/i });
    fireEvent.click(submitBtn);

    await waitFor(() => {
      expect(checksService.updateChecks).toHaveBeenCalledWith(
        'proj-1',
        expect.objectContaining({
          requirePassingChecks: true,
          contextName: 'ci/codedock-check',
        })
      );
    });
  });
});
