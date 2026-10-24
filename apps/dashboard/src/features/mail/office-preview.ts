import JSZip from 'jszip';
import {
  INCOMPLETE_NOTE,
  MAX_COLUMNS,
  MAX_ROWS,
  MAX_SHEETS,
  STYLES,
  type XmlNode,
  columnIndex,
  descendants,
  elements,
  escapeHtml,
  plainText,
  readXml,
  renderTable,
  textOf,
} from './office-preview-xml';

export { escapeHtml } from './office-preview-xml';

const XLSX_TYPES = [
  'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
  'application/vnd.ms-excel',
];
const DOCX_TYPES = [
  'application/vnd.openxmlformats-officedocument.wordprocessingml.document',
  'application/msword',
];
const CSV_TYPES = ['text/csv', 'application/csv'];

export type OfficePreviewKind = 'xlsx' | 'docx' | 'csv';

export const getOfficePreviewKind = (
  mimeType: string,
  filename: string
): OfficePreviewKind | null => {
  const type = mimeType.toLowerCase().split(';')[0]?.trim() ?? '';
  const extension = filename.toLowerCase().split('.').pop() ?? '';
  if (XLSX_TYPES.includes(type) || extension === 'xlsx') return 'xlsx';
  if (DOCX_TYPES.includes(type) || extension === 'docx') return 'docx';
  if (CSV_TYPES.includes(type) || extension === 'csv') return 'csv';
  return null;
};

const parseXlsx = async (data: Uint8Array): Promise<string> => {
  const zip = await JSZip.loadAsync(data);
  const workbook = await readXml(zip, 'xl/workbook.xml');
  const relationships = await readXml(zip, 'xl/_rels/workbook.xml.rels');
  const sharedStringsXml = await readXml(zip, 'xl/sharedStrings.xml');
  const sharedStrings = descendants(sharedStringsXml.root, 'si').map(textOf);
  const targets = new Map(
    descendants(relationships.root, 'Relationship').map((relationship) => [
      relationship.attributes.get('Id') ?? '',
      relationship.attributes.get('Target') ?? '',
    ])
  );
  const sheets = descendants(workbook.root, 'sheet');
  let incomplete =
    workbook.truncated ||
    relationships.truncated ||
    sharedStringsXml.truncated ||
    sheets.length > MAX_SHEETS;

  const rendered: string[] = [];
  for (const sheet of sheets.slice(0, MAX_SHEETS)) {
    const target = targets.get(sheet.attributes.get('id') ?? '') ?? '';
    const path = target.startsWith('/') ? target.slice(1) : `xl/${target.replace(/^\.?\//, '')}`;
    const { root, truncated } = await readXml(zip, path);
    incomplete = incomplete || truncated;
    const rows: string[][] = [];

    for (const row of descendants(root, 'row')) {
      const rowNumber = Number(row.attributes.get('r'));
      while (Number.isInteger(rowNumber) && rows.length < Math.min(rowNumber - 1, MAX_ROWS)) {
        rows.push([]);
      }
      if (rows.length >= MAX_ROWS) {
        incomplete = true;
        break;
      }
      const cells: string[] = [];
      let column = -1;
      for (const cell of elements(row).filter((child) => child.name === 'c')) {
        const reference = cell.attributes.get('r');
        column = reference ? columnIndex(reference) : column + 1;
        if (column < 0) continue;
        if (column >= MAX_COLUMNS) {
          incomplete = true;
          break;
        }
        const type = cell.attributes.get('t');
        const value = descendants(cell, 'v').map(plainText).join('');
        let text = value;
        if (type === 's') text = sharedStrings[Number(value)] ?? '';
        else if (type === 'inlineStr') text = descendants(cell, 'is').map(textOf).join('');
        else if (type === 'b') text = value === '1' ? 'TRUE' : 'FALSE';
        while (cells.length < column) cells.push('');
        cells[column] = text;
      }
      rows.push(cells);
    }

    rendered.push(
      `<h2>${escapeHtml(sheet.attributes.get('name') ?? 'Sheet')}</h2>${renderTable(rows)}`
    );
  }
  return (
    (rendered.join('') || '<p class="note">This workbook has no sheets.</p>') +
    (incomplete ? INCOMPLETE_NOTE : '')
  );
};

const cellText = (cell: XmlNode): string => descendants(cell, 'p').map(textOf).join('\n');

const renderDocxBlocks = (node: XmlNode): string =>
  elements(node)
    .map((child) => {
      if (child.name === 'tbl') {
        return renderTable(
          descendants(child, 'tr').map((row) => descendants(row, 'tc').map(cellText))
        );
      }
      if (child.name !== 'p') return renderDocxBlocks(child);
      const text = textOf(child);
      if (!text.trim()) return '';
      const style = descendants(child, 'pStyle')[0]?.attributes.get('val') ?? '';
      const tag = /^(heading|title)/i.test(style) ? 'h2' : 'p';
      return `<${tag}>${escapeHtml(text)}</${tag}>`;
    })
    .join('');

const parseDocx = async (data: Uint8Array): Promise<string> => {
  const zip = await JSZip.loadAsync(data);
  const { root, truncated } = await readXml(zip, 'word/document.xml');
  return (
    (renderDocxBlocks(root) || '<p class="note">This document has no text.</p>') +
    (truncated ? INCOMPLETE_NOTE : '')
  );
};

const parseCsv = (text: string): string => {
  const rows: string[][] = [];
  let row: string[] = [];
  let cell = '';
  let quoted = false;
  for (let i = 0; i < text.length; i++) {
    const char = text[i];
    if (quoted) {
      if (char === '"' && text[i + 1] === '"') {
        cell += '"';
        i++;
      } else if (char === '"') quoted = false;
      else cell += char;
    } else if (char === '"') quoted = true;
    else if (char === ',') {
      row.push(cell);
      cell = '';
    } else if (char === '\n' || char === '\r') {
      if (char === '\r' && text[i + 1] === '\n') i++;
      row.push(cell);
      rows.push(row);
      row = [];
      cell = '';
    } else cell += char;
    if (rows.length > MAX_ROWS) break;
  }
  if (cell || row.length) rows.push([...row, cell]);
  const incomplete = rows.length > MAX_ROWS || rows.some((r) => r.length > MAX_COLUMNS);
  return (
    renderTable(rows.slice(0, MAX_ROWS).map((r) => r.slice(0, MAX_COLUMNS))) +
    (incomplete ? INCOMPLETE_NOTE : '')
  );
};

/**
 * Builds a self-contained, script-free HTML page for an Office/CSV attachment.
 * All dynamic content is strictly escaped. Sandboxed CSP header applied.
 */
export const renderOfficePreview = async (
  kind: OfficePreviewKind,
  data: Uint8Array,
  title: string
): Promise<string> => {
  const body =
    kind === 'xlsx'
      ? await parseXlsx(data)
      : kind === 'docx'
        ? await parseDocx(data)
        : parseCsv(new TextDecoder().decode(data));

  return `<!doctype html><html><head><meta charset="utf-8"><meta http-equiv="Content-Security-Policy" content="default-src 'none'; style-src 'unsafe-inline'"><title>${escapeHtml(
    title
  )}</title><style>${STYLES}</style></head><body>${body}</body></html>`;
};
