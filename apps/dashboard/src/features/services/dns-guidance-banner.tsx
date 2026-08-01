import { Check, Copy, Info } from 'lucide-react';
import { useState } from 'react';
import { Button } from '#/components/ui/button';

interface DnsGuidanceBannerProps {
  serverIp: string;
  hasDnsProvider: boolean;
}

export function DnsGuidanceBanner({ serverIp, hasDnsProvider }: DnsGuidanceBannerProps) {
  const [copied, setCopied] = useState(false);

  const handleCopy = async () => {
    if (!serverIp) return;

    try {
      await navigator.clipboard.writeText(serverIp);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    } catch {
      setCopied(false);
    }
  };

  return (
    <div className="rounded-xl border border-blue-500/20 bg-blue-500/5 p-4 text-sm">
      <div className="flex items-start gap-3">
        <div className="mt-0.5 flex h-7 w-7 shrink-0 items-center justify-center rounded-lg bg-blue-500/10 text-blue-600">
          <Info className="h-4 w-4" />
        </div>
        <div className="flex-1 space-y-2">
          <div className="flex items-center justify-between">
            <h4 className="font-semibold text-blue-900 dark:text-blue-200">
              Manual DNS Setup Instructions
            </h4>
            {!hasDnsProvider && (
              <span className="rounded-md border border-amber-500/30 bg-amber-500/10 px-2 py-0.5 font-medium text-[10px] text-amber-600 uppercase">
                No DNS Provider Configured
              </span>
            )}
          </div>
          <p className="text-muted-foreground text-xs leading-relaxed">
            {!hasDnsProvider
              ? 'Automated DNS provider API credentials are not configured. To direct traffic to this service, add an A record in your DNS provider pointing to your server IP.'
              : 'For domains requiring manual setup, create an A record in your DNS registrar dashboard using the details below.'}
          </p>

          <div className="mt-3 grid grid-cols-1 gap-3 rounded-lg border bg-background/60 p-3 sm:grid-cols-3">
            <div>
              <span className="text-[10px] text-muted-foreground uppercase tracking-wider">
                Type
              </span>
              <p className="font-mono font-semibold text-xs">A</p>
            </div>
            <div>
              <span className="text-[10px] text-muted-foreground uppercase tracking-wider">
                Host / Name
              </span>
              <p className="font-mono font-semibold text-xs">@ (or subdomain)</p>
            </div>
            <div>
              <span className="text-[10px] text-muted-foreground uppercase tracking-wider">
                Points To (Server IP)
              </span>
              <div className="mt-0.5 flex items-center gap-1.5">
                <code className="font-bold font-mono text-xs">{serverIp || 'Not configured'}</code>
                {serverIp && (
                  <Button
                    type="button"
                    variant="ghost"
                    size="icon"
                    className="h-5 w-5"
                    onClick={handleCopy}
                    aria-label="Copy server IP to clipboard"
                  >
                    {copied ? (
                      <Check className="h-3 w-3 text-emerald-500" />
                    ) : (
                      <Copy className="h-3 w-3 text-muted-foreground" />
                    )}
                  </Button>
                )}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
