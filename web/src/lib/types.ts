export interface User {
  id: number;
  username: string;
  email?: string;
  role: 'user' | 'admin' | 'super_admin';
  avatar?: string;
  created_at?: string;
  updated_at?: string;
}

export interface Post {
  id: number;
  title: string;
  content: string;
  excerpt?: string;
  slug: string;
  status: 'draft' | 'published' | 'archived';
  author_id: number;
  author?: User;
  created_at?: string;
  updated_at?: string;
  published_at?: string;
}

export interface AuthResponse {
  message: string;
  user: User;
  token?: string;
}

export interface ApiResponse<T = any> {
  success: boolean;
  message: string;
  data?: T;
  error?: string;
}