import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { Gauge, Globe, Loader2, RefreshCw, ScanSearch, ShieldCheck } from 'lucide-react';
import { useEffect, useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import type { ServerRateLimitConfig } from '#/interfaces/server';
import { serverService } from '#/services/servers';

interface ServerSecurityTabProps {
  serverId: string;
}

export function ServerSecurityTab({ serverId }: ServerSecurityTabProps) {
  const queryClient = useQueryClient();

  const { data: rateLimit, isLoading: isRateLimitLoading } = useQuery({
    queryKey: ['server-rate-limit', serverId],
    queryFn: () => serverService.getRateLimit(serverId),
  });

  const {
    mutate: runScan,
    data: scanResult,
    isPending: isScanning,
  } = useMutation({
    mutationFn: () => serverService.scanPorts(serverId),
    onSuccess: () => {
      toast.success('Port scan completed successfully');
    },
    onError: () => {
      toast.error('Failed to run port scan');
    },
  });

  const { mutate: saveRateLimit, isPending: isSavingRateLimit } = useMutation({
    mutationFn: (payload: ServerRateLimitConfig) =>
      serverService.updateRateLimit(serverId, payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['server-rate-limit', serverId] });
      toast.success('Rate limit settings saved');
    },
    onError: () => {
      toast.error('Failed to save rate limit settings');
    },
  });

  const [rps, setRps] = useState<number>(50);
  const [burst, setBurst] = useState<number>(20);
  const [whitelistText, setWhitelistText] = useState<string>('');

  useEffect(() => {
    if (rateLimit) {
      setRps(rateLimit.rps);
      setBurst(rateLimit.burst);
      setWhitelistText((rateLimit.whitelist || []).join(', '));
    }
  }, [rateLimit]);

  const handleSaveRateLimit = (e: React.FormEvent) => {
    e.preventDefault();
    const whitelist = whitelistText
      .split(',')
      .map((s) => s.trim())
      .filter(Boolean);
    saveRateLimit({
      rps: Number(rps) || 50,
      burst: Number(burst) || 20,
      whitelist,
    });
  };

  return (
    <div className="space-y-6">
      <div className="rounded-2xl bg-card">
        <div className="flex flex-wrap items-center justify-between gap-3 border-border/50 border-b px-5 py-4">
          <div className="flex min-w-0 items-center gap-3">
            <div className="flex h-9 w-9 items-center justify-center rounded-xl bg-primary/10">
              <ScanSearch className="size-[18px] text-primary" />
            </div>
            <div>
              <h2 className="font-semibold text-[15px] text-foreground">Exposed ports</h2>
              <p className="text-muted-foreground text-xs">
                Listeners reachable from outside this node
              </p>
            </div>
          </div>
          <Button variant="outline" size="sm" onClick={() => runScan()} disabled={isScanning}>
            {isScanning ? (
              <Loader2 className="mr-1.5 size-3.5 animate-spin" />
            ) : (
              <RefreshCw className="mr-1.5 size-3.5" />
            )}
            {scanResult ? 'Scan again' : 'Run port scan'}
          </Button>
        </div>

        <div className="p-5">
          {!scanResult && !isScanning && (
            <div className="flex flex-col items-center justify-center rounded-xl border border-border/70 border-dashed p-8 text-center">
              <div className="flex size-10 items-center justify-center rounded-xl bg-muted">
                <ScanSearch className="size-5 text-muted-foreground" />
              </div>
              <h3 className="mt-3 font-medium text-foreground text-sm">No recent port scan</h3>
              <p className="mt-1 max-w-sm text-muted-foreground text-xs">
                Run an exposure audit to discover listening sockets, public services, and internal
                bindings.
              </p>
              <Button variant="secondary" size="sm" className="mt-4" onClick={() => runScan()}>
                Scan open ports
              </Button>
            </div>
          )}

          {isScanning && !scanResult && (
            <div className="flex h-36 items-center justify-center">
              <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
            </div>
          )}

          {scanResult && (
            <div className="divide-y divide-border/40">
              {scanResult.listeners.map((l) => (
                <div
                  key={`${l.protocol}-${l.port}`}
                  className="flex flex-col gap-2 py-3 sm:flex-row sm:items-center sm:justify-between"
                >
                  <div className="flex items-center gap-3">
                    <span className="font-mono font-semibold text-foreground text-sm">
                      {l.port}/{l.protocol}
                    </span>
                    <span className="rounded-md bg-muted px-2 py-0.5 font-mono text-[11px] text-muted-foreground">
                      {l.process}
                    </span>
                    <span className="font-mono text-muted-foreground text-xs">{l.state}</span>
                  </div>
                  <div className="flex items-center gap-2">
                    {l.exposed ? (
                      <span className="inline-flex items-center gap-1 rounded-full bg-amber-500/10 px-2 py-0.5 font-medium text-[11px] text-amber-500">
                        <Globe className="size-3" />
                        Publicly exposed
                      </span>
                    ) : (
                      <span className="inline-flex items-center gap-1 rounded-full bg-emerald-500/10 px-2 py-0.5 font-medium text-[11px] text-emerald-500">
                        <ShieldCheck className="size-3" />
                        Internal only
                      </span>
                    )}
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>

      <div className="rounded-2xl bg-card">
        <div className="flex flex-wrap items-center justify-between gap-3 border-border/50 border-b px-5 py-4">
          <div className="flex min-w-0 items-center gap-3">
            <div className="flex h-9 w-9 items-center justify-center rounded-xl bg-warning/10">
              <Gauge className="size-[18px] text-warning" />
            </div>
            <div>
              <h2 className="font-semibold text-[15px] text-foreground">Rate limiting</h2>
              <p className="text-muted-foreground text-xs">Edge request policy for this node</p>
            </div>
          </div>
        </div>

        <form onSubmit={handleSaveRateLimit} className="space-y-4 p-5">
          {isRateLimitLoading ? (
            <div className="flex h-24 items-center justify-center">
              <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
            </div>
          ) : (
            <>
              <div className="grid gap-4 sm:grid-cols-2">
                <div className="space-y-1.5">
                  <Label htmlFor="rps">Requests per second (RPS)</Label>
                  <Input
                    id="rps"
                    type="number"
                    min="1"
                    max="10000"
                    value={rps}
                    onChange={(e) => setRps(Number(e.target.value))}
                    disabled={isSavingRateLimit}
                  />
                  <p className="text-[11px] text-muted-foreground">
                    Maximum average request rate per client IP
                  </p>
                </div>

                <div className="space-y-1.5">
                  <Label htmlFor="burst">Burst limit</Label>
                  <Input
                    id="burst"
                    type="number"
                    min="1"
                    max="5000"
                    value={burst}
                    onChange={(e) => setBurst(Number(e.target.value))}
                    disabled={isSavingRateLimit}
                  />
                  <p className="text-[11px] text-muted-foreground">
                    Maximum temporary spike allowed beyond RPS limit
                  </p>
                </div>
              </div>

              <div className="space-y-1.5">
                <Label htmlFor="whitelist">IP Allowlist</Label>
                <Input
                  id="whitelist"
                  type="text"
                  placeholder="127.0.0.1, 10.0.0.0/8"
                  value={whitelistText}
                  onChange={(e) => setWhitelistText(e.target.value)}
                  disabled={isSavingRateLimit}
                />
                <p className="text-[11px] text-muted-foreground">
                  Comma-separated list of IP addresses or CIDRs exempt from rate limits
                </p>
              </div>

              <div className="flex justify-end pt-2">
                <Button type="submit" disabled={isSavingRateLimit}>
                  {isSavingRateLimit && <Loader2 className="mr-1.5 size-3.5 animate-spin" />}
                  Save configuration
                </Button>
              </div>
            </>
          )}
        </form>
      </div>
    </div>
  );
}
