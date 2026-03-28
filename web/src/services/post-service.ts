import api from "@/lib/api";
import { Post } from "@/lib/types";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";

export const usePosts = (params?: Record<string, unknown>) => {
  return useQuery({
    queryKey: ["posts", params],
    queryFn: async () => {
      const endpoint = params?.admin ? "/admin/posts" : "/posts";
      const response = await api.get(endpoint, { params });
      return response as Post[];
    },
  });
};

export const usePost = (idOrSlug: string | number) => {
  return useQuery({
    queryKey: ["post", idOrSlug],
    queryFn: async () => {
      const response = await api.get(`/posts/${idOrSlug}`);
      return response as Post;
    },
    enabled: !!idOrSlug,
  });
};

export const useCreatePost = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (data: Partial<Post>) => {
      const response = await api.post("/admin/posts", data);
      return response as Post;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["posts"] });
      toast.success("Post created successfully!");
    },
  });
};

export const useUpdatePost = (id: number) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (data: Partial<Post>) => {
      const response = await api.put(`/admin/posts/${id}`, data);
      return response as Post;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["posts"] });
      queryClient.invalidateQueries({ queryKey: ["post", id] });
      toast.success("Post updated successfully!");
    },
  });
};

export const useDeletePost = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (id: number) => {
      return await api.delete(`/admin/posts/${id}`);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["posts"] });
      toast.success("Post deleted successfully.");
    },
  });
};
