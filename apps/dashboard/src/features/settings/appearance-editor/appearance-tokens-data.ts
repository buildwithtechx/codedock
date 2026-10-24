export type TokenCategory = 'surfaces' | 'text' | 'controls' | 'status' | 'effects' | 'scales';

export interface TokenDefinition {
  key: string;
  label: string;
  category: TokenCategory;
  defaultValue: string;
}

export const TOKEN_DEFINITIONS: TokenDefinition[] = [
  // Surfaces
  { key: '--th-bg-page', label: 'Page Background', category: 'surfaces', defaultValue: '#000000' },
  { key: '--th-bg-card', label: 'Card Surface', category: 'surfaces', defaultValue: '#0a0a0a' },
  {
    key: '--th-bg-subtle',
    label: 'Subtle Background',
    category: 'surfaces',
    defaultValue: 'rgba(255, 255, 255, .025)',
  },
  {
    key: '--th-bg-hover',
    label: 'Hover State',
    category: 'surfaces',
    defaultValue: 'rgba(255, 255, 255, .04)',
  },
  {
    key: '--th-bg-inset',
    label: 'Inset Surface',
    category: 'surfaces',
    defaultValue: 'rgba(255, 255, 255, .015)',
  },
  {
    key: '--th-bd-default',
    label: 'Default Border',
    category: 'surfaces',
    defaultValue: 'rgba(255, 255, 255, .08)',
  },
  {
    key: '--th-bd-subtle',
    label: 'Subtle Border',
    category: 'surfaces',
    defaultValue: 'rgba(255, 255, 255, .05)',
  },
  {
    key: '--th-bd-strong',
    label: 'Strong Border',
    category: 'surfaces',
    defaultValue: 'rgba(255, 255, 255, .14)',
  },
  { key: '--th-card-bd', label: 'Card Border', category: 'surfaces', defaultValue: 'transparent' },
  {
    key: '--th-card-bd-hover',
    label: 'Card Border Hover',
    category: 'surfaces',
    defaultValue: 'rgba(255, 255, 255, .12)',
  },
  {
    key: '--th-divider',
    label: 'Divider Lines',
    category: 'surfaces',
    defaultValue: 'rgba(255, 255, 255, .08)',
  },

  // Text
  { key: '--th-fg-default', label: 'Primary Text', category: 'text', defaultValue: '#ededed' },
  { key: '--th-fg-muted', label: 'Muted Text', category: 'text', defaultValue: '#a9adba' },
  { key: '--th-fg-subtle', label: 'Subtle Text', category: 'text', defaultValue: '#71717a' },

  // Controls
  { key: '--th-accent', label: 'Accent / Brand', category: 'controls', defaultValue: '#8b5cf6' },
  {
    key: '--th-accent-hover',
    label: 'Accent Hover',
    category: 'controls',
    defaultValue: '#9d74f8',
  },
  {
    key: '--th-accent-fg',
    label: 'Accent Foreground',
    category: 'controls',
    defaultValue: '#ffffff',
  },

  // Status
  { key: '--th-status-active', label: 'Active', category: 'status', defaultValue: '#22c55e' },
  { key: '--th-status-queued', label: 'Queued', category: 'status', defaultValue: '#f59e0b' },
  { key: '--th-status-failed', label: 'Failed', category: 'status', defaultValue: '#ef4444' },

  // Effects
  { key: '--th-radius', label: 'Corner Radius', category: 'effects', defaultValue: '0.5rem' },
  {
    key: '--th-shadow',
    label: 'Card Shadow',
    category: 'effects',
    defaultValue: '0 4px 12px rgba(0,0,0,0.4)',
  },

  // Scales
  { key: '--th-scale-base', label: 'Scale Base', category: 'scales', defaultValue: '1rem' },
  { key: '--th-spacing', label: 'Spacing Unit', category: 'scales', defaultValue: '1rem' },
];

export const TOKEN_CATEGORIES: { id: TokenCategory; label: string }[] = [
  { id: 'surfaces', label: 'Surfaces' },
  { id: 'text', label: 'Text' },
  { id: 'controls', label: 'Controls' },
  { id: 'status', label: 'Status' },
  { id: 'effects', label: 'Effects' },
  { id: 'scales', label: 'Scales' },
];

export const MODE_DEFAULTS: Record<'dark' | 'dim' | 'light', Record<string, string>> = {
  dark: {
    '--th-bg-page': '#000000',
    '--th-bg-card': '#0a0a0a',
    '--th-bg-subtle': 'rgba(255, 255, 255, .025)',
    '--th-bg-hover': 'rgba(255, 255, 255, .04)',
    '--th-bg-inset': 'rgba(255, 255, 255, .015)',
    '--th-bd-default': 'rgba(255, 255, 255, .08)',
    '--th-bd-subtle': 'rgba(255, 255, 255, .05)',
    '--th-bd-strong': 'rgba(255, 255, 255, .14)',
    '--th-card-bd': 'transparent',
    '--th-card-bd-hover': 'rgba(255, 255, 255, .12)',
    '--th-divider': 'rgba(255, 255, 255, .08)',
    '--th-fg-default': '#ededed',
    '--th-fg-muted': '#a9adba',
    '--th-fg-subtle': '#71717a',
    '--th-accent': '#8b5cf6',
    '--th-accent-hover': '#9d74f8',
    '--th-accent-fg': '#ffffff',
    '--th-status-active': '#22c55e',
    '--th-status-queued': '#f59e0b',
    '--th-status-failed': '#ef4444',
    '--th-radius': '0.5rem',
    '--th-shadow': '0 4px 12px rgba(0,0,0,0.4)',
    '--th-scale-base': '1rem',
    '--th-spacing': '1rem',
  },
  dim: {
    '--th-bg-page': '#12151a',
    '--th-bg-card': '#1a1d24',
    '--th-bg-subtle': 'rgba(255, 255, 255, .03)',
    '--th-bg-hover': 'rgba(255, 255, 255, .05)',
    '--th-bg-inset': 'rgba(0, 0, 0, .2)',
    '--th-bd-default': 'rgba(255, 255, 255, .1)',
    '--th-bd-subtle': 'rgba(255, 255, 255, .06)',
    '--th-bd-strong': 'rgba(255, 255, 255, .18)',
    '--th-card-bd': 'transparent',
    '--th-card-bd-hover': 'rgba(255, 255, 255, .15)',
    '--th-divider': 'rgba(255, 255, 255, .1)',
    '--th-fg-default': '#f1f3f7',
    '--th-fg-muted': '#9ba3af',
    '--th-fg-subtle': '#6b7280',
    '--th-accent': '#8b5cf6',
    '--th-accent-hover': '#9d74f8',
    '--th-accent-fg': '#ffffff',
    '--th-status-active': '#10b981',
    '--th-status-queued': '#f59e0b',
    '--th-status-failed': '#f43f5e',
    '--th-radius': '0.5rem',
    '--th-shadow': '0 4px 12px rgba(0,0,0,0.3)',
    '--th-scale-base': '1rem',
    '--th-spacing': '1rem',
  },
  light: {
    '--th-bg-page': '#f9fafb',
    '--th-bg-card': '#ffffff',
    '--th-bg-subtle': 'rgba(0, 0, 0, .02)',
    '--th-bg-hover': 'rgba(0, 0, 0, .04)',
    '--th-bg-inset': 'rgba(0, 0, 0, .015)',
    '--th-bd-default': 'rgba(0, 0, 0, .08)',
    '--th-bd-subtle': 'rgba(0, 0, 0, .04)',
    '--th-bd-strong': 'rgba(0, 0, 0, .16)',
    '--th-card-bd': 'transparent',
    '--th-card-bd-hover': 'rgba(0, 0, 0, .12)',
    '--th-divider': 'rgba(0, 0, 0, .06)',
    '--th-fg-default': '#0f172a',
    '--th-fg-muted': '#64748b',
    '--th-fg-subtle': '#94a3b8',
    '--th-accent': '#7c3aed',
    '--th-accent-hover': '#6d28d9',
    '--th-accent-fg': '#ffffff',
    '--th-status-active': '#16a34a',
    '--th-status-queued': '#d97706',
    '--th-status-failed': '#dc2626',
    '--th-radius': '0.5rem',
    '--th-shadow': '0 4px 12px rgba(0,0,0,0.06)',
    '--th-scale-base': '1rem',
    '--th-spacing': '1rem',
  },
};

export function getEffectiveThemeTokens(
  mode: 'dark' | 'dim' | 'light',
  overrides: Record<string, string> = {}
): Record<string, string> {
  const base = MODE_DEFAULTS[mode] || MODE_DEFAULTS.dark;
  return { ...base, ...overrides };
}
