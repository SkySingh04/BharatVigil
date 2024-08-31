package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
    Firewall struct {
        Rules []struct {
            ID           int      `yaml:"id"`
            Application  string   `yaml:"application"`
            AllowedDomains []string `yaml:"allowed_domains"`
            BlockedDomains []string `yaml:"blocked_domains"`
            AllowedIPs   []string `yaml:"allowed_ips"`
            BlockedIPs   []string `yaml:"blocked_ips"`
            Protocols    []string `yaml:"protocols"`
        } `yaml:"rules"`
    } `yaml:"firewall"`
    
    Monitoring struct {
        Enable          bool   `yaml:"enable"`
        LogFile         string `yaml:"log_file"`
        AlertThresholds struct {
            AbnormalTraffic  int `yaml:"abnormal_traffic"`
            BlockedAttempts  int `yaml:"blocked_attempts"`
        } `yaml:"alert_thresholds"`
    } `yaml:"monitoring"`
    
    AIML struct {
        ModelEndpoint          string `yaml:"model_endpoint"`
        EnableAnomalyDetection bool   `yaml:"enable_anomaly_detection"`
    } `yaml:"ai_ml"`
    
    Endpoints []struct {
        ID        string `yaml:"id"`
        OS        string `yaml:"os"`
        IPAddress string `yaml:"ip_address"`
        Hostname  string `yaml:"hostname"`
    } `yaml:"endpoints"`
    
    WebConsole struct {
        Port       int      `yaml:"port"`
        AllowedIPs []string `yaml:"allowed_ips"`
        BlockedIPs []string `yaml:"blocked_ips"`
        AdminUsers []string `yaml:"admin_users"`
    } `yaml:"web_console"`
    
    Logging struct {
        LogLevel   string `yaml:"log_level"`
        LogFile    string `yaml:"log_file"`
        MaxSize    int    `yaml:"max_size"`
        MaxBackups int    `yaml:"max_backups"`
        MaxAge     int    `yaml:"max_age"`
    } `yaml:"logging"`
    
    Network struct {
        EnableDeepPacketInspection bool   `yaml:"enable_deep_packet_inspection"`
        SelfSignedTLSCert          string `yaml:"self_signed_tls_cert"`
        ProxyEnabled               bool   `yaml:"proxy_enabled"`
    } `yaml:"network"`
}

func LoadConfig(filename string) (*Config, error) {
    config := &Config{}
    data, err := os.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    err = yaml.Unmarshal(data, config)
    if err != nil {
        return nil, err
    }
    return config, nil
}
