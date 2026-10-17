import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { type BuildAdmissionStatus, checksService } from '#/services/checks';
import { DeploymentAdmissionCard } from './deployment-admission-card';

const mockAdmission: BuildAdmissionStatus = {
  canAdmit: true,
  activeBuilds: 2,
  maxConcurrentBuilds: 4,
  queuePosition: 0,
  estimatedWaitSeconds: 0,
  reason: 'Capacity available for immediate execution',
};

describe('DeploymentAdmissionCard', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders admission card with active slot usage and ready badge', () => {
    render(<DeploymentAdmissionCard projectId="proj-1" initialAdmission={mockAdmission} />);

    expect(screen.getByText('Build Admission & Concurrency')).toBeDefined();
    expect(screen.getByText('Admission Ready')).toBeDefined();
    expect(screen.getByText(/2 of 4 active/i)).toBeDefined();
    expect(screen.getByText('Capacity available for immediate execution')).toBeDefined();
  });

  it('renders queued badge and wait estimate when capacity full', () => {
    const queuedAdmission: BuildAdmissionStatus = {
      canAdmit: false,
      activeBuilds: 4,
      maxConcurrentBuilds: 4,
      queuePosition: 2,
      estimatedWaitSeconds: 45,
      reason: 'Maximum concurrent builds reached; run queued',
    };

    render(<DeploymentAdmissionCard projectId="proj-1" initialAdmission={queuedAdmission} />);

    expect(screen.getByText('Queued (#2)')).toBeDefined();
    expect(screen.getByText('~45s wait')).toBeDefined();
    expect(screen.getByText(/4 of 4 active/i)).toBeDefined();
  });

  it('refreshes admission data when refresh button is clicked', async () => {
    const updatedAdmission: BuildAdmissionStatus = {
      ...mockAdmission,
      activeBuilds: 3,
    };
    vi.spyOn(checksService, 'getAdmission').mockResolvedValue(updatedAdmission);

    render(<DeploymentAdmissionCard projectId="proj-1" initialAdmission={mockAdmission} />);

    const refreshBtn = screen.getByRole('button', { name: /refresh admission/i });
    fireEvent.click(refreshBtn);

    await waitFor(() => {
      expect(checksService.getAdmission).toHaveBeenCalledWith('proj-1');
      expect(screen.getByText(/3 of 4 active/i)).toBeDefined();
    });
  });
});
