import { format } from 'date-fns';
import { Check, Copy, Key, Plus, ShieldCheck, Trash2 } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Badge } from '@codedock/ui/badge';
import { Button } from '@codedock/ui/button';
import { useCreateToken, useListTokens } from '#/hooks/use-profile';
import { SettingsSection } from '#/features/settings/settings-section';
import type { CreatePATRequest } from '#/interfaces/users';
import { useListCanvasSummaries } from '#/hooks/use-canvas';
import { ApiKeyDeleteDialog } from './components/api-key-delete-dialog';
import { ApiKeyCreateForm } from './components/api-key-create-form';

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

      {showForm && (
        <ApiKeyCreateForm
          name={name}
          setName={setName}
          readOnly={readOnly}
          setReadOnly={setReadOnly}
          expiryDays={expiryDays}
          setExpiryDays={setExpiryDays}
          projectScope={projectScope}
          setProjectScope={setProjectScope}
          allowedProjects={allowedProjects}
          toggleProject={toggleProject}
          projects={projects}
          isPending={createToken.isPending}
          onCancel={resetForm}
          onSubmit={handleCreate}
        />
      )}

      {isLoading ? (
        <div className="space-y-2">
          {[1, 2].map((i) => (
            <div key={i} className="h-16 animate-pulse rounded-xl bg-muted/40" />
          ))}
        </div>
      ) : tokens.length === 0 ? (
        <div className="rounded-xl border border-border/50 border-dashed py-8 text-center text-muted-foreground text-sm">
          No personal access tokens yet. Create one to authenticate API requests.
        </div>
      ) : (
        <div className="divide-y divide-border/50 rounded-xl border border-border/50">
          {tokens.map((token) => (
            <div key={token.id} className="flex items-center justify-between p-3.5">
              <div className="space-y-1">
                <div className="flex items-center gap-2">
                  <span className="font-medium text-foreground text-sm">{token.name}</span>
                  <Badge
                    variant={token.accessLevel === 'read' ? 'secondary' : 'default'}
                    className="text-[10px]"
                  >
                    {token.accessLevel === 'read' ? 'Read-only' : 'Read/Write'}
                  </Badge>
                  {token.projectScope === 'specific' && (
                    <Badge variant="outline" className="text-[10px]">
                      {token.allowedProjects?.length ?? 0} project
                      {(token.allowedProjects?.length ?? 0) === 1 ? '' : 's'}
                    </Badge>
                  )}
                </div>
                <div className="flex items-center gap-3 text-muted-foreground text-xs">
                  <span>
                    Prefix: <code className="font-mono">{token.prefix}••••</code>
                  </span>
                  <span>Expires: {formatDate(token.expiresAt)}</span>
                  <span>Created: {formatDate(token.createdAt)}</span>
                </div>
              </div>
              <Button
                variant="ghost"
                size="icon"
                onClick={() => setDeleteId(token.id)}
                className="text-muted-foreground hover:text-destructive"
                aria-label={`Revoke token ${token.name}`}
              >
                <Trash2 className="size-4" />
              </Button>
            </div>
          ))}
        </div>
      )}

      <ApiKeyDeleteDialog deleteId={deleteId} onClose={() => setDeleteId(null)} />
    </SettingsSection>
  );
}
