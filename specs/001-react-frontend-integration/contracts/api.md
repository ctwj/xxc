# API Contracts: Moss CMS Headless API

**Date**: 2026-04-29
**Feature**: 001-react-frontend-integration
**Base URL**: `https://api.moss-cms.com` (生产环境)

## Overview

本文档定义 Moss CMS 作为 Headless CMS 对外提供的 REST API 接口规范。

## Authentication

### JWT Token

所有需要认证的接口需要在请求头中携带 Token：

```
Cookie: token=<jwt_token>
```

或在请求头中：

```
Authorization: Bearer <jwt_token>
```

## Public APIs (无需认证)

### GET /api/articles

获取文章列表

**Query Parameters**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 否 | 页码，默认 1 |
| pageSize | int | 否 | 每页数量，默认 20，最大 100 |
| category | string | 否 | 分类 slug 筛选 |
| tag | string | 否 | 标签 slug 筛选 |

**Response**:
```json
{
  "data": [
    {
      "id": 1,
      "slug": "article-slug",
      "title": "文章标题",
      "description": "文章摘要",
      "thumbnail": "https://example.com/thumb.jpg",
      "views": 100,
      "createTime": "2026-04-29T10:00:00Z",
      "category": {
        "id": 1,
        "slug": "category-slug",
        "name": "分类名称"
      },
      "tags": [
        { "id": 1, "slug": "tag-slug", "name": "标签名称" }
      ]
    }
  ],
  "total": 100,
  "page": 1,
  "pageSize": 20,
  "hasMore": true
}
```

**Status Codes**:
- 200: 成功
- 400: 参数错误

---

### GET /api/articles/:slug

获取文章详情

**Path Parameters**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| slug | string | 是 | 文章 slug |

**Response**:
```json
{
  "id": 1,
  "slug": "article-slug",
  "title": "文章标题",
  "content": "<p>文章内容 HTML</p>",
  "description": "文章摘要",
  "thumbnail": "https://example.com/thumb.jpg",
  "keywords": "关键词1,关键词2",
  "views": 100,
  "createTime": "2026-04-29T10:00:00Z",
  "category": {
    "id": 1,
    "slug": "category-slug",
    "name": "分类名称"
  },
  "tags": [
    { "id": 1, "slug": "tag-slug", "name": "标签名称" }
  ],
  "extends": [
    { "key": "language", "value": "Go" }
  ],
  "res": [
    { "key": "download_links", "value": [...] }
  ]
}
```

**Status Codes**:
- 200: 成功
- 404: 文章不存在

---

### GET /api/categories

获取分类列表

**Response**:
```json
[
  {
    "id": 1,
    "slug": "category-slug",
    "name": "分类名称",
    "description": "分类描述",
    "articleCount": 10
  }
]
```

---

### GET /api/tags

获取标签列表

**Response**:
```json
[
  {
    "id": 1,
    "slug": "tag-slug",
    "name": "标签名称",
    "articleCount": 5
  }
]
```

---

### GET /api/search

搜索文章

**Query Parameters**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| keyword | string | 是 | 搜索关键词 |
| page | int | 否 | 页码，默认 1 |

**Response**:
```json
{
  "data": [
    {
      "id": 1,
      "slug": "article-slug",
      "title": "文章标题",
      "description": "文章摘要",
      "thumbnail": "https://example.com/thumb.jpg",
      "views": 100,
      "createTime": "2026-04-29T10:00:00Z"
    }
  ],
  "keyword": "搜索词",
  "total": 10
}
```

---

## Auth APIs

### POST /api/auth/login

用户登录

**Request Body**:
```json
{
  "username": "user@example.com",
  "password": "password123"
}
```

**Response**:
```json
{
  "success": true,
  "user": {
    "id": 1,
    "username": "user",
    "email": "user@example.com",
    "role": "user"
  }
}
```

**Headers**:
```
Set-Cookie: token=<jwt_token>; HttpOnly; Secure; SameSite=Lax; Path=/; Max-Age=604800
```

**Status Codes**:
- 200: 登录成功
- 401: 用户名或密码错误

---

### POST /api/auth/register

用户注册

**Request Body**:
```json
{
  "username": "newuser",
  "email": "newuser@example.com",
  "password": "password123"
}
```

**Response**:
```json
{
  "success": true,
  "message": "注册成功"
}
```

**Status Codes**:
- 200: 注册成功
- 400: 参数错误或用户名已存在

---

### POST /api/auth/logout

用户登出

**Response**:
```json
{
  "success": true
}
```

**Headers**:
```
Set-Cookie: token=; Max-Age=0
```

---

### GET /api/auth/me

获取当前用户信息

**Response**:
```json
{
  "id": 1,
  "username": "user",
  "email": "user@example.com",
  "role": "user"
}
```

**Status Codes**:
- 200: 成功
- 401: 未登录

---

## Favorites APIs (需要认证)

### GET /api/favorites

获取用户收藏列表

**Response**:
```json
{
  "data": [
    {
      "id": 1,
      "articleId": 1,
      "createTime": "2026-04-29T10:00:00Z",
      "article": {
        "id": 1,
        "slug": "article-slug",
        "title": "文章标题",
        "thumbnail": "https://example.com/thumb.jpg"
      }
    }
  ],
  "total": 5
}
```

---

### POST /api/favorites

添加收藏

**Request Body**:
```json
{
  "articleId": 1
}
```

**Response**:
```json
{
  "success": true,
  "id": 1
}
```

**Status Codes**:
- 200: 收藏成功
- 400: 已收藏或文章不存在
- 401: 未登录

---

### DELETE /api/favorites/:id

取消收藏

**Path Parameters**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | int | 是 | 收藏记录 ID |

**Response**:
```json
{
  "success": true
}
```

---

## Webhook APIs

### POST /api/webhook/revalidate

触发 ISR 重新验证

**Request Body**:
```json
{
  "secret": "webhook-secret-key",
  "slug": "article-slug",
  "type": "article"
}
```

**Response**:
```json
{
  "revalidated": true,
  "path": "/article/article-slug"
}
```

**Status Codes**:
- 200: 成功
- 401: secret 无效

---

## Error Response Format

所有错误响应使用统一格式：

```json
{
  "error": "错误类型",
  "message": "错误详细信息",
  "code": "ERROR_CODE"
}
```

## CORS Configuration

后端需要配置以下 CORS 头：

```
Access-Control-Allow-Origin: https://your-domain.vercel.app
Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS
Access-Control-Allow-Headers: Origin, Content-Type, Accept, Authorization
Access-Control-Allow-Credentials: true
```

## Rate Limiting

- 公开 API: 100 请求/分钟/IP
- 认证 API: 10 请求/分钟/IP
- Webhook API: 无限制（通过 secret 验证）
