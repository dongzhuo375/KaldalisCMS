"use client";

import { useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import { useTranslations, useFormatter } from 'next-intl';
import { Link } from '@/i18n/routing';
import api from "@/lib/api";
import { Post } from "@/lib/types";
import { Button } from "@/components/ui/button";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { Calendar, User, ArrowLeft, Clock, Tag } from "lucide-react";

export default function PostDetailPage() {
  const params = useParams();
  const router = useRouter();
  const t = useTranslations(); // Use generic for common keys
  const format = useFormatter();
  const [post, setPost] = useState<Post | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(false);

  // params.id might be string or string[]
  const id = Array.isArray(params.id) ? params.id[0] : params.id;

  useEffect(() => {
    if (!id) return;

    const fetchPost = async () => {
      try {
        const data = await api.get<Post>(`/posts/${id}`);
        setPost(data as unknown as Post);
      } catch (err) {
        console.error("Failed to fetch post:", err);
        setError(true);
      } finally {
        setLoading(false);
      }
    };

    fetchPost();
  }, [id]);

  if (loading) {
    return (
      <div className="max-w-4xl mx-auto py-12 px-4 space-y-8 animate-pulse">
        <div className="h-8 w-32 bg-muted rounded" />
        <div className="h-64 w-full bg-muted rounded-xl" />
        <div className="space-y-4">
           <div className="h-10 w-3/4 bg-muted rounded" />
           <div className="h-6 w-1/2 bg-muted rounded" />
           <div className="h-40 w-full bg-muted rounded" />
        </div>
      </div>
    );
  }

  if (error || !post) {
    return (
      <div className="max-w-4xl mx-auto py-20 text-center space-y-4">
        <h2 className="text-2xl font-bold text-foreground">{t('errors.not_found')}</h2>
        <Button onClick={() => router.back()} variant="outline">
          <ArrowLeft className="mr-2 h-4 w-4" /> {t('common.back')}
        </Button>
      </div>
    );
  }

  return (
    <article className="max-w-4xl mx-auto py-12 px-4 space-y-8">
      {/* Back Navigation */}
      <div>
        <Link href="/posts">
          <Button variant="ghost" className="pl-0 hover:pl-2 transition-all text-muted-foreground hover:text-foreground">
            <ArrowLeft className="mr-2 h-4 w-4" /> {t('common.back')}
          </Button>
        </Link>
      </div>

      {/* Header */}
      <header className="space-y-6">
        <div className="flex flex-wrap items-center gap-3 text-sm text-muted-foreground">
          <Badge variant="secondary" className="rounded-full px-3">
            {post.category?.name || "Uncategorized"}
          </Badge>
          <span className="flex items-center gap-1">
            <Calendar className="w-4 h-4" />
            {post.created_at ? format.dateTime(new Date(post.created_at), { dateStyle: 'long' }) : '-'}
          </span>
          {/* Estimated read time placeholder */}
          <span className="flex items-center gap-1">
            <Clock className="w-4 h-4" />
            5 min read
          </span>
        </div>

        <h1 className="text-4xl md:text-5xl font-extrabold tracking-tight text-foreground leading-tight">
          {post.title}
        </h1>

        <div className="flex items-center gap-4 py-4 border-b border-border">
          <Avatar className="h-12 w-12 border-2 border-background shadow-sm">
            <AvatarImage src={post.author?.avatar} alt={post.author?.username} />
            <AvatarFallback>{post.author?.username?.[0]?.toUpperCase() || "A"}</AvatarFallback>
          </Avatar>
          <div>
             <div className="font-semibold text-foreground">{post.author?.username || "Unknown Author"}</div>
             <div className="text-sm text-muted-foreground">{post.author?.role || "Contributor"}</div>
          </div>
        </div>
      </header>

      {/* Featured Image */}
      {post.cover && (
        <div className="relative w-full aspect-video rounded-2xl overflow-hidden shadow-lg bg-muted">
          <img 
            src={post.cover} 
            alt={post.title} 
            className="w-full h-full object-cover"
          />
        </div>
      )}

      {/* Content */}
      <div className="prose prose-slate dark:prose-invert max-w-none lg:prose-lg leading-loose">
        {/* Simple rendering for now, can be replaced with ReactMarkdown */}
        <div className="whitespace-pre-wrap font-serif text-foreground/90">
          {post.content}
        </div>
      </div>

      {/* Tags */}
      {post.tags && post.tags.length > 0 && (
        <div className="pt-8 border-t border-border">
          <div className="flex flex-wrap gap-2">
            {post.tags.map((tag: any, i: number) => (
              <Badge key={i} variant="outline" className="text-muted-foreground border-border hover:bg-muted">
                # {tag.name || tag}
              </Badge>
            ))}
          </div>
        </div>
      )}
    </article>
  );
}
