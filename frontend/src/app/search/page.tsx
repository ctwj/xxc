"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";

export default function SearchPage() {
  const router = useRouter();
  const [keyword, setKeyword] = useState("");
  const [results, setResults] = useState<any[]>([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(false);
  const [searched, setSearched] = useState(false);

  const handleSearch = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!keyword.trim()) return;

    setLoading(true);
    try {
      const response = await fetch(
        `${process.env.NEXT_PUBLIC_API_URL}/api/search?keyword=${encodeURIComponent(keyword)}`
      );
      const data = await response.json();
      setResults(data.data || []);
      setTotal(data.total || 0);
      setSearched(true);
    } catch (error) {
      console.error("Search failed:", error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="container py-8">
      <h1 className="text-3xl font-bold mb-8">搜索</h1>

      <form onSubmit={handleSearch} className="mb-8">
        <div className="flex gap-4">
          <input
            type="text"
            value={keyword}
            onChange={(e) => setKeyword(e.target.value)}
            placeholder="输入关键词搜索文章..."
            className="flex-1 rounded-md border bg-background px-4 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary"
          />
          <button
            type="submit"
            disabled={loading}
            className="rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90 disabled:opacity-50"
          >
            {loading ? "搜索中..." : "搜索"}
          </button>
        </div>
      </form>

      {searched && (
        <div>
          <p className="text-sm text-muted-foreground mb-4">
            找到 {total} 篇相关文章
          </p>

          <div className="grid gap-4">
            {results.map((article) => (
              <Link
                key={article.id}
                href={`/article/${article.slug}`}
                className="group rounded-lg border bg-card p-4 shadow-sm transition-colors hover:bg-accent"
              >
                <h2 className="font-semibold text-lg mb-2 group-hover:text-primary">
                  {article.title}
                </h2>
                {article.description && (
                  <p className="text-sm text-muted-foreground line-clamp-2">
                    {article.description}
                  </p>
                )}
                <div className="mt-4 flex items-center gap-4 text-xs text-muted-foreground">
                  <span>
                    {new Date(article.createTime * 1000).toLocaleDateString("zh-CN")}
                  </span>
                  <span>{article.views} 阅读</span>
                </div>
              </Link>
            ))}
          </div>

          {results.length === 0 && (
            <div className="text-center py-12 text-muted-foreground">
              未找到相关文章
            </div>
          )}
        </div>
      )}
    </div>
  );
}