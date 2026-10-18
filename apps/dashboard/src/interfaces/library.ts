export type VisibilityFilter = 'all' | 'public' | 'private';
export type SortBy = 'updated' | 'name';

export interface ProviderConnection {
  provider: string;
  connected: boolean;
  accountName?: string;
}

export interface LibraryRepo {
  id: string;
  name: string;
  fullName: string;
  cloneUrl: string;
  htmlUrl?: string;
  private: boolean;
  defaultBranch: string;
  updatedAt?: string;
}

export interface PendingImport {
  repositoryUrl: string;
  branch: string;
  name: string;
}

export interface ImportTarget extends PendingImport {
  projectId: string;
}

export interface LocalImport {
  id: string;
  name: string;
  root: string;
  fileCount: number;
  framework: string;
  packageManager: string;
  createdAt: string;
}

export const GIT_PROVIDERS = [
  { id: 'github', name: 'GitHub', icon: '/git-providers/github-icon.svg' },
  { id: 'gitlab', name: 'GitLab', icon: '/git-providers/gitlab-icon.svg' },
  { id: 'bitbucket', name: 'Bitbucket', icon: '/git-providers/bitbucket-icon.svg' },
  { id: 'gitea', name: 'Gitea', icon: '/git-providers/gitea-icon.svg' },
];
