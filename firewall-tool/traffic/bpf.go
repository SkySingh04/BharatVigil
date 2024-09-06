package traffic

import (
	"os"
	"os/signal"

	"github.com/dropbox/goebpf"
	"go.uber.org/zap"
)

func Setup_bpf(logger *zap.Logger) {
	logger.Info("Initializing eBPF setup...")

	// Specify Interface Name
	interfaceName := "lo"
	logger.Info("Using interface for eBPF program", zap.String("interface", interfaceName))

	// IP BlockList
	// Add the IPs you want to be blocked
	ipList := []string{
		"12.12.11.32",
	}
	logger.Info("IP addresses to be blocked", zap.Strings("ipList", ipList))

	// Load XDP Into App
	bpf := goebpf.NewDefaultEbpfSystem()
	err := bpf.LoadElf("bpf/xdp.elf")
	if err != nil {
		logger.Fatal("Failed to load eBPF ELF file", zap.Error(err))
	}

	blacklist := bpf.GetMapByName("blacklist")
	if blacklist == nil {
		logger.Fatal("eBPF map 'blacklist' not found")
	}
	logger.Info("eBPF map 'blacklist' successfully loaded")

	xdp := bpf.GetProgramByName("firewall")
	if xdp == nil {
		logger.Fatal("Program 'firewall' not found in ELF file")
	}
	logger.Info("eBPF program 'firewall' successfully loaded")

	err = xdp.Load()
	if err != nil {
		logger.Fatal("Failed to load eBPF program into the kernel", zap.Error(err))
	}

	err = xdp.Attach(interfaceName)
	if err != nil {
		logger.Fatal("Failed to attach eBPF program to interface", zap.String("interface", interfaceName), zap.Error(err))
	}
	logger.Info("eBPF program attached successfully to interface", zap.String("interface", interfaceName))

	err = BlockIPAddress(ipList, blacklist)
	if err != nil {
		logger.Fatal("Failed to add IP addresses to the blacklist", zap.Error(err))
	}
	logger.Info("IP addresses successfully added to blacklist", zap.Strings("blockedIPs", ipList))

	defer func() {
		logger.Info("Detaching eBPF program from interface", zap.String("interface", interfaceName))
		err := xdp.Detach()
		if err != nil {
			logger.Error("Error detaching eBPF program", zap.Error(err))
		} else {
			logger.Info("eBPF program detached successfully")
		}
	}()

	// Handle Ctrl+C to stop the program
	ctrlC := make(chan os.Signal, 1)
	signal.Notify(ctrlC, os.Interrupt)
	logger.Info("eBPF Program loaded successfully into the Kernel. Press CTRL+C to stop.")
	<-ctrlC
}

// The Function That adds the IPs to the blacklist map
func BlockIPAddress(ipAddreses []string, blacklist goebpf.Map) error {
	for index, ip := range ipAddreses {
		err := blacklist.Insert(goebpf.CreateLPMtrieKey(ip), index)
		if err != nil {
			return err
		}
	}
	return nil
}
