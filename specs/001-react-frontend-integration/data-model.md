# Data Model: Next.js Frontend Integration

**Date**: 2026-04-29
**Feature**: 001-react-frontend-integration

## Entity Overview

本功能主要复用现有 Moss CMS 的数据模型，新增少量字段以支持前端需求。

## Core Entities

### Article (文章)

**来源**: 现有 `domain/core/entity/article.go`

**字段**:
| 字段 | 类型 | 说明 | 前端使用 |
|------|------|------|----------|
| ID | uint | 主键 | - |
| Slug | string | URL标识 | 文章链接 |
| Title | string | 标题 | 页面标题 |
| Content | string | 内容 | 文章正文 |
| Description | string | 摘要 | SEO description |
| Thumbnail | string | 缩略图 | 卡片图片 |
| Keywords | string | 关键词 | SEO keywords |
| CategoryID | uint | 分类ID | 分类关联 |
| Views | int | 浏览量 | 显示浏览数 |
| Status | bool | 发布状态 | 过滤未发布 |
| CreateTime | time.Time | 创建时间 | 显示时间 |
| Extends | []Extend | 扩展字段 | 侧边栏信息 |
| Res | []Res | 资源字段 | 下载链接 |

**新增字段**: 无（复用现有）

**前端类型映射**:
```typescript
interface Article {
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
  createTime: string
  extends: Extend[]
  res: Res[]
  category?: Category
  tags?: Tag[]
}
```

### Category (分类)

**来源**: 现有 `domain/core/entity/category.go`

**字段**:
| 字段 | 类型 | 说明 |
|------|------|------|
| ID | uint | 主键 |
| Slug | string | URL标识 |
| Name | string | 名称 |
| Title | string | SEO标题 |
| Description | string | 描述 |
| Keywords | string | 关键词 |
| ParentID | uint | 父分类ID |

**前端类型映射**:
```typescript
interface Category {
  id: number
  slug: string
  name: string
  title: string
  description: string
  keywords: string
  parentId: number
}
```

### Tag (标签)

**来源**: 现有 `domain/core/entity/tag.go`

**字段**:
| 字段 | 类型 | 说明 |
|------|------|------|
| ID | uint | 主键 |
| Slug | string | URL标识 |
| Name | string | 名称 |
| Title | string | SEO标题 |
| Description | string | 描述 |
| Keywords | string | 关键词 |

**前端类型映射**:
```typescript
interface Tag {
  id: number
  slug: string
  name: string
  title: string
  description: string
  keywords: string
}
```

### User (用户)

**来源**: 现有 `domain/config/entity/admin.go`

**字段**:
| 字段 | 类型 | 说明 | 前端使用 |
|------|------|------|----------|
| ID | uint | 主键 | 用户标识 |
| Username | string | 用户名 | 登录凭证 |
| Password | string | 密码哈希 | - |
| Email | string | 邮箱 | 可选登录凭证 |
| Role | string | 角色 | 权限判断 |

**前端类型映射**:
```typescript
interface User {
  id: number
  username: string
  email?: string
  role: string
}
```

### Favorite (收藏) - 新增

**说明**: 用户收藏文章的关系表

**字段**:
| 字段 | 类型 | 说明 |
|------|------|------|
| ID | uint | 主键 |
| UserID | uint | 用户ID |
| ArticleID | uint | 文章ID |
| CreateTime | time.Time | 收藏时间 |

**前端类型映射**:
```typescript
interface Favorite {
  id: number
  userId: number
  articleId: number
  createTime: string
  article?: Article // 关联文章详情
}
```

## API Response Types

### 文章列表响应
```typescript
interface ArticleListResponse {
  data: Article[]
  total: number
  page: number
  pageSize: number
  hasMore: boolean
}
```

### 搜索响应
```typescript
interface SearchResponse {
  data: Article[]
  keyword: string
  total: number
}
```

### 认证响应
```typescript
interface AuthResponse {
  success: boolean
  user?: User
  error?: string
}
```

## Validation Rules

### Article
- `slug`: 必填，唯一，URL安全字符
- `title`: 必填，最大200字符
- `content`: 必填
- `status`: 默认 false（未发布）

### User
- `username`: 必填，唯一，3-50字符
- `password`: 必填，最少6字符，存储时哈希

### Favorite
- `userID`: 必填，存在
- `articleID`: 必填，存在且已发布
- 唯一约束: (userID, articleID)

## State Transitions

### Article 发布状态
```
草稿 (status=false) ──发布──→ 已发布 (status=true)
已发布 (status=true) ──撤回──→ 草稿 (status=false)
```

### 用户认证状态
```
未登录 ──登录成功──→ 已登录
已登录 ──Token过期──→ 未登录
已登录 ──主动退出──→ 未登录
```

## Relationships

```
User ──<Favorite>── Article
Article ──<Category>── Category (多对一)
Article ──<Mapping>── Tag (多对多，通过 Mapping 表)
```

## Database Schema Changes

**新增表**: `favorites`

```sql
CREATE TABLE favorites (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    article_id INTEGER NOT NULL,
    create_time DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, article_id),
    FOREIGN KEY (user_id) REFERENCES admins(id),
    FOREIGN KEY (article_id) REFERENCES articles(id)
);
```

**现有表变更**: 无