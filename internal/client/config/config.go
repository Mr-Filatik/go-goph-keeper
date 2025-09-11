// Package config предоставляет функционал для загрузки конфигурации приложения.
package config

import "strings"

// Костанты - значения по умолчанию.
const (
	DefaultServerAddress string = "localhost:8080" // адрес сервера
)

// Config - структура, содержащая основные параметры приложения.
type Config struct {
	ServerAddress string // Aдрес сервера
}

// Initialize создаёт и иницализирует объект *Config.
// Значения присваиваются в следующем порядке (переприсваивают):
//   - значения по умолчания;
//   - значения из флагов командной строки;
//   - значения из переменных окружения.
func Initialize() *Config {
	envsConf := GetConfigEnvsFromOS()
	flagsConf, _ := GetConfigFlagsFromOS()

	config := CreateConfigDefault().
		OverrideConfigFromEnvs(envsConf).
		OverrideConfigFromFlags(flagsConf).
		ValidateConfig()

	return config
}

// CreateConfigDefault создаёт конфиг с дефолтными значениями.
func CreateConfigDefault() *Config {
	config := &Config{
		ServerAddress: DefaultServerAddress,
	}

	return config
}

// ValidateConfig приводит значения конфига к правильному виду.
func (c *Config) ValidateConfig() *Config {
	if c == nil {
		return c
	}

	c.ServerAddress = "http://" + stripHTTPPrefix(c.ServerAddress)

	return c
}

// stripHTTPPrefix обрезает префиксы, для добавления если отсутствует.
func stripHTTPPrefix(addr string) string {
	if strings.HasPrefix(addr, "http://") {
		return addr[7:]
	}

	if strings.HasPrefix(addr, "https://") {
		return addr[8:]
	}

	return addr
}

// overrideConfigCustomValues переопределяет основной конфиг новыми значениями.
// func (c *Config) overrideConfigCustomValues(conf *OtherConfig) *Config {
// 	if c == nil || conf == nil {
// 		return c
// 	}

// 	// logic

// 	return c
// }
