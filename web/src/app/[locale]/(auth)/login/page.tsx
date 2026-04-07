"use client";

import { useState } from "react";
import { useRouter } from "@/i18n/routing";
import { Link } from '@/i18n/routing';
import { useAuthStore } from "@/store/useAuthStore";
import api from "@/lib/api";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useTranslations } from 'next-intl';
import { motion } from "framer-motion";
import { Loader2 } from "lucide-react";

export default function LoginPage() {
  const router = useRouter();
  const setLogin = useAuthStore((state) => state.setLogin);
  const t = useTranslations();

  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setLoading(true);
    setError("");

    const formData = new FormData(e.currentTarget);
    const payload = {
      username: formData.get("username") as string,
      password: formData.get("password") as string,
    };

    try {
      const res: any = await api.post("/users/login", payload);
      const userData = res.user;

      if (!userData || !userData.role) {
        throw new Error("返回数据格式错误，未找到用户信息或角色");
      }

      setLogin(userData);

      if (userData.role === 'admin' || userData.role === 'super_admin') {
        router.replace("/admin/dashboard");
      } else {
        router.replace("/");
      }

    } catch (err: any) {
      console.error("登录错误:", err);
      setError(err.response?.data?.message || err.message || t('auth.login_failed'));
    } finally {
      setLoading(false);
    }
  };

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.5 }}
      className="w-full max-w-sm"
    >
      <div className="bg-white/80 dark:bg-card/80 backdrop-blur-xl rounded-3xl p-8 shadow-xl shadow-black/5 border border-border/50">
        {/* Header */}
        <div className="text-center mb-8">
          <h1 className="text-2xl font-serif font-medium text-foreground mb-2">
            {t('common.app_name')}
          </h1>
          <p className="text-muted-foreground text-sm">
            {t('auth.login_title')}
          </p>
        </div>

        {/* Form */}
        <form onSubmit={handleSubmit} className="space-y-5">
          {error && (
            <div className="p-3 text-sm text-red-600 bg-red-50 dark:bg-red-900/20 rounded-xl border border-red-200 dark:border-red-800">
              {error}
            </div>
          )}

          <div className="space-y-2">
            <Label htmlFor="username" className="text-sm font-medium">
              {t('auth.username')}
            </Label>
            <Input
              id="username"
              name="username"
              placeholder="admin"
              required
              disabled={loading}
              className="h-11 rounded-xl bg-muted/50 border-border/50 focus:border-accent"
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="password" className="text-sm font-medium">
              {t('auth.password')}
            </Label>
            <Input
              id="password"
              name="password"
              type="password"
              required
              disabled={loading}
              className="h-11 rounded-xl bg-muted/50 border-border/50 focus:border-accent"
            />
          </div>

          <Button
            className="w-full h-11 rounded-xl font-medium mt-2"
            type="submit"
            disabled={loading}
          >
            {loading ? (
              <>
                <Loader2 className="w-4 h-4 mr-2 animate-spin" />
                {t('auth.logging_in')}
              </>
            ) : (
              t('auth.sign_in')
            )}
          </Button>
        </form>

        {/* Footer */}
        <div className="mt-6 text-center text-sm text-muted-foreground">
          {t('auth.dont_have_account')}{" "}
          <Link href="/register" className="text-accent hover:underline font-medium">
            {t('auth.sign_up')}
          </Link>
        </div>
      </div>
    </motion.div>
  );
}
