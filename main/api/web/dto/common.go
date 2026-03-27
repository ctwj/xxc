package dto

type MessageResult struct {
	Success bool   `json:"success"`
	Data    any    `json:"data"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type Captcha struct {
	Base64 string `json:"base64"`
	ID     string `json:"id"`
}

func NewCaptcha(base64, id string) *Captcha {
	return &Captcha{Base64: base64, ID: id}
}

// UserStats 用户统计信息
type UserStats struct {
	Total int64 `json:"total"`
	Today int64 `json:"today"`
}
