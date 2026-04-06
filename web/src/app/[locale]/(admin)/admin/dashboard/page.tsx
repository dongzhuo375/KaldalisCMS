"use client";

import { 
  FileText, 
  Activity, 
  Clock, 
  Loader2,
  Database,
  ShieldCheck,
  Cpu,
  Zap,
  Plus,
  Terminal,
  ChevronRight,
  Code
} from "lucide-react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { cn, getImageUrl } from "@/lib/utils";
import { useTranslations } from 'next-intl';
import { useSystemStatus, useReadyz } from "@/services/system-service";
import { useAdminPosts } from "@/services/post-service";
import { PostStatus } from "@/lib/types";
import { useMedia } from "@/services/media-service";
import { motion, AnimatePresence } from "framer-motion";
import { Link } from "@/i18n/routing";
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";

export default function DashboardPage() {
  const t = useTranslations('admin');
  const { data: status, isLoading: statusLoading } = useSystemStatus();
  const { data: health, isLoading: healthLoading } = useReadyz();
  const { data: postsData = [], isLoading: postsLoading } = useAdminPosts();
  const { data: mediaData, isLoading: mediaLoading } = useMedia({ page_size: 1 });

  const posts = Array.isArray(postsData) ? postsData : [];
  const isLoading = statusLoading || healthLoading || postsLoading || mediaLoading;

  if (isLoading) {
    return (
      <div className="h-full space-y-10 pb-20">
        <header className="flex justify-between items-end gap-6">
          <div className="space-y-2">
            <Skeleton className="h-12 w-64" />
            <Skeleton className="h-4 w-48" />
          </div>
          <Skeleton className="h-12 w-40 rounded-full" />
        </header>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          {[1, 2, 3, 4].map((i) => (
            <Skeleton key={i} className="h-32 w-full rounded-3xl" />
          ))}
        </div>
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          <Skeleton className="lg:col-span-2 h-[400px] rounded-3xl" />
          <div className="space-y-8">
            <Skeleton className="h-[200px] rounded-3xl" />
            <Skeleton className="h-[200px] rounded-3xl" />
          </div>
        </div>
      </div>
    );
  }

  const dbStatus = health?.checks?.database?.status === 'ok';
  const engineMode = health?.mode || 'unknown';

  const stats = [
    {
      title: "Database Link",
      value: dbStatus ? "Connected" : "Disconnected",
      status: dbStatus ? "success" : "error",
      icon: Database,
      detail: dbStatus ? "PostgreSQL Active" : "Check logs"
    },
    {
      title: "API Engine",
      value: health?.status === 'ok' ? "Healthy" : "Degraded",
      status: health?.status === 'ok' ? "success" : "warning",
      icon: Zap,
      detail: `Mode: ${engineMode.toUpperCase()}`
    },
    {
      title: "Content Store",
      value: posts.length.toString(),
      status: "neutral",
      icon: FileText,
      detail: "Managed Articles"
    },
    {
      title: "Media Assets",
      value: mediaData?.total?.toString() || "0",
      status: "neutral",
      icon: Activity,
      detail: "Uploaded Files"
    }
  ];

  return (
    <div className="h-full space-y-10 custom-scrollbar pb-20">
      
      {/* Header */}
      <header className="flex flex-col md:flex-row md:items-end justify-between gap-6">
        <motion.div 
          initial={{ opacity: 0, y: 10 }}
          animate={{ opacity: 1, y: 0 }}
          className="space-y-2"
        >
          <h1 className="text-4xl md:text-5xl font-serif font-medium tracking-tight text-foreground">
            {status?.site_name || "Console"}
          </h1>
          <p className="text-muted-foreground font-medium flex items-center gap-2">
            <span className={cn(
              "w-2 h-2 rounded-full animate-pulse",
              dbStatus ? "bg-emerald-500" : "bg-rose-500"
            )} />
            {dbStatus ? "System is fully operational" : "System issue detected"} • {status?.version || 'v2.4.0'}
          </p>
        </motion.div>

        <Link href="/admin/posts/new">
          <Button className="rounded-full bg-primary text-primary-foreground h-12 px-6 font-bold shadow-xl shadow-primary/10 hover:scale-105 transition-transform">
            <Plus className="w-4 h-4 mr-2" /> New Article
          </Button>
        </Link>
      </header>

      {/* Stats Grid: Real Health Data */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        {stats.map((stat, index) => (
          <motion.div
            key={index}
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: index * 0.1 }}
          >
            <Card className="border-border bg-white/50 dark:bg-slate-900/50 backdrop-blur-sm shadow-none hover:border-accent/20 transition-colors group">
              <CardHeader className="flex flex-row items-center justify-between pb-2 space-y-0">
                <CardTitle className="text-[10px] font-bold uppercase tracking-[0.2em] text-muted-foreground group-hover:text-accent transition-colors">
                  {stat.title}
                </CardTitle>
                <stat.icon className={cn(
                  "h-4 w-4",
                  stat.status === 'success' ? "text-emerald-500" : 
                  stat.status === 'error' ? "text-rose-500" : 
                  "text-muted-foreground"
                )} strokeWidth={2.5} />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-serif font-medium">{stat.value}</div>
                <div className="text-[10px] font-bold text-muted-foreground/60 uppercase tracking-tighter mt-1 italic">
                  {stat.detail}
                </div>
              </CardContent>
            </Card>
          </motion.div>
        ))}
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        
        {/* Recent Content */}
        <Card className="lg:col-span-2 border-border shadow-none bg-transparent">
          <CardHeader className="px-0 flex flex-row items-center justify-between">
            <CardTitle className="text-2xl font-serif font-medium">Recent Activity</CardTitle>
            <Link href="/admin/posts" className="text-xs font-bold uppercase text-accent hover:underline">View All</Link>
          </CardHeader>
          <CardContent className="px-0 pt-2">
            <div className="space-y-1">
              {posts.length === 0 ? (
                <p className="text-muted-foreground text-sm py-10 text-center border border-dashed rounded-xl">No articles found.</p>
              ) : posts.map((item, i) => (
                <motion.div 
                  key={item.id}
                  initial={{ opacity: 0, x: -10 }}
                  animate={{ opacity: 1, x: 0 }}
                  transition={{ delay: 0.4 + (i * 0.05) }}
                  className="flex items-center justify-between p-4 rounded-xl hover:bg-muted/50 transition-colors group border border-transparent hover:border-border"
                >
                  <div className="flex items-center gap-4">
                    <div className="w-10 h-10 rounded-full bg-accent/5 overflow-hidden flex items-center justify-center text-accent group-hover:bg-accent group-hover:text-white transition-all">
                      {item.cover ? (
                        <img src={getImageUrl(item.cover)} alt="" className="w-full h-full object-cover" />
                      ) : (
                        <FileText className="w-5 h-5" />
                      )}
                    </div>
                    <div>
                      <h4 className="font-bold text-sm group-hover:text-accent transition-colors">{item.title}</h4>
                      <p className="text-[10px] text-muted-foreground uppercase tracking-widest font-bold">
                        {new Date(item.created_at).toLocaleDateString(undefined, { month: 'long', day: 'numeric' })}
                      </p>
                    </div>
                  </div>
                  <Badge variant="outline" className={cn(
                    "rounded-full border-0 px-3 py-1 text-[10px] font-bold uppercase",
                    item.status === PostStatus.PUBLISHED ? "bg-accent/10 text-accent" : "bg-muted text-muted-foreground"
                  )}>
                    {item.status === PostStatus.PUBLISHED ? 'Published' : 'Draft'}
                  </Badge>
                </motion.div>
              ))}
            </div>
          </CardContent>
        </Card>

        {/* System Telemetry & Logs */}
        <div className="space-y-8">
          <section className="space-y-4">
            <div className="flex items-center gap-2">
              <Terminal className="w-4 h-4 text-accent" />
              <h3 className="text-xl font-serif font-medium">Telemetry</h3>
            </div>
            
            <div className="p-6 rounded-2xl border border-border bg-slate-950 text-emerald-500 font-mono text-[10px] leading-relaxed overflow-hidden relative group">
               <div className="absolute top-2 right-2 flex gap-1 opacity-30 group-hover:opacity-100 transition-opacity">
                  <div className="w-2 h-2 rounded-full bg-rose-500" />
                  <div className="w-2 h-2 rounded-full bg-yellow-500" />
                  <div className="w-2 h-2 rounded-full bg-emerald-500" />
               </div>
               <pre className="custom-scrollbar overflow-x-auto">
                 {JSON.stringify(health, null, 2)}
               </pre>
            </div>
          </section>

          <section className="space-y-4">
            <div className="flex items-center justify-between">
              <h3 className="text-xl font-serif font-medium">Session</h3>
              <div className="w-2 h-2 rounded-full bg-emerald-500" />
            </div>
            <div className="space-y-4 p-6 rounded-2xl border border-border bg-white/30 dark:bg-slate-900/30 backdrop-blur-sm">
               <div className="flex justify-between items-center text-[10px] font-bold uppercase tracking-widest">
                  <span className="text-muted-foreground">Process Identity</span>
                  <span className="text-foreground">Node_01</span>
               </div>
               <div className="flex justify-between items-center text-[10px] font-bold uppercase tracking-widest">
                  <span className="text-muted-foreground">Up Since</span>
                  <span className="text-foreground">24 Mar 2026</span>
               </div>
               <div className="flex justify-between items-center text-[10px] font-bold uppercase tracking-widest">
                  <span className="text-muted-foreground">Region</span>
                  <span className="text-foreground">Global_Edge</span>
               </div>
            </div>
          </section>
        </div>
      </div>
    </div>
  );
}
