import { toast } from 'sonner';
import { env } from '#/env';
import { getAuthHeaders, handleAuthFailure, refreshAuthSession } from '#/lib/auth-refresh';

const API_BASE_URL = env.VITE_API_URL;

export class ApiError extends Error {
  public status: number;
  public data: unknown;

  constructor(status: number, message: string, data?: unknown) {
    super(message);
    this.status = status;
    this.data = data;
    this.name = 'ApiError';
  }
}

async function handleResponseError(response: Response): Promise<never> {
  const isJson = response.headers.get('content-type')?.includes('application/json');
  const data = isJson ? await response.json() : await response.text();
  const errorMessage = data?.message || data?.error || response.statusText || 'An error occurred';
  toast.error(errorMessage);
  throw new ApiError(response.status, errorMessage, data);
}

async function prepareHeaders(options?: RequestInit, body?: unknown): Promise<Headers> {
  const headers = new Headers(options?.headers || {});
  if (!headers.has('Content-Type') && !(body instanceof FormData)) {
    headers.set('Content-Type', 'application/json');
  }
  const authHeaders = getAuthHeaders();
  for (const [key, value] of Object.entries(authHeaders)) {
    if (!headers.has(key)) {
      headers.set(key, value);
    }
  }
  return headers;
}

export const apiClient = {
  async fetch<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
    const url = `${API_BASE_URL}${endpoint}`;
    const headers = await prepareHeaders(options, options.body);

    const response = await fetch(url, {
      ...options,
      headers,
      credentials: 'include',
    });

    if (response.status === 401) {
      const newToken = await refreshAuthSession(API_BASE_URL);
      if (newToken) {
        headers.set('Authorization', `Bearer ${newToken}`);
        const retryResponse = await fetch(url, {
          ...options,
          headers,
          credentials: 'include',
        });
        if (retryResponse.ok) {
          if (retryResponse.status === 204) return {} as T;
          const isJson = retryResponse.headers.get('content-type')?.includes('application/json');
          return (isJson ? await retryResponse.json() : await retryResponse.text()) as T;
        }
      }
      handleAuthFailure();
      throw new ApiError(401, 'Session expired. Please log in again.');
    }

    if (response.status === 204) {
      return {} as T;
    }

    if (!response.ok) {
      return handleResponseError(response);
    }

    const isJson = response.headers.get('content-type')?.includes('application/json');
    return (isJson ? await response.json() : await response.text()) as T;
  },

  get<T>(endpoint: string, options?: RequestInit) {
    return this.fetch<T>(endpoint, { ...options, method: 'GET' });
  },

  async getBlob(endpoint: string, options?: RequestInit): Promise<Blob> {
    const url = `${API_BASE_URL}${endpoint}`;
    const headers = await prepareHeaders(options);
    const response = await fetch(url, {
      ...options,
      method: 'GET',
      headers,
      credentials: 'include',
    });

    if (response.status === 401) {
      const newToken = await refreshAuthSession(API_BASE_URL);
      if (newToken) {
        headers.set('Authorization', `Bearer ${newToken}`);
        const retryResponse = await fetch(url, {
          ...options,
          method: 'GET',
          headers,
          credentials: 'include',
        });
        if (retryResponse.ok) {
          return retryResponse.blob();
        }
      }
      handleAuthFailure();
      throw new ApiError(401, 'Session expired. Please log in again.');
    }

    if (!response.ok) {
      return handleResponseError(response);
    }
    return response.blob();
  },

  async postBlob(endpoint: string, body?: unknown, options?: RequestInit): Promise<Blob> {
    const url = `${API_BASE_URL}${endpoint}`;
    const headers = await prepareHeaders(options, body);
    const response = await fetch(url, {
      ...options,
      method: 'POST',
      body: body instanceof FormData ? body : JSON.stringify(body),
      headers,
      credentials: 'include',
    });

    if (response.status === 401) {
      const newToken = await refreshAuthSession(API_BASE_URL);
      if (newToken) {
        headers.set('Authorization', `Bearer ${newToken}`);
        const retryResponse = await fetch(url, {
          ...options,
          method: 'POST',
          body: body instanceof FormData ? body : JSON.stringify(body),
          headers,
          credentials: 'include',
        });
        if (retryResponse.ok) {
          return retryResponse.blob();
        }
      }
      handleAuthFailure();
      throw new ApiError(401, 'Session expired. Please log in again.');
    }

    if (!response.ok) {
      return handleResponseError(response);
    }
    return response.blob();
  },

  post<T>(endpoint: string, body?: unknown, options?: RequestInit) {
    return this.fetch<T>(endpoint, {
      ...options,
      method: 'POST',
      body: body instanceof FormData ? body : JSON.stringify(body),
    });
  },

  put<T>(endpoint: string, body?: unknown, options?: RequestInit) {
    return this.fetch<T>(endpoint, {
      ...options,
      method: 'PUT',
      body: body instanceof FormData ? body : JSON.stringify(body),
    });
  },

  delete<T>(endpoint: string, options?: RequestInit) {
    return this.fetch<T>(endpoint, { ...options, method: 'DELETE' });
  },
};
