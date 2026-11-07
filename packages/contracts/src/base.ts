/** Standard API response wrapper returned by Codedock HTTP endpoints. */
export interface BaseResponse<T = unknown> {
  status: 'success' | 'error' | 'warning';
  message: string;
  data: T;
  path?: string;
  executionTime?: number;
}

/** Paginated collection wrapper for list queries. */
export interface PaginatedData<T> {
  records: T[];
  total: number;
  page: number;
  totalPages: number;
}
