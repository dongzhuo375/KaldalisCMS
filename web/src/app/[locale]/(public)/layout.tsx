import React from "react";
import { SiteHeader } from "@/components/site/site-header";

export default function PublicLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="flex min-h-screen flex-col bg-slate-50">
      <SiteHeader />
      
      {/* 🟢 修改重点：
          1. mx-auto: 居中
          2. max-w-7xl: 限制最大宽度 (约 1280px)，不让它无限拉伸
          3. px-4: 手机端左右留点缝隙，不贴边
      */}
      <main className="flex-1 container mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
        {children}
      </main>

      <footer className="border-t bg-white py-6 md:py-0">
        <div className="container mx-auto max-w-7xl flex flex-col items-center justify-between gap-4 px-4 md:h-24 md:flex-row">
          <p className="text-center text-sm leading-loose text-muted-foreground md:text-left">
            Built by <span className="font-medium underline underline-offset-4">Kaldalis Team</span>.
          </p>
        </div>
      </footer>
    </div>
  );
}
