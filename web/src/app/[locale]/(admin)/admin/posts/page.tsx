"use client";

import { useState, useEffect } from "react";
import { Link } from "@/i18n/routing";
import { useTranslations, useFormatter } from 'next-intl';
import { 
  DropdownMenu, 
  DropdownMenuContent, 
  DropdownMenuItem, 
  DropdownMenuLabel, 
  DropdownMenuTrigger,
  DropdownMenuSeparator,
  DropdownMenuCheckboxItem
} from "@/components/ui/dropdown-menu";
import { Button } from "@/components/ui/button";
import { 
  MoreHorizontal, 
  Plus, 
  Pencil, 
  Trash, 
  FileText, 
  Search, 
  Terminal,
  Filter,
  CheckCircle2,
  FileClock,
  Archive,
  ChevronLeft,
  ChevronRight,
  SlidersHorizontal,
  ArrowUpDown,
  Download,
  Loader2,
  Calendar as CalendarIcon,
  X,
  LayoutGrid,
  List,
  SortAsc,
  Image as ImageIcon
} from "lucide-react";
import { Post } from "@/lib/types";
import api from "@/lib/api";
import { Input } from "@/components/ui/input";
import { Badge } from "@/components/ui/badge";
import { Checkbox } from "@/components/ui/checkbox";
import { cn } from "@/lib/utils";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";

export default function PostsPage() {
  const t = useTranslations('admin');
  const format = useFormatter();
  const now = new Date();
  const [posts, setPosts] = useState<Post[]>([]);
  const [loading, setLoading] = useState(true);
  const [searchTerm, setSearchTerm] = useState("");
  const [statusFilter, setStatusFilter] = useState<"all" | "published" | "draft" | "archived">("all");
  const [selectedPosts, setSelectedPosts] = useState<number[]>([]);

  // Pagination (mock)
  const [currentPage, setCurrentPage] = useState(1);
  const itemsPerPage = 8;

  const fetchPosts = async () => {
    setLoading(true);
    try {
      const data = await api.get<Post[]>("/posts");
      setPosts(data as unknown as Post[]);
    } catch (error) {
      console.error("Failed to fetch posts", error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchPosts();
  }, []);

  const handleDelete = async (id: number) => {
    if (!confirm("Are you sure you want to delete this post?")) return;
    try {
      await api.delete(`/posts/${id}`);
      setPosts(posts.filter(p => p.id !== id));
      setSelectedPosts(selectedPosts.filter(pid => pid !== id));
    } catch (error) {
      console.error("Delete failed", error);
    }
  };

  const filteredPosts = posts.filter(post => {
    const matchesSearch = post.title.toLowerCase().includes(searchTerm.toLowerCase()) ||
                          post.slug.toLowerCase().includes(searchTerm.toLowerCase());
    
    if (statusFilter === "all") return matchesSearch;
    if (statusFilter === "published") return matchesSearch && post.status === 1;
    if (statusFilter === "draft") return matchesSearch && post.status === 0;
    if (statusFilter === "archived") return matchesSearch && post.status === 2;
    
    return matchesSearch;
  });

  // Pagination logic
  const totalPages = Math.ceil(filteredPosts.length / itemsPerPage);
  const paginatedPosts = filteredPosts.slice(
    (currentPage - 1) * itemsPerPage,
    currentPage * itemsPerPage
  );

  const toggleSelectAll = () => {
    if (selectedPosts.length === paginatedPosts.length) {
      setSelectedPosts([]);
    } else {
      setSelectedPosts(paginatedPosts.map(p => p.id));
    }
  };

  const toggleSelectPost = (id: number) => {
    if (selectedPosts.includes(id)) {
      setSelectedPosts(selectedPosts.filter(pid => pid !== id));
    } else {
      setSelectedPosts([...selectedPosts, id]);
    }
  };

  return (
    <div className="h-full flex flex-col gap-6 text-slate-200 font-sans">
      
      {/* 1. Header Section */}
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4 shrink-0">
        <div>
           <h1 className="text-2xl font-bold tracking-tight text-white mb-1">Article Management</h1>
           <p className="text-slate-400 text-sm">
             Manage, edit, and publish your content.
           </p>
        </div>
        
        <Link href="/admin/posts/new">
          <Button className="h-10 bg-[#ad2bee] hover:bg-[#9225c9] text-white border-0 shadow-[0_4px_12px_rgba(173,43,238,0.3)] transition-all hover:scale-105 font-medium px-6">
            <Plus className="mr-2 h-4 w-4" /> Create New Post
          </Button>
        </Link>
      </div>

      {/* 2. Toolbar */}
      <div className="flex flex-col sm:flex-row gap-3 shrink-0">
         <div className="relative flex-1">
            <Search className="absolute left-3 top-2.5 h-4 w-4 text-slate-500 pointer-events-none" />
            <Input 
              placeholder="Search by title, author, or tag..."
              className="pl-10 h-10 bg-[#0d0b14]/50 border-slate-800 text-slate-200 focus-visible:ring-[#ad2bee]/30 rounded-lg"
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
            />
         </div>
         <div className="flex items-center gap-3">
             <DropdownMenu>
               <DropdownMenuTrigger asChild>
                 <Button variant="outline" className="h-10 border-slate-800 bg-[#0d0b14]/50 text-slate-300 hover:bg-slate-900 hover:text-white min-w-[120px] justify-between">
                   {statusFilter === 'all' ? "All Status" : t(statusFilter as any)}
                   <Filter className="ml-2 h-3.5 w-3.5 opacity-50" /> 
                 </Button>
               </DropdownMenuTrigger>
               <DropdownMenuContent align="end" className="w-48 bg-[#1e1b24] border-slate-800 text-slate-200">
                 <DropdownMenuCheckboxItem checked={statusFilter === 'all'} onCheckedChange={() => setStatusFilter('all')}>All Status</DropdownMenuCheckboxItem>
                 <DropdownMenuCheckboxItem checked={statusFilter === 'published'} onCheckedChange={() => setStatusFilter('published')}>Published</DropdownMenuCheckboxItem>
                 <DropdownMenuCheckboxItem checked={statusFilter === 'draft'} onCheckedChange={() => setStatusFilter('draft')}>Draft</DropdownMenuCheckboxItem>
                 <DropdownMenuCheckboxItem checked={statusFilter === 'archived'} onCheckedChange={() => setStatusFilter('archived')}>Archived</DropdownMenuCheckboxItem>
               </DropdownMenuContent>
             </DropdownMenu>

             <Button variant="outline" className="h-10 border-slate-800 bg-[#0d0b14]/50 text-slate-300 hover:bg-slate-900 hover:text-white min-w-[120px] justify-between">
               Categories
               <Filter className="ml-2 h-3.5 w-3.5 opacity-50" />
             </Button>

             <Button variant="outline" className="h-10 border-slate-800 bg-[#0d0b14]/50 text-slate-300 hover:bg-slate-900 hover:text-white min-w-[140px] justify-between">
               Sort: Newest
               <SortAsc className="ml-2 h-3.5 w-3.5 opacity-50" />
             </Button>
         </div>
      </div>

      {/* 3. Table Container */}
      <div className="bg-[#0d0b14]/40 border border-slate-800/60 rounded-xl overflow-hidden flex flex-col flex-1 shadow-2xl relative">
        
        {/* Table Header */}
        <div className="grid grid-cols-[40px_minmax(300px,1fr)_120px_150px_180px_60px] gap-4 px-6 py-3 border-b border-slate-800/60 bg-[#0d0b14]/20 text-[11px] font-bold text-slate-500 uppercase tracking-wider items-center sticky top-0 z-10 backdrop-blur-md">
            <Checkbox 
                checked={paginatedPosts.length > 0 && selectedPosts.length === paginatedPosts.length}
                onCheckedChange={toggleSelectAll}
                className="border-slate-600 data-[state=checked]:bg-[#ad2bee] data-[state=checked]:border-[#ad2bee]"
            />
            <div>Article</div>
            <div>Status</div>
            <div>Tags</div>
            <div>Last Edited</div>
            <div className="text-right">Actions</div>
        </div>

        {/* Table Body - Scrollable Area */}
        <div className="flex-1 overflow-y-auto custom-scrollbar">
            {loading ? (
                <div className="flex flex-col items-center justify-center h-full py-20">
                    <Loader2 className="h-8 w-8 animate-spin text-[#ad2bee]" />
                </div>
            ) : filteredPosts.length === 0 ? (
                <div className="flex flex-col items-center justify-center h-full py-20 text-slate-500 gap-4">
                    <FileText className="h-12 w-12 opacity-20" />
                    <p>No articles found.</p>
                </div>
            ) : (
                <div className="divide-y divide-slate-800/40">
                    {paginatedPosts.map((post) => (
                        <div 
                          key={post.id} 
                          className={cn(
                            "grid grid-cols-[40px_minmax(300px,1fr)_120px_150px_180px_60px] gap-4 px-6 py-4 items-center hover:bg-white/[0.02] transition-colors group",
                            selectedPosts.includes(post.id) && "bg-[#ad2bee]/5 hover:bg-[#ad2bee]/10"
                          )}
                        >
                            <Checkbox 
                                checked={selectedPosts.includes(post.id)}
                                onCheckedChange={() => toggleSelectPost(post.id)}
                                className="border-slate-600 data-[state=checked]:bg-[#ad2bee] data-[state=checked]:border-[#ad2bee]"
                            />
                            
                            {/* Article Info */}
                            <div className="flex items-center gap-4 min-w-0">
                                <div className="h-10 w-16 bg-slate-800 rounded overflow-hidden shrink-0 relative border border-slate-700/50 group-hover:border-[#ad2bee]/30 transition-colors">
                                   {post.cover ? (
                                     <img src={post.cover} alt="" className="h-full w-full object-cover" />
                                   ) : (
                                     <div className="h-full w-full flex items-center justify-center bg-gradient-to-br from-slate-800 to-slate-900">
                                       <ImageIcon className="h-4 w-4 text-slate-600" />
                                     </div>
                                   )}
                                </div>
                                <div className="min-w-0 flex flex-col gap-0.5">
                                    <Link href={`/admin/posts/${post.id}/edit`} className="block">
                                        <span className="text-sm font-medium text-slate-200 group-hover:text-[#ad2bee] transition-colors truncate block cursor-pointer">
                                            {post.title}
                                        </span>
                                    </Link>
                                    <span className="text-xs text-slate-500 truncate">
                                        by <span className="text-slate-400">{post.author?.username || "Admin"}</span>
                                    </span>
                                </div>
                            </div>

                            {/* Status */}
                            <div>
                                <Badge variant="outline" className={cn(
                                    "border-0 px-2.5 py-1 text-[11px] font-medium rounded-full flex w-fit items-center gap-1.5 transition-transform",
                                    post.status === 1 && "bg-emerald-500/10 text-emerald-400 ring-1 ring-emerald-500/20",
                                    post.status === 0 && "bg-yellow-500/10 text-yellow-400 ring-1 ring-yellow-500/20",
                                    post.status === 2 && "bg-slate-500/10 text-slate-400 ring-1 ring-slate-500/20"
                                )}>
                                    <span className={cn("relative flex h-1.5 w-1.5 rounded-full", 
                                          post.status === 1 ? "bg-emerald-500" : 
                                          post.status === 0 ? "bg-yellow-500" : "bg-slate-500"
                                      )}></span>
                                    {post.status === 1 ? "Published" : post.status === 0 ? "Draft" : "Archived"}
                                </Badge>
                            </div>

                            {/* Tags */}
                            <div className="flex flex-wrap gap-1">
                                <Badge variant="secondary" className="bg-slate-800 text-slate-400 border border-slate-700 text-[10px] h-5 hover:bg-slate-700">#Design</Badge>
                                <Badge variant="secondary" className="bg-slate-800 text-slate-400 border border-slate-700 text-[10px] h-5 hover:bg-slate-700">#UX</Badge>
                            </div>

                            {/* Last Edited */}
                            <div className="flex flex-col">
                                <span className="text-xs text-slate-300 font-medium">
                                    {post.created_at ? format.dateTime(new Date(post.created_at), { month: 'short', day: '2-digit', year: 'numeric' }) : '-'}
                                </span>
                                <span className="text-[10px] text-slate-600">
                                    {post.created_at ? format.relativeTime(new Date(post.created_at), now) : ''}
                                </span>
                            </div>

                            {/* Actions */}
                            <div className="text-right pr-2">
                                <DropdownMenu>
                                    <DropdownMenuTrigger asChild>
                                        <Button variant="ghost" size="icon" className="h-8 w-8 text-slate-500 hover:text-white hover:bg-white/5 rounded-md">
                                            <MoreHorizontal className="h-4 w-4" />
                                        </Button>
                                    </DropdownMenuTrigger>
                                    <DropdownMenuContent align="end" className="w-[160px] bg-[#1e1b24] border-slate-800 text-slate-200">
                                        <Link href={`/admin/posts/${post.id}/edit`}>
                                            <DropdownMenuItem className="cursor-pointer focus:bg-slate-800 focus:text-white">
                                                <Pencil className="mr-2 h-3.5 w-3.5" /> Edit
                                            </DropdownMenuItem>
                                        </Link>
                                        <DropdownMenuItem onClick={() => handleDelete(post.id)} className="text-rose-400 focus:text-rose-300 focus:bg-rose-950/30 cursor-pointer">
                                            <Trash className="mr-2 h-3.5 w-3.5" /> Delete
                                        </DropdownMenuItem>
                                    </DropdownMenuContent>
                                </DropdownMenu>
                            </div>
                        </div>
                    ))}
                </div>
            )}
        </div>

        {/* Footer Pagination */}
        <div className="border-t border-slate-800/60 p-4 bg-[#0d0b14]/60 flex items-center justify-between backdrop-blur-md shrink-0 z-20">
            <div className="text-xs text-slate-500">
                Showing <span className="text-white font-medium">{(currentPage - 1) * itemsPerPage + 1}-{Math.min(currentPage * itemsPerPage, filteredPosts.length)}</span> of <span className="text-white font-medium">{filteredPosts.length}</span> articles
            </div>
            <div className="flex items-center gap-1">
                <Button 
                    variant="outline" 
                    size="icon" 
                    className="h-8 w-8 border-slate-800 bg-transparent hover:bg-white/5 text-slate-400 hover:text-white"
                    onClick={() => setCurrentPage(p => Math.max(1, p - 1))}
                    disabled={currentPage === 1}
                >
                    <ChevronLeft className="h-4 w-4" />
                </Button>
                {Array.from({ length: totalPages }).map((_, i) => (
                    <button
                        key={i}
                        onClick={() => setCurrentPage(i + 1)}
                        className={cn(
                            "h-8 w-8 text-xs rounded-md font-medium transition-all duration-200 flex items-center justify-center",
                            currentPage === i + 1 
                            ? "bg-[#ad2bee] text-white shadow-lg shadow-[#ad2bee]/30" 
                            : "text-slate-500 hover:bg-white/5 hover:text-slate-300"
                        )}
                    >
                        {i + 1}
                    </button>
                ))}
                <Button 
                    variant="outline" 
                    size="icon" 
                    className="h-8 w-8 border-slate-800 bg-transparent hover:bg-white/5 text-slate-400 hover:text-white"
                    onClick={() => setCurrentPage(p => Math.min(totalPages, p + 1))}
                    disabled={currentPage === totalPages}
                >
                    <ChevronRight className="h-4 w-4" />
                </Button>
            </div>
        </div>

      </div>

      <style jsx global>{`
        .custom-scrollbar::-webkit-scrollbar {
          width: 6px;
        }
        .custom-scrollbar::-webkit-scrollbar-track {
          background: rgba(0,0,0,0.1);
        }
        .custom-scrollbar::-webkit-scrollbar-thumb {
          background: rgba(255, 255, 255, 0.1);
          border-radius: 10px;
        }
        .custom-scrollbar::-webkit-scrollbar-thumb:hover {
          background: rgba(255, 255, 255, 0.2);
        }
      `}</style>
    </div>
  );
}