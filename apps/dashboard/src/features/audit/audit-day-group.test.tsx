import { fireEvent, render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { AuditDayGroup } from './audit-day-group';
import type { AuditEventRow } from './api';

const mockEvents: AuditEventRow[] = [
  {
    id: 'evt-1',
    organizationId: 'org-1',
    eventType: 'deployment.created',
    actorUserId: 'usr-1',
    actor: { id: 'usr-1', name: 'Alice Developer', email: 'alice@example.com' },
    resourceType: 'deployment',
    resourceId: 'dep-1',
    resourceName: 'web-prod',
    source: 'dashboard',
    ipAddress: '127.0.0.1',
    userAgent: 'Mozilla/5.0',
    createdAt: '2026-10-16T12:00:00Z',
  },
];

describe('AuditDayGroup', () => {
  it('renders day header and event row sentence', () => {
    const onSelect = vi.fn();
    render(
      <AuditDayGroup
        dayKey="2026-10-16"
        events={mockEvents}
        onSelectEvent={onSelect}
      />
    );

    expect(screen.getByText(/Alice Developer/i)).toBeDefined();
    expect(screen.getByText(/web-prod/i)).toBeDefined();

    const row = screen.getByRole('button');
    fireEvent.click(row);
    expect(onSelect).toHaveBeenCalledWith(mockEvents[0]);
  });
});
