import { useMutation } from '@tanstack/react-query';
import { useState } from 'react';
import { Button } from '@codedock/ui/button';
import { Input } from '@codedock/ui/input';
import { Label } from '@codedock/ui/label';
import { useCreate as useCreateServiceVar } from '#/hooks/use-services';
import type { OperationalCanvas } from './topology-types';

const ENGINE_EXPORTED_KEYS: Record<string, string[]> = {
  postgres: [
    'DATABASE_URL',
    'POSTGRES_URL',
    'POSTGRES_HOST',
    'POSTGRES_PORT',
    'POSTGRES_USER',
    'POSTGRES_PASSWORD',
    'POSTGRES_DB',
  ],
  postgresql: [
    'DATABASE_URL',
    'POSTGRES_URL',
    'POSTGRES_HOST',
    'POSTGRES_PORT',
    'POSTGRES_USER',
    'POSTGRES_PASSWORD',
    'POSTGRES_DB',
  ],
  mysql: [
    'DATABASE_URL',
    'MYSQL_URL',
    'MYSQL_HOST',
    'MYSQL_PORT',
    'MYSQL_USER',
    'MYSQL_PASSWORD',
    'MYSQL_DATABASE',
  ],
  mariadb: [
    'DATABASE_URL',
    'MYSQL_URL',
    'MYSQL_HOST',
    'MYSQL_PORT',
    'MYSQL_USER',
    'MYSQL_PASSWORD',
    'MYSQL_DATABASE',
  ],
  redis: ['REDIS_URL', 'REDIS_HOST', 'REDIS_PORT', 'REDIS_PASSWORD'],
  dragonfly: ['REDIS_URL', 'REDIS_HOST', 'REDIS_PORT', 'REDIS_PASSWORD'],
  keydb: ['REDIS_URL', 'REDIS_HOST', 'REDIS_PORT', 'REDIS_PASSWORD'],
  mongodb: [
    'DATABASE_URL',
    'MONGO_URL',
    'MONGO_HOST',
    'MONGO_PORT',
    'MONGO_USER',
    'MONGO_PASSWORD',
    'MONGO_DB',
  ],
  clickhouse: [
    'CLICKHOUSE_URL',
    'CLICKHOUSE_HOST',
    'CLICKHOUSE_PORT',
    'CLICKHOUSE_USER',
    'CLICKHOUSE_PASSWORD',
  ],
  timescaledb: [
    'DATABASE_URL',
    'POSTGRES_URL',
    'POSTGRES_HOST',
    'POSTGRES_PORT',
    'POSTGRES_USER',
    'POSTGRES_PASSWORD',
    'POSTGRES_DB',
  ],
};

export function CanvasBindingForm({ canvas }: { canvas: OperationalCanvas }) {
  const createServiceVarMutation = useCreateServiceVar();
  const [sourceId, setSource] = useState('');
  const [target, setTarget] = useState('');
  const [key, setKey] = useState('DATABASE_URL');
  const [sourceKey, setSourceKey] = useState('DATABASE_URL');
  const [review, setReview] = useState(false);
  const [reviewedRevision, setReviewedRevision] = useState('');
  const source = canvas.databases.find((database) => database.id === sourceId);
  const engineKeys = source
    ? (ENGINE_EXPORTED_KEYS[source.engine] ?? ['DATABASE_URL'])
    : ['DATABASE_URL'];
  const value = source ? `\u0024{${source.id}.${sourceKey}}` : '';

  const handleSourceChange = (newSourceId: string) => {
    setSource(newSourceId);
    const selectedSource = canvas.databases.find((d) => d.id === newSourceId);
    const availableKeys = selectedSource
      ? (ENGINE_EXPORTED_KEYS[selectedSource.engine] ?? ['DATABASE_URL'])
      : ['DATABASE_URL'];
    const defaultKey = availableKeys[0] ?? 'DATABASE_URL';
    setSourceKey(defaultKey);
  };

  const save = useMutation({
    mutationFn: async () => {
      if (reviewedRevision !== canvas.revision)
        throw new Error('The environment changed. Review this binding again.');
      if (
        !source ||
        !target ||
        !/^[A-Za-z_][A-Za-z0-9_]*$/.test(key) ||
        !/^[A-Za-z_][A-Za-z0-9_]*$/.test(sourceKey)
      )
        throw new Error('Select a database and application and enter valid variable keys.');
      if (source && !engineKeys.includes(sourceKey)) {
        throw new Error(`Select a valid database variable key for ${source.engine}.`);
      }
      await createServiceVarMutation.mutateAsync({
        serviceId: target,
        payload: {
          key,
          value,
          isSecret: true,
          expectedTopologyRevision: reviewedRevision,
        },
      });
    },
    onSuccess: () => {
      setReview(false);
    },
  });
  return (
    <details className="border-b p-3 text-sm">
      <summary className="cursor-pointer">Bind a database variable to an application</summary>
      <fieldset
        disabled={save.isPending}
        onChange={() => setReview(false)}
        className="mt-3 grid gap-3 sm:grid-cols-4"
      >
        <div className="space-y-1">
          <Label htmlFor="binding-source">Database</Label>
          <select
            id="binding-source"
            className="w-full rounded border bg-background p-2"
            value={sourceId}
            onChange={(event) => handleSourceChange(event.target.value)}
          >
            <option value="">Select</option>
            {canvas.databases.map((database) => (
              <option key={database.id} value={database.id}>
                {database.name} ({database.engine})
              </option>
            ))}
          </select>
        </div>
        <div className="space-y-1">
          <Label htmlFor="binding-target">Application</Label>
          <select
            id="binding-target"
            className="w-full rounded border bg-background p-2"
            value={target}
            onChange={(event) => setTarget(event.target.value)}
          >
            <option value="">Select</option>
            {canvas.apps.map((app) => (
              <option key={app.id} value={app.id}>
                {app.name}
              </option>
            ))}
          </select>
        </div>
        <div className="space-y-1">
          <Label htmlFor="binding-source-key">Database key</Label>
          <select
            id="binding-source-key"
            className="w-full rounded border bg-background p-2"
            value={sourceKey}
            onChange={(event) => setSourceKey(event.target.value)}
          >
            {engineKeys.map((k) => (
              <option key={k} value={k}>
                {k}
              </option>
            ))}
          </select>
        </div>
        <div className="space-y-1">
          <Label htmlFor="binding-key">Application key</Label>
          <Input id="binding-key" value={key} onChange={(event) => setKey(event.target.value)} />
        </div>
      </fieldset>
      <div className="mt-3 flex items-center gap-3">
        {review ? (
          <>
            <code>
              {key}={value}
            </code>
            <span>
              This replaces the application's variable and takes effect on its next deployment.
            </span>
            <Button size="sm" disabled={save.isPending} onClick={() => save.mutate()}>
              Apply binding
            </Button>
          </>
        ) : (
          <Button
            size="sm"
            variant="outline"
            disabled={!source || !target}
            onClick={() => {
              setReviewedRevision(canvas.revision);
              setReview(true);
            }}
          >
            Review binding
          </Button>
        )}
      </div>
      {save.isError && (
        <p role="alert" className="mt-2 text-destructive">
          {save.error.message}
        </p>
      )}
    </details>
  );
}
