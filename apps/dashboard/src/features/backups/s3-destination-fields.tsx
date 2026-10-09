import { useState } from 'react';
import { Checkbox } from '#/components/ui/checkbox';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import { cn } from '#/lib/utils';
import type { S3Provider, S3ProviderId } from './s3-providers';
import { providerById, s3Providers } from './s3-providers';

export function S3ProviderMark({
  provider,
  className,
}: {
  provider: S3Provider;
  className?: string;
}) {
  const [failed, setFailed] = useState(false);
  if (!failed) {
    return (
      <img
        src={provider.logo}
        alt=""
        aria-hidden="true"
        onError={() => setFailed(true)}
        className={cn('shrink-0 object-contain', className)}
      />
    );
  }
  const FallbackIcon = provider.icon;
  return <FallbackIcon className={cn('shrink-0 text-muted-foreground', className)} />;
}

export interface S3FieldValues {
  provider: S3ProviderId;
  region: string;
  accountId: string;
  endpoint: string;
  bucket: string;
  pathPrefix: string;
  isDefault: boolean;
  accessKeyId: string;
  secretAccessKey: string;
}

export function S3DestinationFields({
  values,
  editing,
  onChange,
}: {
  values: S3FieldValues;
  editing: boolean;
  onChange: (patch: Partial<S3FieldValues>) => void;
}) {
  const activeProvider = providerById(values.provider);
  return (
    <>
      <div>
        <span className="mb-2 block font-medium text-sm">Provider</span>
        <div className="grid grid-cols-2 gap-2 sm:grid-cols-3">
          {s3Providers.map((option) => {
            const selected = option.id === values.provider;
            return (
              <button
                key={option.id}
                type="button"
                onClick={() => onChange({ provider: option.id })}
                className={cn(
                  'flex items-center gap-2 rounded-xl border px-3 py-2.5 text-left transition-all',
                  selected
                    ? 'border-primary/50 bg-primary/[0.06] shadow-sm'
                    : 'border-border/50 hover:border-border hover:bg-muted/30'
                )}
              >
                <S3ProviderMark provider={option} className="size-5" />
                <span className="truncate font-medium text-[13px]">{option.label}</span>
              </button>
            );
          })}
        </div>
      </div>
      <div className="grid gap-4 sm:grid-cols-2">
        {activeProvider.mode === 'accountId' && (
          <div className="space-y-1.5">
            <Label htmlFor="s3-account">Account ID</Label>
            <Input
              id="s3-account"
              value={values.accountId}
              onChange={(event) => onChange({ accountId: event.target.value })}
              placeholder="a1b2c3d4e5f6"
              className="font-mono"
            />
          </div>
        )}
        {activeProvider.mode === 'endpoint' && (
          <div className="space-y-1.5 sm:col-span-2">
            <Label htmlFor="s3-endpoint">Endpoint URL</Label>
            <Input
              id="s3-endpoint"
              value={values.endpoint}
              onChange={(event) => onChange({ endpoint: event.target.value })}
              placeholder="https://minio.example.com:9000"
              className="font-mono"
            />
          </div>
        )}
        {activeProvider.mode !== 'accountId' && (
          <div className="space-y-1.5">
            <Label htmlFor="s3-region">Region</Label>
            <Input
              id="s3-region"
              value={values.region}
              onChange={(event) => onChange({ region: event.target.value })}
              placeholder={activeProvider.regionPlaceholder}
              className="font-mono"
            />
          </div>
        )}
        <div className="space-y-1.5">
          <Label htmlFor="s3-bucket">Bucket</Label>
          <Input
            id="s3-bucket"
            value={values.bucket}
            onChange={(event) => onChange({ bucket: event.target.value })}
            placeholder="my-backups"
            required
          />
        </div>
        <div className="space-y-1.5">
          <Label htmlFor="s3-prefix">Path prefix</Label>
          <Input
            id="s3-prefix"
            value={values.pathPrefix}
            onChange={(event) => onChange({ pathPrefix: event.target.value })}
            placeholder="codedock/prod"
            className="font-mono"
          />
        </div>
        <div className="space-y-1.5">
          <Label htmlFor="s3-key">Access key ID</Label>
          <Input
            id="s3-key"
            value={values.accessKeyId}
            onChange={(event) => onChange({ accessKeyId: event.target.value })}
            placeholder={editing ? '(stored — leave blank to keep)' : 'AKIA…'}
            className="font-mono"
          />
        </div>
        <div className="space-y-1.5 sm:col-span-2">
          <Label htmlFor="s3-secret">Secret access key</Label>
          <Input
            id="s3-secret"
            type="password"
            value={values.secretAccessKey}
            onChange={(event) => onChange({ secretAccessKey: event.target.value })}
            placeholder={editing ? '(stored — leave blank to keep)' : '••••••••'}
            className="font-mono"
          />
        </div>
      </div>
      <label
        htmlFor="s3-default"
        className="flex cursor-pointer items-center gap-2 text-muted-foreground text-sm"
      >
        <Checkbox
          id="s3-default"
          checked={values.isDefault}
          onCheckedChange={(value) => onChange({ isDefault: value === true })}
        />
        Use as the default destination for new policies
      </label>
    </>
  );
}
