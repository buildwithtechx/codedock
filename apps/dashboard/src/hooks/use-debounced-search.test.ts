import { act, renderHook } from '@testing-library/react';
import { afterEach, describe, expect, it, vi } from 'vitest';
import { useDebouncedSearch } from './use-debounced-search';

afterEach(() => {
  vi.useRealTimers();
});

describe('useDebouncedSearch', () => {
  it('returns the initial value immediately', () => {
    const { result } = renderHook(() => useDebouncedSearch('billing'));
    expect(result.current).toBe('billing');
  });

  it('settles to the latest value after the delay', () => {
    vi.useFakeTimers();
    const { result, rerender } = renderHook(({ value }) => useDebouncedSearch(value), {
      initialProps: { value: 'a' },
    });
    rerender({ value: 'ab' });
    rerender({ value: 'abc' });
    expect(result.current).toBe('a');
    act(() => {
      vi.advanceTimersByTime(300);
    });
    expect(result.current).toBe('abc');
  });

  it('does not settle before the delay elapses', () => {
    vi.useFakeTimers();
    const { result, rerender } = renderHook(({ value }) => useDebouncedSearch(value), {
      initialProps: { value: 'a' },
    });
    rerender({ value: 'ab' });
    act(() => {
      vi.advanceTimersByTime(299);
    });
    expect(result.current).toBe('a');
  });
});
