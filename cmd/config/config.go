package config

import (
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Config struct {
	App      AppConfigs
	Server   ServerConfigs
	Database DatabaseConfigs
	Cache    CacheConfigs
}

type AppConfigs struct {
	Name    string
	Prefork bool
}

type ServerConfigs struct {
	Port     int
	Instance int
	Domain   string
}

type DatabaseConfigs struct {
	URI      string
	Database string
}

type CacheConfigs struct {
	URI      string
	Database string
	Username string
	Password string
}

var (
	config *Config
	once   sync.Once
)

func LoadConfigs(log *zap.Logger) *Config {
	once.Do(func() {

		if err := godotenv.Load(".env"); err != nil {
			log.Fatal("Critical error: Unable to load environment variables", zap.Error(err))
			os.Exit(1)
		}

		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("cmd/config")
		viper.AutomaticEnv()

		if err := viper.ReadInConfig(); err != nil {
			log.Fatal("Critical error: Unable to load YAML configurations", zap.Error(err))
			os.Exit(1)
		}

		viper.Set("app.name", getEnv("NAME", viper.GetString("app.name")))
		viper.Set("app.prefork", getEnvAsBool("PREFORK", viper.GetBool("app.prefork")))

		viper.Set("server.port", getEnvAsInt("PORT", viper.GetInt("server.port")))
		viper.Set("server.instance", getEnvAsInt("INSTANCE", viper.GetInt("server.instance")))
		viper.Set("server.domain", getEnv("DOMAIN", viper.GetString("server.domain")))

		viper.Set("database.uri", getEnv("DB_URI", viper.GetString("database.uri")))
		viper.Set("database.database", getEnv("DB_NAME", viper.GetString("database.database")))

		viper.Set("cache.uri", getEnv("CACHE_URI", viper.GetString("cache.uri")))
		viper.Set("cache.database", getEnv("CACHE_DB", viper.GetString("cache.database")))
		viper.Set("cache.username", getEnv("CACHE_USER", viper.GetString("cache.username")))
		viper.Set("cache.password", getEnv("CACHE_PASS", viper.GetString("cache.password")))

		var conf Config
		if err := viper.Unmarshal(&conf); err != nil {
			log.Fatal("Critical error: Failed to parse configurations", zap.Error(err))
			os.Exit(1)
		}

		config = &conf
		log.Info("Configurations successfully loaded.")
	})

	return config
}

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if valueStr, exists := os.LookupEnv(key); exists {
		if value, err := strconv.Atoi(valueStr); err == nil {
			return value
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if valueStr, exists := os.LookupEnv(key); exists {
		if value, err := strconv.ParseBool(valueStr); err == nil {
			return value
		}
	}
	return defaultValue
}
