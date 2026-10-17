import { act, fireEvent, render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import type { SupportTicket } from '#/services/cloud-support';
import { SupportConversation } from './support-conversation';

describe('SupportConversation', () => {
  const sampleTicket: SupportTicket = {
    id: 'SUP-0001',
    subject: 'Cannot connect to managed database',
    category: 'deployment',
    status: 'open',
    hasDiagnosticBundle: true,
    messages: [
      {
        id: '1',
        sender: 'user',
        senderName: 'Morgan',
        content: 'Port 5432 is timing out on my server.',
        createdAt: '2026-10-16 10:00',
      },
      {
        id: '2',
        sender: 'staff',
        senderName: 'CodeDock Support',
        content: 'Please ensure host firewall allows inbound TCP on 5432.',
        createdAt: '2026-10-16 10:15',
      },
    ],
    createdAt: '2026-10-16',
    updatedAt: '2026-10-16',
  };

  it('renders messages with user and staff attribution', () => {
    render(<SupportConversation ticket={sampleTicket} onBack={vi.fn()} onSendReply={vi.fn()} />);

    expect(screen.getByText('SUP-0001')).toBeTruthy();
    expect(screen.getByText('Cannot connect to managed database')).toBeTruthy();
    expect(screen.getByText('Diagnostics attached')).toBeTruthy();
    expect(screen.getByText('Morgan')).toBeTruthy();
    expect(screen.getByText('CodeDock Staff')).toBeTruthy();
    expect(screen.getByText('Port 5432 is timing out on my server.')).toBeTruthy();
    expect(
      screen.getByText('Please ensure host firewall allows inbound TCP on 5432.')
    ).toBeTruthy();
  });

  it('handles sending replies', async () => {
    const handleSend = vi.fn().mockResolvedValue(undefined);
    render(<SupportConversation ticket={sampleTicket} onBack={vi.fn()} onSendReply={handleSend} />);

    const textarea = screen.getByPlaceholderText('Type your response or update...');
    fireEvent.change(textarea, { target: { value: 'Firewall is now opened and confirmed.' } });

    const sendBtn = screen.getByRole('button', { name: /Send response/ });
    await act(async () => {
      fireEvent.click(sendBtn);
    });

    expect(handleSend).toHaveBeenCalledWith('Firewall is now opened and confirmed.');
  });

  it('hides reply box when ticket is closed', () => {
    render(
      <SupportConversation
        ticket={{ ...sampleTicket, status: 'closed' }}
        onBack={vi.fn()}
        onSendReply={vi.fn()}
      />
    );

    expect(screen.getByText(/This ticket is marked as closed/)).toBeTruthy();
    expect(screen.queryByPlaceholderText('Type your response or update...')).toBeNull();
  });
});
