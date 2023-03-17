package core

import (
	"fmt"
	"io"
	"net/http"

	dns "github.com/alibabacloud-go/alidns-20150109/v2/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	"github.com/kjhch/ali-dns/internal/config"
)

type DnsService struct {
	conf      *config.DnsConf
	dnsClient *dns.Client
}

func NewDnsService(conf *config.DnsConf) *DnsService {
	return &DnsService{
		conf:      conf,
		dnsClient: newClient(conf.AccessKeyId, conf.AccessKeySecret),
	}
}

func (s *DnsService) Ddns() {
	record := s.describeDomainRecord(s.conf.DomainName, s.conf.RR)
	fmt.Printf("%v: %v\n", *record.RecordId, *record.Value)
	externalIp := queryExternalIp()
	fmt.Println(externalIp)
	fmt.Println("equal?:", externalIp == *record.Value)
	if externalIp != *record.Value {
		resp := s.updateDomainRecord(*record.RecordId, externalIp, s.conf.RR)
		fmt.Println(resp)
	}
}

func (s *DnsService) describeDomainRecord(domainName, RR string) *dns.DescribeDomainRecordsResponseBodyDomainRecordsRecord {
	result, err := s.dnsClient.DescribeDomainRecords(&dns.DescribeDomainRecordsRequest{
		DomainName: newVal(domainName),
		RRKeyWord:  newVal(RR),
		Type:       newVal("A"),
	})
	if err != nil {
		panic(err)
	}
	return result.Body.DomainRecords.Record[0]
}

func (s *DnsService) updateDomainRecord(id, externalIp, RR string) *dns.UpdateDomainRecordResponse {
	resp, err := s.dnsClient.UpdateDomainRecord(&dns.UpdateDomainRecordRequest{
		RecordId: newVal(id),
		Value:    newVal(externalIp),
		RR:       newVal(RR),
		Type:     newVal("A"),
	})
	if err != nil {
		panic(err)
	}
	return resp
}

func newClient(id, sk string) *dns.Client {
	config := &openapi.Config{
		AccessKeyId:     newVal(id),
		AccessKeySecret: newVal(sk),
		Endpoint:        newVal("alidns.cn-shanghai.aliyuncs.com"),
	}
	dnsClient, err := dns.NewClient(config)
	if err != nil {
		panic(err)
	}
	return dnsClient
}

func queryExternalIp() string {
	resp, err := http.Get("http://myexternalip.com/raw")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return string(bs)
}

func newVal[T any](val T) *T {
	return &val
}
