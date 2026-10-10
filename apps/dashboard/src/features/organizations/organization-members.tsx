import { Loader2, UserPlus, Users } from 'lucide-react';
import { useState } from 'react';
import { Button } from '#/components/ui/button';
import { useAuthStore } from '#/stores/auth-store';
import { useListOrganizationMembers } from './hooks';
import { InviteMemberModal } from './invite-member-modal';
import { MemberRow } from './member-row';
import { PendingInvitationRow } from './pending-invitation-row';

export function OrganizationMembers({
  organizationId,
  organizationName,
}: {
  organizationId: string;
  organizationName?: string;
}) {
  const { data: members = [], isLoading } = useListOrganizationMembers(organizationId);
  const [inviteOpen, setInviteOpen] = useState(false);
  const user = useAuthStore((s) => s.user);

  const currentUserMember = members.find((m) => m.userId === user?.id || m.email === user?.email);
  const isCurrentUserOwner = currentUserMember?.permission === 'owner';
  const canManageTeam =
    user?.role === 'admin' ||
    currentUserMember?.permission === 'owner' ||
    currentUserMember?.permission === 'admin';

  const activeMembers = members.filter(
    (m) => m.status === 'active' || m.status === 'accepted' || (!m.status && Boolean(m.userId))
  );

  const pendingMembers = members.filter(
    (m) =>
      m.status === 'pending' ||
      m.status === 'invited' ||
      (!m.userId && m.status !== 'active' && m.status !== 'accepted')
  );

  const displayedActive =
    activeMembers.length === 0 && pendingMembers.length === 0 && members.length > 0
      ? members
      : activeMembers;

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
        <div>
          <h1 className="font-bold text-foreground text-xl">{organizationName || 'Team'}</h1>
          <p className="mt-1 text-muted-foreground text-sm">
            Manage who has access to this workspace and its deployed projects.
          </p>
        </div>
        {canManageTeam && (
          <Button
            onClick={() => setInviteOpen(true)}
            className="shrink-0 gap-2 self-start sm:self-auto"
          >
            <UserPlus className="size-4" />
            Invite member
          </Button>
        )}
      </div>

      {isLoading ? (
        <div className="flex items-center justify-center py-16">
          <Loader2 className="size-5 animate-spin text-muted-foreground" />
        </div>
      ) : (
        <>
          <div className="rounded-2xl border border-border/50 bg-card">
            <div className="flex items-center justify-between border-border/50 border-b px-5 py-3.5">
              <h2 className="font-semibold text-foreground text-sm">
                Active members ({displayedActive.length})
              </h2>
            </div>
            {displayedActive.length === 0 ? (
              <div className="flex flex-col items-center justify-center px-4 py-12 text-center">
                <div className="flex size-10 items-center justify-center rounded-xl bg-muted text-muted-foreground">
                  <Users className="size-5" />
                </div>
                <p className="mt-3 font-medium text-foreground text-sm">No active members found</p>
                <p className="mt-1 text-muted-foreground text-xs">
                  Invite teammates to collaborate on this workspace.
                </p>
              </div>
            ) : (
              <div className="divide-y divide-border/40">
                {displayedActive.map((member) => (
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

          {pendingMembers.length > 0 && (
            <div className="rounded-2xl border border-border/50 bg-card">
              <div className="flex items-center justify-between border-border/50 border-b px-5 py-3.5">
                <h2 className="font-semibold text-foreground text-sm">
                  Pending invitations ({pendingMembers.length})
                </h2>
              </div>
              <div className="divide-y divide-border/40">
                {pendingMembers.map((member) => (
                  <PendingInvitationRow
                    key={member.id}
                    member={member}
                    organizationId={organizationId}
                    canManageTeam={canManageTeam}
                  />
                ))}
              </div>
            </div>
          )}
        </>
      )}

      <InviteMemberModal
        organizationId={organizationId}
        open={inviteOpen}
        onOpenChange={setInviteOpen}
        isCurrentUserOwner={isCurrentUserOwner}
      />
    </div>
  );
}
