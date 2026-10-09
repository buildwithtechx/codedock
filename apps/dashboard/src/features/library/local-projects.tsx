import { FolderOpen, Plus, Trash2 } from 'lucide-react';
import { useState } from 'react';
import { Button } from '#/components/ui/button';
import type { LocalImport } from './types';

const STORAGE_KEY = 'codedock.library.local-imports';

export function readLocalImports(): LocalImport[] {
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    const parsed: unknown = raw ? JSON.parse(raw) : [];
    return Array.isArray(parsed) ? (parsed as LocalImport[]) : [];
  } catch {
    return [];
  }
}

export function rememberLocalImport(entry: Omit<LocalImport, 'id' | 'createdAt'>): void {
  try {
    const imports = readLocalImports();
    imports.unshift({
      ...entry,
      id:
        typeof crypto !== 'undefined' && 'randomUUID' in crypto
          ? crypto.randomUUID()
          : String(Date.now()),
      createdAt: new Date().toISOString(),
    });
    localStorage.setItem(STORAGE_KEY, JSON.stringify(imports.slice(0, 20)));
  } catch {
    return;
  }
}

export function LocalProjects({ onImportFolder }: { onImportFolder: () => void }) {
  const [imports, setImports] = useState<LocalImport[]>(() => readLocalImports());

  const handleDelete = (id: string) => {
    const next = imports.filter((entry) => entry.id !== id);
    setImports(next);
    try {
      localStorage.setItem(STORAGE_KEY, JSON.stringify(next));
    } catch {
      return;
    }
  };

  return (
    <div className="rounded-2xl border border-border/50 bg-card">
      <div className="flex items-center justify-between border-border/50 border-b px-5 py-4">
        <div className="flex items-center gap-3">
          <div className="flex h-9 w-9 items-center justify-center rounded-xl bg-primary/10">
            <FolderOpen className="size-[18px] text-primary" />
          </div>
          <div>
            <h2 className="font-semibold text-[15px] text-foreground">Local projects</h2>
            <p className="text-muted-foreground text-xs">
              {imports.length === 0
                ? 'Folders you imported from this browser'
                : `${imports.length} remembered import${imports.length === 1 ? '' : 's'}`}
            </p>
          </div>
        </div>
        <Button size="sm" onClick={onImportFolder}>
          <Plus className="size-3.5" />
          Import
        </Button>
      </div>

      {imports.length === 0 ? (
        <div className="px-6 py-12 text-center">
          <h3 className="font-medium text-foreground/80 text-lg">No local imports yet</h3>
          <p className="mx-auto mt-1 max-w-sm text-muted-foreground text-sm leading-relaxed">
            Pick a folder above to pack and deploy it. Your recent imports are remembered here for
            quick reference.
          </p>
        </div>
      ) : (
        <div className="divide-y divide-border/50">
          {imports.map((entry) => (
            <div
              key={entry.id}
              className="group flex items-center gap-4 px-5 py-3.5 transition-colors hover:bg-muted/40"
            >
              <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-muted/60 transition-colors group-hover:bg-muted">
                <FolderOpen className="size-[18px] text-muted-foreground" />
              </div>
              <div className="min-w-0 flex-1">
                <p className="truncate font-medium text-foreground text-sm">{entry.name}</p>
                <div className="mt-0.5 flex items-center gap-2">
                  <p className="truncate font-mono text-muted-foreground text-xs">{entry.root}</p>
                  <span className="text-muted-foreground/40">·</span>
                  <span className="shrink-0 text-muted-foreground text-xs">
                    {entry.framework} · {entry.fileCount} files
                  </span>
                </div>
              </div>
              <button
                type="button"
                onClick={() => handleDelete(entry.id)}
                aria-label={`Forget ${entry.name}`}
                className="rounded-lg p-1.5 text-muted-foreground/50 opacity-0 transition-all hover:bg-destructive/10 hover:text-destructive group-hover:opacity-100"
              >
                <Trash2 className="size-3.5" />
              </button>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
