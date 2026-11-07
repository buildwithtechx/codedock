import { fireEvent, render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { DeploymentResourceNotice } from './deployment-resource-notice';

describe('DeploymentResourceNotice', () => {
  it('renders null when issue is missing', () => {
    const { container } = render(<DeploymentResourceNotice issue={null} />);
    expect(container.firstChild).toBeNull();
  });

  it('renders container-oom guidance', () => {
    render(<DeploymentResourceNotice issue="container-oom" />);
    expect(screen.getByText('Container ran out of memory')).toBeTruthy();
    expect(screen.getByText(/Docker confirmed an out-of-memory kill/)).toBeTruthy();
  });

  it('renders process-heap guidance', () => {
    render(<DeploymentResourceNotice issue="process-heap" />);
    expect(screen.getByText('Application heap limit reached')).toBeTruthy();
    expect(screen.getByText(/--max-old-space-size/)).toBeTruthy();
  });

  it('renders storage-full guidance', () => {
    render(<DeploymentResourceNotice issue="storage-full" />);
    expect(screen.getByText('Storage allocation failed')).toBeTruthy();
    expect(screen.getByText(/ran out of disk space or inodes/)).toBeTruthy();
  });

  it('triggers onViewLogs callback when clicked', () => {
    const onViewLogs = vi.fn();
    render(<DeploymentResourceNotice issue="process-killed" onViewLogs={onViewLogs} />);
    expect(screen.getByText('Build process was terminated')).toBeTruthy();

    const button = screen.getByRole('button', { name: 'View Logs' });
    fireEvent.click(button);
    expect(onViewLogs).toHaveBeenCalledTimes(1);
  });
});
