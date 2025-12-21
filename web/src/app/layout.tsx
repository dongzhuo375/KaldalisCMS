import type { Metadata } from "next";
// 如果你用了 next/font，可以保留字体设置，没有的话可以去掉
import { Inter } from "next/font/google"; 
import "./globals.css"; // 👈 必须引入全局样式，否则 Tailwind 不生效

const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: "Kaldalis CMS",
  description: "A modern content management system built with Go and Next.js",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="zh-CN">
      <body className={inter.className}>
        {children}
      </body>
    </html>
  );
}
