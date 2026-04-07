"use client";

import { useEffect, useState } from "react";
import { useTranslations, useFormatter } from 'next-intl';
import { Link } from '@/i18n/routing';
import api from "@/lib/api";
import { Post } from "@/lib/types";
import { Calendar, ArrowRight, Image as ImageIcon } from "lucide-react";
import { motion } from "framer-motion";

export default function PostsPage() {
  const t = useTranslations('posts');
  const format = useFormatter();
  const [posts, setPosts] = useState<Post[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchPosts = async () => {
      try {
        const data = await api.get<Post[]>("/posts");
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
    <div className="min-h-[calc(100vh-4rem)] px-6 md:px-12 lg:px-20 py-16">
      <div className="max-w-4xl mx-auto">
        {/* Header */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.6 }}
          className="mb-16"
        >
          <h1 className="text-4xl md:text-5xl font-serif font-medium text-foreground mb-4">
            {t('title')}
          </h1>
          <p className="text-lg text-muted-foreground">
            {t('subtitle')}
          </p>
        </motion.div>

        {/* Posts List */}
        {loading ? (
          <div className="space-y-8">
            {[...Array(4)].map((_, i) => (
              <div key={i} className="animate-pulse flex gap-6">
                <div className="w-32 h-24 bg-muted rounded-xl shrink-0" />
                <div className="flex-1 space-y-3">
                  <div className="h-6 bg-muted rounded w-3/4" />
                  <div className="h-4 bg-muted rounded w-1/4" />
                  <div className="h-4 bg-muted rounded w-full" />
                </div>
              </div>
            ))}
          </div>
        ) : posts.length > 0 ? (
          <div className="space-y-1">
            {posts.map((post, index) => (
              <motion.article
                key={post.id}
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.4, delay: index * 0.05 }}
              >
                <Link
                  href={`/posts/${post.id}`}
                  className="group flex gap-6 py-6 border-b border-border/50 hover:bg-muted/30 -mx-4 px-4 rounded-xl transition-colors"
                >
                  {/* Cover */}
                  <div className="w-28 h-20 md:w-36 md:h-24 bg-muted rounded-xl overflow-hidden shrink-0 flex items-center justify-center">
                    {post.cover ? (
                      <img
                        src={post.cover}
                        alt={post.title}
                        className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
                      />
                    ) : (
                      <ImageIcon className="w-8 h-8 text-muted-foreground/30" />
                    )}
                  </div>

                  {/* Content */}
                  <div className="flex-1 min-w-0">
                    <h2 className="text-lg md:text-xl font-serif font-medium text-foreground group-hover:text-accent transition-colors line-clamp-2 mb-2">
                      {post.title}
                    </h2>
                    <div className="flex items-center gap-3 text-sm text-muted-foreground mb-2">
                      <span className="flex items-center gap-1.5">
                        <Calendar className="w-3.5 h-3.5" />
                        {post.created_at ? format.dateTime(new Date(post.created_at), { dateStyle: 'medium' }) : '-'}
                      </span>
                      <span>·</span>
                      <span>{post.author?.username || "Admin"}</span>
                    </div>
                    <p className="text-muted-foreground text-sm line-clamp-2 hidden md:block">
                      {post.content?.slice(0, 150) || "No preview available..."}
                    </p>
                  </div>

                  {/* Arrow */}
                  <div className="hidden md:flex items-center text-muted-foreground group-hover:text-accent transition-colors">
                    <ArrowRight className="w-5 h-5" />
                  </div>
                </Link>
              </motion.article>
            ))}
          </div>
        ) : (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            className="text-center py-20 border border-dashed border-border rounded-2xl"
          >
            <p className="text-muted-foreground">{t('no_posts')}</p>
          </motion.div>
        )}
      </div>
    </div>
  );
}
