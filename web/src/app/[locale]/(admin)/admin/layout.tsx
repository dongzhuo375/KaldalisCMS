"use client";

import React, { useState, useEffect } from "react";
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
  MoreHorizontal,
  ChevronLeft,
  ChevronRight,
  Menu,
  Command
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
import { Button } from "@/components/ui/button";
import { useParams } from "next/navigation";
import { motion, AnimatePresence } from "framer-motion";

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
  
  // Initialize from localStorage if available
  const [isCollapsed, setIsCollapsed] = useState(false);

  useEffect(() => {
    const saved = localStorage.getItem("sidebar-collapsed");
    if (saved !== null) setIsCollapsed(saved === "true");
    
    // Shortcut listener
    const handleKeyDown = (e: KeyboardEvent) => {
      if ((e.metaKey || e.ctrlKey) && e.key === 'b') {
        e.preventDefault();
        toggleSidebar();
      }
    };
    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, []);

  const toggleSidebar = () => {
    setIsCollapsed(prev => {
      const next = !prev;
      localStorage.setItem("sidebar-collapsed", String(next));
      return next;
    });
  };

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

  if (!isLoggedIn) return null;

  const navItems = [
    { name: t('dashboard'), href: "/admin/dashboard", icon: LayoutDashboard },
    { name: t('content'), href: "/admin/posts", icon: FileText },
    { name: t('media'), href: "/admin/media", icon: Image },
    { name: t('users'), href: "/admin/users", icon: Users },
    { name: t('analytics'), href: "/admin/analytics", icon: BarChart3 },
  ];

  return (
    <div className="flex h-screen w-full bg-transparent text-foreground selection:bg-accent/30 overflow-hidden">
      
      {/* Sidebar: Creative Collapsible Design */}
      <motion.aside 
        initial={false}
        animate={{ 
          width: isCollapsed ? 80 : 280,
          boxShadow: isCollapsed ? "0 0 0 rgba(0,0,0,0)" : "20px 0 50px rgba(0,0,0,0.02)"
        }}
        transition={{ duration: 0.5, ease: [0.4, 0, 0.2, 1] }}
        className="hidden md:flex flex-col border-r border-border bg-white/60 dark:bg-slate-900/60 backdrop-blur-2xl z-20 relative group/sidebar"
      >
        {/* Creative Trigger: The entire right edge of the sidebar */}
        <div 
          onClick={toggleSidebar}
          className="absolute -right-1 top-0 w-2 h-full cursor-col-resize hover:bg-accent/10 transition-colors z-30 group"
          title={isCollapsed ? "Expand Sidebar (Cmd+B)" : "Collapse Sidebar (Cmd+B)"}
        >
           <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 opacity-0 group-hover:opacity-100 transition-opacity bg-accent rounded-full p-0.5">
              {isCollapsed ? <ChevronRight className="w-3 h-3 text-white" /> : <ChevronLeft className="w-3 h-3 text-white" />}
           </div>
        </div>

        {/* Brand Section: Logo also triggers toggle */}
        <div className="flex h-24 items-center px-6 overflow-hidden">
          <div 
            onClick={toggleSidebar}
            className="flex items-center gap-4 group/logo cursor-pointer shrink-0"
          >
            <motion.div 
              whileTap={{ scale: 0.9 }}
              className="flex items-center justify-center w-10 h-10 rounded-2xl bg-primary text-primary-foreground font-serif italic text-2xl font-bold shadow-lg shadow-primary/20 group-hover/logo:rotate-12 transition-transform shrink-0"
            >
              K
            </motion.div>
            <AnimatePresence>
              {!isCollapsed && (
                <motion.div
                  initial={{ opacity: 0, x: -10 }}
                  animate={{ opacity: 1, x: 0 }}
                  exit={{ opacity: 0, x: -10 }}
                  className="flex flex-col"
                >
                  <span className="font-serif text-xl font-bold tracking-tight whitespace-nowrap">Kaldalis</span>
                  <span className="text-[9px] font-bold uppercase tracking-[0.2em] text-accent">Headless Engine</span>
                </motion.div>
              )}
            </AnimatePresence>
          </div>
        </div>
        
        {/* Navigation */}
        <nav className="flex-1 space-y-2 px-4 py-6 overflow-y-auto custom-scrollbar">
          {navItems.map((item) => {
            const isActive = pathname === item.href || pathname.startsWith(`${item.href}/`);
            return (
              <motion.div key={item.href} whileHover={{ x: 4 }} whileTap={{ scale: 0.98 }} className="relative group/nav">
                <Link 
                  href={item.href} 
                  className={cn(
                    "flex items-center rounded-2xl px-4 py-3 text-sm font-medium transition-all duration-300",
                    isCollapsed ? "justify-center" : "gap-4",
                    isActive 
                      ? "bg-primary text-primary-foreground shadow-xl shadow-primary/20" 
                      : "text-muted-foreground hover:bg-muted hover:text-foreground"
                  )}
                >
                  <item.icon className={cn(
                    "w-5 h-5 shrink-0 transition-transform duration-300",
                    isActive ? "scale-110" : "group-hover/nav:scale-110"
                  )} strokeWidth={isActive ? 2.5 : 2} />
                  
                  {!isCollapsed && (
                    <motion.span
                      initial={{ opacity: 0 }}
                      animate={{ opacity: 1 }}
                      className="whitespace-nowrap"
                    >
                      {item.name}
                    </motion.span>
                  )}

                  {/* Tooltip for Collapsed State */}
                  {isCollapsed && (
                    <div className="absolute left-[calc(100%+1.5rem)] px-3 py-2 bg-slate-900 text-white text-[10px] font-bold uppercase tracking-[0.2em] rounded-xl opacity-0 group-hover/nav:opacity-100 translate-x-[-10px] group-hover/nav:translate-x-0 transition-all pointer-events-none whitespace-nowrap z-50 shadow-2xl">
                      {item.name}
                      <div className="absolute top-1/2 -left-1 -translate-y-1/2 border-4 border-transparent border-r-slate-900" />
                    </div>
                  )}
                </Link>
              </motion.div>
            );
          })}
        </nav>
        
        {/* Footer */}
        <div className="mt-auto p-4 space-y-4">
          <div className={cn(
            "flex items-center bg-white dark:bg-slate-800 rounded-3xl border border-border transition-all shadow-sm overflow-hidden",
            isCollapsed ? "p-2 justify-center" : "p-3 justify-between"
          )}>
             <div className="flex items-center gap-3 shrink-0">
               <Avatar className="h-10 w-10 border-2 border-background shadow-sm">
                 <AvatarImage src={user?.avatar || ""} />
                 <AvatarFallback className="bg-muted text-foreground text-xs font-bold">
                   {user?.username?.[0]?.toUpperCase() || "A"}
                 </AvatarFallback>
               </Avatar>
               {!isCollapsed && (
                 <div className="flex flex-col min-w-0">
                   <span className="text-xs font-bold truncate">{user?.username || "Admin"}</span>
                   <div className="flex items-center gap-1.5">
                      <div className="w-1.5 h-1.5 rounded-full bg-emerald-500 animate-pulse" />
                      <span className="text-[9px] text-muted-foreground uppercase tracking-widest font-bold">Online</span>
                   </div>
                 </div>
               )}
             </div>
             
             {!isCollapsed && (
               <DropdownMenu>
                 <DropdownMenuTrigger asChild>
                   <Button variant="ghost" size="icon" className="h-8 w-8 rounded-full hover:bg-muted">
                      <MoreHorizontal className="w-4 h-4 text-muted-foreground" />
                   </Button>
                 </DropdownMenuTrigger>
                 <DropdownMenuContent align="end" side="top" className="w-56 bg-white dark:bg-slate-900 border-border shadow-2xl rounded-2xl p-2">
                   <DropdownMenuItem className="cursor-pointer rounded-xl py-3" onClick={() => router.push("/admin/settings")}>
                     <Settings className="mr-2 h-4 w-4" /> {t('settings')}
                   </DropdownMenuItem>
                   <DropdownMenuSeparator className="bg-border" />
                   <DropdownMenuItem onClick={handleLogout} className="text-accent focus:text-white focus:bg-accent cursor-pointer rounded-xl py-3">
                     <LogOut className="mr-2 h-4 w-4" /> {t('logout_text') || 'Logout'}
                   </DropdownMenuItem>
                 </DropdownMenuContent>
               </DropdownMenu>
             )}
          </div>
        </div>
      </motion.aside>

      {/* Main Viewport */}
      <div className="flex flex-1 flex-col relative overflow-hidden bg-transparent">
        {/* Background Grid Layer */}
        <div className="absolute inset-0 z-0 pointer-events-none opacity-[0.15]">
          <div className="absolute inset-0 bg-[linear-gradient(to_right,var(--border)_1px,transparent_1px),linear-gradient(to_bottom,var(--border)_1px,transparent_1px)] bg-[size:60px_60px]" />
        </div>

        {/* IMPORTANT: This is the scrollable container. Refactored to fix scrolling. */}
        <main className="flex-1 overflow-y-auto overflow-x-hidden custom-scrollbar relative z-10 flex flex-col">
          <div className="p-6 md:p-12 max-w-7xl w-full mx-auto flex-1">
            <AnimatePresence mode="wait">
              <motion.div
                key={pathname}
                initial={{ opacity: 0, y: 15 }}
                animate={{ opacity: 1, y: 0 }}
                exit={{ opacity: 0, y: -15 }}
                transition={{ duration: 0.4, ease: "easeOut" }}
                className="min-h-full"
              >
                {children}
              </motion.div>
            </AnimatePresence>
          </div>
          
          {/* Subtle Page Footer inside scroll view */}
          <footer className="p-12 text-center opacity-20 mt-auto">
             <p className="text-[10px] font-bold uppercase tracking-[0.5em]">Kaldalis Headless Content System</p>
          </footer>
        </main>
      </div>

      <style jsx global>{`
        .custom-scrollbar::-webkit-scrollbar {
          width: 5px;
        }
        .custom-scrollbar::-webkit-scrollbar-track {
          background: transparent;
        }
        .custom-scrollbar::-webkit-scrollbar-thumb {
          background: var(--border);
          border-radius: 10px;
        }
        .custom-scrollbar::-webkit-scrollbar-thumb:hover {
          background: var(--accent);
        }
      `}</style>
    </div>
  );
}
