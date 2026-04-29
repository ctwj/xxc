import { api } from "@/lib/api";
import { Article } from "@/types";

export const articleService = {
  list: async (page = 1, pageSize = 20, category?: string, tag?: string) => {
    return api.getArticles({ page, pageSize, category, tag });
  },

  getBySlug: async (slug: string) => {
    return api.getArticleBySlug(slug);
  },

  search: async (keyword: string, page = 1) => {
    return api.searchArticles(keyword, page);
  },
};