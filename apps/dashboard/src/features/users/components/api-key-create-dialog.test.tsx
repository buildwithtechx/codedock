import { act, fireEvent, render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { ApiKeyCreateDialog } from './api-key-create-dialog';

const mutation = vi.hoisted(() => ({ mutate: vi.fn() }));
vi.mock('#/features/profile', () => ({
  useCreateToken: () => ({ mutate: mutation.mutate, isPending: false }),
}));
vi.mock('#/hooks/use-canvas', () => ({
  useListCanvasSummaries: () => ({
    data: { data: [{ id: 'project', name: 'Project' }] },
    isLoading: false,
    isError: false,
  }),
}));
vi.mock('sonner', () => ({ toast: { error: vi.fn() } }));

describe('ApiKeyCreateDialog', () => {
  it('submits selected projects with no expiry and displays the returned credential', () => {
    const onSuccess = vi.fn();
    render(<ApiKeyCreateDialog open onOpenChange={vi.fn()} onSuccess={onSuccess} />);
    fireEvent.change(screen.getByPlaceholderText('Production deploys'), {
      target: { value: 'automation' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'SPECIFIC PROJECTS' }));
    fireEvent.click(screen.getByRole('checkbox', { name: 'Project' }));
    fireEvent.click(screen.getByRole('button', { name: 'NO EXPIRATION' }));
    fireEvent.click(screen.getByRole('button', { name: 'Create Key' }));
    const [request, callbacks] = mutation.mutate.mock.calls[0];
    expect(request.payload).toEqual({
      name: 'automation',
      accessLevel: 'read',
      projectScope: 'specific',
      allowedProjects: ['project'],
    });
    act(() => callbacks.onSuccess({ data: { plain: 'vpt_credential' } }));
    expect(onSuccess).toHaveBeenCalledWith('vpt_credential');
  });
});
