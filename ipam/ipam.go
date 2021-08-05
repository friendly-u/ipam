package ipam

import (
	"bytes"
	"container/ring"
	"math/big"
	"net"
	"sort"
	"sync"
)

var emptyIP = net.ParseIP("0.0.0.0")

// Allocator manage IP address in rr mode
type Allocator struct {
	ipRangeRing  *ring.Ring
	allocatedMap map[string]bool
	Size         int64
	whereToStart net.IP
	lock         sync.Mutex
}

// NewAllocator create a new Allocator
func NewAllocator(ranges []*IPRange, allocated []net.IP) *Allocator {
	a := &Allocator{
		ipRangeRing:  ring.New(len(ranges)),
		allocatedMap: make(map[string]bool),
		Size:         0,
		whereToStart: emptyIP,
		lock:         sync.Mutex{},
	}

	for i := 0; i < a.ipRangeRing.Len(); i++ {
		a.ipRangeRing.Value = ranges[i]
		a.ipRangeRing = a.ipRangeRing.Next()
		a.Size = a.Size + ranges[i].Size()
	}

	sort.Slice(allocated, func(i, j int) bool {
		return bytes.Compare(allocated[i], allocated[j]) < 0
	})
	for _, ip := range allocated {
		if containsInRanges(ranges, ip) {
			a.allocatedMap[ip.String()] = true
			a.whereToStart = ip
		}
	}

	return a
}

// Allocate allocate an IP
func (a *Allocator) Allocate() *net.IP {
	a.lock.Lock()
	defer a.lock.Unlock()
	for !a.full() {
		ip := a.next()
		if !a.allocated(ip) {
			a.markAllocated(ip)
			return ip
		}
	}
	return nil
}

// Release release an ip
func (a *Allocator) Release(ip net.IP) {
	a.lock.Lock()
	defer a.lock.Unlock()
	if ip == nil {
		return
	}
	delete(a.allocatedMap, ip.String())
}

func (a *Allocator) allocated(ip *net.IP) bool {
	if _, found := a.allocatedMap[ip.String()]; found {
		return true
	}
	return false
}

func (a *Allocator) markAllocated(ip *net.IP) {
	a.allocatedMap[ip.String()] = true
}

func (a *Allocator) full() bool {
	if a.ipRangeRing.Len() == 0 {
		return true
	}
	if int64(len(a.allocatedMap)) >= a.Size {
		return true
	}
	return false
}

func (a *Allocator) next() *net.IP {
	if a.ipRangeRing.Len() == 0 {
		return nil
	}

	r := a.ipRangeRing.Value.(*IPRange)
	if a.whereToStart.Equal(emptyIP) {
		a.whereToStart = r.Start
		return &r.Start
	}

	for !r.Contains(a.whereToStart) {
		a.ipRangeRing = a.ipRangeRing.Next()
		r = a.ipRangeRing.Value.(*IPRange)
	}

	ip := nextIP(a.whereToStart)
	if !r.Contains(ip) {
		a.ipRangeRing = a.ipRangeRing.Next()
		r = a.ipRangeRing.Value.(*IPRange)
		a.whereToStart = r.Start
		return &r.Start
	}

	a.whereToStart = ip
	return &ip
}

func containsInRanges(ranges []*IPRange, ip net.IP) bool {
	for _, r := range ranges {
		if r.Contains(ip) {
			return true
		}
	}
	return false
}

func nextIP(ip net.IP) net.IP {
	// Convert to big.Int and increment
	ipb := big.NewInt(0).SetBytes([]byte(ip))
	ipb.Add(ipb, big.NewInt(1))

	// Add leading zeros
	b := ipb.Bytes()
	b = append(make([]byte, len(ip)-len(b)), b...)
	return net.IP(b)
}
