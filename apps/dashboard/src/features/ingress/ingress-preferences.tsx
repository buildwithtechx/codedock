import { Loader2 } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Button } from '@codedock/ui/button';
import { Input } from '@codedock/ui/input';
import { Label } from '@codedock/ui/label';
import { Switch } from '@codedock/ui/switch';
import type { IngressConnection } from '#/interfaces/ingress';
import { ingressService } from '#/services/ingress';

interface IngressPreferencesProps {
  connection: IngressConnection;
  onClose: () => void;
  onSaved: () => void;
}

export function IngressPreferences({ connection, onClose, onSaved }: IngressPreferencesProps) {
  const [defaultForDomains, setDefaultForDomains] = useState(connection.defaultForDomains);
  const [policyId, setPolicyId] = useState(connection.defaultAccessPolicyId ?? '');
  const [isSubmitting, setIsSubmitting] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSubmitting(true);
    try {
      await ingressService.update({
        id: connection.id,
        defaultForDomains,
        defaultAccessPolicyId: policyId.trim() || null,
      });
      toast.success('Tunnel preferences saved');
      onSaved();
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'Failed to update preferences');
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4 p-5" data-testid="ingress-preferences-form">
      <div className="border-border/50 border-b pb-3">
        <h3 className="font-semibold text-base text-foreground">Tunnel Preferences</h3>
        <p className="text-muted-foreground text-xs">{connection.name}</p>
      </div>

      <div className="flex items-center justify-between rounded-xl bg-muted/40 p-3">
        <div>
          <p className="font-medium text-foreground text-sm">Default for domains</p>
          <p className="text-muted-foreground text-xs">Route newly added domains via this tunnel</p>
        </div>
        <Switch
          checked={defaultForDomains}
          onCheckedChange={setDefaultForDomains}
          disabled={isSubmitting}
        />
      </div>

      <div className="space-y-2">
        <Label htmlFor="policy-id">Default Access Policy ID (Optional)</Label>
        <Input
          id="policy-id"
          placeholder="e.g. pol_bypass_team"
          value={policyId}
          onChange={(e) => setPolicyId(e.target.value)}
          disabled={isSubmitting}
        />
      </div>

      <div className="flex items-center justify-end gap-2 pt-2">
        <Button type="button" variant="ghost" onClick={onClose} disabled={isSubmitting}>
          Cancel
        </Button>
        <Button type="submit" disabled={isSubmitting}>
          {isSubmitting && <Loader2 className="mr-1.5 size-4 animate-spin" />}
          Save preferences
        </Button>
      </div>
    </form>
  );
}
