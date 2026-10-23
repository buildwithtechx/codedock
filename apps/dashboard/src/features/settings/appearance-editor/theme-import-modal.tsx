import { ArrowUpToLine, FileJson, UploadCloud } from 'lucide-react';
import { useRef, useState } from 'react';
import { toast } from 'sonner';
import { Button } from '@codedock/ui/button';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@codedock/ui/dialog';
import { Textarea } from '@codedock/ui/textarea';
import type { CustomTheme } from '#/interfaces/appearance';
import { parseThemeFile, validateThemeFile } from '#/lib/theme-export-import';

interface ThemeImportModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onImport: (theme: CustomTheme) => void;
}

export function ThemeImportModal({ open, onOpenChange, onImport }: ThemeImportModalProps) {
  const [jsonText, setJsonText] = useState('');
  const [error, setError] = useState<string | null>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const handleApplyJson = () => {
    setError(null);
    if (!jsonText.trim()) {
      setError('Please paste theme JSON or upload a file');
      return;
    }

    try {
      const parsed = JSON.parse(jsonText);
      const res = validateThemeFile(parsed);
      if (!res.valid) {
        setError(res.error);
        return;
      }
      onImport({
        id: `theme-${Date.now()}`,
        name: res.theme.name,
        mode: res.theme.mode,
        tokens: res.theme.tokens,
      });
      toast.success(`Theme "${res.theme.name}" imported successfully`);
      onOpenChange(false);
      setJsonText('');
    } catch {
      setError('Invalid JSON syntax');
    }
  };

  const handleFileUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    try {
      const parsed = await parseThemeFile(file);
      onImport({
        id: `theme-${Date.now()}`,
        name: parsed.name,
        mode: parsed.mode,
        tokens: parsed.tokens,
      });
      toast.success(`Imported theme "${parsed.name}"`);
      onOpenChange(false);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to parse file');
    } finally {
      if (fileInputRef.current) fileInputRef.current.value = '';
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <div className="flex items-center gap-2">
            <UploadCloud className="size-5 text-primary" />
            <DialogTitle>Import Theme</DialogTitle>
          </div>
          <DialogDescription>
            Upload a <code className="font-mono text-xs">codedock-theme.json</code> file or paste
            theme JSON tokens below.
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-4 py-2">
          <input
            ref={fileInputRef}
            type="file"
            accept=".json,application/json"
            onChange={handleFileUpload}
            className="hidden"
          />

          <button
            type="button"
            onClick={() => fileInputRef.current?.click()}
            className="flex w-full flex-col items-center justify-center rounded-xl border border-dashed border-border/80 bg-muted/10 p-6 transition-colors hover:border-primary/50 hover:bg-muted/20"
          >
            <FileJson className="size-8 text-muted-foreground mb-2" />
            <span className="font-medium text-foreground text-xs">Click to browse theme JSON</span>
            <span className="text-[11px] text-muted-foreground mt-1">
              Supports Codedock and OpenShip theme formats
            </span>
          </button>

          <div className="relative">
            <div className="absolute inset-0 flex items-center">
              <span className="w-full border-t border-border/60" />
            </div>
            <div className="relative flex justify-center text-[10px] uppercase">
              <span className="bg-popover px-2 text-muted-foreground">Or paste JSON</span>
            </div>
          </div>

          <Textarea
            value={jsonText}
            onChange={(e) => {
              setJsonText(e.target.value);
              setError(null);
            }}
            placeholder='{\n  "version": 1,\n  "name": "Custom Palette",\n  "mode": "dark",\n  "tokens": { ... }\n}'
            className="h-28 font-mono text-xs"
          />

          {error && <p className="text-destructive text-xs">{error}</p>}
        </div>

        <DialogFooter className="gap-2 sm:gap-0">
          <Button variant="ghost" onClick={() => onOpenChange(false)} size="sm">
            Cancel
          </Button>
          <Button onClick={handleApplyJson} size="sm" className="gap-1.5">
            <ArrowUpToLine className="size-3.5" />
            <span>Apply Theme</span>
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
