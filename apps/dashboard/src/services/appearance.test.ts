import { beforeEach, describe, expect, it, vi } from 'vitest';
import type { AppearancePreferences } from '#/interfaces/appearance';
import {
  APPEARANCE_CACHE_POINTER,
  APPEARANCE_CHANGED,
  appearanceCacheKey,
  applyAppearance,
  bootAppearance,
  cacheAppearance,
  cachedAppearance,
  resolveAppearanceMode,
} from './appearance';

describe('appearance-runtime', () => {
  const samplePref: AppearancePreferences = {
    activeTheme: 'custom-1',
    mode: 'dim',
    radius: '0.75rem',
    fontFamily: 'Inter, sans-serif',
    tokens: {
      '--primary': '#8b5cf6',
      '--malicious': 'red; background: black',
    },
    customThemes: [
      {
        id: 'custom-1',
        name: 'Custom Dim',
        tokens: {
          '--card': '#1f232b',
        },
      },
    ],
  };

  beforeEach(() => {
    localStorage.clear();
    document.documentElement.className = '';
    document.documentElement.removeAttribute('data-theme');
    document.documentElement.removeAttribute('style');
  });

  it('generates consistent cache keys', () => {
    expect(appearanceCacheKey('user-123')).toBe('codedock-appearance-v1:user-123');
  });

  it('handles caching and retrieving appearance preferences', () => {
    const key = appearanceCacheKey('user-1');
    expect(cachedAppearance(key)).toBeNull();

    cacheAppearance(key, samplePref, 'dim');
    expect(localStorage.getItem(APPEARANCE_CACHE_POINTER)).toBe(key);
    expect(cachedAppearance(key)).toEqual(samplePref);
  });

  it('returns null on invalid JSON in cached appearance', () => {
    const key = appearanceCacheKey('user-corrupt');
    localStorage.setItem(key, '{ invalid json');
    expect(cachedAppearance(key)).toBeNull();
  });

  it('resolves appearance modes properly', () => {
    expect(resolveAppearanceMode({ activeTheme: 'default', mode: 'dim' }, false)).toBe('dim');
    expect(resolveAppearanceMode({ activeTheme: 'default', mode: 'light' }, true)).toBe('light');
    expect(resolveAppearanceMode({ activeTheme: 'default', mode: 'dark' }, false)).toBe('dark');
    expect(resolveAppearanceMode({ activeTheme: 'default', mode: 'system' }, true)).toBe('dim');
    expect(resolveAppearanceMode({ activeTheme: 'default', mode: 'system' }, false)).toBe('light');
  });

  it('applies appearance styles and variables cleanly without injection', () => {
    const dispatchSpy = vi.spyOn(window, 'dispatchEvent');

    applyAppearance(samplePref, 'dim');

    const root = document.documentElement;
    expect(root.getAttribute('data-theme')).toBe('dim');
    expect(root.classList.contains('dim')).toBe(true);
    expect(root.style.getPropertyValue('--radius')).toBe('0.75rem');
    expect(root.style.getPropertyValue('--font-sans')).toBe('Inter, sans-serif');
    expect(root.style.getPropertyValue('--primary')).toBe('#8b5cf6');
    expect(root.style.getPropertyValue('--card')).toBe('#1f232b');
    expect(root.style.getPropertyValue('--malicious')).toBe('');

    expect(dispatchSpy).toHaveBeenCalledWith(expect.objectContaining({ type: APPEARANCE_CHANGED }));
  });

  it('boots appearance from cached storage safely', () => {
    const key = appearanceCacheKey('user-boot');
    cacheAppearance(key, samplePref, 'dim');

    bootAppearance();

    const root = document.documentElement;
    expect(root.getAttribute('data-theme')).toBe('dim');
    expect(root.classList.contains('dim')).toBe(true);
  });
});
