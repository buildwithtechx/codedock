import { create } from 'zustand';
import type { CodedockApp, CodedockConnection, CodedockDeployment } from '@/interfaces';

interface AppStoreState {
  connection: CodedockConnection | null;
  apps: CodedockApp[];
  deployments: CodedockDeployment[];
  selectedApp: CodedockApp | null;
  isLoading: boolean;
  setConnection: (conn: CodedockConnection | null) => void;
  setApps: (apps: CodedockApp[]) => void;
  setDeployments: (deps: CodedockDeployment[]) => void;
  setSelectedApp: (app: CodedockApp | null) => void;
  setLoading: (loading: boolean) => void;
}

export const useAppStore = create<AppStoreState>((set) => ({
  connection: null,
  apps: [],
  deployments: [],
  selectedApp: null,
  isLoading: false,
  setConnection: (connection) => set({ connection }),
  setApps: (apps) => set({ apps }),
  setDeployments: (deployments) => set({ deployments }),
  setSelectedApp: (selectedApp) => set({ selectedApp }),
  setLoading: (isLoading) => set({ isLoading }),
}));
