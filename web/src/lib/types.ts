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

// Post status constants
export const PostStatus = {
  DRAFT: 0,
  PUBLISHED: 1,
  ARCHIVED: 2,
} as const;

export type PostStatusType = typeof PostStatus[keyof typeof PostStatus];

export interface Post {
  id: number;
  title: string;
  slug: string;
  content: string;
  cover?: string;
  status: PostStatusType;
  author_id: number;
  author?: User;
  category_id?: number;
  category?: Category;
  tags?: Tag[];
  created_at: string;
  updated_at: string;
}

export interface CreatePostDTO {
  title: string;
  content: string;
  cover?: string;
  tags?: string[];
  category_id?: number;
}

export interface UpdatePostDTO {
  title?: string;
  content?: string;
  cover?: string;
  tags?: string[];
  category_id?: number;
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

export interface HealthCheckResult {
  status: 'ok' | 'fail' | 'skip';
  detail?: string;
}

export interface HealthResponse {
  status: string;
  mode: string;
  checks: Record<string, HealthCheckResult>;
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

export interface RegisterDTO {
  username: string;
  email: string;
  password: string;
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
