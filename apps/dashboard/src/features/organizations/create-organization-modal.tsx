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
import { useOrganizationStore } from '#/stores/organization-store';
import { useCreateOrganization } from './hooks';

export function CreateOrganizationModal({
  open,
  onOpenChange,
  onCreated,
}: {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onCreated?: (orgId: string) => void;
}) {
  const [name, setName] = useState('');
  const { mutateAsync: createOrg, isPending } = useCreateOrganization();
  const setActiveOrg = useOrganizationStore((s) => s.setActiveOrganizationId);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim()) return;
    try {
      const res = await createOrg({ name: name.trim() });
      toast.success('Organization created');
      const newOrgId = (res as any)?.id;
      if (newOrgId) {
        setActiveOrg(newOrgId);
        onCreated?.(newOrgId);
      }
      onOpenChange(false);
      setName('');
    } catch (err: unknown) {
      toast.error(err instanceof Error ? err.message : 'Failed to create organization');
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <form onSubmit={handleSubmit}>
          <DialogHeader>
            <DialogTitle>Create Organization</DialogTitle>
            <DialogDescription>
              Create a new organization to manage projects and resources with your team.
            </DialogDescription>
          </DialogHeader>
          <div className="py-6">
            <div className="space-y-2">
              <Label htmlFor="org-name" className="font-medium text-sm">
                Organization Name
              </Label>
              <Input
                id="org-name"
                value={name}
                onChange={(e) => setName(e.target.value)}
                placeholder="Acme Corp"
                disabled={isPending}
                autoFocus
              />
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
            <Button type="submit" disabled={!name.trim() || isPending}>
              {isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              Create
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
