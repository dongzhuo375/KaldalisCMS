"use client";

import { useTranslations } from 'next-intl';
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { ArrowUpRight, ArrowDownRight, Users, Eye, MousePointerClick, Clock, BarChart3, Globe } from "lucide-react";
import { cn } from "@/lib/utils";

export default function AnalyticsPage() {
  const t = useTranslations('admin');

  const stats = [
    {
      title: t('total_views'),
      value: "1.2M",
      change: "+12%",
      trend: "up",
      icon: Eye,
      color: "text-blue-400"
    },
    {
      title: "Unique Visitors",
      value: "843K",
      change: "+8%",
      trend: "up",
      icon: Users,
      color: "text-emerald-400"
    },
    {
      title: "Bounce Rate",
      value: "42%",
      change: "-2%",
      trend: "down", // Good for bounce rate usually, but visually down arrow
      icon: MousePointerClick,
      color: "text-rose-400"
    },
    {
      title: "Avg. Session",
      value: "4m 12s",
      change: "+30s",
      trend: "up",
      icon: Clock,
      color: "text-purple-400"
    }
  ];

  return (
    <div className="h-full flex flex-col gap-6 text-slate-200 font-sans overflow-y-auto custom-scrollbar">
      
      {/* Header */}
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4 shrink-0">
        <div>
           <h1 className="text-2xl font-bold tracking-tight text-white mb-1">{t('analytics_dashboard')}</h1>
           <p className="text-slate-400 text-sm">
             Monitor your website performance and user engagement.
           </p>
        </div>
      </div>

      {/* KPI Cards */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        {stats.map((stat, i) => (
          <Card key={i} className="bg-[#0d0b14]/40 border-slate-800/60 shadow-lg">
             <CardHeader className="flex flex-row items-center justify-between pb-2">
                <CardTitle className="text-sm font-medium text-slate-400">{stat.title}</CardTitle>
                <stat.icon className={cn("h-4 w-4", stat.color)} />
             </CardHeader>
             <CardContent>
                <div className="text-2xl font-bold text-white">{stat.value}</div>
                <p className="text-xs flex items-center mt-1 text-slate-500">
                   {stat.trend === 'up' ? (
                     <ArrowUpRight className="h-3 w-3 mr-1 text-emerald-500" />
                   ) : (
                     <ArrowDownRight className="h-3 w-3 mr-1 text-rose-500" />
                   )}
                   <span className={stat.trend === 'up' ? "text-emerald-500" : "text-rose-500"}>{stat.change}</span>
                   <span className="ml-1">vs last month</span>
                </p>
             </CardContent>
          </Card>
        ))}
      </div>

      {/* Charts Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
         
         {/* Traffic Chart (Mock CSS) */}
         <Card className="bg-[#0d0b14]/40 border-slate-800/60 shadow-lg col-span-1 lg:col-span-2">
            <CardHeader>
               <CardTitle className="text-base text-white flex items-center gap-2">
                 <BarChart3 className="h-4 w-4 text-blue-400" /> {t('traffic_overview')}
               </CardTitle>
            </CardHeader>
            <CardContent>
               <div className="h-64 w-full flex items-end justify-between gap-2 pt-4 px-2">
                  {[40, 65, 45, 80, 55, 70, 40, 60, 50, 75, 85, 95, 60, 40, 65, 55, 80, 90, 75, 65, 50, 45, 60, 70].map((h, i) => (
                    <div key={i} className="w-full bg-slate-800/50 rounded-t-sm relative group">
                       <div 
                         className="absolute bottom-0 left-0 right-0 bg-gradient-to-t from-blue-600 to-blue-400 opacity-60 group-hover:opacity-100 transition-all duration-500 rounded-t-sm"
                         style={{ height: `${h}%` }}
                       ></div>
                    </div>
                  ))}
               </div>
               <div className="flex justify-between mt-2 text-xs text-slate-500 font-mono">
                  <span>00:00</span>
                  <span>06:00</span>
                  <span>12:00</span>
                  <span>18:00</span>
                  <span>23:59</span>
               </div>
            </CardContent>
         </Card>

         {/* Top Sources */}
         <Card className="bg-[#0d0b14]/40 border-slate-800/60 shadow-lg">
            <CardHeader>
               <CardTitle className="text-base text-white flex items-center gap-2">
                 <Globe className="h-4 w-4 text-purple-400" /> Top Sources
               </CardTitle>
            </CardHeader>
            <CardContent>
               <div className="space-y-4">
                  {[
                    { name: "Google", val: 65, col: "bg-blue-500" },
                    { name: "Direct", val: 20, col: "bg-emerald-500" },
                    { name: "Twitter / X", val: 10, col: "bg-slate-500" },
                    { name: "GitHub", val: 5, col: "bg-purple-500" },
                  ].map((source, i) => (
                    <div key={i} className="space-y-1">
                       <div className="flex justify-between text-xs text-slate-300">
                          <span>{source.name}</span>
                          <span>{source.val}%</span>
                       </div>
                       <div className="h-1.5 w-full bg-slate-800 rounded-full overflow-hidden">
                          <div className={`h-full ${source.col}`} style={{ width: `${source.val}%` }}></div>
                       </div>
                    </div>
                  ))}
               </div>
            </CardContent>
         </Card>

         {/* Device Stats */}
         <Card className="bg-[#0d0b14]/40 border-slate-800/60 shadow-lg">
            <CardHeader>
               <CardTitle className="text-base text-white flex items-center gap-2">
                 <MousePointerClick className="h-4 w-4 text-rose-400" /> Device Type
               </CardTitle>
            </CardHeader>
            <CardContent className="flex items-center justify-center py-6">
                <div className="flex gap-8">
                   <div className="text-center space-y-2">
                      <div className="w-16 h-32 border-2 border-slate-700 rounded-lg mx-auto flex items-center justify-center bg-slate-800/20">
                         <span className="text-xs font-bold text-white">45%</span>
                      </div>
                      <p className="text-xs text-slate-500">Mobile</p>
                   </div>
                   <div className="text-center space-y-2">
                      <div className="w-24 h-32 border-2 border-slate-700 rounded-lg mx-auto flex items-center justify-center bg-slate-800/20">
                         <span className="text-xs font-bold text-white">50%</span>
                      </div>
                      <p className="text-xs text-slate-500">Desktop</p>
                   </div>
                   <div className="text-center space-y-2">
                      <div className="w-20 h-32 border-2 border-slate-700 rounded-lg mx-auto flex items-center justify-center bg-slate-800/20">
                         <span className="text-xs font-bold text-white">5%</span>
                      </div>
                      <p className="text-xs text-slate-500">Tablet</p>
                   </div>
                </div>
            </CardContent>
         </Card>

      </div>
    </div>
  );
}