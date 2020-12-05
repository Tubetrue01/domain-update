package config

import (
	"flag"
	"os"
)

type Config struct {
	AccessKeyID     string
	AccessKeySecret string
	Domain          string
	Email           string
	EmailAuthCode   string
}

var command = &Config{}

func init() {
	flag.StringVar(&command.AccessKeyID, "k", "", "AccessKeyId")
	flag.StringVar(&command.AccessKeySecret, "s", "", "AccessKeySecret")
	flag.StringVar(&command.Domain, "d", "", "Domain")
	flag.StringVar(&command.Email, "e", "", "QQ Email")
	flag.StringVar(&command.EmailAuthCode, "a", "", "authCode of Email")
	flag.Parse()

	if !(command.AccessKeyID != "" && command.AccessKeySecret != "" && command.Domain != "" ||
		command.Email != "" && command.EmailAuthCode != "") {

		flag.Usage()
		os.Exit(1)

	}
}

// ObtainCommand 返回初始化之后的 Config 对象
func ObtainCommand() *Config {
	return command
}
