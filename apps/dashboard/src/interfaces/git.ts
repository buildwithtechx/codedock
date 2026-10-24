import type { BaseResponse } from './base';

export interface GitStatus {
  provider: string;
  connected: boolean;
}

export interface GitRepo {
  id: string;
  name: string;
  fullName: string;
  cloneUrl: string;
  private: boolean;
  defaultBranch: string;
  updatedAt: string;
}

export interface GitBranch {
  name: string;
  sha: string;
}

export interface ConnectGitRequest {
  provider: string;
  code?: string;
  accessToken?: string;
  accountName?: string;
}

export interface RepoTreeEntry {
  path: string;
  type: string;
}

export interface MonorepoCandidateApp {
  name: string;
  rootDirectory: string;
  framework: string;
  internalPort: number;
  buildEngine: string;
  packageManager?: string;
  installCommand?: string;
  buildCommand?: string;
  startCommand?: string;
  dockerfilePath?: string;
}

export interface RepositoryInspection {
  framework: string;
  packageManager: string;
  installCommand: string;
  buildCommand: string;
  startCommand: string;
  internalPort: number;
  staticOutput: string;
  buildEngine: string;
  dockerfilePath: string;
  warnings: string[];
  candidateApps?: MonorepoCandidateApp[];
  isMonorepo?: boolean;
}

export type GetGitStatusResponse = BaseResponse<GitStatus[]>;
export type ListGitReposResponse = BaseResponse<GitRepo[]>;
export type ListGitBranchesResponse = BaseResponse<GitBranch[]>;
export type ConnectGitResponse = BaseResponse<GitStatus>;
