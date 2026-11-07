import { describe, expect, it } from 'vitest';
import { commandEnvironmentWarning, inspectCommandEnvironment } from './command-environment';

describe('command-environment parser', () => {
  it('parses simple key=value pairs', () => {
    const res = inspectCommandEnvironment('PORT=3000 NODE_ENV=production npm start');
    expect(res.keys).toEqual(['PORT', 'NODE_ENV']);
    expect(res.variables).toEqual({ PORT: '3000', NODE_ENV: 'production' });
    expect(res.command).toBe('npm start');
    expect(res.migratable).toBe(true);
  });

  it('handles empty input gracefully', () => {
    const res = inspectCommandEnvironment('');
    expect(res.keys).toEqual([]);
    expect(res.migratable).toBe(false);
  });

  it('handles export prefix syntax', () => {
    const res = inspectCommandEnvironment('export PORT=8080 && node server.js');
    expect(res.keys).toEqual(['PORT']);
    expect(res.variables).toEqual({ PORT: '8080' });
    expect(res.command).toBe('node server.js');
    expect(res.migratable).toBe(true);
  });

  it('generates warnings for embedded environment variables', () => {
    const warn = commandEnvironmentWarning('FOO=bar ./run', 'Build command');
    expect(warn).toContain('Build command embeds environment variables: FOO');
    expect(commandEnvironmentWarning('node index.js', 'Start')).toBeUndefined();
  });
});
