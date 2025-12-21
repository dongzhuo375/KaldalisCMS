import { NextResponse } from 'next/server'
import type { NextRequest } from 'next/server'

// 🔴 原来的名字: export function middleware(request: NextRequest)
// 🟢 新的名字: export function proxy(request: NextRequest)
export function proxy(request: NextRequest) {
  const { pathname } = request.nextUrl

  // 1. 获取认证 Token
  const token = request.cookies.get('auth_token')?.value

  // 2. 定义受保护的路径
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

// Config 保持不变
export const config = {
  matcher: ['/((?!api|_next/static|_next/image|favicon.ico|.*\\..*).*)'],
}
