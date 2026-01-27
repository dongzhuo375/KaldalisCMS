"use client";

import { useState } from "react";
import { useAuthStore } from "@/store/useAuthStore";
import { useLocale, useTranslations } from 'next-intl';
import { useRouter, usePathname, Link } from '@/i18n/routing';
import { Button } from "@/components/ui/button";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { LogOut, User, Settings, Globe, ChevronUp } from "lucide-react";

export default function FloatingMenu() {
  const { user, isLoggedIn, logout } = useAuthStore();
  const t = useTranslations();
  const locale = useLocale();
  const router = useRouter();
  const pathname = usePathname();
  const [isOpen, setIsOpen] = useState(false);

  const languages = [
    {code: 'zh-CN', name: '简体中文', flag: '🇨🇳'},
    {code: 'en', name: 'English', flag: '🇺🇸'},
  ];

  const switchLanguage = (newLocale: string) => {
    const currentPath = pathname.replace(/^\/[^\/]*/, '');
    router.push(`/${newLocale}${currentPath}`);
  };

  const handleLogout = async () => {
    try {
      // 调用后端清除 Cookie
      await fetch('/api/users/logout', { method: 'POST' });
    } catch (e) {
      console.error("Logout error", e);
    }
    // 清除前端状态并刷新
    logout();
    window.location.href = "/login";
  };

  return (
    <div className="fixed bottom-6 right-6 z-50 flex flex-col items-end">
      {/* 展开的菜单 */}
      {isOpen && (
        <div className="mb-4 flex flex-col gap-2 animate-in slide-in-from-bottom-2 fade-in-0 duration-200">
          {/* 语言切换 */}
          <div className="bg-white rounded-lg shadow-lg border p-2 min-w-[140px]">
            <div className="flex items-center gap-2 px-3 py-1 text-sm font-medium text-muted-foreground border-b mb-1">
              <Globe className="h-4 w-4" />
              {t('navigation.language')}
            </div>
            {languages.map((language) => (
              <div
                key={language.code}
                onClick={() => switchLanguage(language.code)}
                className={`flex items-center px-3 py-2 text-sm cursor-pointer rounded-md transition-colors hover:bg-accent ${
                  locale === language.code ? 'bg-accent' : ''
                }`}
              >
                <span className="mr-2">{language.flag}</span>
                {language.name}
              </div>
            ))}
          </div>
          
          {/* 如果已登录，显示用户选项 */}
          {isLoggedIn && user && (
            <div className="bg-white rounded-lg shadow-lg border p-2 min-w-[200px]">
              {/* 用户信息 */}
              <div className="px-3 py-2 text-sm border-b">
                <div className="font-medium">{user.username}</div>
                <div className="text-muted-foreground text-xs">{user.email || user.role}</div>
              </div>
              
              {/* 管理员入口 */}
              {(user.role === 'admin' || user.role === 'super_admin') && (
                <div 
                  onClick={() => router.push('/admin/dashboard')}
                  className="flex items-center gap-2 px-3 py-2 text-sm cursor-pointer rounded-md transition-colors hover:bg-accent"
                >
                  <Settings className="h-4 w-4" />
                  {t('navigation.enter_admin')}
                </div>
              )}
              
              {/* 个人资料 */}
              <div className="flex items-center gap-2 px-3 py-2 text-sm cursor-pointer rounded-md transition-colors hover:bg-accent">
                <User className="mr-2 h-4 w-4" />
                {t('navigation.personal_profile')}
              </div>
              
              <div className="border-t my-1"></div>
              
              {/* 退出登录 */}
              <div 
                onClick={handleLogout} 
                className="flex items-center gap-2 px-3 py-2 text-sm cursor-pointer rounded-md transition-colors hover:bg-accent text-red-600 hover:text-red-600"
              >
                <LogOut className="mr-2 h-4 w-4" />
                {t('navigation.logout_text')}
              </div>
            </div>
          )}
          
          {/* 如果未登录，显示登录注册 */}
          {!isLoggedIn && (
            <div className="bg-white rounded-lg shadow-lg border p-2 flex gap-2">
              <Button variant="ghost" size="sm" asChild>
                <Link href="/login">{t('common.login')}</Link>
              </Button>
              <Button size="sm" asChild>
                <Link href="/register">{t('common.register')}</Link>
              </Button>
            </div>
          )}
        </div>
      )}

      {/* 浮动按钮 */}
      <Button
        onClick={() => setIsOpen(!isOpen)}
        size="lg"
        className="h-14 w-14 rounded-full shadow-lg bg-blue-600 hover:bg-blue-700 text-white border-0 transition-all duration-200"
      >
        {isLoggedIn && user ? (
          <Avatar className="h-8 w-8">
            <AvatarImage src={user.avatar} alt={user.username} />
            <AvatarFallback className="bg-white text-blue-600 text-sm font-medium">
              {user.username?.[0]?.toUpperCase() || "U"}
            </AvatarFallback>
          </Avatar>
        ) : (
          <ChevronUp 
            className={`h-6 w-6 transition-transform duration-200 ${
              isOpen ? 'rotate-180' : ''
            }`}
          />
        )}
      </Button>
    </div>
  );
}