import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import type { InstallAppInput, OneClickDeployRequest } from '#/interfaces/templates';
import { useOrganizationStore } from '#/stores/organization-store';
import { appsApi } from './api';

export function useInstalledApps() {
  const organizationId = useOrganizationStore((state) => state.activeOrganizationId);
  return useQuery({
    queryKey: ['apps', 'installed', organizationId],
    queryFn: () => appsApi.listInstalled(organizationId as string),
    enabled: Boolean(organizationId),
  });
}

export function useAppCatalog() {
  return useQuery({
    queryKey: ['apps', 'catalog'],
    queryFn: () => appsApi.listCatalog(),
  });
}

export function useCatalogApp(appId: string) {
  return useQuery({
    queryKey: ['apps', 'catalog', appId],
    queryFn: () => appsApi.getCatalogApp(appId),
    enabled: appId.length > 0,
    retry: false,
  });
}

export function useReviewCatalogInstall() {
  return useMutation({
    mutationFn: (payload: InstallAppInput) => appsApi.reviewInstall(payload),
  });
}

export function useDeployCatalogApp() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: OneClickDeployRequest) => appsApi.deployCatalogApp(payload),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['apps'] });
      await queryClient.invalidateQueries({ queryKey: ['projects'] });
      await queryClient.invalidateQueries({ queryKey: ['canvas'] });
    },
  });
}
