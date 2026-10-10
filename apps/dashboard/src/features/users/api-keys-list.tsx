import { format } from 'date-fns';
import { Check, Copy, Key, Lock, Plus, ShieldCheck, Trash2 } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Badge } from '#/components/ui/badge';
import { Button } from '#/components/ui/button';
import { Checkbox } from '#/components/ui/checkbox';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '#/components/ui/select';
import { useCreateToken, useListTokens } from '#/features/profile';
import { SettingsSection } from '#/features/settings/settings-section';
import type { CreatePATRequest } from '#/features/users';
import { useListCanvasSummaries } from '#/hooks/use-canvas';
import { ApiKeyDeleteDialog } from './components/api-key-delete-dialog';

const EXPIRY_OPTIONS = [
  { value: 'none', days: null, label: 'No expiration' },
  { value: '30', days: 30, label: '30 days' },
  { value: '90', days: 90, label: '90 days' },
  { value: '365', days: 365, label: '1 year' },
] as const;

function formatDate(iso: string | undefined): string {
  if (!iso) return 'Never expires';
  return format(new Date(iso), 'MMM d, yyyy');
}

export function ApiKeysList() {
  const { data: tokensResponse, isLoading } = useListTokens();
  const createToken = useCreateToken();
  const projectsQuery = useListCanvasSummaries();
  const projects = projectsQuery.data?.data ?? [];

  const [showForm, setShowForm] = useState(false);
  const [name, setName] = useState('');
  const [readOnly, setReadOnly] = useState(false);
  const [expiryDays, setExpiryDays] = useState<number | null>(30);
  const [projectScope, setProjectScope] = useState<'all' | 'specific'>('all');
  const [allowedProjects, setAllowedProjects] = useState<string[]>([]);
  const [newToken, setNewToken] = useState<string | null>(null);
  const [copied, setCopied] = useState(false);
  const [deleteId, setDeleteId] = useState<string | null>(null);

  const tokens = tokensResponse?.data ?? [];

  const resetForm = () => {
    setShowForm(false);
    setName('');
    setReadOnly(false);
    setExpiryDays(30);
    setProjectScope('all');
    setAllowedProjects([]);
  };

  const handleCreate = () => {
    if (!name.trim()) {
      toast.error('Give the token a name');
      return;
    }
    if (projectScope === 'specific' && allowedProjects.length === 0) {
      toast.error('Select at least one project or choose all projects');
      return;
    }
    const payload: CreatePATRequest = {
      name: name.trim(),
      accessLevel: readOnly ? 'read' : 'read_write',
      projectScope,
      allowedProjects: projectScope === 'specific' ? allowedProjects : [],
    };
    if (expiryDays !== null) {
      const date = new Date();
      date.setDate(date.getDate() + expiryDays);
      payload.expiresAt = date.toISOString();
    }
    createToken.mutate(
      { payload },
      {
        onSuccess: (data) => {
          setNewToken(data.data.plain);
          resetForm();
        },
        onError: (err) => {
          toast.error(err.message || 'Failed to create token');
        },
      }
    );
  };

  const copyToken = async () => {
    if (!newToken) return;
    try {
      await navigator.clipboard.writeText(newToken);
      setCopied(true);
      toast.success('Token copied to clipboard');
      setTimeout(() => setCopied(false), 1500);
    } catch {
      toast.error('Copy failed, select the token manually');
    }
  };

  const toggleProject = (projectId: string, checked: boolean) => {
    setAllowedProjects((current) =>
      checked ? [...current, projectId] : current.filter((id) => id !== projectId)
    );
  };

  return (
    <SettingsSection
      icon={<Key className="size-4 text-primary" />}
      title="Personal access tokens"
      description="Tokens authenticate with the API as your user account."
      action={
        !showForm ? (
          <Button size="sm" onClick={() => setShowForm(true)} className="gap-1.5">
            <Plus className="size-3.5" />
            New token
          </Button>
        ) : undefined
      }
    >
      <div className="mb-4 flex gap-2.5 rounded-xl border border-border/50 bg-muted/30 p-3">
        <ShieldCheck className="mt-0.5 h-4 w-4 shrink-0 text-muted-foreground" />
        <div className="text-muted-foreground text-xs leading-relaxed">
          A token <span className="font-medium text-foreground">acts as you</span>, with full access
          to everything you can reach. Prefer{' '}
          <span className="font-medium text-foreground">read-only</span> tokens and{' '}
          <span className="font-medium text-foreground">limit them to specific projects</span>{' '}
          whenever possible.
        </div>
      </div>

      {newToken && (
        <div className="mb-4 rounded-xl border border-primary/30 bg-primary/[0.04] p-4">
          <p className="mb-1 font-medium text-foreground text-sm">Copy your new token</p>
          <p className="mb-3 text-muted-foreground text-xs">
            This is the only time the full token is shown. Store it somewhere secure.
          </p>
          <div className="flex items-center gap-2">
            <code className="min-w-0 flex-1 truncate rounded-lg bg-muted px-3 py-2 font-mono text-foreground text-xs">
              {newToken}
            </code>
            <Button size="sm" onClick={copyToken} className="shrink-0 gap-1.5">
              {copied ? <Check className="h-3.5 w-3.5" /> : <Copy className="h-3.5 w-3.5" />}
              {copied ? 'Copied' : 'Copy'}
            </Button>
          </div>
          <button
            type="button"
            onClick={() => setNewToken(null)}
            className="mt-3 font-medium text-muted-foreground text-xs transition-colors hover:text-foreground"
          >
            Dismiss
          </button>
        </div>
      )}

      {showForm ? (
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
              <div className="mt-3 max-h-40 space-y-2 overflow-y-auto rounded-lg border border-border/50 p-3">
                {projectsQuery.isLoading && (
                  <p className="text-muted-foreground text-sm">Loading projects...</p>
                )}
                {projectsQuery.isError && (
                  <div className="space-y-2">
                    <p className="text-sm">Could not load projects.</p>
                    <Button
                      type="button"
                      variant="outline"
                      size="sm"
                      onClick={() => projectsQuery.refetch()}
                    >
                      Retry projects
                    </Button>
                  </div>
                )}
                {projects.map((project) => (
                  <div key={project.id} className="flex items-center gap-2 text-sm">
                    <Checkbox
                      id={`token-project-${project.id}`}
                      checked={allowedProjects.includes(project.id)}
                      onCheckedChange={(v) => toggleProject(project.id, v === true)}
                    />
                    <Label
                      htmlFor={`token-project-${project.id}`}
                      className="cursor-pointer font-normal"
                    >
                      {project.name}
                    </Label>
                  </div>
                ))}
              </div>
            )}
          </div>

          <div className="flex items-center gap-2">
            <Button onClick={handleCreate} disabled={createToken.isPending} className="gap-2">
              <Plus className="h-4 w-4" />
              {createToken.isPending ? 'Creating...' : 'Create token'}
            </Button>
            <Button variant="ghost" onClick={resetForm}>
              Cancel
            </Button>
          </div>
        </div>
      ) : (
        <Button variant="secondary" onClick={() => setShowForm(true)} className="mb-4 gap-2">
          <Plus className="h-4 w-4" />
          New token
        </Button>
      )}

      {isLoading ? (
        <p className="py-3 text-muted-foreground text-sm">Loading tokens...</p>
      ) : tokens.length === 0 ? (
        <p className="py-2 text-muted-foreground text-sm">
          No active tokens. Create one to access the API programmatically.
        </p>
      ) : (
        <div className="divide-y divide-border/50">
          {tokens.map((token) => (
            <div key={token.id} className="flex items-center justify-between gap-3 py-3">
              <div className="flex min-w-0 items-center gap-3">
                <Key className="h-4 w-4 shrink-0 text-muted-foreground" />
                <div className="min-w-0">
                  <div className="flex flex-wrap items-center gap-2">
                    <p className="truncate font-medium text-foreground text-sm">{token.name}</p>
                    {token.accessLevel === 'read' && (
                      <Badge variant="secondary" className="gap-1 text-[10px]">
                        <ShieldCheck className="h-3 w-3" />
                        Read-only
                      </Badge>
                    )}
                    {token.projectScope === 'specific' && (
                      <Badge className="gap-1 bg-primary/15 text-[10px] text-primary hover:bg-primary/15">
                        <Lock className="h-3 w-3" />
                        Scoped
                      </Badge>
                    )}
                  </div>
                  <p className="mt-0.5 font-mono text-muted-foreground text-xs">
                    {token.prefix}...{' '}
                    <span className="font-sans">Created {formatDate(token.createdAt)}</span>
                  </p>
                </div>
              </div>
              <div className="flex shrink-0 items-center gap-3">
                <span className="text-muted-foreground text-xs">
                  Expires {formatDate(token.expiresAt)}
                </span>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => setDeleteId(token.id)}
                  className="shrink-0 gap-1.5 text-muted-foreground hover:border-destructive/50 hover:text-destructive"
                >
                  <Trash2 className="h-3.5 w-3.5" />
                  Revoke
                </Button>
              </div>
            </div>
          ))}
        </div>
      )}

      <ApiKeyDeleteDialog deleteId={deleteId} onClose={() => setDeleteId(null)} />
    </SettingsSection>
  );
}
