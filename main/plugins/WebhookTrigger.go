package plugins

import (
	"bytes"
	"encoding/json"
	"net/http"
	"moss/domain/core/entity"
	"moss/infrastructure/support/log"
)

// WebhookTrigger triggers webhook calls when articles are published/updated
type WebhookTrigger struct {
	WebhookURL     string
	Secret         string
	EnableOnCreate bool
	EnableOnUpdate bool
}

// NewWebhookTrigger creates a new webhook trigger
func NewWebhookTrigger(webhookURL, secret string) *WebhookTrigger {
	return &WebhookTrigger{
		WebhookURL:     webhookURL,
		Secret:         secret,
		EnableOnCreate: true,
		EnableOnUpdate: true,
	}
}

// ArticleCreateAfter handles article creation event
func (p *WebhookTrigger) ArticleCreateAfter(article *entity.Article) {
	if !p.EnableOnCreate {
		return
	}
	p.triggerWebhook(article.Slug, "article")
}

// ArticleUpdateAfter handles article update event
func (p *WebhookTrigger) ArticleUpdateAfter(article *entity.Article) {
	if !p.EnableOnUpdate {
		return
	}
	p.triggerWebhook(article.Slug, "article")
}

func (p *WebhookTrigger) triggerWebhook(slug, contentType string) {
	if p.WebhookURL == "" {
		return
	}

	payload := map[string]string{
		"secret": p.Secret,
		"slug":   slug,
		"type":   contentType,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Error("Failed to marshal webhook payload", log.Err(err))
		return
	}

	req, err := http.NewRequest("POST", p.WebhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Error("Failed to create webhook request", log.Err(err))
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error("Failed to send webhook", log.Err(err))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		log.Warn("Webhook returned error status", log.Any("status", resp.StatusCode))
	} else {
		log.Info("Webhook triggered successfully", log.String("slug", slug))
	}
}