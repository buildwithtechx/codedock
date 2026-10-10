import { Loader2, Trash2 } from 'lucide-react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '#/components/ui/select';
import { useRemoveOrganizationMember, useUpdateOrganizationMember } from './hooks';
import type { OrganizationMember } from './interfaces';

export function MemberRow({
  member,
  organizationId,
  isCurrentUserOwner,
  canManageTeam,
  ownerCount,
  currentUserId,
  currentUserEmail,
}: {
  member: OrganizationMember;
  organizationId: string;
  isCurrentUserOwner: boolean;
  canManageTeam: boolean;
  ownerCount: number;
  currentUserId?: string;
  currentUserEmail?: string;
}) {
  const { mutateAsync: updateMember, isPending: isUpdating } =
    useUpdateOrganizationMember(organizationId);
  const { mutateAsync: removeMember, isPending: isRemoving } =
    useRemoveOrganizationMember(organizationId);

  const isCurrentUser =
    (Boolean(member.userId) && member.userId === currentUserId) ||
    member.email === currentUserEmail;

  const handlePermissionChange = async (newPermission: string) => {
    try {
      await updateMember({ memberId: member.userId!, payload: { permission: newPermission } });
      toast.success('Member role updated');
    } catch (err: unknown) {
      const message = err instanceof Error ? err.message : 'Failed to update member role';
      toast.error(message);
    }
  };

  const handleRemove = async () => {
    if (!window.confirm(`Are you sure you want to remove ${member.email}?`)) return;
    try {
      await removeMember(member.id);
      toast.success('Member removed');
    } catch (err: unknown) {
      const message = err instanceof Error ? err.message : 'Failed to remove member';
      toast.error(message);
    }
  };

  const isPending = isUpdating || isRemoving;
  const canEditRole =
    isCurrentUserOwner && member.permission !== 'owner' && Boolean(member.userId) && !isCurrentUser;
  const initial = member.email ? member.email.slice(0, 1).toUpperCase() : 'U';

  return (
    <div className="flex items-center gap-4 px-5 py-4 transition-colors hover:bg-muted/30">
      <div className="flex size-9 shrink-0 items-center justify-center rounded-full bg-muted font-medium text-foreground text-sm">
        {initial}
      </div>

      <div className="min-w-0 flex-1">
        <div className="flex items-center gap-2">
          <p className="truncate font-medium text-foreground text-sm">{member.email}</p>
          {isCurrentUser && (
            <span className="rounded-md bg-muted px-1.5 py-0.5 font-normal text-muted-foreground text-xs">
              You
            </span>
          )}
        </div>
        <div className="flex items-center gap-2 text-muted-foreground text-xs">
          <span>
            {member.invitedAt ? new Date(member.invitedAt).toLocaleDateString() : 'Active'}
          </span>
          {member.status !== 'active' && member.status !== 'accepted' && (
            <span className="rounded-full bg-amber-500/10 px-2 py-0.5 font-medium text-amber-600 dark:text-amber-400">
              {member.status}
            </span>
          )}
        </div>
      </div>

      <div className="flex items-center gap-3">
        {canEditRole ? (
          <Select
            defaultValue={member.permission}
            onValueChange={handlePermissionChange}
            disabled={isPending}
          >
            <SelectTrigger className="h-8 w-28 text-xs">
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="owner">Owner</SelectItem>
              <SelectItem value="admin">Admin</SelectItem>
              <SelectItem value="member">Member</SelectItem>
            </SelectContent>
          </Select>
        ) : (
          <span className="rounded-md border border-border/50 bg-muted/40 px-2.5 py-1 font-medium text-muted-foreground text-xs uppercase tracking-wider">
            {member.permission}
          </span>
        )}

        <Button
          variant="ghost"
          size="icon"
          onClick={handleRemove}
          disabled={
            isPending ||
            !canManageTeam ||
            (member.permission === 'owner' && (ownerCount <= 1 || !isCurrentUserOwner))
          }
          className="size-8 text-muted-foreground hover:bg-destructive/10 hover:text-destructive"
          aria-label={`Remove ${member.email}`}
        >
          {isRemoving ? <Loader2 className="size-4 animate-spin" /> : <Trash2 className="size-4" />}
        </Button>
      </div>
    </div>
  );
}
