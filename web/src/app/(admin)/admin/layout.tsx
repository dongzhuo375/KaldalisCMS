// src/app/(admin)/admin/layout.tsx
import React from "react";

export default function AdminLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="flex h-screen w-full overflow-hidden bg-slate-50">
      {/* 1. 左侧侧边栏 (Sidebar) */}
      <aside className="hidden w-64 flex-col border-r bg-white md:flex">
        <div className="flex h-16 items-center border-b px-6 text-lg font-bold tracking-tight text-slate-900">
          Kaldalis CMS
        </div>
        <nav className="flex-1 space-y-1 p-4">
          {/* 这里的链接之后可以用 Link 组件替换 */}
          <a href="/admin/dashboard" className="flex items-center rounded-lg bg-slate-100 px-3 py-2 text-sm font-medium text-slate-900">
            仪表盘
          </a>
          <a href="/admin/posts" className="flex items-center rounded-lg px-3 py-2 text-sm font-medium text-slate-600 hover:bg-slate-50 hover:text-slate-900">
            文章管理
          </a>
          <a href="/admin/themes" className="flex items-center rounded-lg px-3 py-2 text-sm font-medium text-slate-600 hover:bg-slate-50 hover:text-slate-900">
            主题设置
          </a>
        </nav>
        <div className="border-t p-4">
          <div className="text-xs text-slate-400">v1.0.0-alpha</div>
        </div>
      </aside>

      {/* 2. 右侧主内容区 */}
      <div className="flex flex-1 flex-col">
        {/* 右侧顶部栏 (Header) */}
        <header className="flex h-16 items-center justify-between border-b bg-white px-8">
          <div className="text-sm font-medium text-slate-500">
            后台管理 / <span className="text-slate-900">仪表盘</span>
          </div>
          <div className="flex items-center gap-4">
            <div className="h-8 w-8 rounded-full bg-slate-200" /> {/* 头像占位 */}
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
