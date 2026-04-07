"use client";

import { useState } from "react";
import { useTranslations } from "next-intl";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Search,
  Upload,
  Image as ImageIcon,
  Trash2,
  File,
  MoreHorizontal,
  Grid,
  List as ListIcon,
  Loader2,
  ExternalLink,
} from "lucide-react";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { cn, getImageUrl } from "@/lib/utils";
import {
  useMedia,
  useUploadMedia,
  useDeleteMedia,
} from "@/services/media-service";
import { toast } from "sonner";

export default function MediaPage() {
  const t = useTranslations("admin");
  const [viewMode, setViewMode] = useState<"grid" | "list">("grid");
  const [search, setSearch] = useState("");

  const { data, isLoading } = useMedia({ q: search });
  const files = data?.items || [];
  const uploadMutation = useUploadMedia();
  const deleteMutation = useDeleteMedia();

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      uploadMutation.mutate(file);
    }
  };

  const handleDelete = (id: number) => {
    if (confirm("Are you sure you want to delete this file?")) {
      deleteMutation.mutate(id);
    }
  };

  const formatSize = (bytes: number) => {
    if (bytes === 0) return "0 Bytes";
    const k = 1024;
    const sizes = ["Bytes", "KB", "MB", "GB"];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + " " + sizes[i];
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col md:flex-row md:items-end justify-between gap-4">
        <div>
          <h1 className="text-3xl font-serif font-medium tracking-tight">
            {t("media_library")}
          </h1>
          <p className="text-sm text-muted-foreground mt-1">
            Manage your images, documents, and other media assets.
          </p>
        </div>

        <div className="relative">
          <input
            type="file"
            id="media-upload"
            className="hidden"
            onChange={handleFileChange}
            accept="image/*,application/pdf"
          />
          <Button
            className="rounded-full bg-accent text-accent-foreground hover:bg-accent/90 h-11 px-6 font-bold shadow-lg shadow-accent/10 hover:shadow-xl transition-all"
            onClick={() => document.getElementById("media-upload")?.click()}
            disabled={uploadMutation.isPending}
          >
            {uploadMutation.isPending ? (
              <Loader2 className="mr-2 h-4 w-4 animate-spin" />
            ) : (
              <Upload className="mr-2 h-4 w-4" />
            )}
            {t("upload_file")}
          </Button>
        </div>
      </div>

      {/* Toolbar */}
      <div className="flex flex-col sm:flex-row gap-3">
        <div className="relative flex-1">
          <Search className="absolute left-4 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground pointer-events-none" />
          <Input
            placeholder={t("search")}
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            className="pl-11 h-11 bg-white/50 dark:bg-white/[0.03] border-border rounded-xl"
          />
        </div>
        <div className="flex items-center gap-2 bg-muted rounded-lg p-1">
          <button
            onClick={() => setViewMode("grid")}
            className={cn(
              "p-2 rounded-md transition-colors",
              viewMode === "grid"
                ? "bg-background text-foreground shadow-sm"
                : "text-muted-foreground hover:text-foreground"
            )}
          >
            <Grid className="h-4 w-4" />
          </button>
          <button
            onClick={() => setViewMode("list")}
            className={cn(
              "p-2 rounded-md transition-colors",
              viewMode === "list"
                ? "bg-background text-foreground shadow-sm"
                : "text-muted-foreground hover:text-foreground"
            )}
          >
            <ListIcon className="h-4 w-4" />
          </button>
        </div>
      </div>

      {/* Content */}
      <div className="rounded-xl border border-border bg-white/60 dark:bg-white/[0.03] min-h-[400px] relative">
        {isLoading ? (
          <div className="absolute inset-0 flex items-center justify-center">
            <Loader2 className="h-8 w-8 animate-spin text-accent" />
          </div>
        ) : files.length === 0 ? (
          <div className="absolute inset-0 flex flex-col items-center justify-center text-muted-foreground">
            <ImageIcon className="h-12 w-12 opacity-20 mb-3" />
            <p className="text-sm">No files found.</p>
          </div>
        ) : viewMode === "grid" ? (
          <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-4 p-4">
            {files.map((file) => (
              <div
                key={file.id}
                className="group relative aspect-square rounded-xl border border-border overflow-hidden bg-muted/30 hover:border-accent/30 transition-colors"
              >
                {file.mime_type.startsWith("image/") ? (
                  <img
                    src={getImageUrl(file.url)}
                    alt={file.filename}
                    className="w-full h-full object-cover"
                  />
                ) : (
                  <div className="w-full h-full flex flex-col items-center justify-center gap-2 text-muted-foreground">
                    <File className="h-8 w-8 opacity-30" />
                    <span className="text-xs font-mono uppercase opacity-50">
                      {file.mime_type.split("/")[1]}
                    </span>
                  </div>
                )}

                {/* Hover overlay */}
                <div className="absolute inset-0 bg-black/60 opacity-0 group-hover:opacity-100 transition-opacity flex flex-col justify-end p-3">
                  <p className="text-xs font-medium text-white truncate mb-1">
                    {file.filename}
                  </p>
                  <div className="flex justify-between items-center">
                    <span className="text-[10px] text-white/60">
                      {formatSize(file.size)}
                    </span>
                    <div className="flex gap-1.5">
                      <a
                        href={getImageUrl(file.url)}
                        target="_blank"
                        rel="noreferrer"
                        className="h-6 w-6 flex items-center justify-center rounded-full hover:bg-white/20 text-white"
                      >
                        <ExternalLink className="h-3.5 w-3.5" />
                      </a>
                      <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                          <button className="h-6 w-6 flex items-center justify-center rounded-full hover:bg-white/20 text-white">
                            <MoreHorizontal className="h-4 w-4" />
                          </button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align="end">
                          <DropdownMenuItem
                            onClick={() => {
                              navigator.clipboard.writeText(
                                getImageUrl(file.url)
                              );
                              toast.success("Link copied to clipboard");
                            }}
                            className="cursor-pointer"
                          >
                            Copy URL
                          </DropdownMenuItem>
                          <DropdownMenuItem
                            onClick={() => handleDelete(file.id)}
                            className="text-destructive focus:text-destructive cursor-pointer"
                          >
                            <Trash2 className="mr-2 h-3.5 w-3.5" />{" "}
                            {t("delete")}
                          </DropdownMenuItem>
                        </DropdownMenuContent>
                      </DropdownMenu>
                    </div>
                  </div>
                </div>
              </div>
            ))}
          </div>
        ) : (
          <div className="divide-y divide-border">
            {files.map((file) => (
              <div
                key={file.id}
                className="flex items-center justify-between py-3 px-4 hover:bg-muted/50 transition-colors group"
              >
                <div className="flex items-center gap-4">
                  <div className="h-10 w-10 rounded-lg bg-muted overflow-hidden flex items-center justify-center shrink-0">
                    {file.mime_type.startsWith("image/") ? (
                      <img
                        src={getImageUrl(file.url)}
                        alt=""
                        className="h-full w-full object-cover"
                      />
                    ) : (
                      <File className="h-5 w-5 text-muted-foreground" />
                    )}
                  </div>
                  <div>
                    <p className="text-sm font-medium group-hover:text-accent transition-colors">
                      {file.filename}
                    </p>
                    <p className="text-xs text-muted-foreground">
                      {file.mime_type} •{" "}
                      {new Date(file.created_at).toLocaleDateString()}
                    </p>
                  </div>
                </div>
                <div className="flex items-center gap-4">
                  <span className="text-xs text-muted-foreground font-mono">
                    {formatSize(file.size)}
                  </span>
                  <div className="flex gap-1">
                    <Button
                      variant="ghost"
                      size="icon"
                      className="h-8 w-8 text-muted-foreground hover:text-accent"
                      onClick={() => {
                        navigator.clipboard.writeText(getImageUrl(file.url));
                        toast.success("Link copied to clipboard");
                      }}
                    >
                      <ExternalLink className="h-4 w-4" />
                    </Button>
                    <Button
                      variant="ghost"
                      size="icon"
                      className="h-8 w-8 text-muted-foreground hover:text-destructive"
                      onClick={() => handleDelete(file.id)}
                    >
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
