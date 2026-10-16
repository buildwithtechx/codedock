import { Cloud, Loader2 } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import { Switch } from '#/components/ui/switch';
import { ingressService, type IngressConnectInput } from '#/services/ingress';

interface IngressConnectFormProps {
  serverId: string;
  onBusyChange?: (busy: boolean) => void;
  onClose: () => void;
  onSaved: () => void;
}

export function IngressConnectForm({
  serverId,
  onBusyChange,
  onClose,
  onSaved,
}: IngressConnectFormProps) {
  const [name, setName] = useState('');
  const [credentialId, setCredentialId] = useState('');
  const [accountId, setAccountId] = useState('');
  const [connector, setConnector] = useState<'managed' | 'existing'>('managed');
  const [tunnelId, setTunnelId] = useState('');
  const [defaultForDomains, setDefaultForDomains] = useState(true);
  const [isSubmitting, setIsSubmitting] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim() || !credentialId.trim() || !accountId.trim()) {
      toast.error('Please fill in tunnel name, credential ID, and account ID');
      return;
    }

    setIsSubmitting(true);
    onBusyChange?.(true);
    try {
      const payload: IngressConnectInput = {
        serverId,
        credentialId: credentialId.trim(),
        accountId: accountId.trim(),
        name: name.trim(),
        connector,
        tunnelId: connector === 'existing' && tunnelId.trim() ? tunnelId.trim() : undefined,
        defaultForDomains,
      };

      await ingressService.connect(payload);
      toast.success('Cloudflare Tunnel ingress connection created');
      onSaved();
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'Failed to connect ingress tunnel');
    } finally {
      setIsSubmitting(false);
      onBusyChange?.(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4 p-5" data-testid="ingress-connect-form">
      <div className="flex items-center gap-3 border-border/50 border-b pb-4">
        <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-primary/10">
          <Cloud className="size-5 text-primary" />
        </div>
        <div>
          <h3 className="font-semibold text-base text-foreground">Connect Cloudflare Tunnel</h3>
          <p className="text-muted-foreground text-xs">Route traffic to this host through Cloudflare Tunnel</p>
        </div>
      </div>

      <div className="space-y-2">
        <Label htmlFor="tunnel-name">Tunnel Name</Label>
        <Input
          id="tunnel-name"
          placeholder="e.g. edge-tunnel-primary"
          value={name}
          onChange={(e) => setName(e.target.value)}
          disabled={isSubmitting}
          required
        />
      </div>

      <div className="grid grid-cols-2 gap-3">
        <div className="space-y-2">
          <Label htmlFor="credential-id">Credential ID</Label>
          <Input
            id="credential-id"
            placeholder="cf-cred-..."
            value={credentialId}
            onChange={(e) => setCredentialId(e.target.value)}
            disabled={isSubmitting}
            required
          />
        </div>
        <div className="space-y-2">
          <Label htmlFor="account-id">Account ID</Label>
          <Input
            id="account-id"
            placeholder="Cloudflare Account ID"
            value={accountId}
            onChange={(e) => setAccountId(e.target.value)}
            disabled={isSubmitting}
            required
          />
        </div>
      </div>

      <div className="space-y-2">
        <Label>Connector Mode</Label>
        <div className="grid grid-cols-2 gap-2">
          <Button
            type="button"
            variant={connector === 'managed' ? 'default' : 'outline'}
            size="sm"
            onClick={() => setConnector('managed')}
            disabled={isSubmitting}
          >
            Managed daemon
          </Button>
          <Button
            type="button"
            variant={connector === 'existing' ? 'default' : 'outline'}
            size="sm"
            onClick={() => setConnector('existing')}
            disabled={isSubmitting}
          >
            Existing tunnel
          </Button>
        </div>
      </div>

      {connector === 'existing' && (
        <div className="space-y-2">
          <Label htmlFor="existing-tunnel-id">Existing Tunnel ID</Label>
          <Input
            id="existing-tunnel-id"
            placeholder="UUID of existing tunnel"
            value={tunnelId}
            onChange={(e) => setTunnelId(e.target.value)}
            disabled={isSubmitting}
          />
        </div>
      )}

      <div className="flex items-center justify-between rounded-xl bg-muted/40 p-3">
        <div>
          <p className="font-medium text-foreground text-sm">Default for domains</p>
          <p className="text-muted-foreground text-xs">Automatically assign new domains to this tunnel</p>
        </div>
        <Switch
          checked={defaultForDomains}
          onCheckedChange={setDefaultForDomains}
          disabled={isSubmitting}
        />
      </div>

      <div className="flex items-center justify-end gap-2 pt-2">
        <Button type="button" variant="ghost" onClick={onClose} disabled={isSubmitting}>
          Cancel
        </Button>
        <Button type="submit" disabled={isSubmitting}>
          {isSubmitting && <Loader2 className="mr-1.5 size-4 animate-spin" />}
          Connect Tunnel
        </Button>
      </div>
    </form>
  );
}
