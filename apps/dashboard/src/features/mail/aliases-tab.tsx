import { Badge } from '@codedock/ui/badge';
import { Button } from '@codedock/ui/button';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '@codedock/ui/dialog';
import { Input } from '@codedock/ui/input';
import { Label } from '@codedock/ui/label';
import { AlertCircle, ArrowRight, Loader2, Plus, ShieldCheck, Trash2 } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { useCreateAlias, useDeleteAlias, useListAliases } from '#/hooks/use-mail';

export function AliasesTab({
  serverId,
  domainId,
  domainName,
}: {
  serverId: string;
  domainId: string;
  domainName: string;
}) {
  const { data: aliases, isLoading, isError, error } = useListAliases(serverId, domainId);
  const createMutation = useCreateAlias(serverId, domainId);
  const deleteMutation = useDeleteAlias(serverId, domainId);

  const [createOpen, setCreateOpen] = useState(false);
  const [sourcePrefix, setSourcePrefix] = useState('');
  const [destination, setDestination] = useState('');

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    const cleanSource = sourcePrefix.trim();
    const cleanDest = destination.trim();
    if (!cleanSource || !cleanDest) {
      toast.error('Source alias and destination mailbox are required');
      return;
    }
    const fullSource = cleanSource.includes('@') ? cleanSource : `${cleanSource}@${domainName}`;

    try {
      await createMutation.mutateAsync({
        source: fullSource,
        destination: cleanDest,
      });
      toast.success(`Created alias ${fullSource} -> ${cleanDest}`);
      setSourcePrefix('');
      setDestination('');
      setCreateOpen(false);
    } catch (err: any) {
      toast.error(err?.message || 'Failed to create alias');
    }
  };

  const handleDelete = async (aliasId: string, source: string) => {
    if (!confirm(`Delete alias ${source}?`)) return;
    try {
      await deleteMutation.mutateAsync(aliasId);
      toast.success(`Deleted alias ${source}`);
    } catch (err: any) {
      toast.error(err?.message || 'Failed to delete alias');
    }
  };

  return (
    <div className="space-y-4">
      <div className="flex flex-wrap items-center justify-between gap-3">
        <div>
          <h3 className="font-semibold text-foreground text-sm">Mail Aliases</h3>
          <p className="text-muted-foreground text-xs">
            Forward mail received at alias addresses to active destination mailboxes.
          </p>
        </div>
        <Button size="sm" onClick={() => setCreateOpen(true)} className="h-8 gap-1.5 text-xs">
          <Plus className="size-3.5" />
          Create alias
        </Button>
      </div>

      {/* Strict Sender Notice */}
      <div className="flex items-start gap-2.5 rounded-xl border border-border/60 bg-muted/30 p-3.5 text-muted-foreground text-xs">
        <ShieldCheck className="mt-0.5 size-4 shrink-0 text-primary" />
        <p className="leading-relaxed">
          <strong className="text-foreground">Strict direct alias ownership:</strong> Authenticated
          mailboxes can send outbound mail using an alias only when they are its direct destination.
          External forwards and catch-alls only route incoming mail.
        </p>
      </div>

      {isLoading ? (
        <div className="flex items-center justify-center py-8 text-muted-foreground text-xs">
          <Loader2 className="mr-2 size-4 animate-spin" />
          Loading aliases...
        </div>
      ) : isError ? (
        <div className="flex items-center gap-2 rounded-xl border border-destructive/30 bg-destructive/10 p-4 text-destructive text-xs">
          <AlertCircle className="size-4 shrink-0" />
          <span>
            Failed to load aliases: {error instanceof Error ? error.message : 'Unknown error'}
          </span>
        </div>
      ) : (aliases || []).length === 0 ? (
        <div className="rounded-xl border border-border/60 border-dashed p-6 text-center text-muted-foreground text-xs">
          No aliases configured for {domainName}. Create an alias to route email to existing
          mailboxes.
        </div>
      ) : (
        <div className="divide-y divide-border/40 overflow-hidden rounded-xl border border-border/60 bg-card shadow-xs">
          {(aliases || []).map((alias) => (
            <div key={alias.id} className="flex items-center justify-between p-3.5 text-xs">
              <div className="flex min-w-0 items-center gap-2">
                <span className="font-medium font-mono text-foreground">{alias.source}</span>
                <ArrowRight className="size-3 shrink-0 text-muted-foreground" />
                <span className="truncate font-mono text-muted-foreground">
                  {alias.destination}
                </span>
                {alias.active && (
                  <Badge variant="outline" className="ml-1 text-[10px]">
                    Active
                  </Badge>
                )}
              </div>
              <Button
                size="sm"
                variant="ghost"
                disabled={deleteMutation.isPending}
                onClick={() => handleDelete(alias.id, alias.source)}
                className="h-7 w-7 p-0 text-muted-foreground hover:text-destructive"
              >
                <Trash2 className="size-3.5" />
              </Button>
            </div>
          ))}
        </div>
      )}

      {/* Create Alias Modal */}
      <Dialog open={createOpen} onOpenChange={setCreateOpen}>
        <DialogContent className="max-w-md">
          <DialogHeader>
            <DialogTitle>Create Mail Alias</DialogTitle>
            <DialogDescription>
              Create a virtual address that delivers to a destination mailbox on {domainName}.
            </DialogDescription>
          </DialogHeader>

          <form onSubmit={handleCreate} className="space-y-4 pt-2">
            <div className="space-y-1.5">
              <Label className="text-xs">Source Alias</Label>
              <div className="flex items-center gap-1.5">
                <Input
                  placeholder="support"
                  value={sourcePrefix}
                  onChange={(e) => setSourcePrefix(e.target.value)}
                  className="h-8 font-mono text-xs"
                  required
                />
                <span className="font-mono text-muted-foreground text-xs">@{domainName}</span>
              </div>
            </div>

            <div className="space-y-1.5">
              <Label className="text-xs">Destination Mailbox</Label>
              <Input
                placeholder="alice@example.com"
                value={destination}
                onChange={(e) => setDestination(e.target.value)}
                className="h-8 font-mono text-xs"
                required
              />
              <p className="text-[11px] text-muted-foreground">
                Signing in as this mailbox grants permission to send mail with this alias.
              </p>
            </div>

            <div className="flex justify-end gap-2 border-border/40 border-t pt-2">
              <Button
                type="button"
                variant="outline"
                size="sm"
                onClick={() => setCreateOpen(false)}
                className="h-8 text-xs"
              >
                Cancel
              </Button>
              <Button
                type="submit"
                size="sm"
                disabled={createMutation.isPending}
                className="h-8 text-xs"
              >
                {createMutation.isPending && <Loader2 className="mr-1 size-3.5 animate-spin" />}
                Create Alias
              </Button>
            </div>
          </form>
        </DialogContent>
      </Dialog>
    </div>
  );
}
