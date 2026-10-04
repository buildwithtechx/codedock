import { fireEvent, render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { EnvBulkModal } from './env-bulk-modal';

describe('EnvBulkModal', () => {
  it('renders nothing when isOpen is false', () => {
    const { container } = render(
      <EnvBulkModal isOpen={false} onClose={vi.fn()} onImport={vi.fn()} />
    );
    expect(container.firstChild).toBeNull();
  });

  it('renders modal dialog when isOpen is true', () => {
    render(<EnvBulkModal isOpen={true} onClose={vi.fn()} onImport={vi.fn()} />);
    expect(screen.getByText('Bulk Import Environment Variables')).toBeDefined();
    expect(screen.getByText('Parse & Import')).toBeDefined();
    expect(screen.getByText('Cancel')).toBeDefined();
  });

  it('parses input and calls onImport with parsed key-value entries', () => {
    const onImport = vi.fn();
    const onClose = vi.fn();
    render(<EnvBulkModal isOpen={true} onClose={onClose} onImport={onImport} />);

    const textarea = screen.getByPlaceholderText(/DATABASE_URL/);
    fireEvent.change(textarea, {
      target: { value: 'API_KEY=secret_123\nPORT=4000\n# Comment\nDEBUG=true' },
    });

    const submitBtn = screen.getByText('Parse & Import');
    fireEvent.click(submitBtn);

    expect(onImport).toHaveBeenCalledWith([
      { key: 'API_KEY', value: 'secret_123' },
      { key: 'PORT', value: '4000' },
      { key: 'DEBUG', value: 'true' },
    ]);
    expect(onClose).toHaveBeenCalled();
  });

  it('calls onClose when Cancel button is clicked', () => {
    const onClose = vi.fn();
    render(<EnvBulkModal isOpen={true} onClose={onClose} onImport={vi.fn()} />);

    const cancelBtn = screen.getByText('Cancel');
    fireEvent.click(cancelBtn);

    expect(onClose).toHaveBeenCalled();
  });
});
