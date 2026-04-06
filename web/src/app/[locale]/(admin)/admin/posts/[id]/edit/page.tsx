"use client";

import { useParams } from "next/navigation";
import { PostEditor } from "@/components/admin/post-editor";
import { useAdminPost } from "@/services/post-service";
import { Loader2, AlertCircle } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Link } from "@/i18n/routing";

export default function EditPostPage() {
  const params = useParams();
  const id = params?.id as string;
  
  const { data: post, isLoading, isError, error } = useAdminPost(id);

  if (isLoading) {
    return (
      <div className="flex h-[80vh] w-full items-center justify-center">
        <Loader2 className="h-10 w-10 animate-spin text-indigo-500" />
      </div>
    );
  }

  if (isError || !post) {
    return (
      <div className="flex h-[80vh] w-full flex-col items-center justify-center gap-4 text-center">
        <AlertCircle className="h-12 w-12 text-rose-500 opacity-50" />
        <div className="space-y-1">
          <h2 className="text-xl font-bold text-white">Failed to load post</h2>
          <p className="text-slate-400">{(error as { message?: string })?.message || "The post could not be found or retrieved."}</p>
        </div>
        <Link href="/admin/posts">
          <Button variant="outline" className="mt-4 border-slate-800">
            Back to Content
          </Button>
        </Link>
      </div>
    );
  }

  return (
    <PostEditor
      mode="edit"
      initialData={post}
    />
  );
}
