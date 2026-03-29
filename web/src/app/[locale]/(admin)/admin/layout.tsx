"use client";

import React from "react";
import { Link, usePathname, useRouter } from "@/i18n/routing";
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
  MoreHorizontal
} from "lucide-react";

import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { useParams } from "next/navigation";

export default function AdminLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const pathname = usePathname();
  const router = useRouter();
  const params = useParams();
  const rawLocale = params?.locale as string;
  const locale = (rawLocale && rawLocale !== 'undefined') ? rawLocale : 'zh-CN';
  
  const t = useTranslations('admin');
  const { user, isLoggedIn, logout } = useAuthStore();

  // 登录守卫
  React.useEffect(() => {
    if (!isLoggedIn) {
      router.replace("/login");
    }
  }, [isLoggedIn, router]);

  const handleLogout = async () => {
    try {
      await api.post("/users/logout");
    } catch (error) {
      console.error("登出请求失败:", error);
    }
    logout();
    router.push("/login");
  };

  if (!isLoggedIn) {
    return null;
  }

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
    <div className="flex h-screen w-full overflow-hidden bg-background text-foreground selection:bg-accent/30">
      {/* Sidebar: Sharp, minimal, high-contrast */}
      <aside className="hidden w-64 flex-col border-r border-border bg-white dark:bg-slate-900 md:flex z-20">
        <div className="flex h-20 items-center px-8 border-b border-border">
          <Link href="/" className="flex items-center gap-3 group">
            <div className="flex items-center justify-center w-8 h-8 rounded-full bg-accent text-white font-serif italic text-xl font-bold transition-transform group-hover:scale-110">K</div>
            <span className="font-serif text-2xl font-medium tracking-tight">Kaldalis</span>
          </Link>
        </div>
        
        <nav className="flex-1 space-y-1 px-4 py-6">
          {navItems.map((item) => {
            const isActive = pathname === item.href || pathname.startsWith(`${item.href}/`);
            return (
              <Link 
                key={item.href}
                href={item.href} 
                className={cn(
                  "flex items-center gap-3 rounded-lg px-4 py-2.5 text-sm font-medium transition-all duration-200",
                  isActive 
                    ? "bg-accent text-white shadow-lg shadow-accent/20" 
                    : "text-muted-foreground hover:bg-muted hover:text-foreground"
                )}
              >
                <item.icon className={cn(
                  "w-4 h-4 transition-colors",
                  isActive ? "text-white" : "text-muted-foreground group-hover:text-foreground"
                )} strokeWidth={isActive ? 2.5 : 2} />
                {item.name}
              </Link>
            );
          })}
        </nav>
        
        {/* Sidebar Footer */}
        <div className="mt-auto border-t border-border p-4 space-y-4">
          <Link 
            href="/admin/settings" 
            className={cn(
              "flex items-center gap-3 rounded-lg px-4 py-2 text-sm font-medium transition-colors",
              pathname.startsWith("/admin/settings")
                ? "text-accent bg-accent/5"
                : "text-muted-foreground hover:text-foreground"
            )}
          >
            <Settings className="w-4 h-4" />
            {t('settings')}
          </Link>

          <div className="flex items-center justify-between bg-muted/50 rounded-xl p-3 border border-border">
             <div className="flex items-center gap-3">
               <Avatar className="h-8 w-8 border border-border">
                 <AvatarImage src={user?.avatar || ""} />
                 <AvatarFallback className="bg-primary text-primary-foreground text-[10px]">
                   {user?.username?.[0]?.toUpperCase() || "A"}
                 </AvatarFallback>
               </Avatar>
               <div className="flex flex-col min-w-0">
                 <span className="text-xs font-bold truncate">{user?.username || "Admin"}</span>
                 <span className="text-[10px] text-muted-foreground uppercase tracking-widest font-medium">Administrator</span>
               </div>
             </div>
             
             <DropdownMenu>
               <DropdownMenuTrigger asChild>
                 <button className="text-muted-foreground hover:text-foreground transition-colors p-1">
                    <MoreHorizontal className="w-4 h-4" />
                 </button>
               </DropdownMenuTrigger>
               <DropdownMenuContent align="end" side="top" className="w-56 bg-white dark:bg-slate-900 border-border shadow-xl">
                 <DropdownMenuLabel className="text-[10px] font-bold uppercase tracking-widest text-muted-foreground">Session</DropdownMenuLabel>
                 <DropdownMenuItem className="cursor-pointer">
                   <Settings className="mr-2 h-4 w-4" /> Account Settings
                 </DropdownMenuItem>
                 <DropdownMenuSeparator className="bg-border" />
                 <DropdownMenuItem onClick={handleLogout} className="text-accent focus:text-white focus:bg-accent cursor-pointer">
                   <LogOut className="mr-2 h-4 w-4" /> {t('logout_text') || 'Logout'}
                 </DropdownMenuItem>
               </DropdownMenuContent>
             </DropdownMenu>
          </div>
        </div>
      </aside>

      {/* Main Content Area */}
      <div className="flex flex-1 flex-col bg-background relative overflow-hidden">
        {/* Subtle Background Elements */}
        <div className="absolute inset-0 z-0 pointer-events-none opacity-40">
          <div className="absolute inset-0 bg-grain" />
          <div className="absolute inset-0 bg-[linear-gradient(to_right,var(--border)_1px,transparent_1px),linear-gradient(to_bottom,var(--border)_1px,transparent_1px)] bg-[size:40px_40px] [mask-image:radial-gradient(ellipse_60%_50%_at_50%_0%,#000_70%,transparent_100%)]" />
        </div>

        <main className="flex-1 overflow-hidden p-6 md:p-10 relative z-10 flex flex-col">
          {children}
        </main>
      </div>
    </div>
  );
}
