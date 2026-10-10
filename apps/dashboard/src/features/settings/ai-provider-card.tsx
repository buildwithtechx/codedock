import { Check, CheckCircle2, ChevronDown, Star } from 'lucide-react';
import { Button } from '#/components/ui/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '#/components/ui/dropdown-menu';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import { SettingsSection } from './settings-section';

export type AiProviderItem = {
  id: string;
  name: string;
  models: { id: string; name: string }[];
  icon: string;
  keyField: string;
  modelField: string;
};

export function AiProviderCard({
  provider,
  currentModel,
  currentKey,
  isDefault,
  isPending,
  onModelChange,
  onKeyChange,
  onSetDefault,
  onSave,
}: {
  provider: AiProviderItem;
  currentModel: string;
  currentKey: string;
  isDefault: boolean;
  isPending: boolean;
  onModelChange: (modelId: string) => void;
  onKeyChange: (key: string) => void;
  onSetDefault: () => void;
  onSave: () => void;
}) {
  const isConfigured = Boolean(currentKey && currentKey.trim().length > 0);
  const activeModel = provider.models.find((m) => m.id === currentModel) || provider.models[0];

  return (
    <SettingsSection
      collapsible
      defaultOpen={isDefault || isConfigured}
      icon={<img src={provider.icon} alt={provider.name} className="h-5 w-5 object-contain" />}
      title={provider.name}
      description={activeModel ? `Model: ${activeModel.name}` : 'Not configured'}
      action={
        <div className="flex items-center gap-2">
          {isDefault && (
            <span className="inline-flex items-center gap-1 rounded-full bg-amber-500/10 px-2.5 py-0.5 font-medium text-[10px] text-amber-600 dark:text-amber-400">
              <Star className="size-3 fill-amber-500 text-amber-500" /> Default
            </span>
          )}
          {isConfigured ? (
            <span className="inline-flex items-center gap-1 rounded-full bg-emerald-500/10 px-2.5 py-0.5 font-medium text-[10px] text-emerald-600 dark:text-emerald-400">
              <CheckCircle2 className="size-3" /> Configured
            </span>
          ) : (
            <span className="inline-flex items-center gap-1 rounded-full bg-muted px-2.5 py-0.5 font-medium text-[10px] text-muted-foreground">
              Not configured
            </span>
          )}
        </div>
      }
    >
      <div className="space-y-4">
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
          <div className="space-y-1.5">
            <Label className="text-xs">Model</Label>
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="outline" className="w-full justify-between font-normal text-xs">
                  <span className="truncate">
                    {activeModel ? `${activeModel.name} (${activeModel.id})` : 'Select model'}
                  </span>
                  <ChevronDown className="h-3.5 w-3.5 opacity-60" />
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="start" className="max-h-60 w-64 overflow-y-auto">
                {provider.models.map((m) => (
                  <DropdownMenuItem
                    key={m.id}
                    onSelect={() => onModelChange(m.id)}
                    className="flex cursor-pointer items-center justify-between py-1.5"
                  >
                    <div className="flex flex-col">
                      <span className="font-medium text-xs">{m.name}</span>
                      <span className="font-mono text-[10px] text-muted-foreground">{m.id}</span>
                    </div>
                    {currentModel === m.id && <Check className="h-3.5 w-3.5 text-emerald-500" />}
                  </DropdownMenuItem>
                ))}
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
          <div className="space-y-1.5">
            <Label className="text-xs">API Key</Label>
            <Input
              type="password"
              placeholder={`API key for ${provider.name}`}
              value={currentKey}
              onChange={(e) => onKeyChange(e.target.value)}
              className="bg-muted/30 font-mono text-xs"
            />
          </div>
        </div>

        <div className="flex items-center justify-between border-border/40 border-t pt-3">
          <Button
            type="button"
            variant="ghost"
            size="sm"
            onClick={onSetDefault}
            disabled={isDefault || !isConfigured}
            className="gap-1.5 text-muted-foreground text-xs hover:text-foreground"
          >
            <Star className={`size-3.5 ${isDefault ? 'fill-current' : ''}`} />
            {isDefault ? 'Default provider' : 'Make default'}
          </Button>
          <Button size="sm" onClick={onSave} disabled={isPending} className="gap-1.5 text-xs">
            <Check className="size-3.5" />
            Save {provider.name}
          </Button>
        </div>
      </div>
    </SettingsSection>
  );
}
