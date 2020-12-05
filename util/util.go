package util

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"io/ioutil"
	"net/http"
	"org.tubetrue01/domain-update/config"
	"org.tubetrue01/domain-update/notify"
)

var cache = make(map[string]string, 1)

// ObtainIpFromPool 从缓存中获取 ip 地址
func ObtainIpFromPool() (ip string, isExists bool) {
	ip, isExists = cache["ip"]
	return
}

// UpdateIpPool 更新本地缓存
func UpdateIpPool(ip string) {
	cache["ip"] = ip
}

// ObtainPubIp 获取公网 ip 地址
func ObtainPubIp() string {
	resp, err := http.Get("https://myexternalip.com/raw")
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	content, _ := ioutil.ReadAll(resp.Body)
	return string(content)
}

// Update 更新域名解析以及本地缓存
func Update(subDomain *alidns.Record, config *config.Config, externalIp string) {
	notify.UpdateDomain(subDomain, config)
	UpdateIpPool(externalIp)
}
