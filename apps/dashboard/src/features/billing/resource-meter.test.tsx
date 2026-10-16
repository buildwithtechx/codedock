import { render, screen } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import { ResourceMeter, ResourceMetersOverview } from './resource-meter';

describe('ResourceMeter', () => {
  it('renders progress and limit for bounded resource', () => {
    render(<ResourceMeter label="System Memory" used={4} max={16} unit="GB" />);

    expect(screen.getByText('System Memory')).toBeDefined();
    expect(screen.getByText(/4/)).toBeDefined();
    expect(screen.getByText(/Limit: 16 GB/)).toBeDefined();
    expect(screen.getByText('25% used')).toBeDefined();

    const meter = screen.getByRole('meter', { name: 'System Memory' });
    expect(meter.getAttribute('aria-valuenow')).toBe('4');
    expect(meter.getAttribute('aria-valuemax')).toBe('16');
  });

  it('renders unlimited note when max is omitted or null', () => {
    render(<ResourceMeter label="Bandwidth" used={120} max={null} unit="GB" />);

    expect(screen.getByText('Bandwidth')).toBeDefined();
    expect(screen.getByText(/Limit: Unlimited/)).toBeDefined();
    expect(screen.queryByText(/% used/)).toBeNull();
  });

  it('renders overview containing 4 resource meters', () => {
    render(
      <ResourceMetersOverview
        stats={{
          cpuUsed: 2,
          cpuLimit: 4,
          memoryUsedGb: 6,
          memoryLimitGb: 8,
          storageUsedGb: 50,
          storageLimitGb: 100,
          bandwidthUsedGb: 200,
          bandwidthLimitGb: 1000,
        }}
      />
    );

    expect(screen.getByText('Resource Utilization')).toBeDefined();
    expect(screen.getByText('Dedicated vCPU')).toBeDefined();
    expect(screen.getByText('System Memory')).toBeDefined();
    expect(screen.getByText('NVMe Storage')).toBeDefined();
    expect(screen.getByText('Network Egress')).toBeDefined();
  });
});
