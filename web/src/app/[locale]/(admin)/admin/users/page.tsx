"use client";

import { useTranslations } from "next-intl";
import { Users, Construction } from "lucide-react";

export default function UsersPage() {
  const t = useTranslations("admin");

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-serif font-medium tracking-tight">
          {t("user_management")}
        </h1>
        <p className="text-sm text-muted-foreground mt-1">
          Manage user accounts, roles and permissions.
        </p>
      </div>

      <div className="flex flex-col items-center justify-center py-24 border border-dashed border-border rounded-2xl">
        <div className="w-14 h-14 rounded-2xl bg-muted/50 flex items-center justify-center mb-4">
          <Users className="w-7 h-7 text-muted-foreground/40" />
        </div>
        <h2 className="text-lg font-serif font-medium mb-1">Coming Soon</h2>
        <p className="text-sm text-muted-foreground max-w-sm text-center">
          User management requires backend API support. This feature is on the
          roadmap.
        </p>
        <div className="flex items-center gap-1.5 mt-4 text-xs text-muted-foreground/60">
          <Construction className="w-3.5 h-3.5" />
          <span>Pending: User CRUD API, role assignment endpoints</span>
        </div>
      </div>
    </div>
  );
}
