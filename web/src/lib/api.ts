import axios from 'axios';
import Cookies from 'js-cookie';
import { toast } from 'sonner';

export const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

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

// 🟢 通用请求拦截器
const requestInterceptor = (config: any) => {
  const csrfToken = Cookies.get('kaldalis_csrf');
  if (csrfToken) {
    config.headers['X-CSRF-Token'] = csrfToken;
  }
  return config;
};

api.interceptors.request.use(requestInterceptor);
sysApi.interceptors.request.use(requestInterceptor);

// 🔵 响应拦截器
api.interceptors.response.use(
  (response) => response.data,
  (error) => {
    const message = error.response?.data?.message || error.message || "Unknown error occurred";
    const status = error.response?.status;
    const url = error.config?.url || "";

    if (status !== 401 && status !== 404) {
      toast.error(message);
    }

    if (status !== 401 || (!url.includes('/users/profile') && !url.includes('/system/status'))) {
      console.error(`[API Error] ${status}: ${message} (${url})`);
    }

    return Promise.reject(error);
  }
);

sysApi.interceptors.response.use(
  (response) => response.data,
  (error) => {
    // 系统探测实例不抛出 UI 错误，也不在控制台报错（除非是 500）
    if (error.response?.status >= 500) {
      console.error(`[System API Error] ${error.response.status}: ${error.message}`);
    }
    return Promise.reject(error);
  }
);

export default api;

