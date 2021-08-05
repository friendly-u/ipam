package ipam

import (
	"fmt"

	log "github.com/gogap/logrus"
)

var IpAllocator *Allocator

func InitIPAM(ipamConfigFile string) {
	log.Info("Start init ipam...")
	if err := initAllocator(ipamConfigFile); err != nil {
		log.Error("Failed to init allocator, ", err)
		return
	}
	log.Info("Complete init ipam")
}

func initAllocator(ipamConfigFile string) error {
	if IpAllocator != nil {
		return nil
	}

	c, err := LoadConfigFromFile(ipamConfigFile)
	if err != nil {
		return fmt.Errorf("load ipam config failed: %s", err.Error())
	}
	ranges, err := c.IPRanges()
	if err != nil {
		return fmt.Errorf("get ip ranges failed: %s", err.Error())
	}

	IpAllocator = NewAllocator(ranges, nil)
	return nil
}
