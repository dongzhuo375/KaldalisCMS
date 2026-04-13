import api from "@/lib/api";
import { User, LoginDTO, RegisterDTO } from "@/lib/types";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useAuthStore } from "@/store/useAuthStore";
import { toast } from "sonner";

export const useProfile = () => {
  return useQuery({
    queryKey: ["profile"],
    queryFn: async () => {
      const response = await api.get("/users/profile");
      return response as unknown as User;
    },
    retry: false,
  });
};

export const useLogin = () => {
  const queryClient = useQueryClient();
  const setLogin = useAuthStore((state) => state.setLogin);

  return useMutation({
    mutationFn: async (data: LoginDTO) => {
      const response = await api.post("/users/login", data);
      return response as unknown as User;
    },
    onSuccess: (user) => {
      setLogin(user);
      queryClient.setQueryData(["profile"], user);
      toast.success("Login successful!");
    },
  });
};

export const useRegister = () => {
  return useMutation({
    mutationFn: async (data: RegisterDTO) => {
      return await api.post("/users/register", data);
    },
    onSuccess: () => {
      toast.success("Registration successful! You can now login.");
    },
  });
};

export const useLogout = () => {
  const queryClient = useQueryClient();
  const logout = useAuthStore((state) => state.logout);

  return useMutation({
    mutationFn: async () => {
      return await api.post("/users/logout");
    },
    onSuccess: () => {
      logout();
      queryClient.removeQueries();
      toast.success("Logged out successfully.");
    },
  });
};
