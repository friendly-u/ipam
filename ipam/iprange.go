package ipam

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"
	"net"
)

// IPRange represent a range of IP addresses
type IPRange struct {
	Start net.IP
	End   net.IP
}

// NewIPRange creates an IPRange
func NewIPRange(start, end string) (*IPRange, error) {
	startIP := net.ParseIP(start)
	if startIP == nil {
		return nil, fmt.Errorf("invalid ip address %s", start)
	}
	endIP := net.ParseIP(end)
	if endIP == nil {
		return nil, fmt.Errorf("invalid ip address %s", end)
	}
	if bytes.Compare(endIP, startIP) < 0 {
		return nil, errors.New("ip range end is less than start")
	}
	r := &IPRange{
		Start: startIP,
		End:   endIP,
	}
	return r, nil
}

// Contains returns true if ip is in this range
func (r *IPRange) Contains(ip net.IP) bool {
	if bytes.Compare(ip, r.Start) >= 0 && bytes.Compare(ip, r.End) <= 0 {
		return true
	}
	return false
}

// Size return total ip addresses in this range
func (r *IPRange) Size() int64 {
	endInt := big.NewInt(0).SetBytes([]byte(r.End))
	startInt := big.NewInt(0).SetBytes([]byte(r.Start))
	s := big.NewInt(0).Sub(endInt, startInt)
	return s.Int64() + 1
}
