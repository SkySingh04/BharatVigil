package main

import (
    "firewall-tool/cmd"
    "firewall-tool/utils/log"
    "firewall-tool/utils/tshark"
    "fmt"
)

func main() {
	// Check if tshark is installed
    if !tshark.IsTsharkInstalled() {
        fmt.Println("tshark is not installed. Installing...")
        if err := tshark.InstallTshark(); err != nil {
            fmt.Println("Failed to install tshark:", err)
            return
        }
        fmt.Println("tshark installed successfully.")
    }

    // Check if tshark is running
    if !tshark.IsTsharkRunning() {
        fmt.Println("tshark is not running. Starting...")
        if err := tshark.StartTshark(); err != nil {
            fmt.Println("Failed to start tshark:", err)
            return
        }
        fmt.Println("tshark started successfully.")
    } else {
        fmt.Println("tshark is already running.")
    }

    logger, err := log.New()
    if err != nil {
        fmt.Println("Failed to start the logger for the CLI", err)
        return
    }

    // Execute the root command to start the CLI and associated services
    cmd.Execute(logger)

}