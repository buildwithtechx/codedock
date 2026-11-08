import type { ProjectSecurityResult } from '@codedock/contracts';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { securityService } from '#/services/security';
import { ProjectSecurityCard } from './project-security-card';

vi.mock('#/services/security', () => ({
  securityService: {
    getAdvisories: vi.fn(),
    getProjectSecurity: vi.fn(),
    checkProject: vi.fn(),
    setProjectPolicy: vi.fn(),
    applyRemediation: vi.fn(),
  },
}));

const createWrapper = () => {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: { retry: false },
      mutations: { retry: false },
    },
  });
  return ({ children }: { children: React.ReactNode }) => (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  );
};

const mockProjectResult: ProjectSecurityResult = {
  projectId: 'proj-123',
  deploymentId: 'dep-456',
  checkedAt: '2026-11-07T12:00:00Z',
  autoRefresh: false,
  findings: [
    {
      advisoryId: 'cve-2026-node-1',
      title: 'Node.js Prototype Pollution Vulnerability',
      severity: 'critical',
      action: 'https://example.com/cve-2026-node-1',
      confidence: 'confirmed',
      serviceName: 'web-api',
      evidence: 'node 20.10.0 matches <20.18.2',
      canRefresh: true,
      runtimeUpgradeRequired: false,
      imageChangeRequired: false,
    },
  ],
};

describe('ProjectSecurityCard', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders up to date state when there are no findings', () => {
    vi.mocked(securityService.getProjectSecurity).mockResolvedValue({
      status: 'success',
      message: 'OK',
      data: {
        projectId: 'proj-safe',
        autoRefresh: false,
        findings: [],
      },
    });

    render(<ProjectSecurityCard projectId="proj-safe" />, {
      wrapper: createWrapper(),
    });

    expect(screen.getByText('Up to date')).toBeDefined();
    expect(
      screen.getByText(/No known CVEs or runtime security advisories apply to this workload./i)
    ).toBeDefined();
  });

  it('renders findings and remediation button when vulnerabilities are detected', () => {
    vi.mocked(securityService.getProjectSecurity).mockResolvedValue({
      status: 'success',
      message: 'OK',
      data: mockProjectResult,
    });

    render(<ProjectSecurityCard projectId="proj-123" initialResult={mockProjectResult} />, {
      wrapper: createWrapper(),
    });

    expect(screen.getByText('CRITICAL')).toBeDefined();
    expect(screen.getByText('Node.js Prototype Pollution Vulnerability')).toBeDefined();
    expect(screen.getByText(/node 20.10.0 matches <20.18.2/i)).toBeDefined();
    expect(screen.getByText('Refresh & Redeploy')).toBeDefined();
  });

  it('calls applyRemediation when Refresh & Redeploy is clicked', async () => {
    vi.mocked(securityService.getProjectSecurity).mockResolvedValue({
      status: 'success',
      message: 'OK',
      data: mockProjectResult,
    });
    vi.mocked(securityService.applyRemediation).mockResolvedValue({
      status: 'success',
      message: 'Remediation queued',
      data: {
        id: 'ref-1',
        projectId: 'proj-123',
        deploymentId: 'dep-456',
        advisoryId: 'cve-2026-node-1',
        status: 'queued',
        createdAt: '2026-11-07T12:00:00Z',
      },
    });

    render(<ProjectSecurityCard projectId="proj-123" initialResult={mockProjectResult} />, {
      wrapper: createWrapper(),
    });

    const refreshBtn = screen.getByText('Refresh & Redeploy');
    fireEvent.click(refreshBtn);

    await waitFor(() => {
      expect(securityService.applyRemediation).toHaveBeenCalledWith('proj-123', {
        advisoryId: 'cve-2026-node-1',
        deploymentId: 'dep-456',
        serviceId: undefined,
        runtimePatchId: undefined,
      });
    });
  });

  it('renders precision semver runtime patch transition and passes runtimePatchId on apply', async () => {
    const patchedResult: ProjectSecurityResult = {
      projectId: 'proj-123',
      deploymentId: 'dep-456',
      checkedAt: '2026-11-07T12:00:00Z',
      autoRefresh: false,
      findings: [
        {
          advisoryId: 'cve-2026-node-1',
          title: 'Node.js Prototype Pollution Vulnerability',
          severity: 'high',
          action: 'https://example.com/cve-2026-node-1',
          confidence: 'confirmed',
          serviceName: 'web-api',
          serviceId: 'svc-node-1',
          evidence: 'node 20.10.0 matches <20.18.2',
          canRefresh: true,
          runtimePatch: {
            reviewId: '111122223333444455556666777788889999aaaabbbbccccddddeeeeffff0000',
            runtime: 'node',
            currentVersion: '20.10.0',
            targetVersion: '20.18.2',
            buildImage: 'node:20.18.2-bookworm-slim',
            runtimeImage: 'node:20.18.2-bookworm-slim',
          },
        },
      ],
    };

    vi.mocked(securityService.getProjectSecurity).mockResolvedValue({
      status: 'success',
      message: 'OK',
      data: patchedResult,
    });
    vi.mocked(securityService.applyRemediation).mockResolvedValue({
      status: 'success',
      message: 'Remediation queued',
      data: {
        id: 'ref-1',
        projectId: 'proj-123',
        deploymentId: 'dep-456',
        advisoryId: 'cve-2026-node-1',
        status: 'queued',
        createdAt: '2026-11-07T12:00:00Z',
      },
    });

    render(<ProjectSecurityCard projectId="proj-123" initialResult={patchedResult} />, {
      wrapper: createWrapper(),
    });

    expect(screen.getByText('Runtime Patch:')).toBeDefined();
    expect(screen.getByText('20.10.0')).toBeDefined();
    expect(screen.getByText('20.18.2')).toBeDefined();

    const refreshBtn = screen.getByText('Refresh & Redeploy');
    fireEvent.click(refreshBtn);

    await waitFor(() => {
      expect(securityService.applyRemediation).toHaveBeenCalledWith('proj-123', {
        advisoryId: 'cve-2026-node-1',
        deploymentId: 'dep-456',
        serviceId: 'svc-node-1',
        runtimePatchId: '111122223333444455556666777788889999aaaabbbbccccddddeeeeffff0000',
      });
    });
  });

  it('calls checkProject when Check now is clicked', async () => {
    vi.mocked(securityService.getProjectSecurity).mockResolvedValue({
      status: 'success',
      message: 'OK',
      data: mockProjectResult,
    });
    vi.mocked(securityService.checkProject).mockResolvedValue({
      status: 'success',
      message: 'Security checked',
      data: mockProjectResult,
    });

    render(<ProjectSecurityCard projectId="proj-123" initialResult={mockProjectResult} />, {
      wrapper: createWrapper(),
    });

    const checkBtn = screen.getByText('Check now');
    fireEvent.click(checkBtn);

    await waitFor(() => {
      expect(securityService.checkProject).toHaveBeenCalledWith('proj-123');
    });
  });
});
