import { act, fireEvent, render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import type { CloudWorkspaceSummary } from '#/interfaces/server';
import {
  ManagedServerActionFeedback,
  type ManagedServerActions,
} from './managed-server-action-feedback';

describe('ManagedServerActionFeedback', () => {
  const baseServer: CloudWorkspaceSummary = {
    id: 'ws-prod-1',
    serverId: 'srv-1',
    name: 'Production Worker',
    planTierId: 'starter',
    subscriptionStatus: 'active',
    projectCount: 0,
    state: 'running',
  };

  it('renders guided deletion with subscription management link', async () => {
    const handleRemove = vi.fn();
    const handleCancel = vi.fn();
    const actions: ManagedServerActions = {
      busy: false,
      error: null,
      remove: handleRemove,
    };

    render(
      <ManagedServerActionFeedback
        server={baseServer}
        actions={actions}
        deleting={true}
        onCancelDelete={handleCancel}
      />
    );

    expect(screen.getByText('Delete Server')).toBeTruthy();
    expect(
      screen.getByText(/If this server has an active subscription, manage or cancel it in billing/)
    ).toBeTruthy();

    const link = screen.getByRole('link', { name: /Manage subscription/ });
    expect(link.getAttribute('href')).toBe('/billing/overview?workspaceId=ws-prod-1');

    const deleteBtn = screen.getByRole('button', { name: 'Confirm deletion' });
    await act(async () => {
      fireEvent.click(deleteBtn);
    });
    expect(handleRemove).toHaveBeenCalledOnce();

    const cancelBtn = screen.getByRole('button', { name: 'Cancel' });
    fireEvent.click(cancelBtn);
    expect(handleCancel).toHaveBeenCalledOnce();
  });

  it('disables deletion when server still has active projects', () => {
    const actions: ManagedServerActions = {
      busy: false,
      remove: vi.fn(),
    };

    render(
      <ManagedServerActionFeedback
        server={{ ...baseServer, projectCount: 2 }}
        actions={actions}
        deleting={true}
        onCancelDelete={vi.fn()}
      />
    );

    expect(
      screen.getByText(/Cannot delete server while it contains 2 active project\(s\)/)
    ).toBeTruthy();

    const deleteBtn = screen.getByRole('button', { name: 'Confirm deletion' });
    expect(deleteBtn.hasAttribute('disabled')).toBe(true);
  });

  it('renders error banner when action fails', () => {
    const actions: ManagedServerActions = {
      error: 'Failed to communicate with provider API',
      remove: vi.fn(),
    };

    render(
      <ManagedServerActionFeedback
        server={baseServer}
        actions={actions}
        deleting={false}
        onCancelDelete={vi.fn()}
      />
    );

    expect(screen.getByRole('alert')).toBeTruthy();
    expect(screen.getByText('Failed to communicate with provider API')).toBeTruthy();
  });
});
