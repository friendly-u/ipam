package main

import (
	"flag"

	"git.ucloudadmin.com/leesin/ipam/ipam"
	log "github.com/gogap/logrus"
)

func main() {
	ipamConfigFile := flag.String("ipam-config", "ipam.json", "config file for IPAM")
	flag.Parse()

	//初始化 IPAM
	ipam.InitIPAM(*ipamConfigFile)

	// 分配 IP
	ip := ipam.IpAllocator.Allocate()
	if ip == nil {
		log.Error("No can allocate ip.")
	}
	log.Info("Allocate your ip is ", ip)

	// 释放 IP
	ipam.IpAllocator.Release(*ip)

	// 获取总的IP数
	log.Info("所有IP总数: ", ipam.IpAllocator.Size)
}
