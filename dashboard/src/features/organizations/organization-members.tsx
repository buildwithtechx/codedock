import { Loader2, Plus, Trash2 } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '#/components/ui/dialog';
import { Input } from '#/components/ui/input';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '#/components/ui/select';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '#/components/ui/table';
import type { OrganizationMember } from '#/features/organizations';
import {
  useInviteOrganizationMember,
  useListOrganizationMembers,
  useRemoveOrganizationMember,
  useUpdateOrganizationMember,
} from '#/features/organizations';
import { useAuthStore } from '#/stores/auth-store';

export function OrganizationMembers({ organizationId }: { organizationId: string }) {
  const { data: members, isLoading } = useListOrganizationMembers(organizationId);
  const [inviteOpen, setInviteOpen] = useState(false);
  const user = useAuthStore((s) => s.user);

  const currentUserMember = members?.find((m) => m.userId === user?.id || m.email === user?.email);
  const isCurrentUserOwner = currentUserMember?.permission === 'owner';

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <div>
          <h3 className="font-semibold text-lg">Organization Members</h3>
          <p className="text-muted-foreground text-sm">
            Manage who has access to this organization and its projects.
          </p>
        </div>
        <Button onClick={() => setInviteOpen(true)} size="sm">
          <Plus className="mr-2 h-4 w-4" />
          Invite Member
        </Button>
      </div>

      <div className="rounded-md border">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Email</TableHead>
              <TableHead>Permission</TableHead>
              <TableHead>Status</TableHead>
              <TableHead className="w-25 text-right">Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {isLoading ? (
              <TableRow>
                <TableCell colSpan={4} className="h-24 text-center">
                  <Loader2 className="mx-auto h-6 w-6 animate-spin text-muted-foreground" />
                </TableCell>
              </TableRow>
            ) : members?.length === 0 ? (
              <TableRow>
                <TableCell colSpan={4} className="h-24 text-center text-muted-foreground">
                  No members found.
                </TableCell>
              </TableRow>
            ) : (
              (() => {
                const ownerCount = members?.filter((m) => m.permission === 'owner').length || 0;
                return members?.map((member) => (
                  <MemberRow
                    key={member.id}
                    member={member}
                    organizationId={organizationId}
                    isCurrentUserOwner={isCurrentUserOwner}
                    ownerCount={ownerCount}
                  />
                ));
              })()
            )}
          </TableBody>
        </Table>
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

function MemberRow({
  member,
  organizationId,
  isCurrentUserOwner,
  ownerCount,
}: {
  member: OrganizationMember;
  organizationId: string;
  isCurrentUserOwner: boolean;
  ownerCount: number;
}) {
  const { mutateAsync: updateMember, isPending: isUpdating } =
    useUpdateOrganizationMember(organizationId);
  const { mutateAsync: removeMember, isPending: isRemoving } =
    useRemoveOrganizationMember(organizationId);

  const handlePermissionChange = async (newPermission: string) => {
    try {
      await updateMember({ memberId: member.userId!, payload: { permission: newPermission } });
      toast.success('Permission updated');
    } catch (err: any) {
      toast.error(err.message || 'Failed to update permission');
    }
  };

  const handleRemove = async () => {
    if (!confirm('Are you sure you want to remove this member?')) return;
    try {
      await removeMember(member.id);
      toast.success('Member removed');
    } catch (err: any) {
      toast.error(err.message || 'Failed to remove member');
    }
  };

  const isPending = isUpdating || isRemoving;
  const canEditRole = isCurrentUserOwner && member.permission !== 'owner';

  return (
    <TableRow>
      <TableCell className="font-medium">{member.email}</TableCell>
      <TableCell>
        <Select
          defaultValue={member.permission}
          onValueChange={handlePermissionChange}
          disabled={isPending || !canEditRole}
        >
          <SelectTrigger className="h-8 w-35">
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            {isCurrentUserOwner && <SelectItem value="owner">Owner</SelectItem>}
            <SelectItem value="admin">Admin</SelectItem>
            <SelectItem value="member">Member</SelectItem>
          </SelectContent>
        </Select>
      </TableCell>
      <TableCell>
        <span
          className={`inline-flex items-center rounded-full px-2 py-0.5 font-medium text-xs ${
            member.status === 'active'
              ? 'bg-green-100 text-green-700'
              : 'bg-yellow-100 text-yellow-700'
          }`}
        >
          {member.status}
        </span>
      </TableCell>
      <TableCell className="text-right">
        <Button
          variant="ghost"
          size="icon"
          onClick={handleRemove}
          disabled={
            isPending || (member.permission === 'owner' && (ownerCount <= 1 || !isCurrentUserOwner))
          }
          className="text-destructive hover:bg-destructive/10 hover:text-destructive"
        >
          {isRemoving ? (
            <Loader2 className="h-4 w-4 animate-spin" />
          ) : (
            <Trash2 className="h-4 w-4" />
          )}
        </Button>
      </TableCell>
    </TableRow>
  );
}

function InviteMemberModal({
  organizationId,
  open,
  onOpenChange,
  isCurrentUserOwner,
}: {
  organizationId: string;
  open: boolean;
  onOpenChange: (open: boolean) => void;
  isCurrentUserOwner: boolean;
}) {
  const [email, setEmail] = useState('');
  const [permission, setPermission] = useState('member');
  const { mutateAsync: inviteMember, isPending } = useInviteOrganizationMember(organizationId);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!email) return;
    try {
      await inviteMember({ email, permission });
      toast.success('Member invited successfully');
      onOpenChange(false);
      setEmail('');
      setPermission('member');
    } catch (err: any) {
      toast.error(err.message || 'Failed to invite member');
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <form onSubmit={handleSubmit}>
          <DialogHeader>
            <DialogTitle>Invite Member</DialogTitle>
            <DialogDescription>Invite a new member to join this organization.</DialogDescription>
          </DialogHeader>
          <div className="space-y-4 py-4">
            <div className="space-y-2">
              <label htmlFor="email" className="font-medium text-sm">
                Email Address
              </label>
              <Input
                id="email"
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                placeholder="colleague@example.com"
                disabled={isPending}
              />
            </div>
            <div className="space-y-2">
              <label htmlFor="permission" className="font-medium text-sm">
                Role
              </label>
              <Select value={permission} onValueChange={setPermission} disabled={isPending}>
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  {isCurrentUserOwner && <SelectItem value="owner">Owner</SelectItem>}
                  <SelectItem value="admin">Admin</SelectItem>
                  <SelectItem value="member">Member</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>
          <DialogFooter>
            <Button
              type="button"
              variant="outline"
              onClick={() => onOpenChange(false)}
              disabled={isPending}
            >
              Cancel
            </Button>
            <Button type="submit" disabled={!email || isPending}>
              {isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              Send Invite
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
