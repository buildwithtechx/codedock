export interface Organization {
  id: string;
  name: string;
  createdAt: string;
  updatedAt: string;
}

export interface OrganizationMember {
  id: string;
  organizationId: string;
  userId: string;
  permission: string;
  status: string;
  invitedAt: string;
  acceptedAt: string;
}

export interface CreateOrganizationRequest {
  name: string;
}
