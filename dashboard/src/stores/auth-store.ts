import { create } from 'zustand';
import type { User } from '#/interfaces';

export interface AuthState {
  token: string | null;
  refreshToken: string | null;
  user: User | null;
  isAuthenticated: boolean;
  setAuth: (token: string, refreshToken: string | null, user: User) => void;
  logout: () => void;
}

const isBrowser = typeof window !== 'undefined';

const getInitialAuth = () => {
  if (!isBrowser) {
    return { token: null, refreshToken: null, user: null, isAuthenticated: false };
  }

  localStorage.removeItem('codedock_auth_token');
  localStorage.removeItem('codedock_refresh_token');

  const storedUser = localStorage.getItem('codedock_auth_user');
  let parsedUser = null;
  if (storedUser) {
    try {
      parsedUser = JSON.parse(storedUser);
    } catch (e) {
      console.error('Failed to parse stored user:', e);
      localStorage.removeItem('codedock_auth_user');
    }
  }

  return {
    token: null,
    refreshToken: null,
    user: parsedUser,
    isAuthenticated: !!parsedUser,
  };
};

export const useAuthStore = create<AuthState>((set) => ({
  ...getInitialAuth(),
  setAuth: (token: string, refreshToken: string | null, user: User) => {
    set({ token, refreshToken, user, isAuthenticated: true });
  },
  logout: () => {
    set({ token: null, refreshToken: null, user: null, isAuthenticated: false });
  },
}));

if (isBrowser) {
  useAuthStore.subscribe((state) => {
    if (state.user) {
      localStorage.setItem('codedock_auth_user', JSON.stringify(state.user));
    } else {
      localStorage.removeItem('codedock_auth_user');
    }
  });
}
