package main

import (
	"firewall-tool/cmd"
	"firewall-tool/utils/db"
	"firewall-tool/utils/log"
	"firewall-tool/utils/tshark"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"go.uber.org/zap"
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
	// Check if SQLite3 is installed
	if !db.IsSQLite3Installed() {
		fmt.Println("SQLite3 is not installed. Installing...")
		if err := db.InstallSQLite3(); err != nil {
			fmt.Println("Failed to install SQLite3:", err)
			return
		}
		fmt.Println("SQLite3 installed successfully.")
	}

	// Check if tshark is running
	pcapFile := "./capture.pcap"
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

	// Start the command.sh script in a separate goroutine
	// go startCommandScript(logger)

	// Execute the root command to start the CLI and associated services
	cmd.Execute(logger)
}

func startCommandScript(logger *zap.Logger) {
	scriptPath := "../model_pipeline_files/command.sh" // Adjust the path to your command.sh file
	absPath, err := filepath.Abs(scriptPath)
	if err != nil {
		logger.Fatal("Failed to find command.sh script", zap.Error(err))
	}

	for {
		cmd := exec.Command("sh", absPath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		logger.Info("Starting command.sh script...")
		if err := cmd.Start(); err != nil {
			logger.Error("Failed to start command.sh script", zap.Error(err))
			return
		}

		// Wait for the script to finish
		if err := cmd.Wait(); err != nil {
			logger.Error("command.sh script exited with error", zap.Error(err))
		} else {
			logger.Info("command.sh script stopped.")
		}
	}
}
