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

type hardwareAddrSlice []net.HardwareAddr

func (s hardwareAddrSlice) fill(n int) {
	slice := make(hardwareAddrSlice, n)

	for i := 0; i < n; i++ {
		buf := make([]byte, 6)
		_, err := rand.Read(buf)
		if err != nil {
			panic(err)
		}
		slice = append(slice, buf)
	}

	copy(s, slice)
}

func (s hardwareAddrSlice) rand() net.HardwareAddr {
	return s[rand.Intn(len(s))]
}

func newMacAddrsSpoofer(n int) *hardwareAddrSlice {
	hardwareAddrs := make(hardwareAddrSlice, n)
	hardwareAddrs.fill(n)
	return &hardwareAddrs
}

func newRandomPayload(n int) []byte {
	buf := make([]byte, n)
	rand.Read(buf)
	return buf
}
