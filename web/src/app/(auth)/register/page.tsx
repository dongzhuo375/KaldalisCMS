"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import api from "@/lib/api";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";

export default function RegisterPage() {
  const router = useRouter();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setLoading(true);
    setError("");

    const formData = new FormData(e.currentTarget);
    const username = formData.get("username") as string;
    const email = formData.get("email") as string; // 新增邮箱
    const password = formData.get("password") as string;
    const confirmPassword = formData.get("confirmPassword") as string;

    // 1. 前端校验
    if (password !== confirmPassword) {
      setError("两次输入的密码不一致");
      setLoading(false);
      return;
    }

    try {
      // 2. 调用后端接口
      // 目标路径: /users/register
      // 发送数据: { username, email, password }
      await api.post("/users/register", { 
        username, 
        email, 
        password 
      });

      // 3. 注册成功，跳转登录
      // 这里的 alert 仅作演示，实际可以用 toast 组件
      alert("注册成功！请登录");
      router.push("/login"); 
    } catch (err: any) {
      console.error(err);
      setError(err.response?.data?.message || "注册失败，请检查输入");
    } finally {
      setLoading(false);
    }
  };

  return (
    <Card className="w-full max-w-md shadow-lg">
      <CardHeader className="space-y-1">
        <CardTitle className="text-2xl font-bold">注册新账号</CardTitle>
        <CardDescription>加入 Kaldalis CMS 开启内容创作</CardDescription>
      </CardHeader>
      <form onSubmit={handleSubmit}>
        <CardContent className="space-y-4">
          {error && <div className="text-sm font-medium text-red-500">{error}</div>}
          
          {/* 用户名 */}
          <div className="space-y-2">
            <Label htmlFor="username">用户名</Label>
            <Input id="username" name="username" placeholder="yourname" required />
          </div>

          {/* 新增：邮箱 */}
          <div className="space-y-2">
            <Label htmlFor="email">邮箱地址</Label>
            <Input id="email" name="email" type="email" placeholder="name@example.com" required />
          </div>

          {/* 密码 */}
          <div className="space-y-2">
            <Label htmlFor="password">密码</Label>
            <Input id="password" name="password" type="password" required />
          </div>

          {/* 确认密码 */}
          <div className="space-y-2">
            <Label htmlFor="confirmPassword">确认密码</Label>
            <Input id="confirmPassword" name="confirmPassword" type="password" required />
          </div>
        </CardContent>
        <CardFooter className="flex flex-col space-y-2">
          <Button className="w-full" type="submit" disabled={loading}>
            {loading ? "注册中..." : "立即注册"}
          </Button>
          <div className="text-sm text-center text-slate-500">
            已有账号？{" "}
            <Link href="/login" className="text-blue-600 hover:underline">
              去登录
            </Link>
          </div>
        </CardFooter>
      </form>
    </Card>
  );
}
