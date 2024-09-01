package cmd

import (
	"firewall-tool/server"
	"firewall-tool/traffic"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	disableConsole bool
	disableMLModel bool
	disableServer  bool
	disableEBPF    bool
	configPath     string // Variable to hold the path to the config file
	currentConfig  *traffic.Config
	logger         *zap.Logger
)

var rootCmd = &cobra.Command{
	Use:   "bharat-rakshak",
	Short: "BharatVigil - A Centralized Application-Context Aware Firewall",
	Run: func(cmd *cobra.Command, args []string) {
		logger.Info("Starting BharatVigil tool...")

		// Load the initial config
		var err error
		currentConfig, err = traffic.LoadConfig(configPath)
		if err != nil {
			logger.Fatal("Error loading config", zap.Error(err))
		}

		// Start watching the config file for changes
		go watchConfigFile(configPath)

		var wg sync.WaitGroup

		// Start the web console if not disabled
		if !disableConsole {
			wg.Add(1)
			go func() {
				defer wg.Done()
				startWebConsole(logger)
			}()
		}

		// Start the ML model if not disabled
		if !disableMLModel {
			wg.Add(1)
			go func() {
				defer wg.Done()
				startMLModel(logger)
			}()
		}

		// Start the backend server if not disabled
		if !disableServer {
			wg.Add(1)
			go func() {
				defer wg.Done()
				server.StartServer(logger) // Start the Gin server
			}()
		}

		// Setup eBPF if not disabled
		if !disableEBPF {
			traffic.Setup_bpf(logger)
		} else {
			logger.Info("eBPF is disabled.")
		}

		// Keep the main goroutine alive until all other goroutines have completed
		wg.Wait()
		logger.Info("BharatVigil tool stopped.")
	},
}



// Execute executes the root command.
func Execute(l *zap.Logger) {
	logger = l
	if err := rootCmd.Execute(); err != nil {
		logger.Error("Error executing root command", zap.Error(err))
		os.Exit(1)
	}
}

func init() {
	// Define flags for disabling web console, ML model, backend server, eBPF, and setting config path
	rootCmd.Flags().BoolVar(&disableConsole, "disable-console", false, "Disable the web console")
	rootCmd.Flags().BoolVar(&disableMLModel, "disable-ml", false, "Disable the ML model")
	rootCmd.Flags().BoolVar(&disableServer, "disable-server", false, "Disable the backend server")
	rootCmd.Flags().BoolVar(&disableEBPF, "disable-ebpf", false, "Disable the eBPF program") // New eBPF flag
	rootCmd.Flags().StringVar(&configPath, "config", "../config.yaml", "Path to the configuration file")
}

func startWebConsole(logger *zap.Logger) {
	consoleDir := "../web-console"
	absPath, err := filepath.Abs(consoleDir)
	if err != nil {
		logger.Fatal("Failed to find web console directory", zap.Error(err))
	}

	cmd := exec.Command("npm", "run", "dev") // or "yarn dev"
	cmd.Dir = absPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	logger.Info("Starting Web Console...")
	if err := cmd.Start(); err != nil {
		logger.Fatal("Failed to start the web console", zap.Error(err))
	}

	// Wait for the web console process to finish
	if err := cmd.Wait(); err != nil {
		logger.Error("Web console exited with error", zap.Error(err))
	} else {
		logger.Info("Web console stopped.")
	}
}

func startMLModel(logger *zap.Logger) {
	modelDir := "../ml-model"
	absPath, err := filepath.Abs(modelDir)
	if err != nil {
		logger.Fatal("Failed to find ML model directory", zap.Error(err))
	}

	cmd := exec.Command("python", "model.py") // Adjust the command as per your ML model
	cmd.Dir = absPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	logger.Info("Starting ML Model...")
	if err := cmd.Start(); err != nil {
		logger.Fatal("Failed to start the ML model", zap.Error(err))
	}

	// Wait for the ML model process to finish
	if err := cmd.Wait(); err != nil {
		logger.Error("ML model exited with error", zap.Error(err))
	} else {
		logger.Info("ML model stopped.")
	}
}

func watchConfigFile(configPath string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logger.Fatal("Failed to create watcher", zap.Error(err))
	}
	defer watcher.Close()

	err = watcher.Add(configPath)
	if err != nil {
		logger.Fatal("Failed to watch config file", zap.Error(err))
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				logger.Info("Config file changed, reloading...")

				// Load the new config
				newConfig, err := traffic.LoadConfig(configPath)
				if err != nil {
					logger.Error("Error reloading config", zap.Error(err))
					continue
				}

				// Log the changes
				logConfigChanges(currentConfig, newConfig)

				// Update the current config to the new one
				currentConfig = newConfig

				logger.Info("Config reloaded successfully.")
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			logger.Error("Error watching config file", zap.Error(err))
		}
	}
}

func logConfigChanges(oldConfig, newConfig *traffic.Config) {
	// Convert configs to JSON strings for comparison
	oldConfigStr := oldConfig.String()
	newConfigStr := newConfig.String()

	// Use a diff library to show changes
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(oldConfigStr, newConfigStr, false)

	// Log the differences
	for _, diff := range diffs {
		switch diff.Type {
		case diffmatchpatch.DiffInsert:
			logger.Info("Added", zap.String("change", diff.Text))
		case diffmatchpatch.DiffDelete:
			logger.Info("Removed", zap.String("change", diff.Text))
		case diffmatchpatch.DiffEqual:
			// Equal parts can be skipped or logged depending on your preference
		}
	}
}
