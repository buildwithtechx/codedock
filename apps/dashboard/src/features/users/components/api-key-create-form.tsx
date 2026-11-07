import { Lock } from 'lucide-react';
import { Button } from '@codedock/ui/button';
import { Checkbox } from '@codedock/ui/checkbox';
import { Input } from '@codedock/ui/input';
import { Label } from '@codedock/ui/label';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@codedock/ui/select';

export const EXPIRY_OPTIONS = [
  { value: 'none', days: null, label: 'No expiration' },
  { value: '30', days: 30, label: '30 days' },
  { value: '90', days: 90, label: '90 days' },
  { value: '365', days: 365, label: '1 year' },
] as const;

interface ApiKeyCreateFormProps {
  name: string;
  setName: (name: string) => void;
  readOnly: boolean;
  setReadOnly: (ro: boolean) => void;
  expiryDays: number | null;
  setExpiryDays: (days: number | null) => void;
  projectScope: 'all' | 'specific';
  setProjectScope: (scope: 'all' | 'specific') => void;
  allowedProjects: string[];
  toggleProject: (projectId: string, checked: boolean) => void;
  projects: Array<{ id: string; name: string }>;
  isPending: boolean;
  onCancel: () => void;
  onSubmit: () => void;
}

export function ApiKeyCreateForm({
  name,
  setName,
  readOnly,
  setReadOnly,
  expiryDays,
  setExpiryDays,
  projectScope,
  setProjectScope,
  allowedProjects,
  toggleProject,
  projects,
  isPending,
  onCancel,
  onSubmit,
}: ApiKeyCreateFormProps) {
  return (
    <div className="mb-4 space-y-3 rounded-xl border border-border/50 p-4">
      <Input
        value={name}
        onChange={(e) => setName(e.target.value)}
        placeholder="Token name, e.g. CI deploys"
        className="bg-muted/30"
      />
      <div className="flex flex-wrap items-center gap-4">
        <div className="flex select-none items-center gap-2 text-foreground text-sm">
          <Checkbox
            id="token-read-only"
            checked={readOnly}
            onCheckedChange={(v) => setReadOnly(v === true)}
          />
          <Label htmlFor="token-read-only" className="cursor-pointer font-normal">
            Read-only
          </Label>
        </div>
        <Select
          value={expiryDays === null ? 'none' : String(expiryDays)}
          onValueChange={(v) => setExpiryDays(v === 'none' ? null : Number(v))}
        >
          <SelectTrigger className="w-40 bg-muted/30">
            <SelectValue placeholder="Expiration" />
          </SelectTrigger>
          <SelectContent>
            {EXPIRY_OPTIONS.map((o) => (
              <SelectItem key={o.value} value={o.value}>
                {o.label}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>

      <div className="rounded-lg border border-border/50 bg-muted/[0.04] p-3">
        <div className="grid grid-cols-2 gap-2">
          <button
            type="button"
            onClick={() => setProjectScope('all')}
            className={`flex h-9 items-center justify-center gap-2 rounded-lg border font-medium text-xs transition-colors ${
              projectScope === 'all'
                ? 'border-primary/50 bg-primary/10 text-primary'
                : 'border-border/50 bg-transparent text-muted-foreground hover:bg-background/50 hover:text-foreground'
            }`}
          >
            All projects
          </button>
          <button
            type="button"
            onClick={() => setProjectScope('specific')}
            className={`flex h-9 items-center justify-center gap-2 rounded-lg border font-medium text-xs transition-colors ${
              projectScope === 'specific'
                ? 'border-primary/50 bg-primary/10 text-primary'
                : 'border-border/50 bg-transparent text-muted-foreground hover:bg-background/50 hover:text-foreground'
            }`}
          >
            <Lock className="h-3.5 w-3.5" />
            Specific projects
          </button>
        </div>
        {projectScope === 'specific' && (
          <div className="mt-3 space-y-2 border-border/40 border-t pt-3">
            <p className="text-muted-foreground text-xs">Choose accessible projects:</p>
            {projects.length === 0 ? (
              <p className="py-2 text-muted-foreground text-xs">No projects available.</p>
            ) : (
              <div className="grid max-h-44 gap-2 overflow-y-auto sm:grid-cols-2">
                {projects.map((project) => {
                  const checked = allowedProjects.includes(project.id);
                  return (
                    <label
                      key={project.id}
                      className="flex cursor-pointer items-center gap-2 rounded-lg border border-border/40 bg-muted/20 px-2.5 py-1.5 text-xs hover:bg-muted/40"
                    >
                      <Checkbox
                        checked={checked}
                        onCheckedChange={(val) => toggleProject(project.id, val === true)}
                      />
                      <span className="truncate text-foreground">{project.name}</span>
                    </label>
                  );
                })}
              </div>
            )}
          </div>
        )}
      </div>

      <div className="flex justify-end gap-2 pt-1">
        <Button size="sm" variant="ghost" onClick={onCancel}>
          Cancel
        </Button>
        <Button size="sm" disabled={isPending} onClick={onSubmit}>
          {isPending ? 'Generating...' : 'Generate token'}
        </Button>
      </div>
    </div>
  );
}
