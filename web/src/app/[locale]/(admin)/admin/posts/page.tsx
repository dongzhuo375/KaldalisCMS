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
  List
} from "lucide-react";
import { Post } from "@/lib/types";
import api from "@/lib/api";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
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
  const [posts, setPosts] = useState<Post[]>([]);
  const [loading, setLoading] = useState(true);
  const [searchTerm, setSearchTerm] = useState("");
  const [statusFilter, setStatusFilter] = useState<"all" | "published" | "draft" | "archived">("all");
  const [selectedPosts, setSelectedPosts] = useState<number[]>([]);
  const [viewMode, setViewMode] = useState<"list" | "grid">("list");

  // Pagination (mock)
  const [currentPage, setCurrentPage] = useState(1);
  const itemsPerPage = 10;

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

  // Calculate stats
  const stats = {
    total: posts.length,
    published: posts.filter(p => p.status === 1).length,
    drafts: posts.filter(p => p.status === 0).length,
  };

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
    <div className="h-full flex flex-col space-y-6">
      
      {/* 1. Header Section */}
      <div className="flex flex-col md:flex-row md:items-end justify-between gap-4">
        <div>
           <div className="flex items-center gap-2 mb-1">
             <h1 className="text-3xl font-bold tracking-tight text-white">{t('posts')}</h1>
             <Badge variant="outline" className="border-emerald-500/30 text-emerald-400 bg-emerald-500/10 font-mono text-xs">
               v2.0
             </Badge>
           </div>
           <p className="text-slate-400 text-sm max-w-lg">
             Manage your blog content, track status, and publish updates.
           </p>
        </div>
        
        <div className="flex items-center gap-3">
           <Link href="/admin/posts/new">
            <Button className="h-10 bg-emerald-600 hover:bg-emerald-500 text-white border-0 shadow-[0_0_20px_rgba(16,185,129,0.2)] transition-all hover:scale-105">
              <Plus className="mr-2 h-4 w-4" /> {t('create')}
            </Button>
          </Link>
        </div>
      </div>

      {/* 2. Stats Grid (Optional but adds layout depth) */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <Card className="bg-slate-900/40 border-slate-800/60 backdrop-blur-sm">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium text-slate-400">Total Posts</CardTitle>
            <FileText className="h-4 w-4 text-emerald-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-white">{stats.total}</div>
          </CardContent>
        </Card>
        <Card className="bg-slate-900/40 border-slate-800/60 backdrop-blur-sm">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium text-slate-400">Published</CardTitle>
            <CheckCircle2 className="h-4 w-4 text-emerald-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-white">{stats.published}</div>
          </CardContent>
        </Card>
        <Card className="bg-slate-900/40 border-slate-800/60 backdrop-blur-sm">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium text-slate-400">Drafts</CardTitle>
            <FileClock className="h-4 w-4 text-yellow-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-white">{stats.drafts}</div>
          </CardContent>
        </Card>
      </div>

      {/* 3. Toolbar Card (Separated) */}
      <Card className="bg-slate-900/60 border-slate-800/60 backdrop-blur-md p-4 flex flex-col sm:flex-row justify-between items-center gap-4">
         <div className="flex items-center gap-3 w-full sm:w-auto">
             <div className="relative w-full sm:w-72">
                <Search className="absolute left-3 top-2.5 h-4 w-4 text-slate-500 pointer-events-none" />
                <Input 
                  placeholder={t('search')}
                  className="pl-10 h-10 bg-slate-950/50 border-slate-800 text-slate-200 focus-visible:ring-emerald-500/30"
                  value={searchTerm}
                  onChange={(e) => setSearchTerm(e.target.value)}
                />
             </div>
             <DropdownMenu>
               <DropdownMenuTrigger asChild>
                 <Button variant="outline" className="h-10 border-slate-800 bg-slate-950/50 text-slate-300 hover:bg-slate-800 hover:text-white">
                   <Filter className="mr-2 h-4 w-4" /> 
                   {statusFilter === 'all' ? t('filter') : t(statusFilter as any)}
                 </Button>
               </DropdownMenuTrigger>
               <DropdownMenuContent align="start" className="w-56 bg-slate-900 border-slate-800 text-slate-200">
                 <DropdownMenuLabel>Status</DropdownMenuLabel>
                 <DropdownMenuSeparator className="bg-slate-800" />
                 <DropdownMenuCheckboxItem checked={statusFilter === 'all'} onCheckedChange={() => setStatusFilter('all')}>All</DropdownMenuCheckboxItem>
                 <DropdownMenuCheckboxItem checked={statusFilter === 'published'} onCheckedChange={() => setStatusFilter('published')}>Published</DropdownMenuCheckboxItem>
                 <DropdownMenuCheckboxItem checked={statusFilter === 'draft'} onCheckedChange={() => setStatusFilter('draft')}>Draft</DropdownMenuCheckboxItem>
                 <DropdownMenuCheckboxItem checked={statusFilter === 'archived'} onCheckedChange={() => setStatusFilter('archived')}>Archived</DropdownMenuCheckboxItem>
               </DropdownMenuContent>
             </DropdownMenu>
         </div>
         
         <div className="flex items-center gap-2">
            <div className="bg-slate-950/50 rounded-lg p-1 border border-slate-800 hidden sm:flex">
                <button 
                  onClick={() => setViewMode('list')}
                  className={cn(
                    "p-2 rounded-md transition-all",
                    viewMode === 'list' ? "bg-slate-800 text-white shadow-sm" : "text-slate-500 hover:text-slate-300"
                  )}
                >
                    <List className="w-4 h-4" />
                </button>
                <button 
                  onClick={() => setViewMode('grid')}
                  className={cn(
                    "p-2 rounded-md transition-all",
                    viewMode === 'grid' ? "bg-slate-800 text-white shadow-sm" : "text-slate-500 hover:text-slate-300"
                  )}
                >
                    <LayoutGrid className="w-4 h-4" />
                </button>
            </div>
            <Button variant="outline" className="h-10 border-slate-800 bg-slate-950/50 text-slate-400 hover:bg-slate-900 hover:text-white">
              <Download className="mr-2 h-4 w-4" /> Export
            </Button>
         </div>
      </Card>

      {/* 4. Content Area (Separated) */}
      <Card className="bg-slate-900/60 border-slate-800/60 backdrop-blur-md flex-1 overflow-hidden flex flex-col shadow-xl">
        {/* Table Content */}
        <div className="flex-1 overflow-auto relative min-h-[400px]">
            {loading ? (
                <div className="absolute inset-0 flex flex-col items-center justify-center z-20">
                    <Loader2 className="h-10 w-10 animate-spin text-emerald-500" />
                </div>
            ) : filteredPosts.length === 0 ? (
                <div className="flex flex-col items-center justify-center h-full text-slate-500 gap-4">
                    <div className="bg-slate-900/50 p-6 rounded-full border border-slate-800">
                        <FileText className="h-10 w-10 opacity-30" />
                    </div>
                    <p>No posts found matching your filters.</p>
                </div>
            ) : (
            <Table>
                <TableHeader className="bg-slate-950/50 sticky top-0 z-10 backdrop-blur-sm">
                    <TableRow className="hover:bg-transparent border-slate-800/80">
                        <TableHead className="w-[40px] pl-6">
                            <Checkbox 
                                checked={paginatedPosts.length > 0 && selectedPosts.length === paginatedPosts.length}
                                onCheckedChange={toggleSelectAll}
                                className="border-slate-600 data-[state=checked]:bg-emerald-600 data-[state=checked]:border-emerald-600"
                            />
                        </TableHead>
                        <TableHead className="w-[80px] text-xs font-semibold text-slate-400 uppercase tracking-wider">{t('id')}</TableHead>
                        <TableHead className="text-xs font-semibold text-slate-400 uppercase tracking-wider">
                            <div className="flex items-center gap-2 cursor-pointer hover:text-emerald-400 transition-colors group">
                                {t('title')} <ArrowUpDown className="h-3 w-3 text-slate-600 group-hover:text-emerald-500" />
                            </div>
                        </TableHead>
                        <TableHead className="w-[140px] text-xs font-semibold text-slate-400 uppercase tracking-wider">{t('status')}</TableHead>
                        <TableHead className="w-[200px] text-xs font-semibold text-slate-400 uppercase tracking-wider">{t('author')}</TableHead>
                        <TableHead className="w-[150px] text-xs font-semibold text-slate-400 uppercase tracking-wider text-right">{t('date')}</TableHead>
                        <TableHead className="w-[60px]"></TableHead>
                    </TableRow>
                </TableHeader>
                <TableBody className="border-slate-800/50">
                    {paginatedPosts.map((post, index) => (
                        <TableRow 
                          key={post.id} 
                          className={cn(
                            "border-slate-800/50 hover:bg-slate-800/30 group transition-all duration-200",
                            selectedPosts.includes(post.id) && "bg-emerald-950/10 hover:bg-emerald-950/20"
                          )}
                          style={{
                            animation: `fadeIn 0.3s ease-out forwards`,
                            animationDelay: `${index * 0.05}s`,
                            opacity: 0 
                          }}
                        >
                            <TableCell className="pl-6">
                                <Checkbox 
                                    checked={selectedPosts.includes(post.id)}
                                    onCheckedChange={() => toggleSelectPost(post.id)}
                                    className="border-slate-600 data-[state=checked]:bg-emerald-600 data-[state=checked]:border-emerald-600"
                                />
                            </TableCell>
                            <TableCell className="font-mono text-xs text-slate-500 group-hover:text-slate-400 transition-colors">#{post.id}</TableCell>
                            <TableCell>
                                <div className="flex flex-col gap-1.5 py-1">
                                    <span className="font-medium text-[15px] text-slate-200 group-hover:text-emerald-400 transition-colors truncate max-w-[300px] md:max-w-[400px]">
                                        {post.title}
                                    </span>
                                    <span className="text-xs text-slate-600 font-mono truncate max-w-[300px] flex items-center gap-1.5 group-hover:text-slate-500">
                                        <Terminal className="w-3 h-3" /> 
                                        <span>/posts/{post.slug}</span>
                                    </span>
                                </div>
                            </TableCell>
                            <TableCell>
                                <Badge variant="outline" className={cn(
                                    "border-0 px-2.5 py-1 text-xs font-medium rounded-full flex w-fit items-center gap-2 transition-transform group-hover:scale-105",
                                    post.status === 1 && "bg-emerald-500/10 text-emerald-400 ring-1 ring-emerald-500/20",
                                    post.status === 0 && "bg-yellow-500/10 text-yellow-400 ring-1 ring-yellow-500/20",
                                    post.status === 2 && "bg-slate-500/10 text-slate-400 ring-1 ring-slate-500/20"
                                )}>
                                    <span className={cn("relative flex h-1.5 w-1.5")}>
                                      {post.status === 1 && <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>}
                                      <span className={cn("relative inline-flex rounded-full h-1.5 w-1.5", 
                                          post.status === 1 ? "bg-emerald-500" : 
                                          post.status === 0 ? "bg-yellow-500" : "bg-slate-500"
                                      )}></span>
                                    </span>
                                    {post.status === 1 ? t('published') : post.status === 0 ? t('draft') : t('archived')}
                                </Badge>
                            </TableCell>
                            <TableCell>
                                <div className="flex items-center gap-3">
                                    <div className="h-8 w-8 rounded-full bg-gradient-to-br from-slate-700 to-slate-800 flex items-center justify-center text-[10px] font-bold text-slate-300 ring-2 ring-slate-900 border border-slate-700/50 shadow-sm">
                                        {post.author?.username?.[0]?.toUpperCase() || "A"}
                                    </div>
                                    <div className="flex flex-col">
                                      <span className="text-sm text-slate-200 font-medium">{post.author?.username || "Admin"}</span>
                                    </div>
                                </div>
                            </TableCell>
                            <TableCell className="text-slate-400 text-sm font-mono text-right group-hover:text-slate-300 transition-colors">
                                {post.created_at ? format.dateTime(new Date(post.created_at), { dateStyle: 'medium' }) : '-'}
                            </TableCell>
                            <TableCell className="text-right pr-4">
                                <DropdownMenu>
                                    <DropdownMenuTrigger asChild>
                                        <Button variant="ghost" className="h-8 w-8 p-0 text-slate-500 hover:text-white hover:bg-slate-800 data-[state=open]:bg-slate-800 data-[state=open]:text-white rounded-md transition-colors">
                                            <span className="sr-only">Open menu</span>
                                            <MoreHorizontal className="h-4 w-4" />
                                        </Button>
                                    </DropdownMenuTrigger>
                                    <DropdownMenuContent align="end" className="w-[180px] bg-slate-900 border-slate-800 text-slate-200 shadow-xl">
                                        <DropdownMenuLabel className="text-xs text-slate-500 font-mono uppercase tracking-wider">Actions</DropdownMenuLabel>
                                        <Link href={`/admin/posts/${post.id}/edit`}>
                                            <DropdownMenuItem className="cursor-pointer focus:bg-slate-800 focus:text-white">
                                                <Pencil className="mr-2 h-4 w-4" /> {t('cmd_edit')}
                                            </DropdownMenuItem>
                                        </Link>
                                        <Link href={`/posts/${post.id}`} target="_blank">
                                            <DropdownMenuItem className="cursor-pointer focus:bg-slate-800 focus:text-white">
                                                <FileText className="mr-2 h-4 w-4" /> {t('cmd_view')}
                                            </DropdownMenuItem>
                                        </Link>
                                        <DropdownMenuSeparator className="bg-slate-800" />
                                        <DropdownMenuItem onClick={() => handleDelete(post.id)} className="text-rose-400 focus:text-rose-300 focus:bg-rose-950/30 cursor-pointer">
                                            <Trash className="mr-2 h-4 w-4" /> {t('cmd_delete')}
                                        </DropdownMenuItem>
                                    </DropdownMenuContent>
                                </DropdownMenu>
                            </TableCell>
                        </TableRow>
                    ))}
                </TableBody>
            </Table>
            )}
            
             {/* Inline Style for keyframes */}
             <style jsx global>{`
              @keyframes fadeIn {
                from { opacity: 0; transform: translateY(10px); }
                to { opacity: 1; transform: translateY(0); }
              }
            `}</style>
        </div>

        {/* Pagination Footer */}
        <div className="border-t border-slate-800/60 p-4 bg-slate-900/40 flex items-center justify-between backdrop-blur-sm">
            <div className="text-xs text-slate-500 font-mono">
                Showing <span className="text-emerald-400 font-medium">{(currentPage - 1) * itemsPerPage + 1}-{Math.min(currentPage * itemsPerPage, filteredPosts.length)}</span> of <span className="text-emerald-400 font-medium">{filteredPosts.length}</span>
            </div>
            <div className="flex items-center gap-2">
                <Button 
                    variant="outline" 
                    size="sm" 
                    className="h-8 w-8 p-0 border-slate-800 bg-slate-950 hover:bg-slate-900 disabled:opacity-50 transition-colors"
                    onClick={() => setCurrentPage(p => Math.max(1, p - 1))}
                    disabled={currentPage === 1}
                >
                    <ChevronLeft className="h-4 w-4" />
                </Button>
                <div className="flex items-center gap-1">
                   {Array.from({ length: totalPages }).map((_, i) => (
                       <button
                          key={i}
                          onClick={() => setCurrentPage(i + 1)}
                          className={cn(
                              "h-8 w-8 text-xs rounded-md font-medium transition-all duration-200",
                              currentPage === i + 1 
                                ? "bg-emerald-600 text-white shadow-lg shadow-emerald-900/50 scale-105" 
                                : "text-slate-500 hover:bg-slate-800 hover:text-slate-300"
                          )}
                       >
                           {i + 1}
                       </button>
                   ))}
                </div>
                <Button 
                    variant="outline" 
                    size="sm" 
                    className="h-8 w-8 p-0 border-slate-800 bg-slate-950 hover:bg-slate-900 disabled:opacity-50 transition-colors"
                    onClick={() => setCurrentPage(p => Math.min(totalPages, p + 1))}
                    disabled={currentPage === totalPages}
                >
                    <ChevronRight className="h-4 w-4" />
                </Button>
            </div>
        </div>
      </Card>
    </div>
  );
}