"use client";

import { useState } from "react";
import { useTranslations } from 'next-intl';
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Globe, Shield, Save, Mail, Layout, Code } from "lucide-react";
import { cn } from "@/lib/utils";

export default function SettingsPage() {
  const t = useTranslations('admin');
  const [activeTab, setActiveTab] = useState("general");
  const [loading, setLoading] = useState(false);

  const handleSave = () => {
    setLoading(true);
    // Simulate API call
    setTimeout(() => {
      setLoading(false);
      alert("Settings saved successfully (Mock)");
    }, 1000);
  };

  const tabs = [
    { id: "general", label: t('site_information'), icon: Globe },
    { id: "seo", label: t('seo_settings'), icon: SearchIcon },
    { id: "security", label: t('security_settings'), icon: Shield },
    { id: "advanced", label: t('advanced_settings'), icon: Code },
  ];

  return (
    <div className="h-full flex flex-col gap-6 text-slate-200 font-sans">
      
      {/* Header */}
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4 shrink-0">
        <div>
           <h1 className="text-2xl font-bold tracking-tight text-white mb-1">{t('system_settings')}</h1>
           <p className="text-slate-400 text-sm">
             Configure global site settings and preferences.
           </p>
        </div>
        
        <Button 
          onClick={handleSave}
          disabled={loading}
          className="h-10 bg-[#ad2bee] hover:bg-[#9225c9] text-white border-0 shadow-[0_4px_12px_rgba(173,43,238,0.3)] transition-all hover:scale-105 font-medium px-6"
        >
            <Save className="mr-2 h-4 w-4" /> 
            {loading ? t('loading') : t('save_changes')}
        </Button>
      </div>

      <div className="flex flex-col md:flex-row gap-8 flex-1 overflow-hidden">
        
        {/* Sidebar Tabs */}
        <div className="w-full md:w-64 flex flex-col gap-2 shrink-0">
          {tabs.map(tab => (
            <button
              key={tab.id}
              onClick={() => setActiveTab(tab.id)}
              className={cn(
                "flex items-center gap-3 px-4 py-3 text-sm font-medium rounded-lg transition-all text-left",
                activeTab === tab.id 
                  ? "bg-[#ad2bee]/10 text-[#ad2bee] border border-[#ad2bee]/20 shadow-[inset_0_0_10px_rgba(173,43,238,0.1)]" 
                  : "text-slate-400 hover:text-slate-200 hover:bg-slate-800/50"
              )}
            >
              <tab.icon className="h-4 w-4" />
              {tab.label}
            </button>
          ))}
        </div>

        {/* Content Area */}
        <div className="flex-1 overflow-y-auto custom-scrollbar pb-10">
          <div className="max-w-2xl space-y-8">
            
            {activeTab === 'general' && (
              <div className="space-y-6 animate-in fade-in slide-in-from-right-4 duration-500">
                <Card className="bg-[#0d0b14]/40 border-slate-800/60 shadow-xl">
                  <CardHeader>
                    <CardTitle className="text-lg text-white">General Information</CardTitle>
                    <CardDescription className="text-slate-500">Basic site details.</CardDescription>
                  </CardHeader>
                  <CardContent className="space-y-4">
                    <div className="space-y-2">
                      <Label className="text-slate-300">Site Title</Label>
                      <Input defaultValue="Kaldalis CMS" className="bg-[#0d0b14] border-slate-800 focus-visible:ring-[#ad2bee]/50" />
                    </div>
                    <div className="space-y-2">
                      <Label className="text-slate-300">Tagline</Label>
                      <Input defaultValue="A modern content management system." className="bg-[#0d0b14] border-slate-800 focus-visible:ring-[#ad2bee]/50" />
                    </div>
                    <div className="space-y-2">
                      <Label className="text-slate-300">Site URL</Label>
                      <Input defaultValue="https://kaldalis.com" className="bg-[#0d0b14] border-slate-800 focus-visible:ring-[#ad2bee]/50" />
                    </div>
                  </CardContent>
                </Card>

                <Card className="bg-[#0d0b14]/40 border-slate-800/60 shadow-xl">
                  <CardHeader>
                    <CardTitle className="text-lg text-white">Contact Info</CardTitle>
                    <CardDescription className="text-slate-500">Displayed in footer and contact forms.</CardDescription>
                  </CardHeader>
                  <CardContent className="space-y-4">
                    <div className="space-y-2">
                      <Label className="text-slate-300">Admin Email</Label>
                      <div className="relative">
                        <Mail className="absolute left-3 top-2.5 h-4 w-4 text-slate-500" />
                        <Input defaultValue="admin@kaldalis.com" className="pl-9 bg-[#0d0b14] border-slate-800 focus-visible:ring-[#ad2bee]/50" />
                      </div>
                    </div>
                  </CardContent>
                </Card>
              </div>
            )}

            {activeTab === 'seo' && (
              <div className="space-y-6 animate-in fade-in slide-in-from-right-4 duration-500">
                 <Card className="bg-[#0d0b14]/40 border-slate-800/60 shadow-xl">
                  <CardHeader>
                    <CardTitle className="text-lg text-white">Global SEO</CardTitle>
                    <CardDescription className="text-slate-500">Default meta tags for the site.</CardDescription>
                  </CardHeader>
                  <CardContent className="space-y-4">
                    <div className="space-y-2">
                      <Label className="text-slate-300">Meta Description</Label>
                      <textarea className="flex w-full rounded-md border border-slate-800 bg-[#0d0b14] px-3 py-2 text-sm text-slate-200 placeholder:text-slate-500 focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-[#ad2bee]/50 disabled:cursor-not-allowed disabled:opacity-50 min-h-[100px]" defaultValue="Kaldalis is a high-performance CMS built for developers and content creators." />
                    </div>
                    <div className="space-y-2">
                      <Label className="text-slate-300">Keywords</Label>
                      <Input defaultValue="cms, golang, nextjs, react" className="bg-[#0d0b14] border-slate-800 focus-visible:ring-[#ad2bee]/50" />
                    </div>
                  </CardContent>
                </Card>
              </div>
            )}

            {/* Other tabs can be empty placeholders for now */}
            {activeTab === 'security' && (
                <div className="flex flex-col items-center justify-center py-20 text-slate-500 animate-in fade-in">
                    <Shield className="h-12 w-12 opacity-20 mb-4" />
                    <p>Security settings are managed via environment variables.</p>
                </div>
            )}
             {activeTab === 'advanced' && (
                <div className="flex flex-col items-center justify-center py-20 text-slate-500 animate-in fade-in">
                    <Code className="h-12 w-12 opacity-20 mb-4" />
                    <p>Advanced configuration.</p>
                </div>
            )}

          </div>
        </div>
      </div>
    </div>
  );
}

// Helper icon
function SearchIcon(props: any) {
  return (
    <svg
      {...props}
      xmlns="http://www.w3.org/2000/svg"
      width="24"
      height="24"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    >
      <circle cx="11" cy="11" r="8" />
      <path d="m21 21-4.3-4.3" />
    </svg>
  )
}
