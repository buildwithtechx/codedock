import { format } from 'date-fns';
import { Calendar, Clock, FolderOpen, Key, Plus, Trash2 } from 'lucide-react';
import { useState } from 'react';
import { PageHeader } from '#/components/layout/page-header';
import { Button } from '#/components/ui/button';
import { useListTokens } from '#/features/profile';
import { ApiKeyCreateDialog } from './components/api-key-create-dialog';
import { ApiKeyDeleteDialog } from './components/api-key-delete-dialog';
import { ApiKeyNewDialog } from './components/api-key-new-dialog';

export function ApiKeysList() {
  const { data: tokensResponse, isLoading } = useListTokens();

  const [isCreateOpen, setIsCreateOpen] = useState(false);
  const [isNewKeyOpen, setIsNewKeyOpen] = useState(false);
  const [deleteId, setDeleteId] = useState<string | null>(null);
  const [newKeyPlain, setNewKeyPlain] = useState('');

  const handleCreateSuccess = (plainKey: string) => {
    setNewKeyPlain(plainKey);
    setIsCreateOpen(false);
    setIsNewKeyOpen(true);
  };

  const tokens = tokensResponse?.data || [];

  return (
    <div className="space-y-6">
      <PageHeader
        title="API access"
        description={
          isLoading
            ? 'Loading access tokens...'
            : `${tokens.length} API key${tokens.length === 1 ? '' : 's'}`
        }
        action={
          <Button onClick={() => setIsCreateOpen(true)} className="gap-2">
            <Plus className="h-4 w-4" />
            Create API key
          </Button>
        }
      />

      <div className="grid grid-cols-1 gap-6">
        {isLoading ? (
          <div className="flex min-h-[24rem] items-center justify-center">
            <span className="font-medium text-muted-foreground text-sm">
              Loading access tokens...
            </span>
          </div>
        ) : tokens.length === 0 ? (
          <div className="flex min-h-[28rem] flex-col items-center justify-center text-center">
            <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-card text-primary">
              <Key className="h-5 w-5 text-primary" />
            </div>
            <h2 className="mt-5 font-medium text-foreground/90 text-xl tracking-[-0.02em]">
              No API keys
            </h2>
            <p className="mt-2 max-w-sm text-muted-foreground/75 text-sm leading-6">
              Create an API key to access Codedock programmatically.
            </p>
            <Button onClick={() => setIsCreateOpen(true)} className="mt-6 gap-2">
              <Plus className="h-4 w-4" />
              Create API key
            </Button>
          </div>
        ) : (
          tokens.map((token) => (
            <div key={token.id} className="rounded-2xl bg-card p-6">
              <div className="mb-6 flex items-start justify-between">
                <div>
                  <div className="flex items-center gap-3">
                    <h2 className="font-semibold text-base text-foreground">{token.name}</h2>
                    <div className="rounded-full bg-primary/10 px-2 py-0.5 font-medium text-[10px] text-primary uppercase tracking-widest">
                      Active
                    </div>
                  </div>
                  <p className="mt-2 font-mono text-[10px] text-muted-foreground uppercase tracking-widest">
                    {token.prefix}
                  </p>
                </div>
                <Button
                  variant="outline"
                  onClick={() => setDeleteId(token.id)}
                  className="h-9 w-9 bg-transparent p-0 text-muted-foreground transition-colors hover:bg-destructive/10 hover:text-destructive"
                >
                  <Trash2 className="h-4 w-4" />
                </Button>
              </div>

              <div className="grid grid-cols-1 gap-4 md:grid-cols-4">
                <div className="rounded-xl bg-muted/50 p-4">
                  <div className="mb-3 flex items-center gap-2">
                    <Key className="h-3 w-3 text-muted-foreground" />
                    <span className="font-bold text-[10px] text-muted-foreground uppercase tracking-[0.15em]">
                      ACCESS
                    </span>
                  </div>
                  <div className="font-medium text-foreground/90 text-sm">
                    {token.accessLevel === 'read_write' ? 'Read and Write' : 'Read'}
                  </div>
                </div>

                <div className="rounded-xl bg-muted/50 p-4">
                  <div className="mb-3 flex items-center gap-2">
                    <FolderOpen className="h-3 w-3 text-muted-foreground" />
                    <span className="font-bold text-[10px] text-muted-foreground uppercase tracking-[0.15em]">
                      PROJECTS
                    </span>
                  </div>
                  <div className="font-medium text-foreground/90 text-sm">
                    {token.projectScope === 'specific' ? 'Specific projects' : 'All projects'}
                  </div>
                </div>

                <div className="rounded-xl bg-muted/50 p-4">
                  <div className="mb-3 flex items-center gap-2">
                    <Calendar className="h-3 w-3 text-muted-foreground" />
                    <span className="font-bold text-[10px] text-muted-foreground uppercase tracking-[0.15em]">
                      EXPIRATION
                    </span>
                  </div>
                  <div className="font-medium text-foreground/90 text-sm">
                    {token.expiresAt
                      ? `Expires ${format(new Date(token.expiresAt), 'MMM d, hh:mm a')}`
                      : 'No expiration'}
                  </div>
                </div>

                <div className="rounded-xl bg-muted/50 p-4">
                  <div className="mb-3 flex items-center gap-2">
                    <Clock className="h-3 w-3 text-muted-foreground" />
                    <span className="font-bold text-[10px] text-muted-foreground uppercase tracking-[0.15em]">
                      LAST USED
                    </span>
                  </div>
                  <div className="font-medium text-foreground/90 text-sm">Never</div>
                </div>
              </div>
            </div>
          ))
        )}
      </div>

      <ApiKeyCreateDialog
        open={isCreateOpen}
        onOpenChange={setIsCreateOpen}
        onSuccess={handleCreateSuccess}
      />
      <ApiKeyNewDialog
        open={isNewKeyOpen}
        onOpenChange={setIsNewKeyOpen}
        newKeyPlain={newKeyPlain}
        onClose={() => {
          setNewKeyPlain('');
          setIsNewKeyOpen(false);
        }}
      />
      <ApiKeyDeleteDialog deleteId={deleteId} onClose={() => setDeleteId(null)} />
    </div>
  );
}
