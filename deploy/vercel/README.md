# Vercel 前端部署指南

## 1. 在 Vercel 创建项目

1. 登录 [vercel.com](https://vercel.com)
2. 点击 "Add New" → "Project"
3. 选择 "Import Git Repository" 或跳过（我们用 API 部署）
4. 记下项目名称

## 2. 获取 Vercel 凭证

### 方式一：通过 Vercel CLI

```bash
# 安装 Vercel CLI
npm i -g vercel

# 登录
vercel login

# 获取组织 ID (team slug)
vercel teams list
# 输出类似: team_xxx  your-team-name

# 创建临时项目并链接获取 project ID
mkdir -p /tmp/vercel-setup && cd /tmp/vercel-setup
echo '{}' > package.json
vercel link
# 选择你的团队和项目，或创建新项目

# 查看项目 ID
cat .vercel/project.json
# 输出: {"projectId":"xxx_xxx","orgId":"team_xxx"}
```

### 方式二：通过 Vercel 网页

1. **VERCEL_ORG_ID (团队 ID)**:
   - 打开 [vercel.com/account/settings](https://vercel.com/account/settings)
   - 在 "General" 页面找到 "Team ID" 或 "User ID"
   - 个人账户格式: `user_xxx`
   - 团队账户格式: `team_xxx`

2. **VERCEL_PROJECT_ID (项目 ID)**:
   - 打开项目设置页面: `vercel.com/your-username/project-name/settings`
   - 在 "General" 页面找到 "Project ID"
   - 格式: `prj_xxx`

3. **VERCEL_TOKEN (API Token)**:
   - 打开 [vercel.com/account/tokens](https://vercel.com/account/tokens)
   - 点击 "Create Token"
   - 设置名称如 "GitHub Actions Deploy"
   - 选择有效期 (建议 "No Expiration")
   - 复制生成的 token

## 3. 配置 GitHub Secrets

在仓库 Settings → Secrets and variables → Actions 添加：

| Secret | 说明 | 格式示例 |
|--------|------|----------|
| `VERCEL_TOKEN` | Vercel API Token | `xxxxxxxxxxxxxxxxxxxx` |
| `VERCEL_ORG_ID` | 组织/团队 ID | `team_xxx` 或 `user_xxx` |
| `VERCEL_PROJECT_ID` | 项目 ID | `prj_xxx` |

**注意**: `.vercel` 目录被 `.gitignore` 忽略，所以这些值必须通过 GitHub Secrets 配置。

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
