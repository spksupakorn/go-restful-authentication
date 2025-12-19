package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Server  ServerConfig
	MongoDB MongoDBConfig
	JWT     JWTConfig
	Logging LoggingConfig
}

type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Host string `mapstructure:"host"`
	Env  string `mapstructure:"env"`
}

type MongoDBConfig struct {
	URI      string `mapstructure:"uri"`
	Database string `mapstructure:"database"`
	Timeout  int    `mapstructure:"timeout"`
}

type JWTConfig struct {
	SecretKey            string `mapstructure:"secret_key"`
	AccessTokenDuration  int    `mapstructure:"access_token_duration"`  // in hours
	RefreshTokenDuration int    `mapstructure:"refresh_token_duration"` // in hours
}

type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

func LoadConfig() (*Config, error) {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}

	viper.SetConfigName("config." + env)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./internal/config")
	viper.AddConfigPath(".")

	// Read environment variables
	viper.AutomaticEnv()

	// Set defaults
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.env", "development")
	viper.SetDefault("mongodb.uri", "mongodb://localhost:27017")
	viper.SetDefault("mongodb.database", "user_management")
	viper.SetDefault("mongodb.timeout", 10)
	viper.SetDefault("jwt.secret_key", "your-secret-key-change-this")
	viper.SetDefault("jwt.access_token_duration", 24)
	viper.SetDefault("jwt.refresh_token_duration", 168)
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.format", "json")

	// Allow environment variables to override
	viper.BindEnv("server.port", "PORT")
	viper.BindEnv("server.host", "HOST")
	viper.BindEnv("server.env", "APP_ENV")
	viper.BindEnv("mongodb.uri", "MONGODB_URI")
	viper.BindEnv("mongodb.database", "MONGODB_DATABASE")
	viper.BindEnv("mongodb.timeout", "MONGODB_TIMEOUT")
	viper.BindEnv("jwt.secret_key", "JWT_SECRET_KEY")
	viper.BindEnv("jwt.access_token_duration", "JWT_ACCESS_TOKEN_DURATION")
	viper.BindEnv("jwt.refresh_token_duration", "JWT_REFRESH_TOKEN_DURATION")
	viper.BindEnv("logging.level", "LOG_LEVEL")
	viper.BindEnv("logging.format", "LOG_FORMAT")

	// Try to read config file, but don't fail if it doesn't exist
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("Config file not found, using defaults and environment variables")
		} else {
			return nil, err
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
