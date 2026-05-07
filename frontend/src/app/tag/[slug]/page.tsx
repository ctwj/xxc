import { api } from "@/lib/api";
import { Article } from "@/types";
import Link from "next/link";

export const revalidate = 60;

async function getTag(slug: string) {
  try {
    const tags = await api.getTags();
    return tags.find((t) => t.slug === slug);
  } catch (error) {
    console.error("Failed to fetch tag:", error);
    return null;
  }
}

async function getTagArticles(slug: string) {
  try {
    return await api.getArticles({ page: 1, pageSize: 50, tag: slug });
  } catch (error) {
    console.error("Failed to fetch tag articles:", error);
    return { data: [], total: 0 };
  }
}

export default async function TagPage({
  params,
}: {
  params: Promise<{ slug: string }>;
}) {
  const { slug } = await params;
  const tag = await getTag(slug);
  const { data: articles } = await getTagArticles(slug);

  if (!tag) {
    return (
      <div className="container py-8">
        <div className="text-center py-12 text-muted-foreground">
          标签不存在
        </div>
      </div>
    );
  }

  return (
    <div className="container py-8">
      <header className="mb-8">
        <nav className="mb-4 text-sm text-muted-foreground">
          <Link href="/" className="hover:text-foreground">
            首页
          </Link>
          <span className="mx-2">/</span>
          <Link href="/tags" className="hover:text-foreground">
            标签
          </Link>
          <span className="mx-2">/</span>
          <span className="text-foreground">{tag.name}</span>
        </nav>
        <h1 className="text-3xl font-bold">
          标签: {tag.title || tag.name}
        </h1>
        {tag.description && (
          <p className="mt-2 text-muted-foreground">{tag.description}</p>
        )}
      </header>

      <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
        {articles.map((article: Article) => (
          <Link
            key={article.id}
            href={`/article/${article.slug}`}
            className="group rounded-lg border bg-card text-card-foreground shadow-sm transition-colors hover:bg-accent"
          >
            {article.thumbnail && (
              <div className="aspect-video overflow-hidden rounded-t-lg">
                <img
                  src={article.thumbnail}
                  alt={article.title}
                  className="object-cover w-full h-full transition-transform group-hover:scale-105"
                />
              </div>
            )}
            <div className="p-4">
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
            </div>
          </Link>
        ))}
      </div>

      {articles.length === 0 && (
        <div className="text-center py-12 text-muted-foreground">
          该标签下暂无文章
        </div>
      )}
    </div>
  );
}