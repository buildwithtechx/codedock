import { ArrowLeft, Loader2, Plus, Trash2 } from 'lucide-react';
import type React from 'react';
import { Button } from '#/components/ui/button';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';

interface ServiceConfigFormProps {
  name: string;
  onNameChange: (v: string) => void;
  image: string;
  onImageChange: (v: string) => void;
  port: number;
  onPortChange: (v: number) => void;
  volume: string;
  onVolumeChange: (v: string) => void;
  envRows: { key: string; value: string }[];
  onAddEnvRow: () => void;
  onUpdateEnvRow: (index: number, field: 'key' | 'value', val: string) => void;
  onRemoveEnvRow: (index: number) => void;
  onBack: () => void;
  onCancel: () => void;
  onSubmit: (e: React.FormEvent) => void;
  isSubmitting: boolean;
}

export function ServiceConfigForm({
  name,
  onNameChange,
  image,
  onImageChange,
  port,
  onPortChange,
  volume,
  onVolumeChange,
  envRows,
  onAddEnvRow,
  onUpdateEnvRow,
  onRemoveEnvRow,
  onBack,
  onCancel,
  onSubmit,
  isSubmitting,
}: ServiceConfigFormProps) {
  return (
    <form onSubmit={onSubmit} className="space-y-4">
      <button
        type="button"
        onClick={onBack}
        className="inline-flex items-center gap-1.5 text-muted-foreground text-xs hover:text-foreground"
      >
        <ArrowLeft className="size-3.5" />
        Back to catalog
      </button>

      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
        <div className="space-y-1.5">
          <Label htmlFor="svc-name" className="text-xs">
            Service Name
          </Label>
          <Input
            id="svc-name"
            value={name}
            onChange={(e) => onNameChange(e.target.value)}
            placeholder="e.g. postgres"
            required
          />
        </div>

        <div className="space-y-1.5">
          <Label htmlFor="svc-port" className="text-xs">
            Port
          </Label>
          <Input
            id="svc-port"
            type="number"
            value={port}
            onChange={(e) => onPortChange(Number(e.target.value))}
            placeholder="5432"
            required
          />
        </div>
      </div>

      <div className="space-y-1.5">
        <Label htmlFor="svc-image" className="text-xs">
          Docker Image
        </Label>
        <Input
          id="svc-image"
          value={image}
          onChange={(e) => onImageChange(e.target.value)}
          placeholder="e.g. postgres:16-alpine"
          required
        />
      </div>

      {volume && (
        <div className="space-y-1.5">
          <Label htmlFor="svc-vol" className="text-xs">
            Data Volume
          </Label>
          <Input
            id="svc-vol"
            value={volume}
            onChange={(e) => onVolumeChange(e.target.value)}
            placeholder="e.g. pgdata:/var/lib/postgresql/data"
          />
        </div>
      )}

      <div className="space-y-2 border-border/50 border-t pt-2">
        <div className="flex items-center justify-between">
          <Label className="font-semibold text-xs">Environment Variables</Label>
          <Button
            type="button"
            variant="outline"
            size="sm"
            onClick={onAddEnvRow}
            className="h-7 gap-1 text-[11px]"
          >
            <Plus className="size-3" />
            Add Var
          </Button>
        </div>

        {envRows.length === 0 ? (
          <p className="text-muted-foreground text-xs">No environment variables configured.</p>
        ) : (
          <div className="space-y-2">
            {envRows.map((row, idx) => (
              <div key={idx} className="flex items-center gap-2">
                <Input
                  value={row.key}
                  onChange={(e) => onUpdateEnvRow(idx, 'key', e.target.value)}
                  placeholder="KEY"
                  className="font-mono text-xs"
                />
                <Input
                  value={row.value}
                  onChange={(e) => onUpdateEnvRow(idx, 'value', e.target.value)}
                  placeholder="VALUE"
                  className="font-mono text-xs"
                />
                <button
                  type="button"
                  onClick={() => onRemoveEnvRow(idx)}
                  className="text-muted-foreground hover:text-destructive"
                >
                  <Trash2 className="size-3.5" />
                </button>
              </div>
            ))}
          </div>
        )}
      </div>

      <div className="flex justify-end gap-2 border-border/50 border-t pt-4">
        <Button type="button" variant="outline" onClick={onCancel}>
          Cancel
        </Button>
        <Button type="submit" disabled={isSubmitting || !name.trim() || !image.trim()}>
          {isSubmitting ? (
            <>
              <Loader2 className="mr-2 size-4 animate-spin" />
              Adding...
            </>
          ) : (
            'Add Service'
          )}
        </Button>
      </div>
    </form>
  );
}
