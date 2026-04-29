import { Article, Category, Tag, ApiResponse } from "@/types";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:9008";

interface FetchOptions extends RequestInit {
  token?: string;
}

class ApiClient {
  private baseUrl: string;

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl;
  }

  private async request<T>(
    endpoint: string,
    options: FetchOptions = {}
  ): Promise<T> {
    const { token, ...fetchOptions } = options;

    const headers = new Headers({
      "Content-Type": "application/json",
    });

    if (options.headers) {
      Object.entries(options.headers).forEach(([key, value]) => {
        headers.set(key, value);
      });
    }

    if (token) {
      headers.set("Authorization", `Bearer ${token}`);
    }

    const response = await fetch(`${this.baseUrl}${endpoint}`, {
      ...fetchOptions,
      headers,
      credentials: "include",
    });

    if (!response.ok) {
      const error = await response.json().catch(() => ({
        message: "An error occurred",
      }));
      throw new Error(error.message || `HTTP error! status: ${response.status}`);
    }

    return response.json();
  }

  async get<T>(endpoint: string, options?: FetchOptions): Promise<T> {
    return this.request<T>(endpoint, { ...options, method: "GET" });
  }

  async post<T>(endpoint: string, data?: unknown, options?: FetchOptions): Promise<T> {
    return this.request<T>(endpoint, {
      ...options,
      method: "POST",
      body: data ? JSON.stringify(data) : undefined,
    });
  }

  async put<T>(endpoint: string, data?: unknown, options?: FetchOptions): Promise<T> {
    return this.request<T>(endpoint, {
      ...options,
      method: "PUT",
      body: data ? JSON.stringify(data) : undefined,
    });
  }

  async delete<T>(endpoint: string, options?: FetchOptions): Promise<T> {
    return this.request<T>(endpoint, { ...options, method: "DELETE" });
  }

  // Article methods
  async getArticles(params: {
    page?: number;
    pageSize?: number;
    category?: string;
    tag?: string;
  }): Promise<ApiResponse<Article[]>> {
    const searchParams = new URLSearchParams();
    if (params.page) searchParams.set("page", String(params.page));
    if (params.pageSize) searchParams.set("pageSize", String(params.pageSize));
    if (params.category) searchParams.set("category", params.category);
    if (params.tag) searchParams.set("tag", params.tag);
    return this.get(`/api/articles?${searchParams}`);
  }

  async getArticleBySlug(slug: string): Promise<Article> {
    return this.get(`/api/articles/${slug}`);
  }

  async searchArticles(keyword: string, page = 1): Promise<ApiResponse<Article[]>> {
    return this.get(`/api/search?keyword=${encodeURIComponent(keyword)}&page=${page}`);
  }

  // Category methods
  async getCategories(): Promise<Category[]> {
    return this.get("/api/categories");
  }

  // Tag methods
  async getTags(): Promise<Tag[]> {
    return this.get("/api/tags");
  }

  // Auth methods
  async login(username: string, password: string): Promise<{ success: boolean; user?: { username: string; role: string }; error?: string }> {
    return this.post("/api/auth/login", { username, password });
  }

  async logout(): Promise<{ success: boolean }> {
    return this.post("/api/auth/logout");
  }

  async getCurrentUser(): Promise<{ id: number; username: string; role: string }> {
    return this.get("/api/auth/me");
  }

  // Favorite methods
  async getFavorites(): Promise<ApiResponse<{ id: number; articleId: number; article: Article }[]>> {
    return this.get("/api/favorites");
  }

  async addFavorite(articleId: number): Promise<{ success: boolean }> {
    return this.post("/api/favorites", { articleId });
  }

  async removeFavorite(id: number): Promise<{ success: boolean }> {
    return this.delete(`/api/favorites/${id}`);
  }
}

export const api = new ApiClient(API_BASE_URL);

// Export named APIs for backward compatibility
export const articleApi = {
  list: (page = 1, pageSize = 20, category?: string, tag?: string) =>
    api.getArticles({ page, pageSize, category, tag }),
  get: (slug: string) => api.getArticleBySlug(slug),
  search: (keyword: string, page = 1) => api.searchArticles(keyword, page),
};

export const categoryApi = {
  list: () => api.getCategories(),
};

export const tagApi = {
  list: () => api.getTags(),
};

export const authApi = {
  login: (username: string, password: string) => api.login(username, password),
  logout: () => api.logout(),
  me: () => api.getCurrentUser(),
};

export const favoriteApi = {
  list: () => api.getFavorites(),
  add: (articleId: number) => api.addFavorite(articleId),
  remove: (id: number) => api.removeFavorite(id),
};