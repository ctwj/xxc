import { api } from "@/lib/api";
import Link from "next/link";

export const revalidate = 60;

async function getCategories() {
  try {
    return await api.getCategories();
  } catch (error) {
    console.error("Failed to fetch categories:", error);
    return [];
  }
}

export default async function CategoriesPage() {
  const categories = await getCategories();

  return (
    <div className="container py-8">
      <h1 className="text-3xl font-bold mb-8">分类</h1>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        {categories.map((category) => (
          <Link
            key={category.id}
            href={`/category/${category.slug}`}
            className="group rounded-lg border bg-card p-6 shadow-sm transition-colors hover:bg-accent"
          >
            <h2 className="text-xl font-semibold mb-2 group-hover:text-primary">
              {category.title || category.name}
            </h2>
            {category.description && (
              <p className="text-sm text-muted-foreground line-clamp-2">
                {category.description}
              </p>
            )}
            <p className="mt-4 text-sm text-muted-foreground">
              {category.articleCount} 篇文章
            </p>
          </Link>
        ))}
      </div>

      {categories.length === 0 && (
        <div className="text-center py-12 text-muted-foreground">
          暂无分类
        </div>
      )}
    </div>
  );
}