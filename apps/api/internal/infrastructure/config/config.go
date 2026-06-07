package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

// Config holds all application configuration.
type Config struct {
	App      AppConfig
	DB       DBConfig
	Redis    RedisConfig
	JWT      JWTConfig
	S3       S3Config
	Midtrans MidtransConfig
	Google   GoogleConfig
}

type AppConfig struct {
	Name string
	Env  string
	Port string
	URL  string
}

type DBConfig struct {
	Host         string
	Port         string
	Name         string
	User         string
	Password     string
	MaxIdleConns int
	MaxOpenConns int
	MaxLifetime  time.Duration
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type JWTConfig struct {
	AccessSecret  string
	RefreshSecret string
	AccessExpiry  time.Duration
	RefreshExpiry time.Duration
}

type S3Config struct {
	Endpoint  string
	Bucket    string
	AccessKey string
	SecretKey string
	Region    string
	UseSSL    bool
}

type MidtransConfig struct {
	ServerKey string
	ClientKey string
	Env       string // "sandbox" or "production"
}

type GoogleConfig struct {
	ClientID     string
	ClientSecret string
}

// Load reads configuration from environment variables and .env file.
func Load() (*Config, error) {
	readEnvFile()
	viper.AutomaticEnv()

	cfg := &Config{
		App: AppConfig{
			Name: getStringOrDefault("APP_NAME", "rent-outdoor-api"),
			Env:  getStringOrDefault("APP_ENV", "development"),
			Port: getStringOrDefault("APP_PORT", "8080"),
			URL:  getStringOrDefault("APP_URL", "http://localhost:8080"),
		},
		DB: DBConfig{
			Host:         getStringOrDefault("DB_HOST", "localhost"),
			Port:         getStringOrDefault("DB_PORT", "3306"),
			Name:         getStringOrDefault("DB_NAME", "rent_outdoor"),
			User:         getStringOrDefault("DB_USER", "root"),
			Password:     viper.GetString("DB_PASSWORD"),
			MaxIdleConns: viper.GetInt("DB_MAX_IDLE_CONNS"),
			MaxOpenConns: viper.GetInt("DB_MAX_OPEN_CONNS"),
			MaxLifetime:  time.Duration(viper.GetInt("DB_MAX_LIFETIME_MINUTES")) * time.Minute,
		},
		Redis: RedisConfig{
			Host:     getStringOrDefault("REDIS_HOST", "localhost"),
			Port:     getStringOrDefault("REDIS_PORT", "6379"),
			Password: viper.GetString("REDIS_PASSWORD"),
			DB:       viper.GetInt("REDIS_DB"),
		},
		JWT: JWTConfig{
			AccessSecret:  viper.GetString("JWT_ACCESS_SECRET"),
			RefreshSecret: viper.GetString("JWT_REFRESH_SECRET"),
			AccessExpiry:  parseDuration("JWT_ACCESS_EXPIRY", 15*time.Minute),
			RefreshExpiry: parseDuration("JWT_REFRESH_EXPIRY", 7*24*time.Hour),
		},
		S3: S3Config{
			Endpoint:  viper.GetString("S3_ENDPOINT"),
			Bucket:    viper.GetString("S3_BUCKET"),
			AccessKey: viper.GetString("S3_ACCESS_KEY"),
			SecretKey: viper.GetString("S3_SECRET_KEY"),
			Region:    getStringOrDefault("S3_REGION", "us-east-1"),
			UseSSL:    viper.GetBool("S3_USE_SSL"),
		},
		Midtrans: MidtransConfig{
			ServerKey: viper.GetString("MIDTRANS_SERVER_KEY"),
			ClientKey: viper.GetString("MIDTRANS_CLIENT_KEY"),
			Env:       getStringOrDefault("MIDTRANS_ENV", "sandbox"),
		},
		Google: GoogleConfig{
			ClientID:     viper.GetString("GOOGLE_CLIENT_ID"),
			ClientSecret: viper.GetString("GOOGLE_CLIENT_SECRET"),
		},
	}

	// Set defaults
	if cfg.DB.MaxIdleConns == 0 {
		cfg.DB.MaxIdleConns = 10
	}
	if cfg.DB.MaxOpenConns == 0 {
		cfg.DB.MaxOpenConns = 100
	}
	if cfg.DB.MaxLifetime == 0 {
		cfg.DB.MaxLifetime = 5 * time.Minute
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func readEnvFile() {
	for _, candidate := range []string{
		".env",
		filepath.Join("apps", "api", ".env"),
	} {
		if _, err := os.Stat(candidate); err != nil {
			continue
		}

		viper.SetConfigFile(candidate)
		_ = viper.ReadInConfig()
		return
	}
}

func (c *Config) validate() error {
	if c.JWT.AccessSecret == "" && c.App.Env == "production" {
		return fmt.Errorf("JWT_ACCESS_SECRET is required in production")
	}
	if c.JWT.RefreshSecret == "" && c.App.Env == "production" {
		return fmt.Errorf("JWT_REFRESH_SECRET is required in production")
	}
	return nil
}

// DSN returns the MySQL data source name.
func (c *DBConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.User, c.Password, c.Host, c.Port, c.Name)
}

func getStringOrDefault(key, defaultVal string) string {
	val := viper.GetString(key)
	if val == "" {
		return defaultVal
	}
	return val
}

func parseDuration(key string, defaultVal time.Duration) time.Duration {
	val := viper.GetString(key)
	if val == "" {
		return defaultVal
	}
	d, err := time.ParseDuration(val)
	if err != nil {
		return defaultVal
	}
	return d
}
