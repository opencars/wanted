package config

import (
	"log"
	"strconv"

	"github.com/BurntSushi/toml"
)

// Settings is decoded configuration file.
type Settings struct {
	DB Database `toml:"database"`
}

// Database contains configuration details for database.
type Database struct {
	Network    string `toml:"network"`
	Host       string `toml:"host"`
	Port       int    `toml:"port"`
	User       string `toml:"username"`
	Password   string `toml:"password"`
	Name       string `toml:"database"`
	MaxRetries int    `toml:"max_retries"`
	Pool       int    `toml:"pool"`
}

// Address return API address in "host:port" format.
func (db *Database) Address() string {
	return db.Host + ":" + strconv.Itoa(db.Port)
}

// New reads application configuration from specified file path.
func New(path string) (*Settings, error) {
	config := new(Settings)
	if _, err := toml.DecodeFile(path, config); err != nil {
		log.Fatal(err)
	}

	return config, nil
}
