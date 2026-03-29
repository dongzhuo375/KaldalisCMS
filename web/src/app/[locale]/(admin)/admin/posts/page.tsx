"use client";

import { useState } from "react";
import { Link } from "@/i18n/routing";
import { useFormatter, useLocale } from 'next-intl';
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
  SortAsc,
  Calendar
} from "lucide-react";
import { usePosts, useDeletePost } from "@/services/post-service";
import { Input } from "@/components/ui/input";
import { Badge } from "@/components/ui/badge";
import { cn } from "@/lib/utils";
import { motion, AnimatePresence } from "framer-motion";

export default function PostsPage() {
  const format = useFormatter();
  const params = useParams();
  const rawLocale = params?.locale as string;
  const locale = (rawLocale && rawLocale !== 'undefined') ? rawLocale : 'zh-CN';
  const now = new Date();
  
  const [searchTerm, setSearchTerm] = useState("");
  const [statusFilter, setStatusFilter] = useState<"all" | "published" | "draft" | "archived">("all");
  const [currentPage, setCurrentPage] = useState(1);
  const itemsPerPage = 8;

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
    deletePostMutation.mutate(id);
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

  return (
    <div className="h-full flex flex-col gap-10 text-foreground font-sans pb-20">
      
      {/* Header Section */}
      <div className="flex flex-col md:flex-row md:items-end justify-between gap-6 shrink-0">
        <div className="space-y-2">
           <h1 className="text-4xl md:text-5xl font-serif font-medium tracking-tight">Article Archive</h1>
           <p className="text-muted-foreground font-medium">
             Manage your written works and creative thoughts.
           </p>
        </div>
        
        <Link href="/admin/posts/new">
          <Button className="rounded-full bg-primary text-primary-foreground h-12 px-8 font-bold shadow-xl shadow-primary/10 hover:scale-105 transition-transform">
            <Plus className="mr-2 h-4 w-4" /> Create New Post
          </Button>
        </Link>
      </div>

      {/* Toolbar */}
      <div className="flex flex-col sm:flex-row gap-4 shrink-0">
         <div className="relative flex-1">
            <Search className="absolute left-4 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground pointer-events-none" />
            <Input 
              placeholder="Search by title or slug..."
              className="pl-12 h-12 bg-white/50 dark:bg-slate-900/50 border-border focus-visible:ring-accent rounded-2xl"
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
            />
         </div>
         <div className="flex items-center gap-3">
             <DropdownMenu>
               <DropdownMenuTrigger asChild>
                 <Button variant="outline" className="h-12 border-border bg-white/50 dark:bg-slate-900/50 rounded-2xl min-w-[140px] justify-between font-bold text-xs uppercase tracking-widest">
                   {statusFilter}
                   <Filter className="ml-2 h-3.5 w-3.5 opacity-50" /> 
                 </Button>
               </DropdownMenuTrigger>
               <DropdownMenuContent align="end" className="w-48 bg-white dark:bg-slate-900 border-border shadow-xl rounded-xl">
                 {['all', 'published', 'draft', 'archived'].map((s) => (
                   <DropdownMenuCheckboxItem 
                    key={s}
                    checked={statusFilter === s} 
                    onCheckedChange={() => setStatusFilter(s as any)}
                    className="capitalize font-medium"
                   >
                     {s}
                   </DropdownMenuCheckboxItem>
                 ))}
               </DropdownMenuContent>
             </DropdownMenu>

             <Button variant="outline" className="h-12 border-border bg-white/50 dark:bg-slate-900/50 rounded-2xl font-bold text-xs uppercase tracking-widest">
               <SortAsc className="mr-2 h-3.5 w-3.5 opacity-50" />
               Latest
             </Button>
         </div>
      </div>

      {/* List Container */}
      <div className="flex-1 relative min-h-[400px]">
        {isLoading ? (
            <div className="flex flex-col items-center justify-center h-64">
                <Loader2 className="h-8 w-8 animate-spin text-accent" />
            </div>
        ) : filteredPosts.length === 0 ? (
            <motion.div 
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              className="flex flex-col items-center justify-center h-64 text-muted-foreground gap-4 border border-dashed rounded-3xl"
            >
                <FileText className="h-12 w-12 opacity-10" />
                <p className="font-medium">No articles found in this archive.</p>
            </motion.div>
        ) : (
            <div className="grid grid-cols-1 gap-4">
                <AnimatePresence mode="popLayout">
                    {paginatedPosts.map((post, i) => {
                        const dateObj = formatDateSafe(post.updated_at);
                        return (
                          <motion.div 
                            key={post.id} 
                            initial={{ opacity: 0, y: 10 }}
                            animate={{ opacity: 1, y: 0 }}
                            exit={{ opacity: 0, scale: 0.95 }}
                            transition={{ delay: i * 0.05 }}
                            className="group flex flex-col md:flex-row md:items-center justify-between p-6 bg-white/40 dark:bg-slate-900/40 backdrop-blur-sm border border-border rounded-3xl hover:border-accent/30 transition-all hover:shadow-xl hover:shadow-accent/5"
                          >
                              <div className="flex items-center gap-6 min-w-0">
                                  <div className="h-16 w-24 bg-muted rounded-2xl overflow-hidden shrink-0 border border-border group-hover:scale-105 transition-transform">
                                     {post.cover ? (
                                       <img src={post.cover} alt="" className="h-full w-full object-cover" />
                                     ) : (
                                       <div className="h-full w-full flex items-center justify-center">
                                         <ImageIcon className="h-6 w-6 text-muted-foreground/30" />
                                       </div>
                                     )}
                                  </div>
                                  <div className="min-w-0 space-y-1">
                                      <Link href={`/admin/posts/${post.id}/edit`}>
                                          <h2 className="text-xl font-serif font-medium text-foreground group-hover:text-accent transition-colors truncate cursor-pointer">
                                              {post.title}
                                          </h2>
                                      </Link>
                                      <div className="flex items-center gap-3 text-[10px] font-bold uppercase tracking-widest text-muted-foreground">
                                          <span className="flex items-center gap-1.5"><Calendar className="w-3 h-3" /> {dateObj ? new Intl.DateTimeFormat(locale, { month: 'long', day: '2-digit' }).format(dateObj) : '-'}</span>
                                          <span className="opacity-20">•</span>
                                          <span>By {post.author?.username || "Admin"}</span>
                                      </div>
                                  </div>
                              </div>

                              <div className="flex items-center gap-6 mt-4 md:mt-0 ml-auto md:ml-0">
                                  <Badge variant="outline" className={cn(
                                      "border-0 px-4 py-1 text-[10px] font-bold uppercase rounded-full",
                                      post.status === 1 ? "bg-accent/10 text-accent" : "bg-muted text-muted-foreground"
                                  )}>
                                      {post.status === 1 ? "Published" : "Draft"}
                                  </Badge>

                                  <DropdownMenu>
                                      <DropdownMenuTrigger asChild>
                                          <Button variant="ghost" size="icon" className="h-10 w-10 text-muted-foreground hover:text-foreground rounded-full hover:bg-muted transition-colors">
                                              <MoreHorizontal className="h-5 w-5" />
                                          </Button>
                                      </DropdownMenuTrigger>
                                      <DropdownMenuContent align="end" className="w-[180px] bg-white dark:bg-slate-900 border-border shadow-2xl rounded-2xl p-2">
                                          <Link href={`/admin/posts/${post.id}/edit`}>
                                              <DropdownMenuItem className="cursor-pointer rounded-xl py-2.5">
                                                  <Pencil className="mr-2 h-4 w-4" /> Edit Article
                                              </DropdownMenuItem>
                                          </Link>
                                          <DropdownMenuItem 
                                            onClick={() => handleDelete(post.id)} 
                                            disabled={deletePostMutation.isPending}
                                            className="text-accent focus:text-white focus:bg-accent cursor-pointer rounded-xl py-2.5"
                                          >
                                              <Trash className="mr-2 h-4 w-4" /> Delete Permanently
                                          </DropdownMenuItem>
                                      </DropdownMenuContent>
                                  </DropdownMenu>
                              </div>
                          </motion.div>
                        );
                    })}
                </AnimatePresence>
            </div>
        )}
      </div>

      {/* Footer Pagination */}
      {totalPages > 1 && (
        <div className="flex items-center justify-between pt-4 border-t border-border">
            <div className="text-[10px] font-bold uppercase tracking-widest text-muted-foreground">
                Page <span className="text-foreground">{currentPage}</span> of <span className="text-foreground">{totalPages}</span>
            </div>
            <div className="flex items-center gap-2">
                <Button 
                    variant="outline" 
                    size="icon" 
                    className="h-10 w-10 border-border bg-transparent rounded-full hover:bg-muted"
                    onClick={() => {
                      setCurrentPage(p => Math.max(1, p - 1));
                      window.scrollTo({ top: 0, behavior: 'smooth' });
                    }}
                    disabled={currentPage === 1}
                >
                    <ChevronLeft className="h-4 w-4" />
                </Button>
                
                <Button 
                    variant="outline" 
                    size="icon" 
                    className="h-10 w-10 border-border bg-transparent rounded-full hover:bg-muted"
                    onClick={() => {
                      setCurrentPage(p => Math.min(totalPages, p + 1));
                      window.scrollTo({ top: 0, behavior: 'smooth' });
                    }}
                    disabled={currentPage === totalPages}
                >
                    <ChevronRight className="h-4 w-4" />
                </Button>
            </div>
        </div>
      )}
    </div>
  );
}
