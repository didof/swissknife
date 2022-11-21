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
	ver = version.Get()
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
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.swissknife)")

	rootCmd.PersistentFlags().StringP("author", "a", "didof", "author name for copyright attribution")
	viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
	viper.SetDefault("author", "didof <didonato.fr>")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".WirePenguin" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".WirePenguin")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		// fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
