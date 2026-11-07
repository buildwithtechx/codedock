import { render, screen } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import { ComposeChangeSummary } from './compose-change-summary';

describe('ComposeChangeSummary', () => {
  it('renders nothing when changes are empty or undefined', () => {
    const { container: c1 } = render(<ComposeChangeSummary />);
    expect(c1.firstChild).toBeNull();

    const { container: c2 } = render(
      <ComposeChangeSummary changes={{ added: [], changed: [], removed: [] }} />
    );
    expect(c2.firstChild).toBeNull();
  });

  it('renders added, changed, and removed badges with service names', () => {
    render(
      <ComposeChangeSummary
        changes={{
          added: ['auth-service'],
          changed: ['api-gateway'],
          removed: ['legacy-worker'],
        }}
      />
    );

    expect(screen.getByText('Detected Changes')).toBeDefined();
    expect(screen.getByText('+auth-service')).toBeDefined();
    expect(screen.getByText('~api-gateway')).toBeDefined();
    expect(screen.getByText('-legacy-worker')).toBeDefined();
    expect(screen.getByText('legacy-worker')).toBeDefined();
    expect(
      screen.getByText(
        /will be stopped, unrouted from edge proxy, and orphaned containers safely pruned/i
      )
    ).toBeDefined();
  });
});
