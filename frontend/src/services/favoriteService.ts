import { api } from "@/lib/api";
import { Article } from "@/types";

export const favoriteService = {
  list: async (page = 1, pageSize = 20) => {
    return api.getFavorites();
  },

  add: async (articleId: number) => {
    return api.addFavorite(articleId);
  },

  remove: async (id: number) => {
    return api.removeFavorite(id);
  },

  checkFavorited: async (articleId: number): Promise<boolean> => {
    try {
      const { data } = await api.getFavorites();
      return data.some((item) => item.articleId === articleId);
    } catch {
      return false;
    }
  },
};