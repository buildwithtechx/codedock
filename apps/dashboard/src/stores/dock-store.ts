import { create } from 'zustand';

interface DockState {
  isOpen: boolean;
  activeTab: 'logs' | 'terminal';
  serviceId?: string;
  deploymentId?: string;
  title?: string;
  height: number;
  openDock: (opts?: {
    serviceId?: string;
    deploymentId?: string;
    tab?: 'logs' | 'terminal';
    title?: string;
  }) => void;
  closeDock: () => void;
  toggleDock: () => void;
  setActiveTab: (tab: 'logs' | 'terminal') => void;
  setHeight: (height: number) => void;
}

export const useDockStore = create<DockState>((set) => ({
  isOpen: false,
  activeTab: 'logs',
  height: 280,
  openDock: (opts) =>
    set((state) => ({
      isOpen: true,
      activeTab: opts?.tab ?? state.activeTab,
      serviceId: opts?.serviceId ?? state.serviceId,
      deploymentId: opts?.deploymentId ?? state.deploymentId,
      title: opts?.title ?? state.title,
    })),
  closeDock: () => set({ isOpen: false }),
  toggleDock: () => set((state) => ({ isOpen: !state.isOpen })),
  setActiveTab: (tab) => set({ activeTab: tab }),
  setHeight: (height) => set({ height }),
}));
