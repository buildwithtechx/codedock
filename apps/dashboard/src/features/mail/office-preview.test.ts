import JSZip from 'jszip';
import { describe, expect, it } from 'vitest';
import { getOfficePreviewKind, renderOfficePreview } from './office-preview';

const zipBytes = async (files: Record<string, string>): Promise<Uint8Array> => {
  const zip = new JSZip();
  for (const [path, content] of Object.entries(files)) {
    zip.file(path, content);
  }
  return zip.generateAsync({ type: 'uint8array', compression: 'DEFLATE' });
};

describe('getOfficePreviewKind', () => {
  it('detects by mime type or extension and ignores legacy binary formats', () => {
    expect(getOfficePreviewKind('application/octet-stream', 'Report.XLSX')).toBe('xlsx');
    expect(getOfficePreviewKind('text/csv; charset=utf-8', 'a')).toBe('csv');
    expect(getOfficePreviewKind('application/msword', 'old.doc')).toBe('docx');
    expect(getOfficePreviewKind('application/pdf', 'file.pdf')).toBeNull();
  });
});

describe('renderOfficePreview', () => {
  it('renders every xlsx sheet with shared strings, gaps and escaped markup', async () => {
    const data = await zipBytes({
      'xl/workbook.xml':
        '<workbook><sheets><sheet name="Sales &amp; Tax" sheetId="1" r:id="rId1"/></sheets></workbook>',
      'xl/_rels/workbook.xml.rels':
        '<Relationships><Relationship Id="rId1" Target="worksheets/sheet1.xml"/></Relationships>',
      'xl/sharedStrings.xml': '<sst><si><t>Item</t></si><si><t>&lt;b&gt;x&lt;/b&gt;</t></si></sst>',
      'xl/worksheets/sheet1.xml':
        '<worksheet><sheetData><row r="1"><c r="A1" t="s"><v>0</v></c><c r="C1"><v>42</v></c></row><row r="2"><c r="A2" t="s"><v>1</v></c></row></sheetData></worksheet>',
    });
    const html = await renderOfficePreview('xlsx', data, 'a.xlsx');
    expect(html).toContain('<h2>Sales &amp; Tax</h2>');
    expect(html).toContain('<tr><td>Item</td><td></td><td>42</td></tr>');
    expect(html).toContain('&lt;b&gt;x&lt;/b&gt;');
    expect(html).not.toContain('<b>');
  });

  it('renders docx paragraphs and tables', async () => {
    const data = await zipBytes({
      'word/document.xml':
        '<w:document><w:body><w:p><w:pPr><w:pStyle w:val="Heading1"/></w:pPr><w:r><w:t>Title</w:t></w:r></w:p><w:p><w:r><w:t xml:space="preserve">Hello </w:t></w:r><w:r><w:t>&lt;i&gt;</w:t></w:r></w:p><w:tbl><w:tr><w:tc><w:p><w:r><w:t>A</w:t></w:r></w:p></w:tc><w:tc><w:p><w:r><w:t>B</w:t></w:r></w:p></w:tc></w:tr></w:tbl></w:body></w:document>',
    });
    const html = await renderOfficePreview('docx', data, 'a.docx');
    expect(html).toContain('<h2>Title</h2>');
    expect(html).toContain('<p>Hello &lt;i&gt;</p>');
    expect(html).toContain('<table><tr><td>A</td><td>B</td></tr></table>');
  });

  it('renders csv with quotes and caps rows', async () => {
    const csv = `a,"b,""c"""\n${Array.from({ length: 600 }, (_, i) => `${i},x`).join('\n')}`;
    const html = await renderOfficePreview('csv', new TextEncoder().encode(csv), 'a.csv');
    expect(html).toContain('<td>b,&quot;c&quot;</td>');
    expect(html).toContain('preview is incomplete');
  });

  it('reads cells without references, phonetic guides and prefixed namespaces', async () => {
    const data = await zipBytes({
      'xl/workbook.xml':
        '<x:workbook><x:sheets><x:sheet name="S" r:id="rId1"/></x:sheets></x:workbook>',
      'xl/_rels/workbook.xml.rels':
        '<Relationships><Relationship Id="rId1" Target="/xl/worksheets/sheet1.xml"/></Relationships>',
      'xl/sharedStrings.xml':
        '<x:sst><x:si><x:t>漢字</x:t><x:rPh sb="0" eb="2"><x:t>カンジ</x:t></x:rPh></x:si></x:sst>',
      'xl/worksheets/sheet1.xml':
        '<x:worksheet><x:sheetData><x:row><x:c t="s"><x:v>0</x:v></x:c><x:c><x:v>8</x:v></x:c></x:row><x:row r="3"><x:c t="inlineStr"><x:is><x:t>&#x110000;</x:t></x:is></x:c></x:row></x:sheetData></x:worksheet>',
    });
    const html = await renderOfficePreview('xlsx', data, 'a.xlsx');
    expect(html).toContain('<tr><td>漢字</td><td>8</td></tr><tr></tr><tr><td>\uFFFD</td></tr>');
  });

  it('keeps textbox text once and ignores tab stop definitions', async () => {
    const data = await zipBytes({
      'word/document.xml':
        '<w:document><w:body><w:p><w:pPr><w:tabs><w:tab w:val="left"/></w:tabs></w:pPr><w:r><w:t>Before </w:t></w:r><w:r><mc:AlternateContent><mc:Choice><w:txbxContent><w:p><w:r><w:t>Box</w:t></w:r></w:p></w:txbxContent></mc:Choice><mc:Fallback><w:txbxContent><w:p><w:r><w:t>Box</w:t></w:r></w:p></w:txbxContent></mc:Fallback></mc:AlternateContent></w:r><w:r><w:t xml:space="preserve"> After</w:t></w:r></w:p></w:body></w:document>',
    });
    expect(await renderOfficePreview('docx', data, 'a.docx')).toContain('<p>Before Box After</p>');
  });

  it('stays fast on unclosed tags and stops reading oversized entries', async () => {
    const unclosed = await zipBytes({
      'word/document.xml': '<w:p>'.repeat(5000) + '<'.repeat(20000) + 'x'.repeat(1000000),
    });
    const start = performance.now();
    await renderOfficePreview('docx', unclosed, 'a.docx');
    expect(performance.now() - start).toBeLessThan(1000);

    const bomb = await zipBytes({
      'word/document.xml': `<w:p><w:r><w:t>${'x'.repeat(64 * 1024 * 1024)}</w:t></w:r></w:p>`,
    });
    const bombStart = performance.now();
    const html = await renderOfficePreview('docx', bomb, 'a.docx');
    expect(html).toContain('preview is incomplete');
    expect(performance.now() - bombStart).toBeLessThan(3000);
  });
});
