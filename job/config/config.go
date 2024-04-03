package config

import (
	"errors"
	"github.com/spf13/viper"
	"os"
)

type Config struct {
	Database Database `mapstructure:"database"`
	Aws      Aws      `mapstructure:"aws"`
}

type Database struct {
	Name                 string `mapstructure:"name"`
	Host                 string `mapstructure:"host"`
	Pass                 string `mapstructure:"pass"`
	User                 string `mapstructure:"user"`
	Port                 string `mapstructure:"port"`
	ProductCollection    string `mapstructure:"products"`
	ObjectInfoCollection string `mapstructure:"objectinfo"`
}

type Aws struct {
	Region    string `mapstructure:"region"`
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
	S3        []S3   `mapstructure:"S3"`
}

type S3 struct {
	BucketName string `mapstructure:"BucketName"`
	ObjectKey  string `mapstructure:"ObjectKey"`
}

// LoadDatabase loads database configuration from environment variables.
// It returns an error if any of the required environment variables are not set.
func (c *Config) LoadDatabase() error {
	c.Database.Name = os.Getenv("DB_NAME")
	c.Database.Host = os.Getenv("DB_HOST")
	c.Database.Pass = os.Getenv("DB_PASS")
	c.Database.User = os.Getenv("DB_USER")
	c.Database.Port = os.Getenv("DB_PORT")
	c.Database.ProductCollection = os.Getenv("DB_PRODUCT_COLLECTION")
	c.Database.ObjectInfoCollection = os.Getenv("DB_OBJECTINFO_COLLECTION")
	for _, env := range []string{"DB_NAME", "DB_HOST", "DB_PASS", "DB_USER", "DB_PORT", "DB_PRODUCT_COLLECTION", "DB_OBJECTINFO_COLLECTION"} {
		if os.Getenv(env) == "" {
			return errors.New(env + " is required")
		}
	}
	return nil
}

// LoadAws loads aws configuration from environment variables.
// It returns an error if any of the required environment variables are not set.
func (c *Config) LoadAws() error {
	c.Aws.Region = os.Getenv("AWS_REGION")
	c.Aws.AccessKey = os.Getenv("AWS_ACCESS_KEY_ID")
	c.Aws.SecretKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
	for _, env := range []string{"AWS_REGION", "AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY"} {
		if os.Getenv(env) == "" {
			return errors.New(env + " is required")
		}
	}
	return nil
}

// LoadS3Objects loads S3 objects from configuration file.
// It returns an error if the S3 objects cannot be unmarshalled.
func (c *Config) LoadS3Objects() error {
	viper.SetTypeByDefaultValue(true)
	viper.SetConfigName("s3-objects")
	viper.SetConfigType("yml")
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	viper.AddConfigPath(pwd)
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	if err := viper.Unmarshal(&c.Aws); err != nil {
		return err
	}
	return nil
}

// LoadConfig loads configuration from file.
// It sets initial values for database and aws configurations.
func LoadConfig() (*Config, error) {
	var (
		cfg Config
		err error
	)
	if err = cfg.LoadS3Objects(); err != nil {
		return nil, err
	}
	if err = cfg.LoadDatabase(); err != nil {
		return nil, err
	}
	if err = cfg.LoadAws(); err != nil {
		return nil, err
	}
	return &cfg, nil
}
