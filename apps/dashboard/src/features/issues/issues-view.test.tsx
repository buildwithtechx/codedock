import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import type { ReactNode } from 'react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { useOrganizationStore } from '#/stores/organization-store';
import type { AttentionIssue } from './interfaces';
import { IssuesView } from './issues-view';

const mocks = vi.hoisted(() => ({ get: vi.fn(), delete: vi.fn(), post: vi.fn() }));
vi.mock('#/lib/api-client', () => ({
  apiClient: mocks,
  getApiBaseUrl: () => 'http://test',
}));

vi.mock('@tanstack/react-router', () => ({
  Link: ({ children }: { children?: ReactNode }) => <span>{children}</span>,
}));

const resolvedIssue: AttentionIssue = {
  id: 'iss-1',
  organizationId: 'org-1',
  kind: 'deployment-failed',
  subject: 'svc-1',
  severity: 'warning',
  status: 'resolved',
  title: 'Deploy failing for web',
  detail: 'latest deployment failed',
  remediation: 'redeploy',
  occurrences: 1,
  firstSeen: '2026-01-01T00:00:00Z',
  lastSeen: '2026-01-01T00:00:00Z',
  updatedAt: '2026-01-01T00:00:00Z',
};

function renderView() {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  const view = render(
    <QueryClientProvider client={client}>
      <IssuesView />
    </QueryClientProvider>
  );
  client.clear();
  return view;
}

async function openHistoryTab() {
  fireEvent.mouseDown(await screen.findByRole('tab', { name: 'History' }));
}

beforeEach(() => {
  vi.clearAllMocks();
  useOrganizationStore.setState({ activeOrganizationId: 'org-1' });
});

describe('monitoring clear history', () => {
  it('offers the clear control only when history is clearable', async () => {
    mocks.get.mockResolvedValue({ data: { issues: [resolvedIssue], canClearHistory: true } });
    renderView();
    await openHistoryTab();
    expect(await screen.findByRole('button', { name: /clear history/i })).toBeTruthy();
  });

  it('hides the clear control without admin permission', async () => {
    mocks.get.mockResolvedValue({ data: { issues: [resolvedIssue], canClearHistory: false } });
    renderView();
    await openHistoryTab();
    await screen.findByText('Deploy failing for web');
    expect(screen.queryByRole('button', { name: /clear history/i })).toBeNull();
  });

  it('requires an inline confirmation before clearing', async () => {
    mocks.get.mockResolvedValue({ data: { issues: [resolvedIssue], canClearHistory: true } });
    mocks.delete.mockResolvedValue({ data: { deleted: 1 }, message: 'Monitoring history cleared' });
    renderView();
    await openHistoryTab();

    fireEvent.click(await screen.findByRole('button', { name: /clear history/i }));
    expect(mocks.delete).not.toHaveBeenCalled();
    expect(
      await screen.findByText(/Clear resolved monitoring history for this organization\?/)
    ).toBeTruthy();

    const confirm = screen.getAllByRole('button', { name: /^clear history$/i }).at(-1);
    expect(confirm).toBeTruthy();
    fireEvent.click(confirm as HTMLElement);

    await waitFor(() =>
      expect(mocks.delete).toHaveBeenCalledWith('/organizations/org-1/attention/history')
    );
  });
});
