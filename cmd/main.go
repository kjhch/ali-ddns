package main

import (
	"flag"

	"github.com/kjhch/ali-dns/internal/config"
	"github.com/kjhch/ali-dns/internal/core"
)

func main() {
	confPath := flag.String("c", "dns.json", "配置文件路径")
	flag.Parse()

	dnsConf := config.LoadFromFile(*confPath)
	core.NewDnsService(dnsConf).Ddns()
}
