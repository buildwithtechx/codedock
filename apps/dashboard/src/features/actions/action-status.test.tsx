import { render, screen } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import { ActionStatus } from './action-status';

describe('ActionStatus', () => {
  it('renders ready status badge', () => {
    render(<ActionStatus status="ready" />);
    expect(screen.getByText('ready')).toBeDefined();
  });

  it('renders running status with loader', () => {
    render(<ActionStatus status="running" />);
    expect(screen.getByText('running')).toBeDefined();
  });

  it('renders failed status variant', () => {
    render(<ActionStatus status="failed" />);
    expect(screen.getByText('failed')).toBeDefined();
  });

  it('renders queued status variant', () => {
    render(<ActionStatus status="queued" />);
    expect(screen.getByText('queued')).toBeDefined();
  });
});
