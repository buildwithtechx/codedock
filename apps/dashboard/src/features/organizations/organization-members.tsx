import { Building2, Loader2, UserPlus, Users } from 'lucide-react';
import { useState } from 'react';
import { Button } from '#/components/ui/button';
import { SettingsSection } from '#/features/settings/settings-section';
import { useAuthStore } from '#/stores/auth-store';
import { CreateOrganizationModal } from './create-organization-modal';
import { useListOrganizationMembers } from './hooks';
import { InviteMemberInline } from './invite-member-inline';
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
  const [createOrgOpen, setCreateOrgOpen] = useState(false);
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
    <SettingsSection
      icon={<Users className="size-4 text-primary" />}
      title={organizationName || 'Team'}
      description="Manage who has access to this workspace and its deployed projects."
      action={
        <div className="flex items-center gap-2">
          <Button
            size="sm"
            variant="outline"
            onClick={() => setCreateOrgOpen(true)}
            className="gap-1.5"
          >
            <Building2 className="size-3.5" />
            New organization
          </Button>
          {canManageTeam && !inviteOpen && (
            <Button size="sm" onClick={() => setInviteOpen(true)} className="gap-1.5">
              <UserPlus className="size-3.5" />
              Invite member
            </Button>
          )}
        </div>
      }
    >
      <CreateOrganizationModal open={createOrgOpen} onOpenChange={setCreateOrgOpen} />
      {isLoading ? (
        <div className="flex items-center justify-center py-12">
          <Loader2 className="size-5 animate-spin text-muted-foreground" />
        </div>
      ) : (
        <div className="space-y-6">
          {inviteOpen && (
            <InviteMemberInline
              organizationId={organizationId}
              onInvited={() => {}}
              onClose={() => setInviteOpen(false)}
            />
          )}

          <div>
            <div className="mb-2 flex items-center justify-between">
              <h3 className="font-semibold text-foreground text-sm">
                Active members ({displayedActive.length})
              </h3>
            </div>
            {displayedActive.length === 0 ? (
              <div className="flex flex-col items-center justify-center rounded-xl border border-border/50 border-dashed py-10 text-center">
                <div className="flex size-10 items-center justify-center rounded-xl bg-muted text-muted-foreground">
                  <Users className="size-5" />
                </div>
                <p className="mt-3 font-medium text-foreground text-sm">No active members found</p>
                <p className="mt-1 text-muted-foreground text-xs">
                  Invite teammates to collaborate on this workspace.
                </p>
              </div>
            ) : (
              <div className="divide-y divide-border/40 rounded-xl border border-border/50 bg-background/50">
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
            <div className="border-border/40 border-t pt-4">
              <div className="mb-2 flex items-center justify-between">
                <h3 className="font-semibold text-foreground text-sm">
                  Pending invitations ({pendingMembers.length})
                </h3>
              </div>
              <div className="divide-y divide-border/40 rounded-xl border border-border/50 bg-background/50">
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
        </div>
      )}
    </SettingsSection>
  );
}
