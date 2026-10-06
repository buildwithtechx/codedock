import { useState } from 'react';
import { Label } from '#/components/ui/label';
import { useGitStatus, useListGitRepos } from '#/hooks/use-git';

export function RepositorySelector({
  onSelect,
}: {
  onSelect: (url: string, branch: string) => void;
}) {
  const [provider, setProvider] = useState('');
  const connections = useGitStatus();
  const repositories = useListGitRepos(provider);
  return (
    <div className="grid gap-3 sm:grid-cols-2">
      <div className="space-y-2">
        <Label htmlFor="setup-provider">Connected Git account</Label>
        <select
          id="setup-provider"
          className="w-full rounded-md border bg-background p-2"
          value={provider}
          onChange={(event) => setProvider(event.target.value)}
        >
          <option value="">Enter a URL manually</option>
          {connections.data?.data
            ?.filter((connection) => connection.connected)
            .map((connection) => (
              <option key={connection.provider} value={connection.provider}>
                {connection.provider}
              </option>
            ))}
        </select>
      </div>
      {provider && (
        <div className="space-y-2">
          <Label htmlFor="setup-repository">Repository</Label>
          <select
            id="setup-repository"
            className="w-full rounded-md border bg-background p-2"
            defaultValue=""
            onChange={(event) => {
              const repository = repositories.data?.data?.find(
                (item) => item.id === event.target.value
              );
              if (repository) onSelect(repository.cloneUrl, repository.defaultBranch);
            }}
          >
            <option value="">Select a repository</option>
            {repositories.data?.data?.map((repository) => (
              <option key={repository.id} value={repository.id}>
                {repository.fullName}
              </option>
            ))}
          </select>
        </div>
      )}
      {(connections.isError || repositories.isError) && (
        <p role="alert" className="text-destructive text-sm">
          Git account or repository list unavailable. You can still enter a repository URL.
        </p>
      )}
    </div>
  );
}
