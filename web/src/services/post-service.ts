import api from "@/lib/api";
import { Post, CreatePostDTO, UpdatePostDTO } from "@/lib/types";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";

// ============================================================================
// Query Keys - Centralized for consistency
// ============================================================================
export const postKeys = {
  // Public
  publicAll: ["public-posts"] as const,
  publicDetail: (idOrSlug: string | number) => ["public-post", idOrSlug] as const,
  // Admin
  adminAll: ["admin-posts"] as const,
  adminDetail: (id: string | number) => ["admin-post", id] as const,
};

// ============================================================================
// Public Hooks - For guest/public access (published posts only)
// ============================================================================

/**
 * Fetch published posts for public display
 */
export const usePublicPosts = () => {
  return useQuery({
    queryKey: postKeys.publicAll,
    queryFn: async () => {
      const response = await api.get("/posts");
      return response as unknown as Post[];
    },
  });
};

/**
 * Fetch a single published post by ID or slug
 */
export const usePublicPost = (idOrSlug: string | number) => {
  return useQuery({
    queryKey: postKeys.publicDetail(idOrSlug),
    queryFn: async () => {
      const response = await api.get(`/posts/${idOrSlug}`);
      return response as unknown as Post;
    },
    enabled: !!idOrSlug,
  });
};

// ============================================================================
// Admin Hooks - For authenticated users with permissions (all posts)
// ============================================================================

/**
 * Fetch all posts for admin management (includes drafts)
 */
export const useAdminPosts = () => {
  return useQuery({
    queryKey: postKeys.adminAll,
    queryFn: async () => {
      const response = await api.get("/admin/posts");
      return response as unknown as Post[];
    },
  });
};

/**
 * Fetch a single post by ID for admin editing
 */
export const useAdminPost = (id: string | number) => {
  return useQuery({
    queryKey: postKeys.adminDetail(id),
    queryFn: async () => {
      const response = await api.get(`/admin/posts/${id}`);
      return response as unknown as Post;
    },
    enabled: !!id,
  });
};

/**
 * Create a new post (always created as draft)
 */
export const useCreatePost = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (data: CreatePostDTO) => {
      const response = await api.post("/admin/posts", data);
      return response as unknown as Post;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: postKeys.adminAll });
      toast.success("Post created successfully!");
    },
  });
};

/**
 * Update an existing post
 */
export const useUpdatePost = (id: number) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (data: UpdatePostDTO) => {
      const response = await api.put(`/admin/posts/${id}`, data);
      return response as unknown as Post;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: postKeys.adminAll });
      queryClient.invalidateQueries({ queryKey: postKeys.adminDetail(id) });
      queryClient.invalidateQueries({ queryKey: postKeys.publicAll });
      toast.success("Post updated successfully!");
    },
  });
};

/**
 * Delete a post
 */
export const useDeletePost = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (id: number) => {
      return await api.delete(`/admin/posts/${id}`);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: postKeys.adminAll });
      queryClient.invalidateQueries({ queryKey: postKeys.publicAll });
      toast.success("Post deleted successfully.");
    },
  });
};

/**
 * Publish a draft post
 */
export const usePublishPost = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (id: number) => {
      return await api.post(`/admin/posts/${id}/publish`);
    },
    onSuccess: (_, id) => {
      queryClient.invalidateQueries({ queryKey: postKeys.adminAll });
      queryClient.invalidateQueries({ queryKey: postKeys.adminDetail(id) });
      queryClient.invalidateQueries({ queryKey: postKeys.publicAll });
      toast.success("Post published successfully!");
    },
  });
};

/**
 * Move a published post back to draft
 */
export const useDraftPost = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (id: number) => {
      return await api.post(`/admin/posts/${id}/draft`);
    },
    onSuccess: (_, id) => {
      queryClient.invalidateQueries({ queryKey: postKeys.adminAll });
      queryClient.invalidateQueries({ queryKey: postKeys.adminDetail(id) });
      queryClient.invalidateQueries({ queryKey: postKeys.publicAll });
      toast.success("Post moved to draft.");
    },
  });
};
