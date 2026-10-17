import { fireEvent, render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import type { AuditEventRow } from './api';
import { AuditDayGroup } from './audit-day-group';

const mockEvents: AuditEventRow[] = [
  {
    id: 'evt-1',
    organizationId: 'org-1',
    userId: 'usr-1',
    action: 'deployment.created',
    eventType: 'deployment.created',
    actorUserId: 'usr-1',
    actor: { id: 'usr-1', name: 'Alice Developer', email: 'alice@example.com' },
    resource: 'deployment',
    resourceType: 'deployment',
    resourceId: 'dep-1',
    resourceName: 'web-prod',
    details: 'Created deployment web-prod',
    source: 'dashboard',
    ipAddress: '127.0.0.1',
    userAgent: 'Mozilla/5.0',
    createdAt: '2026-10-16T12:00:00Z',
  },
];

describe('AuditDayGroup', () => {
  it('renders day header and event row sentence', () => {
    const onSelect = vi.fn();
    render(<AuditDayGroup dayKey="2026-10-16" events={mockEvents} onSelectEvent={onSelect} />);

    expect(screen.getByText(/Alice Developer/i)).toBeDefined();
    expect(screen.getByText(/web-prod/i)).toBeDefined();

    const row = screen.getByRole('button');
    fireEvent.click(row);
    expect(onSelect).toHaveBeenCalledWith(mockEvents[0]);
  });

  it('never displays raw UUID in title and falls back to clean noun', () => {
    const rawUuidEvent: AuditEventRow = {
      id: 'evt-2',
      action: 'project.create',
      userId: 'usr-solomon',
      actorName: 'Solomon Olatunji',
      resource: '91fdae7e-559c-41af-8426-814fe0d5e643',
      details: '',
      ipAddress: '127.0.0.1',
      createdAt: '2026-10-16T12:00:00Z',
    };

    render(<AuditDayGroup dayKey="2026-10-16" events={[rawUuidEvent]} onSelectEvent={vi.fn()} />);

    expect(screen.queryByText(/91fdae7e-559c-41af-8426-814fe0d5e643/i)).toBeNull();
    expect(screen.getByText('Solomon Olatunji created the project')).toBeDefined();
  });
});
