import { watchUploadStream } from './stream-upload.js';
import type { ApiResponse, CliContext } from './types.js';

export class ApiClient {
  private readonly baseUrl: string;
  private readonly token?: string;
  private csrfToken?: string;
  private csrfPromise?: Promise<string>;

  constructor(ctx: CliContext) {
    this.baseUrl = ctx.serverUrl.replace(/\/+$/, '');
    this.token = ctx.token;
  }

  private async request<T>(
    method: string,
    endpoint: string,
    body?: unknown,
    contentType?: string
  ): Promise<ApiResponse<T>> {
    const cleanEndpoint = endpoint.startsWith('/') ? endpoint : `/${endpoint}`;
    const url = `${this.baseUrl}${cleanEndpoint}`;

    const headers: Record<string, string> = {
      'Content-Type': contentType || 'application/json',
      Accept: 'application/json',
    };

    if (this.token) {
      headers.Authorization = `Bearer ${this.token}`;
    }

    const init: RequestInit = {
      method,
      headers,
    };

    let streamError: unknown;
    if (body instanceof ReadableStream) {
      init.body = watchUploadStream(body, (error) => {
        streamError = error;
      });
      Object.assign(init, { duplex: 'half' });
    } else if (body instanceof FormData) {
      delete headers['Content-Type'];
      init.body = body;
    } else if (body !== undefined) {
      init.body = JSON.stringify(body);
    }
    const authRequest = /\/auth\//.test(cleanEndpoint);
    let res: Response;
    try {
      if (method !== 'GET' && method !== 'HEAD' && (!this.token || authRequest)) {
        const token = await this.getCsrfToken();
        headers['X-CSRF-Token'] = token;
        headers.Cookie = `csrf_token=${token}`;
      }
      res = await fetch(url, init);
    } catch (err: unknown) {
      if (streamError !== undefined) throw streamError;
      const msg = err instanceof Error ? err.message : String(err);
      throw new Error(`Failed to connect to Codedock server at ${this.baseUrl}: ${msg}`, {
        cause: err,
      });
    }

    const text = await res.text();
    let data: ApiResponse<T>;
    try {
      data = text ? JSON.parse(text) : {};
    } catch {
      throw new Error(
        `Invalid JSON response from server (HTTP ${res.status}): ${text.slice(0, 150)}`
      );
    }

    if (!res.ok) {
      const errorMsg = data.message || `Request failed with HTTP status ${res.status}`;
      throw new Error(errorMsg);
    }

    return Array.isArray(data) ? { data: data as T } : data;
  }

  private async getCsrfToken(): Promise<string> {
    if (this.csrfToken) return this.csrfToken;
    if (this.csrfPromise) return this.csrfPromise;
    this.csrfPromise = (async () => {
      try {
        const response = await fetch(`${this.baseUrl}/api/auth/csrf`);
        if (!response.ok) throw new Error('Failed to initialize CSRF protection');
        const data = (await response.json()) as { token?: string };
        if (!data.token || data.token === '_echo_csrf_using_sec_fetch_site_')
          throw new Error('Server did not return a CSRF token');
        this.csrfToken = data.token;
        return data.token;
      } finally {
        this.csrfPromise = undefined;
      }
    })();
    return this.csrfPromise;
  }

  async postStream<T>(
    endpoint: string,
    body: ReadableStream<Uint8Array>,
    contentType: string
  ): Promise<ApiResponse<T>> {
    return this.request<T>('POST', endpoint, body, contentType);
  }

  async get<T>(endpoint: string): Promise<ApiResponse<T>> {
    return this.request<T>('GET', endpoint);
  }

  async post<T>(endpoint: string, body?: unknown): Promise<ApiResponse<T>> {
    return this.request<T>('POST', endpoint, body);
  }

  async put<T>(endpoint: string, body?: unknown): Promise<ApiResponse<T>> {
    return this.request<T>('PUT', endpoint, body);
  }

  async patch<T>(endpoint: string, body?: unknown): Promise<ApiResponse<T>> {
    return this.request<T>('PATCH', endpoint, body);
  }

  async del<T>(endpoint: string): Promise<ApiResponse<T>> {
    return this.request<T>('DELETE', endpoint);
  }
}
