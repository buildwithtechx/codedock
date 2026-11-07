import { FileText, Plus, Trash2, Upload } from 'lucide-react';
import type { ReactNode } from 'react';
import { useRef } from 'react';
import { toast } from 'sonner';
import { parseDotenv } from '#/lib/dotenv';

export type KV = { key: string; value: string };

export const SECTION_TONES: Record<string, { bg: string; text: string }> = {
  basics: { bg: 'bg-sky-500/10', text: 'text-sky-500' },
  schedule: { bg: 'bg-violet-500/10', text: 'text-violet-500' },
  environment: { bg: 'bg-amber-500/10', text: 'text-amber-500' },
  servers: { bg: 'bg-emerald-500/10', text: 'text-emerald-500' },
  reliability: { bg: 'bg-blue-500/10', text: 'text-blue-500' },
  triggers: { bg: 'bg-fuchsia-500/10', text: 'text-fuchsia-500' },
  notifications: { bg: 'bg-rose-500/10', text: 'text-rose-500' },
};

export function Section({
  title,
  icon,
  tone,
  children,
}: {
  title: string;
  icon?: ReactNode;
  tone?: { bg: string; text: string };
  children: ReactNode;
}) {
  const t = tone ?? { bg: 'bg-primary/10', text: 'text-primary' };
  return (
    <section className="rounded-2xl border border-border/50 bg-card">
      <div className="flex items-center gap-3 border-border/50 border-b px-5 py-4">
        {icon && (
          <div className={`flex size-9 shrink-0 items-center justify-center rounded-xl ${t.bg}`}>
            <span className={t.text}>{icon}</span>
          </div>
        )}
        <h2 className="font-semibold text-[15px] text-foreground">{title}</h2>
      </div>
      <div className="space-y-4 p-5">{children}</div>
    </section>
  );
}

export function Field({ label, children }: { label?: string; children: ReactNode }) {
  return (
    <label className="block">
      {label && (
        <span className="mb-1.5 block font-medium text-muted-foreground text-sm">{label}</span>
      )}
      {children}
    </label>
  );
}

export function SummaryRow({
  label,
  value,
  mono,
}: {
  label: string;
  value: string;
  mono?: boolean;
}) {
  return (
    <div className="flex items-center justify-between gap-3">
      <span className="shrink-0 text-muted-foreground text-sm">{label}</span>
      <span
        className={`min-w-0 truncate text-right text-foreground text-sm ${mono ? 'font-mono' : ''}`}
      >
        {value}
      </span>
    </div>
  );
}

export function KeyValueEditor({
  rows,
  setRows,
  addLabel,
  secret,
}: {
  rows: KV[];
  setRows: (updater: (prev: KV[]) => KV[]) => void;
  addLabel: string;
  secret?: boolean;
}) {
  const fileRef = useRef<HTMLInputElement>(null);

  const setItem = (index: number, patch: Partial<KV>) => {
    setRows((prev) => prev.map((r, i) => (i === index ? { ...r, ...patch } : r)));
  };

  const removeItem = (index: number) => {
    setRows((prev) => prev.filter((_, i) => i !== index));
  };

  const importEnv = (text: string) => {
    const parsed = parseDotenv(text);
    if (parsed.length === 0) {
      toast.error('No environment variables found in pasted content');
      return;
    }
    setRows((prev) => {
      const existing = new Map(prev.filter((r) => r.key.trim()).map((r) => [r.key, r]));
      for (const p of parsed) existing.set(p.key, p);
      return Array.from(existing.values());
    });
    toast.success(`Imported ${parsed.length} variable${parsed.length === 1 ? '' : 's'}`);
  };

  return (
    <div className="space-y-2">
      {!secret && (
        <div className="flex items-center justify-between gap-2">
          <span className="font-medium text-muted-foreground text-xs uppercase tracking-wider">
            Variables
          </span>
          <div className="flex items-center gap-1.5">
            <button
              type="button"
              onClick={async () => {
                try {
                  importEnv(await navigator.clipboard.readText());
                } catch {
                  toast.error('Could not access clipboard');
                }
              }}
              className="inline-flex items-center gap-1.5 rounded-lg border border-border/60 px-2.5 py-1 text-muted-foreground text-xs font-medium hover:bg-muted/50 hover:text-foreground"
            >
              <FileText className="size-3.5" /> Paste .env
            </button>
            <button
              type="button"
              onClick={() => fileRef.current?.click()}
              className="inline-flex items-center gap-1.5 rounded-lg border border-border/60 px-2.5 py-1 text-muted-foreground text-xs font-medium hover:bg-muted/50 hover:text-foreground"
            >
              <Upload className="size-3.5" /> Upload .env
            </button>
            <input
              ref={fileRef}
              type="file"
              accept=".env,text/plain"
              className="hidden"
              onChange={(e) => {
                const f = e.target.files?.[0];
                if (f) void f.text().then(importEnv);
                e.target.value = '';
              }}
            />
          </div>
        </div>
      )}
      <div className="space-y-1.5">
        {rows.map((r, i) => (
          <div key={i} className="flex items-center gap-2">
            <input
              value={r.key}
              onChange={(e) => setItem(i, { key: e.target.value })}
              placeholder="KEY"
              className="w-1/2 rounded-xl border border-border/60 bg-background px-3 py-2 font-mono text-foreground text-sm outline-none transition-colors focus:border-primary/50"
            />
            <input
              value={r.value}
              onChange={(e) => setItem(i, { value: e.target.value })}
              placeholder="value"
              type={secret ? 'password' : 'text'}
              className="w-1/2 rounded-xl border border-border/60 bg-background px-3 py-2 font-mono text-foreground text-sm outline-none transition-colors focus:border-primary/50"
            />
            <button
              type="button"
              onClick={() => removeItem(i)}
              className="rounded-lg p-2 text-muted-foreground hover:bg-red-500/10 hover:text-red-500"
              aria-label="Remove item"
            >
              <Trash2 className="size-4" />
            </button>
          </div>
        ))}
        <button
          type="button"
          onClick={() => setRows((prev) => [...prev, { key: '', value: '' }])}
          className="flex w-full items-center justify-center gap-1.5 rounded-xl border border-border/60 border-dashed py-2.5 font-medium text-muted-foreground text-sm transition-colors hover:border-primary/40 hover:bg-muted/30 hover:text-foreground"
        >
          <Plus className="size-4" /> {addLabel}
        </button>
      </div>
    </div>
  );
}
