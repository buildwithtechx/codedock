import { describe, expect, it, vi } from 'vitest';
import {
  type CodedockThemeFile,
  downloadThemeFile,
  parseThemeFile,
  validateThemeFile,
} from './theme-export-import';

describe('theme-export-import', () => {
  const validTheme: CodedockThemeFile = {
    version: 1,
    name: 'Tokyo Night Custom',
    mode: 'dark',
    radius: '0.5rem',
    tokens: {
      '--background': '#1a1b26',
      '--primary': '#7aa2f7',
    },
  };

  it('validates a correct theme file', () => {
    const res = validateThemeFile(validTheme);
    expect(res.valid).toBe(true);
    if (res.valid) {
      expect(res.theme.name).toBe('Tokyo Night Custom');
      expect(res.theme.tokens['--primary']).toBe('#7aa2f7');
    }
  });

  it('rejects missing or unsupported version', () => {
    const res = validateThemeFile({ ...validTheme, version: 2 });
    expect(res.valid).toBe(false);
  });

  it('rejects invalid mode', () => {
    const res = validateThemeFile({ ...validTheme, mode: 'neon' });
    expect(res.valid).toBe(false);
  });

  it('rejects unsafe token properties or illegal chars', () => {
    const res = validateThemeFile({
      ...validTheme,
      tokens: {
        '--malicious': 'red; background: black',
      },
    });
    expect(res.valid).toBe(false);
  });

  it('rejects non-prefixed CSS variables', () => {
    const res = validateThemeFile({
      ...validTheme,
      tokens: {
        background: '#000000',
      },
    });
    expect(res.valid).toBe(false);
  });

  it('downloads theme JSON file with slugified filename', () => {
    const clickSpy = vi.fn();
    const createElementSpy = vi.spyOn(document, 'createElement').mockReturnValue({
      set href(_val: string) {},
      set download(_val: string) {},
      click: clickSpy,
    } as unknown as HTMLAnchorElement);
    const revokeSpy = vi.spyOn(URL, 'revokeObjectURL').mockImplementation(() => {});

    downloadThemeFile(validTheme);

    expect(createElementSpy).toHaveBeenCalledWith('a');
    expect(clickSpy).toHaveBeenCalled();
    expect(revokeSpy).toHaveBeenCalled();

    createElementSpy.mockRestore();
    revokeSpy.mockRestore();
  });

  it('parses and validates a File instance', async () => {
    const file = new File([JSON.stringify(validTheme)], 'theme.json', {
      type: 'application/json',
    });
    const parsed = await parseThemeFile(file);
    expect(parsed.name).toBe('Tokyo Night Custom');
  });

  it('throws on corrupt file content', async () => {
    const file = new File(['not json'], 'theme.json', {
      type: 'application/json',
    });
    await expect(parseThemeFile(file)).rejects.toThrow('Selected file is not valid JSON');
  });
});
