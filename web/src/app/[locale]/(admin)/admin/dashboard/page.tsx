"use client";

import { 
  Eye, 
  FileText, 
  Users, 
  Server, 
  Activity, 
  Clock, 
  Loader2,
  TrendingUp,
  Plus
} from "lucide-react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { cn } from "@/lib/utils";
import { useTranslations } from 'next-intl';
import { useSystemStatus } from "@/services/system-service";
import { usePosts } from "@/services/post-service";
import { motion } from "framer-motion";
import { Link } from "@/i18n/routing";

import { Skeleton } from "@/components/ui/skeleton";

export default function DashboardPage() {
  const t = useTranslations('admin');
  const { data: status, isLoading: statusLoading } = useSystemStatus();
  const { data: posts = [], isLoading: postsLoading } = usePosts({ limit: 5, admin: true });

  if (statusLoading || postsLoading) {
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

  const stats = [
    {
      title: "Total Impressions",
      value: "124,592",
      change: "+12.5%",
      icon: Eye,
    },
    {
      title: "Articles Published",
      value: posts.filter(p => p.status === 1).length.toString(),
      change: `+${posts.length} total`,
      icon: FileText,
    },
    {
      title: "Community Growth",
      value: "1.2k",
      change: "Stable",
      icon: Users,
    },
    {
      title: "System Uptime",
      value: "99.9%",
      change: "Online",
      icon: Activity,
    }
  ];

  return (
    <div className="h-full overflow-y-auto space-y-10 custom-scrollbar pb-20">
      
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
            <span className="w-2 h-2 rounded-full bg-accent animate-pulse" />
            System is operational • {status?.version || 'v2.4.0'}
          </p>
        </motion.div>

        <Link href="/admin/posts/new">
          <Button className="rounded-full bg-primary text-primary-foreground h-12 px-6 font-bold shadow-xl shadow-primary/10 hover:scale-105 transition-transform">
            <Plus className="w-4 h-4 mr-2" /> New Article
          </Button>
        </Link>
      </header>

      {/* Stats Grid */}
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
                <stat.icon className="h-4 w-4 text-muted-foreground" strokeWidth={2.5} />
              </CardHeader>
              <CardContent>
                <div className="text-3xl font-serif font-medium">{stat.value}</div>
                <div className="flex items-center gap-1.5 mt-1">
                  <TrendingUp className="w-3 h-3 text-accent" />
                  <span className="text-[10px] font-bold text-accent">{stat.change}</span>
                </div>
              </CardContent>
            </Card>
          </motion.div>
        ))}
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        
        {/* Recent Content */}
        <Card className="lg:col-span-2 border-border shadow-none bg-transparent">
          <CardHeader className="px-0">
            <CardTitle className="text-2xl font-serif font-medium">Recent Articles</CardTitle>
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
                    <div className="w-10 h-10 rounded-full bg-accent/5 flex items-center justify-center text-accent group-hover:bg-accent group-hover:text-white transition-all">
                      <FileText className="w-5 h-5" />
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
                    item.status === 1 ? "bg-accent/10 text-accent" : "bg-muted text-muted-foreground"
                  )}>
                    {item.status === 1 ? 'Published' : 'Draft'}
                  </Badge>
                </motion.div>
              ))}
            </div>
          </CardContent>
        </Card>

        {/* System Activity */}
        <div className="space-y-8">
          <section className="space-y-4">
            <h3 className="text-xl font-serif font-medium">System Health</h3>
            <div className="space-y-4 p-6 rounded-2xl border border-border bg-white/30 dark:bg-slate-900/30 backdrop-blur-sm">
              <ResourceItem label="Core Usage" value={14} />
              <ResourceItem label="Memory" value={42} />
              <ResourceItem label="I/O Speed" value={8} />
            </div>
          </section>

          <section className="space-y-4">
            <h3 className="text-xl font-serif font-medium">Activity Log</h3>
            <div className="relative border-l-2 border-muted ml-2 space-y-6 py-2">
              <LogItem time="Just now" msg="Post 'Refactor' published" />
              <LogItem time="2h ago" msg="New media asset uploaded" />
              <LogItem time="5h ago" msg="Admin session initiated" />
            </div>
          </section>
        </div>
      </div>
    </div>
  );
}

function ResourceItem({ label, value }: { label: string, value: number }) {
  return (
    <div className="space-y-2">
      <div className="flex justify-between text-[10px] font-bold uppercase tracking-widest text-muted-foreground">
        <span>{label}</span>
        <span>{value}%</span>
      </div>
      <div className="h-1 bg-muted rounded-full overflow-hidden">
        <motion.div 
          initial={{ width: 0 }}
          animate={{ width: `${value}%` }}
          transition={{ duration: 1, ease: "easeOut" }}
          className="h-full bg-accent"
        />
      </div>
    </div>
  )
}

function LogItem({ time, msg }: { time: string, msg: string }) {
  return (
    <div className="ml-6 relative">
      <div className="absolute -left-[33px] w-3 h-3 rounded-full bg-background border-2 border-accent" />
      <p className="text-xs font-bold">{msg}</p>
      <p className="text-[10px] text-muted-foreground uppercase font-bold tracking-tighter mt-0.5">{time}</p>
    </div>
  )
}

import { Button } from "@/components/ui/button";
