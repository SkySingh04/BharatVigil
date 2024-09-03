package tshark

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func IsTsharkInstalled() bool {
	cmd := exec.Command("which", "tshark")
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

func InstallTshark() error {
	pkgManager := getPkgManager()
	cmd := exec.Command("sh", "-c", pkgManager+" install -y tshark")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func getPkgManager() string {
	// Determine the package manager based on the OS
	// For simplicity, assuming Linux with apt-get or yum
	if _, err := exec.LookPath("apt-get"); err == nil {
		return "sudo apt-get"
	}
	if _, err := exec.LookPath("yum"); err == nil {
		return "sudo yum"
	}
	return "sudo apt-get" // Default to apt-get if no other package manager is found
}

func IsTsharkRunning() bool {
	cmd := exec.Command("pgrep", "tshark")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(output)) != ""
}

func StartTshark(outputFile string) error {
	// Check if the output file already exists, if it exists, delete it
	if _, err := os.Stat(outputFile); err == nil {
		if err := os.Remove(outputFile); err != nil {
			return fmt.Errorf("failed to remove existing file %s: %v", outputFile, err)
		} else {
			fmt.Printf("removed existing file %s\n", outputFile)
		}
	}

	// Create the file if it doesn't exist
	file, err := os.OpenFile(outputFile, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %v", outputFile, err)
	}
	file.Close()

	// Run tshark with elevated privileges if necessary
	cmd := exec.Command("sudo", "tshark", "-P", "-F", "pcap", "-w", outputFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Start()
}
