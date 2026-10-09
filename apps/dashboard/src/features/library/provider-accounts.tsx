import { Link2 } from 'lucide-react';
import { Button } from '#/components/ui/button';
import { cn } from '#/lib/utils';
import { GIT_PROVIDERS, type ProviderConnection } from './types';

export function ProviderAccounts({
  connections,
  selected,
  onSelect,
  onImportUrl,
  importingUrl = false,
}: {
  connections: ProviderConnection[];
  selected: string;
  onSelect: (provider: string) => void;
  onImportUrl?: () => void;
  importingUrl?: boolean;
}) {
  const connected = new Set(
    connections.filter((item) => item.connected).map((item) => item.provider)
  );
  return (
    <div className="flex min-w-0 items-center gap-2">
      <div className="flex min-w-0 items-center gap-1.5 overflow-x-auto p-0.5">
        {GIT_PROVIDERS.map((provider) => {
          const active = !importingUrl && selected === provider.id;
          const isConnected = connected.has(provider.id);
          return (
            <Button
              key={provider.id}
              type="button"
              variant="ghost"
              onClick={() => onSelect(provider.id)}
              aria-pressed={active}
              className={cn(
                'h-9 shrink-0 rounded-lg px-3',
                active ? 'bg-primary/10 text-primary' : 'text-muted-foreground'
              )}
            >
              <img src={provider.icon} alt="" aria-hidden="true" className="size-5 rounded-full" />
              {provider.name}
              <span
                aria-hidden="true"
                className={cn(
                  'size-1.5 rounded-full',
                  isConnected ? 'bg-emerald-500' : 'bg-muted-foreground/40'
                )}
              />
            </Button>
          );
        })}
      </div>
      {onImportUrl && (
        <div className="flex shrink-0 items-center">
          <Button
            type="button"
            variant="ghost"
            size="icon"
            onClick={onImportUrl}
            aria-label="Import from URL"
            title="Import from URL"
            aria-pressed={importingUrl}
            className={cn('rounded-lg', importingUrl && 'bg-primary/10 text-primary')}
          >
            <Link2 className="size-4" aria-hidden="true" />
          </Button>
        </div>
      )}
    </div>
  );
}
