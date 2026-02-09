"use client";

import { useEffect } from "react";
import { useRouter } from "@/i18n/routing";
import { useAuthStore } from "@/store/useAuthStore";

export default function AuthLayout({ children }: { children: React.ReactNode }) {
  const { isLoggedIn } = useAuthStore();
  const router = useRouter();

  useEffect(() => {
    if (isLoggedIn) {
      router.replace("/");
    }
  }, [isLoggedIn, router]);

  if (isLoggedIn) return null;

  return (
    <div className="flex min-h-screen items-center justify-center bg-slate-100 dark:bg-slate-950 p-4 relative">
        {/* 背景装饰 */}
        <div className="absolute inset-0 z-0 pointer-events-none overflow-hidden">
            <div className="absolute top-1/4 left-1/4 w-96 h-96 bg-blue-500/10 rounded-full blur-[100px]"></div>
            <div className="absolute bottom-1/4 right-1/4 w-96 h-96 bg-purple-500/10 rounded-full blur-[100px]"></div>
        </div>
        
        <div className="relative z-10 w-full flex justify-center">
            {children}
        </div>
    </div>
  );
}