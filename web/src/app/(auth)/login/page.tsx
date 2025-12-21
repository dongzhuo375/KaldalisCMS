// src/app/(auth)/login/page.tsx
"use client";

import { useRouter } from "next/navigation";
import { useAuthStore } from "@/store/useAuthStore";
import api from "@/lib/api";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { useState } from "react";

export default function LoginPage() {
  const router = useRouter();
  const setLogin = useAuthStore((state) => state.setLogin);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setLoading(true);
    setError("");

    const formData = new FormData(e.currentTarget);
    const username = formData.get("username") as string;
    const password = formData.get("password") as string;

    try {
      // 1. 调用 Go 后端登录接口
    
      const res: any = await api.post("/users/login", { username, password });

      // 2. 登录成功：把用户信息存入 Zustand
      // 假设后端返回数据格式为: { code: 200, data: { username: 'admin', role: 'admin' } }
      setLogin(res.data);

      // 3. 跳转到后台仪表盘
      router.push("/admin/dashboard");
    } catch (err: any) {
      setError(err.response?.data?.message || "登录失败，请检查账号密码");
    } finally {
      setLoading(false);
    }
  };

  return (
    <Card className="w-full max-w-md shadow-lg">
      <CardHeader className="space-y-1">
        <CardTitle className="text-2xl font-bold">Kaldalis CMS</CardTitle>
        <CardDescription>请输入您的账号密码进入管理后台</CardDescription>
      </CardHeader>
      <form onSubmit={handleSubmit}>
        <CardContent className="space-y-4">
          {error && <div className="text-sm font-medium text-red-500">{error}</div>}
          <div className="space-y-2">
            <Label htmlFor="username">用户名</Label>
            <Input id="username" name="username" placeholder="admin" required />
          </div>
          <div className="space-y-2">
            <Label htmlFor="password">密码</Label>
            <Input id="password" name="password" type="password" required />
          </div>
        </CardContent>
        <CardFooter>
          <Button className="w-full" type="submit" disabled={loading}>
            {loading ? "登录中..." : "立即登录"}
          </Button>
        </CardFooter>
      </form>
    </Card>
  );
}
