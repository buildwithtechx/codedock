import { Check, Copy, Mail, Trash2 } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import { useRemoveOrganizationMember } from './hooks';
import type { OrganizationMember } from './interfaces';

export function PendingInvitationRow({
  member,
  organizationId,
  canManageTeam,
}: {
  member: OrganizationMember;
  organizationId: string;
  canManageTeam: boolean;
}) {
  const { mutateAsync: cancelInvite, isPending } = useRemoveOrganizationMember(organizationId);
  const [copied, setCopied] = useState(false);

  const handleCopyLink = async () => {
    try {
      const inviteUrl = `${window.location.origin}/invite/${member.id}`;
      await navigator.clipboard.writeText(inviteUrl);
      setCopied(true);
      toast.success('Invitation link copied to clipboard');
      setTimeout(() => setCopied(false), 2000);
    } catch {
      toast.error('Failed to copy invitation link');
    }
  };

  const handleCancel = async () => {
    if (!window.confirm(`Cancel invitation for ${member.email}?`)) return;
    try {
      await cancelInvite(member.id);
      toast.success('Invitation cancelled');
    } catch (err: unknown) {
      const message = err instanceof Error ? err.message : 'Failed to cancel invitation';
      toast.error(message);
    }
  };

  return (
    <div className="flex items-center gap-4 px-5 py-4 transition-colors hover:bg-muted/20">
      <div className="flex size-9 shrink-0 items-center justify-center rounded-full bg-muted/60 text-muted-foreground">
        <Mail className="size-4" />
      </div>

      <div className="min-w-0 flex-1">
        <p className="truncate font-medium text-foreground text-sm">{member.email}</p>
        <p className="mt-0.5 text-muted-foreground text-xs">
          Invited as {member.permission} ·{' '}
          {member.invitedAt ? new Date(member.invitedAt).toLocaleDateString() : 'Pending'}
        </p>
      </div>

      <div className="flex items-center gap-2">
        <Button
          type="button"
          variant="outline"
          size="sm"
          onClick={handleCopyLink}
          className="gap-1.5 text-xs"
        >
          {copied ? <Check className="size-3.5 text-emerald-500" /> : <Copy className="size-3.5" />}
          {copied ? 'Copied' : 'Copy link'}
        </Button>
        {canManageTeam && (
          <Button
            type="button"
            variant="ghost"
            size="sm"
            onClick={handleCancel}
            disabled={isPending}
            className="text-muted-foreground text-xs hover:bg-destructive/10 hover:text-destructive"
          >
            <Trash2 className="size-3.5" />
            Cancel
          </Button>
        )}
      </div>
    </div>
  );
}
