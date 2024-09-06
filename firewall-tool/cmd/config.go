package cmd

import (
	"firewall-tool/traffic"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var configFile string

func reloadConfig(cmd *cobra.Command, args []string) {
	cfg, err := traffic.LoadConfig(configFile)
	if err != nil {
		fmt.Println("Error loading config:", err)
		os.Exit(1)
	}
	fmt.Println("Config reloaded successfully via command")
	//print the config
	fmt.Println(cfg)
}

func init() {
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "config.yaml", "Config file")
	rootCmd.AddCommand(&cobra.Command{
		Use:   "reload-config",
		Short: "Reload the configuration file",
		Run:   reloadConfig,
	})
}
