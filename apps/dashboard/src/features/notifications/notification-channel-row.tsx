import { Check, CheckCircle2, Plus, Zap } from 'lucide-react';
import type React from 'react';
import { useState } from 'react';
import { Button } from '#/components/ui/button';
import { Label } from '#/components/ui/label';
import { Switch } from '#/components/ui/switch';

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
      <div className="flex items-start gap-3.5">
        <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl border border-border/50 bg-background/50 text-foreground">
          {icon}
        </div>
        <div className="min-w-0 flex-1">
          <p className="font-medium text-foreground text-sm">{name}</p>
          <p className="mt-0.5 text-muted-foreground text-xs">{description}</p>
          <div className="mt-2 flex items-center">
            {enabled ? (
              <span className="inline-flex items-center gap-1 rounded-full bg-emerald-500/10 px-2 py-0.5 font-medium text-[10px] text-emerald-600 dark:text-emerald-400">
                <CheckCircle2 className="size-3" /> Active
              </span>
            ) : isConfigured ? (
              <span className="inline-flex items-center gap-1 rounded-full bg-amber-500/10 px-2 py-0.5 font-medium text-[10px] text-amber-600 dark:text-amber-400">
                Disabled
              </span>
            ) : (
              <span className="inline-flex items-center gap-1 rounded-full bg-muted px-2 py-0.5 font-medium text-[10px] text-muted-foreground">
                Not configured
              </span>
            )}
          </div>
        </div>
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

      {isExpanded && (
        <div className="mt-4 space-y-4 border-border/40 border-t pt-4">
          <div className="space-y-3">{children}</div>

          <div className="flex items-center justify-between border-border/40 border-t pt-3">
            <div className="flex items-center gap-3">
              <div className="flex items-center gap-2">
                <Switch
                  checked={enabled}
                  onCheckedChange={onToggle}
                  aria-label={`Enable ${name}`}
                />
                <Label className="text-muted-foreground text-xs">Enable {name}</Label>
              </div>
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
            </div>
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
