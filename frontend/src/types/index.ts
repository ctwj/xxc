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
  // 多媒体字段
  type?: ContentType
  mediaUrls?: string  // JSON string
  videoUrl?: string
  coverUrl?: string
  tag?: string  // First tag name for display
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
  createdAt: string
  article?: Article
}

// Like types
export interface Like {
  id: number
  userId: number
  articleId: number
  type: LikeType  // 1=like, 2=dislike
  createdAt: string
  article?: Article
}

export type LikeType = 0 | 1 | 2  // 0=none, 1=like, 2=dislike

// View History types
export interface ViewHistory {
  id: number
  userId: number
  articleId: number
  viewedAt: string
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

export interface LikeStatusResponse {
  type: LikeType
  likes: number
  dislikes: number
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