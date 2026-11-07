import { FolderGit2, Search } from 'lucide-react';
import { Input } from '@codedock/ui/input';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@codedock/ui/select';

export const DEPLOYMENT_FILTERS = [
  { value: 'all', label: 'All' },
  { value: 'success', label: 'Success' },
  { value: 'failed', label: 'Failed' },
  { value: 'building', label: 'Building' },
  { value: 'active', label: 'In Progress' },
  { value: 'canceled', label: 'Canceled' },
] as const;

export type DeploymentFilterValue = (typeof DEPLOYMENT_FILTERS)[number]['value'];

interface DeploymentFiltersProps {
  search: string;
  projectId: string;
  status: string;
  projects: Array<{ id: string; name: string }>;
  onSearchChange: (value: string) => void;
  onProjectChange: (value: string) => void;
  onStatusChange: (value: string) => void;
}

export function DeploymentFilters({
  search,
  projectId,
  status,
  projects,
  onSearchChange,
  onProjectChange,
  onStatusChange,
}: DeploymentFiltersProps) {
  return (
    <div className="flex flex-col gap-3 sm:flex-row sm:flex-wrap sm:items-center">
      <div className="relative w-full sm:min-w-55 sm:flex-1">
        <Search className="pointer-events-none absolute top-1/2 left-3.5 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
        <Input
          type="text"
          value={search}
          onChange={(event) => onSearchChange(event.target.value)}
          placeholder="Search deployments..."
          aria-label="Search deployments..."
          className="h-10 bg-muted/60 ps-10 pe-4"
        />
      </div>

      <Select value={projectId} onValueChange={onProjectChange}>
        <SelectTrigger className="flex w-full items-center gap-2 bg-card sm:w-44">
          <FolderGit2 className="h-4 w-4 shrink-0 text-muted-foreground" />
          <SelectValue placeholder="All Projects" />
        </SelectTrigger>
        <SelectContent>
          <SelectItem value="all">All Projects</SelectItem>
          {projects.map((project) => (
            <SelectItem key={project.id} value={project.id}>
              {project.name}
            </SelectItem>
          ))}
        </SelectContent>
      </Select>

      <div className="scrollbar-hide scrollbar-none max-w-full shrink-0 overflow-x-auto [&::-webkit-scrollbar]:hidden">
        <div className="inline-flex items-center gap-1">
          {DEPLOYMENT_FILTERS.map((f) => (
            <button
              key={f.value}
              type="button"
              aria-pressed={status === f.value}
              onClick={() => onStatusChange(f.value)}
              className={`inline-flex h-10 shrink-0 items-center whitespace-nowrap rounded-lg px-4 font-medium text-sm transition-colors focus-visible:outline-2 focus-visible:outline-ring ${
                status === f.value
                  ? 'bg-foreground text-background'
                  : 'text-muted-foreground hover:bg-muted hover:text-foreground'
              }`}
            >
              {f.label}
            </button>
          ))}
        </div>
      </div>
    </div>
  );
}
