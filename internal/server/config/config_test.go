package config_test

import (
	"flag"
	"testing"

	"github.com/mr-filatik/go-goph-keeper/internal/server/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

/*
	===== GetConfigEnvs =====
*/

type testGetConfigEnvs struct {
	name string
	args map[string]string
	want config.EnvsConfig
}

func createTestsForGetConfigEnvs() []testGetConfigEnvs {
	tests := []testGetConfigEnvs{
		{
			name: "full values",
			args: map[string]string{
				config.EnvKeyHashKey:       "my-hash-key",
				config.EnvKeyServerAddress: "example.com:8080",
				config.EnvKeyDatabase:      "database.one:8080",
				config.EnvKeyCryptoJWTKey:  "secret-jwt-key",
			},
			want: config.EnvsConfig{
				HashKey:              "my-hash-key",
				HashKeyIsValue:       true,
				ServerAddress:        "example.com:8080",
				ServerAddressIsValue: true,
				Database:             "database.one:8080",
				DatabaseIsValue:      true,
				CryptoJWTKey:         "secret-jwt-key",
				CryptoJWTKeyIsValue:  true,
			},
		},
		{
			name: "partial values",
			args: map[string]string{
				config.EnvKeyServerAddress: "example.com:8080",
				config.EnvKeyCryptoJWTKey:  "secret-jwt-key",
			},
			want: config.EnvsConfig{
				HashKey:              "",
				HashKeyIsValue:       false,
				ServerAddress:        "example.com:8080",
				ServerAddressIsValue: true,
				Database:             "",
				DatabaseIsValue:      false,
				CryptoJWTKey:         "secret-jwt-key",
				CryptoJWTKeyIsValue:  true,
			},
		},
		{
			name: "empty values",
			args: map[string]string{},
			want: config.EnvsConfig{
				HashKey:              "",
				HashKeyIsValue:       false,
				ServerAddress:        "",
				ServerAddressIsValue: false,
				Database:             "",
				DatabaseIsValue:      false,
				CryptoJWTKey:         "",
				CryptoJWTKeyIsValue:  false,
			},
		},
	}

	return tests
}

func TestGetConfigEnvs(t *testing.T) {
	t.Parallel()

	tests := createTestsForGetConfigEnvs()

	for index := range tests {
		internalTest := tests[index]
		t.Run(internalTest.name, func(t *testing.T) {
			t.Parallel()

			mockEnv := func(key string) (string, bool) {
				val, ok := internalTest.args[key]

				return val, ok
			}

			config := config.GetConfigEnvs(mockEnv)

			assert.Equal(t, internalTest.want.HashKey, config.HashKey)
			assert.Equal(t, internalTest.want.HashKeyIsValue, config.HashKeyIsValue)

			assert.Equal(t, internalTest.want.ServerAddress, config.ServerAddress)
			assert.Equal(t, internalTest.want.ServerAddressIsValue, config.ServerAddressIsValue)

			assert.Equal(t, internalTest.want.Database, config.Database)
			assert.Equal(t, internalTest.want.DatabaseIsValue, config.DatabaseIsValue)

			assert.Equal(t, internalTest.want.CryptoJWTKey, config.CryptoJWTKey)
			assert.Equal(t, internalTest.want.CryptoJWTKeyIsValue, config.CryptoJWTKeyIsValue)
		})
	}
}

/*
	===== GetConfigFlags =====
*/

type testGetConfigFlags struct {
	name string
	args []string
	want config.FlagsConfig
}

func createTestsForGetConfigFlags() []testGetConfigFlags {
	tests := []testGetConfigFlags{
		{
			name: "full values",
			args: []string{
				"-" + config.FlagHashKey, "my-hash-key",
				"-" + config.FlagServerAddress, "example.com:8080",
				"-" + config.FlagDatabase, "database.one:8080",
				"-" + config.FlagCryptoJWTKey, "secret-jwt-key",
			},
			want: config.FlagsConfig{
				HashKey:              "my-hash-key",
				HashKeyIsValue:       true,
				ServerAddress:        "example.com:8080",
				ServerAddressIsValue: true,
				Database:             "database.one:8080",
				DatabaseIsValue:      true,
				CryptoJWTKey:         "secret-jwt-key",
				CryptoJWTKeyIsValue:  true,
			},
		},
		{
			name: "partial values",
			args: []string{
				"-" + config.FlagServerAddress, "example.com:8080",
				"-" + config.FlagCryptoJWTKey, "secret-jwt-key",
			},
			want: config.FlagsConfig{
				HashKey:              "",
				HashKeyIsValue:       false,
				ServerAddress:        "example.com:8080",
				ServerAddressIsValue: true,
				Database:             "",
				DatabaseIsValue:      false,
				CryptoJWTKey:         "secret-jwt-key",
				CryptoJWTKeyIsValue:  true,
			},
		},
		{
			name: "empty values",
			args: []string{},
			want: config.FlagsConfig{
				HashKey:              "",
				HashKeyIsValue:       false,
				ServerAddress:        "",
				ServerAddressIsValue: false,
				Database:             "",
				DatabaseIsValue:      false,
				CryptoJWTKey:         "",
				CryptoJWTKeyIsValue:  false,
			},
		},
	}

	return tests
}

func TestGetConfigFlags(t *testing.T) {
	t.Parallel()

	tests := createTestsForGetConfigFlags()

	for index := range tests {
		internalTest := tests[index]
		t.Run(internalTest.name, func(t *testing.T) {
			t.Parallel()

			fs := flag.NewFlagSet("test", flag.ContinueOnError)

			config, err := config.GetConfigFlags(fs, internalTest.args)
			require.NoError(t, err)
			require.NotNil(t, config)

			assert.Equal(t, internalTest.want.HashKey, config.HashKey)
			assert.Equal(t, internalTest.want.HashKeyIsValue, config.HashKeyIsValue)

			assert.Equal(t, internalTest.want.ServerAddress, config.ServerAddress)
			assert.Equal(t, internalTest.want.ServerAddressIsValue, config.ServerAddressIsValue)

			assert.Equal(t, internalTest.want.Database, config.Database)
			assert.Equal(t, internalTest.want.DatabaseIsValue, config.DatabaseIsValue)

			assert.Equal(t, internalTest.want.CryptoJWTKey, config.CryptoJWTKey)
			assert.Equal(t, internalTest.want.CryptoJWTKeyIsValue, config.CryptoJWTKeyIsValue)
		})
	}
}

/*
	===== CreateConfigDefault =====
*/

func TestCreateConfigDefault(t *testing.T) {
	t.Parallel()

	defaultConfig := config.CreateConfigDefault()

	assert.Equal(t, config.DefaultServerAddress, defaultConfig.ServerAddress)
	assert.Equal(t, config.DefaultHashKey, defaultConfig.HashKey)
	assert.Equal(t, config.DefaultCryptoJWTKey, defaultConfig.CryptoJWTKey)
	assert.Equal(t, config.DefaultDatabase, defaultConfig.Database)
}

/*
	===== OverrideConfigFromEnvs =====
*/

func TestOverrideConfigFromEnvs(t *testing.T) {
	t.Parallel()

	configEnvs := &config.EnvsConfig{
		CryptoJWTKey:         "test-crypto-jwt-key",
		CryptoJWTKeyIsValue:  true,
		Database:             "test-database",
		DatabaseIsValue:      true,
		ServerAddress:        "",
		ServerAddressIsValue: false,
		HashKey:              "",
		HashKeyIsValue:       false,
	}

	defaultConfig := config.CreateConfigDefault().
		OverrideConfigFromEnvs(configEnvs)

	assert.Equal(t, config.DefaultServerAddress, defaultConfig.ServerAddress)
	assert.Equal(t, config.DefaultHashKey, defaultConfig.HashKey)

	assert.Equal(t, "test-crypto-jwt-key", defaultConfig.CryptoJWTKey)
	assert.Equal(t, "test-database", defaultConfig.Database)
}

/*
	===== OverrideConfigFromFlags =====
*/

func TestOverrideConfigFromFlags(t *testing.T) {
	t.Parallel()

	configFlags := &config.FlagsConfig{
		CryptoJWTKey:         "test-crypto-jwt-key",
		CryptoJWTKeyIsValue:  true,
		Database:             "test-database",
		DatabaseIsValue:      true,
		ServerAddress:        "",
		ServerAddressIsValue: false,
		HashKey:              "",
		HashKeyIsValue:       false,
	}

	defaultConfig := config.CreateConfigDefault().
		OverrideConfigFromFlags(configFlags)

	assert.Equal(t, config.DefaultServerAddress, defaultConfig.ServerAddress)
	assert.Equal(t, config.DefaultHashKey, defaultConfig.HashKey)

	assert.Equal(t, "test-crypto-jwt-key", defaultConfig.CryptoJWTKey)
	assert.Equal(t, "test-database", defaultConfig.Database)
}

/*
	===== ValidateConfig =====
*/

// func TestValidateConfig(t *testing.T) {
// 	t.Parallel()

// 	defaultConfig := &config.Config{
// 		CryptoJWTKey:  "test-crypto-jwt-key",
// 		Database:      "test-database",
// 		ServerAddress: "localhost:8080",
// 		HashKey:       "hash-key",
// 	}

// 	defaultConfig = defaultConfig.ValidateConfig()

// 	assert.Equal(t, "http://localhost:8080", defaultConfig.ServerAddress)

// 	assert.Equal(t, "hash-key", defaultConfig.HashKey)
// 	assert.Equal(t, "test-crypto-jwt-key", defaultConfig.CryptoJWTKey)
// 	assert.Equal(t, "test-database", defaultConfig.Database)
// }
