import { clsx, type ClassValue } from "clsx"
import { twMerge } from "tailwind-merge"
import { API_BASE_URL } from "./api"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function getImageUrl(url: string | undefined | null) {
  if (!url) return "";
  if (url.startsWith('http')) return url;
  // Use relative path to take advantage of Next.js proxy rewrite
  return url.startsWith('/') ? url : `/${url}`;
}
