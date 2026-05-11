# Feature Specification: Fix Content Display Issues

**Feature Branch**: `010-fix-content-display`
**Created**: 2026-05-11
**Status**: Draft
**Input**: User description: "系统部署后，还是发现问题，1. 正文存在 html标签， 2. 图片视频不显示， 访问 https://www.l9.lc 可以查看实际情况"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - View Article Cards Without HTML Tags (Priority: P1)

用户访问首页，浏览文章卡片列表时，不应看到任何原始 HTML 标签（如 `<p>`, `<br>`, `<strong>`）。卡片内容应显示干净的纯文本摘要。

**Why this priority**: 首页是用户的第一印象，所有 23 篇文章都受此影响，每篇卡片的描述区域都显示了 HTML 标签，严重影响可读性和专业形象。

**Independent Test**: 打开首页 https://www.l9.lc，检查所有文章卡片内容，确认无任何 HTML 标签以文本形式出现。

**Acceptance Scenarios**:

1. **Given** 用户访问首页，**When** 页面加载完成显示文章卡片列表，**Then** 所有卡片的描述文本中不包含 `<p>`, `</p>`, `<br>`, `<strong>`, `<a>` 等任何 HTML 标签
2. **Given** 文章描述包含 HTML 格式内容，**When** 卡片渲染该描述，**Then** 描述以纯文本形式展示，文字格式自然流畅
3. **Given** 文章描述包含超链接文本，**When** 卡片渲染该描述，**Then** 超链接文本以普通文字形式展示（卡片级别不需要可点击链接）

---

### User Story 2 - View Article Images (Priority: P1)

用户浏览文章卡片和文章详情页时，图片应正确加载和显示，不出现损坏图片图标。

**Why this priority**: 约 15/23 篇文章包含图片，所有图片当前返回 404，用户完全看不到任何图片内容。

**Independent Test**: 打开首页或任意含图片的文章详情页，确认图片正常加载显示。

**Acceptance Scenarios**:

1. **Given** 文章包含缩略图，**When** 用户查看首页文章卡片，**Then** 缩略图正确显示在卡片上
2. **Given** 文章详情包含内嵌图片，**When** 用户打开文章详情页，**Then** 所有图片正确加载和显示
3. **Given** 图片 URL 为相对路径，**When** 前端拼接完整 URL，**Then** 图片可正常通过该 URL 访问

---

### User Story 3 - View Article Videos (Priority: P2)

用户查看包含视频的文章详情页时，视频应正确加载和播放。

**Why this priority**: 仅 1 篇文章包含视频，影响范围小于图片问题，但与图片问题根因相同。

**Independent Test**: 打开包含视频的文章详情页，确认视频播放器出现且可正常播放。

**Acceptance Scenarios**:

1. **Given** 文章包含视频，**When** 用户打开文章详情页，**Then** 视频播放器正确显示且视频可正常播放
2. **Given** 视频 URL 为相对路径，**When** 前端拼接完整 URL，**Then** 视频可正常通过该 URL 访问

---

### Edge Cases

- 文章描述恰好被截断在 HTML 标签中间（如 `<p>Some text...</p>` 被截断为 `<p>Some text...`），系统应优雅处理不完整标签
- 文章无图片/视频时，不应显示损坏的媒体占位符
- 媒体 ID 非常长（如 19 位数字）时，URL 路径应正确拼接
- 网络慢或媒体加载失败时，应有合理的加载状态或错误提示

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: 首页文章卡片的描述内容必须以纯文本形式展示，不得出现任何 HTML 标签
- **FR-002**: 后端 API 返回的 `description` 字段应不包含 HTML 标签，或在存储时已将 HTML 转为纯文本
- **FR-003**: 媒体（图片/视频）的访问 URL 必须是公开可访问的，不需要管理员认证
- **FR-004**: Telegram 同步生成的媒体 URL 路径必须与实际注册的路由路径一致
- **FR-005**: 前端拼接媒体完整 URL 时，使用的基础 URL 必须能正确解析到媒体资源
- **FR-006**: 文章卡片组件应正确处理 `html` 和 `markdown` 类型的内容，选择合适的渲染方式
- **FR-007**: 文章详情页的 HTML 内容应正常渲染为富文本（非纯文本）

### Key Entities

- **Article**: 包含 `description`（摘要）、`content`（正文）、`thumbnail`（缩略图）、`mediaUrls`（媒体列表）、`videoUrl`（视频地址）、`contentType`（内容类型）
- **Media URL**: 由 Telegram 同步生成的相对路径，格式为 `/api/telegram/media/{id}`，需要映射到实际可访问的路由

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 首页所有 23 篇文章卡片的描述区域均无可见 HTML 标签
- **SC-002**: 所有包含缩略图的文章卡片图片加载成功率 100%
- **SC-003**: 文章详情页中的内嵌图片和视频可正常加载和显示
- **SC-004**: 媒体资源通过公开 URL 可直接访问（无需认证），HTTP 状态码为 200
- **SC-005**: 修复不影响管理后台的现有功能（管理后台媒体管理正常工作）

## Assumptions

- 媒体资源不需要认证即可公开访问（Telegram 频道内容本身是公开的）
- 文章描述（description）始终以纯文本形式展示即可，卡片级别不需要富文本渲染
- 修复应保持与现有管理后台的兼容性
- 前端使用 Next.js + React，媒体 URL 拼接逻辑使用 `API_BASE_URL` 环境变量
- 现有数据库中已存储的媒体 URL 使用 `/api/telegram/media/{id}` 格式，修复后这些已有数据应继续有效
