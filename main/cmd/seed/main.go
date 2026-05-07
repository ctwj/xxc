package main

import (
	"fmt"
	"log"
	"moss/domain/core/entity"
	"moss/domain/core/repository"
	"moss/domain/core/service"
	"time"

	"golang.org/x/crypto/bcrypt"
	"moss/infrastructure/persistent/db"
	_ "moss/startup"
)

func main() {
	// Wait for db to be ready
	time.Sleep(1 * time.Second)

	// Create test user
	createTestUser()

	// Create test categories
	createTestCategories()

	// Create test tags
	createTestTags()

	// Create test articles
	createTestArticles()

	fmt.Println("\n✅ Test data created successfully!")
}

func createTestUser() {
	fmt.Println("Creating test user...")

	// Check if user already exists
	exists, _ := repository.User.ExistsByUsername("testuser")
	if exists {
		fmt.Println("  - Test user already exists")
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("test123456"), bcrypt.DefaultCost)
	user := &entity.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: string(hashedPassword),
		Role:     "user",
	}

	if err := repository.User.Create(user); err != nil {
		log.Printf("  - Failed to create test user: %v", err)
		return
	}

	fmt.Println("  - Created test user: testuser / test123456")
}

func createTestCategories() {
	fmt.Println("Creating test categories...")

	categories := []struct {
		Name        string
		Slug        string
		Title       string
		Description string
	}{
		{"技术", "tech", "技术分享", "技术相关的文章"},
		{"生活", "life", "生活随笔", "生活相关的文章"},
		{"资讯", "news", "最新资讯", "最新资讯信息"},
	}

	for _, cat := range categories {
		// Check if exists
		_, err := service.Category.GetBySlug(cat.Slug)
		if err == nil {
			fmt.Printf("  - Category '%s' already exists\n", cat.Name)
			continue
		}

		category := &entity.Category{
			Name:        cat.Name,
			Slug:        cat.Slug,
			Title:       cat.Title,
			Description: cat.Description,
		}

		if err := service.Category.Create(category); err != nil {
			log.Printf("  - Failed to create category '%s': %v", cat.Name, err)
		} else {
			fmt.Printf("  - Created category: %s\n", cat.Name)
		}
	}
}

func createTestTags() {
	fmt.Println("Creating test tags...")

	tags := []struct {
		Name        string
		Slug        string
		Title       string
		Description string
	}{
		{"前端", "frontend", "前端开发", "前端开发相关"},
		{"后端", "backend", "后端开发", "后端开发相关"},
		{"Go", "go", "Go语言", "Go语言相关"},
		{"React", "react", "React框架", "React框架相关"},
		{"Next.js", "nextjs", "Next.js框架", "Next.js框架相关"},
	}

	for _, t := range tags {
		// Check if exists
		_, err := service.Tag.GetBySlug(t.Slug)
		if err == nil {
			fmt.Printf("  - Tag '%s' already exists\n", t.Name)
			continue
		}

		tag := &entity.Tag{
			Name:        t.Name,
			Slug:        t.Slug,
			Title:       t.Title,
			Description: t.Description,
		}

		if err := service.Tag.Create(tag); err != nil {
			log.Printf("  - Failed to create tag '%s': %v", t.Name, err)
		} else {
			fmt.Printf("  - Created tag: %s\n", t.Name)
		}
	}
}

func createTestArticles() {
	fmt.Println("Creating test articles...")

	// Get categories
	techCat, _ := service.Category.GetBySlug("tech")
	lifeCat, _ := service.Category.GetBySlug("life")
	newsCat, _ := service.Category.GetBySlug("news")

	// Get tags
	frontendTag, _ := service.Tag.GetBySlug("frontend")
	backendTag, _ := service.Tag.GetBySlug("backend")
	goTag, _ := service.Tag.GetBySlug("go")
	reactTag, _ := service.Tag.GetBySlug("react")
	nextjsTag, _ := service.Tag.GetBySlug("nextjs")

	articles := []struct {
		Title       string
		Slug        string
		Description string
		Content     string
		Thumbnail   string
		CategoryID  int
		TagIDs      []int
		ContentType string
		MediaUrls   string
		VideoUrl    string
	}{
		{
			Title:       "Next.js 15 新特性详解",
			Slug:        "nextjs-15-new-features",
			Description: "Next.js 15 带来了许多令人兴奋的新特性，包括改进的路由系统和更好的性能优化。",
			Content:     "<p>Next.js 15 带来了许多令人兴奋的新特性，包括改进的路由系统和更好的性能优化。本文将详细介绍这些新特性。</p><h2>1. 改进的路由系统</h2><p>新的路由系统更加灵活，支持更复杂的路由配置。</p><h2>2. 性能优化</h2><p>通过改进的编译器和运行时，Next.js 15 提供了更好的性能。</p>",
			Thumbnail:   "https://images.unsplash.com/photo-1555066931-4365d14bab8c?w=800",
			CategoryID:  techCat.ID,
			TagIDs:      []int{frontendTag.ID, reactTag.ID, nextjsTag.ID},
			ContentType: "text",
		},
		{
			Title:       "Go 语言并发编程最佳实践",
			Slug:        "go-concurrency-best-practices",
			Description: "Go 语言的并发编程是其最大的特色之一，本文介绍一些最佳实践。",
			Content:     "<p>Go 语言的并发编程是其最大的特色之一。通过 goroutine 和 channel，我们可以轻松实现高效的并发程序。</p><h2>Goroutine 使用技巧</h2><p>合理使用 goroutine 可以显著提高程序性能。</p><h2>Channel 最佳实践</h2><p>Channel 是 goroutine 之间通信的桥梁。</p>",
			Thumbnail:   "https://images.unsplash.com/photo-1526374965328-7f61d4dc18c5?w=800",
			CategoryID:  techCat.ID,
			TagIDs:      []int{backendTag.ID, goTag.ID},
			ContentType: "text",
		},
		{
			Title:       "今日科技新闻速递",
			Slug:        "tech-news-today",
			Description: "今日科技领域的重要新闻和动态。",
			Content:     "<p>今日科技领域发生了许多重要事件：</p><ul><li>AI 技术持续发展</li><li>云计算市场增长</li><li>开源社区活跃</li></ul>",
			Thumbnail:   "https://images.unsplash.com/photo-1504711434969-e33886168f5c?w=800",
			CategoryID:  newsCat.ID,
			TagIDs:      []int{},
			ContentType: "text",
		},
		{
			Title:       "生活小技巧分享",
			Slug:        "life-tips",
			Description: "一些实用的生活小技巧，让你的生活更加便捷。",
			Content:     "<p>生活中有许多小技巧可以让我们的日常更加便捷：</p><h2>时间管理</h2><p>合理规划时间，提高效率。</p><h2>健康生活</h2><p>保持良好的作息习惯。</p>",
			Thumbnail:   "https://images.unsplash.com/photo-1506784983877-45594efa4cbe?w=800",
			CategoryID:  lifeCat.ID,
			TagIDs:      []int{},
			ContentType: "text",
		},
		{
			Title:       "React 19 新特性预览",
			Slug:        "react-19-preview",
			Description: "React 19 即将发布，让我们提前了解新特性。",
			Content:     "<p>React 19 带来了许多期待已久的新特性：</p><h2>Server Components</h2><p>服务端组件将成为默认选项。</p><h2>改进的状态管理</h2><p>新的状态管理方案更加简洁。</p>",
			Thumbnail:   "https://images.unsplash.com/photo-1633356122544-f134324a6cee?w=800",
			CategoryID:  techCat.ID,
			TagIDs:      []int{frontendTag.ID, reactTag.ID},
			ContentType: "text",
		},
	}

	for _, a := range articles {
		// Check if exists
		_, err := service.Article.GetBySlug(a.Slug)
		if err == nil {
			fmt.Printf("  - Article '%s' already exists\n", a.Title)
			continue
		}

		now := time.Now().Unix()
		article := &entity.Article{
			ArticleBase: entity.ArticleBase{
				Title:       a.Title,
				Slug:        a.Slug,
				Description: a.Description,
				Thumbnail:   a.Thumbnail,
				CategoryID:  a.CategoryID,
				Status:      true,
				Views:       100,
				CreateTime:  now,
			},
		}

		if err := service.Article.Create(article); err != nil {
			log.Printf("  - Failed to create article '%s': %v", a.Title, err)
			continue
		}

		// Create article detail
		detail := &entity.ArticleDetail{
			ArticleID:   article.ID,
			Content:     a.Content,
			ContentType: a.ContentType,
			MediaUrls:   a.MediaUrls,
			VideoUrl:    a.VideoUrl,
		}
		db.DB.Create(detail)

		// Create tag mappings
		for _, tagID := range a.TagIDs {
			mapping := &entity.MappingTag{
				ArticleID: article.ID,
				TagID:     tagID,
			}
			db.DB.Create(mapping)
		}

		fmt.Printf("  - Created article: %s\n", a.Title)
	}
}
