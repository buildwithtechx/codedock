import { fireEvent, render, screen } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { AppearanceSetting } from './appearance-setting';

const mockSetTheme = vi.fn();
vi.mock('@codedock/ui', () => ({
  useTheme: () => ({
    theme: 'dark',
    setTheme: mockSetTheme,
  }),
}));

vi.mock('@tanstack/react-router', () => ({
  Link: ({ children, to, search }: { children: React.ReactNode; to: string; search?: unknown }) => (
    <a href={`${to}?search=${JSON.stringify(search)}`}>{children}</a>
  ),
}));

vi.mock('#/stores/auth-store', () => ({
  useAuthStore: (selector: (s: { user: { id: string } }) => unknown) =>
    selector({ user: { id: 'user-1' } }),
}));

vi.mock('#/services/appearance', () => ({
  appearanceCacheKey: (id: string) => `codedock-appearance-v1:${id}`,
  cachedAppearance: vi.fn().mockReturnValue(null),
  cacheAppearance: vi.fn(),
  applyAppearance: vi.fn(),
  getAppearance: vi.fn().mockResolvedValue(null),
  setAppearance: vi.fn().mockResolvedValue(null),
}));

describe('AppearanceSetting', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders base theme options and customize button', () => {
    render(<AppearanceSetting />);
    expect(screen.getByText('Appearance')).toBeDefined();
    expect(screen.getByText('Light')).toBeDefined();
    expect(screen.getByText('Dim')).toBeDefined();
    expect(screen.getByText('Dark')).toBeDefined();
    expect(screen.getByText('System')).toBeDefined();
    expect(screen.getByText('Customize')).toBeDefined();
  });

  it('clicking a base theme option triggers setTheme', () => {
    render(<AppearanceSetting />);
    const lightRadio = screen.getByDisplayValue('light');
    fireEvent.click(lightRadio);

    expect(mockSetTheme).toHaveBeenCalledWith('light');
  });
});
