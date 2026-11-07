import { Boxes, Container, Database, GitBranch, LayoutTemplate, SquareCode } from 'lucide-react';

export type DeployPathId = 'git' | 'docker' | 'database' | 'compose' | 'one-click' | 'examples';

export interface DeployPath {
  id: DeployPathId;
  title: string;
  description: string;
  hint: string;
  icon: typeof GitBranch;
}

export const DEPLOY_PATHS: DeployPath[] = [
  {
    id: 'git',
    title: 'Git repository',
    description:
      'Deploy source code from GitHub, GitLab, Bitbucket, or Gitea with automatic builds.',
    hint: 'Best for applications you develop and push.',
    icon: GitBranch,
  },
  {
    id: 'docker',
    title: 'Docker image',
    description: 'Run a prebuilt image from any public or private registry without a build step.',
    hint: 'Best for published images and third-party tools.',
    icon: Container,
  },
  {
    id: 'database',
    title: 'Database',
    description: 'Provision PostgreSQL, MySQL, Redis, or another managed database for the project.',
    hint: 'Best for stateful backing services.',
    icon: Database,
  },
  {
    id: 'compose',
    title: 'Docker Compose',
    description: 'Import a Compose file and review every service before Codedock creates them.',
    hint: 'Best for multi-service stacks you already defined.',
    icon: Boxes,
  },
  {
    id: 'one-click',
    title: 'One-click app',
    description: 'Install a verified application from the catalog with guided configuration.',
    hint: 'Best for standard tools with zero setup.',
    icon: LayoutTemplate,
  },
  {
    id: 'examples',
    title: 'Example project',
    description: 'Deploy a runnable starter to learn the platform or fork it into your own code.',
    hint: 'Best for exploring and prototyping.',
    icon: SquareCode,
  },
];

export function deployPathById(id: DeployPathId): DeployPath {
  return DEPLOY_PATHS.find((path) => path.id === id) ?? DEPLOY_PATHS[0];
}
