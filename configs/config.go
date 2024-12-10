package configs

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"sync"
)

type Config struct {
	Enrichers *EnrichersConfig
	Server    *ServerConfig
	Cache     *CacheConfig
}

type ServerConfig struct {
	Host string
	Port int
}

type EnrichersConfig struct {
	Path string
}

type CacheConfig struct {
	Address string
	Password string
	Db int
}

func loadConfig(filePath string) (*Config, error) {
	if filePath == "" {
		return nil, errors.New("config file path is required")
	}
	v := viper.New()
	v.SetConfigFile(filePath)
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	return &config, nil
}

type ConfigManager struct {
	config *Config
	mu     sync.RWMutex
}

func (manager *ConfigManager) LoadConfig(filepath string) error {
	manager.mu.Lock()
	defer manager.mu.Unlock()

	config, err := loadConfig(filepath)
	if err != nil {
		return err
	}

	manager.config = config
	return nil
}

func (manager *ConfigManager) GetConfig(filepath string) (*Config, error) {

	if manager.config == nil {
		err := manager.LoadConfig(filepath)

		if err != nil {
			return nil, err
		}
	}

	return manager.config, nil
}
