package synflood

import (
	"context"
	"time"

	"github.com/didof/swissknife/internal/logger"
	"github.com/didof/swissknife/internal/version"
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

	var cancel context.CancelFunc
	if opts.FloodDurationMilliseconds != -1 {
		ctx, cancel = context.WithTimeout(ctx, time.Millisecond*time.Duration(opts.FloodDurationMilliseconds))
	}

	go func() {
		if err := do(ctx, opts); err != nil {
			cancel()
			log.Fatal("an error occured on flooding process", zap.String("error", err.Error()))
		}
	}()

	for {
		select {
		case <-ctx.Done():
			cancel()
			return
		default:
			continue
		}
	}
}

func do(ctx context.Context, opts SynFloodOptions) error {
	return nil
}
