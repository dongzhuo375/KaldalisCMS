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
  List as ListIcon
} from "lucide-react";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Badge } from "@/components/ui/badge";
import { cn } from "@/lib/utils";

// Mock Data
const MOCK_FILES = [
  { id: 1, name: "hero-background.jpg", type: "image/jpeg", size: "2.4 MB", url: "https://images.unsplash.com/photo-1579546929518-9e396f3cc809", date: "2023-10-24" },
  { id: 2, name: "avatar-profile.png", type: "image/png", size: "156 KB", url: "https://images.unsplash.com/photo-1535713875002-d1d0cf377fde", date: "2023-10-23" },
  { id: 3, name: "project-mockup.png", type: "image/png", size: "4.1 MB", url: "https://images.unsplash.com/photo-1558655146-d09347e92766", date: "2023-10-22" },
  { id: 4, name: "document-v2.pdf", type: "application/pdf", size: "8.5 MB", url: null, date: "2023-10-21" },
  { id: 5, name: "icon-set.svg", type: "image/svg+xml", size: "32 KB", url: null, date: "2023-10-20" },
  { id: 6, name: "banner-design.jpg", type: "image/jpeg", size: "1.2 MB", url: "https://images.unsplash.com/photo-1557683316-973673baf926", date: "2023-10-19" },
  { id: 7, name: "team-photo.jpg", type: "image/jpeg", size: "3.8 MB", url: "https://images.unsplash.com/photo-1522071820081-009f0129c71c", date: "2023-10-18" },
  { id: 8, name: "config.json", type: "application/json", size: "2 KB", url: null, date: "2023-10-17" },
];

export default function MediaPage() {
  const t = useTranslations('admin');
  const [viewMode, setViewMode] = useState<'grid' | 'list'>('grid');
  const [files, setFiles] = useState(MOCK_FILES);

  const handleDelete = (id: number) => {
    if (confirm("Are you sure you want to delete this file?")) {
      setFiles(files.filter(f => f.id !== id));
    }
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
        
        <Button className="h-10 bg-[#ad2bee] hover:bg-[#9225c9] text-white border-0 shadow-[0_4px_12px_rgba(173,43,238,0.3)] transition-all hover:scale-105 font-medium px-6">
            <Upload className="mr-2 h-4 w-4" /> {t('upload_file')}
        </Button>
      </div>

      {/* Toolbar */}
      <div className="flex flex-col sm:flex-row gap-3 shrink-0">
         <div className="relative flex-1">
            <Search className="absolute left-3 top-2.5 h-4 w-4 text-slate-500 pointer-events-none" />
            <Input 
              placeholder={t('search')}
              className="pl-10 h-10 bg-[#0d0b14]/50 border-slate-800 text-slate-200 focus-visible:ring-[#ad2bee]/30 rounded-lg"
            />
         </div>
         <div className="flex items-center gap-3">
             <Button variant="outline" className="h-10 border-slate-800 bg-[#0d0b14]/50 text-slate-300 hover:bg-slate-900 hover:text-white min-w-[100px] justify-between">
               {t('filter')}
               <Filter className="ml-2 h-3.5 w-3.5 opacity-50" /> 
             </Button>

             <div className="flex items-center bg-[#0d0b14]/50 border border-slate-800 rounded-lg p-1">
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

      {/* Storage Indicator */}
      <div className="w-full bg-[#0d0b14]/30 border border-slate-800/50 rounded-lg p-4 flex items-center justify-between">
         <div className="flex items-center gap-4">
            <div className="h-10 w-10 rounded-full bg-slate-800 flex items-center justify-center text-slate-400">
               <ImageIcon className="h-5 w-5" />
            </div>
            <div>
               <div className="text-sm font-medium text-slate-200">{t('storage_used')}</div>
               <div className="text-xs text-slate-500">2.4 GB of 10 GB used</div>
            </div>
         </div>
         <div className="w-32 md:w-64 h-2 bg-slate-800 rounded-full overflow-hidden">
             <div className="h-full bg-emerald-500 w-[24%]" />
         </div>
      </div>

      {/* Content Area */}
      <div className="bg-[#0d0b14]/40 border border-slate-800/60 rounded-xl overflow-hidden flex-1 shadow-2xl relative p-6">
        
        {viewMode === 'grid' ? (
           <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-6">
             {files.map(file => (
               <div key={file.id} className="group relative aspect-square bg-slate-900/50 rounded-xl border border-slate-800 overflow-hidden hover:border-[#ad2bee]/50 transition-all hover:-translate-y-1">
                 {/* Image Preview or Icon */}
                 {file.url ? (
                   <img src={file.url} alt={file.name} className="w-full h-full object-cover" />
                 ) : (
                   <div className="w-full h-full flex flex-col items-center justify-center gap-2 text-slate-500 bg-slate-900">
                      <File className="h-10 w-10 opacity-20" />
                      <span className="text-xs font-mono uppercase opacity-50">{file.type.split('/')[1]}</span>
                   </div>
                 )}
                 
                 {/* Overlay */}
                 <div className="absolute inset-0 bg-black/60 opacity-0 group-hover:opacity-100 transition-opacity flex flex-col justify-end p-4">
                    <p className="text-xs font-medium text-white truncate w-full mb-1">{file.name}</p>
                    <div className="flex justify-between items-center">
                       <span className="text-[10px] text-slate-400">{file.size}</span>
                       <DropdownMenu>
                         <DropdownMenuTrigger asChild>
                           <button className="h-6 w-6 flex items-center justify-center rounded-full hover:bg-white/20 text-white">
                             <MoreHorizontal className="h-4 w-4" />
                           </button>
                         </DropdownMenuTrigger>
                         <DropdownMenuContent align="end" className="bg-[#1e1b24] border-slate-800 text-slate-200">
                           <DropdownMenuItem onClick={() => handleDelete(file.id)} className="text-rose-400 focus:text-rose-300 focus:bg-rose-950/30 cursor-pointer">
                             <Trash2 className="mr-2 h-3.5 w-3.5" /> {t('delete')}
                           </DropdownMenuItem>
                         </DropdownMenuContent>
                       </DropdownMenu>
                    </div>
                 </div>
               </div>
             ))}
             
             {/* Upload Placeholder */}
             <div className="aspect-square bg-slate-900/20 rounded-xl border border-dashed border-slate-700 hover:border-[#ad2bee] hover:bg-[#ad2bee]/5 transition-all cursor-pointer flex flex-col items-center justify-center gap-3 text-slate-500 hover:text-[#ad2bee]">
                <Upload className="h-8 w-8" />
                <span className="text-xs font-medium uppercase tracking-wider">{t('upload_file')}</span>
             </div>
           </div>
        ) : (
           <div className="divide-y divide-slate-800/50">
             {files.map(file => (
               <div key={file.id} className="flex items-center justify-between py-3 px-4 hover:bg-white/[0.02] rounded-lg transition-colors group">
                  <div className="flex items-center gap-4">
                     <div className="h-10 w-10 rounded bg-slate-800 overflow-hidden flex items-center justify-center shrink-0">
                        {file.url ? (
                          <img src={file.url} alt="" className="h-full w-full object-cover" />
                        ) : (
                          <File className="h-5 w-5 text-slate-500" />
                        )}
                     </div>
                     <div>
                        <p className="text-sm font-medium text-slate-200 group-hover:text-[#ad2bee] transition-colors">{file.name}</p>
                        <p className="text-xs text-slate-500">{file.type} • {file.date}</p>
                     </div>
                  </div>
                  <div className="flex items-center gap-6">
                     <span className="text-xs text-slate-400 font-mono">{file.size}</span>
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
             ))}
           </div>
        )}

      </div>
    </div>
  );
}