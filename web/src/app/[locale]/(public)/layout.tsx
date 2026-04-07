import React from "react";
import { SiteHeader } from "@/components/site/site-header";

export default function PublicLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="flex min-h-screen flex-col">
      <SiteHeader />

      <main className="flex-1">
        {children}
      </main>

      {/* Footer sits on top of the wave */}
      <footer className="relative z-20 bg-slate-900 dark:bg-slate-800 py-6 md:py-0">
        <div className="container mx-auto max-w-7xl flex flex-col items-center justify-between gap-4 px-4 md:h-16 md:flex-row">
          <p className="text-center text-sm leading-loose text-slate-400 md:text-left">
            Built by <span className="font-medium text-slate-300">Kaldalis Team</span>.
          </p>
        </div>
      </footer>
    </div>
  );
}
