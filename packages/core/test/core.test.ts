import { describe, expect, it } from 'vitest';
import {
  DEFAULT_ROLLBACK_WINDOW,
  MAX_ROLLBACK_WINDOW,
  clampRollbackWindow,
  decodeDeploySlug,
  encodeDeploySlug,
  formatBytes,
  formatDuration,
  isValidDomain,
  isValidUrl,
  normalizeRollbackWindow,
} from '../src/index';

describe('@codedock/core', () => {
  describe('rollback-window', () => {
    it('has standard defaults and limits', () => {
      expect(DEFAULT_ROLLBACK_WINDOW).toBe(5);
      expect(MAX_ROLLBACK_WINDOW).toBe(20);
    });

    it('clamps rollback window properly', () => {
      expect(clampRollbackWindow(0)).toBe(0);
      expect(clampRollbackWindow(-5)).toBe(0);
      expect(clampRollbackWindow(5)).toBe(5);
      expect(clampRollbackWindow(25)).toBe(20);
      expect(clampRollbackWindow(10, 8)).toBe(8);
      expect(clampRollbackWindow(3.7)).toBe(3);
    });

    it('normalizes unknown values with fallback', () => {
      expect(normalizeRollbackWindow(undefined)).toBe(5);
      expect(normalizeRollbackWindow(null)).toBe(5);
      expect(normalizeRollbackWindow('')).toBe(5);
      expect(normalizeRollbackWindow('  ')).toBe(5);
      expect(normalizeRollbackWindow('invalid')).toBe(5);
      expect(normalizeRollbackWindow('8')).toBe(8);
      expect(normalizeRollbackWindow(12)).toBe(12);
      expect(normalizeRollbackWindow(50)).toBe(20);
      expect(normalizeRollbackWindow(-1)).toBe(0);
      expect(normalizeRollbackWindow(undefined, 7)).toBe(7);
    });
  });

  describe('format', () => {
    it('formats bytes into units', () => {
      expect(formatBytes(0)).toBe('0 B');
      expect(formatBytes(1024)).toBe('1 KB');
      expect(formatBytes(1024 * 1024)).toBe('1 MB');
      expect(formatBytes(1024 * 1024 * 1024)).toBe('1 GB');
      expect(formatBytes(1536, 1)).toBe('1.5 KB');
      expect(formatBytes(-10)).toBe('—');
    });

    it('formats duration', () => {
      expect(formatDuration(0)).toBe('0s');
      expect(formatDuration(45)).toBe('45s');
      expect(formatDuration(90)).toBe('1m 30s');
      expect(formatDuration(120)).toBe('2m');
      expect(formatDuration(3665)).toBe('1h 1m');
    });
  });

  describe('slug', () => {
    it('encodes names into safe deployment slugs', () => {
      expect(encodeDeploySlug('My Awesome App!')).toBe('my-awesome-app');
      expect(encodeDeploySlug('  ---test___app---  ')).toBe('test-app');
      expect(encodeDeploySlug('prod-api-v2')).toBe('prod-api-v2');
      expect(encodeDeploySlug('')).toBe('');
    });

    it('decodes slugs into title words', () => {
      expect(decodeDeploySlug('my-awesome-app')).toBe('My Awesome App');
      expect(decodeDeploySlug('api-service')).toBe('Api Service');
      expect(decodeDeploySlug('')).toBe('');
    });
  });

  describe('validation', () => {
    it('validates domain syntax', () => {
      expect(isValidDomain('example.com')).toBe(true);
      expect(isValidDomain('sub.app.codedock.io')).toBe(true);
      expect(isValidDomain('invalid_domain')).toBe(false);
      expect(isValidDomain('http://example.com')).toBe(false);
      expect(isValidDomain('')).toBe(false);
    });

    it('validates URLs', () => {
      expect(isValidUrl('https://example.com')).toBe(true);
      expect(isValidUrl('http://localhost:8080/api')).toBe(true);
      expect(isValidUrl('ftp://ftp.example.com')).toBe(false);
      expect(isValidUrl('not a url')).toBe(false);
    });
  });
});
