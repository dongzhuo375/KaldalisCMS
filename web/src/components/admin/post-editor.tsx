"use client";

import { useState, useEffect, useCallback } from "react";
import dynamic from "next/dynamic";
import { useRouter } from "@/i18n/routing";
import { useTheme } from "next-themes";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import {
  ArrowLeft,
  Loader2,
  ImageIcon,
  X,
  Link as LinkIcon,
  Rocket,
  Maximize2,
  Minimize2,
  Trash2,
  Save,
  FileText,
} from "lucide-react";
import { Post, PostStatus } from "@/lib/types";
import { cn } from "@/lib/utils";
import { useCreatePost, useUpdatePost } from "@/services/post-service";
import { useUploadMedia } from "@/services/media-service";

const MDEditor = dynamic(() => import("@uiw/react-md-editor"), { ssr: false });

interface PostEditorProps {
  initialData?: Partial<Post>;
  mode: "create" | "edit";
}

export function PostEditor({ initialData, mode }: PostEditorProps) {
  const router = useRouter();
  const { resolvedTheme } = useTheme();
  const [focusMode, setFocusMode] = useState(false);
  const [mounted, setMounted] = useState(false);

  const [formData, setFormData] = useState({
    title: initialData?.title || "",
    slug: initialData?.slug || "",
    content: initialData?.content || "",
    status: initialData?.status ?? PostStatus.DRAFT,
    cover: initialData?.cover || "",
    tags: (initialData?.tags || []).map((t) =>
      typeof t === "string" ? t : t.name
    ),
    excerpt: "",
  });

  const createPost = useCreatePost();
  const updatePost = useUpdatePost(initialData?.id || 0);
  const uploadMedia = useUploadMedia();

  const isSubmitting = createPost.isPending || updatePost.isPending;
  const isPublished = formData.status === PostStatus.PUBLISHED;

  useEffect(() => setMounted(true), []);

  const handleTitleChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    const title = e.target.value;
    const updates: Record<string, unknown> = { title };

    if (mode === "create" && !initialData?.slug) {
      updates.slug = title
        .toLowerCase()
        .replace(/[^a-z0-9]+/g, "-")
        .replace(/(^-|-$)+/g, "");
    }

    setFormData((prev) => ({ ...prev, ...updates }));

    e.target.style.height = "auto";
    e.target.style.height = e.target.scrollHeight + "px";
  };

  const handleChange = (key: string, value: unknown) => {
    setFormData((prev) => ({ ...prev, [key]: value }));
  };

  const handleSave = useCallback(
    async (statusOverride?: number) => {
      const finalStatus = statusOverride ?? formData.status;
      const submissionData = {
        ...formData,
        status: finalStatus,
      };

      if (mode === "create") {
        createPost.mutate(submissionData, {
          onSuccess: (newPost) => {
            if (newPost?.id) {
              router.push(`/admin/posts/${newPost.id}/edit`);
            } else {
              router.push("/admin/posts");
            }
          },
        });
      } else {
        updatePost.mutate(submissionData);
      }
    },
    [formData, mode, createPost, updatePost, router]
  );

  const handleCoverUpload = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      uploadMedia.mutate(file, {
        onSuccess: (asset) => {
          handleChange("cover", asset.url);
        },
      });
    }
  };

  const removeCover = () => handleChange("cover", "");

  const [tagInput, setTagInput] = useState("");
  const handleAddTag = (e: React.KeyboardEvent) => {
    if (e.key === "Enter" && tagInput.trim()) {
      e.preventDefault();
      if (!formData.tags.includes(tagInput.trim())) {
        setFormData((prev) => ({
          ...prev,
          tags: [...prev.tags, tagInput.trim()],
        }));
      }
      setTagInput("");
    }
  };

  const removeTag = (tag: string) => {
    setFormData((prev) => ({
      ...prev,
      tags: prev.tags.filter((t: string) => t !== tag),
    }));
  };

  // Keyboard shortcuts
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if ((e.metaKey || e.ctrlKey) && e.key === "s") {
        e.preventDefault();
        handleSave();
      }
      if (e.key === "Escape" && focusMode) {
        setFocusMode(false);
      }
    };
    window.addEventListener("keydown", handleKeyDown);
    return () => window.removeEventListener("keydown", handleKeyDown);
  }, [handleSave, focusMode]);

  const wordCount = formData.content.trim()
    ? formData.content.trim().split(/\s+/).length
    : 0;
  const charCount = formData.content.length;

  const colorMode = resolvedTheme === "dark" ? "dark" : "light";

  // -- Focus mode: fullscreen overlay --
  if (focusMode) {
    return (
      <div
        className="fixed inset-0 z-50 flex flex-col bg-background"
        data-color-mode={colorMode}
      >
        {/* Focus mode top bar */}
        <div className="flex items-center justify-between px-6 py-3 border-b border-border/50 bg-background/80 backdrop-blur-sm">
          <div className="flex items-center gap-3 text-sm text-muted-foreground">
            <FileText className="h-4 w-4" />
            <span className="font-medium truncate max-w-[300px]">
              {formData.title || "Untitled"}
            </span>
            <span className="text-xs opacity-60">
              {charCount} chars
            </span>
          </div>
          <div className="flex items-center gap-2">
            <Button
              variant="ghost"
              size="sm"
              onClick={() => handleSave(PostStatus.DRAFT)}
              disabled={isSubmitting}
              className="text-muted-foreground hover:text-foreground"
            >
              {isSubmitting ? (
                <Loader2 className="h-4 w-4 animate-spin" />
              ) : (
                <Save className="h-4 w-4 mr-1.5" />
              )}
              Save
            </Button>
            <Button
              variant="ghost"
              size="icon"
              onClick={() => setFocusMode(false)}
              className="text-muted-foreground hover:text-foreground"
            >
              <Minimize2 className="h-4 w-4" />
            </Button>
          </div>
        </div>

        {/* Focus editor area */}
        <div className="flex-1 overflow-y-auto">
          <div className="max-w-3xl mx-auto px-6 py-10">
            <textarea
              placeholder="Title"
              className="w-full bg-transparent border-none font-serif font-medium text-4xl text-foreground placeholder:text-muted-foreground/25 focus:outline-none resize-none leading-tight mb-8"
              rows={1}
              value={formData.title}
              onChange={handleTitleChange}
            />
            <div data-color-mode={colorMode}>
              <MDEditor
                value={formData.content}
                onChange={(val) => handleChange("content", val || "")}
                height="calc(100vh - 280px)"
                visibleDragbar={false}
                preview="live"
                textareaProps={{
                  placeholder: "Write your content...",
                }}
              />
            </div>
          </div>
        </div>
      </div>
    );
  }

  // -- Normal mode: integrated in admin layout --
  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-4">
          <Button
            variant="ghost"
            size="icon"
            onClick={() => router.back()}
            className="rounded-xl text-muted-foreground hover:text-foreground"
          >
            <ArrowLeft className="h-4 w-4" />
          </Button>
          <div>
            <h1 className="text-2xl font-serif font-bold tracking-tight">
              {mode === "create" ? "New Post" : "Edit Post"}
            </h1>
            <p className="text-sm text-muted-foreground mt-0.5">
              {mode === "create"
                ? "Create a new article"
                : `Editing "${initialData?.title || "Untitled"}"`}
            </p>
          </div>
        </div>

        <div className="flex items-center gap-2">
          <Button
            variant="ghost"
            size="icon"
            onClick={() => setFocusMode(true)}
            className="rounded-xl text-muted-foreground hover:text-foreground"
            title="Focus mode"
          >
            <Maximize2 className="h-4 w-4" />
          </Button>

          <Button
            variant="outline"
            size="sm"
            onClick={() => handleSave(PostStatus.DRAFT)}
            disabled={isSubmitting}
            className="rounded-xl"
          >
            <Save className="h-4 w-4 mr-1.5" />
            {isPublished ? "Unpublish" : "Save Draft"}
          </Button>

          <Button
            size="sm"
            onClick={() => handleSave(PostStatus.PUBLISHED)}
            disabled={isSubmitting}
            className="rounded-xl"
          >
            {isSubmitting ? (
              <Loader2 className="h-4 w-4 animate-spin mr-1.5" />
            ) : (
              <Rocket className="h-4 w-4 mr-1.5" />
            )}
            {isPublished ? "Update" : "Publish"}
          </Button>
        </div>
      </div>

      {/* Two column layout */}
      <div className="flex gap-6 items-start">
        {/* Main editor column */}
        <div className="flex-1 min-w-0 space-y-4">
          {/* Title card */}
          <div className="rounded-2xl border border-border bg-white/60 dark:bg-white/[0.03] backdrop-blur-sm p-6">
            <textarea
              placeholder="Post title..."
              className="w-full bg-transparent border-none font-serif font-medium text-3xl text-foreground placeholder:text-muted-foreground/25 focus:outline-none resize-none leading-tight"
              rows={1}
              value={formData.title}
              onChange={handleTitleChange}
            />
          </div>

          {/* Markdown editor card - no backdrop-blur/transforms to avoid cursor misalignment */}
          <div
            className="rounded-2xl border border-border bg-white dark:bg-background overflow-hidden"
            data-color-mode={colorMode}
          >
            {mounted && (
              <MDEditor
                value={formData.content}
                onChange={(val) => handleChange("content", val || "")}
                height={560}
                visibleDragbar={false}
                preview="live"
                textareaProps={{
                  placeholder: "Write your content...",
                }}
              />
            )}
          </div>

          {/* Stats bar */}
          <div className="flex items-center gap-4 px-2 text-xs text-muted-foreground">
            <span>{charCount} characters</span>
            <span>{wordCount} words</span>
            <span className="ml-auto flex items-center gap-1.5">
              <kbd className="px-1.5 py-0.5 bg-muted rounded text-[10px] font-mono">
                Ctrl+S
              </kbd>
              save
            </span>
          </div>
        </div>

        {/* Sidebar */}
        <div className="w-72 shrink-0 space-y-4">
          {/* Status */}
          <div className="rounded-2xl border border-border bg-white/60 dark:bg-white/[0.03] backdrop-blur-sm p-5">
            <label className="text-xs font-bold uppercase tracking-wider text-muted-foreground">
              Status
            </label>
            <div className="mt-3 flex items-center gap-2">
              <div
                className={cn(
                  "w-2 h-2 rounded-full",
                  isPublished
                    ? "bg-emerald-500"
                    : "bg-amber-400"
                )}
              />
              <span className="text-sm font-medium">
                {isPublished ? "Published" : "Draft"}
              </span>
            </div>
          </div>

          {/* Cover Image */}
          <div className="rounded-2xl border border-border bg-white/60 dark:bg-white/[0.03] backdrop-blur-sm p-5">
            <label className="text-xs font-bold uppercase tracking-wider text-muted-foreground">
              Cover Image
            </label>
            <div className="mt-3 relative rounded-xl overflow-hidden bg-muted/50 border border-border/50 aspect-video flex items-center justify-center group cursor-pointer hover:border-accent/30 transition-colors">
              {formData.cover ? (
                <>
                  <img
                    src={formData.cover}
                    alt="Cover"
                    className="w-full h-full object-cover"
                  />
                  <button
                    onClick={removeCover}
                    className="absolute top-2 right-2 p-1.5 rounded-lg bg-black/50 text-white opacity-0 group-hover:opacity-100 transition-opacity hover:bg-black/70"
                  >
                    <Trash2 className="w-3 h-3" />
                  </button>
                </>
              ) : (
                <div className="text-center p-4">
                  <ImageIcon className="w-5 h-5 mx-auto text-muted-foreground/50 mb-1.5" />
                  <span className="text-xs text-muted-foreground/60">
                    Click to upload
                  </span>
                </div>
              )}
              {uploadMedia.isPending && (
                <div className="absolute inset-0 bg-background/80 flex items-center justify-center">
                  <Loader2 className="h-5 w-5 animate-spin text-accent" />
                </div>
              )}
              {!formData.cover && (
                <input
                  type="file"
                  className="absolute inset-0 opacity-0 cursor-pointer"
                  onChange={handleCoverUpload}
                  accept="image/*"
                />
              )}
            </div>
          </div>

          {/* Tags */}
          <div className="rounded-2xl border border-border bg-white/60 dark:bg-white/[0.03] backdrop-blur-sm p-5">
            <label className="text-xs font-bold uppercase tracking-wider text-muted-foreground">
              Tags
            </label>
            <div className="mt-3 rounded-xl border border-border/50 p-3 flex flex-wrap gap-1.5 min-h-[72px] focus-within:border-accent/30 transition-colors bg-background/50">
              {formData.tags.map((tag: string) => (
                <Badge
                  key={tag}
                  variant="secondary"
                  className="px-2.5 py-1 text-xs flex items-center gap-1 rounded-lg"
                >
                  {tag}
                  <X
                    className="w-3 h-3 cursor-pointer hover:text-destructive transition-colors"
                    onClick={() => removeTag(tag)}
                  />
                </Badge>
              ))}
              <input
                className="bg-transparent border-none focus:outline-none text-sm flex-1 min-w-[60px] placeholder:text-muted-foreground/40"
                placeholder="Add tag..."
                value={tagInput}
                onChange={(e) => setTagInput(e.target.value)}
                onKeyDown={handleAddTag}
              />
            </div>
          </div>

          {/* Slug */}
          <div className="rounded-2xl border border-border bg-white/60 dark:bg-white/[0.03] backdrop-blur-sm p-5">
            <label className="text-xs font-bold uppercase tracking-wider text-muted-foreground">
              Slug
            </label>
            <div className="mt-3 rounded-xl border border-border/50 px-3 py-2.5 flex items-center gap-2 focus-within:border-accent/30 transition-colors bg-background/50">
              <LinkIcon className="w-4 h-4 text-muted-foreground/50 shrink-0" />
              <input
                className="bg-transparent border-none focus:outline-none text-sm w-full placeholder:text-muted-foreground/40 font-mono"
                value={formData.slug}
                onChange={(e) => handleChange("slug", e.target.value)}
                placeholder="post-url-slug"
              />
            </div>
          </div>
        </div>
      </div>

      <style jsx global>{`
        /* MDEditor theme integration */
        .w-md-editor {
          border: none !important;
          border-radius: 0 !important;
          box-shadow: none !important;
          background-color: transparent !important;
        }
        .w-md-editor-toolbar {
          background-color: transparent !important;
          border-bottom: 1px solid var(--border) !important;
          padding: 8px 12px !important;
          min-height: 40px;
        }
        .w-md-editor-toolbar ul > li > button {
          color: var(--muted-foreground) !important;
          border-radius: 6px !important;
        }
        .w-md-editor-toolbar ul > li > button:hover {
          color: var(--foreground) !important;
          background-color: var(--muted) !important;
        }
        .w-md-editor-toolbar ul > li > button.active {
          color: var(--foreground) !important;
          background-color: var(--muted) !important;
        }
        .w-md-editor-content {
          background-color: transparent !important;
        }
        .w-md-editor-text {
          background-color: transparent !important;
        }
        /* DO NOT override font-size, line-height, padding, or font-family
           on .w-md-editor-text-input / .w-md-editor-text-pre —
           MDEditor relies on pixel-perfect alignment between the textarea
           and the syntax-highlighted pre overlay. Any override breaks the cursor. */
        .w-md-editor-preview {
          background-color: transparent !important;
          padding: 20px !important;
          border-left: 1px solid var(--border) !important;
        }
        .wmde-markdown {
          background-color: transparent !important;
          font-size: 15px !important;
          line-height: 1.75 !important;
        }
        .wmde-markdown hr {
          border-color: var(--border) !important;
        }
        /* Dark mode adjustments */
        [data-color-mode="dark"] .w-md-editor {
          color: var(--foreground) !important;
        }
        [data-color-mode="dark"] .w-md-editor-text-input {
          color: var(--foreground) !important;
          -webkit-text-fill-color: var(--foreground) !important;
        }
      `}</style>
    </div>
  );
}
