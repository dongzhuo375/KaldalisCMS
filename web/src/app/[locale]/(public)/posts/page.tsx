"use client";

import { useEffect, useState } from "react";
import { useTranslations, useFormatter } from 'next-intl';
import { Link } from '@/i18n/routing';
import api from "@/lib/api";
import { Post } from "@/lib/types";
import { Card, CardContent, CardHeader, CardTitle, CardFooter } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Calendar, User, ArrowRight, Image as ImageIcon } from "lucide-react";

export default function PostsPage() {
  const t = useTranslations('posts');
  const format = useFormatter();
  const [posts, setPosts] = useState<Post[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchPosts = async () => {
      try {
        // api.get returns response.data which is Post[]
        const data = await api.get<Post[]>("/posts");
        // Cast to Post[] because api.get generic typing might be loose
        setPosts(data as unknown as Post[]);
      } catch (error) {
        console.error("Failed to fetch posts:", error);
      } finally {
        setLoading(false);
      }
    };

    fetchPosts();
  }, []);

  return (
    <div className="space-y-12 py-12">
      {/* Header Section */}
      <div className="text-center space-y-4">
        <h1 className="text-4xl font-extrabold tracking-tight lg:text-5xl text-slate-900 dark:text-slate-50">
          {t('title')}
        </h1>
        <p className="text-xl text-slate-600 dark:text-slate-400 max-w-2xl mx-auto">
          {t('subtitle')}
        </p>
      </div>

      {/* Content Section */}
      {loading ? (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 max-w-7xl mx-auto">
          {[...Array(6)].map((_, i) => (
            <Card key={i} className="animate-pulse border-0 shadow-sm bg-white dark:bg-slate-900">
              <div className="h-48 bg-slate-200 dark:bg-slate-800 rounded-t-xl" />
              <CardContent className="space-y-4 p-6">
                <div className="h-6 bg-slate-200 dark:bg-slate-800 rounded w-3/4" />
                <div className="h-4 bg-slate-200 dark:bg-slate-800 rounded w-1/2" />
                <div className="h-20 bg-slate-200 dark:bg-slate-800 rounded w-full" />
              </CardContent>
            </Card>
          ))}
        </div>
      ) : posts.length > 0 ? (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8 max-w-7xl mx-auto">
          {posts.map((post) => (
            <Card 
              key={post.id} 
              className="flex flex-col overflow-hidden border-0 shadow-lg hover:shadow-xl transition-all duration-300 hover:-translate-y-1 bg-white dark:bg-slate-900 ring-1 ring-slate-200 dark:ring-slate-800"
            >
              {/* Cover Image Area */}
              <div className="relative h-48 bg-slate-100 dark:bg-slate-800 flex items-center justify-center overflow-hidden">
                {post.cover ? (
                  <img 
                    src={post.cover} 
                    alt={post.title} 
                    className="w-full h-full object-cover transition-transform duration-500 hover:scale-105" 
                  />
                ) : (
                  <ImageIcon className="w-12 h-12 text-slate-300 dark:text-slate-600" />
                )}
              </div>
              
              <CardHeader className="p-6 pb-2">
                <div className="flex items-center gap-2 text-xs text-slate-500 dark:text-slate-400 mb-2">
                  <span className="flex items-center gap-1 bg-slate-100 dark:bg-slate-800 px-2 py-1 rounded-full">
                    <Calendar className="w-3 h-3" />
                    {post.created_at ? format.dateTime(new Date(post.created_at), { dateStyle: 'medium' }) : '-'}
                  </span>
                </div>
                <CardTitle className="text-xl font-bold leading-tight line-clamp-2 text-slate-900 dark:text-slate-100 group-hover:text-blue-600 dark:group-hover:text-blue-400 transition-colors">
                  <Link href={`/posts/${post.id}`} className="hover:underline decoration-blue-500/30 underline-offset-4">
                    {post.title}
                  </Link>
                </CardTitle>
              </CardHeader>
              
              <CardContent className="px-6 py-2 flex-grow">
                 <p className="text-slate-600 dark:text-slate-400 line-clamp-3 text-sm leading-relaxed">
                   {post.content || "No preview available..."}
                 </p>
              </CardContent>
              
              <CardFooter className="p-6 pt-4 border-t border-slate-100 dark:border-slate-800 flex items-center justify-between">
                <div className="flex items-center gap-2 text-sm text-slate-500 dark:text-slate-400">
                  <User className="w-4 h-4" />
                  <span>{post.author?.username || "Admin"}</span>
                </div>
                <Button variant="ghost" size="sm" asChild className="hover:bg-blue-50 dark:hover:bg-blue-900/20 hover:text-blue-600 dark:hover:text-blue-400 p-0 h-auto font-medium">
                  <Link href={`/posts/${post.id}`} className="flex items-center gap-1">
                    {t('read_more')} <ArrowRight className="w-4 h-4" />
                  </Link>
                </Button>
              </CardFooter>
            </Card>
          ))}
        </div>
      ) : (
        <div className="text-center py-20 bg-slate-50 dark:bg-slate-900 rounded-2xl border border-dashed border-slate-200 dark:border-slate-800 max-w-4xl mx-auto">
          <p className="text-slate-500 dark:text-slate-400 text-lg">{t('no_posts')}</p>
        </div>
      )}
    </div>
  );
}
