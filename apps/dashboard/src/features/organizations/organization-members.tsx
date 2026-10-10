import { Loader2, Plus, Users } from 'lucide-react';
import { useState } from 'react';
import { Button } from '#/components/ui/button';
import { useAuthStore } from '#/stores/auth-store';
import { useListOrganizationMembers } from './hooks';
import { InviteMemberModal } from './invite-member-modal';
import { MemberRow } from './member-row';

export function OrganizationMembers({ organizationId }: { organizationId: string }) {
  const { data: members = [], isLoading } = useListOrganizationMembers(organizationId);
  const [inviteOpen, setInviteOpen] = useState(false);
  const user = useAuthStore((s) => s.user);

  const currentUserMember = members.find((m) => m.userId === user?.id || m.email === user?.email);
  const isCurrentUserOwner = currentUserMember?.permission === 'owner';
  const canManageTeam =
    user?.role === 'admin' ||
    currentUserMember?.permission === 'owner' ||
    currentUserMember?.permission === 'admin';

  const ownerCount =
    members.filter(
      (m) =>
        m.permission === 'owner' &&
        Boolean(m.userId) &&
        (m.status === 'active' || m.status === 'accepted')
    ).length || 0;

  return (
    <div className="space-y-6">
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div className="flex items-center gap-3">
          <div className="flex h-9 w-9 items-center justify-center rounded-xl bg-primary/10 text-primary">
            <Users className="h-4 w-4" />
          </div>
          <div>
            <h2 className="font-semibold text-foreground text-sm">Team Members</h2>
            <p className="text-muted-foreground text-xs">
              Manage who has access to this workspace and its deployed projects.
            </p>
          </div>
        </div>
        {canManageTeam && (
          <Button onClick={() => setInviteOpen(true)} size="sm" className="shrink-0 gap-2">
            <Plus className="size-4" />
            Invite member
          </Button>
        )}
      </div>

      <div className="overflow-hidden rounded-2xl border border-border/50 bg-card">
        <div className="flex items-center justify-between border-border/50 border-b px-5 py-3.5">
          <h3 className="font-semibold text-foreground text-sm">
            Active members ({members.length})
          </h3>
        </div>

        {isLoading ? (
          <div className="flex items-center justify-center py-12">
            <Loader2 className="size-5 animate-spin text-muted-foreground" />
          </div>
        ) : members.length === 0 ? (
          <div className="flex flex-col items-center justify-center px-4 py-12 text-center">
            <div className="flex size-10 items-center justify-center rounded-xl bg-muted text-muted-foreground">
              <Users className="size-5" />
            </div>
            <p className="mt-3 font-medium text-foreground text-sm">No members found</p>
            <p className="mt-1 text-muted-foreground text-xs">
              Invite teammates to collaborate on this workspace.
            </p>
          </div>
        ) : (
          <div className="divide-y divide-border/40">
            {members.map((member) => (
              <MemberRow
                key={member.id}
                member={member}
                organizationId={organizationId}
                isCurrentUserOwner={isCurrentUserOwner}
                canManageTeam={canManageTeam}
                ownerCount={ownerCount}
                currentUserId={user?.id}
                currentUserEmail={user?.email}
              />
            ))}
          </div>
        )}
      </div>

      <InviteMemberModal
        organizationId={organizationId}
        open={inviteOpen}
        onOpenChange={setInviteOpen}
        isCurrentUserOwner={isCurrentUserOwner}
      />
    </div>
  );
}
