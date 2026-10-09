import { CheckCircle2, Loader2, XCircle } from 'lucide-react';
import { useEffect, useMemo, useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import { cn } from '#/lib/utils';
import {
  useCreateS3Destination,
  useCreateSFTPDestination,
  useUpdateS3Destination,
  useVerifyS3Draft,
} from './hooks';
import type { DestinationKind, S3Destination } from './interfaces';
import type { S3FieldValues } from './s3-destination-fields';
import { S3DestinationFields } from './s3-destination-fields';
import type { S3ProviderId } from './s3-providers';
import { buildS3Endpoint, deriveS3Provider } from './s3-providers';
import type { SFTPFieldValues } from './sftp-destination-fields';
import { SFTPDestinationFields } from './sftp-destination-fields';

type TestState = { status: 'idle' | 'testing' | 'ok' | 'fail'; reason?: string };

export function DestinationConfigureForm({
  kind,
  editing,
  onCancel,
  onSaved,
}: {
  kind: DestinationKind;
  editing: S3Destination | null;
  onCancel: () => void;
  onSaved: () => void;
}) {
  const initial = useMemo(
    () =>
      kind === 's3' && editing
        ? deriveS3Provider(editing.endpoint, editing.region)
        : { provider: 'r2' as S3ProviderId, region: '', accountId: '', endpoint: '' },
    [kind, editing]
  );
  const [name, setName] = useState(editing?.name ?? '');
  const [description, setDescription] = useState(editing?.description ?? '');
  const [s3, setS3] = useState<S3FieldValues>({
    provider: initial.provider,
    region: initial.region,
    accountId: initial.accountId,
    endpoint: initial.endpoint,
    bucket: editing?.bucket ?? '',
    pathPrefix: editing?.pathPrefix ?? '',
    isDefault: editing?.isDefault ?? false,
    accessKeyId: '',
    secretAccessKey: '',
  });
  const [sftp, setSftp] = useState<SFTPFieldValues>({
    host: '',
    port: '22',
    username: '',
    pathPrefix: '',
    password: '',
    privateKey: '',
  });
  const [test, setTest] = useState<TestState>({ status: 'idle' });

  const createS3 = useCreateS3Destination();
  const updateS3 = useUpdateS3Destination();
  const createSftp = useCreateSFTPDestination();
  const verifyDraft = useVerifyS3Draft();
  const busy = createS3.isPending || updateS3.isPending || createSftp.isPending;

  useEffect(() => {
    setTest({ status: 'idle' });
  }, [s3, kind]);

  const runTest = async () => {
    if (
      !s3.bucket.trim() ||
      (!s3.accessKeyId.trim() && !editing) ||
      (!s3.secretAccessKey.trim() && !editing)
    ) {
      toast.error('Enter bucket and credentials before testing');
      return;
    }
    setTest({ status: 'testing' });
    try {
      const built = buildS3Endpoint(s3.provider, s3);
      const res = await verifyDraft.mutateAsync({
        provider: s3.provider,
        endpoint: built.endpoint,
        bucket: s3.bucket.trim(),
        region: built.region,
        accessKeyId: s3.accessKeyId.trim() || editing?.accessKeyId || '',
        secretAccessKey: s3.secretAccessKey,
      });
      if (res.data?.ok) setTest({ status: 'ok' });
      else setTest({ status: 'fail', reason: res.data?.reason || 'Could not connect' });
    } catch (error) {
      setTest({
        status: 'fail',
        reason: error instanceof Error ? error.message : 'Verification failed',
      });
    }
  };

  const submit = async (event: React.FormEvent) => {
    event.preventDefault();
    try {
      if (kind === 's3') {
        const built = buildS3Endpoint(s3.provider, s3);
        const payload = {
          projectId: 'global',
          name: name.trim(),
          description: description.trim(),
          provider: s3.provider,
          endpoint: built.endpoint,
          bucket: s3.bucket.trim(),
          region: built.region,
          pathPrefix: s3.pathPrefix.trim(),
          isDefault: s3.isDefault,
          accessKeyId: s3.accessKeyId.trim() || editing?.accessKeyId || '',
          secretAccessKey: s3.secretAccessKey,
        };
        if (editing) {
          await updateS3.mutateAsync({ id: editing.id, payload });
          toast.success(`Destination "${name.trim()}" updated`);
        } else {
          await createS3.mutateAsync({ payload });
          toast.success(`Destination "${name.trim()}" created`);
        }
      } else {
        await createSftp.mutateAsync({
          name: name.trim(),
          description: description.trim(),
          host: sftp.host.trim(),
          port: Number.parseInt(sftp.port, 10) || 22,
          username: sftp.username.trim(),
          password: sftp.password || undefined,
          privateKey: sftp.privateKey.trim() || undefined,
          pathPrefix: sftp.pathPrefix.trim() || undefined,
        });
        toast.success(`Destination "${name.trim()}" created`);
      }
      onSaved();
    } catch (error) {
      toast.error(error instanceof Error ? error.message : 'Failed to save destination');
    }
  };

  return (
    <form onSubmit={submit} className="space-y-5">
      <div className="space-y-1.5">
        <Label htmlFor="destination-name">Name</Label>
        <Input
          id="destination-name"
          value={name}
          onChange={(event) => setName(event.target.value)}
          placeholder="Production backups"
          required
        />
      </div>
      <div className="space-y-1.5">
        <Label htmlFor="destination-description">Description</Label>
        <Input
          id="destination-description"
          value={description}
          onChange={(event) => setDescription(event.target.value)}
          placeholder="Off-server snapshot storage"
        />
      </div>

      {kind === 's3' ? (
        <>
          <S3DestinationFields
            values={s3}
            editing={!!editing}
            onChange={(patch) => setS3((previous) => ({ ...previous, ...patch }))}
          />
          {test.status !== 'idle' && (
            <div
              role="status"
              className={cn(
                'flex items-start gap-2.5 rounded-xl border px-4 py-3 text-sm',
                test.status === 'ok' && 'border-emerald-500/30 bg-emerald-500/10 text-emerald-600',
                test.status === 'fail' &&
                  'border-destructive/30 bg-destructive/10 text-destructive',
                test.status === 'testing' && 'border-border/50 bg-muted/20 text-muted-foreground'
              )}
            >
              {test.status === 'testing' && <Loader2 className="mt-0.5 h-4 w-4 animate-spin" />}
              {test.status === 'ok' && <CheckCircle2 className="mt-0.5 h-4 w-4" />}
              {test.status === 'fail' && <XCircle className="mt-0.5 h-4 w-4" />}
              <span>
                {test.status === 'testing'
                  ? 'Probing the bucket…'
                  : test.status === 'ok'
                    ? 'Connection verified — the bucket is reachable.'
                    : (test.reason ?? 'Verification failed')}
              </span>
            </div>
          )}
        </>
      ) : (
        <SFTPDestinationFields
          values={sftp}
          onChange={(patch) => setSftp((previous) => ({ ...previous, ...patch }))}
        />
      )}

      <div className="flex flex-wrap items-center justify-between gap-3 border-border/40 border-t pt-5">
        {kind === 's3' ? (
          <Button
            type="button"
            variant="outline"
            onClick={runTest}
            disabled={busy || test.status === 'testing' || !name.trim()}
          >
            {test.status === 'testing' ? (
              <Loader2 className="h-4 w-4 animate-spin" />
            ) : (
              <CheckCircle2 className="h-4 w-4" />
            )}
            Test connection
          </Button>
        ) : (
          <span />
        )}
        <div className="ms-auto flex items-center gap-2">
          <Button type="button" variant="ghost" onClick={onCancel} disabled={busy}>
            Cancel
          </Button>
          <Button type="submit" disabled={busy || !name.trim()}>
            {busy ? <Loader2 className="h-4 w-4 animate-spin" /> : null}
            {editing ? 'Save changes' : 'Save destination'}
          </Button>
        </div>
      </div>
    </form>
  );
}
