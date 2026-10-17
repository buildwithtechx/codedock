import {
  buildCreateServerPayload,
  readSshPrivateKey,
  validateSshConnection,
} from '@codedock/setup';
import { AlertCircle, CheckCircle2, ChevronDown, ChevronUp } from 'lucide-react';
import { useRef, useState } from 'react';
import { toast } from 'sonner';
import { Button } from '@codedock/ui/button';
import { Input } from '@codedock/ui/input';
import { Label } from '@codedock/ui/label';
import { ServerSshAdvancedFields } from '#/features/servers/server-ssh-advanced-fields';
import { ServerSshAuthFields } from '#/features/servers/server-ssh-auth-fields';
import { useCreateServer, useTestSSH } from '#/hooks/use-servers';
import type { Server, TestSSHRequest } from '#/interfaces/server';

interface ServerSshFormProps {
  onSuccess: (server: Server, credentials: TestSSHRequest) => void;
  onCancel: () => void;
}

export function ServerSshForm({ onSuccess, onCancel }: ServerSshFormProps) {
  const [name, setName] = useState('');
  const [sshHost, setSshHost] = useState('');
  const [sshPort, setSshPort] = useState('22');
  const [sshUser, setSshUser] = useState('root');
  const [sshAuthMethod, setSshAuthMethod] = useState<'key' | 'password'>('key');
  const [sshPrivateKey, setSshPrivateKey] = useState('');
  const [sshPassword, setSshPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [sshJumpHost, setSshJumpHost] = useState('');
  const [sshTransport, setSshTransport] = useState<'direct' | 'cloudflare'>('direct');
  const [showAdvanced, setShowAdvanced] = useState(false);
  const [testResult, setTestResult] = useState<{ ok: boolean; message: string } | null>(null);

  const fileInputRef = useRef<HTMLInputElement>(null);

  const { mutateAsync: testSSH, isPending: isTesting } = useTestSSH();
  const { mutateAsync: createServer, isPending: isCreating } = useCreateServer();

  const handleFileUpload = async (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    event.target.value = '';
    if (!file) return;

    try {
      setSshPrivateKey(await readSshPrivateKey(file));
      toast.success(`Imported ${file.name}`);
    } catch (error) {
      toast.error(error instanceof Error ? error.message : 'Failed to read private key file');
    }
  };

  const handleTestConnection = async () => {
    setTestResult(null);

    if (!sshHost.trim()) {
      toast.error('Host IP address is required');
      return;
    }

    const validationError = validateSshConnection(
      sshAuthMethod,
      sshPrivateKey,
      sshPassword,
      sshPort
    );

    if (validationError) {
      toast.error(validationError);
      return;
    }

    try {
      const result = await testSSH({
        sshHost: sshHost.trim(),
        sshPort: Number(sshPort),
        sshUser: sshUser.trim() || 'root',
        sshKey: sshAuthMethod === 'key' ? sshPrivateKey : undefined,
        sshPassword: sshAuthMethod === 'password' ? sshPassword : undefined,
      });

      setTestResult({ ok: result.success, message: result.message });

      if (result.success) {
        toast.success('SSH preflight check passed');
      } else {
        toast.error('SSH connection failed');
      }
    } catch (err) {
      const msg = err instanceof Error ? err.message : 'SSH connection failed';
      setTestResult({ ok: false, message: msg });
      toast.error(msg);
    }
  };

  const handleSubmit = async (event: React.FormEvent) => {
    event.preventDefault();

    if (!name.trim()) {
      toast.error('Server name is required');
      return;
    }

    if (!sshHost.trim()) {
      toast.error('Host address is required');
      return;
    }

    const validationError = validateSshConnection(
      sshAuthMethod,
      sshPrivateKey,
      sshPassword,
      sshPort
    );

    if (validationError) {
      toast.error(validationError);
      return;
    }

    try {
      const serverPayload = buildCreateServerPayload(name, {
        host: sshHost,
        port: sshPort,
        user: sshUser,
        authMethod: sshAuthMethod,
        privateKey: sshPrivateKey,
        password: sshPassword,
        jumpHost: sshJumpHost,
        transport: sshTransport,
      });

      const credentials: TestSSHRequest = {
        sshHost: serverPayload.ipAddress,
        sshPort: serverPayload.sshPort,
        sshUser: serverPayload.sshUser,
        sshKey: serverPayload.sshKey,
        sshPassword: serverPayload.sshPassword,
      };

      const server = await createServer({
        ...serverPayload,
        sshHost: serverPayload.ipAddress,
        sshPrivateKey: serverPayload.sshKey,
        sshTransport,
        isLocal: false,
      });

      toast.success('Server added successfully');
      onSuccess(server, credentials);
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'Failed to add server');
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      <div className="grid gap-4 sm:grid-cols-2">
        <div className="space-y-1.5 sm:col-span-2">
          <Label htmlFor="srv-name" className="text-xs">
            Server Name
          </Label>
          <Input
            id="srv-name"
            value={name}
            onChange={(e) => setName(e.target.value)}
            placeholder="eu-central-prod-node"
            required
          />
        </div>

        <div className="space-y-1.5">
          <Label htmlFor="srv-host" className="text-xs">
            Host / IP Address
          </Label>
          <Input
            id="srv-host"
            value={sshHost}
            onChange={(e) => setSshHost(e.target.value)}
            placeholder="198.51.100.1 or node.example.com"
            required
          />
        </div>

        <div className="grid grid-cols-2 gap-2">
          <div className="space-y-1.5">
            <Label htmlFor="srv-port" className="text-xs">
              SSH Port
            </Label>
            <Input
              id="srv-port"
              type="number"
              min={1}
              max={65535}
              value={sshPort}
              onChange={(e) => setSshPort(e.target.value)}
              placeholder="22"
            />
          </div>
          <div className="space-y-1.5">
            <Label htmlFor="srv-user" className="text-xs">
              SSH User
            </Label>
            <Input
              id="srv-user"
              value={sshUser}
              onChange={(e) => setSshUser(e.target.value)}
              placeholder="root"
            />
          </div>
        </div>
      </div>

      <ServerSshAuthFields
        sshAuthMethod={sshAuthMethod}
        setSshAuthMethod={setSshAuthMethod}
        sshPrivateKey={sshPrivateKey}
        setSshPrivateKey={setSshPrivateKey}
        sshPassword={sshPassword}
        setSshPassword={setSshPassword}
        showPassword={showPassword}
        setShowPassword={setShowPassword}
        fileInputRef={fileInputRef}
        handleFileUpload={handleFileUpload}
      />

      <div className="space-y-2">
        <button
          type="button"
          onClick={() => setShowAdvanced(!showAdvanced)}
          className="inline-flex items-center gap-1.5 font-medium text-muted-foreground text-xs hover:text-foreground"
        >
          {showAdvanced ? (
            <ChevronUp className="h-3.5 w-3.5" />
          ) : (
            <ChevronDown className="h-3.5 w-3.5" />
          )}
          Advanced routing (Bastion / Jump Host, Transport)
        </button>

        {showAdvanced && (
          <ServerSshAdvancedFields
            sshJumpHost={sshJumpHost}
            setSshJumpHost={setSshJumpHost}
            sshTransport={sshTransport}
            setSshTransport={setSshTransport}
          />
        )}
      </div>

      {testResult && (
        <div
          className={`flex items-start gap-2.5 rounded-xl border p-3 text-xs ${
            testResult.ok
              ? 'border-emerald-500/30 bg-emerald-500/10 text-emerald-600 dark:text-emerald-400'
              : 'border-destructive/30 bg-destructive/10 text-destructive'
          }`}
        >
          {testResult.ok ? (
            <CheckCircle2 className="mt-0.5 h-4 w-4 shrink-0" />
          ) : (
            <AlertCircle className="mt-0.5 h-4 w-4 shrink-0" />
          )}
          <span className="leading-relaxed">{testResult.message}</span>
        </div>
      )}

      <div className="flex items-center justify-between border-border/60 border-t pt-4">
        <Button
          type="button"
          variant="outline"
          onClick={handleTestConnection}
          disabled={isTesting || !sshHost.trim()}
          className="text-xs"
        >
          {isTesting ? 'Verifying...' : 'Test Connection'}
        </Button>
        <div className="flex items-center gap-2">
          <Button type="button" variant="ghost" onClick={onCancel} className="text-xs">
            Cancel
          </Button>
          <Button type="submit" disabled={isCreating} className="text-xs">
            {isCreating ? 'Adding Server...' : 'Add Server'}
          </Button>
        </div>
      </div>
    </form>
  );
}
