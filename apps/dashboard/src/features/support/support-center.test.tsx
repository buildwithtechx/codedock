import { act, fireEvent, render, screen } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import type { SupportTicket } from '#/interfaces/support';
import * as supportService from '#/services/cloud-support';
import { SupportCenter } from './support-center';

vi.mock('#/services/cloud-support', async (importOriginal) => {
  const actual = await importOriginal<typeof import('#/services/cloud-support')>();
  return {
    ...actual,
    getTickets: vi.fn(),
    createTicket: vi.fn(),
    replyTicket: vi.fn(),
  };
});

describe('SupportCenter', () => {
  const sampleTicket: SupportTicket = {
    id: 'SUP-0001',
    subject: 'Cannot deploy Dockerfile with cache mounts',
    category: 'deployment',
    status: 'open',
    messages: [
      {
        id: '1',
        sender: 'user',
        senderName: 'Taylor',
        content: 'Cache mounts fail with permission denied.',
        createdAt: '2026-10-16',
      },
    ],
    createdAt: '2026-10-16',
    updatedAt: '2026-10-16',
  };

  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('loads and renders support tickets and opens conversation on click', async () => {
    vi.mocked(supportService.getTickets).mockResolvedValue([sampleTicket]);

    await act(async () => {
      render(<SupportCenter />);
    });

    expect(screen.getByText('Support Center')).toBeTruthy();
    expect(screen.getByText('Cannot deploy Dockerfile with cache mounts')).toBeTruthy();

    const ticketRow = screen.getByText('Cannot deploy Dockerfile with cache mounts');
    await act(async () => {
      fireEvent.click(ticketRow);
    });

    expect(screen.getByText('Cache mounts fail with permission denied.')).toBeTruthy();

    const backBtn = screen.getByRole('button', { name: /Back/ });
    await act(async () => {
      fireEvent.click(backBtn);
    });

    expect(screen.getByText('Support Tickets')).toBeTruthy();
  });

  it('opens new ticket modal and handles submission', async () => {
    vi.mocked(supportService.getTickets).mockResolvedValue([]);
    const created: SupportTicket = {
      ...sampleTicket,
      id: 'SUP-0002',
      subject: 'New billing question',
    };
    vi.mocked(supportService.createTicket).mockResolvedValue(created);

    await act(async () => {
      render(<SupportCenter />);
    });

    const newBtn = screen.getByRole('button', { name: /Open your first ticket/ });
    await act(async () => {
      fireEvent.click(newBtn);
    });

    expect(screen.getByText('Create Support Ticket')).toBeTruthy();

    const subjectInput = screen.getByPlaceholderText(/Build timeout/);
    fireEvent.change(subjectInput, { target: { value: 'New billing question' } });

    const messageInput = screen.getByPlaceholderText(/Describe the issue/);
    fireEvent.change(messageInput, { target: { value: 'How does credit rollover work?' } });

    const submitBtn = screen.getByRole('button', { name: /Submit ticket/ });
    await act(async () => {
      fireEvent.click(submitBtn);
    });

    expect(supportService.createTicket).toHaveBeenCalledWith({
      subject: 'New billing question',
      category: 'deployment',
      message: 'How does credit rollover work?',
      attachDiagnostics: false,
    });
  });
});
