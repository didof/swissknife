package swissknife

import (
	"fmt"
	"os"

	"github.com/didof/swissknife/internal/version"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFileName = ".swissknife"
	ver         = version.Get()
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:     "swissknife",
	Short:   "Stuff I've learned along the way",
	Long:    `Here is where I alloc all the cool things I learn`,
	Version: ver.GitVersion,
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", fmt.Sprintf("config file (default is $HOME/%s)", cfgFileName))

	rootCmd.PersistentFlags().StringP("author", "a", "didof", "author name for copyright attribution")
	viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
	viper.SetDefault("author", "didof <didonato.fr>")

	rootCmd.AddCommand(GetSyncFloodCommand())
}

func initConfig() {
	func() {
		if cfgFile != "" {
			viper.SetConfigFile(cfgFile)
			return
		}

		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(cfgFileName)
	}()

	if err := viper.ReadInConfig(); err != nil {
		// fmt.Printf("Can't read config file ($HOME/%s).\n", cfgFileName)
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
