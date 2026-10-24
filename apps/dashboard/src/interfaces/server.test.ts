import { describe, expect, it } from 'vitest';
import { parseServerMetrics, type ServerMetrics } from './server';

describe('parseServerMetrics', () => {
  const sampleMetrics: ServerMetrics = {
    cpu_usage_percentage: 24.5,
    memory_usage_bytes: 1024 * 1024 * 512,
    memory_limit_bytes: 1024 * 1024 * 2048,
    disk_usage_bytes: 1024 * 1024 * 1024 * 10,
    disk_total_bytes: 1024 * 1024 * 1024 * 50,
  };

  it('returns null when metrics is null or undefined', () => {
    expect(parseServerMetrics(null)).toBeNull();
    expect(parseServerMetrics(undefined)).toBeNull();
  });

  it('returns the object directly when given a valid ServerMetrics object', () => {
    const result = parseServerMetrics(sampleMetrics);
    expect(result).toEqual(sampleMetrics);
  });

  it('parses valid raw JSON string', () => {
    const jsonStr = JSON.stringify(sampleMetrics);
    const result = parseServerMetrics(jsonStr);
    expect(result).toEqual(sampleMetrics);
  });

  it('parses base64-encoded JSON string', () => {
    const jsonStr = JSON.stringify(sampleMetrics);
    const base64Str = btoa(jsonStr);
    const result = parseServerMetrics(base64Str);
    expect(result).toEqual(sampleMetrics);
  });

  it('returns null on invalid or corrupt input', () => {
    expect(parseServerMetrics('invalid-json')).toBeNull();
    expect(parseServerMetrics('{broken:json}')).toBeNull();
  });

  it('normalizes partial objects so numeric fields never crash render', () => {
    const result = parseServerMetrics({} as ServerMetrics);
    expect(result).toEqual({
      cpu_usage_percentage: 0,
      memory_usage_bytes: 0,
      memory_limit_bytes: 0,
      disk_usage_bytes: 0,
      disk_total_bytes: 0,
    });
  });

  it('normalizes partial JSON strings and drops non-finite values', () => {
    const result = parseServerMetrics('{"cpu_usage_percentage":12.5}');
    expect(result?.cpu_usage_percentage).toBe(12.5);
    expect(result?.memory_limit_bytes).toBe(0);
  });
});
