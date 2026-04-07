"use client";

import { Link } from '@/i18n/routing';
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
import { ThemeToggle } from "@/components/theme-toggle";
import LanguageSwitcher from "@/components/LanguageSwitcher";

export function SiteHeader() {
  const { user, isLoggedIn, logout } = useAuthStore();
  const t = useTranslations();

  const handleLogout = async () => {
    try {
      await api.post("/users/logout");
    } catch (e) {
      console.error("Logout error", e);
    }
    logout();
    window.location.href = "/login";
  };

  return (
    <header className="sticky top-0 z-50 w-full bg-transparent">
      <div className="container mx-auto max-w-7xl flex h-16 items-center justify-between px-6 md:px-12 lg:px-20">
        {/* Logo */}
        <Link href="/" className="text-xl font-serif font-medium text-foreground hover:text-accent transition-colors">
          {t('common.app_name')}
        </Link>

        {/* Navigation - Hidden on mobile */}
        <nav className="hidden md:flex items-center gap-8 text-sm font-medium text-muted-foreground">
          <Link href="/" className="hover:text-foreground transition-colors">{t('navigation.home')}</Link>
          <Link href="/posts" className="hover:text-foreground transition-colors">{t('navigation.posts_list')}</Link>
          <Link href="/about" className="hover:text-foreground transition-colors">{t('navigation.about_us')}</Link>
        </nav>

        {/* Right side */}
        <div className="flex items-center gap-3">
          <ThemeToggle />
          <LanguageSwitcher />

          {isLoggedIn && user ? (
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="ghost" className="relative h-9 w-9 rounded-full">
                  <Avatar className="h-9 w-9">
                    <AvatarImage src={user.avatar} alt={user.username} />
                    <AvatarFallback className="bg-accent/10 text-accent font-medium">
                      {user.username?.[0]?.toUpperCase() || "U"}
                    </AvatarFallback>
                  </Avatar>
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent className="w-56" align="end">
                <div className="flex items-center justify-start gap-2 p-2">
                  <div className="flex flex-col space-y-1 leading-none">
                    <p className="font-medium">{user.username}</p>
                    <p className="text-xs text-muted-foreground">{user.email || user.role}</p>
                  </div>
                </div>
                <DropdownMenuSeparator />
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
            <Link href="/login">
              <Button size="sm" className="rounded-full px-5 font-medium">
                {t('common.login')}
              </Button>
            </Link>
          )}
        </div>
      </div>
    </header>
  );
}
