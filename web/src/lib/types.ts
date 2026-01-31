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
  slug: string;
  content: string;
  cover?: string;
  status: number; // 0=draft, 1=published, 2=archived
  author_id: number;
  author?: User;
  category_id?: number;
  category?: any; // Define Category type if needed later
  tags?: any[];   // Define Tag type if needed later
  created_at: string;
  updated_at: string;
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
