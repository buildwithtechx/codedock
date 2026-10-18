import { create } from 'zustand';
import type {
  CodedockApp,
  CodedockConnection,
  CodedockDeployment,
  CodedockProject,
} from '@/interfaces';

interface AppStoreState {
  connection: CodedockConnection | null;
  projects: CodedockProject[];
  apps: CodedockApp[];
  deployments: CodedockDeployment[];
  selectedProject: CodedockProject | null;
  selectedApp: CodedockApp | null;
  isLoading: boolean;
  setConnection: (conn: CodedockConnection | null) => void;
  setProjects: (projects: CodedockProject[]) => void;
  setApps: (apps: CodedockApp[]) => void;
  setDeployments: (deps: CodedockDeployment[]) => void;
  setSelectedProject: (project: CodedockProject | null) => void;
  setSelectedApp: (app: CodedockApp | null) => void;
  setLoading: (loading: boolean) => void;
}

export const useAppStore = create<AppStoreState>((set) => ({
  connection: null,
  projects: [],
  apps: [],
  deployments: [],
  selectedProject: null,
  selectedApp: null,
  isLoading: false,
  setConnection: (connection) => set({ connection }),
  setProjects: (projects) => set({ projects }),
  setApps: (apps) => set({ apps }),
  setDeployments: (deployments) => set({ deployments }),
  setSelectedProject: (selectedProject) => set({ selectedProject }),
  setSelectedApp: (selectedApp) => set({ selectedApp }),
  setLoading: (isLoading) => set({ isLoading }),
}));
