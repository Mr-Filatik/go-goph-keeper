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
	FlagHashKey       = "hash-key"
	FlagCryptoJWTKey  = "crypto-jwt-key"
	FlagDatabase      = "database"

	DescriptionServerAddress = "HTTP server run address"
	DescriptionHashKey       = "hash key"
	DescriptionCryptoJWTKey  = "crypto key for JWT"
	DescriptionDatabase      = "connection string for database"
)

// FlagsConfig - структура, содержащая основные переменные окружения для приложения.
type FlagsConfig struct {
	CryptoJWTKey         string // путь до публичного ключа
	HashKey              string // ключ хэширования
	ServerAddress        string // адрес сервера
	Database             string // строка подключения к базе данных
	CryptoJWTKeyIsValue  bool
	HashKeyIsValue       bool
	ServerAddressIsValue bool
	DatabaseIsValue      bool
}

// GetConfigFlags получает конфиг из указанных аргументов.
func GetConfigFlags(flagSet *flag.FlagSet, args []string) (*FlagsConfig, error) {
	config := &FlagsConfig{
		CryptoJWTKey:         "",
		CryptoJWTKeyIsValue:  false,
		HashKey:              "",
		HashKeyIsValue:       false,
		ServerAddress:        "",
		ServerAddressIsValue: false,
		Database:             "",
		DatabaseIsValue:      false,
	}

	argCryptoKey := flagSet.String(FlagCryptoJWTKey, "", DescriptionCryptoJWTKey)
	argHashKey := flagSet.String(FlagHashKey, "", DescriptionHashKey)
	argAddress := flagSet.String(FlagServerAddress, "", DescriptionServerAddress)
	argDatabase := flagSet.String(FlagDatabase, "", DescriptionDatabase)

	if err := flagSet.Parse(args); err != nil {
		return nil, fmt.Errorf("parse argument %w", err)
	}

	if argCryptoKey != nil && *argCryptoKey != "" {
		config.CryptoJWTKey = *argCryptoKey
		config.CryptoJWTKeyIsValue = true
	}

	if argHashKey != nil && *argHashKey != "" {
		config.HashKey = *argHashKey
		config.HashKeyIsValue = true
	}

	if argAddress != nil && *argAddress != "" {
		config.ServerAddress = *argAddress
		config.ServerAddressIsValue = true
	}

	if argDatabase != nil && *argDatabase != "" {
		config.Database = *argDatabase
		config.DatabaseIsValue = true
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

	if conf.CryptoJWTKeyIsValue {
		c.CryptoJWTKey = conf.CryptoJWTKey
	}

	if conf.HashKeyIsValue {
		c.HashKey = conf.HashKey
	}

	if conf.ServerAddressIsValue {
		c.ServerAddress = conf.ServerAddress
	}

	if conf.DatabaseIsValue {
		c.Database = conf.Database
	}

	return c
}
