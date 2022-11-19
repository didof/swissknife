package synflood

import (
	"context"
	"math/rand"
	"net"
	"time"

	"github.com/didof/swissknife/internal/logger"
	"github.com/didof/swissknife/internal/version"
	"github.com/pkg/errors"
	"go.uber.org/zap"
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

	_, err := resolveHost(ctx, host)
	if errors.Is(err, ErrDNSLookup) {
		return errors.Wrap(err, "unable to resolve host")
	}

	return nil
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
