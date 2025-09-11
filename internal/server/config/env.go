// Package config предоставляет функционал для загрузки конфигурации приложения.
package config

import (
	"os"
)

// Ключи для поиска переменных окружения.
const (
	EnvKeyServerAddress = "ADDRESS"
	EnvKeyHashKey       = "HASH_KEY"
	EnvKeyCryptoJWTKey  = "CRYPTO_JWT_KEY"
	EnvKeyDatabase      = "DATABASE"
)

// EnvsConfig - структура, содержащая основные переменные окружения для приложения.
type EnvsConfig struct {
	CryptoJWTKey         string // путь до публичного ключа
	HashKey              string // ключ хэширования
	ServerAddress        string // адрес сервера
	Database             string // строка подключения к базе данных
	CryptoJWTKeyIsValue  bool
	HashKeyIsValue       bool
	ServerAddressIsValue bool
	DatabaseIsValue      bool
}

// EnvReader — интерфейс для чтения переменных окружения.
type EnvReader func(key string) (string, bool)

// GetConfigEnvs получает значения из универсального хранилища.
func GetConfigEnvs(getenv EnvReader) *EnvsConfig {
	config := &EnvsConfig{
		CryptoJWTKey:         "",
		CryptoJWTKeyIsValue:  false,
		HashKey:              "",
		HashKeyIsValue:       false,
		ServerAddress:        "",
		ServerAddressIsValue: false,
		Database:             "",
		DatabaseIsValue:      false,
	}

	envCryptoKey, envIsValue := getenv(EnvKeyCryptoJWTKey)
	if envIsValue && envCryptoKey != "" {
		config.CryptoJWTKey = envCryptoKey
		config.CryptoJWTKeyIsValue = true
	}

	envKey, envIsValue := getenv(EnvKeyHashKey)
	if envIsValue && envKey != "" {
		config.HashKey = envKey
		config.HashKeyIsValue = true
	}

	envAddress, envIsValue := getenv(EnvKeyServerAddress)
	if envIsValue && envAddress != "" {
		config.ServerAddress = envAddress
		config.ServerAddressIsValue = true
	}

	envDatabase, envIsValue := getenv(EnvKeyDatabase)
	if envIsValue && envDatabase != "" {
		config.Database = envDatabase
		config.DatabaseIsValue = true
	}

	return config
}

// GetConfigEnvsFromOS получает значения из переменных окружения ОС.
func GetConfigEnvsFromOS() *EnvsConfig {
	return GetConfigEnvs(func(key string) (string, bool) {
		value, ok := os.LookupEnv(key)

		return value, ok
	})
}

// OverrideConfigFromEnvs переопределяет основной конфиг новыми значениями.
func (c *Config) OverrideConfigFromEnvs(conf *EnvsConfig) *Config {
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
