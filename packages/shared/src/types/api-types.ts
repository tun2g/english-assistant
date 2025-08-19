import type { HttpStatus } from '../constants/api-constants';

export interface ApiResponse<T = unknown> {
  data: T;
  message?: string;
  status: HttpStatus;
  timestamp: string;
}

export interface ApiError {
  message: string;
  status: HttpStatus;
  code?: string;
  details?: Record<string, unknown>;
  timestamp: string;
}

export interface PaginationMeta {
  page: number;
  limit: number;
  total: number;
  totalPages: number;
  hasNext: boolean;
  hasPrev: boolean;
}

export interface PaginatedResponse<T> extends ApiResponse<T[]> {
  meta: PaginationMeta;
}

export interface PaginationParams {
  page?: number;
  limit?: number;
  sort?: string;
  order?: 'asc' | 'desc';
}

export interface ApiRequestConfig {
  timeout?: number;
  retries?: number;
  headers?: Record<string, string>;
}