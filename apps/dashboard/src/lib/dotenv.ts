export interface EnvEntry {
  key: string;
  value: string;
}

const DOUBLE_QUOTE_ESCAPES: Readonly<Record<string, string>> = {
  a: '\x07',
  b: '\b',
  f: '\f',
  n: '\n',
  r: '\r',
  t: '\t',
  v: '\v',
  '"': '"',
  '\\': '\\',
  $: '$',
};

const ENCODE_ESCAPES: Readonly<Record<string, string>> = Object.fromEntries(
  Object.entries(DOUBLE_QUOTE_ESCAPES).map(([escChar, value]) => [value, `\\${escChar}`])
);

function encodeSpecialChars(value: string): string {
  let result = '';
  for (let i = 0; i < value.length; i++) {
    const char = value[i];
    const escaped = ENCODE_ESCAPES[char];
    result += escaped || char;
  }
  return result;
}

function closingQuote(content: string, start: number, quote: string): number {
  let escaped = false;
  for (let i = start; i < content.length; i++) {
    const char = content[i];
    if (char === quote && !escaped) return i;
    escaped = char === '\\' && !escaped;
  }
  return -1;
}

export function parseDotenv(content: string): EnvEntry[] {
  const entries: EnvEntry[] = [];
  let cursor = content.startsWith('\uFEFF') ? 1 : 0;
  while (cursor < content.length) {
    const newline = content.indexOf('\n', cursor);
    const lineEnd = newline < 0 ? content.length : newline;
    const line = content.slice(cursor, lineEnd);
    const assignment = line.match(
      /^[^\S\r\n]*(?:export[^\S\r\n]+)?([A-Za-z_][A-Za-z0-9_]*)[^\S\r\n]*=[^\S\r\n]*/
    );
    const valueStart = cursor + (assignment?.[0].length ?? 0);
    cursor = lineEnd + 1;
    if (!assignment) continue;

    const key = assignment[1]!;
    const quote = content[valueStart];
    if (quote === '"' || quote === "'") {
      const end = closingQuote(content, valueStart + 1, quote);
      const quoted = content.slice(valueStart + 1, end < 0 ? lineEnd : end);
      if (end >= 0) {
        const nextLine = content.indexOf('\n', end + 1);
        cursor = nextLine < 0 ? content.length : nextLine + 1;
      }
      if (quote === "'") {
        entries.push({ key, value: quoted.replace(/\\'/g, "'") });
      } else {
        const decode = quoted.replace(
          /\\([abfnrtv"\\$])/g,
          (_match, char: string) => DOUBLE_QUOTE_ESCAPES[char] || char
        );
        entries.push({ key, value: decode });
      }
    } else {
      const value = content
        .slice(valueStart, lineEnd)
        .replace(/\s+#.*$/, '')
        .trimEnd();
      entries.push({ key, value });
    }
  }
  return entries;
}

export function serializeDotenv(rows: ReadonlyArray<EnvEntry>): string {
  const keys = new Set<string>();
  const lines = rows
    .map(({ key: rawKey, value }) => {
      const key = rawKey.trim();
      if (!/^[A-Za-z_][A-Za-z0-9_]*$/.test(key) || keys.has(key)) {
        return '';
      }
      keys.add(key);
      const trailingBackslashes = value.match(/\\+$/)?.[0].length ?? 0;
      let encoded: string;
      if (value === value.trim() && !/[\s#$]/.test(value) && !/^["'`]/.test(value)) {
        encoded = value;
      } else if (!value.includes("'") && !value.includes('\r') && trailingBackslashes % 2 === 0) {
        encoded = `'${value}'`;
      } else {
        encoded = `"${encodeSpecialChars(value)}"`;
      }
      return `${key}=${encoded}`;
    })
    .filter(Boolean);
  return lines.length ? `${lines.join('\n')}\n` : '';
}
