package ipam

import (
	"encoding/json"
	"io/ioutil"
)

// IPAMConfig is configuration for Allocator
type IPAMConfig struct {
	Ranges []IPRangeConfig `json:"ranges"`
}

// IPRangeConfig is configuration for IPRange
type IPRangeConfig struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

// LoadConfigFromFile loads a json file and create an IPAMConfig
func LoadConfigFromFile(f string) (*IPAMConfig, error) {
	data, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}
	c := IPAMConfig{}
	err = json.Unmarshal(data, &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// IPRanges create list if *IPRanges from config
func (c *IPAMConfig) IPRanges() ([]*IPRange, error) {
	ranges := []*IPRange{}
	for _, r := range c.Ranges {
		ipRange, err := NewIPRange(r.Start, r.End)
		if err != nil {
			return ranges, err
		}
		ranges = append(ranges, ipRange)
	}
	return ranges, nil
}
