"use client";

import { useState } from "react";
import { useRouter, Link } from "@/i18n/routing";
import api from "@/lib/api";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { useTranslations } from 'next-intl';

export default function RegisterPage() {
  const router = useRouter();
  const t = useTranslations();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setLoading(true);
    setError("");

    const formData = new FormData(e.currentTarget);
    const username = formData.get("username") as string;
    const email = formData.get("email") as string;
    const password = formData.get("password") as string;
    const confirmPassword = formData.get("confirmPassword") as string;

    if (password !== confirmPassword) {
      setError("Passwords do not match");
      setLoading(false);
      return;
    }

    try {
      await api.post("/users/register", { 
        username, 
        email, 
        password 
      });

      alert("Registration successful! Please login.");
      router.push("/login"); 
    } catch (err: any) {
      console.error(err);
      setError(err.response?.data?.message || "Registration failed");
    } finally {
      setLoading(false);
    }
  };

  return (
    <Card className="w-full max-w-md shadow-lg">
      <CardHeader className="space-y-1">
        <CardTitle className="text-2xl font-bold">{t('auth.register_title')}</CardTitle>
        <CardDescription>Join Kaldalis CMS and start creating</CardDescription>
      </CardHeader>
      <form onSubmit={handleSubmit}>
        <CardContent className="space-y-4">
          {error && <div className="text-sm font-medium text-red-500 bg-red-50 p-2 rounded border border-red-100">{error}</div>}
          
          <div className="space-y-2">
            <Label htmlFor="username">{t('auth.username')}</Label>
            <Input id="username" name="username" placeholder="yourname" required disabled={loading} />
          </div>

          <div className="space-y-2">
            <Label htmlFor="email">{t('auth.email')}</Label>
            <Input id="email" name="email" type="email" placeholder="name@example.com" required disabled={loading} />
          </div>

          <div className="space-y-2">
            <Label htmlFor="password">{t('auth.password')}</Label>
            <Input id="password" name="password" type="password" required disabled={loading} />
          </div>

          <div className="space-y-2">
            <Label htmlFor="confirmPassword">{t('auth.confirm_password')}</Label>
            <Input id="confirmPassword" name="confirmPassword" type="password" required disabled={loading} />
          </div>
        </CardContent>
        <CardFooter className="flex flex-col space-y-2">
          <Button className="w-full" type="submit" disabled={loading}>
            {loading ? t('common.loading') : t('auth.sign_up')}
          </Button>
          <div className="text-sm text-center text-slate-500">
            {t('auth.already_have_account')}{" "}
            <Link href="/login" className="text-blue-600 hover:underline">
              {t('auth.sign_in')}
            </Link>
          </div>
        </CardFooter>
      </form>
    </Card>
  );
}