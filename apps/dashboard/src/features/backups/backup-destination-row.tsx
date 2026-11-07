import { Link } from '@tanstack/react-router';
import { Loader2, MoreVertical, Pencil, RefreshCw, Star, Trash2 } from 'lucide-react';
import { Button } from '@codedock/ui/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@codedock/ui/dropdown-menu';
import type { S3Destination } from '#/interfaces/backups';
import type { UnifiedDestination } from './destination-display';
import {
  DestinationKindIcon,
  DestinationVerificationBadge,
  formatBytes,
  kindLabel,
} from './destination-display';

export interface BackupDestinationRowProps {
  row: UnifiedDestination;
  stats?: { runs: number; bytes: number; active: number; failed: number };
  verifying: boolean;
  s3Destination?: S3Destination;
  onVerify: (row: UnifiedDestination) => void;
  onEdit: (s3: S3Destination) => void;
  onSetDefault: (row: UnifiedDestination) => void;
  onDelete: (row: UnifiedDestination) => void;
}

export function BackupDestinationRow({
  row,
  stats,
  verifying,
  s3Destination,
  onVerify,
  onEdit,
  onSetDefault,
  onDelete,
}: BackupDestinationRowProps) {
  return (
    <li className="group relative flex items-start gap-3 px-5 py-4 transition-colors hover:bg-muted/20">
      <Link
        to="/backups/$backupId"
        params={{ backupId: row.id }}
        aria-label={row.name}
        className="absolute inset-0 z-0 rounded-xl focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
      />
      <div className="mt-0.5 flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-background text-muted-foreground">
        <DestinationKindIcon kind={row.kind} className="h-5 w-5" />
      </div>
      <div className="min-w-0 flex-1">
        <div className="flex flex-wrap items-center gap-2">
          <p className="truncate font-medium text-foreground text-sm">{row.name}</p>
          {row.isDefault && (
            <span className="inline-flex items-center gap-1 text-muted-foreground text-xs">
              <Star className="h-3 w-3" />
              Default
            </span>
          )}
          <DestinationVerificationBadge destination={row} verifying={verifying} />
        </div>
        <p className="mt-1 truncate text-muted-foreground text-xs" title={row.detail}>
          {kindLabel(row.kind)} · {row.detail}
        </p>
        <div className="mt-1.5 flex flex-wrap items-center gap-x-2 gap-y-1 text-muted-foreground text-xs">
          {stats && stats.runs > 0 ? (
            <>
              <span>
                {formatBytes(stats.bytes)} stored · {stats.runs} runs
              </span>
              {stats.active > 0 && <span className="text-sky-500">· {stats.active} active</span>}
              {stats.failed > 0 && (
                <span className="text-destructive">· {stats.failed} failed</span>
              )}
            </>
          ) : (
            <span>No runs yet</span>
          )}
        </div>
        {row.lastVerifyError && (
          <p className="mt-1 truncate text-destructive text-xs" title={row.lastVerifyError}>
            {row.lastVerifyError}
          </p>
        )}
      </div>
      <div className="relative z-10 flex shrink-0 items-center gap-1">
        <Button
          variant="ghost"
          size="icon"
          onClick={() => onVerify(row)}
          disabled={verifying}
          title="Verify connection"
          aria-label={`Verify connection: ${row.name}`}
        >
          {verifying ? (
            <Loader2 className="h-4 w-4 animate-spin" />
          ) : (
            <RefreshCw className="h-4 w-4" />
          )}
        </Button>
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="ghost" size="icon" aria-label={`Actions: ${row.name}`}>
              <MoreVertical className="h-4 w-4" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end">
            {s3Destination && (
              <DropdownMenuItem onClick={() => onEdit(s3Destination)}>
                <Pencil className="h-4 w-4" />
                Edit
              </DropdownMenuItem>
            )}
            {row.kind === 's3' && !row.isDefault && (
              <DropdownMenuItem onClick={() => onSetDefault(row)}>
                <Star className="h-4 w-4" />
                Set as default
              </DropdownMenuItem>
            )}
            {(s3Destination || (row.kind === 's3' && !row.isDefault)) && <DropdownMenuSeparator />}
            <DropdownMenuItem variant="destructive" onClick={() => onDelete(row)}>
              <Trash2 className="h-4 w-4" />
              Delete
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
    </li>
  );
}
