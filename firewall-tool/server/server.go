package server

import (
	"firewall-tool/traffic"
	"firewall-tool/utils/db"
	"firewall-tool/utils/tshark"
	"net/http"
	"os"
	"sync"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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
