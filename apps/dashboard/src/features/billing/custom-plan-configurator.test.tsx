import { act, fireEvent, render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { CustomPlanConfigurator } from './custom-plan-configurator';

describe('CustomPlanConfigurator', () => {
  it('renders default configuration and calculates initial monthly price', () => {
    render(<CustomPlanConfigurator initialCores={2} initialRamGb={4} initialDiskGb={50} />);

    expect(screen.getByText('Custom Dedicated Cloud Plan')).toBeDefined();
    expect(screen.getByText('2 Cores')).toBeDefined();
    expect(screen.getAllByText('4 GB RAM').length).toBeGreaterThanOrEqual(1);
    expect(screen.getByText('50 GB SSD')).toBeDefined();
    // 10 base + 10 cpu + 12 ram + 5 disk = $37.00
    expect(screen.getByText('$37.00')).toBeDefined();
  });

  it('updates price when resource slider value changes', async () => {
    render(<CustomPlanConfigurator initialCores={2} initialRamGb={4} initialDiskGb={50} />);

    const cpuSlider = screen.getByLabelText('Dedicated vCPU Cores');
    await act(async () => {
      fireEvent.change(cpuSlider, { target: { value: '4' } });
    });

    // 10 base + 20 cpu + 12 ram + 5 disk = $47.00
    expect(screen.getByText('4 Cores')).toBeDefined();
    expect(screen.getByText('$47.00')).toBeDefined();
  });

  it('submits selected quote when form is submitted', async () => {
    const handleSelect = vi.fn();
    render(
      <CustomPlanConfigurator
        initialCores={4}
        initialRamGb={8}
        initialDiskGb={100}
        onSelectPlan={handleSelect}
      />
    );

    const submitBtn = screen.getByRole('button', { name: /configure & deploy/i });
    await act(async () => {
      fireEvent.click(submitBtn);
    });

    expect(handleSelect).toHaveBeenCalledTimes(1);
    expect(handleSelect).toHaveBeenCalledWith(
      expect.objectContaining({
        cores: 4,
        ramGb: 8,
        diskGb: 100,
        monthlyPriceCents: 6400, // 10 base + 20 cpu + 24 ram + 10 disk = 64
      })
    );
  });
});
