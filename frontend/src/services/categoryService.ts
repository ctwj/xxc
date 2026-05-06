import { api } from "@/lib/api";
import { Category } from "@/types";

export const categoryService = {
  list: async (): Promise<Category[]> => {
    return api.getCategories();
  },

  getBySlug: async (slug: string): Promise<Category | undefined> => {
    const categories = await api.getCategories();
    return categories.find((c) => c.slug === slug);
  },
};