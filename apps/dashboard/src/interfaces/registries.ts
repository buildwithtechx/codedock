export interface Registry {
  id: string;
  projectId: string;
  name: string;
  registryUrl: string;
  username: string;
  createdAt: string;
  updatedAt: string;
}

export interface CreateRegistryRequest {
  name: string;
  registryUrl: string;
  username: string;
  passwordToken: string;
}
