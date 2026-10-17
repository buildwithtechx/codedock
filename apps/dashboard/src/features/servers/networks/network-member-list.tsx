import { Trash2 } from 'lucide-react';
import { Button } from '#/components/ui/button';
import { Input } from '#/components/ui/input';
import type { NetworkMember } from '#/services/networks';

interface NetworkMemberListProps {
  members: NetworkMember[];
  onUpdateIp: (index: number, ip: string) => void;
  onRemove: (index: number) => void;
}

export function NetworkMemberList({ members, onUpdateIp, onRemove }: NetworkMemberListProps) {
  return (
    <div className="divide-y divide-border/50 rounded-xl border border-border/60 bg-muted/20">
      {members.map((member, idx) => (
        <div
          key={member.serverId}
          className="flex flex-wrap items-center justify-between gap-3 p-3"
        >
          <div className="flex items-center gap-3">
            <span className="flex h-6 w-6 items-center justify-center rounded-md bg-muted font-mono text-muted-foreground text-xs">
              {idx + 1}
            </span>
            <div>
              <p className="font-medium text-foreground text-xs">
                {member.serverName || member.serverId}
              </p>
              <p className="font-mono text-[11px] text-muted-foreground">{member.privateIp}</p>
            </div>
          </div>

          <div className="flex items-center gap-2">
            <Input
              className="h-8 w-32 font-mono text-xs"
              value={member.privateIp}
              onChange={(e) => onUpdateIp(idx, e.target.value)}
              placeholder="Private IP"
            />
            <Button
              type="button"
              variant="ghost"
              size="icon"
              className="h-8 w-8 text-destructive hover:bg-destructive/10"
              onClick={() => onRemove(idx)}
              disabled={members.length <= 2}
            >
              <Trash2 className="size-3.5" />
            </Button>
          </div>
        </div>
      ))}
    </div>
  );
}
