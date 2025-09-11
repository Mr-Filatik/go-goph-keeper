// Package config предоставляет функционал для загрузки конфигурации приложения.
package config

import (
	"flag"
	"fmt"
	"os"
)

// Ключи для поиска флагов.
const (
	FlagServerAddress = "address"

	DescriptionServerAddress = "HTTP server address"
)

// FlagsConfig - структура, содержащая основные переменные окружения для приложения.
type FlagsConfig struct {
	ServerAddress        string // адрес сервера
	ServerAddressIsValue bool
}

// GetConfigFlags получает конфиг из указанных аргументов.
func GetConfigFlags(flagSet *flag.FlagSet, args []string) (*FlagsConfig, error) {
	config := &FlagsConfig{
		ServerAddress:        "",
		ServerAddressIsValue: false,
	}

	argAddress := flagSet.String(FlagServerAddress, "", DescriptionServerAddress)

	if err := flagSet.Parse(args); err != nil {
		return nil, fmt.Errorf("parse argument %w", err)
	}

	if argAddress != nil && *argAddress != "" {
		config.ServerAddress = *argAddress
		config.ServerAddressIsValue = true
	}

	return config, nil
}

// GetConfigFlagsFromOS получает значения флагов из аргументов запуска приложения в ОС.
func GetConfigFlagsFromOS() (*FlagsConfig, error) {
	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	config, err := GetConfigFlags(fs, os.Args[1:])
	if err != nil {
		return nil, fmt.Errorf("get flag config %w", err)
	}

	return config, nil
}

// OverrideConfigFromFlags переопределяет основной конфиг новыми значениями.
func (c *Config) OverrideConfigFromFlags(conf *FlagsConfig) *Config {
	if c == nil || conf == nil {
		return c
	}

	if conf.ServerAddressIsValue {
		c.ServerAddress = conf.ServerAddress
	}

	return c
}
