import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";
import { AuthProvider } from "@/contexts/AuthContext";

const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: "Moss CMS",
  description: "A modern content management system",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="zh-CN">
      <body className={`${inter.className} antialiased`}>
        <AuthProvider>
          <header className="sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
            <div className="container flex h-14 items-center">
              <a href="/" className="mr-6 flex items-center space-x-2">
                <span className="font-bold text-xl">Moss CMS</span>
              </a>
              <nav className="flex flex-1 items-center justify-end space-x-4">
                <a href="/" className="text-sm font-medium text-muted-foreground hover:text-foreground">
                  首页
                </a>
                <a href="/categories" className="text-sm font-medium text-muted-foreground hover:text-foreground">
                  分类
                </a>
                <a href="/tags" className="text-sm font-medium text-muted-foreground hover:text-foreground">
                  标签
                </a>
                <a href="/search" className="text-sm font-medium text-muted-foreground hover:text-foreground">
                  搜索
                </a>
                <a href="/favorites" className="text-sm font-medium text-muted-foreground hover:text-foreground">
                  收藏
                </a>
                <a href="/login" className="text-sm font-medium text-muted-foreground hover:text-foreground">
                  登录
                </a>
              </nav>
            </div>
          </header>
          <main className="min-h-screen">
            {children}
          </main>
          <footer className="border-t py-6 md:py-8">
            <div className="container flex flex-col items-center justify-center gap-4 text-center text-sm text-muted-foreground">
              <p>© 2024 Moss CMS. All rights reserved.</p>
            </div>
          </footer>
        </AuthProvider>
      </body>
    </html>
  );
}