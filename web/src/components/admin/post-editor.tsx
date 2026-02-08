"use client";

import { useState, useEffect } from "react";
import dynamic from "next/dynamic";
import { useRouter } from "@/i18n/routing";
import { useTranslations } from 'next-intl';
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { Badge } from "@/components/ui/badge";
import { 
  ArrowLeft, 
  Loader2, 
  Globe, 
  ImageIcon, 
  Settings, 
  Type, 
  Send, 
  Save,
  X,
  Bold,
  Italic,
  Link as LinkIcon,
  Quote,
  List as ListIcon,
  Rocket,
  Search
} from "lucide-react";
import { Post } from "@/lib/types";
import { cn } from "@/lib/utils";

// Dynamically import Markdown Editor
const MDEditor = dynamic(
  () => import("@uiw/react-md-editor"),
  { ssr: false }
);

interface PostEditorProps {
  initialData?: Partial<Post>;
  onSubmit: (data: any) => Promise<void>;
  isSubmitting: boolean;
  mode: 'create' | 'edit';
}

export function PostEditor({ initialData, onSubmit, isSubmitting, mode }: PostEditorProps) {
  const t = useTranslations('admin');
  const router = useRouter();
  
  const [formData, setFormData] = useState({
    title: initialData?.title || "",
    slug: initialData?.slug || "",
    content: initialData?.content || "",
    status: initialData?.status?.toString() || "0", 
    cover: initialData?.cover || "",
    tags: ["Design", "UX"] as string[], // Mock tags from design
    excerpt: initialData?.content?.substring(0, 150) || "Exploring the shift towards calm technology and immersive interfaces in modern web development.", 
  });

  // Derived state
  const wordCount = formData.content.split(/\s+/).filter(w => w.length > 0).length;
  const isPublished = formData.status === "1";

  // Auto-generate slug
  useEffect(() => {
    if (mode === 'create' && !initialData?.slug && formData.title) {
      const slug = formData.title
        .toLowerCase()
        .replace(/[^a-z0-9]+/g, '-')
        .replace(/(^-|-$)+/g, '');
      setFormData(prev => ({ ...prev, slug }));
    }
  }, [formData.title, mode, initialData?.slug]);

  const handleChange = (key: string, value: any) => {
    setFormData(prev => ({ ...prev, [key]: value }));
  };

  const handleSubmit = async (statusOverride?: number) => {
    const newStatus = statusOverride ?? parseInt(formData.status);
    const submissionData = {
      ...formData,
      status: newStatus,
      category_id: null, 
      tags: [] 
    };
    
    await onSubmit(submissionData);
    
    // 如果没有报错并走到了这一步，说明提交成功，更新本地状态
    setFormData(prev => ({ ...prev, status: newStatus.toString() }));
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
    setFormData(prev => ({ ...prev, tags: prev.tags.filter(t => t !== tag) }));
  };

  return (
    <div className="flex h-screen bg-[#0d0b14] text-slate-200 overflow-hidden font-sans selection:bg-[#ad2bee]/30">
      
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
            <span>{isPublished ? "PUBLISHED" : "DRAFT"}</span>
            <span className="opacity-30">•</span>
            <span>SAVED 2M AGO</span>
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
                placeholder="The Future of Immersive"
                className="w-full bg-transparent border-none text-6xl font-bold text-white placeholder:text-white/10 focus:outline-none resize-none leading-tight mb-4"
                rows={2}
                value={formData.title}
                onChange={(e) => handleChange('title', e.target.value)}
              />
              <div className="h-1.5 w-24 bg-white rounded-full transition-all group-focus-within:w-32" />
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
                hideToolbar={true} // We have our own sidebar
                textareaProps={{
                  placeholder: "In the realm of digital interfaces...",
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
            className="text-slate-400 hover:text-white text-xs font-bold uppercase tracking-widest px-0"
            onClick={() => handleSubmit(0)}
            disabled={isSubmitting || isPublished}
          >
            {isSubmitting ? <Loader2 className="h-4 w-4 animate-spin mr-2" /> : isPublished ? "Saved" : "Save Draft"}
          </Button>
          
          <Button 
            className="bg-[#ad2bee] hover:bg-[#9225c9] text-white border-0 px-6 rounded-lg font-bold text-xs uppercase tracking-widest shadow-[0_8px_20px_rgba(173,43,238,0.3)] transition-all hover:translate-y-[-2px]"
            onClick={() => handleSubmit(1)}
            disabled={isSubmitting}
          >
            <Rocket className="h-3 w-3 mr-2" />
            {isPublished ? "Update" : "Publish"}
          </Button>
        </div>

        <div className="flex-1 overflow-y-auto p-6 space-y-10 custom-scrollbar">
          
          {/* Post Settings */}
          <div className="space-y-6">
            <h3 className="text-[10px] font-bold text-slate-500 uppercase tracking-[0.2em] flex items-center gap-2">
              <Settings className="w-3 h-3" />
              Post Settings
            </h3>
            
            <div className="space-y-3">
              <Label className="text-[11px] font-semibold text-slate-400">URL Slug</Label>
              <div className="bg-white/5 rounded-lg border border-white/5 px-4 py-3 flex items-center gap-3 focus-within:border-[#ad2bee]/50 transition-colors">
                <LinkIcon className="w-3.5 h-3.5 text-slate-500" />
                <input 
                  className="bg-transparent border-none focus:outline-none w-full text-xs font-medium text-slate-300 placeholder:text-slate-600"
                  value={formData.slug}
                  onChange={(e) => handleChange('slug', e.target.value)}
                />
              </div>
            </div>

            <div className="space-y-3">
              <Label className="text-[11px] font-semibold text-slate-400">Tags</Label>
              <div className="bg-white/5 rounded-lg border border-white/5 p-2 flex flex-wrap gap-2 focus-within:border-[#ad2bee]/50 transition-colors">
                {formData.tags.map(tag => (
                  <Badge 
                    key={tag} 
                    className="bg-[#ad2bee]/10 text-[#ad2bee] hover:bg-[#ad2bee]/20 border-none px-3 py-1 text-[10px] font-bold flex items-center gap-2"
                  >
                    #{tag}
                    <X className="w-3 h-3 cursor-pointer" onClick={() => removeTag(tag)}/>
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

          {/* SEO Metadata */}
          <div className="space-y-6">
            <h3 className="text-[10px] font-bold text-slate-500 uppercase tracking-[0.2em] flex items-center gap-2">
              <Search className="w-3 h-3" />
              SEO Metadata
            </h3>
            
            <div className="space-y-3">
              <Label className="text-[11px] font-semibold text-slate-400">Meta Description</Label>
              <textarea 
                className="w-full h-32 bg-white/5 border border-white/5 rounded-lg p-4 text-[11px] leading-relaxed text-slate-400 focus:outline-none focus:border-[#ad2bee]/50 resize-none transition-colors"
                placeholder="Write a short description..."
                value={formData.excerpt}
                onChange={(e) => handleChange('excerpt', e.target.value)}
              />
            </div>
          </div>

          {/* Cover Image */}
          <div className="space-y-6">
            <h3 className="text-[10px] font-bold text-slate-500 uppercase tracking-[0.2em] flex items-center gap-2">
              <ImageIcon className="w-3 h-3" />
              Cover Image
            </h3>
            
            <div className="group relative rounded-xl overflow-hidden bg-white/5 border border-white/5 transition-all hover:border-[#ad2bee]/30 cursor-pointer aspect-video flex items-center justify-center">
              {formData.cover ? (
                <img src={formData.cover} alt="Cover" className="w-full h-full object-cover" />
              ) : (
                <div className="text-center p-6">
                  <div className="bg-white/5 p-3 rounded-full inline-block mb-3 text-slate-500 group-hover:text-[#ad2bee] transition-colors">
                    <ImageIcon className="w-6 h-6" />
                  </div>
                  <p className="text-[11px] text-slate-500">Click to upload cover</p>
                </div>
              )}
              <input 
                type="file" 
                className="absolute inset-0 opacity-0 cursor-pointer" 
                onChange={() => {/* Handle upload */}} 
              />
            </div>
          </div>
        </div>

        {/* Footer */}
        <div className="p-6 border-t border-white/5">
          <p className="text-[10px] text-slate-600 font-medium text-center">
            Kaldalis CMS v2.4.0 • <span className="hover:text-slate-400 cursor-pointer transition-colors">Help</span>
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