"use client";

import { useState } from "react";
import { useRouter } from "@/i18n/routing";
import { useTranslations } from 'next-intl';
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Rocket, ShieldCheck, Globe, Mail, User, Lock, Loader2, CheckCircle2 } from "lucide-react";
import FluidBackground from "@/components/site/fluid-background";
import api from "@/lib/api";

export default function SetupPage() {
  const t = useTranslations('setup');
  const router = useRouter();
  const [step, setStep] = useState(1);
  const [loading, setLoading] = useState(false);
  const [success, setSuccess] = useState(false);
  const [error, setError] = useState("");

  const [formData, setFormData] = useState({
    siteName: "Kaldalis CMS",
    adminUsername: "",
    adminEmail: "",
    adminPassword: "",
    confirmPassword: ""
  });

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  const handleNext = () => {
    if (step === 1 && !formData.siteName) return;
    setStep(step + 1);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (formData.adminPassword !== formData.confirmPassword) {
      setError("Passwords do not match");
      return;
    }

    setLoading(true);
    setError("");

    try {
      await api.post("/system/setup", {
        site_name: formData.siteName,
        admin_username: formData.adminUsername,
        admin_email: formData.adminEmail,
        admin_password: formData.adminPassword
      });
      setSuccess(true);
    } catch (err: any) {
      setError(err.response?.data?.error || "Setup failed. Please check your backend.");
    } finally {
      setLoading(false);
    }
  };

  if (success) {
    return (
      <div className="relative min-h-screen flex items-center justify-center p-4">
        <FluidBackground />
        <Card className="w-full max-w-md border-0 bg-white/70 dark:bg-slate-900/70 backdrop-blur-xl shadow-2xl animate-in zoom-in-95 duration-500">
          <CardHeader className="text-center pb-2">
            <div className="mx-auto mb-4 h-16 w-16 bg-emerald-500/10 rounded-full flex items-center justify-center">
              <CheckCircle2 className="h-10 w-10 text-emerald-500" />
            </div>
            <CardTitle className="text-2xl font-bold text-slate-900 dark:text-white">{t('success_title')}</CardTitle>
            <CardDescription className="text-slate-600 dark:text-slate-400 mt-2">
              {t('success_desc')}
            </CardDescription>
          </CardHeader>
          <CardContent className="pt-6">
            <Button 
              className="w-full h-12 bg-slate-900 hover:bg-slate-800 dark:bg-white dark:text-slate-900 dark:hover:bg-slate-100 rounded-xl font-bold transition-all"
              onClick={() => router.push("/login")}
            >
              Go to Login
            </Button>
          </CardContent>
        </Card>
      </div>
    );
  }

  return (
    <div className="relative min-h-screen flex items-center justify-center p-4 overflow-hidden">
      <FluidBackground />
      
      <div className="w-full max-w-lg z-10">
        {/* Progress Bar */}
        <div className="flex justify-between mb-8 px-2">
          {[1, 2].map((i) => (
            <div key={i} className="flex flex-col items-center gap-2">
              <div className={`h-2 w-24 md:w-40 rounded-full transition-all duration-500 ${step >= i ? "bg-indigo-600 shadow-[0_0_10px_rgba(79,70,229,0.5)]" : "bg-slate-200 dark:bg-slate-800"}`} />
            </div>
          ))}
        </div>

        <Card className="border-0 bg-white/80 dark:bg-slate-900/80 backdrop-blur-2xl shadow-2xl ring-1 ring-white/20 dark:ring-slate-800/50">
          <CardHeader className="space-y-1">
            <CardTitle className="text-3xl font-extrabold tracking-tight text-slate-900 dark:text-white flex items-center gap-3">
              <Rocket className="h-8 w-8 text-indigo-600" />
              {t('title')}
            </CardTitle>
            <CardDescription className="text-lg text-slate-600 dark:text-slate-400">
              {t('subtitle')}
            </CardDescription>
          </CardHeader>

          <CardContent className="pt-4">
            <form onSubmit={handleSubmit} className="space-y-6">
              {step === 1 ? (
                <div className="space-y-4 animate-in slide-in-from-right-4 duration-500">
                  <div className="space-y-2">
                    <Label className="text-sm font-bold uppercase tracking-widest text-slate-500 dark:text-slate-400">{t('site_name')}</Label>
                    <div className="relative group">
                      <Globe className="absolute left-3 top-3 h-5 w-5 text-slate-400 group-focus-within:text-indigo-600 transition-colors" />
                      <Input 
                        name="siteName"
                        value={formData.siteName}
                        onChange={handleChange}
                        className="pl-10 h-12 rounded-xl bg-white/50 dark:bg-slate-800/50 border-slate-200 dark:border-slate-700 focus:ring-2 focus:ring-indigo-600 transition-all"
                        placeholder="My Awesome Website"
                        required
                      />
                    </div>
                  </div>
                  <Button 
                    type="button"
                    onClick={handleNext}
                    className="w-full h-12 mt-4 bg-indigo-600 hover:bg-indigo-700 text-white rounded-xl font-bold shadow-lg shadow-indigo-600/20 transition-all hover:scale-[1.02]"
                  >
                    Continue to Admin Setup
                  </Button>
                </div>
              ) : (
                <div className="space-y-4 animate-in slide-in-from-right-4 duration-500">
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div className="space-y-2">
                      <Label className="text-sm font-bold uppercase tracking-widest text-slate-500 dark:text-slate-400">{t('admin_username')}</Label>
                      <div className="relative">
                        <User className="absolute left-3 top-3 h-5 w-5 text-slate-400" />
                        <Input 
                          name="adminUsername"
                          value={formData.adminUsername}
                          onChange={handleChange}
                          className="pl-10 h-12 rounded-xl bg-white/50 dark:bg-slate-800/50"
                          placeholder="admin"
                          required
                        />
                      </div>
                    </div>
                    <div className="space-y-2">
                      <Label className="text-sm font-bold uppercase tracking-widest text-slate-500 dark:text-slate-400">{t('admin_email')}</Label>
                      <div className="relative">
                        <Mail className="absolute left-3 top-3 h-5 w-5 text-slate-400" />
                        <Input 
                          name="adminEmail"
                          type="email"
                          value={formData.adminEmail}
                          onChange={handleChange}
                          className="pl-10 h-12 rounded-xl bg-white/50 dark:bg-slate-800/50"
                          placeholder="admin@example.com"
                          required
                        />
                      </div>
                    </div>
                  </div>

                  <div className="space-y-2">
                    <Label className="text-sm font-bold uppercase tracking-widest text-slate-500 dark:text-slate-400">{t('admin_password')}</Label>
                    <div className="relative">
                      <Lock className="absolute left-3 top-3 h-5 w-5 text-slate-400" />
                      <Input 
                        name="adminPassword"
                        type="password"
                        value={formData.adminPassword}
                        onChange={handleChange}
                        className="pl-10 h-12 rounded-xl bg-white/50 dark:bg-slate-800/50"
                        placeholder="••••••••"
                        required
                      />
                    </div>
                  </div>

                  <div className="space-y-2">
                    <Label className="text-sm font-bold uppercase tracking-widest text-slate-500 dark:text-slate-400">Confirm Password</Label>
                    <div className="relative">
                      <Lock className="absolute left-3 top-3 h-5 w-5 text-slate-400" />
                      <Input 
                        name="confirmPassword"
                        type="password"
                        value={formData.confirmPassword}
                        onChange={handleChange}
                        className="pl-10 h-12 rounded-xl bg-white/50 dark:bg-slate-800/50"
                        placeholder="••••••••"
                        required
                      />
                    </div>
                  </div>

                  {error && (
                    <div className="text-sm text-rose-500 bg-rose-500/10 p-3 rounded-lg border border-rose-500/20">
                      {error}
                    </div>
                  )}

                  <div className="flex gap-4 pt-4">
                    <Button 
                      type="button"
                      variant="outline"
                      onClick={() => setStep(1)}
                      className="flex-1 h-12 rounded-xl border-slate-200 dark:border-slate-700"
                    >
                      Back
                    </Button>
                    <Button 
                      type="submit"
                      disabled={loading}
                      className="flex-[2] h-12 bg-indigo-600 hover:bg-indigo-700 text-white rounded-xl font-bold shadow-lg shadow-indigo-600/20"
                    >
                      {loading ? <Loader2 className="h-5 w-5 animate-spin mr-2" /> : <ShieldCheck className="h-5 w-5 mr-2" />}
                      {loading ? t('setting_up') : t('complete_setup')}
                    </Button>
                  </div>
                </div>
              )}
            </form>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}