export type MaskedContainer = {
  id: string;
  name: string;
  image: string;
  ports: string[];
  envKeys: string[];
  env: Record<string, string>;
  volumes: string[];
  labels: Record<string, string>;
  status: string;
  composeProject?: string;
  routes: string[];
};
export type MaskedStack = {
  containers: MaskedContainer[];
  composeProjects: string[];
  platform: string;
  host: string;
};
export type MigrationSource = {
  id: string;
  organizationId: string;
  name: string;
  sshHost: string;
  sshPort: number;
  sshUser: string;
  sshAuthMethod: string;
  fingerprint: string;
  createdAt: string;
  updatedAt: string;
};
export type MigrationRun = {
  id: string;
  organizationId: string;
  userId: string;
  sourceId: string;
  sourceKind: string;
  projectId: string;
  targetServerId: string;
  mode: string;
  status: string;
  phase: string;
  logs: string;
  error: string;
  cancelRequested: boolean;
  createdAt: string;
  updatedAt: string;
};
export type MigrationPrompt = {
  id: string;
  kind: string;
  subject: string;
  detail: string;
  options: { id: string; label: string; description: string }[];
  expiresAt: number;
};
export type MigrationReview = {
  runId: string;
  confirmation: string;
};
export type MigrationPreview = {
  images: string[];
  volumes: string[];
  services: {
    name: string;
    image: string;
    ports: string[];
    volumes: string[];
    status: string;
  }[];
  conflicts: { kind: string; subject: string; detail: string }[];
  warnings: string[];
  downtime: string;
};
