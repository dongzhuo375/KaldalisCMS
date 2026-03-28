"use client";

import { useState } from "react";
import { useTranslations } from 'next-intl';
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { 
  Search, 
  Upload, 
  Image as ImageIcon, 
  Trash2, 
  File, 
  MoreHorizontal, 
  Filter,
  Grid,
  List as ListIcon,
  Loader2,
  ExternalLink
} from "lucide-react";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { cn } from "@/lib/utils";
import { useMedia, useUploadMedia, useDeleteMedia } from "@/services/media-service";
import { toast } from "sonner";

export default function MediaPage() {
  const t = useTranslations('admin');
  const [viewMode, setViewMode] = useState<'grid' | 'list'>('grid');
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
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  };

  return (
    <div className="h-full flex flex-col gap-6 text-slate-200 font-sans">
      {/* Header */}
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4 shrink-0">
        <div>
           <h1 className="text-2xl font-bold tracking-tight text-white mb-1">{t('media_library')}</h1>
           <p className="text-slate-400 text-sm">
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
            className="h-10 bg-indigo-600 hover:bg-indigo-700 text-white border-0 shadow-lg transition-all hover:scale-105 font-medium px-6"
            onClick={() => document.getElementById('media-upload')?.click()}
            disabled={uploadMutation.isPending}
          >
            {uploadMutation.isPending ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : <Upload className="mr-2 h-4 w-4" />}
            {t('upload_file')}
          </Button>
        </div>
      </div>

      {/* Toolbar */}
      <div className="flex flex-col sm:flex-row gap-3 shrink-0">
         <div className="relative flex-1">
            <Search className="absolute left-3 top-2.5 h-4 w-4 text-slate-500 pointer-events-none" />
            <Input 
              placeholder={t('search')}
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              className="pl-10 h-10 bg-slate-900/50 border-slate-800 text-slate-200 focus-visible:ring-indigo-500/30 rounded-lg"
            />
         </div>
         <div className="flex items-center gap-3">
             <Button variant="outline" className="h-10 border-slate-800 bg-slate-900/50 text-slate-300 hover:bg-slate-800 hover:text-white min-w-[100px] justify-between">
               {t('filter')}
               <Filter className="ml-2 h-3.5 w-3.5 opacity-50" /> 
             </Button>

             <div className="flex items-center bg-slate-900/50 border border-slate-800 rounded-lg p-1">
                <button 
                  onClick={() => setViewMode('grid')}
                  className={cn(
                    "p-1.5 rounded transition-all",
                    viewMode === 'grid' ? "bg-slate-800 text-white shadow-sm" : "text-slate-500 hover:text-slate-300"
                  )}
                >
                  <Grid className="h-4 w-4" />
                </button>
                <button 
                  onClick={() => setViewMode('list')}
                  className={cn(
                    "p-1.5 rounded transition-all",
                    viewMode === 'list' ? "bg-slate-800 text-white shadow-sm" : "text-slate-500 hover:text-slate-300"
                  )}
                >
                  <ListIcon className="h-4 w-4" />
                </button>
             </div>
         </div>
      </div>

      {/* Content Area */}
      <div className="bg-slate-900/40 border border-slate-800/60 rounded-xl overflow-hidden flex-1 shadow-2xl relative p-6">

        {isLoading ? (
          <div className="absolute inset-0 flex items-center justify-center">
            <Loader2 className="h-8 w-8 animate-spin text-indigo-500" />
          </div>
        ) : files.length === 0 ? (
          <div className="absolute inset-0 flex flex-col items-center justify-center text-slate-500">
             <ImageIcon className="h-16 w-16 mb-4 opacity-10" />
             <p className="text-sm">No files found.</p>
          </div>
        ) : viewMode === 'grid' ? (
           <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-6">
             {files.map(file => (
               <div key={file.id} className="group relative aspect-square bg-slate-900/50 rounded-xl border border-slate-800 overflow-hidden hover:border-indigo-500/50 transition-all hover:-translate-y-1">
                 {/* Image Preview or Icon */}
                 {file.mime_type.startsWith('image/') ? (
                   <img src={file.url} alt={file.filename} className="w-full h-full object-cover" />
                 ) : (
                   <div className="w-full h-full flex flex-col items-center justify-center gap-2 text-slate-500 bg-slate-900">
                      <File className="h-10 w-10 opacity-20" />
                      <span className="text-xs font-mono uppercase opacity-50">{file.mime_type.split('/')[1]}</span>
                   </div>
                 )}

                 {/* Overlay */}
                 <div className="absolute inset-0 bg-black/60 opacity-0 group-hover:opacity-100 transition-opacity flex flex-col justify-end p-4">
                    <p className="text-xs font-medium text-white truncate w-full mb-1">{file.filename}</p>
                    <div className="flex justify-between items-center">
                       <span className="text-[10px] text-slate-400">{formatSize(file.size)}</span>
                       <div className="flex gap-2">
                         <a href={file.url} target="_blank" rel="noreferrer" className="h-6 w-6 flex items-center justify-center rounded-full hover:bg-white/20 text-white">
                            <ExternalLink className="h-3.5 w-3.5" />
                         </a>
                         <DropdownMenu>
                           <DropdownMenuTrigger asChild>
                             <button className="h-6 w-6 flex items-center justify-center rounded-full hover:bg-white/20 text-white">
                               <MoreHorizontal className="h-4 w-4" />
                             </button>
                           </DropdownMenuTrigger>
                           <DropdownMenuContent align="end" className="bg-slate-900 border-slate-800 text-slate-200">
                             <DropdownMenuItem 
                               onClick={() => {
                                 navigator.clipboard.writeText(file.url);
                                 toast.success("Link copied to clipboard");
                               }} 
                               className="cursor-pointer"
                             >
                               Copy URL
                             </DropdownMenuItem>
                             <DropdownMenuItem onClick={() => handleDelete(file.id)} className="text-rose-400 focus:text-rose-300 focus:bg-rose-950/30 cursor-pointer">
                               <Trash2 className="mr-2 h-3.5 w-3.5" /> {t('delete')}
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
           <div className="divide-y divide-slate-800/50">
             {files.map(file => (
               <div key={file.id} className="flex items-center justify-between py-3 px-4 hover:bg-white/[0.02] rounded-lg transition-colors group">
                  <div className="flex items-center gap-4">
                     <div className="h-10 w-10 rounded bg-slate-800 overflow-hidden flex items-center justify-center shrink-0">
                        {file.mime_type.startsWith('image/') ? (
                          <img src={file.url} alt="" className="h-full w-full object-cover" />
                        ) : (
                          <File className="h-5 w-5 text-slate-500" />
                        )}
                     </div>
                     <div>
                        <p className="text-sm font-medium text-slate-200 group-hover:text-indigo-400 transition-colors">{file.filename}</p>
                        <p className="text-xs text-slate-500">{file.mime_type} • {new Date(file.created_at).toLocaleDateString()}</p>
                     </div>
                  </div>
                  <div className="flex items-center gap-6">
                     <span className="text-xs text-slate-400 font-mono">{formatSize(file.size)}</span>
                     <div className="flex gap-2">
                        <Button 
                          variant="ghost" 
                          size="icon" 
                          className="h-8 w-8 text-slate-500 hover:text-indigo-400 hover:bg-indigo-400/10"
                          onClick={() => {
                            navigator.clipboard.writeText(file.url);
                            toast.success("Link copied to clipboard");
                          }}
                        >
                           <ExternalLink className="h-4 w-4" />
                        </Button>
                        <Button 
                          variant="ghost" 
                          size="icon" 
                          className="h-8 w-8 text-slate-500 hover:text-rose-400 hover:bg-rose-950/30"
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