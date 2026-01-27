"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import {Link} from '@/i18n/routing'; // 引入 Link 用于跳转注册
import { useAuthStore } from "@/store/useAuthStore";
import api from "@/lib/api";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { useTranslations } from 'next-intl';

export default function LoginPage() {
  const router = useRouter();
  // 1. 修复：正确从 Zustand 获取 setLogin 方法
  const setLogin = useAuthStore((state) => state.setLogin);
  const t = useTranslations();
  
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setLoading(true);
    setError("");

    // 2. 修复：获取表单数据并构造 Payload
    const formData = new FormData(e.currentTarget);
    const payload = {
      username: formData.get("username") as string,
      password: formData.get("password") as string,
    };

    try {
      // 3. 请求接口
      // 注意：根据你的 api.ts 拦截器，res 应该是后端返回的 JSON body
      const res: any = await api.post("/users/login", payload);
      
      console.log("登录响应数据:", res); // 调试用，看后端到底返了什么

      // 4. 解析数据
      // 后端返回: { message: "...", user: { username: "admin", role: "admin" } }
      const userData = res.user; // <-- 修正：从 res.user 获取嵌套的用户对象

      if (!userData || !userData.role) {
        throw new Error("返回数据格式错误，未找到用户信息或角色");
      }

      // 5. 存入 Zustand
      setLogin(userData); // 传递正确的 user 对象

      // 6. 核心跳转逻辑
      if (userData.role === 'admin' || userData.role === 'super_admin') {
        router.replace("/admin/dashboard");
      } else {
        // 普通用户去首页
        router.replace("/");
      }

    } catch (err: any) {
      console.error("登录错误:", err);
      // 优先显示后端返回的错误信息
      setError(err.response?.data?.message || err.message || t('auth.login_failed'));
    } finally {
      setLoading(false);
    }
  };

  return (
    <Card className="w-full max-w-md shadow-lg">
      <CardHeader className="space-y-1">
        <CardTitle className="text-2xl font-bold">{t('common.app_name')}</CardTitle>
        <CardDescription>{t('auth.login_title')}</CardDescription>
      </CardHeader>
      
      <form onSubmit={handleSubmit}>
        <CardContent className="space-y-4">
          {error && (
            <div className="p-3 text-sm font-medium text-red-500 bg-red-50 rounded-md border border-red-200">
              {error}
            </div>
          )}
          
          <div className="space-y-2">
            <Label htmlFor="username">{t('auth.username')}</Label>
            <Input id="username" name="username" placeholder="admin" required disabled={loading} />
          </div>
          
          <div className="space-y-2">
            <Label htmlFor="password">{t('auth.password')}</Label>
            <Input id="password" name="password" type="password" required disabled={loading} />
          </div>
        </CardContent>
        
        <CardFooter className="flex flex-col space-y-2">
          <Button className="w-full" type="submit" disabled={loading}>
            {loading ? t('auth.logging_in') : t('auth.sign_in')}
          </Button>
          
          {/* 补充：去注册页面的链接 */}
          <div className="text-sm text-center text-slate-500">
            {t('auth.dont_have_account')}{" "}
            <Link href="/register" className="text-blue-600 hover:underline">
              {t('auth.sign_up')}
            </Link>
          </div>
        </CardFooter>
      </form>
    </Card>
  );
}
