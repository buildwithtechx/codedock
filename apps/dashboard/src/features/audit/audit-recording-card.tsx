import { Loader2 } from 'lucide-react';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@codedock/ui/select';
import { Switch } from '@codedock/ui/switch';

const RETENTION_CHOICES = [7, 30, 90, 180, 365];

interface AuditRecordingCardProps {
  enabled: boolean;
  retentionDays: number;
  saving?: boolean;
  canManage?: boolean;
  onToggle: (enabled: boolean) => void;
  onRetentionChange: (days: number) => void;
}

export function AuditRecordingCard({
  enabled,
  retentionDays,
  saving = false,
  canManage = true,
  onToggle,
  onRetentionChange,
}: AuditRecordingCardProps) {
  return (
    <div className="rounded-2xl border border-border/50 bg-card p-5">
      <div className="flex items-start justify-between gap-4">
        <div className="min-w-0">
          <p className="font-medium text-foreground text-sm">Record activity</p>
          <p className="mt-1 text-muted-foreground text-xs leading-relaxed">
            Keep a record of every change made in this organization. On by default.
          </p>
          {!canManage && (
            <p className="mt-1 text-muted-foreground/70 text-xs">
              Admin permissions required to modify.
            </p>
          )}
        </div>
        <div className="flex shrink-0 items-center gap-2">
          {saving && <Loader2 className="size-3.5 animate-spin text-muted-foreground" />}
          <Switch
            checked={enabled}
            onCheckedChange={onToggle}
            disabled={!canManage || saving}
            aria-label="Record activity"
          />
        </div>
      </div>

      <div className="mt-4 flex flex-wrap items-center gap-3 border-border/40 border-t pt-4 text-xs">
        <span className="text-muted-foreground">Keep entries for</span>
        <div className="w-32">
          <Select
            value={String(retentionDays)}
            onValueChange={(val) => onRetentionChange(Number(val))}
            disabled={!canManage || saving}
          >
            <SelectTrigger size="sm" className="h-8">
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              {RETENTION_CHOICES.map((days) => (
                <SelectItem key={days} value={String(days)}>
                  {days} days
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
        <span className="block w-full text-[11px] text-muted-foreground/70 sm:w-auto">
          Older entries are deleted automatically.
        </span>
      </div>
    </div>
  );
}
