export interface FolderEntry {
  path: string;
  file: File;
}

export interface FolderInspection {
  name: string;
  framework: string;
  packageManager: string;
  fileCount: number;
}

type FileWithPath = File & { webkitRelativePath?: string };

function entryPath(file: File, fallbackRoot: string): string {
  const relative = (file as FileWithPath).webkitRelativePath ?? '';
  if (relative.includes('/')) return relative.split('/').slice(1).join('/');
  return `${fallbackRoot}/${file.name}`;
}

export function collectFolderFiles(files: FileList | File[]): FolderEntry[] {
  const list = Array.from(files);
  const first = list[0] as FileWithPath | undefined;
  const root = first?.webkitRelativePath?.split('/')[0] || `folder-${Date.now().toString(36)}`;
  const entries: FolderEntry[] = [];
  for (const file of list) {
    const path = entryPath(file, root);
    const lowered = path.toLowerCase();
    if (!path || path.endsWith('/')) continue;
    if (lowered.startsWith('.git/') || lowered.startsWith('node_modules/')) continue;
    if (lowered.endsWith('.ds_store') || lowered.endsWith('thumbs.db')) continue;
    entries.push({ path, file });
  }
  return entries;
}

export function detectPackageManager(paths: Set<string>): string {
  if (paths.has('bun.lockb') || paths.has('bun.lock')) return 'bun';
  if (paths.has('pnpm-lock.yaml')) return 'pnpm';
  if (paths.has('yarn.lock')) return 'yarn';
  return 'npm';
}

function frameworkFromPackage(pkg: Record<string, unknown>): string {
  const deps = {
    ...((pkg.dependencies as Record<string, string> | undefined) ?? {}),
    ...((pkg.devDependencies as Record<string, string> | undefined) ?? {}),
  };
  if (deps.next) return 'Next.js';
  if (deps.nuxt || deps['@nuxt/cli']) return 'Nuxt';
  if (deps['@remix-run/node']) return 'Remix';
  if (deps.astro) return 'Astro';
  if (deps.vite) return 'Vite';
  if (deps.react) return 'React';
  if (deps.vue) return 'Vue';
  if (deps['@angular/core']) return 'Angular';
  if (deps.svelte) return 'Svelte';
  if (deps.express || deps.fastify || deps.koa || deps.hono) return 'Node';
  if (deps.nestjs || deps['@nestjs/core']) return 'NestJS';
  return 'Node';
}

export async function inspectFolderFiles(entries: FolderEntry[]): Promise<FolderInspection> {
  const paths = new Set(entries.map((entry) => entry.path));
  const packageManager = detectPackageManager(paths);
  const first = entries[0]?.file as FileWithPath | undefined;
  let name = first?.webkitRelativePath?.split('/')[0] || 'app';
  let framework = 'Unknown';
  const manifest = entries.find((entry) => entry.path === 'package.json');
  if (manifest) {
    try {
      const pkg = JSON.parse(await manifest.file.text()) as Record<string, unknown>;
      if (typeof pkg.name === 'string' && pkg.name) name = pkg.name;
      framework = frameworkFromPackage(pkg);
    } catch {
      framework = 'Node';
    }
  } else if (paths.has('go.mod')) {
    framework = 'Go';
  } else if (paths.has('requirements.txt') || paths.has('pyproject.toml')) {
    framework = 'Python';
  } else if (paths.has('Gemfile')) {
    framework = 'Ruby';
  } else if (paths.has('composer.json')) {
    framework = 'PHP';
  } else if (paths.has('Cargo.toml')) {
    framework = 'Rust';
  } else if ([...paths].some((path) => path.endsWith('.csproj'))) {
    framework = '.NET';
  } else if (paths.has('Dockerfile')) {
    framework = 'Docker';
  }
  return { name, framework, packageManager, fileCount: entries.length };
}

function writeString(view: Uint8Array, offset: number, length: number, value: string): void {
  const bytes = new TextEncoder().encode(value);
  view.set(bytes.slice(0, length), offset);
}

function writeOctal(view: Uint8Array, offset: number, length: number, value: number): void {
  const text = `${value.toString(8).padStart(length - 1, '0')}\0`;
  writeString(view, offset, length, text);
}

function tarHeader(path: string, size: number): Uint8Array {
  const header = new Uint8Array(512);
  writeString(header, 0, 100, path.slice(0, 100));
  writeOctal(header, 100, 8, 0o644);
  writeOctal(header, 108, 8, 0);
  writeOctal(header, 116, 8, 0);
  writeOctal(header, 124, 12, size);
  writeOctal(header, 136, 12, Math.floor(Date.now() / 1000));
  header.fill(0x20, 148, 156);
  header[156] = 0x30;
  writeString(header, 257, 6, 'ustar\0');
  writeString(header, 265, 3, '00');
  let checksum = 0;
  for (const byte of header) checksum += byte;
  writeString(header, 148, 8, `${checksum.toString(8).padStart(6, '0')}\0 `);
  return header;
}

function concat(chunks: Uint8Array[]): Uint8Array {
  const total = chunks.reduce((sum, chunk) => sum + chunk.length, 0);
  const out = new Uint8Array(total);
  let offset = 0;
  for (const chunk of chunks) {
    out.set(chunk, offset);
    offset += chunk.length;
  }
  return out;
}

export async function packTarGz(entries: FolderEntry[]): Promise<Blob> {
  const chunks: Uint8Array[] = [];
  for (const entry of entries) {
    const data = new Uint8Array(await entry.file.arrayBuffer());
    chunks.push(tarHeader(entry.path, data.length));
    chunks.push(data);
    const remainder = data.length % 512;
    if (remainder > 0) chunks.push(new Uint8Array(512 - remainder));
  }
  chunks.push(new Uint8Array(1024));
  const tar = concat(chunks);
  const compressed = new Blob([tar as BlobPart])
    .stream()
    .pipeThrough(new CompressionStream('gzip'));
  return new Response(compressed).blob();
}
