import { createFileRoute, Link } from '@tanstack/react-router';
import { Building, Loader2, Plus, Trash2 } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '#/components/ui/card';
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
  useCreateOrganization,
  useDeleteOrganization,
  useListOrganizations,
} from '#/features/organizations';

export const Route = createFileRoute('/_dashboard/organizations')({
  component: OrganizationsPage,
});

function OrganizationsPage() {
  const { data: orgs, isLoading } = useListOrganizations();
  const [createOpen, setCreateOpen] = useState(false);

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="flex h-12 w-12 shrink-0 items-center justify-center rounded-lg border border-primary/20 bg-primary/10 text-primary">
            <Building className="h-6 w-6" />
          </div>
          <div>
            <h1 className="font-bold text-xl">Organizations</h1>
            <p className="text-muted-foreground text-sm">Manage your organizations and teams.</p>
          </div>
        </div>

        <Button onClick={() => setCreateOpen(true)} className="gap-2">
          <Plus className="h-4 w-4" />
          NEW ORGANIZATION
        </Button>
      </div>

      {isLoading ? (
        <div className="flex justify-center p-12">
          <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
        </div>
      ) : orgs?.length === 0 ? (
        <div className="flex h-64 flex-col items-center justify-center rounded-xl border border border-dashed bg-card/40">
          <Building className="mb-4 h-8 w-8 text-muted-foreground" />
          <h3 className="font-bold text-foreground text-lg tracking-tight">No organizations yet</h3>
          <p className="mt-1 text-center text-muted-foreground text-sm">
            Create an organization to collaborate with your team.
          </p>
          <Button className="mt-6 gap-2" onClick={() => setCreateOpen(true)}>
            <Plus className="h-4 w-4" />
            CREATE ORGANIZATION
          </Button>
        </div>
      ) : (
        <div className="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3">
          {orgs?.map((org) => (
            <OrganizationCard key={org.id} org={org} />
          ))}
        </div>
      )}

      <CreateOrganizationModal open={createOpen} onOpenChange={setCreateOpen} />
    </div>
  );
}

function OrganizationCard({ org }: { org: any }) {
  const { mutateAsync: deleteOrg, isPending } = useDeleteOrganization();

  const handleDelete = async (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();
    if (!confirm('Are you sure you want to delete this organization?')) return;
    try {
      await deleteOrg(org.id);
      toast.success('Organization deleted');
    } catch (err: any) {
      toast.error(err.message || 'Failed to delete organization');
    }
  };

  return (
    <Link to="/organizations/$organizationId" params={{ organizationId: org.id }}>
      <Card className="cursor-pointer transition-colors hover:border-primary/50">
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="font-medium text-sm">{org.name}</CardTitle>
          <Building className="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <CardDescription className="text-xs">
            Created at {new Date(org.createdAt).toLocaleDateString()}
          </CardDescription>
          <div className="mt-4 flex justify-end">
            <Button variant="destructive" size="sm" onClick={handleDelete} disabled={isPending}>
              {isPending ? (
                <Loader2 className="h-4 w-4 animate-spin" />
              ) : (
                <Trash2 className="h-4 w-4" />
              )}
            </Button>
          </div>
        </CardContent>
      </Card>
    </Link>
  );
}

function CreateOrganizationModal({
  open,
  onOpenChange,
}: {
  open: boolean;
  onOpenChange: (open: boolean) => void;
}) {
  const [name, setName] = useState('');
  const { mutateAsync: createOrg, isPending } = useCreateOrganization();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!name) return;
    try {
      await createOrg({ name });
      toast.success('Organization created');
      onOpenChange(false);
      setName('');
    } catch (err: any) {
      toast.error(err.message || 'Failed to create organization');
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <form onSubmit={handleSubmit}>
          <DialogHeader>
            <DialogTitle>Create Organization</DialogTitle>
            <DialogDescription>
              Create a new organization to manage projects with your team.
            </DialogDescription>
          </DialogHeader>
          <div className="py-6">
            <div className="space-y-2">
              <label htmlFor="name" className="font-medium text-sm">
                Name
              </label>
              <Input
                id="name"
                value={name}
                onChange={(e) => setName(e.target.value)}
                placeholder="Acme Corp"
                disabled={isPending}
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
            <Button type="submit" disabled={!name || isPending}>
              {isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              Create
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
