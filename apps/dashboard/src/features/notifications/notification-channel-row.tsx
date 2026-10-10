import { Check, CheckCircle2, Plus, Zap } from 'lucide-react';
import type React from 'react';
import { useState } from 'react';
import { Button } from '#/components/ui/button';
import { PreferenceRow } from '#/features/settings/preference-row';

export function NotificationChannelRow({
  icon,
  name,
  description,
  isConfigured,
  enabled,
  onToggle,
  onSave,
  onTest,
  saving,
  testing,
  disabled,
  children,
}: {
  icon: React.ReactNode;
  name: string;
  description: string;
  isConfigured: boolean;
  enabled: boolean;
  onToggle: (v: boolean) => void;
  onSave: () => void;
  onTest: () => void;
  saving: boolean;
  testing: boolean;
  disabled: boolean;
  children: React.ReactNode;
}) {
  const [isExpanded, setIsExpanded] = useState(enabled);

  return (
    <div className="rounded-xl border border-border/50 p-4">
      <div className="flex flex-wrap items-center justify-between gap-3.5">
        <div className="flex min-w-0 items-center gap-3.5">
          <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl border border-border/50 bg-background/50 text-foreground">
            {icon}
          </div>
          <div className="min-w-0">
            <p className="font-medium text-foreground text-sm">{name}</p>
            <p className="mt-0.5 text-muted-foreground text-xs">{description}</p>
          </div>
        </div>

        <div className="flex shrink-0 items-center gap-2.5">
          {enabled ? (
            <span className="inline-flex items-center gap-1 rounded-full bg-emerald-500/10 px-2.5 py-0.5 font-medium text-[11px] text-emerald-600 dark:text-emerald-400">
              <CheckCircle2 className="size-3" /> Active
            </span>
          ) : isConfigured ? (
            <span className="inline-flex items-center gap-1 rounded-full bg-amber-500/10 px-2.5 py-0.5 font-medium text-[11px] text-amber-600 dark:text-amber-400">
              Disabled
            </span>
          ) : (
            <span className="inline-flex items-center gap-1 rounded-full bg-muted px-2.5 py-0.5 font-medium text-[11px] text-muted-foreground">
              Not configured
            </span>
          )}
          <Button
            variant="outline"
            size="sm"
            onClick={() => setIsExpanded(!isExpanded)}
            className="shrink-0 gap-1 text-xs"
          >
            <Plus className="size-3.5" />
            {isExpanded ? 'Close' : isConfigured ? 'Edit' : 'Add'}
          </Button>
        </div>
      </div>

      {isExpanded && (
        <div className="mt-4 space-y-4 border-border/40 border-t pt-4">
          <div className="space-y-3">{children}</div>

          <PreferenceRow
            title={`Enable ${name}`}
            hint={`Deliver notification events and deployment alerts to ${name}.`}
            checked={enabled}
            onChange={onToggle}
          />

          <div className="flex items-center justify-between border-border/40 border-t pt-3">
            <Button
              size="sm"
              variant="outline"
              className="gap-1.5 text-xs"
              disabled={!enabled || testing || disabled || saving}
              onClick={onTest}
            >
              <Zap className="h-3.5 w-3.5" />
              {testing ? 'Sending…' : 'Test'}
            </Button>
            <Button
              size="sm"
              onClick={onSave}
              disabled={disabled || saving}
              className="gap-1.5 text-xs"
            >
              <Check className="size-3.5" />
              {saving ? 'Saving...' : `Save ${name}`}
            </Button>
          </div>
        </div>
      )}
    </div>
  );
}
