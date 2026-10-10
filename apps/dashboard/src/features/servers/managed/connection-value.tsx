import { Check, Copy, Eye, EyeOff } from 'lucide-react';
import { useState } from 'react';
import { Button } from '#/components/ui/button';

export function ConnectionValue({
  label,
  value,
  secret = false,
}: {
  label: string;
  value: string;
  secret?: boolean;
}) {
  const [copied, setCopied] = useState(false);
  const [revealed, setRevealed] = useState(false);

  const handleCopy = async () => {
    try {
      await navigator.clipboard.writeText(value);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    } catch {
      setCopied(false);
    }
  };

  const displayValue = secret && !revealed ? '••••••••••••••••••••••••••••••••' : value;

  return (
    <div className="min-w-0 space-y-1.5">
      <p className="font-medium text-muted-foreground text-xs">{label}</p>
      <div className="flex min-w-0 items-center gap-2 rounded-xl bg-muted/40 p-3">
        <code dir="ltr" className="min-w-0 flex-1 whitespace-pre-wrap break-all font-mono text-xs">
          {displayValue}
        </code>
        {secret && (
          <Button
            type="button"
            size="icon"
            variant="ghost"
            className="size-7"
            aria-label={revealed ? 'Hide secret' : 'Reveal secret'}
            onClick={() => setRevealed((v) => !v)}
          >
            {revealed ? <EyeOff className="size-3.5" /> : <Eye className="size-3.5" />}
          </Button>
        )}
        <Button
          type="button"
          size="icon"
          variant="ghost"
          className="size-7"
          aria-label={`Copy ${label}`}
          onClick={() => void handleCopy()}
        >
          {copied ? <Check className="size-3.5 text-success" /> : <Copy className="size-3.5" />}
        </Button>
      </div>
    </div>
  );
}
