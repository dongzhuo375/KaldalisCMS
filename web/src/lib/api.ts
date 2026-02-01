import axios from 'axios';
import Cookies from 'js-cookie';

// 创建 axios 实例
const api = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL || '/api/v1',
  withCredentials: true, // 必须开启，否则不会发送 Cookie
  timeout: 10000,
});

// 🟢 请求拦截器 (Request Interceptor)
// 在这里处理 CSRF Token
api.interceptors.request.use(
  (config) => {
    // 1. 尝试从浏览器 Cookie 中获取 CSRF Token
    const csrfToken = Cookies.get('kaldalis_csrf');
        console.log("🚀 [API Debug] URL:", config.url, "CSRF Token:", csrfToken);
    // 2. 如果拿到了，就塞到 Header 里
    // 后端通常识别的 Header key 是 "X-CSRF-Token" 或 "X-Xsrf-Token"
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
    // 统一错误处理
    console.error("API请求错误:", error.response?.data?.message || error.message);
    
    // 如果是 401 (未登录) 且当前不在登录页，跳转登录
    if (error.response?.status === 401) {
       // 注意：Next.js 的 Router 在这里不能直接用，只能用 window.location
       if (typeof window !== 'undefined' && !window.location.pathname.includes('/login')) {
         // window.location.href = '/login'; // 可选：强制踢下线
       }
    }
    
    // 如果是 403 (CSRF 失败或权限不足)
    if (error.response?.status === 403) {
        console.error("权限不足或 CSRF 校验失败");
    }

    return Promise.reject(error);
  }
);

export default api;
