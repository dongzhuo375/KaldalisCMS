// src/lib/api.ts
import axios from 'axios';

const api = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL,
  withCredentials: true, // 关键：允许跨域携带 Cookie/凭证
  timeout: 10000,
});

// 响应拦截器：简单的错误提示处理
api.interceptors.response.use(
  (response) => response.data, // 直接返回后端返回的 JSON 数据部分
  (error) => {
    // 这里先不做复杂的跳转，只把错误往外抛
    console.error("API 错误:", error.response?.data || error.message);
    return Promise.reject(error);
  }
);

export default api;
