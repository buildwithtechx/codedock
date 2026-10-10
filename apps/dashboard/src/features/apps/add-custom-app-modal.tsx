import { AlertCircle, CheckCircle2, FileUp, Loader2, Upload } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '#/components/ui/dialog';
import { appsApi } from './api';

interface AddCustomAppModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onAdded?: () => void;
}

export function AddCustomAppModal({ open, onOpenChange, onAdded }: AddCustomAppModalProps) {
  const [template, setTemplate] = useState<unknown>(null);
  const [name, setName] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [busy, setBusy] = useState(false);

  const evaluate = (text: string) => {
    setTemplate(null);
    setName(null);
    setError(null);
    if (!text.trim()) return;
    try {
      const obj = JSON.parse(text) as Record<string, unknown>;
      const parsedName =
        typeof obj.name === 'string'
          ? obj.name
          : typeof obj['x-codedock'] === 'object' &&
              obj['x-codedock'] !== null &&
              typeof (obj['x-codedock'] as Record<string, unknown>).name === 'string'
            ? ((obj['x-codedock'] as Record<string, unknown>).name as string)
            : typeof obj.id === 'string'
              ? obj.id
              : 'custom-app';
      setTemplate(obj);
      setName(parsedName);
    } catch {
      setError("That file isn't valid JSON.");
    }
  };

  const handleFile = async (file?: File) => {
    if (!file) return;
    try {
      const text = await file.text();
      evaluate(text);
    } catch {
      setError('Could not read the uploaded file.');
    }
  };

  const handleSubmit = async () => {
    if (!template || busy) return;
    setBusy(true);
    try {
      await appsApi.addCustom(template);
      toast.success(`Added "${name}" to your catalog.`);
      onAdded?.();
      onOpenChange(false);
      setTemplate(null);
      setName(null);
      setError(null);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Couldn't add the app.");
    } finally {
      setBusy(false);
    }
  };

  const handleClose = () => {
    if (busy) return;
    onOpenChange(false);
    setTemplate(null);
    setName(null);
    setError(null);
  };

  return (
    <Dialog open={open} onOpenChange={handleClose}>
      <DialogContent className="sm:max-w-130">
        <DialogHeader>
          <DialogTitle>Add a custom app</DialogTitle>
          <DialogDescription>
            Upload an app definition (JSON). It's added to your catalog and marked unverified.{' '}
            <a
              href="https://docs.codedock.run/apps/custom"
              target="_blank"
              rel="noopener noreferrer"
              className="text-primary hover:underline"
            >
              How to write one
            </a>
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-4 py-2">
          <label
            onDragOver={(e) => e.preventDefault()}
            onDrop={(e) => {
              e.preventDefault();
              void handleFile(e.dataTransfer.files?.[0]);
            }}
            className="flex cursor-pointer flex-col items-center justify-center gap-2 rounded-xl border border-border/60 border-dashed bg-muted/20 px-4 py-8 text-center transition-colors hover:border-primary/40 hover:bg-muted/30"
          >
            <Upload className="size-6 text-muted-foreground" />
            <span className="text-foreground text-sm">
              Drop <code className="font-mono text-xs">app.json</code> here, or{' '}
              <span className="text-primary underline">browse</span>
            </span>
            <input
              type="file"
              accept="application/json,.json"
              className="hidden"
              onChange={(e) => void handleFile(e.target.files?.[0])}
            />
          </label>

          {error && (
            <div className="flex items-start gap-2 rounded-xl border border-destructive/40 bg-destructive/10 px-3.5 py-2.5 text-destructive text-sm">
              <AlertCircle className="mt-0.5 size-4 shrink-0" />
              <span>{error}</span>
            </div>
          )}

          {template !== null && !error && (
            <div className="space-y-2">
              <div className="flex items-center gap-2 rounded-xl border border-emerald-500/40 bg-emerald-500/10 px-3.5 py-2.5 text-emerald-600 text-sm dark:text-emerald-400">
                <CheckCircle2 className="size-4 shrink-0" />
                <span>
                  Valid app definition — <span className="font-medium">{name}</span>
                </span>
              </div>
              <div className="flex items-start gap-2 rounded-xl border border-amber-500/40 bg-amber-500/10 px-3.5 py-2.5 text-amber-600 text-xs dark:text-amber-400">
                <AlertCircle className="mt-0.5 size-4 shrink-0" />
                <span>
                  <strong className="font-semibold">Unverified.</strong> This deploys images you
                  provided — not an official, reviewed app. Review the definition and only add apps
                  you trust.
                </span>
              </div>
            </div>
          )}
        </div>

        <DialogFooter className="gap-2 sm:gap-0">
          <Button variant="ghost" onClick={handleClose} disabled={busy}>
            Cancel
          </Button>
          <Button onClick={handleSubmit} disabled={!template || busy} className="gap-2">
            {busy ? (
              <>
                <Loader2 className="size-4 animate-spin" />
                Adding…
              </>
            ) : (
              <>
                <FileUp className="size-4" />
                Add app
              </>
            )}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
