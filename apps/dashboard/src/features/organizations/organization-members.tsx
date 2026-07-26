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
import {
  useInviteOrganizationMember,
  useListOrganizationMembers,
  useRemoveOrganizationMember,
  useUpdateOrganizationMember,
} from '#/features/organizations/hooks';
import type { OrganizationMember } from '#/features/organizations/interfaces';

export function OrganizationMembers({ organizationId }: { organizationId: string }) {
  const { data: members, isLoading } = useListOrganizationMembers(organizationId);
  const [inviteOpen, setInviteOpen] = useState(false);

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
              <TableHead className="w-[100px] text-right">Actions</TableHead>
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
              members?.map((member) => (
                <MemberRow key={member.id} member={member} organizationId={organizationId} />
              ))
            )}
          </TableBody>
        </Table>
      </div>

      <InviteMemberModal
        organizationId={organizationId}
        open={inviteOpen}
        onOpenChange={setInviteOpen}
      />
    </div>
  );
}

function MemberRow({
  member,
  organizationId,
}: {
  member: OrganizationMember;
  organizationId: string;
}) {
  const { mutateAsync: updateMember, isPending: isUpdating } =
    useUpdateOrganizationMember(organizationId);
  const { mutateAsync: removeMember, isPending: isRemoving } =
    useRemoveOrganizationMember(organizationId);

  const handlePermissionChange = async (newPermission: string) => {
    try {
      await updateMember({ memberId: member.id, payload: { permission: newPermission } });
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

  return (
    <TableRow>
      <TableCell className="font-medium">{member.email}</TableCell>
      <TableCell>
        <Select
          defaultValue={member.permission}
          onValueChange={handlePermissionChange}
          disabled={isPending}
        >
          <SelectTrigger className="h-8 w-[140px]">
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="Admin">Admin</SelectItem>
            <SelectItem value="Member">Member</SelectItem>
            <SelectItem value="Viewer">Viewer</SelectItem>
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
          disabled={isPending}
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
}: {
  organizationId: string;
  open: boolean;
  onOpenChange: (open: boolean) => void;
}) {
  const [email, setEmail] = useState('');
  const [permission, setPermission] = useState('Member');
  const { mutateAsync: inviteMember, isPending } = useInviteOrganizationMember(organizationId);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!email) return;
    try {
      await inviteMember({ email, permission });
      toast.success('Member invited successfully');
      onOpenChange(false);
      setEmail('');
      setPermission('Member');
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
                  <SelectItem value="Admin">Admin</SelectItem>
                  <SelectItem value="Member">Member</SelectItem>
                  <SelectItem value="Viewer">Viewer</SelectItem>
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
