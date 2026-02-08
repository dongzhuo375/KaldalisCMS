"use client";

import { useState } from "react";
import { useRouter } from "@/i18n/routing";
import { PostEditor } from "@/components/admin/post-editor";
import api from "@/lib/api";

export default function NewPostPage() {
  const router = useRouter();
  const [isSubmitting, setIsSubmitting] = useState(false);

  const handleSubmit = async (data: any) => {
    setIsSubmitting(true);
    try {
      await api.post("/posts", data);
      // Simple feedback for now
      alert("Post created successfully!");
      router.push("/admin/posts");
    } catch (error: any) {
      console.error("Failed to create post", error);
      alert(error.response?.data?.message || "Failed to create post");
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="container mx-auto max-w-5xl py-6">
      <PostEditor 
        mode="create"
        onSubmit={handleSubmit}
        isSubmitting={isSubmitting}
      />
    </div>
  );
}
