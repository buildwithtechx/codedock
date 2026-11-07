export interface Organization {
  id: string;
  name: string;
  createdAt: string;
  updatedAt: string;
}

export type OrganizationRole = 'member' | 'admin' | 'owner' | 'restricted';

export interface OrganizationMember {
  id: string;
  organizationId: string;
  userId?: string;
  email: string;
  permission: OrganizationRole | string;
  status: string;
  invitedAt: string;
  acceptedAt?: string;
  inviteLink?: string;
}

export interface CreateOrganizationRequest {
  name: string;
}

export interface InviteOrganizationMemberRequest {
  email: string;
  permission: OrganizationRole | string;
  delivery?: 'email' | 'link';
  projects?: string[];
}

export interface UpdateOrganizationMemberRequest {
  permission: OrganizationRole | string;
}
