import { createFileRoute } from '@tanstack/react-router';
import { GithubIntegration, GitProviders } from '#/features/sources';

export const Route = createFileRoute('/_dashboard/sources')({
  validateSearch: (search: Record<string, unknown>): { code?: string } => {
    return {
      code: (search.code as string) || undefined,
    };
  },
  component: () => (
    <div className="space-y-12 pb-12">
      <GithubIntegration />
      <div className="h-px w-full bg-border/50" />
      <GitProviders />
    </div>
  ),
});
