import { create } from 'zustand';

type OrganizationState = {
  activeOrganizationId: string | null;
  setActiveOrganizationId: (organizationId: string | null) => void;
};

const storageKey = 'codedock_active_organization_id';

const initialOrganizationId =
  typeof window === 'undefined' ? null : localStorage.getItem(storageKey) || null;

export const useOrganizationStore = create<OrganizationState>((set) => ({
  activeOrganizationId: initialOrganizationId,
  setActiveOrganizationId: (organizationId) => {
    if (organizationId) {
      localStorage.setItem(storageKey, organizationId);
    } else {
      localStorage.removeItem(storageKey);
    }
    set({ activeOrganizationId: organizationId });
  },
}));
