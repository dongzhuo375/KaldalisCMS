"use client";

import { useAuthStore } from "@/store/useAuthStore";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
import {Link} from '@/i18n/routing';
import { ArrowRight, BookOpen, Users, Shield, Zap, CheckCircle2, Layers } from "lucide-react";
import { useTranslations } from 'next-intl';
import FluidBackground from "@/components/site/fluid-background";

export default function HomePage() {
  const { user, isLoggedIn } = useAuthStore();
  const t = useTranslations();

  return (
    <div className="relative z-0 min-h-[calc(100vh-4rem)] flex flex-col justify-center">
      {/* 3D 流体背景 */}
      <FluidBackground />

      <div className="space-y-24 py-12 md:py-20">
        
        {/* Hero 区域 */}
        <section className="text-center space-y-8 max-w-5xl mx-auto px-4">
          
          {/* 版本通告 Badge */}
          <div className="inline-flex items-center rounded-full border border-slate-200 bg-white/50 px-3 py-1 text-sm font-medium text-slate-800 backdrop-blur-md dark:border-slate-800 dark:bg-slate-900/50 dark:text-slate-200 shadow-sm transition-colors hover:bg-white/80 dark:hover:bg-slate-900/80 cursor-default">
            <span className="flex h-2 w-2 rounded-full bg-blue-600 mr-2"></span>
            VERSION 2.4.0 NOW AVAILABLE
          </div>

          <h1 className="text-5xl font-extrabold tracking-tight lg:text-7xl text-slate-900 dark:text-slate-50 drop-shadow-sm">
            {t('common.welcome')} <span className="text-transparent bg-clip-text bg-gradient-to-r from-blue-600 to-purple-600 dark:from-blue-400 dark:to-purple-400">{t('common.app_name')}</span>
          </h1>
          
          <p className="text-xl md:text-2xl text-slate-600 dark:text-slate-300 max-w-3xl mx-auto leading-relaxed">
            {t('home.hero_subtitle')}
          </p>
          
          {/* 特性标签 */}
          <div className="flex flex-wrap justify-center gap-4 md:gap-8 text-sm font-medium text-slate-500 dark:text-slate-400">
            <div className="flex items-center gap-1.5"><Zap className="w-4 h-4 text-yellow-500" /> 快速</div>
            <div className="flex items-center gap-1.5"><Shield className="w-4 h-4 text-blue-500" /> 安全</div>
            <div className="flex items-center gap-1.5"><Layers className="w-4 h-4 text-purple-500" /> 易于扩展</div>
          </div>

          <div className="flex flex-col sm:flex-row justify-center gap-4 pt-8">
              {isLoggedIn ? (
                <Button size="lg" className="h-12 px-8 text-base rounded-full shadow-lg hover:shadow-xl transition-all hover:scale-105">
                  {t('home.start_exploring')} <ArrowRight className="ml-2 w-4 h-4" />
                </Button>
            ) : (
              <>
                  <Link href="/register">
                    <Button size="lg" className="h-12 px-8 text-base rounded-full shadow-lg bg-slate-900 hover:bg-slate-800 dark:bg-white dark:text-slate-900 dark:hover:bg-slate-100 transition-all hover:scale-105">
                      {t('auth.sign_up')} <ArrowRight className="ml-2 w-4 h-4" />
                    </Button>
                  </Link>
                  <Link href="/login">
                    <Button variant="outline" size="lg" className="h-12 px-8 text-base rounded-full border-slate-300 hover:bg-white/50 dark:border-slate-700 dark:hover:bg-slate-800/50 backdrop-blur-sm transition-all">
                      {t('auth.sign_in')}
                    </Button>
                  </Link>
              </>
            )}
          </div>
        </section>

        {/* 如果登录了，显示用户欢迎卡片 (Glassmorphism) */}
        {isLoggedIn && user && (
          <section className="max-w-4xl mx-auto px-4 w-full animate-in fade-in slide-in-from-bottom-4 duration-700">
            <Card className="border-0 shadow-2xl bg-white/70 dark:bg-slate-900/60 backdrop-blur-xl ring-1 ring-slate-200/50 dark:ring-slate-700/50">
              <CardHeader>
                <CardTitle className="flex items-center gap-2 text-2xl">
                  {t('home.welcome_back')}, <span className="bg-gradient-to-r from-blue-600 to-cyan-500 bg-clip-text text-transparent">{user.username}</span> 👋
                </CardTitle>
                <CardDescription className="text-base">
                  {t('home.current_role')}: <span className="font-mono bg-blue-50 dark:bg-blue-900/30 text-blue-700 dark:text-blue-300 px-2 py-0.5 rounded border border-blue-100 dark:border-blue-800">{user.role}</span>
                </CardDescription>
              </CardHeader>
              <CardContent className="space-y-6">
                <p className="text-slate-600 dark:text-slate-300 text-lg">
                  {t('home.user_welcome_desc')}
                  {user.role === 'admin' && <span className="font-medium text-slate-800 dark:text-slate-200">{t('home.admin_welcome_suffix')}</span>}
                </p>
                <div className="flex flex-wrap gap-4">
                  <Button variant="secondary" className="gap-2 h-10 px-6 shadow-sm hover:shadow-md transition-shadow">
                    <Users className="h-4 w-4" /> {t('navigation.personal_profile')}
                  </Button>
                  {/* 只有管理员显示这个按钮 */}
                  {(user.role === 'admin' || user.role === 'super_admin') && (
                    <Link href="/admin/dashboard">
                      <Button className="gap-2 h-10 px-6 bg-slate-900 text-white hover:bg-slate-800 dark:bg-white dark:text-slate-900 dark:hover:bg-slate-200 shadow-lg hover:shadow-xl transition-all hover:-translate-y-0.5">
                        <Shield className="h-4 w-4" /> {t('navigation.enter_admin')}
                      </Button>
                    </Link>
                  )}
                </div>
              </CardContent>
            </Card>
          </section>
        )}

        {/* 功能特性展示 (Glass Cards) */}
        <section className="grid grid-cols-1 md:grid-cols-3 gap-6 max-w-6xl mx-auto px-4">
          <FeatureCard 
            icon={<BookOpen className="h-8 w-8 text-blue-500" />}
            title={t('home.features.content_management')}
            desc={t('home.features.content_management_desc')}
          />
          <FeatureCard 
            icon={<Users className="h-8 w-8 text-green-500" />}
            title={t('home.features.user_system')}
            desc={t('home.features.user_system_desc')}
          />
          <FeatureCard 
            icon={<Shield className="h-8 w-8 text-purple-500" />}
            title={t('home.features.security')}
            desc={t('home.features.security_desc')}
          />
        </section>
      </div>
    </div>
  );
}

function FeatureCard({ icon, title, desc }: { icon: any, title: string, desc: string }) {
  return (
    <Card className="group hover:shadow-xl transition-all duration-300 border-0 bg-white/50 dark:bg-slate-900/50 backdrop-blur-sm ring-1 ring-slate-200/50 dark:ring-slate-800/50 hover:-translate-y-1">
      <CardHeader>
        <div className="mb-4 p-3 rounded-2xl bg-white dark:bg-slate-800 w-fit shadow-sm group-hover:scale-110 transition-transform duration-300">{icon}</div>
        <CardTitle className="text-xl font-bold text-slate-800 dark:text-slate-100">{title}</CardTitle>
      </CardHeader>
      <CardContent>
        <p className="text-slate-600 dark:text-slate-400 leading-relaxed">{desc}</p>
      </CardContent>
    </Card>
  )
}