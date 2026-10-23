export interface CodedockThemeFile {
  version: 1;
  name: string;
  author?: string;
  mode: 'dark' | 'dim' | 'light';
  radius?: string;
  tokens: Record<string, string>;
}

const ILLEGAL_CHARS = [';', '{', '}', '<', '>', '\\'];

export function validateThemeFile(
  data: unknown
): { valid: true; theme: CodedockThemeFile } | { valid: false; error: string } {
  if (!data || typeof data !== 'object') {
    return { valid: false, error: 'Theme file must be a JSON object' };
  }

  const payload = data as Partial<CodedockThemeFile>;

  if (payload.version !== 1) {
    return { valid: false, error: 'Unsupported theme file version (expected version: 1)' };
  }

  if (typeof payload.name !== 'string' || !payload.name.trim() || payload.name.length > 64) {
    return { valid: false, error: 'Theme name must be a string between 1 and 64 characters' };
  }

  if (
    payload.mode !== 'dark' &&
    payload.mode !== 'dim' &&
    payload.mode !== 'light' &&
    payload.mode !== 'system'
  ) {
    return { valid: false, error: 'Theme mode must be one of: "dark", "dim", "light", "system"' };
  }

  if (payload.radius && typeof payload.radius !== 'string') {
    return { valid: false, error: 'Theme radius must be a valid CSS string' };
  }

  if (!payload.tokens || typeof payload.tokens !== 'object' || Array.isArray(payload.tokens)) {
    return { valid: false, error: 'Theme tokens must be an object of CSS key-value pairs' };
  }

  for (const [key, value] of Object.entries(payload.tokens)) {
    if (!key.startsWith('--') || key.length > 64) {
      return { valid: false, error: `Invalid CSS variable name: ${key}` };
    }
    if (typeof value !== 'string' || value.length > 512) {
      return { valid: false, error: `Invalid CSS value for token ${key}` };
    }
    if (
      ILLEGAL_CHARS.some((char) => key.includes(char) || value.includes(char)) ||
      /url\s*\(/i.test(value)
    ) {
      return { valid: false, error: `Unsafe token declaration detected in ${key}` };
    }
  }

  return {
    valid: true,
    theme: {
      version: 1,
      name: payload.name.trim(),
      author: typeof payload.author === 'string' ? payload.author.slice(0, 64) : undefined,
      mode: payload.mode,
      radius: payload.radius,
      tokens: payload.tokens,
    },
  };
}

export function downloadThemeFile(theme: CodedockThemeFile): void {
  const json = JSON.stringify(theme, null, 2);
  const blob = new Blob([json], { type: 'application/json' });
  const url = URL.createObjectURL(blob);
  const anchor = document.createElement('a');
  const safeName =
    theme.name
      .toLowerCase()
      .replace(/[^a-z0-9]+/g, '-')
      .replace(/^-|-$/g, '') || 'theme';
  anchor.href = url;
  anchor.download = `codedock-theme-${safeName}.json`;
  anchor.click();
  URL.revokeObjectURL(url);
}

export async function parseThemeFile(file: File): Promise<CodedockThemeFile> {
  const text = await file.text();
  let json: unknown;
  try {
    json = JSON.parse(text);
  } catch {
    throw new Error('Selected file is not valid JSON');
  }

  const result = validateThemeFile(json);
  if (!result.valid) {
    throw new Error(result.error);
  }
  return result.theme;
}
