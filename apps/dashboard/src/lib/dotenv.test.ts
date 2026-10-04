import { describe, expect, it } from 'vitest';
import { parseDotenv, serializeDotenv } from './dotenv';

describe('parseDotenv', () => {
  it('parses simple key-value pairs', () => {
    const input = 'FOO=bar\nBAZ=qux';
    const result = parseDotenv(input);
    expect(result).toEqual([
      { key: 'FOO', value: 'bar' },
      { key: 'BAZ', value: 'qux' },
    ]);
  });

  it('handles export prefix and whitespace', () => {
    const input = '  export APP_PORT=8080  \nexport DB_HOST=localhost';
    const result = parseDotenv(input);
    expect(result).toEqual([
      { key: 'APP_PORT', value: '8080' },
      { key: 'DB_HOST', value: 'localhost' },
    ]);
  });

  it('handles single-quoted and double-quoted values', () => {
    const input = `SINGLE='hello world'\nDOUBLE="quoted \\"value\\" with newline\\n"`;
    const result = parseDotenv(input);
    expect(result).toHaveLength(2);
    expect(result[0]).toEqual({ key: 'SINGLE', value: 'hello world' });
    expect(result[1].key).toBe('DOUBLE');
    expect(result[1].value).toContain('quoted "value" with newline');
  });

  it('ignores comments and empty lines', () => {
    const input =
      '# This is a comment\n\nNAME=CodeDock # trailing comment\n\n# Another comment\nVERSION=1.0.0';
    const result = parseDotenv(input);
    expect(result).toEqual([
      { key: 'NAME', value: 'CodeDock' },
      { key: 'VERSION', value: '1.0.0' },
    ]);
  });

  it('handles empty string input', () => {
    expect(parseDotenv('')).toEqual([]);
    expect(parseDotenv('   \n\n  ')).toEqual([]);
  });
});

describe('serializeDotenv', () => {
  it('serializes simple variables', () => {
    const entries = [
      { key: 'PORT', value: '3000' },
      { key: 'NODE_ENV', value: 'production' },
    ];
    const output = serializeDotenv(entries);
    expect(output).toBe('PORT=3000\nNODE_ENV=production\n');
  });

  it('quotes values containing spaces and special characters', () => {
    const entries = [
      { key: 'GREETING', value: 'hello world' },
      { key: 'URL', value: 'https://example.com/api?a=1&b=2' },
    ];
    const output = serializeDotenv(entries);
    expect(output).toContain("GREETING='hello world'");
  });

  it('roundtrips cleanly between parse and serialize', () => {
    const initial = [
      { key: 'APP_NAME', value: 'Codedock Platform' },
      { key: 'PORT', value: '8080' },
      { key: 'SECRET', value: 'super-secret-key-123' },
    ];
    const serialized = serializeDotenv(initial);
    const parsed = parseDotenv(serialized);
    expect(parsed).toEqual(initial);
  });
});
