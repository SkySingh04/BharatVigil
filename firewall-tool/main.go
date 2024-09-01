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

	// Pass the logger to the root command
	cmd.Execute(logger)
}
