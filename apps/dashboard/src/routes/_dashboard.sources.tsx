import { createFileRoute } from '@tanstack/react-router';
import { PageHeader } from '#/components/layout/page-header';
import { GithubIntegration, GitProviders } from '#/features/sources';

export const Route = createFileRoute('/_dashboard/sources')({
  validateSearch: (search: Record<string, unknown>): { code?: string } => {
    return {
      code: (search.code as string) || undefined,
    };
  },
  component: () => (
    <div className="space-y-8 pb-12">
      <PageHeader
        title="Sources"
        description="Connect source control providers and deployment applications for your organization."
      />
      <GithubIntegration />
      <GitProviders />
    </div>
  ),
});
