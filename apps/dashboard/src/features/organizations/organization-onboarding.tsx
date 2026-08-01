import { Building, Loader2 } from 'lucide-react';
import { type FormEvent, useEffect, useState } from 'react';
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
import { useAuthStore } from '#/stores/auth-store';
import { useCreateOrganization, useListOrganizations } from './hooks';

export function OrganizationOnboarding() {
  const user = useAuthStore((state) => state.user);
  const { data: organizations, isLoading, isError } = useListOrganizations();
  const { mutateAsync: createOrganization, isPending } = useCreateOrganization();
  const [name, setName] = useState('');
  const [open, setOpen] = useState(false);

  useEffect(() => {
    if (!isLoading && !isError && organizations?.length === 0) {
      const suggestedName = user?.name ? `${user.name}'s Workspace` : 'My Workspace';
      setName(suggestedName);
      setOpen(true);
    }
  }, [isError, isLoading, organizations, user?.name]);

  const create = async (organizationName: string) => {
    try {
      await createOrganization({ name: organizationName.trim() });
      setOpen(false);
      toast.success('Organization created');
    } catch (error) {
      toast.error(error instanceof Error ? error.message : 'Failed to create organization');
    }
  };

  const handleUseSuggestedName = () => create('');

  const handleSubmit = (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    void create(name);
  };

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogContent showCloseButton={false}>
        <DialogHeader>
          <div className="mb-2 flex h-11 w-11 items-center justify-center rounded-xl border border-primary/20 bg-primary/10 text-primary">
            <Building className="h-5 w-5" />
          </div>
          <DialogTitle>Create your organization</DialogTitle>
          <DialogDescription>
            Organizations keep your projects and team members together.
          </DialogDescription>
        </DialogHeader>
        <form onSubmit={handleSubmit}>
          <div className="space-y-2 py-4">
            <label htmlFor="organization-onboarding-name" className="font-medium text-sm">
              Organization name
            </label>
            <Input
              id="organization-onboarding-name"
              value={name}
              onChange={(event) => setName(event.target.value)}
              disabled={isPending}
              autoFocus
            />
          </div>
          <DialogFooter>
            <Button
              type="button"
              variant="outline"
              onClick={handleUseSuggestedName}
              disabled={isPending}
            >
              {isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              Create one for me
            </Button>
            <Button type="submit" disabled={!name.trim() || isPending}>
              {isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              Create organization
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
