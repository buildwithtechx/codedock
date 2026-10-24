import type JSZip from 'jszip';

export const MAX_ROWS = 500;
export const MAX_COLUMNS = 50;
export const MAX_SHEETS = 20;
export const MAX_XML_BYTES = 8 * 1024 * 1024;
export const MAX_DEPTH = 256;

export const INCOMPLETE_NOTE =
  '<p class="note">This preview is incomplete. Download the file to see everything.</p>';

export const STYLES = `
body { font: 14px/1.5 system-ui, sans-serif; margin: 16px; color: #111; }
h2 { font-size: 16px; margin: 24px 0 8px; }
table { border-collapse: collapse; margin-bottom: 8px; width: 100%; }
td { border: 1px solid #d4d4d8; padding: 4px 8px; white-space: pre-wrap; vertical-align: top; }
tr:first-child td { background: #f4f4f5; font-weight: 600; }
.note { color: #71717a; font-size: 13px; margin-top: 12px; font-style: italic; }
@media (prefers-color-scheme: dark) {
  body { background: #111; color: #eee; }
  td { border-color: #3f3f46; }
  tr:first-child td { background: #27272a; }
}
`;

export const escapeHtml = (value: string): string =>
  value
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#39;');

const XML_ENTITIES = new Map([
  ['lt', '<'],
  ['gt', '>'],
  ['quot', '"'],
  ['apos', "'"],
  ['amp', '&'],
]);

export const decodeXml = (value: string): string =>
  value.replace(/&(#x[0-9a-f]+|#\d+|\w+);/gi, (entity, body: string) => {
    if (!body.startsWith('#')) return XML_ENTITIES.get(body) ?? entity;
    const code = /^#x/i.test(body) ? parseInt(body.slice(2), 16) : parseInt(body.slice(1), 10);
    return code <= 0x10ffff ? String.fromCodePoint(code) : '\uFFFD';
  });

export interface XmlNode {
  name: string;
  attributes: Map<string, string>;
  children: (XmlNode | string)[];
}

const XML_TOKEN = /<(\/?)([^\s/<>]+)((?:\s[^<>]*?)?)(\/?)>|([^<]+)/g;
const XML_ATTRIBUTE = /([^\s=]+)\s*=\s*(?:"([^"]*)"|'([^']*)')/g;

export const localName = (name: string): string => name.slice(name.indexOf(':') + 1);

export const parseXml = (xml: string): XmlNode => {
  const root: XmlNode = { name: '', attributes: new Map(), children: [] };
  const stack: XmlNode[] = [root];
  for (const match of xml.matchAll(XML_TOKEN)) {
    const [, closing, rawName, rawAttributes, selfClosing, text] = match;
    const parent = stack.at(-1) ?? root;
    if (text !== undefined) {
      parent.children.push(decodeXml(text));
      continue;
    }
    if (!rawName || rawName.startsWith('?') || rawName.startsWith('!')) continue;
    const name = localName(rawName);
    if (closing) {
      let index = -1;
      for (let i = stack.length - 1; i >= 0; i--) {
        if (stack[i].name === name) {
          index = i;
          break;
        }
      }
      if (index > 0) stack.length = index;
      continue;
    }
    const node: XmlNode = {
      name,
      attributes: new Map(
        [...(rawAttributes ?? '').matchAll(XML_ATTRIBUTE)].map(
          ([, key = '', doubleQuoted, singleQuoted]) => [
            localName(key),
            decodeXml(doubleQuoted ?? singleQuoted ?? ''),
          ]
        )
      ),
      children: [],
    };
    parent.children.push(node);
    if (!selfClosing && stack.length < MAX_DEPTH) stack.push(node);
  }
  return root;
};

export const elements = (node: XmlNode): XmlNode[] =>
  node.children.filter((child): child is XmlNode => typeof child !== 'string');

export const descendants = (node: XmlNode, name: string): XmlNode[] =>
  elements(node).flatMap((child) => (child.name === name ? [child] : descendants(child, name)));

export const plainText = (node: XmlNode): string =>
  node.children.map((child) => (typeof child === 'string' ? child : plainText(child))).join('');

const SKIPPED_TEXT_ELEMENTS = new Set(['pPr', 'rPr', 'rPh', 'Fallback']);

export const textOf = (node: XmlNode): string =>
  elements(node)
    .map((child) => {
      if (child.name === 't') return plainText(child);
      if (child.name === 'tab') return '\t';
      if (child.name === 'br' || child.name === 'cr') return '\n';
      if (SKIPPED_TEXT_ELEMENTS.has(child.name)) return '';
      return textOf(child);
    })
    .join('');

export type StreamableZipObject = JSZip.JSZipObject & {
  internalStream: (type: 'uint8array') => JSZip.JSZipStreamHelper<Uint8Array>;
};

export const readCapped = (
  file: StreamableZipObject
): Promise<{ chunks: Uint8Array[]; truncated: boolean }> =>
  new Promise((resolve, reject) => {
    const chunks: Uint8Array[] = [];
    let size = 0;
    const stream = file.internalStream('uint8array');
    stream
      .on('data', (chunk) => {
        if (size >= MAX_XML_BYTES) return;
        chunks.push(chunk.subarray(0, MAX_XML_BYTES - size));
        size += chunk.length;
        if (size >= MAX_XML_BYTES) {
          stream.pause();
          resolve({ chunks, truncated: true });
        }
      })
      .on('error', reject)
      .on('end', () => resolve({ chunks, truncated: false }))
      .resume();
  });

export const readXml = async (
  zip: JSZip,
  path: string
): Promise<{ root: XmlNode; truncated: boolean }> => {
  const file = zip.file(path) as StreamableZipObject | null;
  if (!file) return { root: parseXml(''), truncated: false };
  const { chunks, truncated } = await readCapped(file);
  const decoder = new TextDecoder();
  const xml = chunks.map((chunk) => decoder.decode(chunk, { stream: true })).join('');
  return { root: parseXml(xml + decoder.decode()), truncated };
};

export const columnIndex = (reference: string): number => {
  let index = 0;
  for (const char of reference.replace(/[^A-Z]/gi, '').toUpperCase()) {
    index = index * 26 + char.charCodeAt(0) - 64;
  }
  return index - 1;
};

export const renderTable = (rows: string[][]): string =>
  `<table>${rows
    .map((row) => `<tr>${row.map((cell) => `<td>${escapeHtml(cell)}</td>`).join('')}</tr>`)
    .join('')}</table>`;
