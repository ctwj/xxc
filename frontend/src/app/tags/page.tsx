import { api } from "@/lib/api";
import Link from "next/link";

export const revalidate = 60;

async function getTags() {
  try {
    return await api.getTags();
  } catch (error) {
    console.error("Failed to fetch tags:", error);
    return [];
  }
}

export default async function TagsPage() {
  const tags = await getTags();

  return (
    <div className="container py-8">
      <h1 className="text-3xl font-bold mb-8">标签</h1>

      <div className="flex flex-wrap gap-3">
        {tags.map((tag) => (
          <Link
            key={tag.id}
            href={`/tag/${tag.slug}`}
            className="group inline-flex items-center rounded-full border bg-card px-4 py-2 shadow-sm transition-colors hover:bg-accent"
          >
            <span className="font-medium group-hover:text-primary">
              {tag.name}
            </span>
            <span className="ml-2 text-sm text-muted-foreground">
              {tag.articleCount}
            </span>
          </Link>
        ))}
      </div>

      {tags.length === 0 && (
        <div className="text-center py-12 text-muted-foreground">
          暂无标签
        </div>
      )}
    </div>
  );
}