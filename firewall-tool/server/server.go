package server

import (
	"firewall-tool/traffic"
	"firewall-tool/utils/db"
	"firewall-tool/utils/tshark"

	"bytes"
	"encoding/json"
	"firewall-tool/traffic"
	"io"
	"net/http"
	"os"
	"sync"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

var (
	configPath    = "../config.yaml" // Use your actual path to the config file
	configLock    sync.RWMutex       // Ensures thread-safe access to the config
	currentConfig *traffic.Config    // Holds the current configuration
)

// StartServer initializes the Gin server and sets up routes.
func StartServer(logger *zap.Logger) {
	logger.Info("Starting the Gin server...")
	r := gin.Default()

	// Configure CORS to allow all origins
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Allow all origins
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Endpoint to get the current YAML configuration
	r.GET("/config", func(c *gin.Context) {
		configLock.RLock()
		defer configLock.RUnlock()

		c.File(configPath)
	})

	// Endpoint to update the YAML configuration
	r.POST("/config", func(c *gin.Context) {
		configLock.Lock()
		defer configLock.Unlock()

		// Log the incoming request
		logger.Info("Received request to update config")

		// Read and log the body from the request
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			logger.Error("Failed to read request body", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read request body"})
			return
		}
		logger.Info("Request body", zap.String("body", string(body)))
		// Re-set the body to the request
		c.Request.Body = io.NopCloser(bytes.NewReader(body))

		// Parse the new config from the POST request
		var newConfig traffic.Config
		if err := json.Unmarshal(body, &newConfig); err != nil {
			logger.Error("Invalid request to update config", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		logger.Info("New config", zap.Any("config", newConfig))

		// Load the current YAML configuration from file
		currentData, err := os.ReadFile(configPath)
		if err != nil {
			logger.Error("Failed to read config file", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read config"})
			return
		}
		logger.Info("Current config", zap.String("config", string(currentData)))

		// Initialize currentConfig if it’s nil
		if currentConfig == nil {
			currentConfig = &traffic.Config{}
		}

		// Parse the current config
		if err := yaml.Unmarshal(currentData, currentConfig); err != nil {
			logger.Error("Failed to parse current config", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse config"})
			return
		}
		logger.Info("Current config after unmarshaling", zap.Any("config", currentConfig))

		// Merge the new config with the current config
		currentConfig.Firewall.Rules = mergeFirewallRules(currentConfig.Firewall.Rules, newConfig.Firewall.Rules)

		// Update other parts of the config
		// if !newConfig.Monitoring.IsZeroValue() {
		// 	currentConfig.Monitoring = newConfig.Monitoring
		// }
		// if !newConfig.AI_ML.IsZeroValue() {
		// 	currentConfig.AI_ML = newConfig.AI_ML
		// }
		// if len(newConfig.Endpoints) > 0 {
		// 	currentConfig.Endpoints = newConfig.Endpoints
		// }
		// if !newConfig.WebConsole.IsZeroValue() {
		// 	currentConfig.WebConsole = newConfig.WebConsole
		// }
		// if !newConfig.Logging.IsZeroValue() {
		// 	currentConfig.Logging = newConfig.Logging
		// }
		// if !newConfig.Network.IsZeroValue() {
		// 	currentConfig.Network = newConfig.Network
		// }

		// Convert the updated config back to YAML
		updatedConfigData, err := yaml.Marshal(currentConfig)
		if err != nil {
			logger.Error("Failed to marshal updated config", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to serialize config"})
			return
		}

		// Write the updated config back to the file
		err = os.WriteFile(configPath, updatedConfigData, 0644)

		if err != nil {
			logger.Error("Failed to write config file", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update config"})
			return
		}

		logger.Info("Config updated successfully")
		c.JSON(http.StatusOK, gin.H{"status": "Config updated"})
	})

	r.GET("/events", tshark.SseHandler)

	r.GET("/requests", func(c *gin.Context) {
		requests, err := db.GetAllRequests()
		if err != nil {
			logger.Error("Failed to get requests", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get requests"})
			return
		}

		c.JSON(http.StatusOK, requests)
	})

	// Start the server
	if err := r.Run(":8080"); err != nil {
		logger.Fatal("Failed to start the Gin server", zap.Error(err))
	}
}

// Helper function to merge firewall rules
func mergeFirewallRules(existingRules, newRules []traffic.FirewallRule) []traffic.FirewallRule {
	ruleMap := make(map[int]traffic.FirewallRule) // Using ID as the key for merging

	// Add existing rules to map
	for _, rule := range existingRules {
		ruleMap[rule.ID] = rule
	}

	// Add new rules to map, updating if necessary
	for _, rule := range newRules {
		ruleMap[rule.ID] = rule
	}

	// Convert map back to slice
	mergedRules := make([]traffic.FirewallRule, 0, len(ruleMap))
	for _, rule := range ruleMap {
		mergedRules = append(mergedRules, rule)
	}

	return mergedRules
}
