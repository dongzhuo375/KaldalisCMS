"use client"; // 1. 必须变身客户端组件才能交互

import React from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useAuthStore } from "@/store/useAuthStore";
import api from "@/lib/api";

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
    <div className="flex h-screen w-full overflow-hidden bg-slate-50">
      {/* 左侧侧边栏 (Sidebar) */}
      <aside className="hidden w-64 flex-col border-r bg-white md:flex">
        <div className="flex h-16 items-center border-b px-6 text-lg font-bold tracking-tight text-slate-900">
          Kaldalis CMS
        </div>
        <nav className="flex-1 space-y-1 p-4">
          <Link href="/admin/dashboard" className="flex items-center rounded-lg bg-slate-100 px-3 py-2 text-sm font-medium text-slate-900">
            仪表盘
          </Link>
          <Link href="/admin/posts" className="flex items-center rounded-lg px-3 py-2 text-sm font-medium text-slate-600 hover:bg-slate-50 hover:text-slate-900">
            文章管理
          </Link>
          <Link href="/admin/themes" className="flex items-center rounded-lg px-3 py-2 text-sm font-medium text-slate-600 hover:bg-slate-50 hover:text-slate-900">
            主题设置
          </Link>
        </nav>
        <div className="border-t p-4">
          <div className="text-xs text-slate-400">v1.0.0-alpha</div>
        </div>
      </aside>

      {/* 右侧主内容区 */}
      <div className="flex flex-1 flex-col">
        {/* 右侧顶部栏 (Header) */}
        <header className="flex h-16 items-center justify-between border-b bg-white px-8">
          <div className="text-sm font-medium text-slate-500">
            后台管理 / <span className="text-slate-900">仪表盘</span>
          </div>

          <div className="flex items-center gap-4">
            {/* 4. 用户信息与下拉菜单 */}
            <div className="flex items-center gap-2">
              <span className="text-sm text-slate-700 hidden sm:inline-block">
                {user?.username || "管理员"}
              </span>
              
              <DropdownMenu>
                <DropdownMenuTrigger className="focus:outline-none">
                  <Avatar className="cursor-pointer hover:opacity-80 transition-opacity">
                    {/* 如果没有头像，显示名字首字母 */}
                    <AvatarImage src={user?.avatar || ""} />
                    <AvatarFallback className="bg-slate-900 text-white">
                      {user?.username?.[0]?.toUpperCase() || "A"}
                    </AvatarFallback>
                  </Avatar>
                </DropdownMenuTrigger>
                
                <DropdownMenuContent align="end" className="w-56">
                  <DropdownMenuLabel>我的账户</DropdownMenuLabel>
                  <DropdownMenuSeparator />
                  <DropdownMenuItem onClick={() => router.push('/')}>
                    返回前台首页
                  </DropdownMenuItem>
                  <DropdownMenuItem onClick={() => alert("个人设置开发中...")}>
                    个人设置
                  </DropdownMenuItem>
                  <DropdownMenuSeparator />
                  <DropdownMenuItem 
                    onClick={handleLogout} 
                    className="text-red-600 focus:text-red-600 focus:bg-red-50 cursor-pointer"
                  >
                    退出登录
                  </DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
            </div>
          </div>
        </header>

        {/* 具体的页面内容 */}
        <main className="flex-1 overflow-y-auto p-8">
          {children}
        </main>
      </div>
    </div>
  );
}
