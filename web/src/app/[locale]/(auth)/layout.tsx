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
    <div className="flex min-h-screen items-center justify-center p-4">
      {/* SunWaveBackground is already provided by root layout */}
      <div className="relative z-10 w-full flex justify-center">
        {children}
      </div>
    </div>
  );
}
