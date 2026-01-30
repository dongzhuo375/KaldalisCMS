"use client"; // 1. 必须变身客户端组件才能交互

import React from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useAuthStore } from "@/store/useAuthStore";
import api from "@/lib/api";
import { useTranslations } from 'next-intl';

// 引入 shadcn 组件
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";

export default function AdminLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const router = useRouter();
  const t = useTranslations('admin');
  // 2. 从 Store 获取用户信息和清理方法
  const { user, logout } = useAuthStore();

  // 3. 登出核心逻辑
  const handleLogout = async () => {
    try {
      // 调用后端清除 HttpOnly Cookie (路径必须对)
      await api.post("/users/logout");
    } catch (error) {
      console.error("登出请求失败:", error);
      // 即使后端报错，前端也要强制登出
    }

    // 清除前端 Zustand 状态
    logout();

    // 强制刷新跳转 (清除内存残留)
    window.location.href = "/login";
  };

  return (
    <div className="flex h-screen w-full overflow-hidden bg-slate-950 text-slate-200">
      {/* 左侧侧边栏 (Sidebar) */}
      <aside className="hidden w-64 flex-col border-r border-slate-800 bg-slate-950 md:flex">
        <div className="flex h-16 items-center border-b border-slate-800 px-6 text-lg font-bold tracking-tight text-white">
          <div className="flex items-center justify-center w-8 h-8 rounded-lg bg-emerald-500 text-slate-950 font-bold mr-3">K</div>
          Kaldalis
        </div>
        <nav className="flex-1 space-y-1 p-4">
          <Link href="/admin/dashboard" className="flex items-center gap-3 rounded-lg bg-slate-800/50 px-3 py-2 text-sm font-medium text-white border border-slate-700/50">
            <svg className="w-5 h-5 text-emerald-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z" />
            </svg>
            {t('dashboard')}
          </Link>
          <Link href="/admin/posts" className="flex items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium text-slate-400 hover:bg-slate-900 hover:text-white transition-colors">
            <svg className="w-5 h-5 text-slate-500 group-hover:text-slate-300" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
            </svg>
            {t('content')}
          </Link>
          <Link href="/admin/media" className="flex items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium text-slate-400 hover:bg-slate-900 hover:text-white transition-colors">
            <svg className="w-5 h-5 text-slate-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
            </svg>
            {t('media')}
          </Link>
          <Link href="/admin/users" className="flex items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium text-slate-400 hover:bg-slate-900 hover:text-white transition-colors">
            <svg className="w-5 h-5 text-slate-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />
            </svg>
            {t('users')}
          </Link>
          <Link href="/admin/analytics" className="flex items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium text-slate-400 hover:bg-slate-900 hover:text-white transition-colors">
            <svg className="w-5 h-5 text-slate-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
            </svg>
            {t('analytics')}
          </Link>
        </nav>
        
        {/* Settings 放在底部 */}
        <div className="p-4 mt-auto space-y-1">
          <Link href="/admin/settings" className="flex items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium text-slate-400 hover:bg-slate-900 hover:text-white transition-colors">
            <svg className="w-5 h-5 text-slate-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
            </svg>
            {t('settings')}
          </Link>
        </div>

        <div className="border-t border-slate-800 p-4">
           {/* 底部用户信息 */}
           <div className="flex items-center gap-3">
             <Avatar className="h-9 w-9 border border-slate-700">
               <AvatarImage src={user?.avatar || ""} />
               <AvatarFallback className="bg-slate-800 text-slate-200">
                 {user?.username?.[0]?.toUpperCase() || "A"}
               </AvatarFallback>
             </Avatar>
             <div className="flex flex-col">
               <span className="text-sm font-medium text-white">{user?.username || "Admin User"}</span>
               <span className="text-xs text-slate-500">{user?.email || "sysop@kaldalis.io"}</span>
             </div>
           </div>
        </div>
      </aside>

      {/* 右侧主内容区 */}
      <div className="flex flex-1 flex-col bg-slate-950">
        {/* Header 已经被集成到 Dashboard 页面内部了，这里可以移除或简化
            根据设计图，Dashboard 页面自己有一个 Topbar。
            所以这里的 Header 如果存在，应该是非常极简的，或者干脆没有。
            设计图左侧 Sidebar，右侧是一个巨大的 Dashboard 面板。
            Dashboard 面板内部有 "root @ ..." 的 Header。
            所以 AdminLayout 的 Header 应该移除，给 Dashboard 全屏空间。
        */}
        
        {/* 具体的页面内容 */}
        <main className="flex-1 overflow-y-auto p-4 md:p-8">
          {children}
        </main>
      </div>
    </div>
  );
}
