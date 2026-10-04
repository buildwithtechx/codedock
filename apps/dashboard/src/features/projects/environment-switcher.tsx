import { Check, ChevronDown, Layers, Plus } from 'lucide-react';
import { useState } from 'react';
import { Badge } from '#/components/ui/badge';
import { Button } from '#/components/ui/button';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '#/components/ui/dialog';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '#/components/ui/dropdown-menu';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import { useCreateEnvironment, useListByProject } from '#/hooks/use-environments';

interface EnvironmentSwitcherProps {
  projectId: string;
  activeEnvironmentId?: string;
  onSelectEnvironment: (envId: string) => void;
}

export function EnvironmentSwitcher({
  projectId,
  activeEnvironmentId,
  onSelectEnvironment,
}: EnvironmentSwitcherProps) {
  const [createOpen, setCreateOpen] = useState(false);
  const [newEnvName, setNewEnvName] = useState('');
  const { data: envsResponse, isLoading } = useListByProject(projectId);
  const createEnvMutation = useCreateEnvironment();

  const environments = envsResponse?.data || [];
  const activeEnv = environments.find((e) => e.id === activeEnvironmentId) || environments[0];

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newEnvName.trim()) return;

    try {
      const res = await createEnvMutation.mutateAsync({
        projectId,
        name: newEnvName.trim(),
      });
      setNewEnvName('');
      setCreateOpen(false);
      if (res?.data?.id) {
        onSelectEnvironment(res.data.id);
      }
    } catch {
      // Handled by mutation
    }
  };

  return (
    <>
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button
            variant="outline"
            size="sm"
            className="flex h-9 items-center gap-2 rounded-xl border-border/80 bg-background/60 px-3 font-medium text-xs shadow-xs hover:bg-muted/60"
            disabled={isLoading}
          >
            <Layers className="h-3.5 w-3.5 text-primary" />
            <span className="max-w-30 truncate font-semibold">
              {activeEnv?.name || 'Environment'}
            </span>
            {activeEnv && (
              <Badge
                variant="secondary"
                className="h-5 px-1.5 py-0 text-[10px] uppercase tracking-wider"
              >
                {activeEnv.name.toLowerCase().includes('prod') ? 'Prod' : 'Env'}
              </Badge>
            )}
            <ChevronDown className="h-3 w-3 text-muted-foreground opacity-60" />
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="start" className="w-56 rounded-xl p-1">
          <div className="px-2 py-1.5 font-semibold text-[11px] text-muted-foreground uppercase tracking-wider">
            Environments
          </div>
          {environments.map((env) => {
            const isSelected = activeEnv?.id === env.id;
            return (
              <DropdownMenuItem
                key={env.id}
                onClick={() => onSelectEnvironment(env.id)}
                className="flex items-center justify-between gap-2 rounded-lg px-2.5 py-2 text-xs"
              >
                <div className="flex items-center gap-2 truncate">
                  <span className="truncate font-medium">{env.name}</span>
                </div>
                {isSelected && <Check className="h-3.5 w-3.5 text-primary" />}
              </DropdownMenuItem>
            );
          })}
          <DropdownMenuSeparator />
          <DropdownMenuItem
            onClick={() => setCreateOpen(true)}
            className="flex items-center gap-2 rounded-lg px-2.5 py-2 text-primary text-xs"
          >
            <Plus className="h-3.5 w-3.5" />
            <span>New environment</span>
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>

      <Dialog open={createOpen} onOpenChange={setCreateOpen}>
        <DialogContent className="sm:max-w-md">
          <form onSubmit={handleCreate}>
            <DialogHeader>
              <DialogTitle>Create environment</DialogTitle>
              <DialogDescription>
                Environments let you run isolated instances of your services and configurations.
              </DialogDescription>
            </DialogHeader>
            <div className="space-y-4 py-4">
              <div className="space-y-2">
                <Label htmlFor="env-name">Environment name</Label>
                <Input
                  id="env-name"
                  placeholder="e.g. staging, preview, testing"
                  value={newEnvName}
                  onChange={(e) => setNewEnvName(e.target.value)}
                  autoFocus
                />
              </div>
            </div>
            <DialogFooter>
              <Button type="button" variant="outline" onClick={() => setCreateOpen(false)}>
                Cancel
              </Button>
              <Button type="submit" disabled={!newEnvName.trim() || createEnvMutation.isPending}>
                {createEnvMutation.isPending ? 'Creating...' : 'Create environment'}
              </Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>
    </>
  );
}
