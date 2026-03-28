"use client";

import { useState } from "react";
import dynamic from "next/dynamic";
import { useRouter } from "@/i18n/routing";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { Badge } from "@/components/ui/badge";
import { 
  ArrowLeft, 
  Loader2, 
  ImageIcon, 
  Settings, 
  Type, 
  X,
  Bold,
  Italic,
  Link as LinkIcon,
  Quote,
  List as ListIcon,
  Rocket
} from "lucide-react";
import { Post } from "@/lib/types";
import { cn } from "@/lib/utils";
import { useCreatePost, useUpdatePost } from "@/services/post-service";
import { useUploadMedia } from "@/services/media-service";

// Dynamically import Markdown Editor
const MDEditor = dynamic(
  () => import("@uiw/react-md-editor"),
  { ssr: false }
);

interface PostEditorProps {
  initialData?: Partial<Post>;
  mode: 'create' | 'edit';
}

import { useParams } from "next/navigation";

export function PostEditor({ initialData, mode }: PostEditorProps) {
  const router = useRouter();
  const params = useParams();
  const locale = (params?.locale as string) || 'zh-CN';
  
  const [formData, setFormData] = useState({
    title: initialData?.title || "",
    slug: initialData?.slug || "",
    content: initialData?.content || "",
    status: initialData?.status ?? 0, 
    cover: initialData?.cover || "",
    tags: (initialData?.tags || []).map(t => typeof t === 'string' ? t : t.name),
    excerpt: "", 
  });

  const createPost = useCreatePost();
  const updatePost = useUpdatePost(initialData?.id || 0);
  const uploadMedia = useUploadMedia();

  const isSubmitting = createPost.isPending || updatePost.isPending;
  const isPublished = formData.status === 1;

  // Manual slug generation when title changes in create mode
  const handleTitleChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    const title = e.target.value;
    const updates: Record<string, unknown> = { title };

    if (mode === 'create' && !initialData?.slug) {
      updates.slug = title
        .toLowerCase()
        .replace(/[^a-z0-9]+/g, '-')
        .replace(/(^-|-$)+/g, '');
    }

    setFormData(prev => ({ ...prev, ...updates }));
  };

  const handleChange = (key: string, value: unknown) => {
    setFormData(prev => ({ ...prev, [key]: value }));
  };
  const handleSave = async (statusOverride?: number) => {
    const finalStatus = statusOverride ?? formData.status;
    const submissionData = {
      ...formData,
      status: finalStatus,
    };

    if (mode === 'create') {
      createPost.mutate(submissionData, {
        onSuccess: (newPost) => {
          if (newPost?.id) {
            // Explicitly ensuring locale is not 'undefined' in any string interpolation
            const safeLocale = (locale && locale !== 'undefined') ? locale : 'zh-CN';
            // Next-intl router.push usually doesn't need the locale prefix, 
            // but we ensure the path is clean.
            router.push(`/admin/posts/${newPost.id}/edit`);
          } else {
            router.push('/admin/posts');
          }
        }
      });
    } else {
      updatePost.mutate(submissionData);
    }
  };

  const handleCoverUpload = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      uploadMedia.mutate(file, {
        onSuccess: (asset) => {
          handleChange('cover', asset.url);
        }
      });
    }
  };

  const [tagInput, setTagInput] = useState("");
  const handleAddTag = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && tagInput.trim()) {
      e.preventDefault();
      if (!formData.tags.includes(tagInput.trim())) {
        setFormData(prev => ({ ...prev, tags: [...prev.tags, tagInput.trim()] }));
      }
      setTagInput("");
    }
  };

  const removeTag = (tag: string) => {
    setFormData(prev => ({ ...prev, tags: formData.tags.filter((t: string) => t !== tag) }));
  };

  const wordCount = formData.content.split(/\s+/).filter(w => w.length > 0).length;

  return (
    <div className="flex h-screen bg-[#0d0b14] text-slate-200 overflow-hidden font-sans selection:bg-indigo-500/30">

      {/* 1. Left Vertical Toolbar */}
      <aside className="w-16 flex flex-col items-center py-6 border-r border-white/5 gap-8 shrink-0">
        <Button 
          variant="ghost" 
          size="icon"
          onClick={() => router.back()}
          className="text-slate-500 hover:text-white hover:bg-white/5 rounded-xl transition-all"
        >
          <ArrowLeft className="h-5 w-5" />
        </Button>

        <div className="flex flex-col gap-4">
          {[Type, Bold, Italic, LinkIcon, Quote, ListIcon, ImageIcon].map((Icon, i) => (
            <Button 
              key={i}
              variant="ghost" 
              size="icon"
              className="text-slate-500 hover:text-white hover:bg-white/5 transition-all"
            >
              <Icon className="h-5 w-5" />
            </Button>
          ))}
        </div>
      </aside>

      {/* 2. Main Content Area */}
      <main className="flex-1 flex flex-col overflow-hidden relative">

        {/* Status Bar */}
        <div className="flex items-center justify-between px-12 py-6 shrink-0">
          <div className="flex items-center gap-2 text-[10px] font-bold tracking-widest text-slate-500 uppercase">
            <span className={cn(isPublished ? "text-emerald-500" : "text-yellow-500")}>
              {isPublished ? "PUBLISHED" : "DRAFT"}
            </span>
            <span className="opacity-30">•</span>
            <span>{isSubmitting ? "SAVING..." : "READY"}</span>
          </div>
          <div className="text-[10px] font-bold tracking-widest text-slate-500 uppercase">
            {wordCount.toLocaleString()} WORDS
          </div>
        </div>

        {/* Editor Canvas */}
        <div className="flex-1 overflow-y-auto custom-scrollbar px-12 pb-24">
          <div className="max-w-4xl mx-auto">
            {/* Title Section */}
            <div className="mb-12 group">
              <textarea
                placeholder="Untitled Story"
                className="w-full bg-transparent border-none text-6xl font-bold text-white placeholder:text-white/10 focus:outline-none resize-none leading-tight mb-4"
                rows={1}
                value={formData.title}
                onChange={handleTitleChange}
              />
              <div className="h-1 w-24 bg-indigo-500 rounded-full transition-all group-focus-within:w-48" />
            </div>

            {/* Markdown Editor */}
            <div data-color-mode="dark" className="prose prose-invert max-w-none prose-p:text-slate-400 prose-p:text-lg prose-p:leading-relaxed prose-headings:text-white">
              <MDEditor
                value={formData.content}
                onChange={(val) => handleChange('content', val || "")}
                height={800}
                className="bg-transparent border-none shadow-none !bg-transparent md-editor-custom"
                visibleDragbar={false}
                preview="edit"
                hideToolbar={true} 
                textareaProps={{
                  placeholder: "Write your masterpiece here...",
                  className: "text-xl leading-relaxed text-slate-400 font-sans focus:outline-none"
                }}
                style={{ backgroundColor: 'transparent' }}
              />
            </div>
          </div>
        </div>
      </main>

      {/* 3. Right Sidebar (Settings) */}
      <aside className="w-80 border-l border-white/5 bg-[#0d0b14] flex flex-col shrink-0">
        {/* Top Buttons */}
        <div className="p-6 flex items-center justify-between gap-4">
          <Button 
            variant="ghost" 
            className="text-slate-400 hover:text-white text-[10px] font-bold uppercase tracking-widest px-0"
            onClick={() => handleSave(0)}
            disabled={isSubmitting || isPublished}
          >
            {isSubmitting ? <Loader2 className="h-3 w-3 animate-spin mr-2" /> : isPublished ? "Saved" : "Save Draft"}
          </Button>

          <Button 
            className="bg-indigo-600 hover:bg-indigo-700 text-white border-0 px-6 rounded-lg font-bold text-[10px] uppercase tracking-widest shadow-lg shadow-indigo-600/20 transition-all hover:translate-y-[-2px]"
            onClick={() => handleSave(1)}
            disabled={isSubmitting}
          >
            {isSubmitting ? <Loader2 className="h-3 w-3 animate-spin mr-2" /> : <Rocket className="h-3 w-3 mr-2" />}
            {isPublished ? "Update" : "Publish"}
          </Button>
        </div>

        <div className="flex-1 overflow-y-auto p-6 space-y-10 custom-scrollbar">

          {/* Post Settings */}
          <div className="space-y-6">
            <h3 className="text-[10px] font-bold text-slate-600 uppercase tracking-[0.2em] flex items-center gap-2">
              <Settings className="w-3 h-3" />
              Settings
            </h3>

            <div className="space-y-3">
              <Label className="text-[11px] font-bold text-slate-400">URL Slug</Label>
              <div className="bg-white/5 rounded-lg border border-white/5 px-4 py-3 flex items-center gap-3 focus-within:border-indigo-500/50 transition-colors">
                <LinkIcon className="w-3.5 h-3.5 text-slate-600" />
                <input 
                  className="bg-transparent border-none focus:outline-none w-full text-xs font-medium text-slate-300 placeholder:text-slate-700"
                  value={formData.slug}
                  onChange={(e) => handleChange('slug', e.target.value)}
                />
              </div>
            </div>

            <div className="space-y-3">
              <Label className="text-[11px] font-bold text-slate-400">Tags</Label>
              <div className="bg-white/5 rounded-lg border border-white/5 p-2 flex flex-wrap gap-2 focus-within:border-indigo-500/50 transition-colors">
                {formData.tags.map((tag: string) => (
                  <Badge 
                    key={tag} 
                    className="bg-indigo-500/10 text-indigo-400 hover:bg-indigo-500/20 border-none px-3 py-1 text-[10px] font-bold flex items-center gap-1.5"
                  >
                    #{tag}
                    <X className="w-3 h-3 cursor-pointer opacity-50 hover:opacity-100" onClick={() => removeTag(tag)}/>
                  </Badge>
                ))}
                <input 
                  className="bg-transparent border-none focus:outline-none text-[11px] font-medium text-slate-500 px-2 py-1 min-w-[60px]"
                  placeholder="Add..."
                  value={tagInput}
                  onChange={(e) => setTagInput(e.target.value)}
                  onKeyDown={handleAddTag}
                />
              </div>
            </div>
          </div>

          {/* Cover Image */}
          <div className="space-y-6">
            <h3 className="text-[10px] font-bold text-slate-600 uppercase tracking-[0.2em] flex items-center gap-2">
              <ImageIcon className="w-3 h-3" />
              Cover Image
            </h3>

            <div className="group relative rounded-xl overflow-hidden bg-white/5 border border-white/5 transition-all hover:border-indigo-500/30 cursor-pointer aspect-video flex items-center justify-center">
              {formData.cover ? (
                <img src={formData.cover} alt="Cover" className="w-full h-full object-cover" />
              ) : (
                <div className="text-center p-6">
                  <div className="bg-white/5 p-3 rounded-full inline-block mb-3 text-slate-600 group-hover:text-indigo-400 transition-colors">
                    <ImageIcon className="w-6 h-6" />
                  </div>
                  <p className="text-[11px] text-slate-600">Click to upload cover</p>
                </div>
              )}
              {uploadMedia.isPending && (
                <div className="absolute inset-0 bg-black/60 flex items-center justify-center">
                  <Loader2 className="h-6 w-6 animate-spin text-indigo-500" />
                </div>
              )}
              <input 
                type="file" 
                className="absolute inset-0 opacity-0 cursor-pointer" 
                onChange={handleCoverUpload}
                accept="image/*"
              />
            </div>
          </div>
        </div>

        {/* Footer */}
        <div className="p-6 border-t border-white/5">
          <p className="text-[10px] text-slate-700 font-bold tracking-widest text-center uppercase">
            Kaldalis Engine v2.4.0
          </p>
        </div>
      </aside>

      <style jsx global>{`
        .md-editor-custom .w-md-editor-content {
          background-color: transparent !important;
        }
        .md-editor-custom .w-md-editor-text {
          background-color: transparent !important;
        }
        .md-editor-custom .w-md-editor-preview {
          background-color: transparent !important;
          border-left-color: rgba(255,255,255,0.05) !important;
        }
        .md-editor-custom {
          box-shadow: none !important;
        }
        .custom-scrollbar::-webkit-scrollbar {
          width: 4px;
        }
        .custom-scrollbar::-webkit-scrollbar-track {
          background: transparent;
        }
        .custom-scrollbar::-webkit-scrollbar-thumb {
          background: rgba(255, 255, 255, 0.05);
          border-radius: 10px;
        }
        .custom-scrollbar::-webkit-scrollbar-thumb:hover {
          background: rgba(255, 255, 255, 0.1);
        }
      `}</style>
    </div>
  );
}