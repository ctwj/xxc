# Vercel 前端部署指南

## 1. 在 Vercel 创建项目

1. 登录 [vercel.com](https://vercel.com)
2. 点击 "Add New" → "Project"
3. 选择 "Import Git Repository" 或跳过（我们用 API 部署）
4. 记下项目名称

## 2. 获取 Vercel 凭证

```bash
# 安装 Vercel CLI
npm i -g vercel

# 登录
vercel login

# 获取组织 ID
vercel whoami

# 在任意目录创建临时项目并链接
mkdir -p /tmp/vercel-setup && cd /tmp/vercel-setup
echo '{}' > package.json
vercel link
# 选择你的团队和项目，或创建新项目

# 查看项目 ID
cat .vercel/project.json
```

## 3. 配置 GitHub Secrets

在仓库 Settings → Secrets and variables → Actions 添加：

| Secret | 说明 | 获取方式 |
|--------|------|----------|
| `VERCEL_TOKEN` | Vercel API Token | 在 [vercel.com/account/tokens](https://vercel.com/account/tokens) 创建 |
| `VERCEL_ORG_ID` | 组织 ID | 运行 `vercel whoami` |
| `VERCEL_PROJECT_ID` | 项目 ID | 查看 `.vercel/project.json` 中的 `projectId` |

## 4. 触发部署

- **自动触发**: 新 Release 发布时自动部署
- **手动触发**: 在 GitHub Actions 页面手动运行 `deploy-frontend.yml`

## 5. 验证部署

部署完成后访问 Vercel 分配的域名，检查：
- 管理后台是否正常加载
- API 请求是否正确转发到 `api.l9.lc`

## 文件说明

- `vercel.json` - Vercel 配置，API 代理和缓存规则
- `.github/workflows/deploy-frontend.yml` - 自动部署工作流
