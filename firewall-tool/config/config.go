package config

import (
	"os"

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
