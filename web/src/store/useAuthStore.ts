// src/store/useAuthStore.ts
import { create } from 'zustand';
import { persist } from 'zustand/middleware';

interface User {
  username: string;
  role: string;
  avatar?: string;
}

interface AuthState {
  user: User | null;
  isLoggedIn: boolean;
  setLogin: (user: User) => void;
  logout: () => void;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      user: null,
      isLoggedIn: false,
      setLogin: (user) => set({ user, isLoggedIn: true }),
      logout: () => {
        // 这里后续可以调用后端的 logout 接口
        set({ user: null, isLoggedIn: false });
      },
    }),
    {
      name: 'kaldalis-auth-storage', // 存储在 localStorage 里的键名
    }
  )
);
