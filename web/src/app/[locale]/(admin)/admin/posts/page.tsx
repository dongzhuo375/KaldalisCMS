"use client";

import { useState } from "react";
import { Link } from "@/i18n/routing";
import { useFormatter } from 'next-intl';
import { useParams } from "next/navigation";
import { 
  DropdownMenu, 
  DropdownMenuContent, 
  DropdownMenuItem, 
  DropdownMenuTrigger,
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
  Filter,
  ChevronLeft,
  ChevronRight,
  Loader2,
  Image as ImageIcon,
  SortAsc
} from "lucide-react";
import { usePosts, useDeletePost } from "@/services/post-service";
import { Input } from "@/components/ui/input";
import { Badge } from "@/components/ui/badge";
import { Checkbox } from "@/components/ui/checkbox";
import { cn } from "@/lib/utils";

export default function PostsPage() {
  const format = useFormatter();
  const params = useParams();
  const rawLocale = params?.locale as string;
  const locale = (rawLocale && rawLocale !== 'undefined') ? rawLocale : 'zh-CN';
  const now = new Date();
  
  const [searchTerm, setSearchTerm] = useState("");
  const [statusFilter, setStatusFilter] = useState<"all" | "published" | "draft" | "archived">("all");
  const [selectedPosts, setSelectedPosts] = useState<number[]>([]);
  const [currentPage, setCurrentPage] = useState(1);
  const itemsPerPage = 10;

  const { data: posts = [], isLoading } = usePosts({ admin: true });
  const deletePostMutation = useDeletePost();

  const formatDateSafe = (dateStr: string) => {
    if (!dateStr) return null;
    const d = new Date(dateStr);
    if (isNaN(d.getTime()) || d.getFullYear() <= 1) return null;
    return d;
  };

  const handleDelete = async (id: number) => {
    if (!confirm("Are you sure you want to delete this post?")) return;
    deletePostMutation.mutate(id, {
      onSuccess: () => {
        setSelectedPosts(prev => prev.filter(pid => pid !== id));
      }
    });
  };

  const filteredPosts = posts.filter(post => {
    const matchesSearch = post.title.toLowerCase().includes(searchTerm.toLowerCase()) ||
                          post.slug.toLowerCase().includes(searchTerm.toLowerCase());
    
    if (statusFilter === "all") return matchesSearch;
    const statusMap = { published: 1, draft: 0, archived: 2 };
    return matchesSearch && post.status === (statusMap as any)[statusFilter];
  });

  const totalPages = Math.ceil(filteredPosts.length / itemsPerPage);
  const paginatedPosts = filteredPosts.slice(
    (currentPage - 1) * itemsPerPage,
    currentPage * itemsPerPage
  );

  const toggleSelectAll = () => {
    if (selectedPosts.length === paginatedPosts.length && paginatedPosts.length > 0) {
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
      
      {/* Header Section */}
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4 shrink-0">
        <div>
           <h1 className="text-2xl font-bold tracking-tight text-white mb-1">Content Management</h1>
           <p className="text-slate-400 text-sm">
             Draft, schedule, and publish your articles.
           </p>
        </div>
        
        <Link href="/admin/posts/new">
          <Button className="h-10 bg-indigo-600 hover:bg-indigo-700 text-white border-0 shadow-lg transition-all hover:scale-105 font-medium px-6">
            <Plus className="mr-2 h-4 w-4" /> Create New Post
          </Button>
        </Link>
      </div>

      {/* Toolbar */}
      <div className="flex flex-col sm:flex-row gap-3 shrink-0">
         <div className="relative flex-1">
            <Search className="absolute left-3 top-2.5 h-4 w-4 text-slate-500 pointer-events-none" />
            <Input 
              placeholder="Search posts..."
              className="pl-10 h-10 bg-slate-900/50 border-slate-800 text-slate-200 focus-visible:ring-indigo-500/30 rounded-lg"
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
            />
         </div>
         <div className="flex items-center gap-3">
             <DropdownMenu>
               <DropdownMenuTrigger asChild>
                 <Button variant="outline" className="h-10 border-slate-800 bg-slate-900/50 text-slate-300 hover:bg-slate-800 hover:text-white min-w-[120px] justify-between">
                   <span className="capitalize">{statusFilter}</span>
                   <Filter className="ml-2 h-3.5 w-3.5 opacity-50" /> 
                 </Button>
               </DropdownMenuTrigger>
               <DropdownMenuContent align="end" className="w-48 bg-slate-900 border-slate-800 text-slate-200">
                 {['all', 'published', 'draft', 'archived'].map((s) => (
                   <DropdownMenuCheckboxItem 
                    key={s}
                    checked={statusFilter === s} 
                    onCheckedChange={() => setStatusFilter(s as any)}
                    className="capitalize"
                   >
                     {s}
                   </DropdownMenuCheckboxItem>
                 ))}
               </DropdownMenuContent>
             </DropdownMenu>

             <Button variant="outline" className="h-10 border-slate-800 bg-slate-900/50 text-slate-300 hover:bg-slate-800 hover:text-white min-w-[140px] justify-between">
               Sort: Newest
               <SortAsc className="ml-2 h-3.5 w-3.5 opacity-50" />
             </Button>
         </div>
      </div>

      {/* Table Container */}
      <div className="bg-slate-900/40 border border-slate-800/60 rounded-xl overflow-hidden flex flex-col flex-1 shadow-2xl relative">
        
        {/* Table Header */}
        <div className="grid grid-cols-[40px_minmax(300px,1fr)_120px_180px_60px] gap-4 px-6 py-3 border-b border-slate-800/60 bg-slate-900/20 text-[11px] font-bold text-slate-500 uppercase tracking-wider items-center sticky top-0 z-10 backdrop-blur-md">
            <Checkbox 
                checked={paginatedPosts.length > 0 && selectedPosts.length === paginatedPosts.length}
                onCheckedChange={toggleSelectAll}
                className="border-slate-700 data-[state=checked]:bg-indigo-600 data-[state=checked]:border-indigo-600"
            />
            <div>Article</div>
            <div>Status</div>
            <div>Last Edited</div>
            <div className="text-right">Actions</div>
        </div>

        {/* Table Body */}
        <div className="flex-1 overflow-y-auto custom-scrollbar">
            {isLoading ? (
                <div className="flex flex-col items-center justify-center h-full py-20">
                    <Loader2 className="h-8 w-8 animate-spin text-indigo-500" />
                </div>
            ) : filteredPosts.length === 0 ? (
                <div className="flex flex-col items-center justify-center h-full py-20 text-slate-500 gap-4">
                    <FileText className="h-12 w-12 opacity-10" />
                    <p className="text-sm">No articles found.</p>
                </div>
            ) : (
                <div className="divide-y divide-slate-800/40">
                    {paginatedPosts.map((post) => {
                        const dateObj = formatDateSafe(post.updated_at);
                        return (
                          <div 
                            key={post.id} 
                            className={cn(
                              "grid grid-cols-[40px_minmax(300px,1fr)_120px_150px_180px_60px] gap-4 px-6 py-4 items-center hover:bg-white/[0.02] transition-colors group",
                              selectedPosts.includes(post.id) && "bg-indigo-500/5"
                            )}
                          >
                              <Checkbox 
                                  checked={selectedPosts.includes(post.id)}
                                  onCheckedChange={() => toggleSelectPost(post.id)}
                                  className="border-slate-700 data-[state=checked]:bg-indigo-600 data-[state=checked]:border-indigo-600"
                              />
                              
                              {/* Article Info */}
                              <div className="flex items-center gap-4 min-w-0">
                                  <div className="h-10 w-16 bg-slate-800 rounded overflow-hidden shrink-0 relative border border-slate-800 group-hover:border-indigo-500/30 transition-colors">
                                     {post.cover ? (
                                       <img src={post.cover} alt="" className="h-full w-full object-cover" />
                                     ) : (
                                       <div className="h-full w-full flex items-center justify-center bg-slate-900">
                                         <ImageIcon className="h-4 w-4 text-slate-700" />
                                       </div>
                                     )}
                                  </div>
                                  <div className="min-w-0 flex flex-col gap-0.5">
                                      <Link href={`/admin/posts/${post.id}/edit`}>
                                          <span className="text-sm font-medium text-slate-200 group-hover:text-indigo-400 transition-colors truncate block cursor-pointer">
                                              {post.title}
                                          </span>
                                      </Link>
                                      <span className="text-[10px] text-slate-500 uppercase font-bold tracking-tight">
                                          BY <span className="text-slate-400">{post.author?.username || "ADMIN"}</span>
                                      </span>
                                  </div>
                              </div>

                              {/* Status */}
                              <div>
                                  <Badge variant="outline" className={cn(
                                      "border-0 px-2.5 py-1 text-[10px] font-bold uppercase rounded-full flex w-fit items-center gap-1.5",
                                      post.status === 1 && "bg-emerald-500/10 text-emerald-400 ring-1 ring-emerald-500/20",
                                      post.status === 0 && "bg-yellow-500/10 text-yellow-400 ring-1 ring-yellow-500/20",
                                      post.status === 2 && "bg-slate-500/10 text-slate-400 ring-1 ring-slate-500/20"
                                  )}>
                                      <span className={cn("h-1 w-1 rounded-full", 
                                            post.status === 1 ? "bg-emerald-500" : 
                                            post.status === 0 ? "bg-yellow-500" : "bg-slate-500"
                                        )}></span>
                                      {post.status === 1 ? "Published" : post.status === 0 ? "Draft" : "Archived"}
                                  </Badge>
                              </div>

                              {/* Last Edited */}
                              <div className="flex flex-col">
                                  <span className="text-xs text-slate-300 font-medium font-mono">
                                      {dateObj ? 
                                        new Intl.DateTimeFormat(locale, { month: 'short', day: '2-digit', year: 'numeric' }).format(dateObj) 
                                        : '-'}
                                  </span>
                                  <span className="text-[10px] text-slate-600 font-bold uppercase tracking-tighter">
                                      {(dateObj && locale !== 'undefined') ? format.relativeTime(dateObj, now) : ''}
                                  </span>
                              </div>

                              {/* Actions */}
                              <div className="text-right pr-2">
                                  <DropdownMenu>
                                      <DropdownMenuTrigger asChild>
                                          <button className="text-slate-500 hover:text-white p-1 rounded-md hover:bg-slate-800 transition-colors">
                                              <MoreHorizontal className="h-4 w-4" />
                                          </button>
                                      </DropdownMenuTrigger>
                                      <DropdownMenuContent align="end" className="w-[160px] bg-slate-900 border-slate-800 text-slate-200">
                                          <Link href={`/admin/posts/${post.id}/edit`}>
                                              <DropdownMenuItem className="cursor-pointer">
                                                  <Pencil className="mr-2 h-3.5 w-3.5" /> Edit
                                              </DropdownMenuItem>
                                          </Link>
                                          <DropdownMenuItem 
                                            onClick={() => handleDelete(post.id)} 
                                            disabled={deletePostMutation.isPending}
                                            className="text-rose-400 focus:text-rose-300 focus:bg-rose-950/30 cursor-pointer"
                                          >
                                              <Trash className="mr-2 h-3.5 w-3.5" /> Delete
                                          </DropdownMenuItem>
                                      </DropdownMenuContent>
                                  </DropdownMenu>
                              </div>
                          </div>
                        );
                    })}
                </div>
            )}
        </div>

        {/* Footer Pagination */}
        {totalPages > 1 && (
          <div className="border-t border-slate-800/60 p-4 bg-slate-950/60 flex items-center justify-between backdrop-blur-md shrink-0 z-20">
              <div className="text-[10px] text-slate-500 font-bold uppercase tracking-widest">
                  Showing <span className="text-white">{(currentPage - 1) * itemsPerPage + 1}-{Math.min(currentPage * itemsPerPage, filteredPosts.length)}</span> of <span className="text-white">{filteredPosts.length}</span>
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
        )}

      </div>
    </div>
  );
}
