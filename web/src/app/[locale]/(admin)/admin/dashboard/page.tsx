"use client";

import {
  FileText,
  Activity,
  Database,
  Zap,
  Plus,
} from "lucide-react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { cn, getImageUrl } from "@/lib/utils";
import { useTranslations } from "next-intl";
import { useSystemStatus, useReadyz } from "@/services/system-service";
import { useAdminPosts } from "@/services/post-service";
import { PostStatus } from "@/lib/types";
import { useMedia } from "@/services/media-service";
import { motion } from "framer-motion";
import { Link } from "@/i18n/routing";
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";

export default function DashboardPage() {
  const t = useTranslations("admin");
  const { data: status, isLoading: statusLoading } = useSystemStatus();
  const { data: health, isLoading: healthLoading } = useReadyz();
  const { data: postsData = [], isLoading: postsLoading } = useAdminPosts();
  const { data: mediaData, isLoading: mediaLoading } = useMedia({ page_size: 1 });

  const posts = Array.isArray(postsData) ? postsData : [];
  const isLoading = statusLoading || healthLoading || postsLoading || mediaLoading;

  if (isLoading) {
    return (
      <div className="space-y-10 pb-20">
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
        <Skeleton className="h-[400px] rounded-3xl" />
      </div>
    );
  }

  const dbStatus = health?.checks?.database?.status === "ok";

  const stats = [
    {
      title: "Database",
      value: dbStatus ? "Connected" : "Disconnected",
      status: dbStatus ? "success" : "error",
      icon: Database,
      detail: dbStatus ? "PostgreSQL" : "Check logs",
    },
    {
      title: "API",
      value: health?.status === "ok" ? "Healthy" : "Degraded",
      status: health?.status === "ok" ? "success" : "warning",
      icon: Zap,
      detail: `${(health?.mode || "app").toUpperCase()} mode`,
    },
    {
      title: "Posts",
      value: posts.length.toString(),
      status: "neutral",
      icon: FileText,
      detail: `${posts.filter((p) => p.status === PostStatus.PUBLISHED).length} published`,
    },
    {
      title: "Media",
      value: mediaData?.total?.toString() || "0",
      status: "neutral",
      icon: Activity,
      detail: "Files uploaded",
    },
  ];

  return (
    <div className="space-y-10 pb-20">
      {/* Header */}
      <header className="flex flex-col md:flex-row md:items-end justify-between gap-6">
        <motion.div
          initial={{ opacity: 0, y: 10 }}
          animate={{ opacity: 1, y: 0 }}
          className="space-y-2"
        >
          <h1 className="text-4xl md:text-5xl font-serif font-medium tracking-tight">
            {status?.site_name || "Console"}
          </h1>
          <p className="text-muted-foreground font-medium flex items-center gap-2">
            <span
              className={cn(
                "w-2 h-2 rounded-full",
                dbStatus ? "bg-emerald-500 animate-pulse" : "bg-rose-500"
              )}
            />
            {dbStatus ? "All systems operational" : "System issue detected"}
          </p>
        </motion.div>

        <Link href="/admin/posts/new">
          <Button className="rounded-full bg-accent text-accent-foreground h-12 px-6 font-bold shadow-lg shadow-accent/10 hover:shadow-xl transition-all">
            <Plus className="w-4 h-4 mr-2" /> New Article
          </Button>
        </Link>
      </header>

      {/* Stats */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        {stats.map((stat, index) => (
          <motion.div
            key={index}
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: index * 0.08 }}
          >
            <Card className="border-border bg-white/60 dark:bg-white/[0.03] backdrop-blur-sm shadow-none hover:border-accent/20 transition-colors group">
              <CardHeader className="flex flex-row items-center justify-between pb-2 space-y-0">
                <CardTitle className="text-[10px] font-bold uppercase tracking-[0.2em] text-muted-foreground group-hover:text-accent transition-colors">
                  {stat.title}
                </CardTitle>
                <stat.icon
                  className={cn(
                    "h-4 w-4",
                    stat.status === "success"
                      ? "text-emerald-500"
                      : stat.status === "error"
                        ? "text-rose-500"
                        : "text-muted-foreground"
                  )}
                  strokeWidth={2.5}
                />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-serif font-medium">
                  {stat.value}
                </div>
                <div className="text-[10px] font-bold text-muted-foreground/60 uppercase tracking-tighter mt-1 italic">
                  {stat.detail}
                </div>
              </CardContent>
            </Card>
          </motion.div>
        ))}
      </div>

      {/* Recent Posts */}
      <section>
        <div className="flex items-center justify-between mb-6">
          <h2 className="text-2xl font-serif font-medium">Recent Articles</h2>
          <Link
            href="/admin/posts"
            className="text-xs font-bold uppercase text-accent hover:underline"
          >
            View All
          </Link>
        </div>

        {posts.length === 0 ? (
          <div className="text-center py-16 border border-dashed border-border rounded-2xl">
            <FileText className="w-8 h-8 mx-auto text-muted-foreground/30 mb-3" />
            <p className="text-sm text-muted-foreground">No articles yet</p>
            <Link href="/admin/posts/new">
              <Button variant="outline" size="sm" className="mt-4 rounded-xl">
                <Plus className="w-4 h-4 mr-1.5" /> Create your first post
              </Button>
            </Link>
          </div>
        ) : (
          <div className="space-y-1">
            {posts.map((item, i) => (
              <motion.div
                key={item.id}
                initial={{ opacity: 0, x: -10 }}
                animate={{ opacity: 1, x: 0 }}
                transition={{ delay: 0.3 + i * 0.04 }}
              >
                <Link
                  href={`/admin/posts/${item.id}/edit`}
                  className="flex items-center justify-between p-4 rounded-xl hover:bg-muted/50 transition-colors group border border-transparent hover:border-border"
                >
                  <div className="flex items-center gap-4">
                    <div className="w-10 h-10 rounded-full bg-accent/5 overflow-hidden flex items-center justify-center text-accent group-hover:bg-accent group-hover:text-white transition-all shrink-0">
                      {item.cover ? (
                        <img
                          src={getImageUrl(item.cover)}
                          alt=""
                          className="w-full h-full object-cover"
                        />
                      ) : (
                        <FileText className="w-5 h-5" />
                      )}
                    </div>
                    <div className="min-w-0">
                      <h4 className="font-bold text-sm group-hover:text-accent transition-colors truncate">
                        {item.title}
                      </h4>
                      <p className="text-[10px] text-muted-foreground uppercase tracking-widest font-bold">
                        {new Date(item.created_at).toLocaleDateString(
                          undefined,
                          { month: "long", day: "numeric" }
                        )}
                      </p>
                    </div>
                  </div>
                  <Badge
                    variant="outline"
                    className={cn(
                      "rounded-full border-0 px-3 py-1 text-[10px] font-bold uppercase shrink-0",
                      item.status === PostStatus.PUBLISHED
                        ? "bg-accent/10 text-accent"
                        : "bg-muted text-muted-foreground"
                    )}
                  >
                    {item.status === PostStatus.PUBLISHED
                      ? "Published"
                      : "Draft"}
                  </Badge>
                </Link>
              </motion.div>
            ))}
          </div>
        )}
      </section>
    </div>
  );
}
