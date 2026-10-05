import {
  AlertCircle,
  CheckCircle2,
  ChevronDown,
  ChevronUp,
  Eye,
  EyeOff,
  KeyRound,
  Lock,
  Upload,
} from 'lucide-react';
import { useRef, useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';
import { ServerSshAdvancedFields } from '#/features/servers/server-ssh-advanced-fields';
import { useCreateServer, useTestSSH } from '#/hooks/use-servers';
import type { Server } from '#/interfaces/server';

interface ServerSshFormProps {
  onSuccess: (server: Server) => void;
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

    if (file.size > 64 * 1024) {
      toast.error('Key file exceeds maximum size of 64KB');
      return;
    }

    try {
      const text = await file.text();
      if (!text.includes('PRIVATE KEY-----')) {
        toast.error('Select an SSH private key file');
        return;
      }
      setSshPrivateKey(text);
      toast.success(`Imported ${file.name}`);
    } catch {
      toast.error('Failed to read private key file');
    }
  };

  const handleTestConnection = async () => {
    if (!sshHost.trim()) {
      toast.error('Host IP address is required');
      return;
    }

    setTestResult(null);
    try {
      const result = await testSSH({
        sshHost: sshHost.trim(),
        sshPort: Number.parseInt(sshPort, 10) || 22,
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

    if (
      (sshAuthMethod === 'key' && !sshPrivateKey.trim()) ||
      (sshAuthMethod === 'password' && !sshPassword)
    ) {
      toast.error('Provide the selected SSH credential');
      return;
    }
    const port = Number(sshPort);
    if (!Number.isInteger(port) || port < 1 || port > 65535) {
      toast.error('SSH port must be between 1 and 65535');
      return;
    }
    try {
      const server = await createServer({
        name: name.trim(),
        ipAddress: sshHost.trim(),
        sshHost: sshHost.trim(),
        sshPort: Number.parseInt(sshPort, 10) || 22,
        sshUser: sshUser.trim() || 'root',
        sshAuthMethod,
        sshPrivateKey: sshAuthMethod === 'key' ? sshPrivateKey : undefined,
        sshPassword: sshAuthMethod === 'password' ? sshPassword : undefined,
        sshTransport,
        sshJumpHost: sshJumpHost.trim() || undefined,
        isLocal: false,
      });
      toast.success('Server added successfully');
      onSuccess(server);
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

      <div className="space-y-3 rounded-xl border border-border/70 bg-muted/20 p-4">
        <div className="flex items-center justify-between">
          <Label className="font-semibold text-xs">Authentication Method</Label>
          <div className="inline-flex rounded-lg border border-border/80 bg-background p-0.5 text-xs">
            <button
              type="button"
              onClick={() => setSshAuthMethod('key')}
              className={`inline-flex items-center gap-1.5 rounded-md px-3 py-1 font-medium transition-colors ${
                sshAuthMethod === 'key'
                  ? 'bg-primary text-primary-foreground'
                  : 'text-muted-foreground hover:text-foreground'
              }`}
            >
              <KeyRound className="h-3 w-3" />
              Private Key
            </button>
            <button
              type="button"
              onClick={() => setSshAuthMethod('password')}
              className={`inline-flex items-center gap-1.5 rounded-md px-3 py-1 font-medium transition-colors ${
                sshAuthMethod === 'password'
                  ? 'bg-primary text-primary-foreground'
                  : 'text-muted-foreground hover:text-foreground'
              }`}
            >
              <Lock className="h-3 w-3" />
              Password
            </button>
          </div>
        </div>

        {sshAuthMethod === 'key' ? (
          <div className="space-y-2">
            <div className="flex items-center justify-between">
              <span className="text-muted-foreground text-xs">
                Paste OpenSSH / PEM key or upload key file
              </span>
              <Button
                type="button"
                variant="outline"
                size="sm"
                className="h-7 gap-1 text-xs"
                onClick={() => fileInputRef.current?.click()}
              >
                <Upload className="h-3 w-3" />
                Upload file
              </Button>
              <input
                ref={fileInputRef}
                type="file"
                className="hidden"
                accept=".pem,.key,id_rsa,id_ed25519"
                onChange={handleFileUpload}
              />
            </div>
            <textarea
              value={sshPrivateKey}
              onChange={(e) => setSshPrivateKey(e.target.value)}
              placeholder="-----BEGIN OPENSSH PRIVATE KEY-----&#10;...&#10;-----END OPENSSH PRIVATE KEY-----"
              rows={4}
              spellCheck={false}
              className="w-full rounded-lg border border-border bg-background p-2.5 font-mono text-xs outline-none focus:ring-2 focus:ring-primary/20"
            />
          </div>
        ) : (
          <div className="space-y-1.5">
            <Label htmlFor="srv-pass" className="text-xs">
              SSH Password
            </Label>
            <div className="relative">
              <Input
                id="srv-pass"
                type={showPassword ? 'text' : 'password'}
                value={sshPassword}
                onChange={(e) => setSshPassword(e.target.value)}
                placeholder="••••••••••••"
                className="pr-9"
              />
              <button
                type="button"
                onClick={() => setShowPassword(!showPassword)}
                className="absolute top-1/2 right-2.5 -translate-y-1/2 text-muted-foreground hover:text-foreground"
              >
                {showPassword ? (
                  <EyeOff className="h-3.5 w-3.5" />
                ) : (
                  <Eye className="h-3.5 w-3.5" />
                )}
              </button>
            </div>
          </div>
        )}
      </div>

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
