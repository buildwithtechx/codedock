import { describe, expect, it } from 'vitest';
import { isCloudPurchaseIntent } from './billing-intent';

describe('isCloudPurchaseIntent', () => {
  it('accepts valid UUIDv4 strings', () => {
    expect(isCloudPurchaseIntent('c3b12345-6789-4abc-8def-0123456789ab')).toBe(true);
    expect(isCloudPurchaseIntent('C3B12345-6789-4ABC-8DEF-0123456789AB')).toBe(true);
  });

  it('rejects non-uuid and invalid versions', () => {
    expect(isCloudPurchaseIntent('')).toBe(false);
    expect(isCloudPurchaseIntent(null)).toBe(false);
    expect(isCloudPurchaseIntent('not-a-uuid')).toBe(false);
    expect(isCloudPurchaseIntent('c3b12345-6789-1abc-8def-0123456789ab')).toBe(false);
  });
});
