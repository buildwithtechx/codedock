import type {
  CreateMailAliasRequest,
  Fail2banActionResponse,
  Fail2banStatus,
  MailAlias,
} from '@codedock/contracts';
import type { BaseResponse } from '#/interfaces/base';
import type {
  CreateMailServerRequest,
  DkimKey,
  Mailbox,
  MailboxCredentials,
  MailCert,
  MailCheckReport,
  MailDnsResponse,
  MailDomain,
  MailServer,
} from '#/interfaces/mail';
import { apiClient } from '#/lib/api-client';
import { handleApiError } from '#/lib/error';

export const mailService = {
  listServers: async (): Promise<MailServer[]> => {
    try {
      const res = await apiClient.get<BaseResponse<MailServer[]>>('/mail/servers');
      return res.data;
    } catch (err) {
      throw handleApiError(err);
    }
  },

  createServer: async (payload: CreateMailServerRequest): Promise<MailServer> => {
    try {
      const res = await apiClient.post<BaseResponse<MailServer>>('/mail/servers', payload);
      return res.data;
    } catch (err) {
      throw handleApiError(err);
    }
  },

  deleteServer: async (id: string): Promise<void> => {
    try {
      await apiClient.delete<BaseResponse<null>>(`/mail/servers/${id}`);
    } catch (err) {
      throw handleApiError(err);
    }
  },

  listDomains: async (serverId: string): Promise<MailDomain[]> => {
    try {
      const res = await apiClient.get<BaseResponse<MailDomain[]>>(
        `/mail/servers/${serverId}/domains`
      );
      return res.data;
    } catch (err) {
      throw handleApiError(err);
    }
  },

  createDomain: async (serverId: string, name: string): Promise<MailDomain> => {
    try {
      const res = await apiClient.post<BaseResponse<MailDomain>>(
        `/mail/servers/${serverId}/domains`,
        { name }
      );
      return res.data;
    } catch (err) {
      throw handleApiError(err);
    }
  },

  deleteDomain: async (serverId: string, domainId: string): Promise<void> => {
    try {
      await apiClient.delete<BaseResponse<null>>(`/mail/servers/${serverId}/domains/${domainId}`);
    } catch (err) {
      throw handleApiError(err);
    }
  },

  listMailboxes: async (serverId: string, domainId: string): Promise<Mailbox[]> => {
    try {
      const res = await apiClient.get<BaseResponse<Mailbox[]>>(
        `/mail/servers/${serverId}/domains/${domainId}/mailboxes`
      );
      return res.data;
    } catch (err) {
      throw handleApiError(err);
    }
  },

  createMailbox: async (
    serverId: string,
    domainId: string,
    localPart: string
  ): Promise<MailboxCredentials> => {
    try {
      const res = await apiClient.post<BaseResponse<MailboxCredentials>>(
        `/mail/servers/${serverId}/domains/${domainId}/mailboxes`,
        { localPart }
      );
      return res.data;
    } catch (err) {
      throw handleApiError(err);
    }
  },

  deleteMailbox: async (serverId: string, domainId: string, mailboxId: string): Promise<void> => {
    try {
      await apiClient.delete<BaseResponse<null>>(
        `/mail/servers/${serverId}/domains/${domainId}/mailboxes/${mailboxId}`
      );
    } catch (err) {
      throw handleApiError(err);
    }
  },

  createDkim: async (serverId: string, domainId: string): Promise<DkimKey> => {
    try {
      const res = await apiClient.post<BaseResponse<DkimKey>>(
        `/mail/servers/${serverId}/domains/${domainId}/dkim`,
        { selector: 'mail', algorithm: 'Ed25519' }
      );
      return res.data;
    } catch (err) {
      throw handleApiError(err);
    }
  },

  getDns: async (serverId: string, domainId: string): Promise<MailDnsResponse> => {
    try {
      const res = await apiClient.get<BaseResponse<MailDnsResponse>>(
        `/mail/servers/${serverId}/domains/${domainId}/dns`
      );
      return res.data;
    } catch (err) {
      throw handleApiError(err);
    }
  },

  getChecks: async (serverId: string, domainId: string): Promise<MailCheckReport> => {
    try {
      const res = await apiClient.get<BaseResponse<MailCheckReport>>(
        `/mail/servers/${serverId}/domains/${domainId}/checks`
      );
      return res.data;
    } catch (err) {
      throw handleApiError(err);
    }
  },

  listCerts: async (serverId: string): Promise<MailCert[]> => {
    try {
      const res = await apiClient.get<BaseResponse<MailCert[]>>(`/mail/servers/${serverId}/certs`);
      return res.data;
    } catch (err) {
      throw handleApiError(err);
    }
  },

  getFail2banStatus: async (serverId: string, jail?: string): Promise<Fail2banStatus> => {
    try {
      const path = jail
        ? `/mail/servers/${serverId}/fail2ban/status?jail=${encodeURIComponent(jail)}`
        : `/mail/servers/${serverId}/fail2ban/status`;
      const res = await apiClient.get<BaseResponse<Fail2banStatus>>(path);
      return res.data;
    } catch (err) {
      throw handleApiError(err);
    }
  },

  banIP: async (serverId: string, ip: string, jail: string): Promise<Fail2banActionResponse> => {
    try {
      const res = await apiClient.post<BaseResponse<Fail2banActionResponse>>(
        `/mail/servers/${serverId}/fail2ban/ban`,
        { ip, jail }
      );
      return res.data;
    } catch (err) {
      throw handleApiError(err);
    }
  },

  unbanIP: async (serverId: string, ip: string, jail?: string): Promise<Fail2banActionResponse> => {
    try {
      const res = await apiClient.post<BaseResponse<Fail2banActionResponse>>(
        `/mail/servers/${serverId}/fail2ban/unban`,
        { ip, jail }
      );
      return res.data;
    } catch (err) {
      throw handleApiError(err);
    }
  },

  listAliases: async (serverId: string, domainId: string): Promise<MailAlias[]> => {
    try {
      const res = await apiClient.get<BaseResponse<MailAlias[]>>(
        `/mail/servers/${serverId}/domains/${domainId}/aliases`
      );
      return res.data;
    } catch (err) {
      throw handleApiError(err);
    }
  },

  createAlias: async (
    serverId: string,
    domainId: string,
    req: CreateMailAliasRequest
  ): Promise<MailAlias> => {
    try {
      const res = await apiClient.post<BaseResponse<MailAlias>>(
        `/mail/servers/${serverId}/domains/${domainId}/aliases`,
        req
      );
      return res.data;
    } catch (err) {
      throw handleApiError(err);
    }
  },

  deleteAlias: async (serverId: string, domainId: string, aliasId: string): Promise<void> => {
    try {
      await apiClient.delete<BaseResponse<null>>(
        `/mail/servers/${serverId}/domains/${domainId}/aliases/${aliasId}`
      );
    } catch (err) {
      throw handleApiError(err);
    }
  },
};
