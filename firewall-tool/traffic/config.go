package traffic

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

// Config represents the entire configuration structure.
type Config struct {
	Firewall   FirewallConfig   `yaml:"firewall" json:"firewall"`
	Monitoring MonitoringConfig `yaml:"monitoring" json:"monitoring"`
	AI_ML      AI_MLConfig      `yaml:"ai_ml" json:"ai_ml"`
	Endpoints  []EndpointConfig `yaml:"endpoints" json:"endpoints"`
	WebConsole WebConsoleConfig `yaml:"web_console" json:"web_console"`
	Logging    LoggingConfig    `yaml:"logging" json:"logging"`
	Network    NetworkConfig    `yaml:"network" json:"network"`
}

// FirewallConfig represents the firewall rules.
type FirewallConfig struct {
	Rules []FirewallRule `yaml:"rules" json:"rules"`
}

// FirewallRule represents a single firewall rule for an application.
type FirewallRule struct {
	ID             int      `yaml:"id" json:"id"`
	Application    string   `yaml:"application" json:"application"`
	AllowedDomains []string `yaml:"allowed_domains" json:"allowed_domains"`
	BlockedDomains []string `yaml:"blocked_domains" json:"blocked_domains"`
	AllowedIPs     []string `yaml:"allowed_ips" json:"allowed_ips"`
	BlockedIPs     []string `yaml:"blocked_ips" json:"blocked_ips"`
	Protocols      []string `yaml:"protocols" json:"protocols"`
}

// MonitoringConfig represents the monitoring configuration.
type MonitoringConfig struct {
	Enable          bool            `yaml:"enable" json:"enable"`
	LogFile         string          `yaml:"log_file" json:"log_file"`
	AlertThresholds AlertThresholds `yaml:"alert_thresholds" json:"alert_thresholds"`
}

// AlertThresholds represents the thresholds for generating alerts.
type AlertThresholds struct {
	AbnormalTraffic int `yaml:"abnormal_traffic" json:"abnormal_traffic"`
	BlockedAttempts int `yaml:"blocked_attempts" json:"blocked_attempts"`
}

// AI_MLConfig represents the AI/ML configuration.
type AI_MLConfig struct {
	ModelEndpoint          string `yaml:"model_endpoint" json:"model_endpoint"`
	EnableAnomalyDetection bool   `yaml:"enable_anomaly_detection" json:"enable_anomaly_detection"`
}

// EndpointConfig represents a single endpoint configuration.
type EndpointConfig struct {
	ID        string `yaml:"id" json:"id"`
	OS        string `yaml:"os" json:"os"`
	IPAddress string `yaml:"ip_address" json:"ip_address"`
	Hostname  string `yaml:"hostname" json:"hostname"`
}

// WebConsoleConfig represents the configuration for the web console.
type WebConsoleConfig struct {
	Port       int      `yaml:"port" json:"port"`
	AllowedIPs []string `yaml:"allowed_ips" json:"allowed_ips"`
	BlockedIPs []string `yaml:"blocked_ips" json:"blocked_ips"`
	AdminUsers []string `yaml:"admin_users" json:"admin_users"`
}

// LoggingConfig represents the logging configuration.
type LoggingConfig struct {
	LogLevel   string `yaml:"log_level" json:"log_level"`
	LogFile    string `yaml:"log_file" json:"log_file"`
	MaxSize    int    `yaml:"max_size" json:"max_size"`
	MaxBackups int    `yaml:"max_backups" json:"max_backups"`
	MaxAge     int    `yaml:"max_age" json:"max_age"`
}

// NetworkConfig represents the network-related settings.
type NetworkConfig struct {
	EnableDeepPacketInspection bool   `yaml:"enable_deep_packet_inspection" json:"enable_deep_packet_inspection"`
	SelfSignedTLSCert          string `yaml:"self_signed_tls_cert" json:"self_signed_tls_cert"`
	ProxyEnabled               bool   `yaml:"proxy_enabled" json:"proxy_enabled"`
}

// LoadConfig loads the configuration from a YAML file.
func LoadConfig(configPath string) (*Config, error) {
	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(file, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	return &cfg, nil
}

// String converts the Config to a YAML string for easy comparison.
func (c *Config) String() string {
	data, err := yaml.Marshal(c)
	if err != nil {
		return ""
	}
	return string(data)
}

// IsZeroValue checks if the WebConsoleConfig is zero-valued.
func (wc WebConsoleConfig) IsZeroValue() bool {
	return len(wc.AllowedIPs) == 0 &&
		len(wc.BlockedIPs) == 0 &&
		wc.Port == 0 &&
		len(wc.AdminUsers) == 0
}

// IsZeroValue checks if the LoggingConfig is zero-valued.
func (lc LoggingConfig) IsZeroValue() bool {
	return lc.LogLevel == "" &&
		lc.LogFile == "" &&
		lc.MaxSize == 0 &&
		lc.MaxBackups == 0 &&
		lc.MaxAge == 0
}

// IsZeroValue checks if the NetworkConfig is zero-valued.
func (nc NetworkConfig) IsZeroValue() bool {
	return !nc.EnableDeepPacketInspection &&
		nc.SelfSignedTLSCert == "" &&
		!nc.ProxyEnabled
}

// IsZeroValue checks if the MonitoringConfig is zero-valued.
func (mc MonitoringConfig) IsZeroValue() bool {
	return !mc.Enable &&
		mc.LogFile == "" &&
		mc.AlertThresholds.AbnormalTraffic == 0 &&
		mc.AlertThresholds.BlockedAttempts == 0
}

// IsZeroValue checks if the AI_MLConfig is zero-valued.
func (amc AI_MLConfig) IsZeroValue() bool {
	return amc.ModelEndpoint == "" &&
		!amc.EnableAnomalyDetection
}
