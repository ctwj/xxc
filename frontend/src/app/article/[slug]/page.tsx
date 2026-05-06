import { Article } from "@/types";
import { api } from "@/lib/api";
import { notFound } from "next/navigation";
import Link from "next/link";

// Generate static params for ISR
export async function generateStaticParams() {
  try {
    const { data: articles } = await api.getArticles({ page: 1, pageSize: 100 });
    return articles.map((article: Article) => ({
      slug: article.slug,
    }));
  } catch (error) {
    console.error("Failed to generate static params:", error);
    return [];
  }
}

async function getArticle(slug: string): Promise<Article | null> {
  try {
    return await api.getArticleBySlug(slug);
  } catch (error) {
    console.error("Failed to fetch article:", error);
    return null;
  }
}

export default async function ArticlePage({
  params,
}: {
  params: Promise<{ slug: string }>;
}) {
  const { slug } = await params;
  const article = await getArticle(slug);

  if (!article) {
    notFound();
  }

  return (
    <article className="container py-8 max-w-4xl">
      {/* Breadcrumb */}
      <nav className="mb-6 text-sm text-muted-foreground">
        <Link href="/" className="hover:text-foreground">
          首页
        </Link>
        <span className="mx-2">/</span>
        {article.category && (
          <>
            <Link
              href={`/category/${article.category.slug}`}
              className="hover:text-foreground"
            >
              {article.category.name}
            </Link>
            <span className="mx-2">/</span>
          </>
        )}
        <span className="text-foreground">{article.title}</span>
      </nav>

      {/* Article Header */}
      <header className="mb-8">
        <h1 className="text-3xl font-bold mb-4">{article.title}</h1>
        <div className="flex flex-wrap items-center gap-4 text-sm text-muted-foreground">
          <time>
            {new Date(article.createTime * 1000).toLocaleDateString("zh-CN")}
          </time>
          <span>{article.views} 阅读</span>
          {article.category && (
            <Link
              href={`/category/${article.category.slug}`}
              className="inline-flex items-center rounded-full border px-3 py-1 hover:bg-accent"
            >
              {article.category.name}
            </Link>
          )}
        </div>
        {article.tags && article.tags.length > 0 && (
          <div className="mt-4 flex flex-wrap gap-2">
            {article.tags.map((tag) => (
              <Link
                key={tag.id}
                href={`/tag/${tag.slug}`}
                className="inline-flex items-center rounded-full bg-secondary px-3 py-1 text-xs font-medium hover:bg-secondary/80"
              >
                {tag.name}
              </Link>
            ))}
          </div>
        )}
      </header>

      {/* Featured Image */}
      {article.thumbnail && (
        <div className="mb-8 aspect-video overflow-hidden rounded-lg">
          <img
            src={article.thumbnail}
            alt={article.title}
            className="w-full h-full object-cover"
          />
        </div>
      )}

      {/* Article Content */}
      <div
        className="prose prose-neutral dark:prose-invert max-w-none"
        dangerouslySetInnerHTML={{ __html: article.content || "" }}
      />

      {/* Keywords */}
      {article.keywords && (
        <div className="mt-8 pt-8 border-t">
          <h3 className="text-sm font-medium mb-2">关键词</h3>
          <p className="text-sm text-muted-foreground">{article.keywords}</p>
        </div>
      )}
    </article>
  );
}