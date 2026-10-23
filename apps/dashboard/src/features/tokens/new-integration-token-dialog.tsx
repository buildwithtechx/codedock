import { Blocks } from 'lucide-react';
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

export interface IntegrationToken {
  id: string;
  name: string;
  provider: string;
  maskedToken: string;
  createdAt: string;
}

interface NewIntegrationTokenDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onAddToken: (token: IntegrationToken) => void;
}

const INTEGRATION_PROVIDERS = [
  { id: 'docker', name: 'Docker Hub / Container Registry', placeholder: 'dckr_pat_...' },
  { id: 'github', name: 'GitHub Personal Access Token', placeholder: 'ghp_...' },
  { id: 'gitlab', name: 'GitLab Personal Access Token', placeholder: 'glpat-...' },
  { id: 'sentry', name: 'Sentry Auth Token', placeholder: 'sntrys_...' },
  {
    id: 'webhook',
    name: 'Slack / Discord Webhook URL',
    placeholder: 'https://hooks.slack.com/services/...',
  },
];

export function NewIntegrationTokenDialog({
  open,
  onOpenChange,
  onAddToken,
}: NewIntegrationTokenDialogProps) {
  const [provider, setProvider] = useState('docker');
  const [name, setName] = useState('');
  const [tokenValue, setTokenValue] = useState('');

  const activeProvider =
    INTEGRATION_PROVIDERS.find((p) => p.id === provider) || INTEGRATION_PROVIDERS[0];

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim()) {
      toast.error('Token name is required');
      return;
    }
    if (!tokenValue.trim()) {
      toast.error('Token or credential value is required');
      return;
    }

    const masked =
      tokenValue.length > 8 ? `${tokenValue.slice(0, 4)}••••${tokenValue.slice(-4)}` : '••••••••';

    const newItem: IntegrationToken = {
      id: `int-${Date.now()}`,
      name: name.trim(),
      provider: activeProvider.name,
      maskedToken: masked,
      createdAt: new Date().toISOString(),
    };

    onAddToken(newItem);
    toast.success(`Integration token "${name}" saved`);
    setName('');
    setTokenValue('');
    onOpenChange(false);
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <Blocks className="size-4 text-primary" />
            New Integration Token
          </DialogTitle>
          <DialogDescription>
            Credentials used by external integrations such as Docker registries, Sentry, and git
            providers.
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={handleSubmit} className="space-y-4 py-2">
          <div className="space-y-1.5">
            <Label htmlFor="int-provider">Integration Type</Label>
            <Select value={provider} onValueChange={setProvider}>
              <SelectTrigger id="int-provider">
                <SelectValue placeholder="Select type" />
              </SelectTrigger>
              <SelectContent>
                {INTEGRATION_PROVIDERS.map((p) => (
                  <SelectItem key={p.id} value={p.id}>
                    {p.name}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>

          <div className="space-y-1.5">
            <Label htmlFor="int-name">Name</Label>
            <Input
              id="int-name"
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="e.g. Production Docker Hub Registry"
              required
            />
          </div>

          <div className="space-y-1.5">
            <Label htmlFor="int-val">Secret / Access Token</Label>
            <Input
              id="int-val"
              type="password"
              value={tokenValue}
              onChange={(e) => setTokenValue(e.target.value)}
              placeholder={activeProvider.placeholder}
              required
            />
          </div>

          <DialogFooter className="pt-2">
            <Button type="button" variant="outline" size="sm" onClick={() => onOpenChange(false)}>
              Cancel
            </Button>
            <Button type="submit" size="sm">
              Save Integration Token
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
