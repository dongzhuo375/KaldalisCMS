package main

import (
	"KaldalisCMS/internal/infra/auth"
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"
)

// Config holds the application configuration
type Config struct {
	Database struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		DBName   string `mapstructure:"dbname"`
		SSLMode  string `mapstructure:"sslmode"`
		TimeZone string `mapstructure:"timezone"`
	} `mapstructure:"database"`
	Auth auth.Config
}

var AppConfig Config

// InitConfig initializes and loads the configuration
func InitConfig() {
	v := viper.GetViper()
	v.SetConfigName("config")        // name of config file (without extension)
	v.SetConfigType("yaml")          // type of the config file
	v.AddConfigPath("./cmd/configs") // path to look for the config file in

	// Optional: set default values
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.user", "your_user")
	v.SetDefault("database.password", "your_password")
	v.SetDefault("database.dbname", "kaldalis_cms")
	v.SetDefault("database.sslmode", "disable")
	v.SetDefault("database.timezone", "Asia/Shanghai")

	// Read config file
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			log.Println("Config file not found, using defaults and environment variables.")
			// Create empty config file if not exists so WriteConfig works? 
            // viper.SafeWriteConfig() might be needed if file doesn't exist.
		} else {
			// Config file was found but another error was produced
			log.Fatalf("Fatal error config file: %s \n", err)
		}
	}

	// Environment variables
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Unmarshal the config into the AppConfig struct
	if err := v.Unmarshal(&AppConfig); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	// auth配置
	refinedAuth, err := auth.LoadConfig(v)
	if err != nil {
		log.Fatalf("Failed to refine auth config: %v", err)
	}
	AppConfig.Auth = *refinedAuth

	log.Println("Configuration loaded successfully.")
}

// GetDatabaseDSN constructs the DSN from the loaded configuration
func GetDatabaseDSN() string {
	db := AppConfig.Database
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		db.Host, db.Port, db.User, db.Password, db.DBName, db.SSLMode, db.TimeZone)
}

// SaveDatabaseConfig updates the database configuration and writes it to the file
func SaveDatabaseConfig(host string, port int, user, pass, dbname string) error {
    v := viper.GetViper()
    v.Set("database.host", host)
    v.Set("database.port", port)
    v.Set("database.user", user)
    v.Set("database.password", pass)
    v.Set("database.dbname", dbname)
    
    // Update the in-memory struct as well
    AppConfig.Database.Host = host
    AppConfig.Database.Port = port
    AppConfig.Database.User = user
    AppConfig.Database.Password = pass
    AppConfig.Database.DBName = dbname

    return v.WriteConfig()
}