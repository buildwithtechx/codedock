import { Key } from 'lucide-react';
import { useEffect, useState } from 'react';
import { Button } from '#/components/ui/button';
import { Textarea } from '#/components/ui/textarea';
import type { ManagedSshStatus } from '#/interfaces/server';
import { ConnectionValue } from './connection-value';
import { getManagedSshStatus, setManagedSsh, setManagedSshKey } from './managed-api';
import { ManagedControlPanel } from './managed-control-panel';

export function ManagedServerSsh({ serverId }: { serverId: string }) {
  const [status, setStatus] = useState<ManagedSshStatus | null>(null);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [notice, setNotice] = useState<string | null>(null);
  const [confirmDisable, setConfirmDisable] = useState(false);
  const [publicKey, setPublicKey] = useState('');

  const fetchStatus = async () => {
    setLoading(true);
    setError(null);
    try {
      const data = await getManagedSshStatus(serverId);
      setStatus(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load SSH status');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    void fetchStatus();
  }, [serverId]);

  const toggleSsh = async (enabled: boolean) => {
    setSaving(true);
    setError(null);
    setNotice(null);
    try {
      const res = await setManagedSsh(serverId, {
        enabled,
        expectedEnabled: status?.enabled ?? !enabled,
        confirm: true,
      });
      setStatus(res.status);
      setConfirmDisable(false);
      setNotice(enabled ? 'SSH access enabled.' : 'SSH access disabled.');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to update SSH status');
    } finally {
      setSaving(false);
    }
  };

  const handleSaveKey = async () => {
    if (!publicKey.trim()) return;
    setSaving(true);
    setError(null);
    setNotice(null);
    try {
      await setManagedSshKey(serverId, publicKey.trim());
      setPublicKey('');
      setNotice('Public SSH key added successfully.');
      void fetchStatus();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to add SSH key');
    } finally {
      setSaving(false);
    }
  };

  const enabled = status?.enabled === true;

  return (
    <ManagedControlPanel
      title="SSH Access"
      description="Configure terminal bridge access, bastion proxy credentials, and authorized keys."
      icon={Key}
      busy={loading || saving}
      error={error}
      refresh={() => void fetchStatus()}
    >
      <div className="flex flex-wrap items-center justify-between gap-3 rounded-xl bg-muted/30 p-4">
        <div>
          <span className="font-medium text-foreground text-sm">
            {enabled ? 'SSH Access is Enabled' : 'SSH Access is Disabled'}
          </span>
          <p className="mt-0.5 text-muted-foreground text-xs">
            {enabled
              ? 'Authorized SSH connections are forwarded through the secure edge bastion.'
              : 'All inbound SSH terminal connections are actively blocked.'}
          </p>
        </div>
        <Button
          variant={enabled ? 'secondary' : 'default'}
          size="sm"
          disabled={loading || saving}
          onClick={() => (enabled ? setConfirmDisable(true) : void toggleSsh(true))}
        >
          {enabled ? 'Disable SSH' : 'Enable SSH'}
        </Button>
      </div>

      {confirmDisable && (
        <div className="space-y-3 rounded-xl bg-warning/10 p-4">
          <p className="text-foreground text-sm">
            Disabling SSH will terminate any open shell sessions and revoke temporary bastion
            tunnels.
          </p>
          <div className="flex flex-wrap gap-2">
            <Button
              variant="destructive"
              size="sm"
              disabled={saving}
              onClick={() => void toggleSsh(false)}
            >
              Confirm Disable
            </Button>
            <Button
              variant="ghost"
              size="sm"
              disabled={saving}
              onClick={() => setConfirmDisable(false)}
            >
              Cancel
            </Button>
          </div>
        </div>
      )}

      {notice && (
        <p role="status" className="font-medium text-sm text-success">
          {notice}
        </p>
      )}

      {enabled && (
        <>
          <div className="space-y-4 rounded-xl bg-muted/20 p-4">
            <h3 className="font-medium text-foreground text-sm">Connection Credentials</h3>
            <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
              <ConnectionValue label="SSH User" value={status?.connection?.user ?? 'root'} />
              <ConnectionValue
                label="Bastion Host"
                value={status?.connection?.bastion ?? 'ssh.codedock.run'}
              />
            </div>
            <ConnectionValue
              label="Connect via SSH Command"
              value={status?.connection?.command ?? 'ssh -J ssh.codedock.run root@198.51.100.10'}
            />
          </div>

          <div className="space-y-3">
            <label htmlFor="ssh-public-key" className="font-medium text-foreground text-sm">
              Add Authorized Public Key
            </label>
            <p className="text-muted-foreground text-xs">
              Paste an OpenSSH public key (ssh-ed25519 or ssh-rsa) to grant interactive access.
            </p>
            <Textarea
              id="ssh-public-key"
              rows={3}
              placeholder="ssh-ed25519 AAAAC3NzaC1lZDI1NTE5..."
              value={publicKey}
              onChange={(e) => setPublicKey(e.target.value)}
              className="font-mono text-xs"
            />
            <Button
              size="sm"
              variant="secondary"
              disabled={saving || !publicKey.trim()}
              onClick={() => void handleSaveKey()}
            >
              Add public key
            </Button>
          </div>
        </>
      )}
    </ManagedControlPanel>
  );
}
