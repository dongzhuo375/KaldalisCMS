import api, { sysApi } from "@/lib/api";
import { SystemStatus, SetupDTO, HealthResponse } from "@/lib/types";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";

export interface CheckDBRequest {
  host: string;
  port: number;
  user: string;
  pass: string;
  name: string;
}

export const useSystemStatus = () => {
  return useQuery({
    queryKey: ["system-status"],
    queryFn: async () => {
      const response = await api.get("/system/status");
      return response as unknown as SystemStatus;
    },
    refetchOnWindowFocus: false,
  });
};

export const useReadyz = () => {
  return useQuery({
    queryKey: ["readyz"],
    queryFn: async () => {
      const response = await sysApi.get("/readyz");
      return response as unknown as HealthResponse;
    },
    refetchInterval: 5000,
  });
};

export const useCheckDB = () => {
  return useMutation({
    mutationFn: async (data: CheckDBRequest) => {
      const response = await api.post("/system/check-db", data);
      return response;
    },
    onSuccess: () => {
      toast.success("Database connection check passed!");
    },
  });
};

export const useSetup = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (data: Record<string, unknown>) => {
      return await api.post("/system/setup", data);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["system-status"] });
      toast.success("System installation succeeded! Restarting...");
    },
  });
};
