"use client";

import { useEffect, useState } from "react";
import { usePathname, useRouter } from "@/i18n/routing";
import api from "@/lib/api";
import { Loader2 } from "lucide-react";

export function SystemStatusGuard({ children }: { children: React.ReactNode }) {
  const [isReady, setIsReady] = useState(false);
  const pathname = usePathname();
  const router = useRouter();

  useEffect(() => {
    const checkStatus = async () => {
      // Don't check if we are already on the setup page to avoid recursion
      if (pathname.includes('/setup')) {
        setIsReady(true);
        return;
      }

      try {
        const response = await api.get("/system/status");
        // Backend returns { installed: boolean, site_name: string }
        if (response && response.installed === false) {
          if (!pathname.includes('/setup')) {
            router.replace("/setup");
          } else {
            setIsReady(true);
          }
        } else {
          // 系统已安装
          if (pathname.includes('/setup')) {
            router.replace("/"); // 已安装则不允许回 setup
          } else {
            setIsReady(true);
          }
        }
      } catch (error) {
        console.error("Failed to check system status:", error);
        // If API fails, we assume it's ready or at least let the app try to load
        // to avoid being stuck on a loader forever
        setIsReady(true);
      }
    };

    checkStatus();
  }, [pathname, router]);

  if (!isReady && !pathname.includes('/setup')) {
    return (
      <div className="flex h-screen w-screen items-center justify-center bg-slate-950">
        <div className="flex flex-col items-center gap-4">
          <Loader2 className="h-8 w-8 animate-spin text-indigo-500" />
          <p className="text-slate-400 font-mono text-sm tracking-widest animate-pulse">VERIFYING SYSTEM STATUS...</p>
        </div>
      </div>
    );
  }

  return <>{children}</>;
}
