"use client";

import { useState } from "react";
import { useRouter } from "@/i18n/routing";
import { useTranslations } from 'next-intl';
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { motion, AnimatePresence } from "framer-motion";
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
  ArrowLeft,
  RefreshCw,
  ShieldAlert,
  Users,
  UserPlus
} from "lucide-react";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Checkbox } from "@/components/ui/checkbox";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import SunWaveBackground from "@/components/site/sun-wave-background";
import { useCheckDB, useSetup } from "@/services/system-service";
import { cn } from "@/lib/utils";

const setupSchema = z.object({
  dbHost: z.string().min(1, "Host is required"),
  dbPort: z.string().regex(/^\d+$/, "Port must be a number"),
  dbUser: z.string().min(1, "User is required"),
  dbPass: z.string(),
  dbName: z.string().min(1, "Database name is required"),
  siteName: z.string().min(1, "Site name is required"),
  adminUsername: z.string().min(3, "Username must be at least 3 characters"),
  adminEmail: z.string().email("Invalid email address"),
  adminPassword: z.string().min(8, "Password must be at least 8 characters"),
  confirmPassword: z.string(),
  adminFullAccess: z.boolean(),
  adminCanDelete: z.boolean(),
  userCanUpload: z.boolean(),
  allowAnonymousRead: z.boolean(),
}).refine((data) => data.adminPassword === data.confirmPassword, {
  message: "Passwords do not match",
  path: ["confirmPassword"],
});

type SetupFormValues = z.infer<typeof setupSchema>;

export default function SetupPage() {
  const t = useTranslations('setup');
  const router = useRouter();
  const [step, setStep] = useState(1);
  const [dbVerified, setDbVerified] = useState(false);

  const {
    register,
    handleSubmit,
    watch,
    setValue,
    formState: { errors },
    trigger,
  } = useForm<SetupFormValues>({
    resolver: zodResolver(setupSchema),
    defaultValues: {
      dbHost: "localhost",
      dbPort: "5432",
      dbUser: "postgres",
      dbPass: "",
      dbName: "kaldalis_cms",
      siteName: "Kaldalis CMS",
      adminFullAccess: true,
      adminCanDelete: true,
      userCanUpload: true,
      allowAnonymousRead: true,
    },
  });

  const checkDBMutation = useCheckDB();
  const setupMutation = useSetup();

  const handleTestConnection = async () => {
    const fields: (keyof SetupFormValues)[] = ['dbHost', 'dbPort', 'dbUser', 'dbPass', 'dbName'];
    const isValid = await trigger(fields);
    if (!isValid) return;

    const values = watch();
    checkDBMutation.mutate({
      host: values.dbHost,
      port: Number(values.dbPort),
      user: values.dbUser,
      pass: values.dbPass,
      name: values.dbName,
    }, {
      onSuccess: () => setDbVerified(true),
      onError: () => setDbVerified(false),
    });
  };

  const onNextStep = async () => {
    let isValid = false;
    if (step === 1) {
      isValid = await trigger(['dbHost', 'dbPort', 'dbUser', 'dbPass', 'dbName']);
      if (isValid && !dbVerified) {
        // Force re-verify or show error
        return;
      }
    } else if (step === 2) {
      isValid = await trigger(['adminUsername', 'adminEmail', 'adminPassword', 'confirmPassword']);
    }
    
    if (isValid) setStep(step + 1);
  };

  const onSubmit = (data: SetupFormValues) => {
    setupMutation.mutate({
      ...data,
      db_host: data.dbHost,
      db_port: Number(data.dbPort),
      db_user: data.dbUser,
      db_pass: data.dbPass,
      db_name: data.dbName,
      admin_username: data.adminUsername,
      admin_email: data.adminEmail,
      admin_password: data.adminPassword,
    });
  };

  if (setupMutation.isSuccess) {
    return (
      <div className="relative min-h-screen flex items-center justify-center p-4">
        <SunWaveBackground />
        <motion.div initial={{ scale: 0.9, opacity: 0 }} animate={{ scale: 1, opacity: 1 }}>
          <Card className="w-full max-w-md border-0 bg-white/70 dark:bg-slate-900/70 backdrop-blur-xl shadow-2xl">
            <CardHeader className="text-center pb-2">
              <div className="mx-auto mb-4 h-16 w-16 bg-emerald-500/10 rounded-full flex items-center justify-center">
                <CheckCircle2 className="h-10 w-10 text-emerald-500" />
              </div>
              <CardTitle className="text-2xl font-bold">{t('success_title')}</CardTitle>
              <CardDescription className="mt-2">{t('success_desc')}</CardDescription>
            </CardHeader>
            <CardContent className="pt-6">
              <Button 
                className="w-full h-12 bg-indigo-600 hover:bg-indigo-700 text-white rounded-xl font-bold transition-all"
                onClick={() => router.push("/login")}
              >
                Go to Login
              </Button>
            </CardContent>
          </Card>
        </motion.div>
      </div>
    );
  }

  return (
    <div className="relative min-h-screen flex items-center justify-center p-4 overflow-hidden bg-slate-950">
      <SunWaveBackground />
      
      <div className="w-full max-w-xl z-10">
        <div className="flex justify-between mb-12 px-4">
          {[1, 2, 3].map((i) => (
            <div key={i} className="flex flex-col items-center gap-3">
              <div className={cn(
                "h-1 w-24 md:w-36 rounded-full transition-all duration-700",
                step >= i ? "bg-indigo-500 shadow-[0_0_15px_rgba(99,102,241,0.6)]" : "bg-slate-800"
              )} />
              <span className={cn(
                "text-[10px] font-bold uppercase tracking-[0.2em]",
                step === i ? "text-indigo-400" : "text-slate-600"
              )}>
                {i === 1 ? "Database" : i === 2 ? "Admin" : "Site"}
              </span>
            </div>
          ))}
        </div>

        <Card className="border-0 bg-slate-900/80 backdrop-blur-2xl shadow-2xl ring-1 ring-white/5">
          <CardHeader>
            <CardTitle className="text-3xl font-bold tracking-tight flex items-center gap-3 text-white">
              {step === 1 && <Database className="h-8 w-8 text-indigo-400" />}
              {step === 2 && <ShieldCheck className="h-8 w-8 text-indigo-400" />}
              {step === 3 && <Rocket className="h-8 w-8 text-indigo-400" />}
              {t('title')}
            </CardTitle>
            <CardDescription className="text-slate-400">
              {step === 1 && "Configure your PostgreSQL connection"}
              {step === 2 && "Create the master administrator account"}
              {step === 3 && "Almost there! Site name & RBAC roles"}
            </CardDescription>
          </CardHeader>

          <CardContent className="pt-4">
            <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
              <AnimatePresence mode="wait">
                {step === 1 && (
                  <motion.div 
                    key="step1"
                    initial={{ x: 20, opacity: 0 }}
                    animate={{ x: 0, opacity: 1 }}
                    exit={{ x: -20, opacity: 0 }}
                    className="space-y-4"
                  >
                    <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
                      <div className="md:col-span-3 space-y-2">
                        <Label className="text-[10px] font-bold uppercase tracking-widest text-slate-500">Host</Label>
                        <div className="relative">
                          <Server className="absolute left-3 top-3 h-4 w-4 text-slate-500" />
                          <Input {...register("dbHost")} className="pl-10 bg-white/5 border-white/10" placeholder="localhost" />
                        </div>
                        {errors.dbHost && <p className="text-rose-500 text-[10px]">{errors.dbHost.message}</p>}
                      </div>
                      <div className="space-y-2">
                        <Label className="text-[10px] font-bold uppercase tracking-widest text-slate-500">Port</Label>
                        <Input {...register("dbPort")} className="bg-white/5 border-white/10" placeholder="5432" />
                      </div>
                    </div>

                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                      <div className="space-y-2">
                        <Label className="text-[10px] font-bold uppercase tracking-widest text-slate-500">Username</Label>
                        <Input {...register("dbUser")} className="bg-white/5 border-white/10" placeholder="postgres" />
                      </div>
                      <div className="space-y-2">
                        <Label className="text-[10px] font-bold uppercase tracking-widest text-slate-500">Password</Label>
                        <div className="relative">
                          <KeyRound className="absolute left-3 top-3 h-4 w-4 text-slate-500" />
                          <Input type="password" {...register("dbPass")} className="pl-10 bg-white/5 border-white/10" placeholder="••••••••" />
                        </div>
                      </div>
                    </div>

                    <div className="space-y-2">
                      <Label className="text-[10px] font-bold uppercase tracking-widest text-slate-500">Database Name</Label>
                      <Input {...register("dbName")} className="bg-white/5 border-white/10" placeholder="kaldalis_cms" />
                    </div>

                    <div className="pt-4 space-y-4">
                      <Button 
                        type="button" 
                        onClick={handleTestConnection} 
                        disabled={checkDBMutation.isPending}
                        variant={dbVerified ? "outline" : "secondary"}
                        className={cn(
                          "w-full h-12 rounded-xl font-bold transition-all",
                          dbVerified && "border-emerald-500/50 text-emerald-400 bg-emerald-500/5"
                        )}
                      >
                        {checkDBMutation.isPending ? <Loader2 className="h-5 w-5 animate-spin mr-2" /> : dbVerified ? <CheckCircle2 className="h-5 w-5 mr-2" /> : <RefreshCw className="h-5 w-5 mr-2" />}
                        {checkDBMutation.isPending ? "Testing..." : dbVerified ? "Connection Verified" : "Test Connection"}
                      </Button>

                      <Button 
                        type="button" 
                        onClick={onNextStep} 
                        disabled={!dbVerified}
                        className="w-full h-12 bg-indigo-600 hover:bg-indigo-700 text-white rounded-xl font-bold shadow-lg disabled:opacity-50"
                      >
                        Next Step <ArrowRight className="ml-2 h-4 w-4" />
                      </Button>
                    </div>
                  </motion.div>
                )}

                {step === 2 && (
                  <motion.div 
                    key="step2"
                    initial={{ x: 20, opacity: 0 }}
                    animate={{ x: 0, opacity: 1 }}
                    exit={{ x: -20, opacity: 0 }}
                    className="space-y-4"
                  >
                    <div className="space-y-2">
                      <Label className="text-[10px] font-bold uppercase tracking-widest text-slate-500">{t('admin_username')}</Label>
                      <div className="relative">
                        <User className="absolute left-3 top-3 h-4 w-4 text-slate-500" />
                        <Input {...register("adminUsername")} className="pl-10 bg-white/5 border-white/10" placeholder="admin" />
                      </div>
                      {errors.adminUsername && <p className="text-rose-500 text-[10px]">{errors.adminUsername.message}</p>}
                    </div>
                    <div className="space-y-2">
                      <Label className="text-[10px] font-bold uppercase tracking-widest text-slate-500">{t('admin_email')}</Label>
                      <div className="relative">
                        <Mail className="absolute left-3 top-3 h-4 w-4 text-slate-500" />
                        <Input {...register("adminEmail")} type="email" className="pl-10 bg-white/5 border-white/10" placeholder="admin@example.com" />
                      </div>
                      {errors.adminEmail && <p className="text-rose-500 text-[10px]">{errors.adminEmail.message}</p>}
                    </div>
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                      <div className="space-y-2">
                        <Label className="text-[10px] font-bold uppercase tracking-widest text-slate-500">{t('admin_password')}</Label>
                        <Input type="password" {...register("adminPassword")} className="bg-white/5 border-white/10" placeholder="••••••••" />
                        {errors.adminPassword && <p className="text-rose-500 text-[10px]">{errors.adminPassword.message}</p>}
                      </div>
                      <div className="space-y-2">
                        <Label className="text-[10px] font-bold uppercase tracking-widest text-slate-500">Confirm</Label>
                        <Input type="password" {...register("confirmPassword")} className="bg-white/5 border-white/10" placeholder="••••••••" />
                        {errors.confirmPassword && <p className="text-rose-500 text-[10px]">{errors.confirmPassword.message}</p>}
                      </div>
                    </div>

                    <div className="flex gap-4 pt-4">
                      <Button type="button" variant="outline" onClick={() => setStep(1)} className="flex-1 h-12 rounded-xl border-white/10">
                        <ArrowLeft className="mr-2 h-4 w-4" /> Back
                      </Button>
                      <Button type="button" onClick={onNextStep} className="flex-[2] h-12 bg-indigo-600 hover:bg-indigo-700 text-white rounded-xl font-bold">
                        Next Step <ArrowRight className="ml-2 h-4 w-4" />
                      </Button>
                    </div>
                  </motion.div>
                )}

                {step === 3 && (
                  <motion.div 
                    key="step3"
                    initial={{ x: 20, opacity: 0 }}
                    animate={{ x: 0, opacity: 1 }}
                    exit={{ x: -20, opacity: 0 }}
                    className="space-y-6"
                  >
                    <div className="space-y-2">
                      <Label className="text-[10px] font-bold uppercase tracking-widest text-slate-500">{t('site_name')}</Label>
                      <div className="relative">
                        <Globe className="absolute left-3 top-3 h-5 w-5 text-indigo-400" />
                        <Input {...register("siteName")} className="pl-12 h-14 text-lg rounded-xl bg-white/5 border-white/10" placeholder="Kaldalis CMS" />
                      </div>
                      {errors.siteName && <p className="text-rose-500 text-[10px]">{errors.siteName.message}</p>}
                    </div>

                    <div className="space-y-3 p-5 rounded-xl bg-white/5 border border-white/10">
                      <div className="flex items-center gap-2 mb-2">
                        <ShieldCheck className="h-4 w-4 text-indigo-400" />
                        <h4 className="text-[10px] font-bold uppercase tracking-widest text-indigo-400">Security Configuration</h4>
                      </div>
                      
                      <div className="space-y-4">
                        <RBACItem 
                          id="adminFullAccess" 
                          label="Master Super-Admin" 
                          desc="Wildcard access to all system resources"
                          checked={watch("adminFullAccess")}
                          onCheckedChange={(v) => setValue("adminFullAccess", v as boolean)}
                        />
                        <RBACItem 
                          id="adminCanDelete" 
                          label="Staff Delete Power" 
                          desc="Allow editors to delete posts & media"
                          checked={watch("adminCanDelete")}
                          onCheckedChange={(v) => setValue("adminCanDelete", v as boolean)}
                        />
                        <RBACItem 
                          id="userCanUpload" 
                          label="User Media Upload" 
                          desc="Allow registered users to upload files"
                          checked={watch("userCanUpload")}
                          onCheckedChange={(v) => setValue("userCanUpload", v as boolean)}
                        />
                        <RBACItem 
                          id="allowAnonymousRead" 
                          label="Public Guest Access" 
                          desc="Allow visitors to view content without login"
                          checked={watch("allowAnonymousRead")}
                          onCheckedChange={(v) => setValue("allowAnonymousRead", v as boolean)}
                        />
                      </div>
                    </div>

                    <div className="bg-indigo-500/10 border border-indigo-500/20 rounded-xl p-4 flex gap-3">
                      <ShieldAlert className="h-5 w-5 text-indigo-400 shrink-0 mt-0.5" />
                      <p className="text-[10px] leading-relaxed text-slate-400">
                        <span className="font-bold text-indigo-400">INITIALIZATION NOTICE:</span> Submitting will establish your RBAC hierarchy and persist roles in the security engine.
                      </p>
                    </div>

                    <div className="flex gap-4">
                      <Button type="button" variant="outline" onClick={() => setStep(2)} className="flex-1 h-12 rounded-xl border-white/10">
                        Back
                      </Button>
                      <Button 
                        type="submit" 
                        disabled={setupMutation.isPending} 
                        className="flex-[2] h-12 bg-indigo-600 hover:bg-indigo-700 text-white rounded-xl font-bold shadow-xl shadow-indigo-600/20"
                      >
                        {setupMutation.isPending ? <Loader2 className="h-5 w-5 animate-spin mr-2" /> : <ShieldCheck className="h-5 w-5 mr-2" />}
                        {setupMutation.isPending ? t('setting_up') : t('complete_setup')}
                      </Button>
                    </div>
                  </motion.div>
                )}
              </AnimatePresence>
            </form>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}

function RBACItem({ id, label, desc, checked, onCheckedChange }: { id: string, label: string, desc: string, checked: boolean, onCheckedChange: (v: boolean) => void }) {
  return (
    <div className="flex items-start space-x-3">
      <Checkbox id={id} checked={checked} onCheckedChange={onCheckedChange} className="mt-1 border-white/20 data-[state=checked]:bg-indigo-500" />
      <div className="grid gap-1 leading-none">
        <Label htmlFor={id} className="text-xs font-semibold cursor-pointer text-slate-200">
          {label}
        </Label>
        <p className="text-[10px] text-slate-500">{desc}</p>
      </div>
    </div>
  );
}
