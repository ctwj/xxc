# Research: 前端双平台部署 (Vercel + Cloudflare)

**Feature**: 009-dual-deploy-vercel-cloudflare
**Date**: 2026-05-10

## 研究任务

### 1. Cloudflare Pages 构建配置

**Decision**: 使用 Cloudflare Pages 原生 Git 集成，自动检测 Vite 项目。

**Rationale**: Cloudflare Pages 原生支持 Vite 项目，可以自动检测构建命令和输出目录。

**配置参数**:
- 构建命令: `pnpm build`
- 输出目录: `dist`
- Node.js 版本: 18.x

**Alternatives considered**:
- 使用 Cloudflare Workers: 过于复杂，静态站点不需要
- 使用 GitHub Actions 部署: 增加复杂度，原生 Git 集成更简单

### 2. 双平台配置文件兼容性

**Decision**: 两个平台使用独立的配置文件，互不干扰。

**Rationale**:
- Vercel 使用 `vercel.json`
- Cloudflare 使用 `wrangler.toml`
- 两个配置文件可以共存，不会冲突

**Alternatives considered**:
- 统一配置文件: 不可能，两个平台配置格式不同

### 3. 环境变量管理

**Decision**: 在两个平台的 Dashboard 中分别配置环境变量。

**Rationale**:
- 敏感信息不应存储在代码仓库中
- 两个平台都支持加密环境变量
- 可以通过 Dashboard 或 CLI 配置

**配置方式**:
- Vercel: Dashboard → Settings → Environment Variables
- Cloudflare: Dashboard → Pages → Settings → Environment Variables

### 4. 自定义域名配置

**Decision**: 两个平台使用不同的子域名。

**Rationale**:
- 避免域名冲突
- 可以通过 DNS 负载均衡分流
- 便于监控各平台状态

**域名方案**:
- Vercel: `admin.example.com` 或 `vercel.example.com`
- Cloudflare: `cf.example.com` 或 `cloudflare.example.com`

### 5. 构建缓存策略

**Decision**: 两个平台都启用构建缓存。

**Rationale**:
- 减少构建时间
- 节省构建额度
- 提高部署速度

**配置**:
- Vercel: 默认启用构建缓存
- Cloudflare: 默认启用构建缓存

## 结论

双平台部署方案清晰可行：
1. Cloudflare Pages 原生支持 Vite 项目
2. 配置文件独立，不会冲突
3. 环境变量通过 Dashboard 配置
4. 使用不同子域名避免冲突
