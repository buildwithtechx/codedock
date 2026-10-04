import { FileText, Import, X } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import { Card } from '#/components/ui/card';
import { type EnvEntry, parseDotenv } from '#/lib/dotenv';

interface EnvBulkModalProps {
  isOpen: boolean;
  onClose: () => void;
  onImport: (entries: EnvEntry[]) => void;
}

export function EnvBulkModal({ isOpen, onClose, onImport }: EnvBulkModalProps) {
  const [content, setContent] = useState('');

  if (!isOpen) return null;

  const handleImport = () => {
    if (!content.trim()) {
      toast.error('Please enter environment variable content');
      return;
    }

    const parsed = parseDotenv(content);
    if (parsed.length === 0) {
      toast.error('No valid KEY=VALUE assignments found in input');
      return;
    }

    onImport(parsed);
    toast.success(`Imported ${parsed.length} environment variables`);
    setContent('');
    onClose();
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 p-4 backdrop-blur-xs">
      <Card className="w-full max-w-xl p-6 shadow-2xl">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2.5">
            <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-primary/10 text-primary">
              <FileText className="h-4 w-4" />
            </div>
            <div>
              <h3 className="font-semibold text-sm">Bulk Import Environment Variables</h3>
              <p className="text-muted-foreground text-xs">
                Paste raw standard .env file contents below
              </p>
            </div>
          </div>
          <Button variant="ghost" size="icon" className="h-8 w-8" onClick={onClose}>
            <X className="h-4 w-4" />
          </Button>
        </div>

        <div className="mt-4 space-y-2">
          <textarea
            value={content}
            onChange={(e) => setContent(e.target.value)}
            placeholder={'DATABASE_URL="postgres://..."\nCACHE_PORT=6379\nAPI_SECRET=sk_live_...'}
            rows={10}
            spellCheck={false}
            className="w-full rounded-xl border border-border bg-background p-3 font-mono text-xs outline-none focus:ring-2 focus:ring-primary/20"
          />
          <p className="text-[11px] text-muted-foreground">
            Supports comments (#), single and double quoted values, and multiline escapes.
          </p>
        </div>

        <div className="mt-5 flex items-center justify-end gap-2 border-border/60 border-t pt-4">
          <Button variant="ghost" size="sm" onClick={onClose} className="text-xs">
            Cancel
          </Button>
          <Button size="sm" onClick={handleImport} className="gap-1.5 text-xs">
            <Import className="h-3.5 w-3.5" />
            Parse & Import
          </Button>
        </div>
      </Card>
    </div>
  );
}
