"use client";

import { 
  Eye, 
  FileText, 
  Users, 
  Server, 
  Code, 
  PenTool, 
  Megaphone, 
  Shield, 
  Activity, 
  Clock, 
  Terminal,
  Search,
  Loader2
} from "lucide-react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { cn } from "@/lib/utils";
import { useTranslations } from 'next-intl';
import { useSystemStatus } from "@/services/system-service";
import { usePosts } from "@/services/post-service";

export default function DashboardPage() {
  const t = useTranslations('admin');
  const { data: status, isLoading: statusLoading } = useSystemStatus();
  const { data: posts = [], isLoading: postsLoading } = usePosts({ limit: 5 });

  if (statusLoading || postsLoading) {
    return (
      <div className="h-full flex items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-indigo-500" />
      </div>
    );
  }

  // Derived stats
  const stats = [
    {
      title: t('total_views'),
      value: "124,592",
      change: "+12.5%",
      trend: "up",
      icon: Eye,
      color: "text-emerald-400",
      gradient: "from-emerald-500/20 to-transparent"
    },
    {
      title: t('articles_published'),
      value: posts.length.toString(),
      change: `+${posts.filter(p => p.status === 1).length} ${t('stable')}`,
      trend: "up",
      icon: FileText,
      color: "text-blue-400",
      gradient: "from-blue-500/20 to-transparent"
    },
    {
      title: t('active_users'),
      value: "1.2k",
      change: "Stable",
      trend: "neutral",
      icon: Users,
      color: "text-rose-400",
      gradient: "from-rose-500/20 to-transparent"
    },
    {
      title: t('server_load'),
      value: "28%",
      change: t('stable'),
      trend: "neutral",
      icon: Server,
      color: "text-purple-400",
      gradient: "from-purple-500/20 to-transparent"
    }
  ];

  const activityLog = [
    {
      id: 1,
      message: "System initialized",
      time: status?.version || "v1.0.0",
      type: "success"
    },
    {
      id: 2,
      message: "Admin session started",
      time: "Just now",
      type: "info"
    }
  ];

  return (
    <div className="h-full overflow-y-auto bg-slate-950 text-slate-200 p-6 rounded-xl font-sans space-y-8 custom-scrollbar">
      
      {/* 顶部 CLI 风格 Header */}
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4 border-b border-slate-800 pb-6">
        <div className="flex items-center gap-3 font-mono text-sm text-indigo-400">
          <Terminal className="w-5 h-5" />
          <span>root @ kaldalis-cms : ~/dashboard</span>
          <span className="animate-pulse block w-2 h-4 bg-slate-500 ml-1"></span>
        </div>
        
        <div className="flex items-center gap-4">
          <div className="relative hidden md:block">
            <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-slate-500" />
            <input 
              type="text" 
              placeholder="Search telemetry..." 
              className="bg-slate-900 border border-slate-800 rounded-full pl-9 pr-4 py-1.5 text-sm focus:outline-none focus:border-slate-600 w-64"
            />
          </div>
          <div className="flex gap-2">
            <div className="w-2 h-2 rounded-full bg-indigo-500"></div>
            <div className="w-2 h-2 rounded-full bg-purple-500"></div>
            <div className="w-2 h-2 rounded-full bg-cyan-500"></div>
          </div>
        </div>
      </div>

      {/* Hero 区域 */}
      <div className="space-y-2">
        <h1 className="text-3xl font-bold tracking-tight text-white">{status?.site_name || t('system_overview')}</h1>
        <div className="flex items-center gap-2 text-slate-400 text-sm">
          <span>{t('current_build')} {status?.version || 'v2.4.0'}</span>
          <span className="text-slate-700">•</span>
          <span className="flex items-center gap-1.5 text-emerald-400">
            <span className="relative flex h-2 w-2">
              <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
              <span className="relative inline-flex rounded-full h-2 w-2 bg-emerald-500"></span>
            </span>
            {t('all_systems_operational')}
          </span>
        </div>
      </div>

      {/* 统计卡片 Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        {stats.map((stat, index) => (
          <Card key={index} className="bg-slate-900 border-slate-800 relative overflow-hidden group hover:border-slate-700 transition-colors">
            <div className={`absolute bottom-0 left-0 right-0 h-0.5 bg-gradient-to-r ${stat.gradient}`}></div>
            <CardHeader className="flex flex-row items-center justify-between pb-2">
              <CardTitle className="text-xs font-bold uppercase tracking-widest text-slate-500">
                {stat.title}
              </CardTitle>
              <stat.icon className={`h-6 w-6 opacity-20 ${stat.color} group-hover:opacity-40 transition-opacity`} />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold text-white tracking-tight">{stat.value}</div>
              <p className={`text-[10px] mt-1 font-mono ${
                stat.trend === 'up' ? 'text-emerald-400' : 
                stat.trend === 'down' ? 'text-rose-400' : 'text-slate-400'
              }`}>
                {stat.trend === 'up' ? '↗' : stat.trend === 'down' ? '↘' : '✓'} {stat.change}
              </p>
            </CardContent>
          </Card>
        ))}
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        
        {/* 左侧：最近内容 (占据 2/3) */}
        <Card className="lg:col-span-2 bg-slate-900 border-slate-800">
          <CardHeader className="flex flex-row items-center justify-between">
            <div className="flex items-center gap-2 border-l-2 border-indigo-500 pl-3">
              <CardTitle className="text-lg text-white">{t('recent_content')}</CardTitle>
            </div>
            <div className="hidden sm:block px-3 py-1 bg-slate-950 rounded border border-slate-800 font-mono text-[10px] text-indigo-400">
              ls -la ./content
            </div>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              <div className="grid grid-cols-12 text-[10px] font-bold text-slate-500 uppercase tracking-[0.2em] pb-2 border-b border-slate-800">
                <div className="col-span-6 pl-2">{t('title')}</div>
                <div className="col-span-3">{t('author')}</div>
                <div className="col-span-2">{t('status')}</div>
                <div className="col-span-1 text-right">{t('date')}</div>
              </div>
              
              {posts.length === 0 ? (
                <div className="py-8 text-center text-slate-500 text-sm">No content found.</div>
              ) : posts.map((item) => (
                <div key={item.id} className="grid grid-cols-12 items-center py-3 hover:bg-white/5 rounded-lg transition-colors group">
                  <div className="col-span-6 flex items-center gap-3 pl-2">
                    <div className="p-2 rounded bg-slate-800 text-slate-500 group-hover:text-white transition-colors">
                      <FileText className="w-4 h-4" />
                    </div>
                    <span className="font-medium text-slate-200 group-hover:text-indigo-400 transition-colors truncate">{item.title}</span>
                  </div>
                  <div className="col-span-3 text-sm text-slate-400">{item.author?.username || 'Admin'}</div>
                  <div className="col-span-2">
                    <Badge variant="outline" className={cn(
                      "border-0 px-2 py-0.5 text-[10px] font-bold uppercase",
                      item.status === 1 && "bg-emerald-500/10 text-emerald-400",
                      item.status === 0 && "bg-yellow-500/10 text-yellow-400",
                      item.status === 2 && "bg-slate-500/10 text-slate-400"
                    )}>
                      {item.status === 1 ? 'Published' : item.status === 0 ? 'Draft' : 'Archived'}
                    </Badge>
                  </div>
                  <div className="col-span-1 text-right text-[10px] font-mono text-slate-500">
                    {new Date(item.created_at).toLocaleDateString(undefined, { month: 'short', day: 'numeric' })}
                  </div>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>

        {/* 右侧 Sidebar (占据 1/3) */}
        <div className="space-y-6">
          
          {/* 系统资源 */}
          <Card className="bg-slate-900 border-slate-800">
            <CardHeader>
              <CardTitle className="flex items-center gap-2 text-white text-sm font-bold uppercase tracking-widest">
                <Activity className="w-4 h-4 text-emerald-400" />
                {t('system_resources')}
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-6">
              <ResourceBar label="CPU_USAGE" value={14} color="bg-emerald-500" />
              <ResourceBar label="RAM_ALLOC" value={42} color="bg-indigo-500" />
              <ResourceBar label="STORAGE_IO" value={8} color="bg-cyan-500" />
            </CardContent>
          </Card>

          {/* 活动日志 */}
          <Card className="bg-slate-900 border-slate-800">
            <CardHeader>
              <CardTitle className="flex items-center gap-2 text-white text-sm font-bold uppercase tracking-widest">
                <Clock className="w-4 h-4 text-indigo-400" />
                {t('activity_log')}
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="relative border-l border-slate-800 ml-2 space-y-6">
                {activityLog.map((log) => (
                  <div key={log.id} className="ml-6 relative">
                    <div className={cn(
                      "absolute -left-[31px] w-2.5 h-2.5 rounded-full border-2 border-slate-950",
                      log.type === 'info' && "bg-indigo-500",
                      log.type === 'success' && "bg-emerald-500",
                      log.type === 'neutral' && "bg-slate-500",
                      log.type === 'error' && "bg-rose-500"
                    )}></div>
                    
                    <div>
                      <p className="text-sm text-slate-300">{log.message}</p>
                      <p className="text-[10px] text-slate-500 font-mono mt-0.5 uppercase">{log.time}</p>
                    </div>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}

// 简单的进度条组件
function ResourceBar({ label, value, color }: { label: string, value: number, color: string }) {
  return (
    <div className="space-y-2">
      <div className="flex justify-between text-[10px] font-mono text-slate-500 font-bold">
        <span>{label}</span>
        <span>{value}%</span>
      </div>
      <div className="h-1.5 bg-slate-800 rounded-full overflow-hidden">
        <div 
          className={`h-full ${color} transition-all duration-1000 ease-in-out`} 
          style={{ width: `${value}%` }}
        ></div>
      </div>
    </div>
  )
}