package synflood

import (
	"context"
	"fmt"
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

	for {
		select {
		case <-ctx.Done():
			return
		default:
			time.Sleep(time.Second)
			fmt.Println("tick")
		}
	}
}
