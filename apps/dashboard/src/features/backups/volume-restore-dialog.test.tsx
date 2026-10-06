import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { VolumeRestoreDialog } from './volume-restore-dialog';

const mocks = vi.hoisted(() => ({
  get: vi.fn(),
  post: vi.fn(),
  delete: vi.fn(),
  success: vi.fn(),
}));
vi.mock('#/lib/api-client', () => ({
  apiClient: mocks,
  ApiError: class extends Error {
    status = 409;
  },
}));
vi.mock('sonner', () => ({ toast: { success: mocks.success } }));

describe('VolumeRestoreDialog', () => {
  it('interrupts the POST before server cancellation is registered', async () => {
    mocks.get.mockResolvedValue({ data: { volumeName: 'owned', timeoutSeconds: 30 } });
    let signal: AbortSignal | undefined;
    mocks.post.mockImplementation(
      (_url, _body, options) =>
        new Promise((_resolve, reject) => {
          signal = options.signal;
          signal?.addEventListener('abort', () => reject(new Error('aborted')));
        })
    );
    mocks.delete.mockResolvedValue({});
    render(<VolumeRestoreDialog recordId="record" onClose={vi.fn()} />);
    await screen.findByText(/Target:/);
    fireEvent.change(screen.getByLabelText('Type the volume name to confirm overwrite'), {
      target: { value: 'owned' },
    });
    fireEvent.click(screen.getByRole('button', { name: 'Restore volume' }));
    await waitFor(() => expect(mocks.post).toHaveBeenCalled());
    fireEvent.click(screen.getByRole('button', { name: 'Interrupt restore' }));
    expect(signal?.aborted).toBe(true);
    await screen.findByRole('alert');
    expect(mocks.success).not.toHaveBeenCalled();
    expect(mocks.delete).toHaveBeenCalledWith('/backup-records/record/volume-restore');
  });
});
