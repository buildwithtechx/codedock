import { AlertTriangle, CheckCircle2, CircleDashed, Cloud, Loader2, Server } from 'lucide-react';
import { Button } from '#/components/ui/button';
import { cn } from '#/lib/utils';
import type { DestinationKind, S3Destination, SFTPDestination } from './interfaces';

export interface UnifiedDestination {
  kind: DestinationKind;
  id: string;
  name: string;
  description: string;
  pathPrefix?: string;
  isDefault: boolean;
  lastVerifiedAt?: string;
  lastVerifyError?: string;
  createdAt: string;
  provider: string;
  detail: string;
  credential: string;
}

export function toUnifiedS3(destination: S3Destination): UnifiedDestination {
  const locator = [destination.bucket, destination.region, destination.endpoint]
    .filter(Boolean)
    .join(' · ');
  return {
    kind: 's3',
    id: destination.id,
    name: destination.name,
    description: destination.description,
    pathPrefix: destination.pathPrefix,
    isDefault: destination.isDefault ?? false,
    lastVerifiedAt: destination.lastVerifiedAt,
    lastVerifyError: destination.lastVerifyError,
    createdAt: destination.createdAt,
    provider: destination.provider || 's3',
    detail: locator || destination.bucket,
    credential: destination.accessKeyId ? 'Access key stored' : 'No credentials',
  };
}

export function toUnifiedSftp(destination: SFTPDestination): UnifiedDestination {
  const port = destination.port > 0 ? destination.port : 22;
  const prefix = destination.pathPrefix ? `:${destination.pathPrefix}` : '';
  return {
    kind: 'sftp',
    id: destination.id,
    name: destination.name,
    description: destination.description,
    pathPrefix: destination.pathPrefix,
    isDefault: false,
    lastVerifiedAt: destination.lastVerifiedAt,
    lastVerifyError: destination.lastVerifyError,
    createdAt: destination.createdAt,
    provider: 'sftp',
    detail: `${destination.username || '?'}@${destination.host || '?'}:${port}${prefix}`,
    credential: 'SSH credentials stored',
  };
}

export function kindLabel(kind: DestinationKind): string {
  return kind === 's3' ? 'S3 compatible' : 'SFTP server';
}

export function DestinationKindIcon({
  kind,
  className,
}: {
  kind: DestinationKind;
  className?: string;
}) {
  if (kind === 'sftp') return <Server className={className} />;
  return <Cloud className={className} />;
}

export function providerLabel(provider: string): string {
  const normalized = provider.toLowerCase();
  if (normalized === 'r2' || normalized.includes('cloudflare')) return 'Cloudflare R2';
  if (normalized === 's3' || normalized.includes('aws') || normalized.includes('amazon'))
    return 'AWS S3';
  if (normalized === 'b2' || normalized.includes('backblaze')) return 'Backblaze B2';
  if (normalized.includes('wasabi')) return 'Wasabi';
  if (normalized === 'do' || normalized.includes('digitalocean') || normalized.includes('spaces'))
    return 'DigitalOcean Spaces';
  if (normalized.includes('minio')) return 'MinIO';
  if (normalized.includes('gcp') || normalized.includes('google')) return 'Google Cloud';
  return provider || 'S3 compatible';
}

export function formatBytes(bytes: number): string {
  if (!bytes || bytes <= 0) return '0 B';
  const units = ['B', 'KB', 'MB', 'GB', 'TB'];
  const index = Math.min(Math.floor(Math.log(bytes) / Math.log(1024)), units.length - 1);
  return `${(bytes / 1024 ** index).toFixed(1)} ${units[index]}`;
}

export function formatDuration(start?: string, end?: string): string {
  if (!start || !end) return '—';
  const diffMs = new Date(end).getTime() - new Date(start).getTime();
  if (diffMs <= 0) return '< 1s';
  const seconds = Math.floor(diffMs / 1000);
  if (seconds < 60) return `${seconds}s`;
  return `${Math.floor(seconds / 60)}m ${seconds % 60}s`;
}

export function DestinationVerificationBadge({
  destination,
  verifying = false,
}: {
  destination: Pick<UnifiedDestination, 'lastVerifiedAt' | 'lastVerifyError'>;
  verifying?: boolean;
}) {
  const failed = Boolean(destination.lastVerifyError) && !verifying;
  const verified = !failed && !verifying && Boolean(destination.lastVerifiedAt);
  const title = verifying
    ? 'Verifying connection'
    : failed
      ? destination.lastVerifyError
      : verified && destination.lastVerifiedAt
        ? `Last verified ${new Date(destination.lastVerifiedAt).toLocaleString()}`
        : 'Not verified yet';
  return (
    <span
      title={title}
      className={cn(
        'inline-flex items-center gap-1 rounded-full px-2 py-0.5 font-medium text-xs',
        verifying && 'bg-muted text-muted-foreground',
        failed && 'bg-destructive/10 text-destructive',
        verified && 'bg-emerald-500/10 text-emerald-500',
        !verifying && !failed && !verified && 'bg-muted text-muted-foreground'
      )}
    >
      {verifying ? (
        <Loader2 className="h-3 w-3 animate-spin" />
      ) : failed ? (
        <AlertTriangle className="h-3 w-3" />
      ) : verified ? (
        <CheckCircle2 className="h-3 w-3" />
      ) : (
        <CircleDashed className="h-3 w-3" />
      )}
      {verifying ? 'Verifying' : failed ? 'Failed' : verified ? 'Verified' : 'Not verified'}
    </span>
  );
}

export function DestinationLoadError({
  message,
  onRetry,
  busy,
}: {
  message: string;
  onRetry: () => void;
  busy: boolean;
}) {
  return (
    <div
      role="alert"
      className="mb-5 flex flex-wrap items-center justify-between gap-3 rounded-xl bg-destructive/10 px-4 py-3 text-destructive text-sm"
    >
      <span>{message}</span>
      <Button
        variant="ghost"
        size="sm"
        onClick={onRetry}
        disabled={busy}
        className="h-8 text-destructive hover:bg-destructive/10 hover:text-destructive"
      >
        {busy ? <Loader2 className="h-3.5 w-3.5 animate-spin" /> : 'Retry'}
      </Button>
    </div>
  );
}
