"use client";

import React, { useState, useEffect } from "react";
import { Link, usePathname, useRouter } from "@/i18n/routing";
import { useAuthStore } from "@/store/useAuthStore";
import api from "@/lib/api";
import { useTranslations } from "next-intl";
import { cn } from "@/lib/utils";
import {
  LayoutDashboard,
  FileText,
  Image,
  Users,
  BarChart3,
  Settings,
  LogOut,
  PanelLeftClose,
  PanelLeft,
} from "lucide-react";

import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
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
  const locale =
    rawLocale && rawLocale !== "undefined" ? rawLocale : "zh-CN";

  const t = useTranslations("admin");
  const { user, isLoggedIn, logout } = useAuthStore();

  const [collapsed, setCollapsed] = useState(false);

  useEffect(() => {
    const saved = localStorage.getItem("sidebar-collapsed");
    if (saved !== null) setCollapsed(saved === "true");

    const handleKeyDown = (e: KeyboardEvent) => {
      if ((e.metaKey || e.ctrlKey) && e.key === "b") {
        e.preventDefault();
        toggle();
      }
    };
    window.addEventListener("keydown", handleKeyDown);
    return () => window.removeEventListener("keydown", handleKeyDown);
  }, []);

  const toggle = () => {
    setCollapsed((prev) => {
      const next = !prev;
      localStorage.setItem("sidebar-collapsed", String(next));
      return next;
    });
  };

  React.useEffect(() => {
    if (!isLoggedIn) {
      router.replace("/login");
    }
  }, [isLoggedIn, router]);

  const handleLogout = async () => {
    try {
      await api.post("/users/logout");
    } catch (error) {
      console.error("Logout failed:", error);
    }
    logout();
    router.push("/login");
  };

  if (!isLoggedIn) return null;

  const navItems = [
    { name: t("dashboard"), href: "/admin/dashboard", icon: LayoutDashboard },
    { name: t("content"), href: "/admin/posts", icon: FileText },
    { name: t("media"), href: "/admin/media", icon: Image },
    { name: t("users"), href: "/admin/users", icon: Users },
    { name: t("analytics"), href: "/admin/analytics", icon: BarChart3 },
  ];

  return (
    <div className="flex h-screen w-full bg-transparent text-foreground overflow-hidden">
      {/* Sidebar */}
      <aside
        className={cn(
          "hidden md:flex flex-col border-r border-border bg-background shrink-0 overflow-hidden transition-[width] duration-300 ease-in-out",
          collapsed ? "w-16" : "w-[240px]"
        )}
      >
        {/* Logo + collapse toggle */}
        <div className="h-16 flex items-center justify-between px-4 border-b border-border shrink-0">
          <Link
            href="/admin/dashboard"
            className="flex items-center gap-3 min-w-0"
          >
            <div className="w-8 h-8 rounded-lg bg-accent text-accent-foreground flex items-center justify-center font-serif text-lg font-bold shrink-0">
              K
            </div>
            {!collapsed && (
              <span className="font-serif text-base font-bold tracking-tight truncate">
                Kaldalis
              </span>
            )}
          </Link>
          {!collapsed && (
            <button
              onClick={toggle}
              className="p-1.5 rounded-md text-muted-foreground hover:bg-muted hover:text-foreground transition-colors shrink-0"
              title="Collapse sidebar (Cmd+B)"
            >
              <PanelLeftClose className="w-4 h-4" />
            </button>
          )}
        </div>

        {/* Expand button when collapsed */}
        {collapsed && (
          <div className="px-3 pt-3 shrink-0">
            <button
              onClick={toggle}
              className="w-full flex items-center justify-center p-2 rounded-lg text-muted-foreground hover:bg-muted hover:text-foreground transition-colors"
              title="Expand sidebar (Cmd+B)"
            >
              <PanelLeft className="w-4 h-4" />
            </button>
          </div>
        )}

        {/* Nav */}
        <nav className="flex-1 py-3 px-3 space-y-0.5 overflow-y-auto overflow-x-hidden">
          {navItems.map((item) => {
            const isActive =
              pathname === item.href ||
              pathname.startsWith(`${item.href}/`);
            return (
              <Link
                key={item.href}
                href={item.href}
                className={cn(
                  "flex items-center rounded-lg text-sm font-medium transition-colors relative group",
                  collapsed
                    ? "justify-center p-2.5"
                    : "gap-3 px-3 py-2.5",
                  isActive
                    ? "bg-accent/10 text-accent font-semibold"
                    : "text-muted-foreground hover:bg-muted hover:text-foreground"
                )}
              >
                <item.icon className="w-[18px] h-[18px] shrink-0" />
                {!collapsed && <span className="truncate">{item.name}</span>}

                {collapsed && (
                  <div className="absolute left-full ml-2 px-2.5 py-1.5 bg-foreground text-background text-xs font-medium rounded-md opacity-0 group-hover:opacity-100 pointer-events-none transition-opacity whitespace-nowrap z-50">
                    {item.name}
                  </div>
                )}
              </Link>
            );
          })}
        </nav>

        {/* User at the very bottom */}
        <div className="p-3 border-t border-border shrink-0">
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <button
                className={cn(
                  "flex items-center w-full rounded-lg hover:bg-muted transition-colors text-left",
                  collapsed
                    ? "justify-center p-2"
                    : "gap-3 px-3 py-2"
                )}
              >
                <Avatar className="h-7 w-7 shrink-0">
                  <AvatarImage src={user?.avatar || ""} />
                  <AvatarFallback className="bg-muted text-foreground text-xs font-bold">
                    {user?.username?.[0]?.toUpperCase() || "A"}
                  </AvatarFallback>
                </Avatar>
                {!collapsed && (
                  <div className="flex flex-col min-w-0">
                    <span className="text-sm font-medium truncate">
                      {user?.username || "Admin"}
                    </span>
                    <span className="text-[10px] text-muted-foreground truncate">
                      {user?.role?.replace("_", " ") || "user"}
                    </span>
                  </div>
                )}
              </button>
            </DropdownMenuTrigger>
            <DropdownMenuContent
              align={collapsed ? "center" : "start"}
              side="top"
              className="w-48"
            >
              <DropdownMenuItem
                className="cursor-pointer"
                onClick={() => router.push("/admin/settings")}
              >
                <Settings className="mr-2 h-4 w-4" /> {t("settings")}
              </DropdownMenuItem>
              <DropdownMenuSeparator />
              <DropdownMenuItem
                onClick={handleLogout}
                className="text-destructive focus:text-destructive cursor-pointer"
              >
                <LogOut className="mr-2 h-4 w-4" />{" "}
                {t("logout_text") || "Logout"}
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </aside>

      {/* Main */}
      <div className="flex flex-1 flex-col overflow-hidden">
        <main className="flex-1 overflow-y-auto">
          <div className="p-6 md:p-10 max-w-7xl w-full mx-auto">
            {children}
          </div>
        </main>
      </div>
    </div>
  );
}
