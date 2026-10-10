import {
  ArrowRight,
  Cloud,
  Globe,
  HardDrive,
  Loader2,
  Server,
  SlidersHorizontal,
} from 'lucide-react';
import { Button } from '#/components/ui/button';

interface WizardSummaryProps {
  appName: string;
  targetMode: 'self-hosted' | 'cloud';
  serverName: string;
  domain: string;
  hostPort: string;
  serviceCount: number;
  isDeploying: boolean;
  onDeploy: () => void;
  onReview: () => void;
  disabled?: boolean;
}

export function WizardSummary({
  appName,
  targetMode,
  serverName,
  domain,
  hostPort,
  serviceCount,
  isDeploying,
  onDeploy,
  onReview,
  disabled = false,
}: WizardSummaryProps) {
  const publicEndpoint = domain
    ? `https://${domain}`
    : hostPort
      ? `http://localhost:${hostPort}`
      : 'Auto-generated';

  return (
    <div className="space-y-4">
      <div className="rounded-2xl border border-border/60 bg-card p-4">
        <h3 className="font-semibold text-[11px] text-muted-foreground uppercase tracking-wider">
          Install Summary
        </h3>
        <div className="mt-3 space-y-2.5 text-xs">
          <div className="flex items-start justify-between gap-2">
            <span className="flex items-center gap-1.5 text-muted-foreground">
              {targetMode === 'cloud' ? (
                <Cloud className="size-3.5" />
              ) : (
                <Server className="size-3.5" />
              )}
              Destination
            </span>
            <span className="font-medium text-foreground">
              {targetMode === 'cloud' ? 'Codedock Cloud' : serverName || 'Local Server'}
            </span>
          </div>

          <div className="flex items-start justify-between gap-2">
            <span className="flex items-center gap-1.5 text-muted-foreground">
              <Globe className="size-3.5" />
              Route
            </span>
            <span className="max-w-42.5 truncate font-mono text-[11px] text-foreground">
              {publicEndpoint}
            </span>
          </div>

          <div className="flex items-start justify-between gap-2">
            <span className="flex items-center gap-1.5 text-muted-foreground">
              <HardDrive className="size-3.5" />
              Services
            </span>
            <span className="font-medium text-foreground">
              {serviceCount} {serviceCount === 1 ? 'container' : 'containers'}
            </span>
          </div>
        </div>
      </div>

      <div className="space-y-2">
        <Button
          className="w-full gap-2 py-5 font-semibold text-sm shadow-md"
          onClick={onDeploy}
          disabled={disabled || isDeploying}
        >
          {isDeploying ? (
            <>
              <Loader2 className="size-4 animate-spin" />
              Installing {appName}…
            </>
          ) : (
            <>
              Install {appName}
              <ArrowRight className="size-4" />
            </>
          )}
        </Button>
        <Button
          variant="outline"
          className="w-full gap-1.5"
          onClick={onReview}
          disabled={disabled || isDeploying}
        >
          <SlidersHorizontal className="size-3.5" />
          Review install plan
        </Button>
      </div>
    </div>
  );
}
