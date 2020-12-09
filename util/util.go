package util

import (
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
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
	resp, err := http.Get("http://ip.cip.cc")
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	content, _ := ioutil.ReadAll(resp.Body)
	return strings.Trim(string(content), "\n")
}

// IsValidIp 验证 ip 是否有效
func IsValidIp(ip string) (isMatch bool, err error) {
	isMatch, err = regexp.MatchString("^((25[0-5]|2[0-4]\\d|[01]?\\d\\d?)\\.){3}(25[0-5]|2[0-4]\\d|[01]?\\d\\d?)$", ip)
	return
}
