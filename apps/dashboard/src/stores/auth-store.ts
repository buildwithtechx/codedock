import { create } from 'zustand';
import type { User } from '#/interfaces/users';

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

  const storedToken = localStorage.getItem('codedock_auth_token');
  const storedRefreshToken = localStorage.getItem('codedock_refresh_token');
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
    token: storedToken,
    refreshToken: storedRefreshToken,
    user: parsedUser,
    isAuthenticated: !!(storedToken && parsedUser),
  };
};

export const useAuthStore = create<AuthState>((set) => ({
  ...getInitialAuth(),
  setAuth: (token: string, refreshToken: string | null, user: User) => {
    if (isBrowser) {
      localStorage.setItem('codedock_auth_token', token);
      if (refreshToken) {
        localStorage.setItem('codedock_refresh_token', refreshToken);
      }
      localStorage.setItem('codedock_auth_user', JSON.stringify(user));
    }
    set({ token, refreshToken, user, isAuthenticated: true });
  },
  logout: () => {
    if (isBrowser) {
      localStorage.removeItem('codedock_auth_token');
      localStorage.removeItem('codedock_refresh_token');
      localStorage.removeItem('codedock_auth_user');
    }
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
    if (state.token) {
      localStorage.setItem('codedock_auth_token', state.token);
    } else {
      localStorage.removeItem('codedock_auth_token');
    }
    if (state.refreshToken) {
      localStorage.setItem('codedock_refresh_token', state.refreshToken);
    } else {
      localStorage.removeItem('codedock_refresh_token');
    }
  });
}
