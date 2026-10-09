import { GitBranch, Loader2, Search } from 'lucide-react';
import { useMemo, useState } from 'react';
import { Button } from '#/components/ui/button';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '#/components/ui/dialog';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '#/components/ui/select';
import { useGetStatus, useListGitRepos } from '#/hooks/use-git';
import type { GitRepo } from '#/interfaces/git';

interface ProjectGitRepoDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onSelect: (repository: GitRepo) => void;
  title: string;
  description: string;
}

export function ProjectGitRepoDialog({
  open,
  onOpenChange,
  onSelect,
  title,
  description,
}: ProjectGitRepoDialogProps) {
  const { data: statusRes } = useGetStatus();
  const connections = useMemo(
    () => (statusRes?.data ?? []).filter((connection) => connection.connected),
    [statusRes]
  );
  const [provider, setProvider] = useState('');
  const [query, setQuery] = useState('');
  const { data: reposRes, isLoading } = useListGitRepos(provider);
  const repositories = useMemo(() => {
    const repos = reposRes?.data ?? [];
    const needle = query.trim().toLowerCase();
    if (!needle) return repos;
    return repos.filter((repository) => repository.fullName.toLowerCase().includes(needle));
  }, [reposRes, query]);

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-h-[80vh] overflow-hidden sm:max-w-xl">
        <DialogHeader>
          <DialogTitle>{title}</DialogTitle>
          <DialogDescription>{description}</DialogDescription>
        </DialogHeader>
        <div className="space-y-2">
          <Label>Git account</Label>
          <Select value={provider} onValueChange={setProvider}>
            <SelectTrigger>
              <SelectValue placeholder="Select a connected account" />
            </SelectTrigger>
            <SelectContent>
              {connections.map((connection) => (
                <SelectItem key={connection.provider} value={connection.provider}>
                  {connection.provider}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
        {provider && (
          <div className="relative">
            <Search className="pointer-events-none absolute top-1/2 left-3 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
            <Input
              value={query}
              onChange={(event) => setQuery(event.target.value)}
              placeholder="Filter repositories"
              className="pl-9"
            />
          </div>
        )}
        <div className="max-h-72 overflow-y-auto rounded-xl border border-border/60">
          {!provider ? (
            <p className="p-5 text-center text-muted-foreground text-sm">
              Connect a git account in Settings to browse its repositories.
            </p>
          ) : isLoading ? (
            <div className="flex items-center justify-center gap-2 p-8 text-muted-foreground text-sm">
              <Loader2 className="h-4 w-4 animate-spin" />
              Loading repositories
            </div>
          ) : repositories.length === 0 ? (
            <p className="p-5 text-center text-muted-foreground text-sm">
              No repositories found for this account.
            </p>
          ) : (
            <ul className="divide-y divide-border/50">
              {repositories.map((repository) => (
                <li key={repository.id}>
                  <button
                    type="button"
                    onClick={() => onSelect(repository)}
                    className="flex w-full items-center gap-3 px-4 py-2.5 text-left transition-colors hover:bg-muted/50"
                  >
                    <GitBranch className="h-4 w-4 shrink-0 text-muted-foreground" />
                    <span className="min-w-0 flex-1">
                      <span className="block truncate font-medium text-sm">
                        {repository.fullName}
                      </span>
                      <span className="block text-muted-foreground text-xs">
                        {repository.private ? 'Private' : 'Public'} · {repository.defaultBranch}
                      </span>
                    </span>
                  </button>
                </li>
              ))}
            </ul>
          )}
        </div>
        <div className="flex justify-end">
          <Button variant="ghost" onClick={() => onOpenChange(false)}>
            Cancel
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  );
}
