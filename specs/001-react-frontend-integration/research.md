# Research: Next.js Frontend Integration

**Date**: 2026-04-29
**Feature**: 001-react-frontend-integration

## Research Tasks

### 1. Next.js ISR 最佳实践

**Decision**: 使用 Next.js 15 App Router + generateStaticParams

**Rationale**:
- App Router 是 Next.js 推荐的现代架构
- `generateStaticParams` 支持动态路由的静态生成
- 内置 `revalidatePath` / `revalidateTag` 支持 ISR 按需更新
- 与 Vercel 深度集成，部署简单

**Alternatives Considered**:
- Pages Router: 较老的架构，不推荐新项目使用
- 纯 SSG: 无法支持按需更新
- SSR: 性能不如 ISR

**Implementation Notes**:
```typescript
// app/article/[slug]/page.tsx
export async function generateStaticParams() {
  const articles = await fetch('/api/articles').then(res => res.json())
  return articles.map(article => ({ slug: article.slug }))
}

export const revalidate = 3600 // 1小时后自动重新验证
```

### 2. JWT 认证在 Next.js 中的实现

**Decision**: 使用 `jose` 库 + HttpOnly Cookie

**Rationale**:
- `jose` 是纯 JavaScript JWT 库，支持 Edge Runtime
- HttpOnly Cookie 防止 XSS 攻击
- Next.js API Routes 可以安全处理 Cookie

**Alternatives Considered**:
- NextAuth.js: 功能强大但复杂，对于简单 JWT 认证过度设计
- localStorage: 不安全，容易受 XSS 攻击
- Session: 需要后端存储，不适合 Headless 架构

**Implementation Notes**:
```typescript
// app/api/auth/login/route.ts
import { serialize } from 'cookie'
import { SignJWT } from 'jose'

export async function POST(request: Request) {
  const { email, password } = await request.json()
  // 验证用户...
  const token = await new SignJWT({ userId: user.id })
    .setProtectedHeader({ alg: 'HS256' })
    .setExpirationTime('7d')
    .sign(secret)

  return new Response(JSON.stringify({ success: true }), {
    headers: {
      'Set-Cookie': serialize('token', token, {
        httpOnly: true,
        secure: process.env.NODE_ENV === 'production',
        sameSite: 'lax',
        path: '/',
        maxAge: 60 * 60 * 24 * 7
      })
    }
  })
}
```

### 3. Webhook 触发 ISR 更新

**Decision**: 使用 Next.js API Route + Vercel Webhook

**Rationale**:
- Next.js 内置 `revalidatePath` / `revalidateTag` API
- Vercel 支持 Webhook 触发重新部署
- 简单可靠，无需额外基础设施

**Alternatives Considered**:
- 定时轮询: 有延迟，浪费资源
- On-Demand ISR SDK: 需要 Vercel Pro 计划
- 手动触发: 运维成本高

**Implementation Notes**:
```typescript
// app/api/revalidate/route.ts
import { revalidatePath, revalidateTag } from 'next/cache'

export async function POST(request: Request) {
  const { secret, slug } = await request.json()

  // 验证 Webhook secret
  if (secret !== process.env.REVALIDATE_SECRET) {
    return Response.json({ error: 'Invalid secret' }, { status: 401 })
  }

  if (slug) {
    revalidatePath(`/article/${slug}`)
  } else {
    revalidatePath('/', 'layout') // 更新首页
  }

  return Response.json({ revalidated: true })
}
```

### 4. Go 后端 JWT 实现

**Decision**: 使用 `golang-jwt/jwt/v5` 库

**Rationale**:
- 最流行的 Go JWT 库
- 支持 HS256, RS256 等算法
- 与 Fiber 中间件集成简单

**Alternatives Considered**:
- 自实现 JWT: 安全风险高
- 其他库 (jwt-go): 已停止维护

**Implementation Notes**:
```go
// middleware/jwt.go
import "github.com/golang-jwt/jwt/v5"

func JWTMiddleware(secret string) fiber.Handler {
    return func(c *fiber.Ctx) error {
        tokenString := c.Cookies("token")
        if tokenString == "" {
            return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
        }

        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return []byte(secret), nil
        })

        if err != nil || !token.Valid {
            return c.Status(401).JSON(fiber.Map{"error": "invalid token"})
        }

        claims := token.Claims.(jwt.MapClaims)
        c.Locals("userId", claims["userId"])

        return c.Next()
    }
}
```

### 5. CORS 配置

**Decision**: 使用 Fiber CORS 中间件

**Rationale**:
- Fiber 内置 CORS 支持
- 配置简单，支持 credentials

**Implementation Notes**:
```go
// main.go
import "github.com/gofiber/fiber/v2/middleware/cors"

app.Use(cors.New(cors.Config{
    AllowOrigins: "https://your-domain.vercel.app",
    AllowHeaders: "Origin, Content-Type, Accept, Authorization",
    AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
    AllowCredentials: true,
}))
```

### 6. xxc.zip 组件迁移策略

**Decision**: 渐进式迁移，适配 Next.js App Router

**Rationale**:
- 保持原有 UI 设计和交互
- 最小化代码重写
- 利用 Next.js 优化性能

**迁移清单**:
| 原组件 | Next.js 目标 | 改动 |
|--------|--------------|------|
| `HomePage.tsx` | `app/page.tsx` | 移除 React Router，使用 Next.js Link |
| `LoginPage.tsx` | `app/login/page.tsx` | API 调用改为 Moss CMS |
| `RegisterPage.tsx` | `app/register/page.tsx` | API 调用改为 Moss CMS |
| `FavoritePage.tsx` | `app/favorites/page.tsx` | API 调用改为 Moss CMS |
| `InfoCard.tsx` | `components/home/InfoCard.tsx` | 保持不变 |
| `CardStack.tsx` | `components/home/CardStack.tsx` | 保持不变 |
| `AuthContext.tsx` | `contexts/AuthContext.tsx` | 替换 localStorage 为 Cookie |

## Open Questions (Resolved)

| 问题 | 决策 |
|------|------|
| Next.js 版本 | 15.x (最新稳定版) |
| 渲染模式 | ISR (Incremental Static Regeneration) |
| 认证方式 | JWT + HttpOnly Cookie |
| ISR 触发 | Webhook 按需触发 |
| API 风格 | REST API |
| 部署平台 | Vercel |

## Dependencies

### 前端 (Next.js)
```json
{
  "dependencies": {
    "next": "^15.0.0",
    "react": "^18.0.0",
    "react-dom": "^18.0.0",
    "jose": "^5.0.0",
    "framer-motion": "^12.0.0",
    "@radix-ui/react-*": "latest",
    "tailwindcss": "^3.4.0",
    "clsx": "^2.1.0",
    "sonner": "^2.0.0"
  },
  "devDependencies": {
    "typescript": "^5.9.0",
    "@types/react": "^19.0.0",
    "@types/node": "^20.0.0",
    "jest": "^29.0.0",
    "@testing-library/react": "^14.0.0",
    "playwright": "^1.40.0"
  }
}
```

### 后端 (Go)
```go
// 新增依赖
import (
    "github.com/golang-jwt/jwt/v5"
    "github.com/gofiber/fiber/v2/middleware/cors"
)
```
