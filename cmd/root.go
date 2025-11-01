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
	Short: "A queueing tool with std-in/out as a first-class citizen.",
	Long:  "Piper is a simple queueing tool that can use either redis as or sqlite. Default is sqlite, to use a redis url must be set e.g. with the env \"PIPER_REDIS_URL\" following the format redis://<user>:<password>@<host>:<port>/<db>",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/piper/config.yaml)")

	rootCmd.AddCommand(insertcmd.NewCmdInsert())
	rootCmd.AddCommand(streamcmd.NewCmdStream())
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Find home directory.
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)
	confDir := path.Join(home, ".config/piper")
	if _, err := os.Stat(confDir); os.IsNotExist(err) {
		err := os.MkdirAll(confDir, 0o755)
		if err != nil {
			panic(err)
		}
	}

	viper.SetEnvPrefix("PIPER") //ENVs start with PIPER_
	viper.BindEnv("REDIS_URL")
	viper.BindEnv("SQLITE_FILE")
	viper.SetDefault("sqlite_file", confDir+"/piper.sqlite")

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {

		// Search config in home directory with name ".piper" (without extension).
		viper.AddConfigPath(confDir)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
