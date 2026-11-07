import { describe, expect, it } from 'vitest';
import {
  BUILD_CACHE_BUDGET_MAX_GIB,
  BUILD_CACHE_BUDGET_MIN_GIB,
  budgetBytesToGib,
  parseBudgetGib,
} from './server-build-cache';

describe('server-build-cache', () => {
  const GIB = 1024 ** 3;

  describe('budgetBytesToGib', () => {
    it('returns empty string for null or undefined', () => {
      expect(budgetBytesToGib(null)).toBe('');
      expect(budgetBytesToGib(undefined)).toBe('');
    });

    it('formats exact GiB as integer string', () => {
      expect(budgetBytesToGib(24 * GIB)).toBe('24');
    });

    it('rounds fractional GiB to two decimals', () => {
      expect(budgetBytesToGib(2.5 * GIB)).toBe('2.5');
      expect(budgetBytesToGib(2.555 * GIB)).toBe('2.56');
    });
  });

  describe('parseBudgetGib', () => {
    it('parses valid integer and decimal numbers in GiB into bytes', () => {
      expect(parseBudgetGib('24')).toBe(24 * GIB);
      expect(parseBudgetGib('24.5')).toBe(Math.round(24.5 * GIB));
      expect(parseBudgetGib('24,5')).toBe(Math.round(24.5 * GIB));
    });

    it('returns null for values below min (1 GiB)', () => {
      expect(parseBudgetGib('0')).toBeNull();
      expect(parseBudgetGib('0.5')).toBeNull();
      expect(parseBudgetGib('-5')).toBeNull();
    });

    it('returns null for values above max (65536 GiB)', () => {
      expect(parseBudgetGib(String(BUILD_CACHE_BUDGET_MAX_GIB + 1))).toBeNull();
    });

    it('accepts boundary values min and max', () => {
      expect(parseBudgetGib(String(BUILD_CACHE_BUDGET_MIN_GIB))).toBe(
        BUILD_CACHE_BUDGET_MIN_GIB * GIB
      );
      expect(parseBudgetGib(String(BUILD_CACHE_BUDGET_MAX_GIB))).toBe(
        BUILD_CACHE_BUDGET_MAX_GIB * GIB
      );
    });

    it('returns null for non-numeric or malformed strings', () => {
      expect(parseBudgetGib('abc')).toBeNull();
      expect(parseBudgetGib('')).toBeNull();
      expect(parseBudgetGib('1.2.3')).toBeNull();
      expect(parseBudgetGib('10.999')).toBeNull(); // more than 2 decimal places
    });
  });
});
