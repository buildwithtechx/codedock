import type { ApiResponse, CliContext } from './types.js';

export class ApiClient {
  private readonly baseUrl: string;
  private readonly token?: string;

  constructor(ctx: CliContext) {
    this.baseUrl = ctx.serverUrl.replace(/\/+$/, '');
    this.token = ctx.token;
  }

  private async request<T>(
    method: string,
    endpoint: string,
    body?: unknown
  ): Promise<ApiResponse<T>> {
    const cleanEndpoint = endpoint.startsWith('/') ? endpoint : `/${endpoint}`;
    const url = `${this.baseUrl}${cleanEndpoint}`;

    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
      Accept: 'application/json',
    };

    if (this.token) {
      headers.Authorization = `Bearer ${this.token}`;
    }

    const init: RequestInit = {
      method,
      headers,
    };

    if (body !== undefined) {
      init.body = JSON.stringify(body);
    }

    let res: Response;
    try {
      res = await fetch(url, init);
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : String(err);
      throw new Error(`Failed to connect to Codedock server at ${this.baseUrl}: ${msg}`);
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

    return data;
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
