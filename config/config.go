package config

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/viper"
)

var (
	ApplicationDir string
)

const (
	// TemplateDir stores the name of the directory that contains templates
	TemplateDir = "../templates"
	// TemplateExt stores the extension used for the template files
	TemplateExt = ".gohtml"

	// StaticDir stores the name of the directory that will serve static files
	StaticDir = "static"

	// StaticPrefix stores the URL prefix used when serving static files
	StaticPrefix = "files"
)

type environment string

const (
	// EnvLocal represents the local environment
	EnvLocal environment = "local"

	// EnvTest represents the test environment
	EnvTest environment = "test"

	// EnvDevelop represents the development environment
	EnvDevelop environment = "dev"

	// EnvStaging represents the staging environment
	EnvStaging environment = "staging"

	// EnvQA represents the qa environment
	EnvQA environment = "qa"

	// EnvProduction represents the production environment
	EnvProduction environment = "prod"
)

// SwitchEnvironment sets the environment variable used to dictate which environment the application is
// currently running in.
// This must be called prior to loading the configuration in order for it to take effect.
func SwitchEnvironment(env environment) {
	if err := os.Setenv("COUPLOAN_APP_ENVIRONMENT", string(env)); err != nil {
		panic(err)
	}
}

type (
	// Config stores complete configuration
	Config struct {
		HTTP          HTTPConfig
		App           AppConfig
		Cache         CacheConfig
		Database      DatabaseConfig
		Loki          LokiConfig
		Pricing       Pricing
		CaptchaSolver CaptchaSolver
		Tests         Tests
	}

	// HTTPConfig stores HTTP configuration
	HTTPConfig struct {
		Hostname       string
		Port           uint16
		ReadTimeout    time.Duration
		WriteTimeout   time.Duration
		IdleTimeout    time.Duration
		TrustedProxies []string
		ProxyHeader    string
		TLS            struct {
			Enabled     bool
			Certificate string
			Key         string
		}
	}

	// AppConfig stores application configuration
	AppConfig struct {
		Name          string
		Environment   environment
		EncryptionKey string
		Timeout       time.Duration
		LogLevel      string
	}

	// CacheConfig stores the cache configuration
	CacheConfig struct {
		Hostname     string
		Port         uint16
		Username     string
		Password     string
		Database     int
		TestDatabase int
		Expiration   struct {
			StaticFile time.Duration
			Page       time.Duration
		}
	}

	// DatabaseConfig stores the database configuration
	DatabaseConfig struct {
		Hostname     string
		Port         uint16
		User         string
		Password     string
		Database     string
		TestDatabase string
	}

	LokiConfig struct {
		// Url of the loki server including http:// or https://
		Url      string
		Username string
		Password string
		App      string
	}

	Pricing struct {
		PricePerRegister int
	}

	CaptchaSolver struct {
		CapSolver  CapSolver
		CapMonster CapMonster
	}

	CapSolver struct {
		ApiKey string
	}

	CapMonster struct {
		ApiKey string
	}

	Tests struct {
		Proxy string
	}
)

// GetConfig loads and returns configuration
func GetConfig() (Config, error) {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	ApplicationDir = filepath.Dir(ex)

	var c Config

	// Load the config file
	v := viper.New()

	// Defaults
	v.SetDefault("app.name", "COUPLOAN")

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("config")
	v.AddConfigPath("../config")
	v.AddConfigPath("../../config")
	v.AddConfigPath("../../../config")

	// Load env variables
	v.SetEnvPrefix("COUPLOAN")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		return c, err
	}

	if err := v.Unmarshal(&c); err != nil {
		return c, err
	}

	return c, nil
}

func EnvDefault(key string, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
