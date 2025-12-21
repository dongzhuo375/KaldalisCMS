// src/middleware.ts
import { NextResponse } from 'next/server'
import type { NextRequest } from 'next/server'

export function middleware(request: NextRequest) {
  const { pathname } = request.nextUrl

  // 🔴 错误: const token = request.cookies.get('auth_token')?.value
  // 🟢 正确: 必须和你浏览器里的名字一模一样
  const token = request.cookies.get('kaldalis_auth')?.value

  // 定义受保护的路径
  const isAuthPage = pathname === '/login' || pathname === '/register'
  const isAdminPage = pathname.startsWith('/admin')

  // Case A: 未登录进后台 -> 踢回登录页
  if (isAdminPage && !token) {
    return NextResponse.redirect(new URL('/login', request.url))
  }

  // Case B: 已登录进登录页 -> 踢回后台
  if (isAuthPage && token) {
    return NextResponse.redirect(new URL('/admin/dashboard', request.url))
  }

  return NextResponse.next()
}

export const config = {
  matcher: ['/((?!api|_next/static|_next/image|favicon.ico|.*\\..*).*)'],
}
