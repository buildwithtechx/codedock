import { Check, Copy } from 'lucide-react';
import { useEffect, useRef, useState } from 'react';
import { installCommand } from '../lib/product-links';

export function InstallCommand() {
  const [status, setStatus] = useState<'idle' | 'copied' | 'failed'>('idle');
  const timer = useRef<ReturnType<typeof setTimeout> | undefined>(undefined);
  useEffect(() => () => clearTimeout(timer.current), []);
  const copy = async () => {
    clearTimeout(timer.current);
    try {
      await navigator.clipboard.writeText(installCommand);
      setStatus('copied');
    } catch {
      setStatus('failed');
    }
    timer.current = setTimeout(() => setStatus('idle'), 3000);
  };
  return (
    <div className="w-full max-w-xl">
      <div className="flex min-w-0 items-center gap-3 rounded-xl border border-primary/25 bg-card/80 p-4 shadow-lg backdrop-blur-xl">
        <span aria-hidden="true" className="font-mono text-primary">
          $
        </span>
        <code className="min-w-0 flex-1 overflow-x-auto whitespace-nowrap py-1 text-left text-[11px] sm:text-sm">
          {installCommand}
        </code>
        <button
          type="button"
          onClick={copy}
          aria-label="Copy install command"
          className="shrink-0 rounded-lg p-2 text-muted-foreground hover:bg-muted hover:text-foreground"
        >
          {status === 'copied' ? (
            <Check className="size-4 text-emerald-400" />
          ) : (
            <Copy className="size-4" />
          )}
        </button>
      </div>
      <p aria-live="polite" className="mt-2 min-h-5 text-muted-foreground text-xs">
        {status === 'copied'
          ? 'Install command copied.'
          : status === 'failed'
            ? 'Clipboard unavailable. Select and copy the command above.'
            : 'Run on your Linux server. No manual environment setup.'}
      </p>
    </div>
  );
}
