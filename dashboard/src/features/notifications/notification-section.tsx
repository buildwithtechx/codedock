import { Zap } from 'lucide-react';
import type React from 'react';
import { Button } from '#/components/ui/button';
import { Switch } from '#/components/ui/switch';

export type NotificationSectionProps = {
  icon: React.ReactNode;
  title: string;
  provider: string;
  enabled: boolean;
  onToggle: (v: boolean) => void;
  onTest: () => void;
  testing: boolean;
  children: React.ReactNode;
};

export const NotificationSection = ({
  icon,
  title,
  provider: _provider,
  enabled,
  onToggle,
  onTest,
  testing,
  children,
}: NotificationSectionProps) => (
  <div className="rounded-xl border border-border/50 bg-card/40 p-6">
    <div className="mb-4 flex items-center justify-between">
      <div className="flex items-center gap-3">
        <div className="flex h-9 w-9 items-center justify-center rounded-lg bg-primary/10 text-primary">
          {icon}
        </div>
        <span className="font-semibold text-sm">{title}</span>
      </div>
      <div className="flex items-center gap-3">
        <Button
          size="sm"
          variant="outline"
          className="text-xs"
          disabled={!enabled || testing}
          onClick={onTest}
        >
          <Zap className="mr-1.5 h-3.5 w-3.5" />
          {testing ? 'Sending…' : 'Test'}
        </Button>
        <Switch checked={enabled} onCheckedChange={onToggle} />
      </div>
    </div>
    <div className={`space-y-4 ${!enabled ? 'pointer-events-none opacity-40' : ''}`}>
      {children}
    </div>
  </div>
);
