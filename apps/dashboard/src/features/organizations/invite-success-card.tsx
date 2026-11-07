import { Check, Copy, Link, Send } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Button } from '@codedock/ui/button';
import { Input } from '@codedock/ui/input';

interface InviteSuccessCardProps {
  created: { url: string; email: string; linkOnly: boolean };
  onClose: () => void;
}

export function InviteSuccessCard({ created, onClose }: InviteSuccessCardProps) {
  const [copied, setCopied] = useState(false);

  const copyLink = async () => {
    try {
      await navigator.clipboard.writeText(created.url);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
      toast.success('Invitation link copied');
    } catch {
      toast.error('Failed to copy link');
    }
  };

  return (
    <div className="space-y-4 rounded-xl border border-border/50 bg-card p-5">
      <div className="flex items-start gap-3.5">
        <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-emerald-500/10 text-emerald-500">
          {created.linkOnly ? <Link className="size-5" /> : <Send className="size-5" />}
        </div>
        <div className="min-w-0">
          <h3 className="font-semibold text-foreground text-sm">
            {created.linkOnly ? 'Invitation link ready' : 'Invitation sent'}
          </h3>
          <p className="mt-0.5 break-words text-muted-foreground text-xs">
            Recipient: <span className="font-medium text-foreground">{created.email}</span>
          </p>
        </div>
      </div>

      <div className="flex flex-col gap-2 sm:flex-row">
        <Input
          value={created.url}
          readOnly
          onFocus={(e) => e.currentTarget.select()}
          className="min-w-0 flex-1 font-mono text-xs"
        />
        <Button variant="outline" size="sm" onClick={copyLink} className="shrink-0 gap-1.5">
          {copied ? <Check className="size-3.5 text-emerald-500" /> : <Copy className="size-3.5" />}
          {copied ? 'Copied' : 'Copy link'}
        </Button>
      </div>

      <div className="flex justify-end border-border/40 border-t pt-3">
        <Button size="sm" onClick={onClose}>
          Done
        </Button>
      </div>
    </div>
  );
}
