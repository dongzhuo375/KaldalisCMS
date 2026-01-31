"use client";

import { useState, useEffect, useRef } from "react";
import dynamic from "next/dynamic";
import { useRouter } from "@/i18n/routing";
import { useTranslations } from 'next-intl';
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Badge } from "@/components/ui/badge";
import { Separator } from "@/components/ui/select"; // Using SelectSeparator as generic separator if needed, or div
import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from "@/components/ui/sheet" // Shadcn Sheet for settings if we want it collapsible, or just a div side panel.
import { 
  ArrowLeft, 
  Loader2, 
  Globe, 
  ImageIcon, 
  Settings, 
  Calendar, 
  Hash, 
  Type, 
  Eye, 
  Send, 
  Save,
  PanelRightOpen,
  X
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
  const [showSettings, setShowSettings] = useState(true); // Default show right sidebar
  
  const [formData, setFormData] = useState({
    title: initialData?.title || "",
    slug: initialData?.slug || "",
    content: initialData?.content || "",
    status: initialData?.status?.toString() || "0", 
    cover: initialData?.cover || "",
    tags: [] as string[], // Mock tags array
    excerpt: "", // Missing from Post type but useful for SEO
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

  const handleSubmit = (statusOverride?: number) => {
    onSubmit({
      ...formData,
      status: statusOverride ?? parseInt(formData.status),
      category_id: 1, 
      tags: [1] // Backend expects IDs for now
    });
  };

  // Tag handling (Mock)
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
    <div className="flex h-screen flex-col bg-slate-950 overflow-hidden text-slate-200 animate-in fade-in duration-500">
      
      {/* 1. Top Navigation Bar */}
      <header className="flex h-16 items-center justify-between border-b border-slate-800 bg-slate-950 px-6 shrink-0">
        <div className="flex items-center gap-4">
          <Button 
            variant="ghost" 
            size="icon"
            onClick={() => router.back()}
            className="text-slate-400 hover:text-white hover:bg-slate-800 rounded-full"
          >
            <ArrowLeft className="h-5 w-5" />
          </Button>
          <div className="flex flex-col">
             <span className="text-sm font-medium text-slate-400">
                {mode === 'create' ? "New Post" : "Editing Post"}
             </span>
             <span className="text-xs text-slate-600 flex items-center gap-1">
                {isPublished ? <span className="w-1.5 h-1.5 rounded-full bg-emerald-500"/> : <span className="w-1.5 h-1.5 rounded-full bg-yellow-500"/>}
                {isPublished ? "Published" : "Draft"}
             </span>
          </div>
        </div>

        <div className="flex items-center gap-3">
          <span className="text-xs text-slate-500 font-mono hidden sm:inline-block mr-4">
             {wordCount} words
          </span>
          
          <Button 
            variant="ghost" 
            className="text-slate-400 hover:text-white"
            onClick={() => handleSubmit(0)} // Save as Draft
            disabled={isSubmitting}
          >
            {isSubmitting ? <Loader2 className="h-4 w-4 animate-spin" /> : <Save className="h-4 w-4 mr-2" />}
            Save Draft
          </Button>
          
          <Button 
            className="bg-[#ad2bee] hover:bg-[#9225c9] text-white border-0 shadow-[0_0_15px_rgba(173,43,238,0.3)] transition-all hover:scale-105"
            onClick={() => handleSubmit(1)} // Publish
            disabled={isSubmitting}
          >
             {isPublished ? "Update" : "Publish"}
             <Send className="h-4 w-4 ml-2" />
          </Button>
          
          <Button
            variant="ghost"
            size="icon"
            onClick={() => setShowSettings(!showSettings)}
            className={cn("ml-2 text-slate-400 hover:text-white", showSettings && "bg-slate-800 text-white")}
          >
            <PanelRightOpen className={cn("h-5 w-5 transition-transform", showSettings && "rotate-180")} />
          </Button>
        </div>
      </header>

      {/* 2. Main Workspace (Flex Row) */}
      <div className="flex flex-1 overflow-hidden">
        
        {/* Left Formatting Bar (Optional/Conceptual) */}
        {/* For now, relying on MDEditor's toolbar, but reserving space or integrating here would be next step */}

        {/* Center: Editor Canvas */}
        <main className="flex-1 overflow-y-auto relative bg-[#0d1117] custom-scrollbar">
           <div className="mx-auto max-w-3xl py-12 px-8 min-h-full">
              
              {/* Title Input (as part of the document) */}
              <input
                type="text"
                placeholder="Post Title"
                className="w-full bg-transparent border-none text-4xl sm:text-5xl font-bold text-white placeholder:text-slate-700 focus:outline-none mb-8 resize-none"
                value={formData.title}
                onChange={(e) => handleChange('title', e.target.value)}
                autoFocus
              />

              {/* Markdown Editor */}
              <div data-color-mode="dark" className="prose prose-invert max-w-none">
                <MDEditor
                  value={formData.content}
                  onChange={(val) => handleChange('content', val || "")}
                  height={600}
                  className="bg-transparent border-none shadow-none !bg-transparent"
                  visibleDragbar={false}
                  preview="edit"
                  textareaProps={{
                    placeholder: "Tell your story...",
                    className: "text-lg leading-relaxed text-slate-300 font-serif"
                  }}
                  toolbarHeight={40}
                  style={{ backgroundColor: 'transparent', minHeight: '400px' }}
                />
              </div>
           </div>
        </main>

        {/* Right: Settings Sidebar */}
        <aside 
          className={cn(
            "w-80 bg-slate-950 border-l border-slate-800 overflow-y-auto transition-all duration-300 ease-in-out shrink-0",
            !showSettings && "w-0 border-l-0 opacity-0 overflow-hidden"
          )}
        >
          <div className="p-6 space-y-8">
            
            {/* Post Settings Group */}
            <div className="space-y-4">
               <h3 className="text-xs font-semibold text-slate-500 uppercase tracking-wider flex items-center gap-2">
                 <Settings className="w-3 h-3" />
                 Post Settings
               </h3>
               
               <div className="space-y-3">
                 <Label className="text-slate-400 text-xs">URL Slug</Label>
                 <div className="bg-slate-900 rounded-md border border-slate-800 px-3 py-2 text-sm text-slate-300 flex items-center gap-2">
                    <Globe className="w-3 h-3 text-slate-500 shrink-0" />
                    <input 
                      className="bg-transparent border-none focus:outline-none w-full text-slate-300 text-xs font-mono truncate"
                      value={formData.slug}
                      onChange={(e) => handleChange('slug', e.target.value)}
                    />
                 </div>
               </div>

               <div className="space-y-3">
                 <Label className="text-slate-400 text-xs">Publish Date</Label>
                 <div className="bg-slate-900 rounded-md border border-slate-800 px-3 py-2 text-sm text-slate-300 flex items-center gap-2 cursor-not-allowed opacity-70">
                    <Calendar className="w-3 h-3 text-slate-500 shrink-0" />
                    <span className="text-xs">Immediately</span>
                 </div>
               </div>

               <div className="space-y-3">
                 <Label className="text-slate-400 text-xs">Tags</Label>
                 <div className="bg-slate-900 rounded-md border border-slate-800 p-2 min-h-[80px]">
                    <div className="flex flex-wrap gap-2 mb-2">
                       {formData.tags.map(tag => (
                         <Badge key={tag} variant="secondary" className="bg-slate-800 text-slate-300 hover:bg-slate-700 h-6 text-[10px] gap-1 pr-1">
                           {tag}
                           <X className="w-3 h-3 cursor-pointer hover:text-white" onClick={() => removeTag(tag)}/>
                         </Badge>
                       ))}
                    </div>
                    <input 
                      className="bg-transparent border-none focus:outline-none w-full text-slate-300 text-xs"
                      placeholder="Add tag and press Enter..."
                      value={tagInput}
                      onChange={(e) => setTagInput(e.target.value)}
                      onKeyDown={handleAddTag}
                    />
                 </div>
               </div>
            </div>

            <div className="h-px bg-slate-800 my-4" />

            {/* SEO & Meta Group */}
            <div className="space-y-4">
               <h3 className="text-xs font-semibold text-slate-500 uppercase tracking-wider flex items-center gap-2">
                 <Type className="w-3 h-3" />
                 Meta Data
               </h3>
               
               <div className="space-y-3">
                 <Label className="text-slate-400 text-xs">Excerpt</Label>
                 <textarea 
                    className="w-full h-24 bg-slate-900 border border-slate-800 rounded-md p-3 text-xs text-slate-300 focus:outline-none focus:border-indigo-500 resize-none"
                    placeholder="Write a short excerpt for SEO..."
                    value={formData.excerpt}
                    onChange={(e) => handleChange('excerpt', e.target.value)}
                 />
               </div>
            </div>

            <div className="h-px bg-slate-800 my-4" />

            {/* Feature Image Group */}
            <div className="space-y-4">
               <h3 className="text-xs font-semibold text-slate-500 uppercase tracking-wider flex items-center gap-2">
                 <ImageIcon className="w-3 h-3" />
                 Feature Image
               </h3>
               
               <div className="border-2 border-dashed border-slate-800 rounded-lg p-4 text-center hover:bg-slate-900/50 transition-colors cursor-pointer group relative overflow-hidden">
                  {formData.cover ? (
                    <>
                      {/* eslint-disable-next-line @next/next/no-img-element */}
                      <img src={formData.cover} alt="Cover" className="w-full h-32 object-cover rounded-md mb-2" />
                      <button 
                        onClick={(e) => { e.stopPropagation(); handleChange('cover', ''); }}
                        className="absolute top-2 right-2 bg-black/50 p-1 rounded-full text-white opacity-0 group-hover:opacity-100 transition-opacity"
                      >
                        <X className="w-3 h-3" />
                      </button>
                    </>
                  ) : (
                    <div className="flex flex-col items-center gap-2 py-4">
                       <div className="bg-slate-800 p-2 rounded-full text-slate-400 group-hover:text-white transition-colors">
                         <ImageIcon className="w-5 h-5" />
                       </div>
                       <span className="text-xs text-slate-500">Click to add image url</span>
                    </div>
                  )}
                  {/* Temporary Input overlay for URL */}
                  {!formData.cover && (
                    <input 
                      className="absolute inset-0 opacity-0 cursor-pointer"
                      onChange={(e) => {
                         const url = prompt("Enter Image URL:");
                         if(url) handleChange('cover', url);
                      }} 
                    />
                  )}
               </div>
            </div>

          </div>
        </aside>

      </div>
    </div>
  );
}