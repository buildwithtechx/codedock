import { Button } from '#/components/ui/button';
import { Input } from '#/components/ui/input';
import type { RuntimeVolume } from './runtime-types';

export function RuntimeVolumeFields({
  volumes,
  onChange,
}: {
  volumes: RuntimeVolume[];
  onChange: (volumes: RuntimeVolume[]) => void;
}) {
  const update = (index: number, changes: Partial<RuntimeVolume>) =>
    onChange(volumes.map((volume, i) => (i === index ? { ...volume, ...changes } : volume)));
  return (
    <fieldset className="space-y-3">
      <legend className="font-medium text-sm">Persistent cluster storage</legend>
      {volumes.map((volume, index) => (
        <div
          key={`${index}-${volumes.length}`}
          className="grid grid-cols-2 gap-2 rounded border p-3"
        >
          <Input
            aria-label="Volume name"
            value={volume.name}
            onChange={(event) => update(index, { name: event.target.value })}
          />
          <Input
            aria-label="Mount path"
            value={volume.mountPath}
            onChange={(event) => update(index, { mountPath: event.target.value })}
          />
          <Input
            aria-label="Storage GiB"
            type="number"
            min={1}
            value={volume.sizeGiB}
            onChange={(event) => update(index, { sizeGiB: Number(event.target.value) })}
          />
          <Input
            aria-label="Storage class"
            placeholder="Default storage class"
            value={volume.storageClass}
            onChange={(event) => update(index, { storageClass: event.target.value })}
          />
          <label className="flex items-center gap-2 text-sm">
            <input
              type="checkbox"
              checked={volume.shared}
              onChange={(event) => update(index, { shared: event.target.checked })}
            />
            Shared ReadWriteMany
          </label>
          <Button variant="outline" onClick={() => onChange(volumes.filter((_, i) => i !== index))}>
            Remove
          </Button>
        </div>
      ))}
      <Button
        variant="outline"
        onClick={() =>
          onChange([
            ...volumes,
            {
              name: `data-${volumes.length + 1}`,
              mountPath: '/data',
              sizeGiB: 10,
              storageClass: '',
              shared: false,
            },
          ])
        }
      >
        Add persistent volume
      </Button>
    </fieldset>
  );
}
