import { Loader2, Plus, Users } from 'lucide-react';
import { useState } from 'react';
import { Button } from '#/components/ui/button';
import { SettingsSection } from '#/features/settings/settings-section';
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
      <SettingsSection
        icon={<Users className="size-4 text-primary" />}
        title="Team Members"
        description="Manage who has access to this workspace and its deployed projects."
        action={
          canManageTeam ? (
            <Button onClick={() => setInviteOpen(true)} size="sm" className="shrink-0 gap-2">
              <Plus className="size-4" />
              Invite member
            </Button>
          ) : undefined
        }
      >
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
          <div className="-mx-5 -my-5 divide-y divide-border/40">
            {members.map((member) => (
              <div key={member.id} className="px-5">
                <MemberRow
                  member={member}
                  organizationId={organizationId}
                  isCurrentUserOwner={isCurrentUserOwner}
                  canManageTeam={canManageTeam}
                  ownerCount={ownerCount}
                  currentUserId={user?.id}
                  currentUserEmail={user?.email}
                />
              </div>
            ))}
          </div>
        )}
      </SettingsSection>

      <InviteMemberModal
        organizationId={organizationId}
        open={inviteOpen}
        onOpenChange={setInviteOpen}
        isCurrentUserOwner={isCurrentUserOwner}
      />
    </div>
  );
}
