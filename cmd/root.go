package cmd

import (
	"fmt"
	"os"
	"path"

	insertcmd "github.com/AlfredBerg/piper/cmd/insert"
	streamcmd "github.com/AlfredBerg/piper/cmd/stream"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "piper",
	Short: "A queueing tool with std-in/out as a first-class citizen",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/.piper.yaml)")

	rootCmd.AddCommand(insertcmd.NewCmdInsert())
	rootCmd.AddCommand(streamcmd.NewCmdStream())
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetEnvPrefix("PIPER") //ENVs start with PIPER_
	viper.MustBindEnv("REDIS_URL")

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".piper" (without extension).
		viper.AddConfigPath(path.Join(home, ".config/"))
		viper.SetConfigType("yaml")
		viper.SetConfigName(".piper")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
