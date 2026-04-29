import { Metadata } from "next";

export const metadata: Metadata = {
  title: "收藏夹 - Moss CMS",
};

export default function FavoritesPage() {
  // This page requires authentication
  // For now, show a placeholder
  return (
    <div className="container py-8">
      <h1 className="text-3xl font-bold mb-8">我的收藏</h1>

      <div className="rounded-lg border bg-card p-8 text-center">
        <p className="text-muted-foreground">
          请先登录以查看收藏的文章
        </p>
      </div>
    </div>
  );
}