package cmd

import (
	"github.com/charmbracelet/log"
	"github.com/remeh/userdir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/noahdotpy/prism-sharer/core"
)

var rootCmd = &cobra.Command{
	Use:   "prism-sharer",
	Short: "Prism Sharer allows you to share stuff (like world saves) between Prism Launcher instances",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	cobra.OnInitialize(initCobra)
	rootCmd.PersistentFlags().StringVarP(&core.ConfigFile, "config", "c", "", "config file to use")
	rootCmd.PersistentFlags().BoolVarP(&core.IsVerboseLog, "verbose", "v", false, "extra log messages")
}

func initCobra() {
	if core.IsVerboseLog {
		log.SetLevel(log.DebugLevel)
	}

	initConfig()
}

func initConfig() {
	if core.ConfigFile != "" {
		viper.SetConfigFile(core.ConfigFile)
	} else {
		configDir := userdir.GetConfigHome() + "/prism-sharer"

		viper.AddConfigPath(configDir)
		viper.SetConfigName("config.json")
	}

	viper.SetConfigType("json")

	dataHome := userdir.GetDataHome()

	viper.SetDefault("storeDir", dataHome+"/prism-sharer/")
	viper.SetDefault("instancesDir", dataHome+"/PrismLauncher/instances/")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Can't read config: %v", err)
	}

	if err := viper.Unmarshal(&core.Config); err != nil {
		log.Fatalf("Can't unmarshall config: %v", err)
	}
}
