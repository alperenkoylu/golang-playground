package cache

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"gopkg.in/yaml.v2"
)

var Cache = make(map[string]interface{})

type Config interface {
	LoadFromCache(cache map[string]interface{}) error
}

func Initialize() error {
	configPath := filepath.Join("config")

	files, err := os.ReadDir(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config directory: %v", err)
	}

	for _, f := range files {
		if f.IsDir() || filepath.Ext(f.Name()) != ".yaml" {
			continue
		}

		filePath := filepath.Join(configPath, f.Name())
		data, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read config file '%s': %v", f.Name(), err)
		}

		var cfg map[string]interface{}
		err = yaml.Unmarshal(data, &cfg)
		if err != nil {
			return fmt.Errorf("failed to parse config file '%s': %v", f.Name(), err)
		}

		key := filepath.Base(f.Name())
		key = key[:len(key)-len(filepath.Ext(f.Name()))]
		Cache[key] = cfg
	}

	configKeys := map[string]reflect.Type{
		"database": reflect.TypeOf(Database{}),
	}

	for key, cfg := range configKeys {
		if val, ok := Cache[key]; ok {
			config := reflect.New(cfg).Interface().(Config)
			err := config.LoadFromCache(val.(map[string]interface{}))
			if err != nil {
				return fmt.Errorf("failed to load %s config: %v", key, err)
			}
		}
	}

	return nil
}
