import { readFileSync, readdirSync, writeFileSync } from 'node:fs';
import { dirname, join } from 'node:path';
import { fileURLToPath } from 'node:url';
import { parseCatalogEntry } from './generate-app-catalog.mjs';

const here = dirname(fileURLToPath(import.meta.url));
const siteUrl = process.env.VITE_SITE_URL || 'https://codedock.run';
const routeTree = readFileSync(join(here, '..', 'src', 'routeTree.gen.ts'), 'utf8');
const composeDir = join(here, '..', '..', '..', 'internal', 'engine', 'compose');

const paths = new Set();
for (const match of routeTree.matchAll(/'(\/[^']*)'/g)) {
  const path = match[1];
  if (path.includes('$') || path.includes('*') || path.includes('(')) continue;
  paths.add(path);
}
for (const file of readdirSync(composeDir).sort()) {
  if (!file.endsWith('.yaml')) continue;
  const entry = parseCatalogEntry(
    file.replace(/\.yaml$/, ''),
    readFileSync(join(composeDir, file), 'utf8')
  );
  if (entry) paths.add(`/templates/${entry.id}`);
}

const today = new Date().toISOString().slice(0, 10);
const urls = [...paths]
  .sort()
  .map((path) => `  <url><loc>${siteUrl}${path}</loc><lastmod>${today}</lastmod></url>`)
  .join('\n');
const outFile = join(here, '..', 'public', 'sitemap.xml');
writeFileSync(
  outFile,
  `<?xml version="1.0" encoding="UTF-8"?>\n<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">\n${urls}\n</urlset>\n`
);
console.log(`wrote ${paths.size} urls to ${outFile}`);
