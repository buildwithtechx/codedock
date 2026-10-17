import { beforeEach, describe, expect, it, vi } from 'vitest';
import type { SupportTicket } from '#/interfaces/support';
import { apiClient } from '#/lib/api-client';
import { createTicket, getTicket, getTickets, replyTicket } from './cloud-support';

vi.mock('#/lib/api-client', () => ({
  apiClient: {
    get: vi.fn(),
    post: vi.fn(),
  },
}));

describe('cloud-support service', () => {
  const mockTicket: SupportTicket = {
    id: 'SUP-0001',
    subject: 'Build failing on Railpack',
    category: 'deployment',
    status: 'open',
    messages: [
      {
        id: 'msg-1',
        sender: 'user',
        senderName: 'Alex',
        content: 'My build fails with exit code 1',
        createdAt: '2026-10-16T12:00:00Z',
      },
    ],
    createdAt: '2026-10-16T12:00:00Z',
    updatedAt: '2026-10-16T12:00:00Z',
  };

  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('fetches support tickets list', async () => {
    vi.mocked(apiClient.get).mockResolvedValueOnce({ tickets: [mockTicket] });
    const result = await getTickets();
    expect(result).toHaveLength(1);
    expect(result[0].id).toBe('SUP-0001');
    expect(apiClient.get).toHaveBeenCalledWith('/cloud/support/mine');
  });

  it('handles empty ticket list on error', async () => {
    vi.mocked(apiClient.get).mockRejectedValueOnce(new Error('Network error'));
    const result = await getTickets();
    expect(result).toEqual([]);
  });

  it('creates a new support ticket', async () => {
    vi.mocked(apiClient.post).mockResolvedValueOnce({ ticket: mockTicket });
    const result = await createTicket({
      subject: 'Build failing',
      category: 'deployment',
      message: 'Help needed',
      attachDiagnostics: true,
    });
    expect(result.id).toBe('SUP-0001');
    expect(apiClient.post).toHaveBeenCalledWith('/cloud/support/mine', {
      subject: 'Build failing',
      category: 'deployment',
      message: 'Help needed',
      attachDiagnostics: true,
    });
  });

  it('fetches a single ticket with conversation', async () => {
    vi.mocked(apiClient.get).mockResolvedValueOnce({ ticket: mockTicket });
    const result = await getTicket('SUP-0001');
    expect(result.id).toBe('SUP-0001');
    expect(apiClient.get).toHaveBeenCalledWith('/cloud/support/mine/SUP-0001');
  });

  it('replies to a support ticket', async () => {
    const updated = {
      ...mockTicket,
      messages: [
        ...mockTicket.messages,
        {
          id: 'msg-2',
          sender: 'user',
          senderName: 'Alex',
          content: 'Here are the extra logs',
          createdAt: '2026-10-16T12:05:00Z',
        },
      ],
    };
    vi.mocked(apiClient.post).mockResolvedValueOnce({ ticket: updated });
    const result = await replyTicket('SUP-0001', 'Here are the extra logs');
    expect(result.messages).toHaveLength(2);
    expect(apiClient.post).toHaveBeenCalledWith('/cloud/support/mine/SUP-0001/replies', {
      message: 'Here are the extra logs',
    });
  });
});
