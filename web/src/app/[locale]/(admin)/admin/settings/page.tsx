"use client";

import { useState } from "react";
import { useTranslations } from "next-intl";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Globe, Shield, Save, Mail, Code, Search } from "lucide-react";
import { cn } from "@/lib/utils";

export default function SettingsPage() {
  const t = useTranslations("admin");
  const [activeTab, setActiveTab] = useState("general");
  const [loading, setLoading] = useState(false);

  const handleSave = () => {
    setLoading(true);
    setTimeout(() => {
      setLoading(false);
      alert("Settings saved successfully (Mock)");
    }, 1000);
  };

  const tabs = [
    { id: "general", label: t("site_information"), icon: Globe },
    { id: "seo", label: t("seo_settings"), icon: Search },
    { id: "security", label: t("security_settings"), icon: Shield },
    { id: "advanced", label: t("advanced_settings"), icon: Code },
  ];

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col md:flex-row md:items-end justify-between gap-4">
        <div>
          <h1 className="text-3xl font-serif font-medium tracking-tight">
            {t("system_settings")}
          </h1>
          <p className="text-sm text-muted-foreground mt-1">
            Configure global site settings and preferences.
          </p>
        </div>

        <Button
          onClick={handleSave}
          disabled={loading}
          className="rounded-full bg-accent text-accent-foreground hover:bg-accent/90 h-11 px-6 font-bold shadow-lg shadow-accent/10 transition-all"
        >
          <Save className="mr-2 h-4 w-4" />
          {loading ? t("loading") : t("save_changes")}
        </Button>
      </div>

      <div className="flex flex-col md:flex-row gap-6">
        {/* Sidebar Tabs */}
        <div className="w-full md:w-56 flex flex-col gap-1 shrink-0">
          {tabs.map((tab) => (
            <button
              key={tab.id}
              onClick={() => setActiveTab(tab.id)}
              className={cn(
                "flex items-center gap-3 px-4 py-2.5 text-sm font-medium rounded-xl transition-colors text-left",
                activeTab === tab.id
                  ? "bg-accent/10 text-accent"
                  : "text-muted-foreground hover:text-foreground hover:bg-muted"
              )}
            >
              <tab.icon className="h-4 w-4" />
              {tab.label}
            </button>
          ))}
        </div>

        {/* Content Area */}
        <div className="flex-1 pb-10">
          <div className="max-w-2xl space-y-6">
            {activeTab === "general" && (
              <div className="space-y-6">
                <Card className="border-border bg-white/60 dark:bg-white/[0.03] shadow-none rounded-xl">
                  <CardHeader>
                    <CardTitle className="text-lg font-serif">
                      General Information
                    </CardTitle>
                    <CardDescription>Basic site details.</CardDescription>
                  </CardHeader>
                  <CardContent className="space-y-4">
                    <div className="space-y-2">
                      <Label>Site Title</Label>
                      <Input
                        defaultValue="Kaldalis CMS"
                        className="bg-background border-border focus-visible:ring-accent/50 rounded-xl"
                      />
                    </div>
                    <div className="space-y-2">
                      <Label>Tagline</Label>
                      <Input
                        defaultValue="A modern content management system."
                        className="bg-background border-border focus-visible:ring-accent/50 rounded-xl"
                      />
                    </div>
                    <div className="space-y-2">
                      <Label>Site URL</Label>
                      <Input
                        defaultValue="https://kaldalis.com"
                        className="bg-background border-border focus-visible:ring-accent/50 rounded-xl"
                      />
                    </div>
                  </CardContent>
                </Card>

                <Card className="border-border bg-white/60 dark:bg-white/[0.03] shadow-none rounded-xl">
                  <CardHeader>
                    <CardTitle className="text-lg font-serif">
                      Contact Info
                    </CardTitle>
                    <CardDescription>
                      Displayed in footer and contact forms.
                    </CardDescription>
                  </CardHeader>
                  <CardContent className="space-y-4">
                    <div className="space-y-2">
                      <Label>Admin Email</Label>
                      <div className="relative">
                        <Mail className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                        <Input
                          defaultValue="admin@kaldalis.com"
                          className="pl-10 bg-background border-border focus-visible:ring-accent/50 rounded-xl"
                        />
                      </div>
                    </div>
                  </CardContent>
                </Card>
              </div>
            )}

            {activeTab === "seo" && (
              <div className="space-y-6">
                <Card className="border-border bg-white/60 dark:bg-white/[0.03] shadow-none rounded-xl">
                  <CardHeader>
                    <CardTitle className="text-lg font-serif">
                      Global SEO
                    </CardTitle>
                    <CardDescription>
                      Default meta tags for the site.
                    </CardDescription>
                  </CardHeader>
                  <CardContent className="space-y-4">
                    <div className="space-y-2">
                      <Label>Meta Description</Label>
                      <textarea
                        className="flex w-full rounded-xl border border-border bg-background px-3 py-2 text-sm placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-accent/50 min-h-[100px]"
                        defaultValue="Kaldalis is a high-performance CMS built for developers and content creators."
                      />
                    </div>
                    <div className="space-y-2">
                      <Label>Keywords</Label>
                      <Input
                        defaultValue="cms, golang, nextjs, react"
                        className="bg-background border-border focus-visible:ring-accent/50 rounded-xl"
                      />
                    </div>
                  </CardContent>
                </Card>
              </div>
            )}

            {activeTab === "security" && (
              <div className="flex flex-col items-center justify-center py-20 text-muted-foreground">
                <Shield className="h-12 w-12 opacity-20 mb-4" />
                <p className="text-sm">
                  Security settings are managed via environment variables.
                </p>
              </div>
            )}

            {activeTab === "advanced" && (
              <div className="flex flex-col items-center justify-center py-20 text-muted-foreground">
                <Code className="h-12 w-12 opacity-20 mb-4" />
                <p className="text-sm">Advanced configuration.</p>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
