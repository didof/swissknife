package synflood

import (
	"math/rand"
	"net"
)

type intSlice []int

func (s intSlice) fill(min, max int) {
	slice := make([]int, max-min)
	for i := min; i < max; i++ {
		s = append(slice, i)
	}
	copy(s, slice)
}

func (s intSlice) rand() int {
	return s[rand.Intn(len(s))]
}

func newPortsSpoofer() *intSlice {
	ports := make(intSlice, 65535-1024)
	ports.fill(1024, 65535)
	return &ports
}

type ipSlice []net.IP

func (s ipSlice) fill(n int) {
	randByte := func() byte {
		return byte(rand.Intn(256))
	}

	slice := make([]net.IP, n)

	for i := 0; i < n; i++ {
		ip := net.IPv4(randByte(), randByte(), randByte(), randByte())
		slice = append(slice, ip)
	}

	copy(s, slice)
}

func (s ipSlice) rand() net.IP {
	return s[rand.Intn(len(s))]
}

func newIpsSpoofer(n int) *ipSlice {
	ips := make(ipSlice, n)
	ips.fill(n)
	return &ips
}
