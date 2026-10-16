import { describe, expect, it } from 'vitest';
import { render } from '@testing-library/react';
import { AppGpuBadge } from './app-gpu-badge';

describe('AppGpuBadge', () => {
  it('renders nothing when gpu prop is missing or falsy', () => {
    const { container: c1 } = render(<AppGpuBadge />);
    expect(c1.firstChild).toBeNull();

    const { container: c2 } = render(<AppGpuBadge gpu={null} />);
    expect(c2.firstChild).toBeNull();

    const { container: c3 } = render(<AppGpuBadge gpu="" />);
    expect(c3.firstChild).toBeNull();
  });

  it('renders GPU Required badge when required', () => {
    const { container } = render(<AppGpuBadge gpu="required" />);
    expect(container.textContent).toContain('GPU Required');
    const badge = container.querySelector('span');
    expect(badge).toBeTruthy();
    expect(badge?.className).toContain('text-amber-600');
  });

  it('renders GPU Recommended badge when recommended', () => {
    const { container } = render(<AppGpuBadge gpu="recommended" />);
    expect(container.textContent).toContain('GPU Recommended');
    const badge = container.querySelector('span');
    expect(badge).toBeTruthy();
    expect(badge?.className).toContain('text-muted-foreground');
  });
});
