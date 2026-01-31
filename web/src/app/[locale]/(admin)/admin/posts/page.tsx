"use client";

import { useState, useEffect } from "react";
import Link from "next/link";
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
  Loader2
} from "lucide-react";
import { Post } from "@/lib/types";
import api from "@/lib/api";
import { Card, CardContent } from "@/components/ui/card";
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
      
      {/* Top Header Area */}
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div>
           <h1 className="text-2xl font-bold tracking-tight text-white">{t('posts')}</h1>
           <p className="text-slate-400 text-sm mt-1">Manage your blog posts, articles and content.</p>
        </div>
        
        <div className="flex items-center gap-2">
           <Button variant="outline" className="border-slate-800 bg-slate-900 text-slate-300 hover:bg-slate-800 hover:text-white">
             <Download className="mr-2 h-4 w-4" /> Export
           </Button>
           <Link href="/admin/posts/new">
            <Button className="bg-emerald-600 hover:bg-emerald-700 text-white border-0">
              <Plus className="mr-2 h-4 w-4" /> {t('create')}
            </Button>
          </Link>
        </div>
      </div>

      {/* Main Card */}
      <Card className="bg-slate-950 border-slate-800 shadow-sm flex-1 flex flex-col overflow-hidden">
        
        {/* Filters Toolbar */}
        <div className="p-4 border-b border-slate-800 flex flex-col sm:flex-row gap-4 justify-between items-center bg-slate-900/50">
           <div className="flex items-center gap-2 w-full sm:w-auto">
             <div className="relative w-full sm:w-72">
                <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-slate-500" />
                <Input 
                  placeholder={t('search')}
                  className="pl-9 bg-slate-950 border-slate-800 text-slate-200 focus-visible:ring-emerald-500/50"
                  value={searchTerm}
                  onChange={(e) => setSearchTerm(e.target.value)}
                />
             </div>
             
             <DropdownMenu>
               <DropdownMenuTrigger asChild>
                 <Button variant="outline" className="border-slate-800 bg-slate-950 text-slate-300 hover:bg-slate-800">
                   <SlidersHorizontal className="mr-2 h-4 w-4" /> {t('filter')}
                 </Button>
               </DropdownMenuTrigger>
               <DropdownMenuContent align="start" className="w-56 bg-slate-900 border-slate-800 text-slate-200">
                 <DropdownMenuLabel>Status</DropdownMenuLabel>
                 <DropdownMenuSeparator className="bg-slate-800" />
                 <DropdownMenuCheckboxItem 
                    checked={statusFilter === 'all'} 
                    onCheckedChange={() => setStatusFilter('all')}
                 >
                   All Statuses
                 </DropdownMenuCheckboxItem>
                 <DropdownMenuCheckboxItem 
                    checked={statusFilter === 'published'} 
                    onCheckedChange={() => setStatusFilter('published')}
                 >
                   Published
                 </DropdownMenuCheckboxItem>
                 <DropdownMenuCheckboxItem 
                    checked={statusFilter === 'draft'} 
                    onCheckedChange={() => setStatusFilter('draft')}
                 >
                   Draft
                 </DropdownMenuCheckboxItem>
                 <DropdownMenuCheckboxItem 
                    checked={statusFilter === 'archived'} 
                    onCheckedChange={() => setStatusFilter('archived')}
                 >
                   Archived
                 </DropdownMenuCheckboxItem>
               </DropdownMenuContent>
             </DropdownMenu>
           </div>
           
           {/* Active Filters Display could go here */}
        </div>

        {/* Table Content */}
        <div className="flex-1 overflow-auto">
            {loading ? (
                <div className="flex flex-col items-center justify-center h-64 text-slate-500 gap-3">
                  <Loader2 className="h-8 w-8 animate-spin text-emerald-500" />
                  <span className="text-sm font-mono">{t('loading')}...</span>
                </div>
            ) : filteredPosts.length === 0 ? (
                <div className="flex flex-col items-center justify-center h-64 text-slate-500 gap-3">
                    <div className="bg-slate-900 p-4 rounded-full">
                        <FileText className="h-8 w-8 opacity-40" />
                    </div>
                    <p className="text-sm font-medium text-slate-400">No posts found</p>
                    <p className="text-xs text-slate-600">Try adjusting your search or filters</p>
                </div>
            ) : (
            <Table>
                <TableHeader className="bg-slate-900/80 sticky top-0 z-10 backdrop-blur-sm">
                    <TableRow className="hover:bg-transparent border-slate-800">
                        <TableHead className="w-[40px] pl-4">
                            <Checkbox 
                                checked={paginatedPosts.length > 0 && selectedPosts.length === paginatedPosts.length}
                                onCheckedChange={toggleSelectAll}
                                className="border-slate-600 data-[state=checked]:bg-emerald-600 data-[state=checked]:border-emerald-600"
                            />
                        </TableHead>
                        <TableHead className="w-[80px] text-xs font-medium text-slate-400 uppercase tracking-wider">{t('id')}</TableHead>
                        <TableHead className="text-xs font-medium text-slate-400 uppercase tracking-wider">
                            <div className="flex items-center gap-1 cursor-pointer hover:text-white">
                                {t('title')} <ArrowUpDown className="h-3 w-3" />
                            </div>
                        </TableHead>
                        <TableHead className="w-[120px] text-xs font-medium text-slate-400 uppercase tracking-wider">{t('status')}</TableHead>
                        <TableHead className="w-[180px] text-xs font-medium text-slate-400 uppercase tracking-wider">{t('author')}</TableHead>
                        <TableHead className="w-[150px] text-xs font-medium text-slate-400 uppercase tracking-wider text-right">{t('date')}</TableHead>
                        <TableHead className="w-[60px]"></TableHead>
                    </TableRow>
                </TableHeader>
                <TableBody className="border-slate-800">
                    {paginatedPosts.map((post) => (
                        <TableRow key={post.id} className="border-slate-800 hover:bg-slate-800/40 group transition-colors data-[state=selected]:bg-slate-800/60" data-state={selectedPosts.includes(post.id) ? "selected" : ""}>
                            <TableCell className="pl-4">
                                <Checkbox 
                                    checked={selectedPosts.includes(post.id)}
                                    onCheckedChange={() => toggleSelectPost(post.id)}
                                    className="border-slate-600 data-[state=checked]:bg-emerald-600 data-[state=checked]:border-emerald-600"
                                />
                            </TableCell>
                            <TableCell className="font-mono text-xs text-slate-500">#{post.id}</TableCell>
                            <TableCell>
                                <div className="flex flex-col gap-1 py-1">
                                    <span className="font-medium text-slate-200 group-hover:text-emerald-400 transition-colors truncate max-w-[300px] md:max-w-[400px]">
                                        {post.title}
                                    </span>
                                    <span className="text-xs text-slate-600 font-mono truncate max-w-[300px] flex items-center gap-1">
                                        <Terminal className="w-3 h-3" /> /{post.slug}
                                    </span>
                                </div>
                            </TableCell>
                            <TableCell>
                                <Badge variant="outline" className={cn(
                                    "border-0 px-2 py-0.5 text-xs font-medium rounded-full flex w-fit items-center gap-1.5",
                                    post.status === 1 && "bg-emerald-500/10 text-emerald-400 ring-1 ring-emerald-500/20",
                                    post.status === 0 && "bg-yellow-500/10 text-yellow-400 ring-1 ring-yellow-500/20",
                                    post.status === 2 && "bg-slate-500/10 text-slate-400 ring-1 ring-slate-500/20"
                                )}>
                                    <div className={cn("w-1.5 h-1.5 rounded-full", 
                                        post.status === 1 ? "bg-emerald-400" : 
                                        post.status === 0 ? "bg-yellow-400" : "bg-slate-400"
                                    )} />
                                    {post.status === 1 ? t('published') : post.status === 0 ? t('draft') : t('archived')}
                                </Badge>
                            </TableCell>
                            <TableCell>
                                <div className="flex items-center gap-2.5">
                                    <div className="h-7 w-7 rounded-full bg-slate-800 flex items-center justify-center text-[10px] font-bold text-slate-300 ring-2 ring-slate-900 border border-slate-700">
                                        {post.author?.username?.[0]?.toUpperCase() || "A"}
                                    </div>
                                    <span className="text-sm text-slate-300 font-medium">{post.author?.username || "Admin"}</span>
                                </div>
                            </TableCell>
                            <TableCell className="text-slate-500 text-sm font-mono text-right">
                                {post.created_at ? format.dateTime(new Date(post.created_at), { dateStyle: 'medium' }) : '-'}
                            </TableCell>
                            <TableCell className="text-right pr-4">
                                <DropdownMenu>
                                    <DropdownMenuTrigger asChild>
                                        <Button variant="ghost" className="h-8 w-8 p-0 text-slate-500 hover:text-white hover:bg-slate-800 data-[state=open]:bg-slate-800 data-[state=open]:text-white">
                                            <span className="sr-only">Open menu</span>
                                            <MoreHorizontal className="h-4 w-4" />
                                        </Button>
                                    </DropdownMenuTrigger>
                                    <DropdownMenuContent align="end" className="w-[180px] bg-slate-900 border-slate-800 text-slate-200">
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
        </div>

        {/* Pagination Footer */}
        <div className="border-t border-slate-800 p-4 bg-slate-900/50 flex items-center justify-between">
            <div className="text-xs text-slate-500 font-mono">
                Showing <span className="text-slate-300 font-medium">{(currentPage - 1) * itemsPerPage + 1}-{Math.min(currentPage * itemsPerPage, filteredPosts.length)}</span> of <span className="text-slate-300 font-medium">{filteredPosts.length}</span>
            </div>
            <div className="flex items-center gap-2">
                <Button 
                    variant="outline" 
                    size="sm" 
                    className="h-8 w-8 p-0 border-slate-800 bg-slate-950 hover:bg-slate-900 disabled:opacity-50"
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
                              "h-8 w-8 text-xs rounded-md font-medium transition-colors",
                              currentPage === i + 1 
                                ? "bg-emerald-600 text-white" 
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
                    className="h-8 w-8 p-0 border-slate-800 bg-slate-950 hover:bg-slate-900 disabled:opacity-50"
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