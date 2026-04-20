import axios, { AxiosError, InternalAxiosRequestConfig } from 'axios';
import Cookies from 'js-cookie';
import { toast } from 'sonner';
import type { ApiErrorEnvelope, ErrorInteractionStrategy, NormalizedApiError } from '@/lib/types';

export const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

const DEFAULT_ERROR_CODE = 'INTERNAL_ERROR';

const strategyByCode: Record<string, ErrorInteractionStrategy> = {
  VALIDATION_FAILED: { toast: false, retryable: false, redirectToLogin: false },
  UNAUTHORIZED: { toast: false, retryable: false, redirectToLogin: true },
  FORBIDDEN: { toast: true, retryable: false, redirectToLogin: false },
  NOT_FOUND: { toast: false, retryable: false, redirectToLogin: false },
  DUPLICATE_RESOURCE: { toast: true, retryable: false, redirectToLogin: false },
  CONFLICT: { toast: true, retryable: false, redirectToLogin: false },
  TIMEOUT: { toast: true, retryable: true, redirectToLogin: false },
  INTERNAL_ERROR: { toast: true, retryable: true, redirectToLogin: false },
};

// 创建 axios 实例 (业务 API)
const api = axios.create({
  baseURL: '/api/v1',
  withCredentials: true,
  timeout: 10000,
});

// 创建系统级 axios 实例 (用于根路径探测，如 healthz/readyz)
export const sysApi = axios.create({
  baseURL: '/',
  withCredentials: true,
  timeout: 5000,
});

const requestInterceptor = (config: InternalAxiosRequestConfig) => {
  const csrfToken = Cookies.get('kaldalis_csrf');
  if (csrfToken) {
    config.headers['X-CSRF-Token'] = csrfToken;
  }
  return config;
};

api.interceptors.request.use(requestInterceptor);
sysApi.interceptors.request.use(requestInterceptor);

function isApiErrorEnvelope(value: unknown): value is ApiErrorEnvelope {
  if (!value || typeof value !== 'object') {
    return false;
  }
  const obj = value as Record<string, unknown>;
  return typeof obj.code === 'string' && typeof obj.message === 'string' && 'details' in obj;
}

export function normalizeApiError(error: unknown): NormalizedApiError {
  const axiosErr = error as AxiosError<ApiErrorEnvelope>;
  const status = axiosErr.response?.status ?? 0;
  const data = axiosErr.response?.data;

  const envelope = isApiErrorEnvelope(data)
    ? data
    : {
        code: DEFAULT_ERROR_CODE,
        message: axiosErr.message || 'Unknown error occurred',
        details: null,
      };

  const details = envelope.details && typeof envelope.details === 'object' ? { ...envelope.details } : {};
  const requestIdHeader = axiosErr.response?.headers?.['x-request-id'];
  const requestId =
    (typeof details.request_id === 'string' ? details.request_id : undefined) ||
    (typeof requestIdHeader === 'string' ? requestIdHeader : undefined);

  return {
    status,
    code: envelope.code || DEFAULT_ERROR_CODE,
    message: envelope.message || 'Unknown error occurred',
    details,
    requestId,
    url: axiosErr.config?.url,
    raw: error,
  };
}

export function getErrorInteractionStrategy(code: string): ErrorInteractionStrategy {
  return strategyByCode[code] || strategyByCode.INTERNAL_ERROR;
}

api.interceptors.response.use(
  (response) => response.data,
  (error) => {
    const normalized = normalizeApiError(error);
    const strategy = getErrorInteractionStrategy(normalized.code);
    const url = normalized.url || '';

    if (strategy.toast) {
      toast.error(normalized.message);
    }

    if (!(normalized.status === 401 && (url.includes('/users/profile') || url.includes('/system/status')))) {
      console.warn(
        `[API Error] status=${normalized.status} code=${normalized.code} request_id=${normalized.requestId || '-'} message=${normalized.message} url=${url}`
      );
    }

    return Promise.reject(normalized);
  }
);

sysApi.interceptors.response.use(
  (response) => response.data,
  (error) => {
    const normalized = normalizeApiError(error);
    if (normalized.status >= 500) {
      console.error(
        `[System API Error] status=${normalized.status} code=${normalized.code} request_id=${normalized.requestId || '-'} message=${normalized.message}`
      );
    }
    return Promise.reject(normalized);
  }
);

export default api;

