package notify

import (
	"fmt"
	"github.com/jordan-wright/email"
	"log"
	"net/smtp"
	"org.tubetrue01/domain-update/config"
	"strings"
)

const (
	smtpHost = "smtp.qq.com"
	smtpPort = "25"
)

// SendEmail 发送信息到指定的邮箱（目前只支持 qq 邮箱）
func SendEmail(config *config.Config, content string) {
	e := email.NewEmail()
	e.From = config.Email
	e.To = []string{config.Email}
	e.Subject = "IP 变更通知"
	e.Text = []byte(fmt.Sprintf("IP 地址已经发生变化，新的 IP 地址为：%s\n", content))

	addr := strings.Join([]string{smtpHost, smtpPort}, ":")

	err := e.Send(addr, smtp.PlainAuth("", config.Email, config.EmailAuthCode, smtpHost))
	if err != nil {
		log.Fatal(err)
	}
}
