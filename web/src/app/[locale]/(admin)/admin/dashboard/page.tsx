// src/app/(admin)/admin/dashboard/page.tsx
export default function DashboardPage() {
  return (
    <div className="space-y-6">
      <div className="grid gap-4 md:grid-cols-3">
        <div className="rounded-xl border bg-white p-6 shadow-sm">
          <div className="text-sm text-slate-500">总文章数</div>
          <div className="text-2xl font-bold">128</div>
        </div>
        <div className="rounded-xl border bg-white p-6 shadow-sm">
          <div className="text-sm text-slate-500">今日访问</div>
          <div className="text-2xl font-bold">1,024</div>
        </div>
        <div className="rounded-xl border bg-white p-6 shadow-sm">
          <div className="text-sm text-slate-500">系统状态</div>
          <div className="text-2xl font-bold text-green-500">运行中</div>
        </div>
      </div>
      
      <div className="rounded-xl border bg-white p-12 text-center text-slate-400">
        这里是图表展示区域（待开发）
      </div>
    </div>
  );
}
