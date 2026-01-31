"use client";

import { useState, useEffect } from "react";
import { useRouter } from "@/i18n/routing";
import { useParams } from "next/navigation";
import { PostEditor } from "@/components/admin/post-editor";
import api from "@/lib/api";
import { Loader2 } from "lucide-react";
import { Post } from "@/lib/types";

export default function EditPostPage() {
  const router = useRouter();
  const params = useParams();
  const id = params?.id as string;
  
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [isLoading, setIsLoading] = useState(true);
  const [post, setPost] = useState<Post | null>(null);

  useEffect(() => {
    if (!id) return;
    
    const fetchPost = async () => {
      try {
        const data = await api.get<Post>(`/posts/${id}`);
        setPost(data as unknown as Post);
      } catch (error) {
        console.error("Failed to fetch post", error);
        alert("Failed to load post data");
        router.push("/admin/posts");
      } finally {
        setIsLoading(false);
      }
    };

    fetchPost();
  }, [id, router]);

  const handleSubmit = async (data: any) => {
    setIsSubmitting(true);
    try {
      await api.put(`/posts/${id}`, data);
      alert("Post updated successfully!");
      router.push("/admin/posts");
    } catch (error: any) {
      console.error("Failed to update post", error);
      alert(error.response?.data?.message || "Failed to update post");
    } finally {
      setIsSubmitting(false);
    }
  };

  if (isLoading) {
    return (
      <div className="flex h-96 items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-emerald-500" />
      </div>
    );
  }

  if (!post) return null;

  return (
    <div className="container mx-auto max-w-5xl py-6">
      <PostEditor 
        mode="edit"
        initialData={post}
        onSubmit={handleSubmit}
        isSubmitting={isSubmitting}
      />
    </div>
  );
}
