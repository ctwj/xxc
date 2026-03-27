package entity

// Email 邮件配置
type Email struct {
	Enable     bool   `json:"enable"`      // 是否启用邮件发送
	Host       string `json:"host"`        // SMTP服务器地址
	Port       int    `json:"port"`        // SMTP端口
	Username   string `json:"username"`    // 发件人邮箱
	Password   string `json:"password"`    // 邮箱密码或授权码
	FromName   string `json:"from_name"`   // 发件人名称
	UseTLS     bool   `json:"use_tls"`     // 是否使用TLS
	VerifyEmail bool  `json:"verify_email"` // 是否启用邮箱验证
}

func NewEmail() *Email {
	return &Email{
		Port:       465,
		UseTLS:     true,
		VerifyEmail: false,
	}
}

func (*Email) ConfigID() string {
	return "email"
}
