import { useAuthStore } from '#/stores/auth-store';

let activeRefreshPromise: Promise<string | null> | null = null;

export function getCookie(name: string): string | null {
  const value = `; ${document.cookie}`;
  const parts = value.split(`; ${name}=`);
  if (parts.length === 2) return parts.pop()?.split(';').shift() ?? null;
  return null;
}

export function getAuthHeaders(): Record<string, string> {
  const headers: Record<string, string> = {};
  const token = useAuthStore.getState().token;
  if (token) {
    headers.Authorization = `Bearer ${token}`;
  }
  const csrfToken = getCookie('csrf_token');
  if (csrfToken) {
    headers['X-CSRF-Token'] = csrfToken;
  }
  return headers;
}

export async function refreshAuthSession(apiBaseUrl: string): Promise<string | null> {
  if (activeRefreshPromise) {
    return activeRefreshPromise;
  }

  activeRefreshPromise = (async () => {
    const authState = useAuthStore.getState();
    if (!authState.refreshToken) {
      return null;
    }

    try {
      const res = await fetch(`${apiBaseUrl}/auth/refresh`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ refreshToken: authState.refreshToken }),
        credentials: 'include',
      });

      if (!res.ok) {
        return null;
      }

      const data = await res.json();
      if (data.data?.token && data.data?.user) {
        authState.setAuth(data.data.token, data.data.refreshToken || null, data.data.user);
        return data.data.token as string;
      }
      return null;
    } catch (error) {
      console.error('Failed to refresh authentication session', error);
      return null;
    } finally {
      activeRefreshPromise = null;
    }
  })();

  return activeRefreshPromise;
}

export function handleAuthFailure(): void {
  useAuthStore.getState().logout();
  const path = window.location.pathname;
  const isAuthPage =
    path.startsWith('/signin') ||
    path.startsWith('/signup') ||
    path.startsWith('/forgot-password') ||
    path.startsWith('/reset-password') ||
    path.startsWith('/onboarding') ||
    path.startsWith('/setup');

  if (!isAuthPage) {
    window.location.href = '/signin';
  }
}
