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
  Rocket,
  ChevronRight
} from "lucide-react";
import { Post, PostStatus } from "@/lib/types";
import { cn } from "@/lib/utils";
import { useCreatePost, useUpdatePost } from "@/services/post-service";
import { useUploadMedia } from "@/services/media-service";
import { useParams } from "next/navigation";
import { motion } from "framer-motion";

// Dynamically import Markdown Editor
const MDEditor = dynamic(
  () => import("@uiw/react-md-editor"),
  { ssr: false }
);

interface PostEditorProps {
  initialData?: Partial<Post>;
  mode: 'create' | 'edit';
}

export function PostEditor({ initialData, mode }: PostEditorProps) {
  const router = useRouter();
  const params = useParams();
  const locale = (params?.locale as string) || 'zh-CN';
  
  const [formData, setFormData] = useState({
    title: initialData?.title || "",
    slug: initialData?.slug || "",
    content: initialData?.content || "",
    status: initialData?.status ?? PostStatus.DRAFT,
    cover: initialData?.cover || "",
    tags: (initialData?.tags || []).map(t => typeof t === 'string' ? t : t.name),
    excerpt: "",
  });

  const createPost = useCreatePost();
  const updatePost = useUpdatePost(initialData?.id || 0);
  const uploadMedia = useUploadMedia();

  const isSubmitting = createPost.isPending || updatePost.isPending;
  const isPublished = formData.status === PostStatus.PUBLISHED;

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
    <div className="flex h-screen bg-background text-foreground overflow-hidden font-sans selection:bg-accent/30 relative">
      
      {/* 1. Left Vertical Toolbar: Minimalist */}
      <aside className="w-20 flex flex-col items-center py-8 border-r border-border bg-white dark:bg-slate-900 z-20">
        <Button 
          variant="ghost" 
          size="icon"
          onClick={() => router.back()}
          className="text-muted-foreground hover:text-foreground rounded-full hover:bg-muted mb-10"
        >
          <ArrowLeft className="h-5 w-5" />
        </Button>

        <div className="flex flex-col gap-6">
          {[Type, Bold, Italic, LinkIcon, Quote, ListIcon, ImageIcon].map((Icon, i) => (
            <Button 
              key={i}
              variant="ghost" 
              size="icon"
              className="text-muted-foreground hover:text-accent transition-all rounded-xl"
            >
              <Icon className="h-5 w-5" />
            </Button>
          ))}
        </div>
      </aside>

      {/* 2. Main Content Area: Immersive Writing Space */}
      <main className="flex-1 flex flex-col overflow-hidden relative bg-background/50 backdrop-blur-sm">
        
        {/* Status Header */}
        <div className="flex items-center justify-between px-16 py-8 shrink-0">
          <div className="flex items-center gap-4 text-[10px] font-bold tracking-[0.2em] text-muted-foreground uppercase">
            <span className={cn(
              "px-3 py-1 rounded-full",
              isPublished ? "bg-accent/10 text-accent" : "bg-muted text-muted-foreground"
            )}>
              {isPublished ? "Published" : "Draft Mode"}
            </span>
            <span className="opacity-20">/</span>
            <span className="flex items-center gap-2">
              {isSubmitting ? <Loader2 className="h-3 w-3 animate-spin" /> : <div className="w-1.5 h-1.5 rounded-full bg-emerald-500" />}
              {isSubmitting ? "Persisting..." : "All changes saved"}
            </span>
          </div>
          
          <div className="text-[10px] font-bold tracking-[0.2em] text-muted-foreground uppercase bg-muted/50 px-4 py-1.5 rounded-full">
            {wordCount.toLocaleString()} Words
          </div>
        </div>

        {/* Editor Canvas */}
        <div className="flex-1 overflow-y-auto custom-scrollbar px-16 pb-32">
          <div className="max-w-3xl mx-auto">
            {/* Title Section */}
            <motion.div 
              initial={{ opacity: 0, y: 10 }}
              animate={{ opacity: 1, y: 0 }}
              className="mb-16 group"
            >
              <textarea
                placeholder="Story Title..."
                className="w-full bg-transparent border-none text-6xl md:text-7xl font-serif font-medium text-foreground placeholder:text-foreground/5 focus:outline-none resize-none leading-[1.1] mb-6"
                rows={1}
                value={formData.title}
                onChange={handleTitleChange}
              />
              <div className="h-px w-20 bg-accent transition-all group-focus-within:w-full opacity-30" />
            </motion.div>

            {/* Markdown Editor */}
            <motion.div 
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              transition={{ delay: 0.2 }}
              data-color-mode="light" 
              className="prose prose-lg dark:prose-invert max-w-none prose-p:text-foreground/80 prose-p:leading-relaxed prose-headings:font-serif prose-headings:font-medium"
            >
              <MDEditor
                value={formData.content}
                onChange={(val) => handleChange('content', val || "")}
                height={1000}
                className="bg-transparent border-none shadow-none !bg-transparent md-editor-minimal"
                visibleDragbar={false}
                preview="edit"
                hideToolbar={true} 
                textareaProps={{
                  placeholder: "Begin your creative journey...",
                  className: "text-xl leading-relaxed text-foreground/80 font-sans focus:outline-none"
                }}
                style={{ backgroundColor: 'transparent' }}
              />
            </motion.div>
          </div>
        </div>
      </main>

      {/* 3. Right Sidebar: Refined Controls */}
      <aside className="w-80 border-l border-border bg-white dark:bg-slate-900 flex flex-col shrink-0 z-20">
        <div className="p-8 space-y-10">
          {/* Actions */}
          <div className="flex items-center gap-3">
            <Button
              variant="outline"
              className="flex-1 rounded-full border-border hover:bg-muted text-xs font-bold uppercase tracking-widest h-12"
              onClick={() => handleSave(PostStatus.DRAFT)}
              disabled={isSubmitting}
            >
              {isPublished ? "To Draft" : "Save"}
            </Button>

            <Button
              className="flex-1 rounded-full bg-primary text-primary-foreground hover:bg-primary/90 shadow-xl shadow-primary/10 text-xs font-bold uppercase tracking-widest h-12"
              onClick={() => handleSave(PostStatus.PUBLISHED)}
              disabled={isSubmitting}
            >
              {isSubmitting ? <Loader2 className="h-3 w-3 animate-spin mr-2" /> : <Rocket className="h-3 w-3 mr-2" />}
              {isPublished ? "Update" : "Publish"}
            </Button>
          </div>

          <div className="space-y-12 overflow-y-auto custom-scrollbar pr-2">
            
            {/* Metadata */}
            <section className="space-y-6">
              <h3 className="text-[10px] font-bold text-muted-foreground uppercase tracking-[0.3em] flex items-center gap-2">
                <Settings className="w-3 h-3" />
                Configuration
              </h3>
              
              <div className="space-y-4">
                <div className="space-y-2">
                  <Label className="text-[11px] font-bold text-foreground/60 uppercase tracking-wider ml-1">URL Extension</Label>
                  <div className="bg-muted/50 rounded-2xl border border-border px-4 py-3 flex items-center gap-3 focus-within:border-accent/50 transition-all">
                    <LinkIcon className="w-3.5 h-3.5 text-muted-foreground" />
                    <input 
                      className="bg-transparent border-none focus:outline-none w-full text-xs font-medium placeholder:text-muted-foreground/30"
                      value={formData.slug}
                      onChange={(e) => handleChange('slug', e.target.value)}
                    />
                  </div>
                </div>

                <div className="space-y-2">
                  <Label className="text-[11px] font-bold text-foreground/60 uppercase tracking-wider ml-1">Taxonomy</Label>
                  <div className="bg-muted/50 rounded-2xl border border-border p-3 flex flex-wrap gap-2 focus-within:border-accent/50 transition-all min-h-[100px]">
                    {formData.tags.map((tag: string) => (
                      <Badge 
                        key={tag} 
                        className="bg-accent text-white hover:bg-accent/90 border-none px-3 py-1 text-[10px] font-bold rounded-full flex items-center gap-1.5 shadow-lg shadow-accent/10"
                      >
                        {tag}
                        <X className="w-3 h-3 cursor-pointer opacity-70 hover:opacity-100" onClick={() => removeTag(tag)}/>
                      </Badge>
                    ))}
                    <input 
                      className="bg-transparent border-none focus:outline-none text-[11px] font-medium text-foreground px-2 py-1 min-w-[80px]"
                      placeholder="Add tag..."
                      value={tagInput}
                      onChange={(e) => setTagInput(e.target.value)}
                      onKeyDown={handleAddTag}
                    />
                  </div>
                </div>
              </div>
            </section>

            {/* Visuals */}
            <section className="space-y-6">
              <h3 className="text-[10px] font-bold text-muted-foreground uppercase tracking-[0.3em] flex items-center gap-2">
                <ImageIcon className="w-3 h-3" />
                Cover Visual
              </h3>
              
              <div className="group relative rounded-3xl overflow-hidden bg-muted/50 border-2 border-dashed border-border transition-all hover:border-accent/30 cursor-pointer aspect-[4/3] flex items-center justify-center">
                {formData.cover ? (
                  <img src={formData.cover} alt="Cover" className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500" />
                ) : (
                  <div className="text-center p-6 space-y-3">
                    <div className="bg-background w-12 h-12 rounded-full flex items-center justify-center mx-auto text-muted-foreground group-hover:text-accent transition-colors shadow-sm">
                      <ImageIcon className="w-5 h-5" />
                    </div>
                    <p className="text-[10px] font-bold uppercase tracking-widest text-muted-foreground/60">Upload Cover</p>
                  </div>
                )}
                {uploadMedia.isPending && (
                  <div className="absolute inset-0 bg-background/80 backdrop-blur-sm flex items-center justify-center">
                    <Loader2 className="h-6 w-6 animate-spin text-accent" />
                  </div>
                )}
                <input 
                  type="file" 
                  className="absolute inset-0 opacity-0 cursor-pointer" 
                  onChange={handleCoverUpload}
                  accept="image/*"
                />
              </div>
            </section>
          </div>
        </div>

        {/* Footer */}
        <div className="mt-auto p-8 border-t border-border flex items-center justify-between">
          <div className="flex items-center gap-2">
            <div className="w-2 h-2 rounded-full bg-accent animate-pulse" />
            <span className="text-[10px] font-bold uppercase tracking-widest text-muted-foreground">Active</span>
          </div>
          <p className="text-[10px] text-muted-foreground/40 font-bold tracking-widest uppercase">
            Kaldalis v2.4
          </p>
        </div>
      </aside>

      <style jsx global>{`
        .md-editor-minimal .w-md-editor-content {
          background-color: transparent !important;
        }
        .md-editor-minimal .w-md-editor-text {
          background-color: transparent !important;
        }
        .md-editor-minimal .w-md-editor-preview {
          background-color: transparent !important;
          border-left: 1px solid var(--border) !important;
        }
        .md-editor-minimal {
          box-shadow: none !important;
        }
        .custom-scrollbar::-webkit-scrollbar {
          width: 4px;
        }
        .custom-scrollbar::-webkit-scrollbar-track {
          background: transparent;
        }
        .custom-scrollbar::-webkit-scrollbar-thumb {
          background: var(--border);
          border-radius: 10px;
        }
        .custom-scrollbar::-webkit-scrollbar-thumb:hover {
          background: var(--accent);
        }
      `}</style>
    </div>
  );
}
