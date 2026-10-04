import { Input } from '#/components/ui/input';
import { Label } from '#/components/ui/label';

interface ServerSshAdvancedFieldsProps {
  sshJumpHost: string;
  setSshJumpHost: (val: string) => void;
  sshTransport: 'direct' | 'cloudflare';
  setSshTransport: (val: 'direct' | 'cloudflare') => void;
}

export function ServerSshAdvancedFields({
  sshJumpHost,
  setSshJumpHost,
  sshTransport,
  setSshTransport,
}: ServerSshAdvancedFieldsProps) {
  return (
    <div className="space-y-3 rounded-xl border border-border/60 bg-muted/10 p-4">
      <div className="space-y-1.5">
        <Label htmlFor="srv-jump" className="text-xs">
          SSH Jump Host / Bastion
        </Label>
        <Input
          id="srv-jump"
          value={sshJumpHost}
          onChange={(e) => setSshJumpHost(e.target.value)}
          placeholder="bastion.example.com:22"
        />
      </div>
      <div className="space-y-1.5">
        <Label htmlFor="srv-trans" className="text-xs">
          Transport Mode
        </Label>
        <select
          id="srv-trans"
          value={sshTransport}
          onChange={(e) => setSshTransport(e.target.value as 'direct' | 'cloudflare')}
          className="w-full rounded-lg border border-border bg-background px-3 py-2 text-xs outline-none focus:ring-2 focus:ring-primary/20"
        >
          <option value="direct">Direct TCP</option>
          <option value="cloudflare">Cloudflare Tunnel (cloudflared)</option>
        </select>
      </div>
    </div>
  );
}
