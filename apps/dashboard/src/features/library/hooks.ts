import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { useMemo, useState } from 'react';
import { libraryApi } from './api';
import type { SortBy, VisibilityFilter } from './types';

export const REPOS_PER_PAGE = 20;

export function useGitConnections() {
  return useQuery({
    queryKey: ['library', 'git', 'status'],
    queryFn: () => libraryApi.gitStatus(),
  });
}

export function useConnectProvider() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { provider: string; accessToken: string; accountName: string }) =>
      libraryApi.connectProvider(payload),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['library', 'git'] });
      await queryClient.invalidateQueries({ queryKey: ['git'] });
    },
  });
}

export function useLibraryRepos(provider: string, enabled: boolean) {
  const [search, setSearch] = useState('');
  const [visibility, setVisibility] = useState<VisibilityFilter>('all');
  const [sort, setSort] = useState<SortBy>('updated');
  const [page, setPage] = useState(1);

  const query = useQuery({
    queryKey: ['library', 'repos', provider],
    queryFn: () => libraryApi.listRepos(provider),
    enabled: enabled && provider.length > 0,
  });

  const filtered = useMemo(() => {
    const repos = query.data ?? [];
    const needle = search.trim().toLowerCase();
    let list = repos;
    if (needle) {
      list = list.filter(
        (repo) =>
          repo.name.toLowerCase().includes(needle) || repo.fullName.toLowerCase().includes(needle)
      );
    }
    if (visibility === 'public') list = list.filter((repo) => !repo.private);
    if (visibility === 'private') list = list.filter((repo) => repo.private);
    return [...list].sort((a, b) => {
      if (sort === 'name') return a.name.localeCompare(b.name);
      const left = a.updatedAt ? Date.parse(a.updatedAt) : 0;
      const right = b.updatedAt ? Date.parse(b.updatedAt) : 0;
      return right - left;
    });
  }, [query.data, search, visibility, sort]);

  const totalPages = Math.max(1, Math.ceil(filtered.length / REPOS_PER_PAGE));
  const safePage = Math.min(page, totalPages);
  const pageItems = filtered.slice((safePage - 1) * REPOS_PER_PAGE, safePage * REPOS_PER_PAGE);

  const resetPage =
    <T>(setter: (value: T) => void) =>
    (value: T) => {
      setter(value);
      setPage(1);
    };

  return {
    repos: pageItems,
    total: filtered.length,
    loading: query.isLoading,
    isError: query.isError,
    refetch: query.refetch,
    search,
    setSearch: resetPage(setSearch),
    visibility,
    setVisibility: resetPage(setVisibility),
    sort,
    setSort: resetPage(setSort),
    page: safePage,
    totalPages,
    setPage,
    counts: {
      total: (query.data ?? []).length,
      publicCount: (query.data ?? []).filter((repo) => !repo.private).length,
      privateCount: (query.data ?? []).filter((repo) => repo.private).length,
    },
  };
}
