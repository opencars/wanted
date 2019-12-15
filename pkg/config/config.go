package config

import (
	"strconv"

	"github.com/BurntSushi/toml"
)

// Settings is decoded configuration file.
type Settings struct {
	DB        Database  `toml:"database"`
	Worker    Worker    `toml:"worker"`
	Cleansing Cleansing `toml:"cleansing"`
}

// Database contains configuration details for database.
type Database struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	User     string `toml:"username"`
	Password string `toml:"password"`
	Name     string `toml:"database"`
}

// Worker contains settings for data processing by cmd/worker.
type Worker struct {
	ResourceID string `toml:"resource_id"`
}

// Cleansing contains rules and setting for data cleansing.
type Cleansing struct {
	Brand BrandCleansing `toml:"brand"`
}

// BrandCleansing contains rules and setting for vehicle brand cleansing.
type BrandCleansing struct {
	Matchers []Matcher `toml:"matchers"`
}

// BrandCleansing contains patterns for vehicle brand cleansing.
type Matcher struct {
	Pattern string `toml:"pattern"`
	Maker   string `toml:"maker"`
	Model   string `toml:"model"`
}

// Address return API address in "host:port" format.
func (db *Database) Address() string {
	return db.Host + ":" + strconv.Itoa(db.Port)
}

// New reads application configuration from specified file path.
func New(path string) (*Settings, error) {
	config := &Settings{}
	if _, err := toml.DecodeFile(path, config); err != nil {
		return nil, err
	}

	return config, nil
}
