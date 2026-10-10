import { Check, CheckCircle2, ChevronDown, Plus, Star } from 'lucide-react';
import { useState } from 'react';
import { Button } from '#/components/ui/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '#/components/ui/dropdown-menu';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import { PreferenceRow } from './preference-row';

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
  const [isExpanded, setIsExpanded] = useState(isDefault || isConfigured);

  return (
    <div className="rounded-xl border border-border/50 p-4">
      <div className="flex items-start gap-3.5">
        <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl border border-border/50 bg-background/50">
          <img src={provider.icon} alt={provider.name} className="h-5 w-5 object-contain" />
        </div>
        <div className="min-w-0 flex-1">
          <p className="font-medium text-foreground text-sm">{provider.name}</p>
          <p className="mt-0.5 text-muted-foreground text-xs">
            {activeModel ? `Model: ${activeModel.name}` : 'Select a model and API key'}
          </p>
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

      <div className="mt-3">
        {isConfigured ? (
          <div className="flex items-center gap-1.5">
            <span className="inline-flex items-center gap-1 rounded-full bg-emerald-500/10 px-2 py-0.5 font-medium text-[10px] text-emerald-600 dark:text-emerald-400">
              <CheckCircle2 className="size-3" /> Configured
            </span>
            {isDefault && (
              <span className="inline-flex items-center gap-1 rounded-full bg-amber-500/10 px-2 py-0.5 font-medium text-[10px] text-amber-600 dark:text-amber-400">
                <Star className="size-3 fill-amber-500 text-amber-500" /> Default
              </span>
            )}
          </div>
        ) : (
          <p className="py-0.5 text-muted-foreground text-xs">None added yet.</p>
        )}
      </div>

      {isExpanded && (
        <div className="mt-4 space-y-4 border-border/40 border-t pt-4">
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
              <DropdownMenuContent align="start" className="max-h-60 w-72 overflow-y-auto">
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

          <PreferenceRow
            title="Default Provider"
            hint={`Set ${provider.name} as the default AI provider across Codedock.`}
            checked={isDefault}
            onChange={onSetDefault}
          />

          <div className="flex justify-end border-border/40 border-t pt-3">
            <Button size="sm" onClick={onSave} disabled={isPending} className="gap-1.5 text-xs">
              <Check className="size-3.5" />
              Save {provider.name}
            </Button>
          </div>
        </div>
      )}
    </div>
  );
}
