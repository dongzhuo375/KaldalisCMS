"use client";

import {Link} from '@/i18n/routing';
import { useAuthStore } from "@/store/useAuthStore";
import { useTranslations } from 'next-intl';
import api from "@/lib/api";
import { Button } from "@/components/ui/button";
import { 
  DropdownMenu, 
  DropdownMenuContent, 
  DropdownMenuItem, 
  DropdownMenuSeparator, 
  DropdownMenuTrigger 
} from "@/components/ui/dropdown-menu";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { LogOut, User, Settings } from "lucide-react";

export function SiteHeader() {
  const { user, isLoggedIn, logout } = useAuthStore();
  const t = useTranslations();

  const handleLogout = async () => {
    try {
      // 1. 调用后端清除 Cookie
      await api.post("/users/logout");
    } catch (e) {
      console.error("Logout error", e);
    }
    // 2. 清除前端状态并刷新
    logout();
    window.location.href = "/login";
  };

  return (
    <header className="sticky top-0 z-50 w-full border-b bg-white/95 dark:bg-slate-900/95 dark:border-slate-800 backdrop-blur supports-[backdrop-filter]:bg-white/60 transition-colors duration-300">
         <div className="container mx-auto max-w-7xl flex h-14 items-center justify-between px-4 sm:px-6 lg:px-8">
            {/* 左侧 Logo */}
        <div className="flex items-center gap-6">
          <Link href="/" className="flex items-center space-x-2">
            <span className="text-xl font-bold sm:inline-block dark:text-slate-200">
              {t('common.app_name')}
            </span>
          </Link>
          <nav className="hidden md:flex items-center gap-6 text-sm font-medium text-slate-600 dark:text-slate-400">
            <Link href="/" className="transition-colors hover:text-slate-900 dark:hover:text-slate-50">{t('navigation.home')}</Link>
            <Link href="/posts" className="transition-colors hover:text-slate-900 dark:hover:text-slate-50">{t('navigation.posts_list')}</Link>
            <Link href="/about" className="transition-colors hover:text-slate-900 dark:hover:text-slate-50">{t('navigation.about_us')}</Link>
          </nav>
        </div>

{/* 右侧 用户区域 */}
         <div className="flex items-center gap-2">
           {isLoggedIn && user ? (
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="ghost" className="relative h-8 w-8 rounded-full">
                  <Avatar className="h-8 w-8">
                    <AvatarImage src={user.avatar} alt={user.username} />
                    <AvatarFallback>{user.username?.[0]?.toUpperCase() || "U"}</AvatarFallback>  
                </Avatar>
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent className="w-56" align="end">
                <div className="flex items-center justify-start gap-2 p-2">
                  <div className="flex flex-col space-y-1 leading-none">
                    <p className="font-medium">{user.username}</p>
                    <p className="text-xs leading-none text-muted-foreground">{user.email || user.role}</p>
                  </div>
                </div>
                <DropdownMenuSeparator />
                {/* 如果是管理员，显示后台入口 */}
                 {(user.role === 'admin' || user.role === 'super_admin') && (
                   <DropdownMenuItem asChild>
                     <Link href="/admin/dashboard" className="cursor-pointer">
                       <Settings className="mr-2 h-4 w-4" /> {t('navigation.enter_admin')}
                     </Link>
                   </DropdownMenuItem>
                 )}
                 <DropdownMenuItem>
                   <User className="mr-2 h-4 w-4" /> {t('navigation.personal_profile')}
                 </DropdownMenuItem>
                 <DropdownMenuSeparator />
                 <DropdownMenuItem onClick={handleLogout} className="text-red-600 focus:text-red-600 cursor-pointer">
                   <LogOut className="mr-2 h-4 w-4" /> {t('navigation.logout_text')}
                 </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          ) : (
            <div className="flex items-center gap-2">
              <Link href="/login">
                <Button variant="ghost" size="sm">{t('common.login')}</Button>
              </Link>
              <Link href="/register">
                <Button size="sm">{t('common.register')}</Button>
              </Link>
            </div>
          )}
        </div>
      </div>
    </header>
  );
}
