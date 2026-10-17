import { fireEvent, render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import type { SupportTicket } from '#/services/cloud-support';
import { SupportTicketList } from './support-ticket-list';

describe('SupportTicketList', () => {
  const sampleTickets: SupportTicket[] = [
    {
      id: 'SUP-0001',
      subject: 'Custom base image question',
      category: 'deployment',
      status: 'open',
      messages: [
        { id: '1', sender: 'user', senderName: 'Sam', content: 'Help', createdAt: '2026-10-16' },
      ],
      createdAt: '2026-10-16',
      updatedAt: '2026-10-16',
    },
    {
      id: 'SUP-0002',
      subject: 'Invoice VAT receipt',
      category: 'billing',
      status: 'resolved',
      messages: [
        { id: '2', sender: 'user', senderName: 'Sam', content: 'Invoice', createdAt: '2026-10-16' },
      ],
      createdAt: '2026-10-16',
      updatedAt: '2026-10-16',
    },
  ];

  it('renders ticket overview with subject, status pills, and category badges', () => {
    const onSelect = vi.fn();
    const onNew = vi.fn();

    render(
      <SupportTicketList tickets={sampleTickets} onSelectTicket={onSelect} onNewTicket={onNew} />
    );

    expect(screen.getByText('Support Tickets')).toBeTruthy();
    expect(screen.getByText('Custom base image question')).toBeTruthy();
    expect(screen.getByText('Invoice VAT receipt')).toBeTruthy();
    expect(screen.getByText('open')).toBeTruthy();
    expect(screen.getByText('resolved')).toBeTruthy();
    expect(screen.getByText('Deployments')).toBeTruthy();
    expect(screen.getByText('Billing')).toBeTruthy();

    const ticketItem = screen.getByText('Custom base image question');
    fireEvent.click(ticketItem);
    expect(onSelect).toHaveBeenCalledWith('SUP-0001');

    const newBtn = screen.getByRole('button', { name: /New ticket/ });
    fireEvent.click(newBtn);
    expect(onNew).toHaveBeenCalledOnce();
  });

  it('renders empty state when no tickets exist', () => {
    const onNew = vi.fn();
    render(<SupportTicketList tickets={[]} onSelectTicket={vi.fn()} onNewTicket={onNew} />);

    expect(screen.getByText('No tickets found')).toBeTruthy();
    const ctaBtn = screen.getByRole('button', { name: /Open your first ticket/ });
    fireEvent.click(ctaBtn);
    expect(onNew).toHaveBeenCalledOnce();
  });
});
