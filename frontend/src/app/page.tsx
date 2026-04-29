import { Article, Category } from "@/types";
import { api } from "@/lib/api";
import Link from "next/link";

// ISR: Revalidate every 60 seconds
export const revalidate = 60;

async function getArticles() {
  try {
    const response = await api.getArticles({ page: 1, pageSize: 20 });
    return response;
  } catch (error) {
    console.error("Failed to fetch articles:", error);
    return { data: [], total: 0 };
  }
}

async function getCategories() {
  try {
    return await api.getCategories();
  } catch (error) {
    console.error("Failed to fetch categories:", error);
    return [];
  }
}

export default async function HomePage() {
  const { data: articles, total } = await getArticles();
  const categories = await getCategories();

  return (
    <div className="container py-8">
      {/* Hero Section */}
      <section className="mb-12">
        <h1 className="text-4xl font-bold tracking-tight mb-4">
          Moss CMS
        </h1>
        <p className="text-xl text-muted-foreground">
          一个现代化的内容管理系统
        </p>
      </section>

      {/* Categories */}
      <section className="mb-12">
        <h2 className="text-2xl font-semibold mb-4">分类</h2>
        <div className="flex flex-wrap gap-2">
          {categories.map((category) => (
            <Link
              key={category.id}
              href={`/category/${category.slug}`}
              className="inline-flex items-center rounded-full border px-4 py-1.5 text-sm font-medium transition-colors hover:bg-accent"
            >
              {category.name}
              <span className="ml-2 text-muted-foreground">
                ({category.articleCount})
              </span>
            </Link>
          ))}
        </div>
      </section>

      {/* Articles Grid */}
      <section>
        <h2 className="text-2xl font-semibold mb-6">最新文章</h2>
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
                <h3 className="font-semibold text-lg mb-2 group-hover:text-primary">
                  {article.title}
                </h3>
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

        {total > 20 && (
          <div className="mt-8 text-center">
            <Link
              href="/articles"
              className="inline-flex items-center justify-center rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90"
            >
              查看更多文章
            </Link>
          </div>
        )}
      </section>
    </div>
  );
}