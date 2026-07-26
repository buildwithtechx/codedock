import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { orgService } from './api';
import type {
  CreateOrganizationRequest,
  InviteOrganizationMemberRequest,
  UpdateOrganizationMemberRequest,
} from './types';

export const useListOrganizations = () =>
  useQuery({ queryKey: ['organizations'], queryFn: () => orgService.list() });

export const useCreateOrganization = () => {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (payload: CreateOrganizationRequest) => orgService.create(payload),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['organizations'] }),
  });
};

export const useDeleteOrganization = () => {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => orgService.delete(id),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['organizations'] }),
  });
};

export const useGetOrganization = (id: string) =>
  useQuery({
    queryKey: ['organization', id],
    queryFn: () => orgService.get(id),
    enabled: !!id,
  });

export const useListOrganizationMembers = (id: string) =>
  useQuery({
    queryKey: ['organizationMembers', id],
    queryFn: () => orgService.listMembers(id),
    enabled: !!id,
  });

export const useInviteOrganizationMember = (id: string) => {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (payload: InviteOrganizationMemberRequest) => orgService.inviteMember(id, payload),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['organizationMembers', id] }),
  });
};

export const useUpdateOrganizationMember = (orgId: string) => {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({
      memberId,
      payload,
    }: {
      memberId: string;
      payload: UpdateOrganizationMemberRequest;
    }) => orgService.updateMember(orgId, memberId, payload),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['organizationMembers', orgId] }),
  });
};

export const useRemoveOrganizationMember = (orgId: string) => {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (memberId: string) => orgService.removeMember(orgId, memberId),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['organizationMembers', orgId] }),
  });
};
