import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { Clock, Shield } from 'lucide-react';
import { toast } from 'sonner';
import { Label } from '#/components/ui/label';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '#/components/ui/select';
import { Switch } from '#/components/ui/switch';
import { auditApi } from './api';

const RETENTION_OPTIONS = [
  { value: '7', label: '7 days' },
  { value: '30', label: '30 days' },
  { value: '90', label: '90 days (default)' },
  { value: '365', label: '365 days' },
];

export function AuditSettingsCard() {
  const queryClient = useQueryClient();
  const { data: settingsRes, isLoading } = useQuery({
    queryKey: ['auditSettings'],
    queryFn: () => auditApi.getSettings(),
  });

  const mutation = useMutation({
    mutationFn: (payload: { enabled?: boolean; retentionDays?: number }) =>
      auditApi.updateSettings(payload),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ['auditSettings'] });
      toast.success('Audit settings updated');
    },
    onError: (err: unknown) => {
      toast.error(err instanceof Error ? err.message : 'Failed to update audit settings');
    },
  });

  const settings = settingsRes?.data;
  const enabled = settings?.enabled ?? true;
  const retentionDays = String(settings?.retentionDays ?? 90);

  return (
    <div className="rounded-2xl border border-border/50 bg-card p-5">
      <div className="flex items-center gap-3">
        <div className="flex size-9 items-center justify-center rounded-xl bg-muted">
          <Shield className="size-[18px] text-muted-foreground" />
        </div>
        <div>
          <h2 className="font-semibold text-[15px] text-foreground">Log retention</h2>
          <p className="text-muted-foreground text-xs">Retention period & recording</p>
        </div>
      </div>

      <div className="mt-4 space-y-4 border-border/40 border-t pt-4">
        <div className="flex items-center justify-between">
          <div className="space-y-0.5">
            <Label htmlFor="audit-logging-switch" className="text-sm">
              Record audit events
            </Label>
            <p className="text-muted-foreground text-xs">Log administrative actions</p>
          </div>
          <Switch
            id="audit-logging-switch"
            checked={enabled}
            disabled={isLoading || mutation.isPending}
            onCheckedChange={(checked) => mutation.mutate({ enabled: checked })}
          />
        </div>

        <div className="space-y-1.5">
          <div className="flex items-center gap-1.5">
            <Clock className="size-3.5 text-muted-foreground" />
            <Label htmlFor="retention-select" className="text-muted-foreground text-xs">
              Retention duration
            </Label>
          </div>
          <Select
            value={retentionDays}
            disabled={isLoading || mutation.isPending}
            onValueChange={(val) => mutation.mutate({ retentionDays: Number.parseInt(val, 10) })}
          >
            <SelectTrigger id="retention-select" className="h-8 text-xs">
              <SelectValue placeholder="Select duration" />
            </SelectTrigger>
            <SelectContent>
              {RETENTION_OPTIONS.map((opt) => (
                <SelectItem key={opt.value} value={opt.value} className="text-xs">
                  {opt.label}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
      </div>
    </div>
  );
}
