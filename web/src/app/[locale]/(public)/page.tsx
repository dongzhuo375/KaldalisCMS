"use client";

import { useAuthStore } from "@/store/useAuthStore";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
import {Link} from '@/i18n/routing';
import { ArrowRight, BookOpen, Users, Shield } from "lucide-react";
import { useTranslations } from 'next-intl';

export default function HomePage() {
  const { user, isLoggedIn } = useAuthStore();
  const t = useTranslations();

  return (
    <div className="space-y-16">
      {/* Hero 区域 */}
      <section className="text-center py-24 space-y-8">
        <h1 className="text-4xl font-extrabold tracking-tight lg:text-6xl text-slate-900">
          {t('common.welcome')} Kaldalis CMS
        </h1>
        <p className="text-xl text-slate-600 max-w-3xl mx-auto">
          A modern content management system built with Go and Next.js.
          Fast, secure, and easy to extend.
        </p>
        <div className="flex justify-center gap-4 pt-4">
           {isLoggedIn ? (
             <Button size="lg" className="rounded-full px-8">
               开始探索
             </Button>
           ) : (
             <>
                <Link href="/register">
                  <Button size="lg" className="rounded-full px-8">{t('auth.sign_up')}</Button>
                </Link>
                <Link href="/login">
                  <Button variant="outline" size="lg" className="rounded-full px-8">{t('auth.sign_in')}</Button>
                </Link>
             </>
           )}
        </div>
      </section>

      {/* 如果登录了，显示用户欢迎卡片 */}
      {isLoggedIn && user && (
        <section className="max-w-4xl mx-auto">
          <Card className="border-l-4 border-l-blue-500 shadow-sm">
            <CardHeader>
              <CardTitle>欢迎回来, {user.username} 👋</CardTitle>
              <CardDescription>
                当前身份: <span className="font-mono bg-slate-100 px-2 py-0.5 rounded text-slate-800">{user.role}</span>
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <p className="text-slate-600">
                您现在位于前台首页。普通用户可以在这里浏览文章、管理个人资料。
                {user.role === 'admin' && " 由于您是管理员，您也可以进入后台管理系统。"}
              </p>
              <div className="flex gap-3">
                 <Button variant="secondary" className="gap-2">
                   <Users className="h-4 w-4" /> 个人资料
                 </Button>
                 {/* 只有管理员显示这个按钮 */}
                 {(user.role === 'admin' || user.role === 'super_admin') && (
                   <Link href="/admin/dashboard">
                     <Button className="gap-2">
                       <Shield className="h-4 w-4" /> 进入后台
                     </Button>
                   </Link>
                 )}
              </div>
            </CardContent>
          </Card>
        </section>
      )}

      {/* 功能特性展示 (占位) */}
      <section className="grid grid-cols-1 md:grid-cols-3 gap-6 pt-10">
        <FeatureCard 
          icon={<BookOpen className="h-8 w-8 text-blue-500" />}
          title="内容管理"
          desc="高效的文章发布与编辑体验，支持 Markdown 与富文本。"
        />
        <FeatureCard 
          icon={<Users className="h-8 w-8 text-green-500" />}
          title="用户系统"
          desc="完善的 RBAC 权限控制，支持多角色分级管理。"
        />
        <FeatureCard 
          icon={<Shield className="h-8 w-8 text-purple-500" />}
          title="安全可靠"
          desc="基于 Go Gin 与 Casbin 构建的坚固后端安全防线。"
        />
      </section>
    </div>
  );
}

// 简单的特性小组件
function FeatureCard({ icon, title, desc }: { icon: any, title: string, desc: string }) {
  return (
    <Card className="hover:shadow-md transition-shadow">
      <CardHeader>
        <div className="mb-2">{icon}</div>
        <CardTitle className="text-xl">{title}</CardTitle>
      </CardHeader>
      <CardContent>
        <p className="text-slate-500">{desc}</p>
      </CardContent>
    </Card>
  )
}
