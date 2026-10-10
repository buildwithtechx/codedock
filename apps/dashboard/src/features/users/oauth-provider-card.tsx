import { Check, CheckCircle2, Plus } from 'lucide-react';
import { useState } from 'react';
import { Button } from '#/components/ui/button';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import { Switch } from '#/components/ui/switch';
import type { SaveOAuthProviderRequest } from '#/interfaces/oauth';

export type OAuthProviderDef = {
  id: string;
  name: string;
  icon?: string;
  fields: {
    key: keyof SaveOAuthProviderRequest;
    label: string;
    placeholder: string;
    type?: string;
  }[];
};

export function OAuthProviderCard({
  provider,
  state,
  isPending,
  isSaving,
  onFieldChange,
  onSave,
}: {
  provider: OAuthProviderDef;
  state: Partial<SaveOAuthProviderRequest>;
  isPending: boolean;
  isSaving: boolean;
  onFieldChange: (key: keyof SaveOAuthProviderRequest, value: unknown) => void;
  onSave: () => void;
}) {
  const isConfigured = Boolean(state.clientId && state.clientId.trim().length > 0);
  const isActive = Boolean(state.enabled);
  const [isExpanded, setIsExpanded] = useState(isActive);

  return (
    <div className="rounded-xl border border-border/50 p-4">
      <div className="flex items-start gap-3.5">
        <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl border border-border/50 bg-background/50">
          {provider.icon ? (
            <img src={provider.icon} alt={provider.name} className="h-5 w-5 object-contain" />
          ) : (
            <span className="font-bold text-xs uppercase">{provider.name.slice(0, 2)}</span>
          )}
        </div>
        <div className="min-w-0 flex-1">
          <p className="font-medium text-foreground text-sm">{provider.name}</p>
          <p className="mt-0.5 text-muted-foreground text-xs">
            Single sign-on authentication with {provider.name}
          </p>
          <div className="mt-2 flex items-center">
            {isActive ? (
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
          <div className="space-y-3">
            {provider.fields.map((f) => (
              <div key={f.key} className="space-y-1.5">
                <Label className="text-xs">{f.label}</Label>
                <Input
                  type={f.type ?? 'text'}
                  value={(state[f.key] as string) ?? ''}
                  onChange={(e) => onFieldChange(f.key, e.target.value)}
                  placeholder={f.placeholder}
                  className="bg-muted/30 font-mono text-xs"
                />
              </div>
            ))}
          </div>

          <div className="flex items-center justify-between border-border/40 border-t pt-3">
            <div className="flex items-center gap-2">
              <Switch
                checked={state.enabled ?? false}
                onCheckedChange={(v) => onFieldChange('enabled', v)}
                aria-label={`Enable ${provider.name} login`}
              />
              <Label className="text-muted-foreground text-xs">Enable {provider.name} login</Label>
            </div>
            <Button size="sm" onClick={onSave} disabled={isPending} className="gap-1.5 text-xs">
              <Check className="size-3.5" />
              {isSaving ? 'Saving...' : `Save ${provider.name}`}
            </Button>
          </div>
        </div>
      )}
    </div>
  );
}
