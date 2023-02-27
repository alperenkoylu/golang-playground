package cache

import "fmt"

type Database struct {
	ConnectionString string `yaml:"connection_string"`
	TimeoutAsSeconds int    `yaml:"timeout_as_seconds"`
}

var DatabaseConfig Database

func (db *Database) LoadFromCache(cache map[string]interface{}) error {
	var ok bool

	DatabaseConfig.ConnectionString, ok = cache["connection_string"].(string)
	if !ok {
		return fmt.Errorf("failed to read database connection_string from cache")
	}

	DatabaseConfig.TimeoutAsSeconds, ok = cache["timeout_as_seconds"].(int)
	if !ok {
		return fmt.Errorf("failed to read database timeout_as_seconds from cache")
	}

	return nil
}
