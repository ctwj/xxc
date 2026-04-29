// Article types
export interface Article {
  id: number
  slug: string
  title: string
  content: string
  description: string
  thumbnail: string
  keywords: string
  categoryId: number
  views: number
  status: boolean
  createTime: number  // Unix timestamp
  extends: Extend[]
  res: Res[]
  category?: Category
  tags?: Tag[]
}

export interface Extend {
  key: string
  value: string
}

export interface Res {
  key: string
  value: unknown
}

// Category types
export interface Category {
  id: number
  slug: string
  name: string
  title: string
  description: string
  keywords: string
  parentId: number
  articleCount?: number
}

// Tag types
export interface Tag {
  id: number
  slug: string
  name: string
  title: string
  description: string
  keywords: string
  articleCount?: number
}

// User types
export interface User {
  id: number
  username: string
  email?: string
  role: string
}

// Favorite types
export interface Favorite {
  id: number
  userId: number
  articleId: number
  createTime: string
  article?: Article
}

// API Response types
export interface ApiResponse<T> {
  data: T
  total: number
  page?: number
  pageSize?: number
  hasMore?: boolean
}

export interface ArticleListResponse {
  data: Article[]
  total: number
  page: number
  pageSize: number
  hasMore: boolean
}

export interface SearchResponse {
  data: Article[]
  keyword: string
  total: number
}

export interface AuthResponse {
  success: boolean
  user?: User
  error?: string
  message?: string
}

export interface FavoriteListResponse {
  data: Favorite[]
  total: number
}

// Content types for InfoCard
export type ContentType =
  | 'text'
  | 'image'
  | 'video'
  | 'image-text'
  | 'video-text'
  | 'images'
  | 'images-text'
  | 'images-video-text'
  | 'long-text'

export interface InfoCardData {
  id: string
  type: ContentType
  title?: string
  content?: string
  mediaUrl?: string
  mediaUrls?: string[]
  videoUrl?: string
  coverUrl?: string
  tag: string
  publishTime: string
  isNew?: boolean
  isFavorited?: boolean
}

// Theme types
export type Theme = 'light' | 'dark'