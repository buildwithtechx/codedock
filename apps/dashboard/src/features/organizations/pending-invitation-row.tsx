import { Mail, Trash2 } from 'lucide-react';
import { toast } from 'sonner';
import { Button } from '@codedock/ui/button';
import { useRemoveOrganizationMember } from '#/hooks/use-organizations';
import type { OrganizationMember } from '#/interfaces/organizations';

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
