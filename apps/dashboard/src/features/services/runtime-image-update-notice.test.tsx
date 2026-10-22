import { fireEvent, render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import type { RuntimeImageCheckResult } from '@codedock/contracts';
import { RuntimeImageUpdateNotice } from './runtime-image-update-notice';

describe('RuntimeImageUpdateNotice', () => {
  it('renders nothing when no updates are available', () => {
    const result: RuntimeImageCheckResult = {
      checkedAt: new Date().toISOString(),
      images: [
        {
          serviceId: 'svc-1',
          serviceName: 'web',
          role: 'image',
          ref: 'node:20-alpine',
          currentDigest: 'sha256:1111',
          latestDigest: 'sha256:1111',
          status: 'current',
        },
      ],
    };

    const { container } = render(
      <RuntimeImageUpdateNotice result={result} busy={false} onRefresh={vi.fn()} />
    );
    expect(container.firstChild).toBeNull();
  });

  it('renders notice and calls onRefresh when update is available', () => {
    const onRefresh = vi.fn();
    const result: RuntimeImageCheckResult = {
      checkedAt: new Date().toISOString(),
      images: [
        {
          serviceId: 'svc-1',
          serviceName: 'web',
          role: 'base',
          ref: 'node:20-alpine',
          currentDigest: 'sha256:11111111111111111111',
          latestDigest: 'sha256:22222222222222222222',
          status: 'update_available',
        },
      ],
    };

    render(<RuntimeImageUpdateNotice result={result} busy={false} onRefresh={onRefresh} />);

    expect(screen.getByText('A newer base image is available with security patches')).toBeTruthy();
    expect(screen.getByText(/web:/)).toBeTruthy();
    expect(screen.getByText(/node:20-alpine/)).toBeTruthy();

    const button = screen.getByRole('button', { name: /Redeploy with latest image/i });
    fireEvent.click(button);
    expect(onRefresh).toHaveBeenCalledTimes(1);
  });
});
