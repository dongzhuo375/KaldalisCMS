import { NextResponse } from 'next/server'
import type { NextRequest } from 'next/server'

export function middleware(request: NextRequest) {
  const { pathname } = request.nextUrl
  
  // 调试日志：确认中间件在工作
  console.log("🛑 中间件正在检查:", pathname); 

  // 1. 获取 Token (身份凭证)
  const token = request.cookies.get('kaldalis_auth')?.value
  
  // 2. 获取 Role (关键！需要后端配合 Set-Cookie "kaldalis_role")
  // 如果后端没种这个 Cookie，默认当作普通用户处理
  const role = request.cookies.get('kaldalis_role')?.value

  // 定义路径特征
  const isAuthPage = pathname === '/login' || pathname === '/register'
  const isAdminPage = pathname.startsWith('/admin')

  // --- 场景 A: 保护后台 (Admin Area) ---
  if (isAdminPage) {
    // 1. 根本没登录 -> 滚去登录
    if (!token) {
      return NextResponse.redirect(new URL('/login', request.url))
    }
    // 2. 登录了，但角色不对 (是普通 User) -> 滚去首页
    if (role !== 'admin' && role !== 'super_admin') {
      return NextResponse.redirect(new URL('/', request.url))
    }
  }

  // --- 场景 B: 自动跳转 (Auth Pages) ---
  // 已登录用户手贱去访问 /login，根据角色自动分流
  if (isAuthPage && token) {
    if (role === 'admin' || role === 'super_admin') {
      return NextResponse.redirect(new URL('/admin/dashboard', request.url))
    } else {
      return NextResponse.redirect(new URL('/', request.url))
    }
  }

  return NextResponse.next()
}

// 排除静态资源和 API，只拦截页面
export const config = {
  matcher: ['/((?!api|_next/static|_next/image|favicon.ico|.*\\..*).*)'],
}
