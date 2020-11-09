// Package ipfilter implements a high performance IP filter.
package ipfilter

import (
	"bytes"
	"net"
	"strings"

	"github.com/sym01/algo/avl"
)

type cidr struct {
	// The length of min and max must be net.IPv6len
	min net.IP
	max net.IP
}

// Compare implements avl.Range .
func (l *cidr) Compare(right avl.Range) int {
	r := right.(*cidr)
	if bytes.Compare(l.max, r.min) < 0 {
		return -1
	}
	if bytes.Compare(l.min, r.max) > 0 {
		return 1
	}
	return 0
}

// Contains implements avl.Range .
func (l *cidr) Contains(right avl.Range) bool {
	r := right.(*cidr)
	if bytes.Compare(l.min, r.min) > 0 {
		return false
	}
	if bytes.Compare(l.max, r.max) < 0 {
		return false
	}

	return true
}

// Union implements avl.Range .
func (l *cidr) Union(right avl.Range) avl.Range {
	r := right.(*cidr)
	ret := &cidr{
		min: l.min,
		max: l.max,
	}
	if bytes.Compare(ret.min, r.min) > 0 {
		ret.min = r.min
	}
	if bytes.Compare(ret.max, r.max) < 0 {
		ret.max = r.max
	}
	return ret
}

// IPFilter is a high performance, AVL-based IP filter.
// The filter can be used for filtering any IPv4, IPv6 and
// IP4-mapped IPv6 addresses.
//
// It's thread-safe for read ops. But if you need to read and write at the same
// time, a RWLock is necessary.
type IPFilter struct {
	tree avl.Tree
}

// Add an IP or a CIDR address into the filter.
// The addr can be a IP, such as "192.168.0.1", or a CIDR notation,
// like "192.0.2.0/24" or "2001:db8::/32" .
//
// To add mutliple IP / CIDR addresses into the filter, you can simply call it
// multi times.
func (f *IPFilter) Add(addr string) error {
	if strings.ContainsRune(addr, '/') {
		return f.addCIDR(addr)
	}

	ip := net.ParseIP(addr)
	if ip == nil {
		return &net.ParseError{Type: "IP address", Text: addr}
	}

	// to a 16-bytes representation
	ip = ip.To16()
	f.tree.Insert(&cidr{
		min: ip,
		max: ip,
	})
	return nil
}

func (f *IPFilter) addCIDR(addr string) error {
	ip, ipNet, err := net.ParseCIDR(addr)
	if err != nil {
		return err
	}

	max := ipNet.IP.To16()
	if len(max) < len(ipNet.Mask) {
		return &net.ParseError{Type: "CIDR address", Text: addr}
	}

	for i := 1; i <= len(ipNet.Mask); i++ {
		max[len(max)-i] |= ^ipNet.Mask[len(ipNet.Mask)-i]
	}
	f.tree.Insert(&cidr{
		min: ip.To16(),
		max: max,
	})
	return nil
}

// Search parses addr as an IP address, and checks if it's in the filter.
// If addr is not an IPv4 or IPv6 address, a non-nil error will be returned.
// If addr is in the filter, it will return (true, nil).
func (f *IPFilter) Search(addr string) (bool, error) {
	ip := net.ParseIP(strings.TrimSpace(addr))
	if ip == nil {
		return false, &net.ParseError{Type: "IP address", Text: addr}
	}

	ip = ip.To16()
	return f.tree.Search(&cidr{ip, ip}), nil
}

// SearchIP checks if the ip is in the filter.
// If addr is not an IPv4 or IPv6 address, a non-nil error will be returned.
// If addr is in the filter, it will return (true, nil).
func (f *IPFilter) SearchIP(ip net.IP) (bool, error) {
	ip = ip.To16()
	if ip == nil {
		return false, &net.ParseError{Type: "IP address", Text: ip.String()}
	}
	return f.tree.Search(&cidr{ip, ip}), nil
}

// NewNonGlobalUnicastIPFilter returns a new IPFilter which can filter out all
// the non-global unicast IPv4 and IPv6 IPs.
func NewNonGlobalUnicastIPFilter() *IPFilter {
	f := new(IPFilter)
	// limited broadcast
	_ = f.Add("255.255.255.255")
	// unspecified
	_ = f.Add("0.0.0.0")
	_ = f.Add("::")
	// loopback
	_ = f.Add("127.0.0.0/8")
	_ = f.Add("::1")
	// multicast
	_ = f.Add("224.0.0.0/4")
	_ = f.Add("ff00::/8")
	// link-local
	_ = f.Add("169.254.0.0/16")
	_ = f.Add("fe80::/10")

	return f
}

// NewNonPublicIPFilter returns a new IPFilter which can filter out all the
// non-public IPv4 and IPv6 IPs, such as private addresses.
func NewNonPublicIPFilter() *IPFilter {
	f := NewNonGlobalUnicastIPFilter()
	// RFC1918
	_ = f.Add("10.0.0.0/8")
	_ = f.Add("172.16.0.0/12")
	_ = f.Add("192.168.0.0/16")

	// IPv6 private
	_ = f.Add("fc00::/7")
	return f
}
