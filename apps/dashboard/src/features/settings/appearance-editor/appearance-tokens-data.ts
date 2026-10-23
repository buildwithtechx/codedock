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
  { key: '--th-fg-default', label: 'Primary Text', category: 'text', defaultValue: '#f8f8fc' },
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
    defaultValue: '0 4px 12px rgba(0,0,0,0.1)',
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
