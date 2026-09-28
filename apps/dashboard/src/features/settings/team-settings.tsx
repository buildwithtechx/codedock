import { Building2 } from 'lucide-react';
import { useListOrganizations } from '#/features/organizations';
import { OrganizationMembers } from '#/features/organizations/organization-members';
import { useOrganizationStore } from '#/stores/organization-store';

export function TeamSettings() {
  const activeOrganizationId = useOrganizationStore((state) => state.activeOrganizationId);
  const { data: organizations, isLoading } = useListOrganizations();
  const organizationId = activeOrganizationId ?? organizations?.[0]?.id;

  if (isLoading) {
    return <div className="h-48 animate-pulse rounded-xl border border-border/80 bg-card" />;
  }

  if (!organizationId) {
    return (
      <div className="flex min-h-56 flex-col items-center justify-center rounded-xl border border-border border-dashed bg-card px-6 text-center">
        <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-primary/10 text-primary">
          <Building2 className="h-5 w-5" />
        </div>
        <p className="mt-4 font-semibold text-sm">No workspace available</p>
        <p className="mt-1 max-w-sm text-muted-foreground text-sm">
          Create a workspace before inviting teammates.
        </p>
      </div>
    );
  }

  return <OrganizationMembers organizationId={organizationId} />;
}
