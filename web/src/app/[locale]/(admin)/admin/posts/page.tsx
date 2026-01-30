// src/app/(admin)/admin/posts/page.tsx
"use client";

import { useState, useEffect } from "react";
import Link from "next/link";
import { 
  Table, 
  TableBody, 
  TableCell, 
  TableHead, 
  TableHeader, 
  TableRow 
} from "@/components/ui/table";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { 
  DropdownMenu, 
  DropdownMenuContent, 
  DropdownMenuItem, 
  DropdownMenuLabel, 
  DropdownMenuTrigger 
} from "@/components/ui/dropdown-menu";
import { MoreHorizontal, Plus, Pencil, Trash } from "lucide-react";
import { Post } from "@/lib/types";
import api from "@/lib/api";

// 临时假数据 (联调前用这个占位)
const MOCK_POSTS: Post[] = [
  { id: 1, title: "Hello World", content: "This is the hello world post content", slug: "hello-world", status: "published", author_id: 1, author: { id: 1, username: "admin", role: "admin" }, created_at: "2023-12-01", updated_at: "2023-12-01" },
  { id: 2, title: "Next.js 14 教程", content: "This is the Next.js 14 tutorial content", slug: "nextjs-14-guide", status: "draft", author_id: 1, author: { id: 1, username: "admin", role: "admin" }, created_at: "2023-12-05", updated_at: "2023-12-06" },
];

export default function PostsPage() {
  const [posts, setPosts] = useState<Post[]>(MOCK_POSTS);
  const [loading, setLoading] = useState(true);

  // 模拟从后端获取数据
  const fetchPosts = async () => {
    setLoading(true);
    try {
      // TODO: 等后端写好 list 接口后，解开这行注释
      // const res: any = await api.get("/posts");
      // setPosts(res.data);
      
      // 暂时用假数据模拟延迟
      setTimeout(() => {
        setPosts(MOCK_POSTS); 
        setLoading(false);
      }, 500);
    } catch (error) {
      console.error("获取文章失败", error);
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchPosts();
  }, []);

  const handleDelete = async (id: number) => {
    if (!confirm("确定要删除这篇文章吗？")) return;
    console.log("正在删除:", id);
    // await api.delete(`/posts/${id}`);
    // fetchPosts(); // 刷新列表
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-3xl font-bold tracking-tight">文章管理</h2>
          <p className="text-slate-500">管理您的博客文章发布与状态。</p>
        </div>
        <Link href="/admin/posts/new">
          <Button>
            <Plus className="mr-2 h-4 w-4" /> 新建文章
          </Button>
        </Link>
      </div>

      <div className="rounded-md border bg-white">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead className="w-[100px]">ID</TableHead>
              <TableHead>标题</TableHead>
              <TableHead>状态</TableHead>
              <TableHead>作者</TableHead>
              <TableHead>发布时间</TableHead>
              <TableHead className="text-right">操作</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {posts.map((post) => (
              <TableRow key={post.id}>
                <TableCell className="font-medium">{post.id}</TableCell>
                <TableCell>{post.title}</TableCell>
                <TableCell>
                  <Badge variant={post.status === "published" ? "default" : "secondary"}>
                    {post.status === "published" ? "已发布" : "草稿"}
                  </Badge>
                </TableCell>
                <TableCell>{post.author?.username}</TableCell>
                <TableCell>{post.created_at}</TableCell>
                <TableCell className="text-right">
                  <DropdownMenu>
                    <DropdownMenuTrigger asChild>
                      <Button variant="ghost" className="h-8 w-8 p-0">
                        <span className="sr-only">打开菜单</span>
                        <MoreHorizontal className="h-4 w-4" />
                      </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent align="end">
                      <DropdownMenuLabel>操作</DropdownMenuLabel>
                      <DropdownMenuItem onClick={() => alert(`编辑 ${post.id}`)}>
                        <Pencil className="mr-2 h-4 w-4" /> 编辑
                      </DropdownMenuItem>
                      <DropdownMenuItem onClick={() => handleDelete(post.id)} className="text-red-600">
                        <Trash className="mr-2 h-4 w-4" /> 删除
                      </DropdownMenuItem>
                    </DropdownMenuContent>
                  </DropdownMenu>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </div>
    </div>
  );
}
