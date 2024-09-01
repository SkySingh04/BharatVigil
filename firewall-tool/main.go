package main

import (
	"firewall-tool/cmd"
	"firewall-tool/utils/log"
	"fmt"
)

func main() {
	logger, err := log.New()
	if err != nil {
		fmt.Println("Failed to start the logger for the CLI", err)
		return
	}

	// Execute the root command to start the CLI and associated services
	cmd.Execute(logger)
}
