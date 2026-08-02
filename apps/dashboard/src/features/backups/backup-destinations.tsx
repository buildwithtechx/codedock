import { ArrowRight, Database, MoreVertical, Plus, Trash } from 'lucide-react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '#/components/ui/dropdown-menu';
import { QueryErrorState } from '#/components/ui/query-error-state';
import { useDeleteS3Destination, useListS3Destinations } from '#/features/backups';

type BackupDestinationsProps = {
  onAddDestination: () => void;
};

export function BackupDestinations({ onAddDestination }: BackupDestinationsProps) {
  const { data: s3Destinations, isLoading, isError, refetch } = useListS3Destinations();
  const deleteS3Dest = useDeleteS3Destination();

  const handleDelete = async (id: string) => {
    try {
      await deleteS3Dest.mutateAsync({ id, projectId: 'global' });
      toast.success('S3 destination deleted successfully');
    } catch {
      toast.error('Failed to delete S3 destination');
    }
  };

  if (isLoading) {
    return <div className="min-h-100 animate-pulse" />;
  }

  if (isError) {
    return (
      <QueryErrorState
        title="Storage destinations are unavailable"
        description="Codedock could not load the configured backup storage."
        onRetry={() => refetch()}
      />
    );
  }

  const list = s3Destinations?.data || [];

  return (
    <div>
      {list.length === 0 ? (
        <section className="flex min-h-[26rem] items-center justify-center px-6 py-12 text-center">
          <div className="max-w-md">
            <div className="relative mx-auto h-28 w-52" aria-hidden="true">
              <span className="absolute top-4 left-11 h-17 w-26 rounded-2xl bg-card" />
              <span className="absolute top-1 left-17 flex h-17 w-26 items-center justify-center rounded-2xl bg-muted text-muted-foreground">
                <Database className="h-6 w-6" />
              </span>
              <span className="absolute top-10 right-5 flex h-10 w-10 items-center justify-center rounded-xl border border-primary/45 border-dashed bg-primary/8 text-primary">
                <Plus className="h-4 w-4" />
              </span>
              <span className="absolute top-1/2 right-14 w-10 border-border border-t border-dashed" />
            </div>
            <h2 className="mt-4 font-medium text-foreground/90 text-xl tracking-[-0.02em]">
              Add a storage destination
            </h2>
            <p className="mt-2 text-muted-foreground/75 text-sm leading-6">
              Connect S3-compatible storage to keep instance backups outside this server.
            </p>
            <Button className="mt-6 gap-2" onClick={onAddDestination}>
              Add destination
              <ArrowRight className="h-4 w-4" />
            </Button>
          </div>
        </section>
      ) : (
        <div className="grid grid-cols-1 gap-4 md:grid-cols-2 xl:grid-cols-3">
          {list.map((dest) => (
            <div
              key={dest.id}
              className="relative flex flex-col justify-between rounded-xl bg-card p-5 transition-colors hover:bg-muted/60"
            >
              <div>
                <div className="flex items-start justify-between">
                  <div>
                    <h2 className="font-semibold text-base">{dest.name}</h2>
                    {dest.description && (
                      <p className="mt-0.5 text-muted-foreground text-xs">{dest.description}</p>
                    )}
                  </div>
                  <DropdownMenu>
                    <DropdownMenuTrigger asChild>
                      <Button variant="ghost" size="icon" className="h-8 w-8 text-muted-foreground">
                        <MoreVertical className="h-4 w-4" />
                      </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent align="end">
                      <DropdownMenuItem
                        className="text-destructive focus:text-destructive"
                        onClick={() => handleDelete(dest.id)}
                      >
                        <Trash className="mr-2 h-4 w-4" />
                        Delete
                      </DropdownMenuItem>
                    </DropdownMenuContent>
                  </DropdownMenu>
                </div>

                <dl className="mt-4 space-y-2 font-mono text-xs">
                  <div className="flex justify-between border-border/40 border-b pb-1">
                    <dt className="text-muted-foreground">Provider</dt>
                    <dd className="font-medium uppercase">{dest.provider}</dd>
                  </div>
                  <div className="flex justify-between border-border/40 border-b pb-1">
                    <dt className="text-muted-foreground">Bucket</dt>
                    <dd className="truncate pl-2 text-foreground">{dest.bucket}</dd>
                  </div>
                  <div className="flex justify-between border-border/40 border-b pb-1">
                    <dt className="text-muted-foreground">Region</dt>
                    <dd className="text-foreground">{dest.region || 'auto'}</dd>
                  </div>
                  <div className="flex justify-between">
                    <dt className="text-muted-foreground">Endpoint</dt>
                    <dd className="max-w-50 truncate pl-2 text-foreground" title={dest.endpoint}>
                      {dest.endpoint}
                    </dd>
                  </div>
                </dl>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
