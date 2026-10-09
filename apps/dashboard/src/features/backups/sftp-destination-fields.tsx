import { Server } from 'lucide-react';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import { Textarea } from '#/components/ui/textarea';

export interface SFTPFieldValues {
  host: string;
  port: string;
  username: string;
  pathPrefix: string;
  password: string;
  privateKey: string;
}

export function SFTPDestinationFields({
  values,
  onChange,
}: {
  values: SFTPFieldValues;
  onChange: (patch: Partial<SFTPFieldValues>) => void;
}) {
  return (
    <div className="grid gap-4 sm:grid-cols-2">
      <div className="space-y-1.5">
        <Label htmlFor="sftp-host">Host</Label>
        <Input
          id="sftp-host"
          value={values.host}
          onChange={(event) => onChange({ host: event.target.value })}
          placeholder="backups.example.com"
          required
        />
      </div>
      <div className="space-y-1.5">
        <Label htmlFor="sftp-port">Port</Label>
        <Input
          id="sftp-port"
          inputMode="numeric"
          value={values.port}
          onChange={(event) => onChange({ port: event.target.value })}
          placeholder="22"
        />
      </div>
      <div className="space-y-1.5">
        <Label htmlFor="sftp-user">Username</Label>
        <Input
          id="sftp-user"
          value={values.username}
          onChange={(event) => onChange({ username: event.target.value })}
          placeholder="backup"
          required
        />
      </div>
      <div className="space-y-1.5">
        <Label htmlFor="sftp-prefix">Remote path</Label>
        <Input
          id="sftp-prefix"
          value={values.pathPrefix}
          onChange={(event) => onChange({ pathPrefix: event.target.value })}
          placeholder="/backups/codedock"
          className="font-mono"
        />
      </div>
      <div className="space-y-1.5">
        <Label htmlFor="sftp-password">Password</Label>
        <Input
          id="sftp-password"
          type="password"
          value={values.password}
          onChange={(event) => onChange({ password: event.target.value })}
          placeholder="Optional when using a key"
          className="font-mono"
        />
      </div>
      <div className="space-y-1.5 sm:col-span-2">
        <Label htmlFor="sftp-key">Private key</Label>
        <Textarea
          id="sftp-key"
          value={values.privateKey}
          onChange={(event) => onChange({ privateKey: event.target.value })}
          placeholder="-----BEGIN OPENSSH PRIVATE KEY-----"
          rows={4}
          className="font-mono"
        />
      </div>
      <p className="flex items-start gap-2 text-muted-foreground text-xs sm:col-span-2">
        <Server className="mt-0.5 h-3.5 w-3.5 shrink-0" />
        Provide a password, a private key, or both. The connection can be verified from the
        destination list after saving.
      </p>
    </div>
  );
}
