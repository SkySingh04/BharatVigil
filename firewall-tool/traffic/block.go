package traffic

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"strings"
)

// blockNetworkTraffic blocks all network traffic for the specified applications.
func BlockNetworkTraffic(cfg *Config) {
	for _, rule := range cfg.Firewall.Rules {
		switch runtime.GOOS {
		case "linux", "darwin": // For Linux and macOS
			blockTrafficLinuxMac(rule)
		case "windows":
			blockTrafficWindows(rule)
		default:
			fmt.Printf("OS not supported for blocking network traffic: %s\n", runtime.GOOS)
		}
	}
}

// blockTrafficLinuxMac blocks network traffic using iptables (Linux) or pfctl (macOS).
func blockTrafficLinuxMac(rule FirewallRule) {
	fmt.Printf("Blocking all network traffic for application: %s (ID: %d)\n", rule.Application, rule.ID)

	// Use pgrep to find the process ID(s) of the application
	cmd := exec.Command("pgrep", "-f", rule.Application)
	output, err := cmd.Output()
	if err != nil {
		log.Printf("Failed to find process for application %s: %v", rule.Application, err)
		return
	}

	// Convert output to string and split by lines (in case multiple PIDs are found)
	pids := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(pids) == 0 {
		log.Printf("No processes found for application %s", rule.Application)
		return
	}

	for _, pid := range pids {
		if pid == "" {
			continue
		}

		// Apply iptables rule to drop all traffic for the specific PID
		cmd := exec.Command("sudo", "iptables", "-A", "OUTPUT", "-m", "owner", "--pid-owner", pid, "-j", "DROP")
		if err := cmd.Run(); err != nil {
			log.Printf("Failed to block traffic for application %s (PID: %s) using iptables: %v", rule.Application, pid, err)
		} else {
			fmt.Printf("Successfully blocked traffic for PID: %s (application: %s)\n", pid, rule.Application)
		}
	}
}
// blockTrafficWindows blocks network traffic using netsh (Windows).
func blockTrafficWindows(rule FirewallRule) {
	fmt.Printf("Blocking all network traffic for application: %s (ID: %d)\n", rule.Application, rule.ID)
	// Use Windows Firewall to block traffic, use full path or executable name if required
	cmd := exec.Command("netsh", "advfirewall", "firewall", "add", "rule", "name=BlockTraffic_"+rule.Application, "dir=out", "action=block", "program=C:\\Path\\To\\"+rule.Application+".exe")
	if err := cmd.Run(); err != nil {
		log.Printf("Failed to block traffic for application %s using netsh: %v", rule.Application, err)
	}
}
