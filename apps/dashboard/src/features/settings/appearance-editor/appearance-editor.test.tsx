import { fireEvent, render, screen } from '@testing-library/react';
import type * as React from 'react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { AppearanceEditor } from './appearance-editor';

const mockSetTheme = vi.fn();
vi.mock('@codedock/ui', () => ({
  useTheme: () => ({
    theme: 'dark',
    resolvedTheme: 'dark',
    setTheme: mockSetTheme,
  }),
}));

vi.mock('@tanstack/react-router', () => ({
  Link: ({ children, to }: { children: React.ReactNode; to: string }) => (
    <a href={to}>{children}</a>
  ),
}));

vi.mock('#/stores/auth-store', () => ({
  useAuthStore: (selector: (s: { user: { id: string; name: string } }) => unknown) =>
    selector({ user: { id: 'user-1', name: 'Solomon Olatunji' } }),
}));

vi.mock('#/services/appearance', () => ({
  appearanceCacheKey: (id: string) => `codedock-appearance-v1:${id}`,
  cachedAppearance: vi.fn().mockReturnValue({
    activeTheme: 'neutral-slate',
    mode: 'dark',
    tokens: {},
    customThemes: [],
  }),
  cacheAppearance: vi.fn(),
  applyAppearance: vi.fn(),
  getAppearance: vi.fn().mockResolvedValue(null),
  setAppearance: vi.fn().mockResolvedValue(null),
}));

describe('AppearanceEditor (OpenShip & Coolify UX/UI)', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders OpenShip appearance header, import and download buttons', () => {
    render(<AppearanceEditor />);
    expect(screen.getByText('Appearance')).toBeDefined();
    expect(screen.getByText('Back to General')).toBeDefined();
    expect(screen.getByText('Import')).toBeDefined();
    expect(screen.getByText('Download theme')).toBeDefined();

    // Mode strip
    expect(screen.getByRole('button', { name: /Light/i })).toBeDefined();
    expect(screen.getByRole('button', { name: /Dim/i })).toBeDefined();
    expect(screen.getByRole('button', { name: /Dark/i })).toBeDefined();
    expect(screen.getByRole('button', { name: /System/i })).toBeDefined();
    expect(screen.getByRole('button', { name: /New theme/i })).toBeDefined();
  });

  it('switches mode when clicking mode strip buttons', () => {
    render(<AppearanceEditor />);
    const lightBtn = screen.getByRole('button', { name: /Light/i });
    fireEvent.click(lightBtn);

    expect(mockSetTheme).toHaveBeenCalledWith('light');
  });

  it('renders design tokens, typography, and live preview card', () => {
    render(<AppearanceEditor />);
    expect(screen.getByText('Design tokens')).toBeDefined();
    expect(screen.getByPlaceholderText('Find a token...')).toBeDefined();
    expect(screen.getByText('Typography')).toBeDefined();
    expect(screen.getByText('Live preview')).toBeDefined();
    expect(screen.getByText('Your workspace')).toBeDefined();
    expect(screen.getByText('Solomon Olatunji')).toBeDefined();
    expect(screen.getByText('build completed')).toBeDefined();
  });

  it('allows creating a new custom theme and shows floating save bar', () => {
    render(<AppearanceEditor />);
    const newThemeBtn = screen.getByRole('button', { name: /New theme/i });
    fireEvent.click(newThemeBtn);

    expect(screen.getByText('Unsaved changes')).toBeDefined();
    expect(screen.getByRole('button', { name: /Save/i })).toBeDefined();
  });
});
