package main

import (
	"firewall-tool/config"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/spf13/cobra"
)

var (
	disableConsole bool
	disableMLModel bool
	configPath     = "../config.yaml" // Path to the config file
	currentConfig  *config.Config
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "bharat-rakshak",
		Short: "BharatVigil - A Centralized Application-Context Aware Firewall",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Starting BharatVigil tool...")

			// Load the initial config
			var err error
			currentConfig, err = config.LoadConfig(configPath)
			if err != nil {
				log.Fatalf("Error loading config: %v", err)
			}

			// Start watching the config file for changes
			go watchConfigFile(configPath)

			var wg sync.WaitGroup

			// Start the web console if not disabled
			if !disableConsole {
				wg.Add(1)
				go func() {
					defer wg.Done()
					startWebConsole()
				}()
			}

			// Start the ML model if not disabled
			if !disableMLModel {
				wg.Add(1)
				go func() {
					defer wg.Done()
					startMLModel()
				}()
			}

			// Keep the main goroutine alive until all other goroutines have completed
			wg.Wait()
			fmt.Println("BharatVigil tool stopped.")
		},
	}

	// Define flags for disabling web console and ML model
	rootCmd.Flags().BoolVar(&disableConsole, "disable-console", false, "Disable the web console")
	rootCmd.Flags().BoolVar(&disableMLModel, "disable-ml", false, "Disable the ML model")

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func startWebConsole() {
	consoleDir := "../web-console"
	absPath, err := filepath.Abs(consoleDir)
	if err != nil {
		log.Fatalf("Failed to find web console directory: %v", err)
	}

	cmd := exec.Command("npm", "run", "dev") // or "yarn dev"
	cmd.Dir = absPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("Starting Web Console...")
	if err := cmd.Start(); err != nil {
		log.Fatalf("Failed to start the web console: %v", err)
	}

	// Wait for the web console process to finish
	if err := cmd.Wait(); err != nil {
		log.Printf("Web console exited with error: %v", err)
	} else {
		fmt.Println("Web console stopped.")
	}
}

func startMLModel() {
	modelDir := "../ml-model"
	absPath, err := filepath.Abs(modelDir)
	if err != nil {
		log.Fatalf("Failed to find ML model directory: %v", err)
	}

	cmd := exec.Command("python", "model.py") // Adjust the command as per your ML model
	cmd.Dir = absPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("Starting ML Model...")
	if err := cmd.Start(); err != nil {
		log.Fatalf("Failed to start the ML model: %v", err)
	}

	// Wait for the ML model process to finish
	if err := cmd.Wait(); err != nil {
		log.Printf("ML model exited with error: %v", err)
	} else {
		fmt.Println("ML model stopped.")
	}
}

func watchConfigFile(configPath string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("Failed to create watcher: %v", err)
	}
	defer watcher.Close()

	err = watcher.Add(configPath)
	if err != nil {
		log.Fatalf("Failed to watch config file: %v", err)
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				fmt.Println("Config file changed, reloading...")

				// Load the new config
				newConfig, err := config.LoadConfig(configPath)
				if err != nil {
					log.Printf("Error reloading config: %v", err)
					continue
				}

				// Log the changes
				logConfigChanges(currentConfig, newConfig)

				// Update the current config to the new one
				currentConfig = newConfig

				fmt.Println("Config reloaded successfully.")
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Printf("Error watching config file: %v", err)
		}
	}
}

func logConfigChanges(oldConfig, newConfig *config.Config) {
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
			log.Printf("Added: %s", diff.Text)
		case diffmatchpatch.DiffDelete:
			log.Printf("Removed: %s", diff.Text)
		case diffmatchpatch.DiffEqual:
			// Equal parts can be skipped or logged depending on your preference
		}
	}
}
