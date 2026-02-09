"use client";

import { useState } from "react";
import { useTranslations } from 'next-intl';
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Avatar, AvatarImage, AvatarFallback } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { 
  Search, 
  Plus, 
  MoreHorizontal, 
  Trash2, 
  Shield, 
  User as UserIcon,
  Mail,
  Calendar,
  Edit2
} from "lucide-react";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";

// Mock Data
const MOCK_USERS = [
  { id: 1, username: "admin", email: "admin@kaldalis.com", role: "super_admin", avatar: "", joined: "2023-01-01" },
  { id: 2, username: "sarah_editor", email: "sarah@kaldalis.com", role: "admin", avatar: "https://images.unsplash.com/photo-1494790108377-be9c29b29330", joined: "2023-03-15" },
  { id: 3, username: "mike_writer", email: "mike@kaldalis.com", role: "user", avatar: "https://images.unsplash.com/photo-1599566150163-29194dcaad36", joined: "2023-04-02" },
  { id: 4, username: "jessica_design", email: "jessica@kaldalis.com", role: "user", avatar: "https://images.unsplash.com/photo-1580489944761-15a19d654956", joined: "2023-05-20" },
  { id: 5, username: "alex_dev", email: "alex@kaldalis.com", role: "user", avatar: "", joined: "2023-06-10" },
];

export default function UsersPage() {
  const t = useTranslations('admin');
  const [users, setUsers] = useState(MOCK_USERS);
  const [searchTerm, setSearchTerm] = useState("");

  const filteredUsers = users.filter(user => 
    user.username.toLowerCase().includes(searchTerm.toLowerCase()) ||
    user.email.toLowerCase().includes(searchTerm.toLowerCase())
  );

  const handleDelete = (id: number) => {
    if (confirm("Are you sure you want to delete this user?")) {
      setUsers(users.filter(u => u.id !== id));
    }
  };

  const getRoleBadge = (role: string) => {
    switch(role) {
      case 'super_admin':
        return <Badge className="bg-purple-500/10 text-purple-400 hover:bg-purple-500/20 border-purple-500/20">Super Admin</Badge>;
      case 'admin':
        return <Badge className="bg-blue-500/10 text-blue-400 hover:bg-blue-500/20 border-blue-500/20">Admin</Badge>;
      default:
        return <Badge variant="outline" className="border-slate-700 text-slate-400">User</Badge>;
    }
  };

  return (
    <div className="h-full flex flex-col gap-6 text-slate-200 font-sans">
      {/* Header */}
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4 shrink-0">
        <div>
           <h1 className="text-2xl font-bold tracking-tight text-white mb-1">{t('user_management')}</h1>
           <p className="text-slate-400 text-sm">
             Manage user access and roles.
           </p>
        </div>
        
        <Button className="h-10 bg-[#ad2bee] hover:bg-[#9225c9] text-white border-0 shadow-[0_4px_12px_rgba(173,43,238,0.3)] transition-all hover:scale-105 font-medium px-6">
            <Plus className="mr-2 h-4 w-4" /> {t('add_user')}
        </Button>
      </div>

      {/* Toolbar */}
      <div className="flex flex-col sm:flex-row gap-3 shrink-0">
         <div className="relative flex-1">
            <Search className="absolute left-3 top-2.5 h-4 w-4 text-slate-500 pointer-events-none" />
            <Input 
              placeholder={t('search')}
              className="pl-10 h-10 bg-[#0d0b14]/50 border-slate-800 text-slate-200 focus-visible:ring-[#ad2bee]/30 rounded-lg"
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
            />
         </div>
      </div>

      {/* Table */}
      <div className="bg-[#0d0b14]/40 border border-slate-800/60 rounded-xl overflow-hidden flex-1 shadow-2xl relative">
        <div className="w-full overflow-auto">
          <table className="w-full text-left text-sm">
            <thead className="bg-[#0d0b14]/20 border-b border-slate-800/60 text-[11px] font-bold text-slate-500 uppercase tracking-wider">
              <tr>
                <th className="px-6 py-3">User</th>
                <th className="px-6 py-3">{t('role')}</th>
                <th className="px-6 py-3">Status</th>
                <th className="px-6 py-3 text-right">{t('actions')}</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-800/40">
              {filteredUsers.map((user) => (
                <tr key={user.id} className="hover:bg-white/[0.02] transition-colors group">
                  <td className="px-6 py-4">
                    <div className="flex items-center gap-3">
                      <Avatar className="h-9 w-9 border border-slate-700">
                        <AvatarImage src={user.avatar} />
                        <AvatarFallback className="bg-slate-800 text-slate-300 font-bold">
                          {user.username[0].toUpperCase()}
                        </AvatarFallback>
                      </Avatar>
                      <div className="flex flex-col">
                        <span className="font-medium text-slate-200 group-hover:text-white transition-colors">{user.username}</span>
                        <div className="flex items-center gap-2 text-xs text-slate-500">
                           <Mail className="h-3 w-3" />
                           {user.email}
                        </div>
                      </div>
                    </div>
                  </td>
                  <td className="px-6 py-4">
                    {getRoleBadge(user.role)}
                  </td>
                  <td className="px-6 py-4">
                     <span className="inline-flex items-center gap-1.5 px-2 py-1 rounded-full text-[10px] font-medium bg-emerald-500/10 text-emerald-400 border border-emerald-500/20">
                        <span className="relative flex h-1.5 w-1.5 rounded-full bg-emerald-500"></span>
                        Active
                     </span>
                  </td>
                  <td className="px-6 py-4 text-right">
                    <DropdownMenu>
                      <DropdownMenuTrigger asChild>
                        <Button variant="ghost" size="icon" className="h-8 w-8 text-slate-500 hover:text-white hover:bg-white/5">
                          <MoreHorizontal className="h-4 w-4" />
                        </Button>
                      </DropdownMenuTrigger>
                      <DropdownMenuContent align="end" className="bg-[#1e1b24] border-slate-800 text-slate-200">
                        <DropdownMenuLabel className="text-xs text-slate-500 uppercase">Actions</DropdownMenuLabel>
                        <DropdownMenuItem className="cursor-pointer focus:bg-slate-800 focus:text-white">
                          <Edit2 className="mr-2 h-3.5 w-3.5" /> {t('edit')}
                        </DropdownMenuItem>
                        <DropdownMenuItem onClick={() => handleDelete(user.id)} className="text-rose-400 focus:text-rose-300 focus:bg-rose-950/30 cursor-pointer">
                          <Trash2 className="mr-2 h-3.5 w-3.5" /> {t('delete')}
                        </DropdownMenuItem>
                      </DropdownMenuContent>
                    </DropdownMenu>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}