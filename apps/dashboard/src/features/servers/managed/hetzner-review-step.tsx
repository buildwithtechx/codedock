import { Loader2 } from 'lucide-react';
import { Button } from '#/components/ui/button';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '#/components/ui/select';
import type { ManagedQuota } from '#/features/servers/managed-types';

export function HetznerReviewStep({
  name,
  image,
  images,
  location,
  tierName,
  tierLabel,
  monthlyPrice,
  quota,
  quotaLoading,
  working,
  error,
  onNameChange,
  onImageChange,
  onBack,
  onConfirm,
}: {
  name: string;
  image: string;
  images: string[];
  location: string;
  tierName: string;
  tierLabel: string;
  monthlyPrice: number;
  quota: ManagedQuota | null;
  quotaLoading: boolean;
  working: boolean;
  error: string | null;
  onNameChange: (name: string) => void;
  onImageChange: (image: string) => void;
  onBack: () => void;
  onConfirm: () => void;
}) {
  const quotaWouldExceed =
    quota !== null && (quota.usedServers + 1 > quota.maxServers || quota.maxServers === 0);

  return (
    <div className="space-y-5">
      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
        <div className="space-y-1.5">
          <Label htmlFor="hetzner-server-name" className="text-xs">
            Server name
          </Label>
          <Input
            id="hetzner-server-name"
            value={name}
            onChange={(e) => onNameChange(e.target.value)}
            placeholder="edge-1"
            disabled={working}
            maxLength={63}
          />
        </div>
        <div className="space-y-1.5">
          <Label htmlFor="hetzner-server-image" className="text-xs">
            Image
          </Label>
          <Select value={image} onValueChange={onImageChange} disabled={working}>
            <SelectTrigger id="hetzner-server-image" className="w-full">
              <SelectValue placeholder="Select image" />
            </SelectTrigger>
            <SelectContent>
              {images.map((entry) => (
                <SelectItem key={entry} value={entry}>
                  {entry}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
      </div>

      <div className="rounded-xl border border-border/70">
        <div className="border-border/60 border-b px-4 py-3">
          <h3 className="font-medium text-foreground text-sm">Provision plan</h3>
        </div>
        <dl className="grid grid-cols-2 gap-x-4 gap-y-3 px-4 py-3 text-xs sm:grid-cols-4">
          <div>
            <dt className="text-muted-foreground">Server type</dt>
            <dd className="mt-0.5 font-medium font-mono text-foreground">{tierName}</dd>
          </div>
          <div>
            <dt className="text-muted-foreground">Shape</dt>
            <dd className="mt-0.5 font-medium text-foreground">{tierLabel}</dd>
          </div>
          <div>
            <dt className="text-muted-foreground">Location</dt>
            <dd className="mt-0.5 font-medium font-mono text-foreground">{location}</dd>
          </div>
          <div>
            <dt className="text-muted-foreground">List price</dt>
            <dd className="mt-0.5 font-medium text-foreground">€{monthlyPrice.toFixed(2)}/mo</dd>
          </div>
        </dl>
      </div>

      <div className="rounded-xl border border-border/70 px-4 py-3 text-xs">
        <p className="text-muted-foreground">
          SSH access uses a dedicated keypair Codedock generates per server and registers with
          Hetzner as the server's deploy key. Password logins stay disabled.
        </p>
      </div>

      <div className="rounded-xl border border-border/70 px-4 py-3 text-xs">
        {quotaLoading ? (
          <p className="flex items-center gap-2 text-muted-foreground">
            <Loader2 className="size-3.5 animate-spin" /> Checking organization quota…
          </p>
        ) : quota ? (
          <p className="text-muted-foreground">
            Quota usage after this server:{' '}
            <span className="font-medium text-foreground">
              {quota.usedServers + 1} of {quota.maxServers} servers
            </span>{' '}
            · {quota.usedMemoryGB} of {quota.maxMemoryGB} GB RAM in use. Hetzner bills the stored
            credential directly.
          </p>
        ) : (
          <p className="text-muted-foreground">Quota information is unavailable.</p>
        )}
        {quotaWouldExceed && (
          <p className="mt-1.5 text-destructive">
            This provision exceeds the server quota. Ask an organization admin to raise it first.
          </p>
        )}
      </div>

      {error && <p className="text-destructive text-xs">{error}</p>}

      <div className="flex items-center justify-between border-border/60 border-t pt-4">
        <Button
          type="button"
          variant="ghost"
          onClick={onBack}
          disabled={working}
          className="text-xs"
        >
          Back
        </Button>
        <Button
          type="button"
          onClick={onConfirm}
          disabled={working || !name.trim() || !image || quotaWouldExceed}
          className="text-xs"
        >
          {working && <Loader2 className="size-3.5 animate-spin" />}
          Review & provision
        </Button>
      </div>
    </div>
  );
}
