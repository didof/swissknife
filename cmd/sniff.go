package swissknife

import (
	"fmt"
	"log"
	"os"

	"github.com/didof/swissknife/sniffer"
	"github.com/spf13/cobra"
)

var sniffCmd = &cobra.Command{
	Use:   "sniff",
	Short: "sniff packets",
	Run: func(cmd *cobra.Command, args []string) {
		device, err := cmd.Flags().GetString("interface")
		if err != nil {
			log.Fatal(err)
			return
		}

		if device == "" {
			fmt.Print("Specify one of the following devices.\n\n")
			sniffer.PrintDevicesInfo()
			os.Exit(1)
		}

		sniffer.Run(device)
	},
}

func init() {
	rootCmd.AddCommand(sniffCmd)

	sniffCmd.Flags().StringP("interface", "i", "", "Network interface")
}
