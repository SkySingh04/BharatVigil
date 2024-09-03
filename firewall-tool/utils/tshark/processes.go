package tshark

import (
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

func StartTshark() error {
    cmd := exec.Command("tshark")
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    return cmd.Start()
}