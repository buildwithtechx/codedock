import {
  AlertTriangle,
  CheckCircle2,
  HardDrive,
  Loader2,
  Plus,
  RefreshCw,
  Trash2,
} from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import { Card, CardContent } from '#/components/ui/card';
import { QueryErrorState } from '#/components/ui/query-error-state';
import {
  useDeleteS3Destination,
  useList,
  useListS3Destinations,
  useVerifyS3Destination,
} from '#/features/backups';

type BackupDestinationsProps = {
  onAddDestination: () => void;
};

export function BackupDestinations({ onAddDestination }: BackupDestinationsProps) {
  const { data: s3Destinations, isLoading, isError, refetch } = useListS3Destinations();
  const { data: backupConfigs } = useList();
  const deleteS3Dest = useDeleteS3Destination();
  const verifyS3Dest = useVerifyS3Destination();

  const [verifyingIds, setVerifyingIds] = useState<Record<string, 'loading' | 'success' | 'error'>>(
    {}
  );

  const list = s3Destinations?.data || [];

  const handleVerify = async (id: string, name: string) => {
    setVerifyingIds((prev) => ({ ...prev, [id]: 'loading' }));
    try {
      const res = await verifyS3Dest.mutateAsync(id);
      if (res.data?.ok) {
        setVerifyingIds((prev) => ({ ...prev, [id]: 'success' }));
        toast.success(`Storage "${name}" verified successfully`);
      } else {
        setVerifyingIds((prev) => ({ ...prev, [id]: 'error' }));
        toast.error(`Verification failed: ${res.data?.reason || 'Could not connect to S3'}`);
      }
    } catch (err) {
      setVerifyingIds((prev) => ({ ...prev, [id]: 'error' }));
      toast.error(err instanceof Error ? err.message : 'Verification failed');
    }
  };

  const handleDelete = async (id: string, name: string) => {
    const isReferenced = backupConfigs?.data?.some((c) => c.s3DestinationId === id);
    if (isReferenced) {
      toast.error(`Cannot delete "${name}": it is currently referenced by a backup policy.`);
      return;
    }

    if (
      !window.confirm(
        `Are you sure you want to delete storage destination "${name}"? This action cannot be undone.`
      )
    ) {
      return;
    }

    try {
      await deleteS3Dest.mutateAsync({ id, projectId: 'global' });
      toast.success('Storage destination deleted');
    } catch (error) {
      toast.error(error instanceof Error ? error.message : 'Failed to delete storage destination');
    }
  };

  const getProviderLogo = (provider: string) => {
    const p = provider.toLowerCase();
    if (p.includes('r2') || p.includes('cloudflare')) {
      return '/dns-providers/cloudflare.svg';
    }
    if (p.includes('aws') || p.includes('s3') || p.includes('amazon')) {
      return '/provider-logos/aws.svg';
    }
    if (p.includes('gcp') || p.includes('google')) {
      return '/provider-logos/gcp.svg';
    }
    if (p.includes('digitalocean') || p.includes('spaces')) {
      return '/provider-logos/digitalocean.svg';
    }
    return null;
  };

  if (isLoading) {
    return (
      <div className="flex min-h-[16rem] items-center justify-center">
        <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
      </div>
    );
  }

  if (isError) {
    return (
      <QueryErrorState
        title="Storage destinations are unavailable"
        description="Could not load configured backup storage destinations."
        onRetry={() => void refetch()}
      />
    );
  }

  return (
    <section className="space-y-4">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="font-semibold text-base text-foreground">Storage destinations</h2>
          <p className="mt-0.5 text-muted-foreground text-xs">
            S3-compatible object storage targets for automated snapshots.
          </p>
        </div>
        <Button size="sm" onClick={onAddDestination} className="gap-1.5">
          <Plus className="h-4 w-4" />
          Add destination
        </Button>
      </div>

      {list.length === 0 ? (
        <Card className="border-border/60 bg-card">
          <CardContent className="flex flex-col items-center justify-center py-12 text-center">
            <div className="flex h-12 w-12 items-center justify-center rounded-2xl bg-muted text-muted-foreground">
              <HardDrive className="h-6 w-6" />
            </div>
            <h3 className="mt-4 font-semibold text-base text-foreground">
              No storage destinations
            </h3>
            <p className="mt-1 max-w-sm text-muted-foreground text-xs leading-5">
              Connect Cloudflare R2, AWS S3, MinIO, or DigitalOcean Spaces to store snapshots
              off-server.
            </p>
            <Button size="sm" className="mt-5 gap-1.5" onClick={onAddDestination}>
              <Plus className="h-4 w-4" />
              Connect first destination
            </Button>
          </CardContent>
        </Card>
      ) : (
        <div className="grid gap-4 sm:grid-cols-2">
          {list.map((dest) => {
            const logo = getProviderLogo(dest.provider);
            const verifyState = verifyingIds[dest.id];

            return (
              <div
                key={dest.id}
                className="flex flex-col justify-between rounded-2xl border border-border/60 bg-card p-5 shadow-sm transition-colors hover:border-border"
              >
                <div>
                  <div className="flex items-start justify-between gap-3">
                    <div className="flex items-center gap-3">
                      <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl border border-border/60 bg-muted/60">
                        {logo ? (
                          <img src={logo} alt={dest.provider} className="h-5 w-5 object-contain" />
                        ) : (
                          <HardDrive className="h-5 w-5 text-muted-foreground" />
                        )}
                      </div>
                      <div>
                        <h3 className="font-semibold text-foreground text-sm">{dest.name}</h3>
                        <p className="font-mono text-muted-foreground text-xs uppercase">
                          {dest.provider}
                        </p>
                      </div>
                    </div>

                    <div className="flex items-center gap-1">
                      {verifyState === 'success' && (
                        <span className="flex items-center gap-1 rounded-full bg-emerald-500/10 px-2 py-0.5 font-medium text-emerald-500 text-xs">
                          <CheckCircle2 className="h-3 w-3" /> Verified
                        </span>
                      )}
                      {verifyState === 'error' && (
                        <span className="flex items-center gap-1 rounded-full bg-rose-500/10 px-2 py-0.5 font-medium text-rose-500 text-xs">
                          <AlertTriangle className="h-3 w-3" /> Failed
                        </span>
                      )}
                    </div>
                  </div>

                  {dest.description && (
                    <p className="mt-3 text-muted-foreground text-xs leading-5">
                      {dest.description}
                    </p>
                  )}

                  <dl className="mt-4 space-y-1.5 border-border/40 border-t pt-3 font-mono text-xs">
                    <div className="flex justify-between">
                      <dt className="text-muted-foreground">Bucket</dt>
                      <dd className="font-medium text-foreground">{dest.bucket}</dd>
                    </div>
                    <div className="flex justify-between">
                      <dt className="text-muted-foreground">Region</dt>
                      <dd className="text-muted-foreground">{dest.region || 'auto'}</dd>
                    </div>
                    <div className="flex justify-between">
                      <dt className="text-muted-foreground">Endpoint</dt>
                      <dd
                        className="max-w-[14rem] truncate text-muted-foreground"
                        title={dest.endpoint}
                      >
                        {dest.endpoint}
                      </dd>
                    </div>
                  </dl>
                </div>

                <div className="mt-5 flex items-center justify-between border-border/40 border-t pt-3">
                  <Button
                    variant="outline"
                    size="sm"
                    className="h-8 text-xs"
                    onClick={() => handleVerify(dest.id, dest.name)}
                    disabled={verifyState === 'loading'}
                  >
                    {verifyState === 'loading' ? (
                      <Loader2 className="mr-1.5 h-3.5 w-3.5 animate-spin" />
                    ) : (
                      <RefreshCw className="mr-1.5 h-3.5 w-3.5" />
                    )}
                    Test connection
                  </Button>

                  <Button
                    variant="ghost"
                    size="sm"
                    className="h-8 text-destructive hover:bg-destructive/10 hover:text-destructive"
                    onClick={() => handleDelete(dest.id, dest.name)}
                  >
                    <Trash2 className="mr-1 h-3.5 w-3.5" />
                    Delete
                  </Button>
                </div>
              </div>
            );
          })}
        </div>
      )}
    </section>
  );
}
