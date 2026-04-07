"use client";

import { useEffect } from "react";
import { usePathname, useRouter } from "@/i18n/routing";
import { useSystemStatus } from "@/services/system-service";
import { Loader2 } from "lucide-react";

export function SystemStatusGuard({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();
  const router = useRouter();
  const { data: status, isLoading, isError } = useSystemStatus();

  useEffect(() => {
    if (isLoading) return;

    const isSetupPage = pathname.includes('/setup');

    if (status && status.installed === false) {
      if (!isSetupPage) {
        router.replace("/setup");
      }
    } else if (status && status.installed === true) {
      if (isSetupPage) {
        router.replace("/");
      }
    }
  }, [status, isLoading, pathname, router]);

  if (isLoading && !pathname.includes('/setup')) {
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
