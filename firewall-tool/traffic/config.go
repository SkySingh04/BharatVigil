package traffic

import (
	"fmt"
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/sergi/go-diff/diffmatchpatch"
	"gopkg.in/yaml.v2"
)

// Config represents the entire configuration structure.
type Config struct {
	Firewall   FirewallConfig   `yaml:"firewall"`
	Monitoring MonitoringConfig `yaml:"monitoring"`
	AI_ML      AI_MLConfig      `yaml:"ai_ml"`
	Endpoints  []EndpointConfig `yaml:"endpoints"`
	WebConsole WebConsoleConfig `yaml:"web_console"`
	Logging    LoggingConfig    `yaml:"logging"`
	Network    NetworkConfig    `yaml:"network"`
}

// FirewallConfig represents the firewall rules.
type FirewallConfig struct {
	Rules []FirewallRule `yaml:"rules"`
}

// FirewallRule represents a single firewall rule for an application.
type FirewallRule struct {
	ID              int      `yaml:"id"`
	Application     string   `yaml:"application"`
	AllowedDomains  []string `yaml:"allowed_domains"`
	BlockedDomains  []string `yaml:"blocked_domains"`
	AllowedIPs      []string `yaml:"allowed_ips"`
	BlockedIPs      []string `yaml:"blocked_ips"`
	Protocols       []string `yaml:"protocols"`
}

// MonitoringConfig represents the monitoring configuration.
type MonitoringConfig struct {
	Enable           bool   `yaml:"enable"`
	LogFile          string `yaml:"log_file"`
	AlertThresholds  AlertThresholds `yaml:"alert_thresholds"`
}

// AlertThresholds represents the thresholds for generating alerts.
type AlertThresholds struct {
	AbnormalTraffic  int `yaml:"abnormal_traffic"`
	BlockedAttempts  int `yaml:"blocked_attempts"`
}

// AI_MLConfig represents the AI/ML configuration.
type AI_MLConfig struct {
	ModelEndpoint        string `yaml:"model_endpoint"`
	EnableAnomalyDetection bool   `yaml:"enable_anomaly_detection"`
}

// EndpointConfig represents a single endpoint configuration.
type EndpointConfig struct {
	ID        string `yaml:"id"`
	OS        string `yaml:"os"`
	IPAddress string `yaml:"ip_address"`
	Hostname  string `yaml:"hostname"`
}

// WebConsoleConfig represents the configuration for the web console.
type WebConsoleConfig struct {
	Port         int      `yaml:"port"`
	AllowedIPs   []string `yaml:"allowed_ips"`
	BlockedIPs   []string `yaml:"blocked_ips"`
	AdminUsers   []string `yaml:"admin_users"`
}

// LoggingConfig represents the logging configuration.
type LoggingConfig struct {
	LogLevel   string `yaml:"log_level"`
	LogFile    string `yaml:"log_file"`
	MaxSize    int    `yaml:"max_size"`
	MaxBackups int    `yaml:"max_backups"`
	MaxAge     int    `yaml:"max_age"`
}

// NetworkConfig represents the network-related settings.
type NetworkConfig struct {
	EnableDeepPacketInspection bool   `yaml:"enable_deep_packet_inspection"`
	SelfSignedTLSCert          string `yaml:"self_signed_tls_cert"`
	ProxyEnabled               bool   `yaml:"proxy_enabled"`
}

// LoadConfig loads the configuration from a YAML file
func LoadConfig(configPath string) (*Config, error) {
	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(file, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// String converts the Config to a YAML string for easy comparison
func (c *Config) String() string {
	data, err := yaml.Marshal(c)
	if err != nil {
		return ""
	}
	return string(data)
}


func WatchConfigFile(configPath string , currentConfig *Config) {
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
				newConfig, err := LoadConfig(configPath)
				if err != nil {
					log.Printf("Error reloading config: %v", err)
					continue
				}

				// Log the changes
				logConfigChanges(currentConfig, newConfig)

				// Update the current config to the new one
				currentConfig = newConfig

				log.Println("Config reloaded successfully")
				log.Println("Blocking all network traffic based on modified config...")
				BlockNetworkTraffic(currentConfig)

			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Printf("Error watching config file: %v", err)
		}
	}
}

func logConfigChanges(oldConfig, newConfig *Config) {
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
