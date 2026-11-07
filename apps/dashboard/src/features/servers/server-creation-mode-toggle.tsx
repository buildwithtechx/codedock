import { Cloud, Cpu, Server as ServerIcon } from 'lucide-react';
import { cn } from '#/lib/utils';

export type ServerCreationMode = 'remote' | 'local' | 'hetzner';

const modes = [
  { id: 'remote', label: 'Remote SSH Server', icon: ServerIcon },
  { id: 'local', label: 'Local Docker Node', icon: Cpu },
  { id: 'hetzner', label: 'Hetzner Cloud', icon: Cloud },
] as const;

export function ServerCreationModeToggle({
  mode,
  onChange,
}: {
  mode: ServerCreationMode;
  onChange: (mode: ServerCreationMode) => void;
}) {
  return (
    <div className="flex items-center gap-2">
      {modes.map((entry) => (
        <button
          key={entry.id}
          type="button"
          onClick={() => onChange(entry.id)}
          className={cn(
            'flex items-center gap-2 rounded-xl border px-4 py-2.5 font-medium text-xs transition-colors',
            mode === entry.id
              ? 'border-primary bg-primary/10 text-primary'
              : 'border-border/80 bg-card text-muted-foreground hover:bg-muted/40 hover:text-foreground'
          )}
        >
          <entry.icon className="h-4 w-4" />
          {entry.label}
        </button>
      ))}
    </div>
  );
}
