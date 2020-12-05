package main

import (
	"flag"
	"log"
	"org.tubetrue01/domain-update/config"
	"org.tubetrue01/domain-update/notify"
	"org.tubetrue01/domain-update/util"
	"os"
	"time"
)

var command = &config.Config{}

func main() {
	//task(command)
	notify.SendEmail(command, "192.168.0.2.1")

}

// task 执行定时任务，根据需要定期更新域名解析
func task(config *config.Config) {
	tick := time.Tick(time.Second * time.Duration(10))
	for {
		select {
		case <-tick:
			pubIp := util.ObtainPubIp()

			if ip, exists := util.ObtainIpFromPool(); exists {
				log.Printf("ip 已存在，当前值为：%s, 当前本机公网 ip 为：%s\n", ip, pubIp)
				if ip != pubIp {
					log.Println("ip 已经更改，准备更新域名解析...")
					subDomain := notify.ObtainDomain(config)
					subDomain.Value = pubIp
					util.Update(subDomain, config, pubIp)
				}
			} else {
				record := notify.ObtainDomain(config)
				ipFromDomain := record.Value

				if pubIp != ipFromDomain {
					log.Println("公网 ip 与域名 ip 不符，即将更新域名 ip...")
					record.Value = pubIp
					util.Update(record, config, pubIp)
				} else {
					util.UpdateIpPool(pubIp)
				}
			}
		}
	}
}

func init() {
	flag.StringVar(&command.AccessKeyID, "k", "", "AccessKeyId")
	flag.StringVar(&command.AccessKeySecret, "s", "", "AccessKeySecret")
	flag.StringVar(&command.Domain, "d", "", "Domain")
	flag.StringVar(&command.Email, "e", "", "QQ Email")
	flag.StringVar(&command.EmailAuthCode, "a", "", "authCode of Email")
	flag.Parse()

	if command.AccessKeyID == "" || command.AccessKeySecret == "" || command.Domain == "" {
		flag.Usage()
		os.Exit(1)
	}

}
