package configs

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config represents the application configuration
type Config struct {
	AppName string        `mapstructure:"appname"`
	Trace   TraceConfig   `mapstructure:"trace"`
	Cache   CacheConfig   `mapstructure:"cache"`
	DB      DBConfig      `mapstructure:"db"`
	Info    InfoConfig    `mapstructure:"info"`
	Batch   BatchConfig   `mapstructure:"batch"`
	APIClients APIClientsConfig `mapstructure:"api_clients"`
}

// TraceConfig holds distributed tracing configuration
type TraceConfig struct {
	Enabled  bool              `mapstructure:"enabled"`
	Processor TraceProcessor   `mapstructure:"processor"`
	Sampler  TraceSampler      `mapstructure:"sampler"`
}

// TraceProcessor defines the trace processor settings
type TraceProcessor struct {
	Type    string                 `mapstructure:"type"`
	Options map[string]interface{} `mapstructure:"options"`
}

// TraceSampler defines the trace sampler settings
type TraceSampler struct {
	Type    string                 `mapstructure:"type"`
	Options map[string]interface{} `mapstructure:"options"`
}

// CacheConfig holds cache configuration
type CacheConfig struct {
	RedisServer          string `mapstructure:"redisserver"`
	RedisPassword        string `mapstructure:"redispassword"`
	RedisDBIndex         int    `mapstructure:"redisdbindex"`
	RedisExpirationTime  string `mapstructure:"redisexpirationtime"`
	LCCapacity           int    `mapstructure:"lccapacity"`
	LCNumShards          int    `mapstructure:"lcnumshards"`
	LCTTL                string `mapstructure:"lcttl"`
	LCEvictionPercentage int    `mapstructure:"lcevictionpercentage"`
	LCMinRefreshDelay    string `mapstructure:"lcminrefreshdelay"`
	LCMaxRefreshDelay    string `mapstructure:"lcmaxrefreshdelay"`
	LCRetryBaseDelay     string `mapstructure:"lcretrybasedelay"`
	LCBatchSize          int    `mapstructure:"lcbatchsize"`
	LCBatchBufferTimeout string `mapstructure:"lcbatchbuffertimeout"`
	IsRedisEnabled       bool   `mapstructure:"isredisenabled"`
	IsLocalCacheEnabled  bool   `mapstructure:"islocalcacheenabled"`
}

// DBConfig holds database configuration
type DBConfig struct {
	Username           string `mapstructure:"username"`
	Password           string `mapstructure:"password"`
	Host               string `mapstructure:"host"`
	Port               string `mapstructure:"port"`
	Database           string `mapstructure:"database"`
	Schema             string `mapstructure:"schema"`
	MaxConns           int    `mapstructure:"maxconns"`
	MinConns           int    `mapstructure:"minconns"`
	MaxConnLifetime    int    `mapstructure:"maxconnlifetime"`
	MaxConnIdleTime    int    `mapstructure:"maxconnidletime"`
	HealthCheckPeriod  int    `mapstructure:"healthcheckperiod"`
	QueryTimeoutLow    string `mapstructure:"QueryTimeoutLow"`
	QueryTimeoutMed    string `mapstructure:"QueryTimeoutMed"`
}

// InfoConfig holds application info for Swagger
type InfoConfig struct {
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`
}

// BatchConfig holds batch job configuration
type BatchConfig struct {
	MaturityIntimation MaturityIntimationConfig `mapstructure:"maturity_intimation"`
}

// MaturityIntimationConfig holds maturity intimation batch job settings
type MaturityIntimationConfig struct {
	Enabled         bool     `mapstructure:"enabled"`
	Schedule        string   `mapstructure:"schedule"`
	DaysInAdvance   int      `mapstructure:"days_in_advance"`
	BatchSize       int      `mapstructure:"batch_size"`
	NotificationChannels []string `mapstructure:"notification_channels"`
}

// APIClientsConfig holds external API clients configuration
type APIClientsConfig struct {
	CBS                CBSAPIConfig                `mapstructure:"cbs"`
	PFMS               PFMSAPIConfig               `mapstructure:"pfms"`
	ECMS               ECMSConfig                  `mapstructure:"ecms"`
	PolicyService      PolicyServiceConfig         `mapstructure:"policy_service"`
	CustomerService    CustomerServiceConfig       `mapstructure:"customer_service"`
	NotificationService NotificationServiceConfig  `mapstructure:"notification_service"`
}

// CBSAPIConfig holds CBS API configuration
type CBSAPIConfig struct {
	Enabled       bool   `mapstructure:"enabled"`
	BaseURL       string `mapstructure:"base_url"`
	APIKey        string `mapstructure:"api_key"`
	Timeout       int    `mapstructure:"timeout"`
	RetryAttempts int    `mapstructure:"retry_attempts"`
	RetryDelay    int    `mapstructure:"retry_delay"`
}

// PFMSAPIConfig holds PFMS API configuration
type PFMSAPIConfig struct {
	Enabled       bool   `mapstructure:"enabled"`
	BaseURL       string `mapstructure:"base_url"`
	APIKey        string `mapstructure:"api_key"`
	Timeout       int    `mapstructure:"timeout"`
	RetryAttempts int    `mapstructure:"retry_attempts"`
	RetryDelay    int    `mapstructure:"retry_delay"`
}

// ECMSConfig holds ECMS configuration
type ECMSConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	BaseURL string `mapstructure:"base_url"`
	APIKey  string `mapstructure:"api_key"`
	Timeout int    `mapstructure:"timeout"`
}

// PolicyServiceConfig holds Policy Service configuration
type PolicyServiceConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	BaseURL string `mapstructure:"base_url"`
	Timeout int    `mapstructure:"timeout"`
}

// CustomerServiceConfig holds Customer Service configuration
type CustomerServiceConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	BaseURL string `mapstructure:"base_url"`
	Timeout int    `mapstructure:"timeout"`
}

// NotificationServiceConfig holds Notification Service configuration
type NotificationServiceConfig struct {
	Enabled        bool   `mapstructure:"enabled"`
	BaseURL        string `mapstructure:"base_url"`
	APIKey         string `mapstructure:"api_key"`
	Timeout        int    `mapstructure:"timeout"`
	SMSEnabled     bool   `mapstructure:"sms_enabled"`
	EmailEnabled   bool   `mapstructure:"email_enabled"`
	WhatsappEnabled bool  `mapstructure:"whatsapp_enabled"`
}

// LoadConfig loads configuration from YAML file
func LoadConfig(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	// Read environment variables
	viper.AutomaticEnv()

	// Read configuration file
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

// LoadConfigFromEnv loads configuration from environment variables
// This is used to override sensitive values like API keys
func LoadConfigFromEnv(config *Config) {
	// Override CBS API key from environment
	if cbsAPIKey := viper.GetString("CBS_API_KEY"); cbsAPIKey != "" {
		config.APIClients.CBS.APIKey = cbsAPIKey
	}

	// Override PFMS API key from environment
	if pfmsAPIKey := viper.GetString("PFMS_API_KEY"); pfmsAPIKey != "" {
		config.APIClients.PFMS.APIKey = pfmsAPIKey
	}

	// Override ECMS API key from environment
	if ecmsAPIKey := viper.GetString("ECMS_API_KEY"); ecmsAPIKey != "" {
		config.APIClients.ECMS.APIKey = ecmsAPIKey
	}

	// Override Notification Service API key from environment
	if notificationAPIKey := viper.GetString("NOTIFICATION_API_KEY"); notificationAPIKey != "" {
		config.APIClients.NotificationService.APIKey = notificationAPIKey
	}

	// Override database password from environment
	if dbPassword := viper.GetString("DB_PASSWORD"); dbPassword != "" {
		config.DB.Password = dbPassword
	}
}
