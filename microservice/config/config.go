package config

import (
	"errors"
	"os"
	"strconv"
)

// Config struct stores the configuration of the application
type Config struct {
	Database Database `mapstructure:"database"`
	Port     string   `mapstructure:"port"`
}

// Database struct stores the configuration of the database
type Database struct {
	Name       string          `mapstructure:"name"`
	Host       string          `mapstructure:"host"`
	Pass       string          `mapstructure:"pass"`
	User       string          `mapstructure:"user"`
	Port       string          `mapstructure:"port"`
	Collection MongoCollection `mapstructure:"collection"`
}

// MongoCollection struct stores the configuration of the MongoDB collection
type MongoCollection struct {
	Name string `mapstructure:"name"`
}

// LoadDatabase loads the database configuration from the environment variables
func (c *Config) LoadDatabase() error {
	c.Database.Name = os.Getenv("DB_NAME")
	c.Database.Host = os.Getenv("DB_HOST")
	c.Database.Pass = os.Getenv("DB_PASS")
	c.Database.User = os.Getenv("DB_USER")
	c.Database.Port = os.Getenv("DB_PORT")
	c.Database.Collection.Name = os.Getenv("DB_PRODUCT_COLLECTION")
	for _, env := range []string{"DB_NAME", "DB_HOST", "DB_PASS", "DB_USER", "DB_PORT", "DB_PRODUCT_COLLECTION"} {
		if os.Getenv(env) == "" {
			return errors.New(env + " is required")
		}
	}
	return nil
}

// LoadConfig loads the configuration from the environment variables
func LoadConfig() (*Config, error) {
	var Config Config
	if err := Config.LoadDatabase(); err != nil {
		return nil, err
	}
	Config.Port = os.Getenv("PORT")
	if Config.Port == "" {
		return nil, errors.New("PORT is required")
	}
	portNum, err := strconv.Atoi(Config.Port)
	if err != nil {
		return nil, errors.New("PORT must be a number")
	}
	if portNum < 1 || portNum > 65535 {
		return nil, errors.New("PORT must be between 1 and 65535")
	}
	return &Config, nil
}
