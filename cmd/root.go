package cmd

import (
	"github.com/charmbracelet/log"

	"github.com/noahdotpy/prism-sharer/config"
	"github.com/remeh/userdir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "prism-sharer",
	Short: "Prism Sharer allows you to share stuff (like world saves) between Prism Launcher instances",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

var (
	configFile   string
	isVerboseLog bool
	loadedConfig config.Config
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	cobra.OnInitialize(initCobra)
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "config file to use")
	rootCmd.PersistentFlags().BoolVarP(&isVerboseLog, "verbose", "v", false, "extra log messages")
}

func initCobra() {
	if isVerboseLog {
		log.SetLevel(log.DebugLevel)
	}

	initConfig()
}

func initConfig() {
	if configFile != "" {
		viper.SetConfigFile(configFile)
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

	if err := viper.Unmarshal(&loadedConfig); err != nil {
		log.Fatalf("Can't unmarshall config: %v", err)
	}
}
