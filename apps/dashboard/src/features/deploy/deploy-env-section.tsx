import { ChevronDown, ChevronUp, ClipboardPaste, Download, Key, Plus, Upload } from 'lucide-react';
import { type ChangeEvent, useId, useRef, useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';

interface DeployEnvSectionProps {
  envVars: string;
  onEnvVarsChange: (v: string) => void;
}

export function DeployEnvSection({ envVars, onEnvVarsChange }: DeployEnvSectionProps) {
  const [open, setOpen] = useState(false);
  const [isEditing, setIsEditing] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);
  const inputId = useId();

  const handleUploadClick = () => {
    fileInputRef.current?.click();
  };

  const handleFileUpload = (e: ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    const reader = new FileReader();
    reader.onload = (event) => {
      const content = (event.target?.result as string) || '';
      if (!content.trim()) {
        toast.info('Uploaded .env file is empty');
        return;
      }
      onEnvVarsChange(envVars.trim() ? `${envVars.trim()}\n${content.trim()}` : content.trim());
      setOpen(true);
      setIsEditing(true);
      toast.success(`Imported ${file.name}`);
    };
    reader.onerror = () => {
      toast.error('Failed to read .env file');
    };
    reader.readAsText(file);
    e.target.value = '';
  };

  const handleDownload = () => {
    if (!envVars.trim()) {
      toast.info('No environment variables to download');
      return;
    }
    const blob = new Blob([envVars], { type: 'text/plain;charset=utf-8' });
    const url = URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = url;
    link.download = '.env';
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    URL.revokeObjectURL(url);
    toast.success('Downloaded .env file');
  };

  const lineCount = envVars.split('\n').filter((l) => l.trim() && !l.trim().startsWith('#')).length;
  const hasVars = lineCount > 0;

  return (
    <section className="rounded-2xl border border-border/60 bg-card p-6">
      <input
        ref={fileInputRef}
        type="file"
        accept=".env,text/plain"
        onChange={handleFileUpload}
        className="hidden"
        aria-label="Upload .env file"
      />

      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div className="flex items-center gap-3.5">
          <div className="flex size-10 shrink-0 items-center justify-center rounded-xl bg-muted text-muted-foreground">
            <Key className="size-5" />
          </div>
          <div>
            <h2 className="font-semibold text-foreground text-sm">Environment Variables</h2>
            <p className="text-muted-foreground text-xs">
              {hasVars ? `${lineCount} set` : 'None set'}
            </p>
          </div>
        </div>

        <div className="flex flex-wrap items-center gap-2">
          <Button
            type="button"
            variant="outline"
            size="sm"
            className="h-8 gap-1.5 text-xs"
            onClick={() => {
              setOpen(true);
              setIsEditing(true);
            }}
          >
            <ClipboardPaste className="size-3.5" />
            Paste .env
          </Button>

          <Button
            type="button"
            variant="outline"
            size="sm"
            className="h-8 gap-1.5 text-xs"
            onClick={handleUploadClick}
          >
            <Upload className="size-3.5" />
            Upload .env
          </Button>

          <Button
            type="button"
            variant="outline"
            size="sm"
            disabled={!hasVars}
            className="h-8 gap-1.5 text-xs"
            onClick={handleDownload}
          >
            <Download className="size-3.5" />
            Download .env
          </Button>

          <button
            type="button"
            onClick={() => setOpen((prev) => !prev)}
            className="p-1 text-muted-foreground transition-colors hover:text-foreground"
            aria-label={open ? 'Collapse environment variables' : 'Expand environment variables'}
          >
            {open ? <ChevronUp className="size-4" /> : <ChevronDown className="size-4" />}
          </button>
        </div>
      </div>

      {open && (
        <div className="mt-5 border-border/50 border-t pt-5">
          {!hasVars && !isEditing ? (
            <div className="flex flex-col items-center justify-center rounded-xl border border-border/60 border-dashed bg-muted/10 py-10 text-center">
              <div className="flex size-10 items-center justify-center rounded-full bg-muted/60 text-muted-foreground">
                <Key className="size-5" />
              </div>
              <h3 className="mt-3 font-medium text-foreground text-sm">No environment variables</h3>
              <p className="mt-1 max-w-sm text-muted-foreground text-xs leading-relaxed">
                Paste a .env block anywhere in this box, or drop a .env file here.
              </p>
              <Button
                type="button"
                size="sm"
                className="mt-4 gap-1.5 rounded-full bg-foreground px-4 text-background text-xs hover:bg-foreground/90"
                onClick={() => setIsEditing(true)}
              >
                <Plus className="size-3.5" />
                Add Variable
              </Button>
            </div>
          ) : (
            <div className="space-y-2">
              <textarea
                id={inputId}
                value={envVars}
                onChange={(e) => onEnvVarsChange(e.target.value)}
                placeholder="KEY=value&#10;DATABASE_URL=postgres://..."
                className="h-32 w-full rounded-xl border border-border/60 bg-background/50 p-3 font-mono text-foreground text-xs focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring"
              />
              <div className="flex items-center justify-between text-[11px] text-muted-foreground">
                <span>Supports standard KEY=value syntax, one per line.</span>
                {hasVars && (
                  <button
                    type="button"
                    onClick={() => {
                      onEnvVarsChange('');
                      setIsEditing(false);
                    }}
                    className="text-destructive hover:underline"
                  >
                    Clear all
                  </button>
                )}
              </div>
            </div>
          )}
        </div>
      )}
    </section>
  );
}
