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

type Mail struct{}

// DoNotify 发送 Email 信息
func (m *Mail) DoNotify(config *config.Config, content interface{}) {
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

// DoNotifyBefore 执行通知前的预处理操作
func (m *Mail) DoNotifyBefore(config *config.Config, content interface{}) {
	m.DoNotify(config, content)
}
