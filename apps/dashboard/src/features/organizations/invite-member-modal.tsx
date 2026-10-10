import { Loader2 } from 'lucide-react';
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
import { Label } from '#/components/ui/label';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '#/components/ui/select';
import { useInviteOrganizationMember } from './hooks';

export function InviteMemberModal({
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
    if (!email.trim()) return;
    try {
      await inviteMember({ email: email.trim(), permission });
      toast.success('Member invited successfully');
      onOpenChange(false);
      setEmail('');
      setPermission('member');
    } catch (err: unknown) {
      const message = err instanceof Error ? err.message : 'Failed to invite member';
      toast.error(message);
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-md">
        <form onSubmit={handleSubmit} className="space-y-4">
          <DialogHeader>
            <DialogTitle>Invite Teammate</DialogTitle>
            <DialogDescription>
              Invite a new collaborator to access this workspace and its deployed projects.
            </DialogDescription>
          </DialogHeader>

          <div className="space-y-4 py-2">
            <div className="space-y-2">
              <Label htmlFor="invite-email">Email Address</Label>
              <Input
                id="invite-email"
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                placeholder="colleague@example.com"
                required
                disabled={isPending}
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="invite-role">Workspace Role</Label>
              <Select value={permission} onValueChange={setPermission} disabled={isPending}>
                <SelectTrigger id="invite-role">
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
            <Button type="submit" disabled={!email.trim() || isPending}>
              {isPending && <Loader2 className="mr-2 size-4 animate-spin" />}
              Send Invite
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
