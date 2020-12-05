package notify

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"log"
	"org.tubetrue01/domain-update/config"
)

// ObtainDomain 获取域名信息
func ObtainDomain(config *config.Config) *alidns.Record {
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

// UpdateDomain 更新远程域名解析
func UpdateDomain(subDomain *alidns.Record, config *config.Config) {
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
