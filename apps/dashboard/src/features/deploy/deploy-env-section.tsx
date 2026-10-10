import { ChevronDown, ChevronUp, Download, Key, Upload } from 'lucide-react';
import { type ChangeEvent, useId, useRef, useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';

interface DeployEnvSectionProps {
  envVars: string;
  onEnvVarsChange: (v: string) => void;
}

export function DeployEnvSection({ envVars, onEnvVarsChange }: DeployEnvSectionProps) {
  const [open, setOpen] = useState(false);
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

  return (
    <section className="rounded-2xl border border-border/60 bg-card p-5">
      <input
        ref={fileInputRef}
        type="file"
        accept=".env,text/plain"
        onChange={handleFileUpload}
        className="hidden"
        aria-label="Upload .env file"
      />

      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div className="flex items-center gap-3">
          <div className="grid size-8 place-items-center rounded-lg bg-muted text-muted-foreground">
            <Key className="size-4" />
          </div>
          <div>
            <h3 className="font-semibold text-foreground text-sm">Environment Variables</h3>
            <p className="text-[11px] text-muted-foreground">
              {lineCount > 0
                ? `${lineCount} variable${lineCount === 1 ? '' : 's'} configured`
                : 'None configured'}
            </p>
          </div>
        </div>

        <div className="flex flex-wrap items-center gap-2">
          <Button
            type="button"
            variant="outline"
            size="sm"
            className="h-7 text-xs"
            onClick={() => setOpen((prev) => !prev)}
          >
            Paste .env
          </Button>

          <Button
            type="button"
            variant="outline"
            size="sm"
            className="h-7 gap-1 text-xs"
            onClick={handleUploadClick}
          >
            <Upload className="size-3" />
            Upload .env
          </Button>

          <Button
            type="button"
            variant="outline"
            size="sm"
            disabled={!envVars.trim()}
            className="h-7 gap-1 text-xs"
            onClick={handleDownload}
          >
            <Download className="size-3" />
            Download
          </Button>

          <button
            type="button"
            onClick={() => setOpen((prev) => !prev)}
            className="p-1 text-muted-foreground hover:text-foreground"
            aria-label={open ? 'Collapse environment variables' : 'Expand environment variables'}
          >
            {open ? <ChevronUp className="size-4" /> : <ChevronDown className="size-4" />}
          </button>
        </div>
      </div>

      {open && (
        <div className="mt-4 space-y-2">
          <textarea
            id={inputId}
            value={envVars}
            onChange={(e) => onEnvVarsChange(e.target.value)}
            placeholder="KEY=value&#10;DATABASE_URL=postgres://..."
            className="h-32 w-full rounded-lg border border-border/60 bg-background p-3 font-mono text-foreground text-xs focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring"
          />
          <div className="flex items-center justify-between text-[11px] text-muted-foreground">
            <span>Supports standard KEY=value syntax, one per line.</span>
            {envVars.trim() && (
              <button
                type="button"
                onClick={() => onEnvVarsChange('')}
                className="text-destructive hover:underline"
              >
                Clear all
              </button>
            )}
          </div>
        </div>
      )}
    </section>
  );
}
