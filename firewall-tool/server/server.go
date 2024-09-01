package server

import (
	"firewall-tool/traffic"
	"net/http"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var (
	configPath = "../config.yaml" // Use your actual path to the config file
	configLock sync.RWMutex       // Ensures thread-safe access to the config
	currentConfig *traffic.Config // Holds the current configuration
)

// StartServer initializes the Gin server and sets up routes.
func StartServer(logger *zap.Logger) {
	logger.Info("Starting the Gin server...")
	r := gin.Default()

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

		var newConfig traffic.Config
		if err := c.BindJSON(&newConfig); err != nil {
			logger.Error("Invalid request to update config", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		// Update the current config in memory
		currentConfig = &newConfig

		// Write the updated config to the file
		err := os.WriteFile(configPath, []byte(newConfig.String()), 0644)
		if err != nil {
			logger.Error("Failed to write config file", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update config"})
			return
		}

		logger.Info("Config updated successfully")
		c.JSON(http.StatusOK, gin.H{"status": "Config updated"})
	})

	// Placeholder for network activity endpoints
	r.GET("/network-activity", func(c *gin.Context) {
		// Implement the network activity logic here
		logger.Info("Sending network activity")
		c.JSON(http.StatusOK, gin.H{"status": "Network activity Sending"})
	})

	// Start the server
	if err := r.Run(":8080"); err != nil {
		logger.Fatal("Failed to start the Gin server", zap.Error(err))
	}
}
