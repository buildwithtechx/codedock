import { Database, MoreVertical, Plus, Trash } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { PageHeader } from '#/components/layout/page-header';
import { Button } from '#/components/ui/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '#/components/ui/dropdown-menu';
import { useDeleteS3Destination, useListS3Destinations } from '#/features/backups';
import { CreateS3DestinationDialog } from './create-s3-destination-dialog';

export function S3DestinationsList() {
  const [isOpen, setIsOpen] = useState(false);
  const { data: s3Destinations, isLoading } = useListS3Destinations();
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
    return <div className="p-6">Loading configuration...</div>;
  }

  const list = s3Destinations?.data || [];

  return (
    <div className="space-y-6">
      <PageHeader
        title="Storage destinations"
        description="Manage compatible object storage used for offsite backups."
        action={
          <CreateS3DestinationDialog
            isOpen={isOpen}
            setIsOpen={setIsOpen}
            trigger={
              <Button className="gap-2">
                <Plus className="h-4 w-4" />
                Add destination
              </Button>
            }
          />
        }
      />

      {list.length === 0 ? (
        <div className="flex flex-col items-center justify-center rounded-xl border border-border border-dashed bg-card py-12 text-center">
          <Database className="mb-3 h-10 w-10 text-muted-foreground/40" />
          <h3 className="font-medium text-base">No S3 Destinations Configured</h3>
          <p className="mt-1 max-w-sm text-muted-foreground text-sm">
            Add your AWS S3, Cloudflare R2, or compatible object storage to enable offsite backups.
          </p>
        </div>
      ) : (
        <div className="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
          {list.map((dest) => (
            <div
              key={dest.id}
              className="relative flex flex-col justify-between rounded-xl border border-border bg-card p-5 transition-colors hover:border-primary/30"
            >
              <div>
                <div className="flex items-start justify-between">
                  <div>
                    <h3 className="font-semibold text-base">{dest.name}</h3>
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

                <div className="mt-4 space-y-2 font-mono text-xs">
                  <div className="flex justify-between border-border/40 border-b pb-1">
                    <span className="text-muted-foreground">Provider:</span>
                    <span className="font-medium uppercase">{dest.provider}</span>
                  </div>
                  <div className="flex justify-between border-border/40 border-b pb-1">
                    <span className="text-muted-foreground">Bucket:</span>
                    <span className="truncate pl-2 text-foreground">{dest.bucket}</span>
                  </div>
                  <div className="flex justify-between border-border/40 border-b pb-1">
                    <span className="text-muted-foreground">Region:</span>
                    <span className="text-foreground">{dest.region || 'auto'}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">Endpoint:</span>
                    <span className="max-w-50 truncate pl-2 text-foreground" title={dest.endpoint}>
                      {dest.endpoint}
                    </span>
                  </div>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
