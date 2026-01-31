"use client";

import { useState, useEffect } from "react";
import dynamic from "next/dynamic";
import { useRouter } from "@/i18n/routing";
import { useTranslations } from 'next-intl';
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent } from "@/components/ui/card";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Save, ArrowLeft, Loader2, Globe, FileText, Image as ImageIcon } from "lucide-react";
import { Post } from "@/lib/types";

// Dynamically import Markdown Editor to avoid SSR issues
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
    status: initialData?.status?.toString() || "0", // Default to Draft (0)
    cover: initialData?.cover || "",
  });

  // Auto-generate slug from title if slug is empty (only in create mode)
  useEffect(() => {
    if (mode === 'create' && !initialData?.slug && formData.title) {
      const slug = formData.title
        .toLowerCase()
        .replace(/[^a-z0-9]+/g, '-')
        .replace(/(^-|-$)+/g, '');
      setFormData(prev => ({ ...prev, slug }));
    }
  }, [formData.title, mode, initialData?.slug]);

  const handleChange = (key: string, value: string) => {
    setFormData(prev => ({ ...prev, [key]: value }));
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit({
      ...formData,
      status: parseInt(formData.status),
      // Mock category/tags for now as per plan
      category_id: 1, 
      tags: [1]
    });
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-500">
      
      {/* Toolbar */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div className="flex items-center gap-2">
          <Button 
            type="button" 
            variant="ghost" 
            onClick={() => router.back()}
            className="text-slate-400 hover:text-white hover:bg-slate-800"
          >
            <ArrowLeft className="h-4 w-4 mr-2" />
            {t('back')}
          </Button>
          <h1 className="text-2xl font-bold text-white tracking-tight">
            {mode === 'create' ? t('create') : t('edit')} Post
          </h1>
        </div>
        
        <div className="flex items-center gap-3">
          <Select 
            value={formData.status} 
            onValueChange={(val) => handleChange('status', val)}
          >
            <SelectTrigger className="w-[140px] bg-slate-900 border-slate-700 text-slate-200">
              <SelectValue placeholder="Status" />
            </SelectTrigger>
            <SelectContent className="bg-slate-900 border-slate-700 text-slate-200">
              <SelectItem value="0">Draft</SelectItem>
              <SelectItem value="1">Published</SelectItem>
              <SelectItem value="2">Archived</SelectItem>
            </SelectContent>
          </Select>
          
          <Button 
            type="submit" 
            className="bg-emerald-600 hover:bg-emerald-500 text-white min-w-[120px]"
            disabled={isSubmitting}
          >
            {isSubmitting ? (
              <>
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                Saving...
              </>
            ) : (
              <>
                <Save className="mr-2 h-4 w-4" />
                {t('save')}
              </>
            )}
          </Button>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        
        {/* Main Editor Area */}
        <div className="lg:col-span-2 space-y-6">
          <Card className="bg-slate-950/50 border-slate-800 backdrop-blur-sm">
            <CardContent className="p-6 space-y-4">
              <div className="space-y-2">
                <Label htmlFor="title" className="text-slate-400">Title</Label>
                <Input 
                  id="title"
                  placeholder="Enter post title..." 
                  className="bg-slate-900/50 border-slate-700 text-lg font-medium text-white placeholder:text-slate-600 focus-visible:ring-emerald-500/50"
                  value={formData.title}
                  onChange={(e) => handleChange('title', e.target.value)}
                  required
                />
              </div>
              
              <div className="space-y-2">
                <Label htmlFor="slug" className="text-slate-400 flex items-center gap-2">
                  <Globe className="w-3 h-3" /> Slug
                </Label>
                <div className="flex items-center gap-2 text-sm text-slate-500 font-mono bg-slate-900/30 p-2 rounded border border-slate-800">
                  <span>/posts/</span>
                  <input 
                    id="slug"
                    className="bg-transparent border-none focus:outline-none text-emerald-400 w-full placeholder:text-slate-700"
                    placeholder="post-url-slug"
                    value={formData.slug}
                    onChange={(e) => handleChange('slug', e.target.value)}
                    required
                  />
                </div>
              </div>

              <div className="space-y-2 pt-2" data-color-mode="dark">
                <Label className="text-slate-400 flex items-center gap-2">
                  <FileText className="w-3 h-3" /> Content (Markdown)
                </Label>
                <div className="border border-slate-700 rounded-lg overflow-hidden">
                  <MDEditor
                    value={formData.content}
                    onChange={(val) => handleChange('content', val || "")}
                    height={500}
                    className="bg-slate-900"
                    style={{ backgroundColor: '#0f172a', color: '#e2e8f0' }}
                    preview="edit"
                  />
                </div>
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Sidebar Settings */}
        <div className="space-y-6">
          <Card className="bg-slate-950/50 border-slate-800 backdrop-blur-sm">
             <CardContent className="p-4 space-y-4">
                <h3 className="font-semibold text-slate-200 mb-4">Post Settings</h3>
                
                <div className="space-y-2">
                  <Label htmlFor="cover" className="text-slate-400 text-xs uppercase tracking-wider">Cover Image URL</Label>
                  <div className="relative">
                    <ImageIcon className="absolute left-3 top-2.5 h-4 w-4 text-slate-500" />
                    <Input 
                      id="cover"
                      placeholder="https://..." 
                      className="pl-9 bg-slate-900 border-slate-700 text-sm text-slate-200"
                      value={formData.cover}
                      onChange={(e) => handleChange('cover', e.target.value)}
                    />
                  </div>
                  {formData.cover && (
                    <div className="mt-2 rounded-lg overflow-hidden border border-slate-800 aspect-video relative">
                      {/* eslint-disable-next-line @next/next/no-img-element */}
                      <img 
                        src={formData.cover} 
                        alt="Cover preview" 
                        className="w-full h-full object-cover"
                        onError={(e) => {
                          (e.target as HTMLImageElement).style.display = 'none';
                        }}
                      />
                    </div>
                  )}
                </div>

                <div className="p-3 bg-yellow-500/10 border border-yellow-500/20 rounded-lg text-xs text-yellow-500/80">
                   Note: Categories and Tags selection will be available once the API is ready.
                </div>
             </CardContent>
          </Card>
        </div>

      </div>
    </form>
  );
}
