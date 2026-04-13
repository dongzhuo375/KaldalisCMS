import api from "@/lib/api";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import { MediaAsset, MediaListResponse, MediaUploadResponse } from "@/lib/types";

export const useMedia = (params?: Record<string, unknown>) => {
  return useQuery({
    queryKey: ["media", params],
    queryFn: async () => {
      const response = await api.get("/media", { params });
      return response as unknown as MediaListResponse;
    },
  });
};

export const useUploadMedia = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (file: File) => {
      const formData = new FormData();
      formData.append("file", file);
      const response = await api.post("/media", formData, {
        headers: {
          "Content-Type": "multipart/form-data",
        },
      });
      return (response as unknown as MediaUploadResponse).asset;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["media"] });
      toast.success("File uploaded successfully!");
    },
  });
};

export const useDeleteMedia = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (id: number) => {
      return await api.delete(`/media/${id}`);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["media"] });
      toast.success("File deleted.");
    },
  });
};
