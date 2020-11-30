package main

import (
	"flag"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"io/ioutil"
	"log"
	"net/http"
	"org.tubetrue01/domain-update/config"
	"time"
)

var cache = make(map[string]string, 1)

var command = &config.Config{}

func main() {
	flag.StringVar(&command.AccessKeyID, "k", "", "AccessKeyId")
	flag.StringVar(&command.AccessKeySecret, "s", "", "AccessKeySecret")
	flag.StringVar(&command.Domain, "d", "", "Domain")
	flag.Parse()

	if command.AccessKeyID == "" || command.AccessKeySecret == "" || command.Domain == "" {
		flag.Usage()
		return
	}

	task(command)
}

// obtainIpFromPool 从缓存中获取 ip 地址
func obtainIpFromPool() (ip string, isExists bool) {
	ip, isExists = cache["ip"]
	return
}

// updateIpPool 更新本地缓存
func updateIpPool(ip string) {
	cache["ip"] = ip
}

// obtainPubIp 获取公网 ip 地址
func obtainPubIp() string {
	resp, err := http.Get("https://myexternalip.com/raw")
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	content, _ := ioutil.ReadAll(resp.Body)
	return string(content)
}

// obtainDomain 获取域名信息
func obtainDomain(config *config.Config) *alidns.Record {
	client, err := alidns.NewClientWithAccessKey("cn-hangzhou", config.AccessKeyID, config.AccessKeySecret)

	request := alidns.CreateDescribeDomainRecordsRequest()
	request.Scheme = "https"

	request.DomainName = config.Domain

	response, err := client.DescribeDomainRecords(request)

	if err != nil {
		log.Println(err.Error())
	}

	record := response.DomainRecords.Record[0]

	return &alidns.Record{
		RR:       record.RR,
		RecordId: record.RecordId,
		Type:     record.Type,
		Value:    record.Value,
		TTL:      record.TTL,
	}
}

// updateDomain 更新远程域名解析
func updateDomain(subDomain *alidns.Record, config *config.Config) {
	client, err := alidns.NewClientWithAccessKey("cn-hangzhou", config.AccessKeyID, config.AccessKeySecret)

	request := alidns.CreateUpdateDomainRecordRequest()
	request.Scheme = "https"
	request.RecordId = subDomain.RecordId
	request.RR = subDomain.RR
	request.Type = subDomain.Type
	request.Value = subDomain.Value
	request.TTL = requests.NewInteger64(subDomain.TTL)

	_, err = client.UpdateDomainRecord(request)
	if err != nil {
		log.Print(err.Error())
	}
}

// update 更新域名解析以及本地缓存
func update(subDomain *alidns.Record, config *config.Config, externalIp string) {
	updateDomain(subDomain, config)
	updateIpPool(externalIp)
}

// task 执行定时任务，根据需要定期更新域名解析
func task(config *config.Config) {
	tick := time.Tick(time.Second * time.Duration(10))
	for {
		select {
		case <-tick:
			pubIp := obtainPubIp()

			if ip, exists := obtainIpFromPool(); exists {
				log.Printf("ip 已存在，当前值为：%s, 当前本机公网 ip 为：%s\n", ip, pubIp)
				if ip != pubIp {
					log.Println("ip 已经更改，准备更新域名解析...")
					subDomain := obtainDomain(config)
					subDomain.Value = pubIp
					update(subDomain, config, pubIp)
				}
			} else {
				record := obtainDomain(config)
				ipFromDomain := record.Value

				if pubIp != ipFromDomain {
					log.Println("公网 ip 与域名 ip 不符，即将更新域名 ip...")
					record.Value = pubIp
					update(record, config, pubIp)
				} else {
					updateIpPool(pubIp)
				}
			}
		}
	}
}
