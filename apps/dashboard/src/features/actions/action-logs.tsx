import { Check, Copy, Terminal } from 'lucide-react';
import { useEffect, useRef, useState } from 'react';
import { Button } from '@codedock/ui/button';

interface ActionLogsProps {
  logs?: string;
  isRunning?: boolean;
}

export function ActionLogs({ logs = '', isRunning = false }: ActionLogsProps) {
  const [copied, setCopied] = useState(false);
  const containerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (isRunning && containerRef.current) {
      containerRef.current.scrollTop = containerRef.current.scrollHeight;
    }
  }, [logs, isRunning]);

  const handleCopy = () => {
    if (!logs) return;
    navigator.clipboard.writeText(logs);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  const lines = logs.trim() ? logs.split('\n') : ['Waiting for runner output...'];

  return (
    <div className="flex flex-col rounded-xl border border-border/80 bg-zinc-950 font-mono text-xs text-zinc-100 shadow-inner">
      <div className="flex items-center justify-between border-zinc-800/80 border-b px-4 py-2 text-zinc-400">
        <div className="flex items-center gap-2">
          <Terminal className="size-4 text-emerald-400" />
          <span className="font-semibold text-zinc-300">Run Output</span>
          {isRunning && (
            <span className="flex items-center gap-1.5 text-[11px] text-blue-400">
              <span className="size-1.5 animate-pulse rounded-full bg-blue-400" />
              Live Streaming
            </span>
          )}
        </div>
        <Button
          type="button"
          variant="ghost"
          size="sm"
          onClick={handleCopy}
          className="h-7 gap-1 px-2 text-xs text-zinc-400 hover:bg-zinc-900 hover:text-zinc-200"
        >
          {copied ? <Check className="size-3.5 text-emerald-400" /> : <Copy className="size-3.5" />}
          <span>{copied ? 'Copied' : 'Copy'}</span>
        </Button>
      </div>
      <div
        ref={containerRef}
        data-testid="actions-terminal-container"
        className="max-h-[420px] min-h-[180px] overflow-y-auto p-4 leading-relaxed"
      >
        {lines.map((line, idx) => (
          <div key={idx} className="flex gap-3 hover:bg-zinc-900/40">
            <span className="w-8 shrink-0 select-none text-right text-zinc-600">{idx + 1}</span>
            <span className="whitespace-pre-wrap break-all">{line}</span>
          </div>
        ))}
      </div>
    </div>
  );
}
