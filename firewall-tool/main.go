package main

import (
	"firewall-tool/cmd"
	"firewall-tool/utils/db"
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
	//Check if SQLite3 is installed
	if !db.IsSQLite3Installed() {
		fmt.Println("SQLite3 is not installed. Installing...")
		if err := db.InstallSQLite3(); err != nil {
			fmt.Println("Failed to install SQLite3:", err)
			return
		}
		fmt.Println("SQLite3 installed successfully.")
	}

	// Check if tshark is running
	pcapFile := "../capture.pcap"
	if !tshark.IsTsharkRunning() {
		fmt.Println("tshark is not running. Starting...")
		if err := tshark.StartTshark(pcapFile); err != nil {
			fmt.Println("Failed to start tshark:", err)
			return
		}
		fmt.Println("tshark started successfully.")
	} else {
		fmt.Println("tshark is already running.")
	}

	// Check if SQLite3 is running
	dbFile := "firewall.db"
	if !db.IsSQLite3Running() {
		fmt.Println("SQLite3 is not running. Starting...")
		if err := db.StartSQLite3(dbFile); err != nil {
			fmt.Println("Failed to start SQLite3:", err)
			return
		}
		fmt.Println("SQLite3 started successfully.")
	} else {
		fmt.Println("SQLite3 is already running.")
	}

	logger, err := log.New()
	if err != nil {
		fmt.Println("Failed to start the logger for the CLI", err)
		return
	}

	// Execute the root command to start the CLI and associated services
	cmd.Execute(logger)

}
