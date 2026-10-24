import { Network } from 'lucide-react';
import { useEffect, useId, useState } from 'react';
import { Button } from '@codedock/ui/button';
import { Checkbox } from '@codedock/ui/checkbox';
import { Textarea } from '@codedock/ui/textarea';
import type { ServerNetworkSettings } from '#/interfaces/server';
import { ConnectionValue } from './connection-value';
import {
  useServerNetworkSettings,
  useUpdateServerNetworkSettings,
} from '#/hooks/use-managed-server';
import { ManagedControlPanel } from './managed-control-panel';

export function ManagedServerNetwork({ serverId }: { serverId: string }) {
  const checkboxId = useId();
  const egressId = useId();

  const {
    data: serverSettings,
    isLoading: loading,
    error: queryError,
    refetch,
  } = useServerNetworkSettings(serverId);
  const updateSettingsMutation = useUpdateServerNetworkSettings();

  const [settings, setSettings] = useState<ServerNetworkSettings | null>(null);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [notice, setNotice] = useState<string | null>(null);

  const [internet, setInternet] = useState(true);
  const [rules, setRules] = useState('*');

  useEffect(() => {
    if (serverSettings) {
      setSettings(serverSettings);
      setInternet(serverSettings.internetAccess ?? true);
      setRules(serverSettings.egress?.length ? serverSettings.egress.join('\n') : '*');
    }
  }, [serverSettings]);

  useEffect(() => {
    if (queryError) {
      setError(
        queryError instanceof Error ? queryError.message : 'Failed to load network settings'
      );
    }
  }, [queryError]);

  const fetchSettings = async () => {
    setError(null);
    await refetch();
  };

  const egress = [
    ...new Set(
      rules
        .split(/\r?\n/)
        .map((line) => line.trim().toLowerCase())
        .filter(Boolean)
    ),
  ].sort();

  const rulesChanged =
    internet && JSON.stringify(egress) !== JSON.stringify([...(settings?.egress ?? ['*'])].sort());

  const hasChanges = settings && (internet !== (settings.internetAccess ?? true) || rulesChanged);

  const handleSave = async () => {
    if (!settings) return;
    setSaving(true);
    setError(null);
    setNotice(null);
    try {
      const updated = await updateSettingsMutation.mutateAsync({
        serverId,
        payload: {
          internetAccess: internet,
          expectedInternetAccess: settings.internetAccess ?? true,
          egress,
          expectedRevision: settings.revision,
          confirm: true,
        },
      });
      setSettings(updated);
      setNotice('Network configuration saved successfully.');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to save network settings');
    } finally {
      setSaving(false);
    }
  };

  return (
    <ManagedControlPanel
      title="Networking"
      description="Manage outbound connectivity, egress allowlists, and internal IP topology."
      icon={Network}
      busy={loading || saving}
      error={error}
      refresh={() => void fetchSettings()}
    >
      <div className="flex items-start gap-3 rounded-xl bg-muted/30 p-4">
        <Checkbox
          id={checkboxId}
          checked={internet}
          disabled={loading || saving}
          onCheckedChange={(checked) => setInternet(Boolean(checked))}
        />
        <div>
          <label
            htmlFor={checkboxId}
            className="cursor-pointer font-medium text-foreground text-sm"
          >
            Outbound internet access
          </label>
          <p className="mt-1 text-muted-foreground text-xs leading-relaxed">
            Allow containerized workloads on this server to reach external internet endpoints and
            public APIs.
          </p>
        </div>
      </div>

      <div className="space-y-2">
        <label htmlFor={egressId} className="font-medium text-foreground text-sm">
          Egress allowlist
        </label>
        <p className="text-muted-foreground text-xs">
          Enter allowed domain names or wildcard '*' (one entry per line, max 50 entries).
        </p>
        <Textarea
          id={egressId}
          dir="ltr"
          rows={4}
          value={rules}
          onChange={(e) => setRules(e.target.value)}
          disabled={loading || saving || !internet}
          className="bg-muted/40 font-mono text-xs"
          placeholder="*"
          spellCheck={false}
        />
      </div>

      {hasChanges && (
        <div className="flex flex-wrap items-center justify-between gap-3 rounded-xl bg-warning/10 p-4">
          <p className="text-foreground text-sm">
            Modifying egress rules applies immediately to running workloads.
          </p>
          <Button size="sm" disabled={saving} onClick={() => void handleSave()}>
            {saving ? 'Saving...' : 'Save changes'}
          </Button>
        </div>
      )}

      {notice && (
        <p role="status" className="font-medium text-sm text-success">
          {notice}
        </p>
      )}

      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
        <ConnectionValue label="Private VPC IP" value={settings?.privateIp ?? '10.0.0.4'} />
        <ConnectionValue
          label="Outbound Gateway IP"
          value={settings?.outboundIp ?? '198.51.100.10'}
        />
      </div>

      <div className="rounded-xl bg-muted/20 p-4">
        <h3 className="font-medium text-foreground text-sm">Outbound Mode</h3>
        <p className="mt-1 text-muted-foreground text-xs">
          {settings?.outboundMode === 'managed'
            ? 'Managed isolated tenant network with dynamic NAT gateway.'
            : 'Standard bridged container networking.'}
        </p>
      </div>
    </ManagedControlPanel>
  );
}
