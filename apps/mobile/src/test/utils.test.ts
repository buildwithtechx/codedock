import { cleanUrl, formatDuration, formatStatus, truncate } from '../utils';

describe('Mobile utils', () => {
  test('cleanUrl trims trailing slashes', () => {
    expect(cleanUrl('https://codedock.example.com///')).toBe('https://codedock.example.com');
  });

  test('formatStatus classifies running and failed', () => {
    const running = formatStatus('running');
    expect(running.isSuccess).toBe(true);

    const failed = formatStatus('failed');
    expect(failed.isDanger).toBe(true);
  });

  test('formatDuration formats seconds and minutes', () => {
    expect(formatDuration(45)).toBe('45s');
    expect(formatDuration(125)).toBe('2m 5s');
  });

  test('truncate caps string length', () => {
    expect(truncate('hello world', 5)).toBe('hell…');
  });
});
