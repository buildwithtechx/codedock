import type { AdvisoriesResponse } from '@codedock/contracts';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { securityService } from '#/services/security';
import { SecurityBanner } from './security-banner';

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

const mockAdvisoriesData: AdvisoriesResponse = {
  revision: 12,
  checkedAt: '2026-11-07T12:00:00Z',
  stale: false,
  advisories: [
    {
      id: 'cve-2026-node-1',
      revision: 12,
      title: 'Critical Node.js Prototype Pollution',
      severity: 'critical',
      message: 'A critical prototype pollution flaw affects Node.js runtime.',
      url: 'https://example.com/cve-2026-node-1',
      action: 'https://example.com/cve-2026-node-1',
    },
  ],
  affectedProjects: [
    {
      projectId: 'proj-123',
      deploymentId: 'dep-456',
      checkedAt: '2026-11-07T12:00:00Z',
      autoRefresh: false,
      findings: [
        {
          advisoryId: 'cve-2026-node-1',
          title: 'Critical Node.js Prototype Pollution',
          severity: 'critical',
          action: 'https://example.com/cve-2026-node-1',
          confidence: 'confirmed',
          serviceName: 'web',
          evidence: 'node 20.10.0',
          canRefresh: true,
        },
      ],
    },
  ],
};

describe('SecurityBanner', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    localStorage.clear();
  });

  it('renders nothing when there are no affected projects', () => {
    vi.mocked(securityService.getAdvisories).mockResolvedValue({
      status: 'success',
      message: 'OK',
      data: {
        revision: 1,
        stale: false,
        advisories: [],
        affectedProjects: [],
      },
    });

    const { container } = render(<SecurityBanner />, {
      wrapper: createWrapper(),
    });

    expect(container.firstChild).toBeNull();
  });

  it('renders banner with advisory title and affected workload count', async () => {
    vi.mocked(securityService.getAdvisories).mockResolvedValue({
      status: 'success',
      message: 'OK',
      data: mockAdvisoriesData,
    });

    render(<SecurityBanner />, {
      wrapper: createWrapper(),
    });

    await waitFor(() => {
      expect(screen.getByText('Critical Node.js Prototype Pollution')).toBeDefined();
      expect(screen.getByText(/1 workload affected by signed security advisories/i)).toBeDefined();
      expect(screen.getByText('Review affected apps')).toBeDefined();
    });
  });

  it('opens review dialog when Review affected apps is clicked', async () => {
    vi.mocked(securityService.getAdvisories).mockResolvedValue({
      status: 'success',
      message: 'OK',
      data: mockAdvisoriesData,
    });
    vi.mocked(securityService.getProjectSecurity).mockResolvedValue({
      status: 'success',
      message: 'OK',
      data: mockAdvisoriesData.affectedProjects[0],
    });

    render(<SecurityBanner />, {
      wrapper: createWrapper(),
    });

    await waitFor(() => {
      expect(screen.getByText('Review affected apps')).toBeDefined();
    });

    fireEvent.click(screen.getByText('Review affected apps'));

    await waitFor(() => {
      expect(screen.getByText('Runtime Security Advisories & Audit')).toBeDefined();
      expect(screen.getByText(/Feed Revision: #12/i)).toBeDefined();
    });
  });

  it('dismisses banner when dismiss button is clicked', async () => {
    vi.mocked(securityService.getAdvisories).mockResolvedValue({
      status: 'success',
      message: 'OK',
      data: mockAdvisoriesData,
    });

    render(<SecurityBanner />, {
      wrapper: createWrapper(),
    });

    await waitFor(() => {
      expect(screen.getByText('Critical Node.js Prototype Pollution')).toBeDefined();
    });

    const dismissBtn = screen.getByLabelText('Dismiss security notice');
    fireEvent.click(dismissBtn);

    await waitFor(() => {
      expect(screen.queryByText('Critical Node.js Prototype Pollution')).toBeNull();
    });
  });
});
