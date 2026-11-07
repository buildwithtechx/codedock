import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import type { CreateMailServerRequest } from '#/interfaces/mail';
import { mailService } from '#/services/mail';

export const useListMailServers = () =>
  useQuery({ queryKey: ['mail-servers'], queryFn: () => mailService.listServers() });

export const useCreateMailServer = () => {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (payload: CreateMailServerRequest) => mailService.createServer(payload),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['mail-servers'] }),
  });
};

export const useDeleteMailServer = () => {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => mailService.deleteServer(id),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['mail-servers'] }),
  });
};

export const useListMailDomains = (serverId?: string) =>
  useQuery({
    queryKey: ['mail-servers', serverId, 'domains'],
    queryFn: () => (serverId ? mailService.listDomains(serverId) : Promise.resolve([])),
    enabled: Boolean(serverId),
  });

export const useCreateMailDomain = (serverId?: string) => {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (name: string) => {
      if (!serverId) throw new Error('No mail server selected');
      return mailService.createDomain(serverId, name);
    },
    onSuccess: () => qc.invalidateQueries({ queryKey: ['mail-servers', serverId, 'domains'] }),
  });
};

export const useDeleteMailDomain = (serverId?: string) => {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (domainId: string) => {
      if (!serverId) throw new Error('No mail server selected');
      return mailService.deleteDomain(serverId, domainId);
    },
    onSuccess: () => qc.invalidateQueries({ queryKey: ['mail-servers', serverId, 'domains'] }),
  });
};

export const useListMailboxes = (serverId?: string, domainId?: string) =>
  useQuery({
    queryKey: ['mail-servers', serverId, 'domains', domainId, 'mailboxes'],
    queryFn: () =>
      serverId && domainId ? mailService.listMailboxes(serverId, domainId) : Promise.resolve([]),
    enabled: Boolean(serverId && domainId),
  });

export const useCreateMailbox = (serverId?: string, domainId?: string) => {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (localPart: string) => {
      if (!serverId || !domainId) throw new Error('No domain selected');
      return mailService.createMailbox(serverId, domainId, localPart);
    },
    onSuccess: () =>
      qc.invalidateQueries({
        queryKey: ['mail-servers', serverId, 'domains', domainId, 'mailboxes'],
      }),
  });
};

export const useDeleteMailbox = (serverId?: string, domainId?: string) => {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (mailboxId: string) => {
      if (!serverId || !domainId) throw new Error('No domain selected');
      return mailService.deleteMailbox(serverId, domainId, mailboxId);
    },
    onSuccess: () =>
      qc.invalidateQueries({
        queryKey: ['mail-servers', serverId, 'domains', domainId, 'mailboxes'],
      }),
  });
};

export const useCreateDkim = (serverId?: string, domainId?: string) => {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: () => {
      if (!serverId || !domainId) throw new Error('No domain selected');
      return mailService.createDkim(serverId, domainId);
    },
    onSuccess: () =>
      qc.invalidateQueries({ queryKey: ['mail-servers', serverId, 'domains', domainId] }),
  });
};

export const useMailDns = (serverId?: string, domainId?: string) =>
  useQuery({
    queryKey: ['mail-servers', serverId, 'domains', domainId, 'dns'],
    queryFn: () =>
      serverId && domainId ? mailService.getDns(serverId, domainId) : Promise.resolve(null),
    enabled: Boolean(serverId && domainId),
  });

export const useMailChecks = (serverId?: string, domainId?: string) =>
  useQuery({
    queryKey: ['mail-servers', serverId, 'domains', domainId, 'checks'],
    queryFn: () =>
      serverId && domainId ? mailService.getChecks(serverId, domainId) : Promise.resolve(null),
    enabled: Boolean(serverId && domainId),
  });

export const useMailCerts = (serverId?: string) =>
  useQuery({
    queryKey: ['mail-servers', serverId, 'certs'],
    queryFn: () => (serverId ? mailService.listCerts(serverId) : Promise.resolve([])),
    enabled: Boolean(serverId),
  });
