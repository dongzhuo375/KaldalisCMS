"use client";

import { useState } from "react";
import { useRouter } from "@/i18n/routing";
import { useTranslations } from 'next-intl';
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { 
  Rocket, 
  ShieldCheck, 
  Globe, 
  Mail, 
  User, 
  Lock, 
  Loader2, 
  CheckCircle2, 
  Database,
  Server,
  KeyRound,
  ArrowRight,
  ArrowLeft
} from "lucide-react";
import FluidBackground from "@/components/site/fluid-background";
import api from "@/lib/api";
import { cn } from "@/lib/utils";

export default function SetupPage() {
  const t = useTranslations('setup');
  const router = useRouter();
  const [step, setStep] = useState(1);
  const [loading, setLoading] = useState(false);
  const [success, setSuccess] = useState(false);
  const [error, setError] = useState("");

  const [formData, setFormData] = useState({
    // Database Config
    dbHost: "localhost",
    dbPort: "5432",
    dbUser: "postgres",
    dbPass: "",
    dbName: "kaldalis_cms",
    // Site & Admin Config
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
    setError("");
    // Basic validation per step
    if (step === 1) {
      if (!formData.dbHost || !formData.dbUser || !formData.dbName) {
        setError("Please fill in all required database fields");
        return;
      }
    }
    if (step === 2) {
      if (!formData.adminUsername || !formData.adminEmail || !formData.adminPassword) {
        setError("Please fill in all admin account fields");
        return;
      }
      if (formData.adminPassword !== formData.confirmPassword) {
        setError("Passwords do not match");
        return;
      }
    }
    setStep(step + 1);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!formData.siteName) {
      setError("Site name is required");
      return;
    }

    setLoading(true);
    setError("");

    try {
      await api.post("/system/setup", {
        db_host: formData.dbHost,
        db_port: parseInt(formData.dbPort),
        db_user: formData.dbUser,
        db_pass: formData.dbPass,
        db_name: formData.dbName,
        site_name: formData.siteName,
        admin_username: formData.adminUsername,
        admin_email: formData.adminEmail,
        admin_password: formData.adminPassword
      });
      setSuccess(true);
    } catch (err: any) {
      setError(err.response?.data?.error || "Setup failed. Check database connectivity.");
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
      
      <div className="w-full max-w-xl z-10">
        {/* Progress Bar */}
        <div className="flex justify-between mb-8 px-4">
          {[1, 2, 3].map((i) => (
            <div key={i} className="flex flex-col items-center gap-2">
              <div className={cn(
                "h-1.5 w-24 md:w-36 rounded-full transition-all duration-500",
                step >= i ? "bg-indigo-600 shadow-[0_0_10px_rgba(79,70,229,0.5)]" : "bg-slate-200 dark:bg-slate-800"
              )} />
              <span className={cn(
                "text-[10px] font-bold uppercase tracking-widest",
                step === i ? "text-indigo-600" : "text-slate-400"
              )}>
                {i === 1 ? "Database" : i === 2 ? "Admin" : "Site"}
              </span>
            </div>
          ))}
        </div>

        <Card className="border-0 bg-white/80 dark:bg-slate-900/80 backdrop-blur-2xl shadow-2xl ring-1 ring-white/20 dark:ring-slate-800/50">
          <CardHeader className="space-y-1">
            <CardTitle className="text-3xl font-extrabold tracking-tight text-slate-900 dark:text-white flex items-center gap-3">
              {step === 1 && <Database className="h-8 w-8 text-indigo-600" />}
              {step === 2 && <ShieldCheck className="h-8 w-8 text-indigo-600" />}
              {step === 3 && <Rocket className="h-8 w-8 text-indigo-600" />}
              {t('title')}
            </CardTitle>
            <CardDescription className="text-lg text-slate-600 dark:text-slate-400">
              {step === 1 && "Configure your PostgreSQL connection"}
              {step === 2 && "Create the master administrator account"}
              {step === 3 && "Almost there! Name your new site"}
            </CardDescription>
          </CardHeader>

          <CardContent className="pt-4">
            {error && (
              <div className="mb-6 text-sm text-rose-500 bg-rose-500/10 p-3 rounded-lg border border-rose-500/20 animate-in shake-in duration-300">
                {error}
              </div>
            )}

            <form onSubmit={handleSubmit} className="space-y-6">
              
              {/* Step 1: Database */}
              {step === 1 && (
                <div className="space-y-4 animate-in slide-in-from-right-4 duration-500">
                  <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
                    <div className="md:col-span-3 space-y-2">
                      <Label className="text-xs font-bold uppercase tracking-widest text-slate-500">Host</Label>
                      <div className="relative">
                        <Server className="absolute left-3 top-3 h-4 w-4 text-slate-400" />
                        <Input name="dbHost" value={formData.dbHost} onChange={handleChange} className="pl-10 bg-white/50 dark:bg-slate-800/50" placeholder="localhost" />
                      </div>
                    </div>
                    <div className="space-y-2">
                      <Label className="text-xs font-bold uppercase tracking-widest text-slate-500">Port</Label>
                      <Input name="dbPort" value={formData.dbPort} onChange={handleChange} className="bg-white/50 dark:bg-slate-800/50" placeholder="5432" />
                    </div>
                  </div>

                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div className="space-y-2">
                      <Label className="text-xs font-bold uppercase tracking-widest text-slate-500">Username</Label>
                      <Input name="dbUser" value={formData.dbUser} onChange={handleChange} className="bg-white/50 dark:bg-slate-800/50" placeholder="postgres" />
                    </div>
                    <div className="space-y-2">
                      <Label className="text-xs font-bold uppercase tracking-widest text-slate-500">Password</Label>
                      <div className="relative">
                        <KeyRound className="absolute left-3 top-3 h-4 w-4 text-slate-400" />
                        <Input type="password" name="dbPass" value={formData.dbPass} onChange={handleChange} className="pl-10 bg-white/50 dark:bg-slate-800/50" placeholder="••••••••" />
                      </div>
                    </div>
                  </div>

                  <div className="space-y-2">
                    <Label className="text-xs font-bold uppercase tracking-widest text-slate-500">Database Name</Label>
                    <Input name="dbName" value={formData.dbName} onChange={handleChange} className="bg-white/50 dark:bg-slate-800/50" placeholder="kaldalis_cms" />
                  </div>

                  <Button type="button" onClick={handleNext} className="w-full h-12 bg-indigo-600 hover:bg-indigo-700 text-white rounded-xl font-bold shadow-lg transition-all">
                    Next: Admin Account <ArrowRight className="ml-2 h-4 w-4" />
                  </Button>
                </div>
              )}

              {/* Step 2: Admin Account */}
              {step === 2 && (
                <div className="space-y-4 animate-in slide-in-from-right-4 duration-500">
                  <div className="space-y-2">
                    <Label className="text-xs font-bold uppercase tracking-widest text-slate-500">{t('admin_username')}</Label>
                    <div className="relative">
                      <User className="absolute left-3 top-3 h-4 w-4 text-slate-400" />
                      <Input name="adminUsername" value={formData.adminUsername} onChange={handleChange} className="pl-10 bg-white/50 dark:bg-slate-800/50" placeholder="admin" required />
                    </div>
                  </div>
                  <div className="space-y-2">
                    <Label className="text-xs font-bold uppercase tracking-widest text-slate-500">{t('admin_email')}</Label>
                    <div className="relative">
                      <Mail className="absolute left-3 top-3 h-4 w-4 text-slate-400" />
                      <Input name="adminEmail" type="email" value={formData.adminEmail} onChange={handleChange} className="pl-10 bg-white/50 dark:bg-slate-800/50" placeholder="admin@example.com" required />
                    </div>
                  </div>
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div className="space-y-2">
                      <Label className="text-xs font-bold uppercase tracking-widest text-slate-500">{t('admin_password')}</Label>
                      <Input name="adminPassword" type="password" value={formData.adminPassword} onChange={handleChange} className="bg-white/50 dark:bg-slate-800/50" placeholder="••••••••" required />
                    </div>
                    <div className="space-y-2">
                      <Label className="text-xs font-bold uppercase tracking-widest text-slate-500">Confirm</Label>
                      <Input name="confirmPassword" type="password" value={formData.confirmPassword} onChange={handleChange} className="bg-white/50 dark:bg-slate-800/50" placeholder="••••••••" required />
                    </div>
                  </div>

                  <div className="flex gap-4 pt-2">
                    <Button type="button" variant="outline" onClick={() => setStep(1)} className="flex-1 h-12 rounded-xl">
                      <ArrowLeft className="mr-2 h-4 w-4" /> Back
                    </Button>
                    <Button type="button" onClick={handleNext} className="flex-[2] h-12 bg-indigo-600 hover:bg-indigo-700 text-white rounded-xl font-bold">
                      Next: Site Settings <ArrowRight className="ml-2 h-4 w-4" />
                    </Button>
                  </div>
                </div>
              )}

              {/* Step 3: Site Info */}
              {step === 3 && (
                <div className="space-y-6 animate-in slide-in-from-right-4 duration-500">
                  <div className="space-y-2">
                    <Label className="text-xs font-bold uppercase tracking-widest text-slate-500">{t('site_name')}</Label>
                    <div className="relative">
                      <Globe className="absolute left-3 top-3 h-5 w-5 text-slate-400" />
                      <Input name="siteName" value={formData.siteName} onChange={handleChange} className="pl-10 h-14 text-lg rounded-xl bg-white/50 dark:bg-slate-800/50" placeholder="Kaldalis CMS" required />
                    </div>
                  </div>

                  <div className="bg-indigo-600/5 border border-indigo-600/10 rounded-xl p-4 text-xs text-slate-500 space-y-2">
                    <p className="font-bold text-indigo-600 uppercase">Ready to Install</p>
                    <p>By clicking the button below, we will initialize the database schema, create your account, and set up the system configuration.</p>
                  </div>

                  <div className="flex gap-4">
                    <Button type="button" variant="outline" onClick={() => setStep(2)} className="flex-1 h-12 rounded-xl">
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