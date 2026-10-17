import { describe, expect, it } from 'vitest';
import { commandEnvironmentWarning, inspectCommandEnvironment } from './command-environment';

describe('command-environment parser', () => {
  it('parses leading variables with prefixes', () => {
    const res = inspectCommandEnvironment('NODE_ENV=production PORT=8080 npm start');
    expect(res.migratable).toBe(true);
    expect(res.keys).toEqual(['NODE_ENV', 'PORT']);
    expect(res.variables).toEqual({ NODE_ENV: 'production', PORT: '8080' });
    expect(res.command).toBe('npm start');
  });

  it('handles cross-env and quotes', () => {
    const res = inspectCommandEnvironment('cross-env FOO=\'bar baz\' BAR="hello" bun run dev');
    expect(res.migratable).toBe(true);
    expect(res.variables).toEqual({ FOO: 'bar baz', BAR: 'hello' });
    expect(res.command).toBe('bun run dev');
  });

  it('rejects shell expansion and marks non-migratable', () => {
    const res = inspectCommandEnvironment('TOKEN=$SECRET ./run.sh');
    expect(res.migratable).toBe(false);
  });

  it('detects chained assignments and denies migration', () => {
    const res = inspectCommandEnvironment('prepare && TOKEN=xyz start');
    expect(res.keys).toContain('TOKEN');
    expect(res.migratable).toBe(false);
  });

  it('formats warning text', () => {
    const warn = commandEnvironmentWarning(
      'cross-env NODE_ENV=production npm run build',
      'Build command'
    );
    expect(warn).toContain('Build command embeds environment variables: NODE_ENV');
  });
});
