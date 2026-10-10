import { Check } from 'lucide-react';
import { Button } from '#/components/ui/button';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import { Switch } from '#/components/ui/switch';
import { SettingsSection } from '#/features/settings/settings-section';
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
  return (
    <SettingsSection
      collapsible
      defaultOpen={Boolean(state.enabled)}
      icon={
        provider.icon ? (
          <img src={provider.icon} alt={provider.name} className="h-5 w-5 object-contain" />
        ) : (
          <span className="font-bold text-xs uppercase">{provider.name.slice(0, 2)}</span>
        )
      }
      title={provider.name}
      description={`Configure single sign-on authentication using ${provider.name}.`}
      action={
        <div className="flex items-center gap-2.5">
          <span className="text-muted-foreground text-xs">
            {state.enabled ? 'Enabled' : 'Disabled'}
          </span>
          <Switch
            checked={state.enabled ?? false}
            onCheckedChange={(v) => onFieldChange('enabled', v)}
          />
        </div>
      }
    >
      <div className="space-y-4">
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
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
        <div className="flex justify-end border-border/40 border-t pt-3">
          <Button size="sm" onClick={onSave} disabled={isPending}>
            <Check className="mr-2 h-4 w-4" />
            {isSaving ? 'Saving...' : `Save ${provider.name}`}
          </Button>
        </div>
      </div>
    </SettingsSection>
  );
}
