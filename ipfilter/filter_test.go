package ipfilter

import (
	"fmt"
	"math/rand"
	"net"
	"regexp"
	"testing"
	"time"
)

func TestNewNonPublicIPFilter(t *testing.T) {
	f := NewNonPublicIPFilter()

	testcases := []struct {
		ip     string
		expect bool
	}{
		{" 127.0.0.255	", true}, // with space
		{"127.0.0.255", true},
		{"10.0.0.255", true},
		{"10.0.0.255", true},
		{"fc12:3456:789a:1::1", true},
		{"0.0.0.0", true},
		{"1.1.1.1", false},             // public ip
		{"::ffff:192.0.2.128", false},  // public IP4-mapped IPv6 address
		{"::ffff:192.168.2.128", true}, // public IP4-mapped IPv6 address
	}
	for _, testcase := range testcases {
		found, err := f.Search(testcase.ip)
		if err != nil {
			t.Errorf("unexpected error for %s: %s", testcase.ip, err)
			continue
		}
		if found != testcase.expect {
			t.Errorf("unexpected result for %s, expect %v, got %v",
				testcase.ip, testcase.expect, found)
		}
	}
}

func TestAdd(t *testing.T) {
	f := new(IPFilter)
	testcases := []struct {
		addr  string
		error bool
	}{
		{" 127.0.0.255	", true},
		{"aa/alfjdsa/adf", true},
		{"10.0.0.255/33", true},
		{"10.0.0.255/24", false},
		{"::/33", false},
	}

	for _, testcase := range testcases {
		if err := f.Add(testcase.addr); testcase.error != (err != nil) {
			t.Errorf("unexpected error for %s: %s", testcase.addr, err)
		}
	}
}

func TestSearchIP(t *testing.T) {
	f := new(IPFilter)
	if found, err := f.Search("abcdefg"); found == true || err == nil {
		t.Errorf("unexpected result for abcdefg, got %v, %s", found, err)
	}
	if found, err := f.SearchIP(nil); found == true || err == nil {
		t.Errorf("unexpected result for <nil>, got %v, %s", found, err)
	}

	f.Add("192.168.1.0/24")
	if found, err := f.Search("192.168.2.2"); found == true || err != nil {
		t.Errorf("unexpected result for 192.168.2.2, got %v, %s", found, err)
	}

	// union
	f.Add("192.168.1.0/16")
	if found, err := f.SearchIP(net.ParseIP("192.168.2.2")); found == false || err != nil {
		t.Errorf("unexpected result for 192.168.2.2, got %v, %s", found, err)
	}
}

func BenchmarkSearch(b *testing.B) {
	f := NewNonPublicIPFilter()
	rand.Seed(time.Now().Unix())
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		ip := fmt.Sprintf("%d.%d.%d.%d",
			rand.Intn(256), rand.Intn(256), rand.Intn(256), rand.Intn(256))

		for pb.Next() {
			f.Search(ip)
		}
	})
}

// BenchmarkNormalRegexBasedFilter for normal regexp-based, unreliable filter.
func BenchmarkNormalRegexBasedFilter(b *testing.B) {
	// incomplete regexp
	f := regexp.MustCompile(`(^192\.168\.|^10\.|^172\.16\.|^fc|^255\.255\.255\.255$|^0\.0\.0\.0$)`)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		ip := fmt.Sprintf("%d.%d.%d.%d",
			rand.Intn(256), rand.Intn(256), rand.Intn(256), rand.Intn(256))

		for pb.Next() {
			f.MatchString(ip)
		}
	})
}
