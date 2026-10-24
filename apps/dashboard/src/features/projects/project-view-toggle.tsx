import { LayoutGrid, List } from 'lucide-react';

export type ProjectView = 'list' | 'grid';

const options: { id: ProjectView; label: string; icon: typeof List }[] = [
  { id: 'list', label: 'List', icon: List },
  { id: 'grid', label: 'Grid', icon: LayoutGrid },
];

export function ProjectViewToggle({
  value,
  onChange,
}: {
  value: ProjectView;
  onChange: (view: ProjectView) => void;
}) {
  return (
    <fieldset className="inline-flex shrink-0 items-center gap-0.5 rounded-xl bg-muted/50 p-1">
      <legend className="sr-only">Project layout</legend>
      {options.map((option) => {
        const active = value === option.id;
        return (
          <button
            key={option.id}
            type="button"
            onClick={() => onChange(option.id)}
            aria-pressed={active}
            aria-label={option.label}
            className={`inline-flex items-center gap-1.5 rounded-lg px-3 py-1.5 font-medium text-[13px] transition-colors ${
              active ? 'bg-card text-foreground' : 'text-muted-foreground hover:text-foreground'
            }`}
          >
            <option.icon className="h-4 w-4" />
            <span className="hidden sm:inline">{option.label}</span>
          </button>
        );
      })}
    </fieldset>
  );
}
