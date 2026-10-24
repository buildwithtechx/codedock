import { fireEvent, render, screen } from '@testing-library/react';
import { beforeEach, describe, expect, it } from 'vitest';
import { MailPanels } from './mail-panels';

describe('MailPanels component', () => {
  beforeEach(() => {
    localStorage.clear();
  });

  it('renders both panes with initial desktop sizes', () => {
    render(
      <MailPanels
        mailList={<div data-testid="list-content">Inbox List</div>}
        readingPane={<div data-testid="reader-content">Thread Body</div>}
        hasThread={false}
      />
    );

    expect(screen.getAllByText('Inbox List')[0]).toBeDefined();
    expect(screen.getAllByText('Thread Body')[0]).toBeDefined();

    const separator = screen.getByRole('separator', { name: /resize email list/i });
    expect(separator).toBeDefined();
    expect(separator.getAttribute('aria-valuenow')).toBe('35');
  });

  it('adjusts width using keyboard navigation on the resize handle', () => {
    render(
      <MailPanels
        mailList={<div>Inbox List</div>}
        readingPane={<div>Thread Body</div>}
        hasThread={true}
      />
    );

    const separator = screen.getByRole('separator', { name: /resize email list/i });

    fireEvent.keyDown(separator, { key: 'ArrowRight' });
    expect(separator.getAttribute('aria-valuenow')).toBe('37');
    expect(localStorage.getItem('mail-split-panel-layout')).toBe('37');

    fireEvent.keyDown(separator, { key: 'ArrowLeft' });
    expect(separator.getAttribute('aria-valuenow')).toBe('35');

    fireEvent.keyDown(separator, { key: 'Home' });
    expect(separator.getAttribute('aria-valuenow')).toBe('20');

    fireEvent.keyDown(separator, { key: 'End' });
    expect(separator.getAttribute('aria-valuenow')).toBe('60');
  });

  it('restores persisted width from localStorage', () => {
    localStorage.setItem('mail-split-panel-layout', '48');

    render(
      <MailPanels
        mailList={<div>Inbox List</div>}
        readingPane={<div>Thread Body</div>}
        hasThread={true}
      />
    );

    const separator = screen.getByRole('separator', { name: /resize email list/i });
    expect(separator.getAttribute('aria-valuenow')).toBe('48');
  });
});
