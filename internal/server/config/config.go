// Package config предоставляет функционал для загрузки конфигурации приложения.
package config

// Костанты - значения по умолчанию.
const (
	DefaultServerAddress string = "localhost:8080" // адрес сервера
	DefaultHashKey       string = ""               // ключ хэширования
	DefaultCryptoJWTKey  string = ""               // путь до ключа JWT
	DefaultDatabase      string = ""               // строка подключения к базе данных
)

// Config - структура, содержащая основные параметры приложения.
type Config struct {
	ServerAddress string // Aдрес сервера
	HashKey       string // Ключ хэширования
	CryptoJWTKey  string // Ключ для JWT
	Database      string // Строка подключения к базе данных
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
		HashKey:       DefaultHashKey,
		CryptoJWTKey:  DefaultCryptoJWTKey,
		Database:      DefaultDatabase,
	}

	return config
}

// ValidateConfig приводит значения конфига к правильному виду.
func (c *Config) ValidateConfig() *Config {
	if c == nil {
		return c
	}

	// c.ServerAddress = "http://" + stripHTTPPrefix(c.ServerAddress)

	return c
}

// stripHTTPPrefix обрезает префиксы, для добавления если отсутствует.
// func stripHTTPPrefix(addr string) string {
// 	if strings.HasPrefix(addr, "http://") {
// 		return addr[7:]
// 	}

// 	if strings.HasPrefix(addr, "https://") {
// 		return addr[8:]
// 	}

// 	return addr
// }

// overrideConfigCustomValues переопределяет основной конфиг новыми значениями.
// func (c *Config) overrideConfigCustomValues(conf *OtherConfig) *Config {
// 	if c == nil || conf == nil {
// 		return c
// 	}

// 	// logic

// 	return c
// }
