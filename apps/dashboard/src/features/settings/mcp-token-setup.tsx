import { Link } from '@tanstack/react-router';
import { Check, Copy, ExternalLink, Key } from 'lucide-react';
import { useState } from 'react';
import { Button } from '#/components/ui/button';

export function McpTokenSetup({ endpoint }: { endpoint: string }) {
  const [copied, setCopied] = useState(false);

  const snippet = JSON.stringify(
    {
      mcpServers: {
        codedock: {
          url: endpoint,
          headers: {
            Authorization: 'Bearer cd_pat_YOUR_TOKEN_HERE',
          },
        },
      },
    },
    null,
    2
  );

  const copy = async () => {
    await navigator.clipboard.writeText(snippet);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <div className="space-y-4">
      <div className="flex items-start gap-3 rounded-xl border border-border/60 bg-muted/20 p-4">
        <div className="flex size-8 shrink-0 items-center justify-center rounded-lg bg-amber-500/10 text-amber-500">
          <Key className="size-4" />
        </div>
        <div className="min-w-0 flex-1">
          <p className="font-medium text-foreground text-xs">Personal Access Token Required</p>
          <p className="mt-0.5 text-muted-foreground text-xs leading-relaxed">
            Generate an API token with appropriate workspace permissions to authenticate requests
            from clients that do not yet support interactive browser authorization.
          </p>
          <div className="mt-3">
            <Button asChild variant="outline" size="sm" className="gap-1.5 text-xs">
              <Link to="/settings" search={{ tab: 'tokens' }}>
                Create API Token
                <ExternalLink className="size-3" />
              </Link>
            </Button>
          </div>
        </div>
      </div>

      <div className="space-y-2">
        <p className="font-medium text-foreground text-xs">Static Token Configuration</p>
        <div className="relative">
          <pre className="overflow-x-auto rounded-xl border border-border/60 bg-muted/30 p-3.5 pe-14 font-mono text-xs leading-relaxed">
            {snippet}
          </pre>
          <Button
            variant="ghost"
            size="icon"
            onClick={copy}
            className="absolute top-2 right-2 size-7 text-muted-foreground hover:text-foreground"
            aria-label="Copy snippet"
          >
            {copied ? (
              <Check className="size-3.5 text-emerald-500" />
            ) : (
              <Copy className="size-3.5" />
            )}
          </Button>
        </div>
      </div>
    </div>
  );
}
