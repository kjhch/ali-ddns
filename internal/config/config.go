package config

import (
	"encoding/json"
	"os"
)

type DnsConf struct {
	AccessKeyId     string
	AccessKeySecret string
	DomainName      string
	RR              string
}

func LoadFromFile(path string) *DnsConf {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	result := new(DnsConf)
	err = decoder.Decode(result)
	if err != nil {
		panic(err)
	}
	return result
}
