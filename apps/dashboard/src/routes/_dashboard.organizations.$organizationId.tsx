import { createFileRoute } from '@tanstack/react-router';
import { Loader2 } from 'lucide-react';
import { useGetOrganization } from '#/features/organizations';
import { OrganizationMembers } from '#/features/organizations/organization-members';

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
      <OrganizationMembers organizationId={organizationId} organizationName={organization.name} />
    </div>
  );
}
