import { create } from 'zustand';
import { getApiBaseUrl } from '#/lib/api-client';
import { useAuthStore } from './auth-store';

type OrganizationState = {
  activeOrganizationId: string | null;
  setActiveOrganizationId: (organizationId: string | null) => void;
  clearActiveOrganizationId: () => void;
};

const getStorageKey = () => {
  const user = useAuthStore.getState().user;
  const baseUrl = getApiBaseUrl() || 'default';
  const userId = user?.id || 'anonymous';
  return `codedock_active_org_${encodeURIComponent(baseUrl)}_${encodeURIComponent(userId)}`;
};

const getInitialOrganizationId = () => {
  if (typeof window === 'undefined') return null;
  return localStorage.getItem(getStorageKey()) || null;
};

export const useOrganizationStore = create<OrganizationState>((set) => ({
  activeOrganizationId: getInitialOrganizationId(),
  setActiveOrganizationId: (organizationId) => {
    if (typeof window !== 'undefined') {
      const key = getStorageKey();
      if (organizationId) {
        localStorage.setItem(key, organizationId);
      } else {
        localStorage.removeItem(key);
      }
      localStorage.removeItem('codedock_active_organization_id');
    }
    set({ activeOrganizationId: organizationId });
  },
  clearActiveOrganizationId: () => {
    if (typeof window !== 'undefined') {
      localStorage.removeItem(getStorageKey());
      localStorage.removeItem('codedock_active_organization_id');
    }
    set({ activeOrganizationId: null });
  },
}));

if (typeof window !== 'undefined') {
  useAuthStore.subscribe((state, prevState) => {
    if (!state.isAuthenticated || !state.user || state.user.id !== prevState?.user?.id) {
      if (!state.user) {
        useOrganizationStore.getState().clearActiveOrganizationId();
      } else {
        const stored = localStorage.getItem(getStorageKey());
        useOrganizationStore.setState({ activeOrganizationId: stored || null });
      }
    }
  });
}
