"use client";

import { PostEditor } from "@/components/admin/post-editor";

export default function NewPostPage() {
  return (
    <div className="h-full w-full">
      <PostEditor mode="create" />
    </div>
  );
}
