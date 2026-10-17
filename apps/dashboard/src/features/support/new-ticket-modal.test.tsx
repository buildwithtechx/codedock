import { act, fireEvent, render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { NewTicketModal } from './new-ticket-modal';

describe('NewTicketModal', () => {
  it('renders ticket form inputs and submits payload', async () => {
    const handleSubmit = vi.fn().mockResolvedValue(undefined);
    const handleClose = vi.fn();

    render(<NewTicketModal onClose={handleClose} onSubmit={handleSubmit} />);

    expect(screen.getByText('Create Support Ticket')).toBeTruthy();

    const subjectInput = screen.getByPlaceholderText(/Build timeout/);
    fireEvent.change(subjectInput, { target: { value: 'PostgreSQL connection timeout' } });

    const messageInput = screen.getByPlaceholderText(/Describe the issue/);
    fireEvent.change(messageInput, {
      target: { value: 'Database logs show 100% pool utilization' },
    });

    const checkbox = screen.getByRole('checkbox');
    fireEvent.click(checkbox);

    const submitBtn = screen.getByRole('button', { name: /Submit ticket/ });
    await act(async () => {
      fireEvent.click(submitBtn);
    });

    expect(handleSubmit).toHaveBeenCalledWith({
      subject: 'PostgreSQL connection timeout',
      category: 'deployment',
      message: 'Database logs show 100% pool utilization',
      attachDiagnostics: true,
    });
    expect(handleClose).toHaveBeenCalledOnce();
  });

  it('closes on cancel button click', () => {
    const handleClose = vi.fn();
    render(<NewTicketModal onClose={handleClose} onSubmit={vi.fn()} />);

    const cancelBtn = screen.getByRole('button', { name: 'Cancel' });
    fireEvent.click(cancelBtn);
    expect(handleClose).toHaveBeenCalledOnce();
  });
});
