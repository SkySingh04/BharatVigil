package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "firewall-tool",
    Short: "A CLI tool to manage the centralized application-context aware firewall",
    Long: `firewall-tool is a command-line tool to manage firewall policies,
monitor network traffic, and interact with the centralized management console.`,
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("Welcome to the firewall tool CLI!")
    },
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}

func init() {
    // Add global flags here if needed
}
