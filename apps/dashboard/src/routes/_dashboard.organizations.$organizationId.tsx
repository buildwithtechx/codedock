import { createFileRoute } from '@tanstack/react-router';
import { Building, Loader2 } from 'lucide-react';
import { OrganizationMembers } from '#/features/organizations/organization-members';
import { useGetOrganization } from '#/hooks/useOrganizations';

export const Route = createFileRoute('/_dashboard/organizations/$organizationId')({
  component: OrganizationDetailRoute,
});

function OrganizationDetailRoute() {
  const { organizationId } = Route.useParams();
  const { data: organization, isLoading } = useGetOrganization(organizationId);

  if (isLoading) {
    return (
      <div className="flex justify-center p-12">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  if (!organization) {
    return (
      <div className="flex h-64 items-center justify-center rounded-lg border border-dashed text-muted-foreground">
        Organization not found
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-3">
        <div className="flex h-12 w-12 shrink-0 items-center justify-center rounded-lg border border-primary/20 bg-primary/10 text-primary">
          <Building className="h-6 w-6" />
        </div>
        <div>
          <h1 className="font-bold text-xl">{organization.name}</h1>
          <p className="text-muted-foreground text-sm">
            Created on {new Date(organization.createdAt).toLocaleDateString()}
          </p>
        </div>
      </div>

      <section className="rounded-lg border bg-card p-6 shadow-sm">
        <OrganizationMembers organizationId={organizationId} />
      </section>
    </div>
  );
}
