package config

import (
	"os"
	"strconv"

	"gopkg.in/yaml.v2"
)

// Settings is decoded configuration file.
type Settings struct {
	DB        Database  `yaml:"database"`
	Worker    Worker    `yaml:"worker"`
	Cleansing Cleansing `yaml:"cleansing"`
}

// Database contains configuration details for database.
type Database struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"username"`
	Password string `yaml:"password"`
	Name     string `yaml:"database"`
	SSLMode  string `yaml:"ssl_mode"`
}

// Worker contains settings for data processing by cmd/worker.
type Worker struct {
	ResourceID string `yaml:"resource_id"`
}

// Cleansing contains rules and setting for data cleansing.
type Cleansing struct {
	Brand BrandCleansing `yaml:"brand"`
}

// BrandCleansing contains rules and setting for vehicle brand cleansing.
type BrandCleansing struct {
	Matchers []Matcher `yaml:"matchers"`
}

// Matcher contains patterns for vehicle brand cleansing.
type Matcher struct {
	Pattern string `yaml:"pattern"`
	Maker   string `yaml:"maker"`
	Model   string `yaml:"model"`
}

// Address return API address in "host:port" format.
func (db *Database) Address() string {
	return db.Host + ":" + strconv.Itoa(db.Port)
}

// New reads application configuration from specified file path.
func New(path string) (*Settings, error) {
	var config Settings

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	if err := yaml.NewDecoder(f).Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
