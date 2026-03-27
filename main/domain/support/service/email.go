package service

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/smtp"
	"strings"

	"moss/domain/config"
)

// SendMail 发送邮件
func SendMail(to, subject, body string) error {
	email := config.Config.Email
	if !email.Enable {
		return errors.New("email service is not enabled")
	}
	if email.Host == "" {
		return errors.New("smtp host is not configured")
	}
	if email.Username == "" {
		return errors.New("smtp username is not configured")
	}

	from := email.Username
	if email.FromName != "" {
		from = fmt.Sprintf("%s <%s>", email.FromName, email.Username)
	}

	msg := fmt.Sprintf("From: %s\r\n", from)
	msg += fmt.Sprintf("To: %s\r\n", to)
	msg += fmt.Sprintf("Subject: %s\r\n", subject)
	msg += "MIME-Version: 1.0\r\n"
	msg += "Content-Type: text/html; charset=UTF-8\r\n"
	msg += "\r\n"
	msg += body

	auth := smtp.PlainAuth("", email.Username, email.Password, email.Host)

	addr := fmt.Sprintf("%s:%d", email.Host, email.Port)

	if email.UseTLS {
		return sendMailWithTLS(addr, auth, email.Username, []string{to}, []byte(msg))
	}

	return smtp.SendMail(addr, auth, email.Username, []string{to}, []byte(msg))
}

// sendMailWithTLS 使用TLS发送邮件
func sendMailWithTLS(addr string, auth smtp.Auth, from string, to []string, msg []byte) error {
	host := strings.Split(addr, ":")[0]
	
	conf := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	conn, err := tls.Dial("tcp", addr, conf)
	if err != nil {
		return fmt.Errorf("dial error: %v", err)
	}

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return fmt.Errorf("create client error: %v", err)
	}
	defer client.Close()

	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("auth error: %v", err)
	}

	if err = client.Mail(from); err != nil {
		return fmt.Errorf("set sender error: %v", err)
	}

	for _, addr := range to {
		if err = client.Rcpt(addr); err != nil {
			return fmt.Errorf("set recipient error: %v", err)
		}
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("get data writer error: %v", err)
	}

	_, err = w.Write(msg)
	if err != nil {
		return fmt.Errorf("write message error: %v", err)
	}

	if err = w.Close(); err != nil {
		return fmt.Errorf("close writer error: %v", err)
	}

	return client.Quit()
}

// SendVerificationEmail 发送验证邮件
func SendVerificationEmail(to, verifyLink string) error {
	subject := "请验证您的邮箱 - Moss"
	body := fmt.Sprintf(`
		<html>
		<body style="font-family: Arial, sans-serif; padding: 20px;">
			<h2>欢迎注册 Moss</h2>
			<p>请点击下方链接验证您的邮箱地址：</p>
			<p><a href="%s" style="background-color: #007bff; color: white; padding: 10px 20px; text-decoration: none; border-radius: 5px;">验证邮箱</a></p>
			<p>如果按钮无法点击，请复制以下链接到浏览器打开：</p>
			<p>%s</p>
			<p>此链接有效期为24小时。</p>
			<p>如果您没有注册账号，请忽略此邮件。</p>
		</body>
		</html>
	`, verifyLink, verifyLink)
	
	return SendMail(to, subject, body)
}
