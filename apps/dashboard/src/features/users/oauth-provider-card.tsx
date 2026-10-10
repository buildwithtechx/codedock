import { Check, CheckCircle2, Plus } from 'lucide-react';
import { useState } from 'react';
import { Button } from '#/components/ui/button';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import { PreferenceRow } from '#/features/settings/preference-row';
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
      <div className="flex flex-wrap items-center justify-between gap-3.5">
        <div className="flex min-w-0 items-center gap-3.5">
          <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl border border-border/50 bg-background/50">
            {provider.icon ? (
              <img
                src={provider.icon}
                alt={provider.name}
                className={`h-5 w-5 object-contain ${
                  provider.id === 'github'
                    ? 'dim:brightness-0 dim:invert dark:brightness-0 dark:invert'
                    : ''
                }`}
              />
            ) : (
              <span className="font-bold text-xs uppercase">{provider.name.slice(0, 2)}</span>
            )}
          </div>
          <div className="min-w-0">
            <p className="font-medium text-foreground text-sm">{provider.name}</p>
            <p className="mt-0.5 text-muted-foreground text-xs">
              Single sign-on authentication with {provider.name}
            </p>
          </div>
        </div>

        <div className="flex shrink-0 items-center gap-2.5">
          {isActive ? (
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

          <PreferenceRow
            title={`Enable ${provider.name}`}
            hint={`Allow users to sign in with their ${provider.name} account.`}
            checked={Boolean(state.enabled)}
            onChange={(checked) => onFieldChange('enabled', checked)}
          />

          <div className="flex justify-end border-border/40 border-t pt-3">
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
