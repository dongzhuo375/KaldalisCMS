export interface User {
  id: number;
  username: string;
  email?: string;
  role: 'user' | 'admin' | 'super_admin';
  avatar?: string;
  created_at?: string;
  updated_at?: string;
}

export interface Category {
  id: number;
  name: string;
  slug: string;
}

export interface Tag {
  id: number;
  name: string;
}

export interface Post {
  id: number;
  title: string;
  slug: string;
  content: string;
  cover?: string;
  status: number; // 0=draft, 1=published, 2=archived
  author_id: number;
  author?: User;
  category_id?: number;
  category?: Category;
  tags?: Tag[];
  created_at: string;
  updated_at: string;
}

export interface AuthResponse {
  message: string;
  user: User;
  token?: string;
}

export interface ApiResponse<T = unknown> {
  success: boolean;
  message: string;
  data?: T;
  error?: string;
}

export interface SystemStatus {
  installed: boolean;
  site_name?: string;
  version?: string;
}

export interface SetupDTO {
  site_name: string;
  admin_username: string;
  admin_email: string;
  admin_password: string;
}

export interface LoginDTO {
  username?: string;
  email?: string;
  password?: string;
}

export interface MediaAsset {
  id: number;
  filename: string;
  url: string;
  size: number;
  mime_type: string;
  created_at: string;
}

export interface MediaListResponse {
  items: MediaAsset[];
  total: number;
  page: number;
  page_size: number;
}

export interface MediaUploadResponse {
  asset: MediaAsset;
}
