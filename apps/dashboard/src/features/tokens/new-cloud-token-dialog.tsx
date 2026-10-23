import { Cloud } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Button } from '@codedock/ui/button';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@codedock/ui/dialog';
import { Input } from '@codedock/ui/input';
import { Label } from '@codedock/ui/label';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@codedock/ui/select';
import { useCreateManagedCredential } from '#/hooks/use-managed';

interface NewCloudTokenDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  organizationId: string;
}

const CLOUD_PROVIDERS = [
  { id: 'hetzner', name: 'Hetzner Cloud', placeholder: 'Hetzner Cloud API token' },
  { id: 'aws', name: 'Amazon Web Services (AWS)', placeholder: 'AWS Access Key / Secret token' },
  { id: 'digitalocean', name: 'DigitalOcean', placeholder: 'DigitalOcean Personal Access Token' },
  { id: 'scaleway', name: 'Scaleway', placeholder: 'Scaleway API secret key' },
];

export function NewCloudTokenDialog({
  open,
  onOpenChange,
  organizationId,
}: NewCloudTokenDialogProps) {
  const [provider, setProvider] = useState('hetzner');
  const [label, setLabel] = useState('');
  const [token, setToken] = useState('');
  const createMutation = useCreateManagedCredential();

  const activeProviderMeta = CLOUD_PROVIDERS.find((p) => p.id === provider) || CLOUD_PROVIDERS[0];

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!label.trim()) {
      toast.error('Token label is required');
      return;
    }
    if (!token.trim()) {
      toast.error('API token is required');
      return;
    }

    createMutation.mutate(
      {
        organizationId,
        payload: {
          provider,
          label: label.trim(),
          token: token.trim(),
        },
      },
      {
        onSuccess: () => {
          toast.success(`Cloud token "${label}" added successfully`);
          setLabel('');
          setToken('');
          onOpenChange(false);
        },
        onError: (err) => {
          toast.error(err instanceof Error ? err.message : 'Failed to save cloud token');
        },
      }
    );
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <Cloud className="size-4 text-primary" />
            New Cloud Token
          </DialogTitle>
          <DialogDescription>
            Add an API token for your cloud provider to automatically provision and manage servers.
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={handleSubmit} className="space-y-4 py-2">
          <div className="space-y-1.5">
            <Label htmlFor="cloud-provider">Provider</Label>
            <Select value={provider} onValueChange={(val) => setProvider(val)}>
              <SelectTrigger id="cloud-provider">
                <SelectValue placeholder="Select cloud provider" />
              </SelectTrigger>
              <SelectContent>
                {CLOUD_PROVIDERS.map((p) => (
                  <SelectItem key={p.id} value={p.id}>
                    {p.name}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>

          <div className="space-y-1.5">
            <Label htmlFor="cloud-label">Label / Name</Label>
            <Input
              id="cloud-label"
              value={label}
              onChange={(e) => setLabel(e.target.value)}
              placeholder="e.g. Production Hetzner Account"
              required
            />
          </div>

          <div className="space-y-1.5">
            <Label htmlFor="cloud-api-token">API Token / Secret Key</Label>
            <Input
              id="cloud-api-token"
              type="password"
              value={token}
              onChange={(e) => setToken(e.target.value)}
              placeholder={activeProviderMeta.placeholder}
              required
            />
            <p className="text-[11px] text-muted-foreground">
              Encrypted at rest using AES-256 GCM in the backend database.
            </p>
          </div>

          <DialogFooter className="pt-2">
            <Button type="button" variant="outline" size="sm" onClick={() => onOpenChange(false)}>
              Cancel
            </Button>
            <Button type="submit" size="sm" disabled={createMutation.isPending}>
              {createMutation.isPending ? 'Adding...' : 'Add Cloud Token'}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
