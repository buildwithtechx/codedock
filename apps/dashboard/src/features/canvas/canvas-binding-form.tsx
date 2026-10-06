import { useMutation, useQueryClient } from '@tanstack/react-query';
import { useState } from 'react';
import { Button } from '#/components/ui/button';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import { serviceVariablesService } from '#/services/service-variables';
import type { OperationalCanvas } from './topology-types';

export function CanvasBindingForm({ canvas }: { canvas: OperationalCanvas }) {
  const client = useQueryClient();
  const [sourceId, setSource] = useState('');
  const [target, setTarget] = useState('');
  const [key, setKey] = useState('DATABASE_URL');
  const [sourceKey, setSourceKey] = useState('DATABASE_URL');
  const [review, setReview] = useState(false);
  const [reviewedRevision, setReviewedRevision] = useState('');
  const source = canvas.databases.find((database) => database.id === sourceId);
  const value = source ? `\u0024{${source.id}.${sourceKey}}` : '';
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
      await serviceVariablesService.create(target, {
        key,
        value,
        isSecret: true,
        expectedTopologyRevision: reviewedRevision,
      });
    },
    onSuccess: async () => {
      setReview(false);
      await client.invalidateQueries({ queryKey: ['canvas'] });
      await client.invalidateQueries({ queryKey: ['variables'] });
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
            onChange={(event) => setSource(event.target.value)}
          >
            <option value="">Select</option>
            {canvas.databases.map((database) => (
              <option key={database.id} value={database.id}>
                {database.name}
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
          <Input
            id="binding-source-key"
            value={sourceKey}
            onChange={(event) => setSourceKey(event.target.value)}
          />
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
