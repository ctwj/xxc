package plugins

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"moss/domain/core/entity"
	"moss/domain/core/repository/context"
	"moss/domain/core/service"
	"moss/domain/core/vo"
	pluginEntity "moss/domain/support/entity"
	"moss/infrastructure/utils/request"
	"strings"

	"go.uber.org/zap"
)

type AISeoPlugin struct {
	ApiURL string `json:"api_url"` // API 地址（兼容 OpenAI 格式）
	ApiKey string `json:"api_key"` // API 密钥
	Model  string `json:"model"`   // 模型名称（如 gpt-4）
	Enable bool   `json:"enable"`  // 是否启用

	// 关键词配置
	MinKeywords int `json:"min_keywords"`  // 最少关键词数（默认 4）
	MaxKeywords int `json:"max_keywords"`  // 最多关键词数（默认 8）
	MinLongTail int `json:"min_long_tail"` // 最少长尾关键词数（默认 2）

	// 标签配置
	MinTags int `json:"min_tags"` // 最少标签数（默认 1）
	MaxTags int `json:"max_tags"` // 最多标签数（默认 3）

	// 内容改写
	EnableRewrite bool `json:"enable_rewrite"` // 是否启用内容改写

	// 特殊分类名称
	OtherCategoryName string `json:"other_category_name"` // 其他分类名称

	// 定时任务
	CronEnable bool   `json:"cron_enable"` // 启用定时任务
	CronExp    string `json:"cron_exp"`    // 定时表达式
	BatchSize  int    `json:"batch_size"`  // 每次处理的文章数量（默认 10）

	// 强制重新生成（手动触发时可选）
	ForceRegenerate bool `json:"force_regenerate"` // 是否强制重新生成（忽略已有标记）

	ctx *pluginEntity.Plugin
}

func NewAISeoPlugin() *AISeoPlugin {
	return &AISeoPlugin{
		Model:             "gpt-3.5-turbo",
		Enable:            true,
		MinKeywords:       4,
		MaxKeywords:       8,
		MinLongTail:       2,
		MinTags:           1,
		MaxTags:           3,
		EnableRewrite:     true,
		OtherCategoryName: "其他软件",
		CronEnable:        false,
		CronExp:           "@every 30m",
		BatchSize:         10,
		ForceRegenerate:   false,
	}
}

// Info 返回插件信息
func (p *AISeoPlugin) Info() *pluginEntity.PluginInfo {
	return &pluginEntity.PluginInfo{
		ID:         "AISeoPlugin",
		About:      "SEO插件：AI 驱动的 SEO 优化，智能分类推荐、关键词提取、描述优化、内容改写、标签生成",
		RunEnable:  true,
		CronEnable: true,
		NoOptions:  false,
		PluginInfoPersistent: pluginEntity.PluginInfoPersistent{
			CronStart: false,
			CronExp:   p.CronExp,
		},
	}
}

// Load 插件加载
func (p *AISeoPlugin) Load(ctx *pluginEntity.Plugin) error {
	p.ctx = ctx
	// 注册文章事件
	service.Article.AddCreateBeforeEvents(p)
	service.Article.AddUpdateAfterEvents(p)
	return nil
}

// Run 插件执行（定时任务或手动触发）
func (p *AISeoPlugin) Run(ctx *pluginEntity.Plugin) error {
	p.ctx = ctx

	if !p.Enable {
		p.ctx.Log.Warn("AISeoPlugin is disabled")
		return nil
	}

	// 检查配置是否完整
	if p.ApiURL == "" || p.ApiKey == "" {
		p.ctx.Log.Error("AISeoPlugin configuration is incomplete: API URL or API Key is not configured")
		return errors.New("API URL or API Key is not configured")
	}

	p.ctx.Log.Info("AISeoPlugin started",
		zap.String("api_url", p.ApiURL),
		zap.String("model", p.Model),
		zap.Bool("force_regenerate", p.ForceRegenerate))

	// 获取未生成的文章
	articles, err := p.getUngeneratedArticles(p.BatchSize)
	if err != nil {
		p.ctx.Log.Error("Failed to get articles", zap.Error(err))
		return err
	}

	if len(articles) == 0 {
		p.ctx.Log.Info("No articles to process")
		return nil
	}

	p.ctx.Log.Info("AISeoPlugin started", zap.Int("articles", len(articles)))

	// 处理每篇文章
	successCount := 0
	for _, article := range articles {
		if err := p.processArticle(article); err != nil {
			p.ctx.Log.Error("Failed to process article",
				zap.Int("article_id", article.ID),
				zap.String("title", article.Title),
				zap.Error(err),
			)
		} else {
			successCount++
		}
	}

	p.ctx.Log.Info("AISeoPlugin completed", zap.Int("success", successCount))

	return nil
}

// ArticleCreateBefore 文章创建前事件
func (p *AISeoPlugin) ArticleCreateBefore(item *entity.Article) error {
	if !p.Enable {
		return nil
	}
	// 不在创建时自动处理，避免影响创建速度
	// 可以通过定时任务或手动触发来处理
	return nil
}

// ArticleUpdateAfter 文章更新后事件
func (p *AISeoPlugin) ArticleUpdateAfter(item *entity.Article) {
	if !p.Enable {
		return
	}
	// 不在更新时自动处理，避免频繁触发
	// 可以通过定时任务或手动触发来处理
}

// isArticleGenerated 检查文章是否已生成
func (p *AISeoPlugin) isArticleGenerated(article *entity.Article) bool {
	if article.Extends == nil {
		return false
	}
	val := article.Extends.Get("ai_seo_generated")
	if val == nil {
		return false
	}
	// 检查值是否为 true
	if boolVal, ok := val.(bool); ok {
		return boolVal
	}
	return false
}

// getUngeneratedArticles 获取未生成的文章
func (p *AISeoPlugin) getUngeneratedArticles(limit int) ([]*entity.Article, error) {
	// 获取文章基础列表
	baseList, err := service.Article.List(context.NewContext(limit*3, "id asc"))
	if err != nil {
		p.ctx.Log.Error("Failed to get article list", zap.Error(err))
		return nil, err
	}

	// 逐个获取完整文章信息
	var result []*entity.Article
	for _, base := range baseList {
		if len(result) >= limit {
			break
		}

		// 获取完整文章（包含详情）
		article, err := service.Article.Get(base.ID)
		if err != nil {
			p.ctx.Log.Warn("Failed to get article detail",
				zap.Int("article_id", base.ID),
				zap.Error(err))
			continue
		}

		// 检查是否需要处理
		needProcess := false
		if p.ForceRegenerate {
			// 强制重新生成模式：处理所有文章
			needProcess = true
		} else {
			// 正常模式：只处理未生成的文章
			if !p.isArticleGenerated(article) {
				needProcess = true
			}
		}

		if needProcess {
			result = append(result, article)
		}
	}

	return result, nil
}

// processArticle 处理单篇文章
func (p *AISeoPlugin) processArticle(article *entity.Article) error {
	// 检查是否已生成
	if !p.ForceRegenerate && p.isArticleGenerated(article) {
		p.ctx.Log.Debug("Article already generated, skipping",
			zap.Int("article_id", article.ID),
			zap.String("title", article.Title))
		return nil
	}

	p.ctx.Log.Info("Processing article",
		zap.Int("article_id", article.ID),
		zap.String("title", article.Title),
		zap.String("current_keywords", article.Keywords),
		zap.String("current_description", article.Description))

	// 调用 AI 生成 SEO 内容
	result, err := p.generateSEOContent(article)
	if err != nil {
		p.ctx.Log.Error("Failed to generate SEO content",
			zap.Int("article_id", article.ID),
			zap.String("title", article.Title),
			zap.Error(err))
		return err
	}

	// 更新文章字段
	p.updateArticle(article, result)

	// 保存文章
	if err := service.Article.Update(article); err != nil {
		p.ctx.Log.Error("Failed to update article",
			zap.Int("article_id", article.ID),
			zap.Error(err))
		return err
	}

	// 标记为已生成（只有成功时才记录）
	p.markAsGenerated(article)

	// 再次保存（更新 extends）
	if err := service.Article.Update(article); err != nil {
		p.ctx.Log.Error("Failed to update article extends",
			zap.Int("article_id", article.ID),
			zap.Error(err))
		return err
	}

	// 输出成功日志
	p.ctx.Log.Info("Article processed successfully",
		zap.Int("article_id", article.ID),
		zap.String("title", article.Title),
		zap.Int("keywords_count", len(result.Keywords)),
		zap.String("new_keywords", strings.Join(result.Keywords, ",")),
		zap.String("new_description", result.Description),
		zap.Int("tags_count", len(result.Tags)),
		zap.Bool("content_rewrited", result.ContentRewrited),
	)

	return nil
}

// AISeoResult SEO 生成结果
type AISeoResult struct {
	CategoryID      int      `json:"category_id"`
	CategoryChanged bool     `json:"category_changed"`
	Keywords        []string `json:"keywords"`
	Description     string   `json:"description"`
	ContentRewrited bool     `json:"content_rewrited"`
	Tags            []string `json:"tags"`
}

// generateSEOContent 生成 SEO 内容
func (p *AISeoPlugin) generateSEOContent(article *entity.Article) (*AISeoResult, error) {
	result := &AISeoResult{}

	// 1. 智能分类推荐
	if article.CategoryID == 0 || p.isOtherCategory(article.CategoryID) {
		categoryID, err := p.recommendCategory(article)
		if err != nil {
			p.ctx.Log.Warn("Failed to recommend category", zap.Error(err))
		} else if categoryID > 0 {
			result.CategoryID = categoryID
			result.CategoryChanged = true
			p.ctx.Log.Info("Category recommended",
				zap.Int("article_id", article.ID),
				zap.Int("category_id", categoryID))
		}
	}

	// 2. 关键词提取
	keywords, err := p.extractKeywords(article)
	if err != nil {
		p.ctx.Log.Warn("Failed to extract keywords", zap.Error(err))
	} else if len(keywords) > 0 {
		result.Keywords = keywords
		p.ctx.Log.Info("Keywords extracted",
			zap.Int("article_id", article.ID),
			zap.Int("count", len(keywords)),
			zap.Strings("keywords", keywords))
	} else {
		p.ctx.Log.Warn("No keywords extracted",
			zap.Int("article_id", article.ID))
	}

	// 3. 描述优化
	description, err := p.optimizeDescription(article, result.Keywords)
	if err != nil {
		p.ctx.Log.Warn("Failed to optimize description", zap.Error(err))
	} else if description != "" {
		result.Description = description
		p.ctx.Log.Info("Description optimized",
			zap.Int("article_id", article.ID),
			zap.String("description", description))
	}

	// 4. 内容改写
	if p.EnableRewrite {
		content, err := p.rewriteContent(article, result.Keywords)
		if err != nil {
			p.ctx.Log.Warn("Failed to rewrite content", zap.Error(err))
		} else if content != "" {
			article.Content = content
			result.ContentRewrited = true
			p.ctx.Log.Info("Content rewritten",
				zap.Int("article_id", article.ID))
		}
	}

	// 5. 标签生成
	tags, err := p.generateTags(article, result.Keywords)
	if err != nil {
		p.ctx.Log.Warn("Failed to generate tags", zap.Error(err))
	} else if len(tags) > 0 {
		result.Tags = tags
		p.ctx.Log.Info("Tags generated",
			zap.Int("article_id", article.ID),
			zap.Strings("tags", tags))
	}

	return result, nil
}

// isOtherCategory 检查是否为"其他"分类
func (p *AISeoPlugin) isOtherCategory(categoryID int) bool {
	if categoryID == 0 {
		return false
	}
	category, err := service.Category.Get(categoryID)
	if err != nil {
		return false
	}
	return category.Name == p.OtherCategoryName
}

// recommendCategory 推荐分类
func (p *AISeoPlugin) recommendCategory(article *entity.Article) (int, error) {
	// 获取所有分类
	categories, err := service.Category.List(context.NewContext(100, "id asc"))
	if err != nil {
		return 0, err
	}

	if len(categories) == 0 {
		return 0, errors.New("no categories available")
	}

	// 构建分类列表字符串
	categoryList := ""
	for _, cat := range categories {
		categoryList += fmt.Sprintf("- %s\n", cat.Name)
	}

	// 构建提示词
	prompt := fmt.Sprintf(`你是一个内容分类专家。任务：根据给定的文章标题和内容，从提供的可选分类列表中，选择一个最合适的分类。

文章标题：%s
文章内容：%s

可选分类列表（每个分类用英文逗号分隔）：%s

要求：
- 只返回一个分类名称，且必须严格从列表中选择。
- 如果没有任何分类适合，请返回 "NONE"。
- 不要返回任何其他内容（如解释、标点等）。`,
		article.Title,
		truncateText(article.Content, 1000),
		categoryList,
	)

	// 调用 AI API
	response, err := p.callAI(prompt)
	if err != nil {
		return 0, err
	}

	// 解析响应
	categoryName := strings.TrimSpace(response)
	if categoryName == "NONE" || categoryName == "" {
		return 0, nil
	}

	// 查找分类 ID
	for _, cat := range categories {
		if cat.Name == categoryName {
			return cat.ID, nil
		}
	}

	return 0, nil
}

// extractKeywords 提取关键词
func (p *AISeoPlugin) extractKeywords(article *entity.Article) ([]string, error) {
	p.ctx.Log.Debug("Starting keyword extraction",
		zap.Int("article_id", article.ID),
		zap.String("title", article.Title))

	// 构建提示词
	prompt := fmt.Sprintf(`你是一个 SEO 关键词提取专家。请从以下文章中提取与内容高度相关、符合 SEO 优化原则的关键词。

文章标题：%s
文章内容：%s

要求：
- 提取的关键词总数在 %d 到 %d 个之间。
- 其中至少 %d 个为长尾关键词（由 3-5 个词组成的短语）。
- 其余为短关键词（由 1-2 个词组成）。
- 关键词之间用英文逗号分隔。
- 只返回关键词列表，不要包含任何其他文字或解释。
- 如果无法满足数量要求，则尽可能接近，但不要强行填充不相关的关键词。`,
		article.Title,
		truncateText(article.Content, 1000),
		p.MinKeywords,
		p.MaxKeywords,
		p.MinLongTail,
	)

	// 调用 AI API
	p.ctx.Log.Debug("Calling AI API for keyword extraction",
		zap.Int("article_id", article.ID))
	response, err := p.callAI(prompt)
	if err != nil {
		p.ctx.Log.Error("AI API call failed for keyword extraction",
			zap.Int("article_id", article.ID),
			zap.Error(err))
		return nil, err
	}

	p.ctx.Log.Debug("AI API response received for keyword extraction",
		zap.Int("article_id", article.ID),
		zap.String("response", response))

	// 解析响应
	keywords := strings.Split(response, ",")
	for i, kw := range keywords {
		keywords[i] = strings.TrimSpace(kw)
	}

	// 过滤空字符串
	var result []string
	for _, kw := range keywords {
		if kw != "" {
			result = append(result, kw)
		}
	}

	p.ctx.Log.Debug("Keywords parsed",
		zap.Int("article_id", article.ID),
		zap.Int("count", len(result)),
		zap.Strings("keywords", result))

	return result, nil
}

// optimizeDescription 优化描述
func (p *AISeoPlugin) optimizeDescription(article *entity.Article, keywords []string) (string, error) {
	p.ctx.Log.Debug("Starting description optimization",
		zap.Int("article_id", article.ID))

	// 构建提示词
	keywordsStr := strings.Join(keywords, "、")
	prompt := fmt.Sprintf(`你是一个 SEO 描述生成专家。请根据以下文章标题、内容和关键词，生成一个优化的 SEO 描述。

文章标题：%s
文章内容：%s
关键词列表（用英文逗号分隔）：%s

要求：
- 描述长度不超过 250 个字符（包括标点）。
- 自然地融入核心关键词，突出文章主题。
- 描述应吸引用户点击，同时简洁明了地概括文章内容。
- 避免关键词堆砌，保持语句通顺。
- 只返回描述内容本身，不要包含任何其他文字。`,
		article.Title,
		truncateText(article.Content, 500),
		keywordsStr,
	)

	// 调用 AI API
	response, err := p.callAI(prompt)
	if err != nil {
		p.ctx.Log.Error("AI API call failed for description optimization",
			zap.Int("article_id", article.ID),
			zap.Error(err))
		return "", err
	}

	// 截断到 250 字符
	runes := []rune(response)
	if len(runes) > 250 {
		response = string(runes[:250])
	}

	p.ctx.Log.Debug("Description optimized",
		zap.Int("article_id", article.ID),
		zap.String("description", response))

	return response, nil
}

// rewriteContent 改写内容
func (p *AISeoPlugin) rewriteContent(article *entity.Article, keywords []string) (string, error) {
	p.ctx.Log.Debug("Starting content rewriting",
		zap.Int("article_id", article.ID),
		zap.String("title", article.Title))

	// 构建提示词
	keywordsStr := strings.Join(keywords, "、")
	prompt := fmt.Sprintf(`你是一个专业的 SEO 内容编辑和文案创作专家。请根据给定的关键词和文章主题，对以下文章进行深度改写、扩写和优化，以大幅提升搜索引擎优化效果。

文章标题：%s
文章内容（可能包含 HTML 标签）：%s
关键词列表（用英文逗号分隔）：%s

重要要求：
1. **内容扩写**：在保持原文核心信息的基础上，大幅扩充内容长度（至少增加30%），增加更多细节、使用场景、注意事项等有价值的信息。
2. **文案调整**：使用更具吸引力和专业性的语言，改善文章的可读性和专业性。
3. **差异化改写**：确保改写后的内容与原文至少有 30% 以上的文字差异，避免雷同。
4. **关键词优化**：在文中自然地融入所有提供的关键词，每个关键词至少出现 2-3 次，避免生硬堆砌。
5. **HTML 结构保留**：严格保留原有的 HTML 标签格式（如 <p>、<h3>、<img> 等），只修改标签内的文本内容。
6. **图片保留**：保留原文中的所有图片标签和属性不变。
7. **内容丰富性**：增加更多实用的信息，如使用技巧、注意事项、常见问题等，使内容更加丰富和有价值。
8. **语言风格**：使用更生动、专业、有吸引力的语言，提升文章质量和用户体验。
9. **结构优化**：如果原文结构较简单，可以适当增加段落和小标题，使内容层次更清晰。
10. **独特性**：确保改写后的内容具有独特性，避免与其他文章雷同。

只返回改写后的完整文章内容，不要添加任何其他文字说明。`,
		article.Title,
		truncateText(article.Content, 2000),
		keywordsStr,
	)

	// 调用 AI API
	p.ctx.Log.Debug("Calling AI API for content rewriting",
		zap.Int("article_id", article.ID))
	response, err := p.callAI(prompt)
	if err != nil {
		p.ctx.Log.Error("AI API call failed for content rewriting",
			zap.Int("article_id", article.ID),
			zap.Error(err))
		return "", err
	}

	p.ctx.Log.Debug("Content rewritten successfully",
		zap.Int("article_id", article.ID),
		zap.Int("original_length", len(article.Content)),
		zap.Int("new_length", len(response)))

	return response, nil
}

// generateTags 生成标签
func (p *AISeoPlugin) generateTags(article *entity.Article, keywords []string) ([]string, error) {
	// 构建提示词
	keywordsStr := strings.Join(keywords, "、")
	prompt := fmt.Sprintf(`你是一个内容标签生成专家。请根据以下文章标题、内容和关键词，生成一组相关的标签。

文章标题：%s
文章内容：%s
关键词列表（用英文逗号分隔）：%s

要求：
- 生成 %d 到 %d 个标签。
- 标签应简洁明了，通常是 1-3 个词组成，便于分类和搜索。
- 标签必须与文章内容高度相关，反映核心主题。
- 使用英文逗号分隔标签。
- 只返回标签列表，不要任何其他文字。`,
		article.Title,
		truncateText(article.Content, 1000),
		keywordsStr,
		p.MinTags,
		p.MaxTags,
	)

	// 调用 AI API
	response, err := p.callAI(prompt)
	if err != nil {
		return nil, err
	}

	// 解析响应
	tags := strings.Split(response, ",")
	for i, tag := range tags {
		tags[i] = strings.TrimSpace(tag)
	}

	// 过滤空字符串
	var result []string
	for _, tag := range tags {
		if tag != "" {
			result = append(result, tag)
		}
	}

	return result, nil
}

// callAI 调用 AI API
func (p *AISeoPlugin) callAI(prompt string) (string, error) {
	if p.ApiURL == "" || p.ApiKey == "" {
		return "", errors.New("API URL or API Key is not configured")
	}

	// 构建请求
	requestBody := map[string]interface{}{
		"model": p.Model,
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"temperature": 0.7,
	}

	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	// 调用 API
	respBody, err := request.New().
		AddHeader("Authorization", "Bearer "+p.ApiKey).
		AddHeader("Content-Type", "application/json").
		PostReturnBody(p.ApiURL+"/chat/completions", bytes.NewReader(bodyBytes))
	if err != nil {
		return "", err
	}

	// 解析响应
	var resp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Error *struct {
			Message string `json:"message"`
		} `json:"error"`
	}

	if err := json.Unmarshal(respBody, &resp); err != nil {
		return "", err
	}

	if resp.Error != nil {
		return "", errors.New(resp.Error.Message)
	}

	if len(resp.Choices) == 0 {
		return "", errors.New("no response from AI")
	}

	return strings.TrimSpace(resp.Choices[0].Message.Content), nil
}

// updateArticle 更新文章字段
func (p *AISeoPlugin) updateArticle(article *entity.Article, result *AISeoResult) {
	p.ctx.Log.Debug("Updating article fields",
		zap.Int("article_id", article.ID),
		zap.String("old_keywords", article.Keywords),
		zap.String("old_description", article.Description))

	// 更新分类
	if result.CategoryID > 0 {
		article.CategoryID = result.CategoryID
		p.ctx.Log.Debug("Category updated",
			zap.Int("article_id", article.ID),
			zap.Int("category_id", result.CategoryID))
	}

	// 更新关键词
	if len(result.Keywords) > 0 {
		oldKeywords := article.Keywords
		article.Keywords = strings.Join(result.Keywords, ",")
		p.ctx.Log.Info("Keywords updated",
			zap.Int("article_id", article.ID),
			zap.String("old_keywords", oldKeywords),
			zap.String("new_keywords", article.Keywords))
	}

	// 更新描述
	if result.Description != "" {
		oldDescription := article.Description
		article.Description = result.Description
		p.ctx.Log.Info("Description updated",
			zap.Int("article_id", article.ID),
			zap.String("old_description", oldDescription),
			zap.String("new_description", article.Description))
	}

	// 更新标签
	if len(result.Tags) > 0 {
		p.updateArticleTags(article, result.Tags)
	}
}

// updateArticleTags 更新文章标签
func (p *AISeoPlugin) updateArticleTags(article *entity.Article, tags []string) {
	// 删除旧标签关联
	oldTags, err := service.Tag.ListByArticleID(context.NewContext(100, ""), article.ID)
	if err == nil {
		for _, tag := range oldTags {
			service.Mapping.DeleteArticleTag(article.ID, tag.ID)
		}
	}

	// 创建新标签
	for _, tagName := range tags {
		// 检查标签是否已存在
		tagID, err := service.Tag.GetIdByNameOrCreate(tagName)
		if err != nil {
			p.ctx.Log.Warn("Failed to create tag", zap.String("tag", tagName), zap.Error(err))
			continue
		}

		// 创建关联
		if err := service.Mapping.CreateArticleTag(article.ID, tagID); err != nil {
			p.ctx.Log.Warn("Failed to create article-tag mapping",
				zap.Int("article_id", article.ID),
				zap.Int("tag_id", tagID),
				zap.Error(err),
			)
		}
	}
}

// markAsGenerated 标记为已生成
func (p *AISeoPlugin) markAsGenerated(article *entity.Article) {
	if article.Extends == nil {
		article.Extends = make(vo.Extends, 0)
	}

	// 移除旧的标记（如果存在）
	newExtends := make(vo.Extends, 0)
	for _, item := range article.Extends {
		if item.Key != "ai_seo_generated" {
			newExtends = append(newExtends, item)
		}
	}

	// 添加新标记（只记录 true）
	newExtends = append(newExtends, vo.ExtendsItem{
		Key:   "ai_seo_generated",
		Value: true,
	})

	article.Extends = newExtends
}

// truncateText 截断文本
func truncateText(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}
	runes := []rune(text)
	if len(runes) <= maxLen {
		return text
	}
	return string(runes[:maxLen])
}