"use client";

import React from "react";
import { Link, usePathname } from "@/i18n/routing";
import { useRouter } from "next/navigation";
import { useAuthStore } from "@/store/useAuthStore";
import api from "@/lib/api";
import { useTranslations } from 'next-intl';
import { cn } from "@/lib/utils";
import { 
  LayoutDashboard, 
  FileText, 
  Image, 
  Users, 
  BarChart3, 
  Settings,
  LogOut,
  ChevronRight
} from "lucide-react";

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
  const pathname = usePathname();
  const t = useTranslations('admin');
  const { user, logout } = useAuthStore();

  const handleLogout = async () => {
    try {
      await api.post("/users/logout");
    } catch (error) {
      console.error("登出请求失败:", error);
    }
    logout();
    window.location.href = "/login";
  };

  const navItems = [
    {
      name: t('dashboard'),
      href: "/admin/dashboard",
      icon: LayoutDashboard,
    },
    {
      name: t('content'),
      href: "/admin/posts",
      icon: FileText,
    },
    {
      name: t('media'),
      href: "/admin/media",
      icon: Image,
    },
    {
      name: t('users'),
      href: "/admin/users",
      icon: Users,
    },
    {
      name: t('analytics'),
      href: "/admin/analytics",
      icon: BarChart3,
    },
  ];

  return (
    <div className="flex h-screen w-full overflow-hidden bg-slate-950 text-slate-200">
      {/* 左侧侧边栏 (Sidebar) */}
      <aside className="hidden w-64 flex-col border-r border-slate-800 bg-slate-950 md:flex">
        <div className="flex h-16 items-center border-b border-slate-800 px-6 text-lg font-bold tracking-tight text-white">
          <div className="flex items-center justify-center w-8 h-8 rounded-lg bg-emerald-500 text-slate-950 font-bold mr-3 shadow-[0_0_15px_rgba(16,185,129,0.3)]">K</div>
          Kaldalis
        </div>
        
        <nav className="flex-1 space-y-1 p-4">
          {navItems.map((item) => {
            const isActive = pathname === item.href || pathname.startsWith(`${item.href}/`);
            return (
              <Link 
                key={item.href}
                href={item.href as any} 
                className={cn(
                  "flex items-center justify-between group rounded-lg px-3 py-2.5 text-sm font-medium transition-all duration-200",
                  isActive 
                    ? "bg-emerald-500/10 text-emerald-400 border border-emerald-500/20 shadow-[inset_0_0_10px_rgba(16,185,129,0.05)]" 
                    : "text-slate-400 hover:bg-slate-900 hover:text-slate-200 border border-transparent"
                )}
              >
                <div className="flex items-center gap-3">
                  <item.icon className={cn(
                    "w-5 h-5 transition-colors",
                    isActive ? "text-emerald-400" : "text-slate-500 group-hover:text-slate-300"
                  )} />
                  {item.name}
                </div>
                {isActive && (
                  <div className="w-1.5 h-1.5 rounded-full bg-emerald-500 shadow-[0_0_8px_rgba(16,185,129,0.6)]"></div>
                )}
              </Link>
            );
          })}
        </nav>
        
        {/* Settings & User Info */}
        <div className="mt-auto border-t border-slate-800 p-4 space-y-4">
          <Link 
            href="/admin/settings" 
            className={cn(
              "flex items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium transition-colors",
              pathname.startsWith("/admin/settings")
                ? "bg-slate-800 text-white"
                : "text-slate-400 hover:bg-slate-900 hover:text-white"
            )}
          >
            <Settings className="w-5 h-5 text-slate-500" />
            {t('settings')}
          </Link>

          <div className="flex items-center justify-between bg-slate-900/40 rounded-xl p-3 border border-slate-800/50">
             <div className="flex items-center gap-3">
               <Avatar className="h-9 w-9 border border-slate-700">
                 <AvatarImage src={user?.avatar || ""} />
                 <AvatarFallback className="bg-slate-800 text-slate-200">
                   {user?.username?.[0]?.toUpperCase() || "A"}
                 </AvatarFallback>
               </Avatar>
               <div className="flex flex-col">
                 <span className="text-sm font-semibold text-white truncate max-w-[100px]">{user?.username || "Admin"}</span>
                 <span className="text-[10px] text-slate-500 font-mono uppercase tracking-tighter">System Admin</span>
               </div>
             </div>
             
             <DropdownMenu>
               <DropdownMenuTrigger asChild>
                 <button className="text-slate-500 hover:text-white p-1 rounded-md hover:bg-slate-800 transition-colors">
                    <MoreHorizontal className="w-4 h-4" />
                 </button>
               </DropdownMenuTrigger>
               <DropdownMenuContent align="end" side="top" className="w-56 bg-slate-900 border-slate-800 text-slate-200">
                 <DropdownMenuLabel className="text-xs text-slate-500 font-mono uppercase">Manage Session</DropdownMenuLabel>
                 <DropdownMenuItem className="cursor-pointer focus:bg-slate-800 focus:text-white">
                   <Settings className="mr-2 h-4 w-4" /> {t('settings')}
                 </DropdownMenuItem>
                 <DropdownMenuSeparator className="bg-slate-800" />
                 <DropdownMenuItem onClick={handleLogout} className="text-rose-400 focus:text-rose-300 focus:bg-rose-950/30 cursor-pointer">
                   <LogOut className="mr-2 h-4 w-4" /> {t('logout_text') || 'Logout'}
                 </DropdownMenuItem>
               </DropdownMenuContent>
             </DropdownMenu>
          </div>
        </div>
      </aside>

      {/* 右侧主内容区 */}
      <div className="flex flex-1 flex-col bg-slate-950 relative overflow-hidden">
        {/* Background Pattern */}
        <div className="absolute inset-0 z-0 pointer-events-none">
          <div className="absolute inset-0 bg-[linear-gradient(to_right,#80808012_1px,transparent_1px),linear-gradient(to_bottom,#80808012_1px,transparent_1px)] bg-[size:24px_24px]"></div>
          <div className="absolute left-0 right-0 top-0 -z-10 m-auto h-[310px] w-[310px] rounded-full bg-emerald-500 opacity-20 blur-[100px]"></div>
          <div className="absolute right-0 bottom-0 -z-10 h-[310px] w-[310px] rounded-full bg-indigo-500 opacity-10 blur-[100px]"></div>
        </div>

        <main className="flex-1 overflow-y-auto p-4 md:p-8 relative z-10">
          {children}
        </main>
      </div>
    </div>
  );
}

// 辅助图标
import { MoreHorizontal } from "lucide-react";