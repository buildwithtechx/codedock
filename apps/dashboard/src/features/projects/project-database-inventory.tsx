import { Link } from '@tanstack/react-router';
import { Database as DatabaseIcon } from 'lucide-react';
import { QueryErrorState } from '#/components/ui/query-error-state';
import { useGetDatabases } from '#/features/databases/hooks';

export function ProjectDatabaseInventory({
  projectId,
  environmentId,
}: {
  projectId: string;
  environmentId?: string;
}) {
  const { data, isLoading, isError, refetch } = useGetDatabases(projectId);
  if (isLoading) return <p className="text-muted-foreground text-sm">Loading databases...</p>;
  if (isError)
    return (
      <QueryErrorState
        title="Databases are unavailable"
        description="Could not load this project's databases."
        onRetry={() => {
          void refetch();
        }}
      />
    );
  const databases = (data?.data || []).filter(
    (database) => !database.environmentId || database.environmentId === environmentId
  );
  if (databases.length === 0) return null;
  return (
    <section className="space-y-3">
      <h2 className="font-semibold text-sm">Databases ({databases.length})</h2>
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        {databases.map((database) => (
          <Link
            key={database.id}
            to="/databases/$databaseId"
            params={{ databaseId: database.id }}
            className="flex items-center gap-3 rounded-xl border bg-card p-4 hover:border-primary/50"
          >
            <DatabaseIcon className="h-5 w-5 text-primary" />
            <div>
              <p className="font-medium">{database.name}</p>
              <p className="text-muted-foreground text-xs">
                {database.engine} ? {database.status}
              </p>
            </div>
          </Link>
        ))}
      </div>
    </section>
  );
}
