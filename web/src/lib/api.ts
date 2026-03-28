import axios from 'axios';
import Cookies from 'js-cookie';
import { toast } from 'sonner';

// 创建 axios 实例
const api = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL || '/api/v1',
  withCredentials: true, // 必须开启，否则不会发送 Cookie
  timeout: 10000,
});

// 🟢 请求拦截器 (Request Interceptor)
api.interceptors.request.use(
  (config) => {
    // 1. 尝试从浏览器 Cookie 中获取 CSRF Token
    // 注意：CSRF Cookie 名在后端配置为 kaldalis_csrf
    const csrfToken = Cookies.get('kaldalis_csrf');
    
    // 2. 如果拿到了，就塞到 Header 里
    if (csrfToken) {
      config.headers['X-CSRF-Token'] = csrfToken;
    }

    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// 🔵 响应拦截器 (Response Interceptor)
api.interceptors.response.use(
  (response) => {
    // 直接返回 data，简化调用
    return response.data;
  },
  (error) => {
    const message = error.response?.data?.message || error.message || "Unknown error occurred";
    const status = error.response?.status;
    const url = error.config?.url || "";

    // 统一错误提示
    // 401: Unauthorized (expected when checking profile if not logged in)
    // 404: Not Found
    if (status !== 401 && status !== 404) {
      toast.error(message);
    }

    // Only log error if it's not a routine profile check failing
    if (status !== 401 || (!url.includes('/users/profile') && !url.includes('/system/status'))) {
      console.error(`[API Error] ${status}: ${message} (${url})`);
    }
    
    // 如果是 401 (未登录) 且当前不在登录页，跳转登录
    if (status === 401) {
       if (typeof window !== 'undefined' && !window.location.pathname.includes('/login') && !window.location.pathname.includes('/setup')) {
         // window.location.href = '/login'; 
       }
    }

    return Promise.reject(error);
  }
);

export default api;
