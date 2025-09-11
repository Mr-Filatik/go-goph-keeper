package config_test

import (
	"flag"
	"testing"

	"github.com/mr-filatik/go-goph-keeper/internal/client/config"
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
				config.EnvKeyServerAddress: "example.com:8080",
			},
			want: config.EnvsConfig{
				ServerAddress:        "example.com:8080",
				ServerAddressIsValue: true,
			},
		},
		{
			name: "partial values",
			args: map[string]string{
				config.EnvKeyServerAddress: "example.com:8080",
			},
			want: config.EnvsConfig{
				ServerAddress:        "example.com:8080",
				ServerAddressIsValue: true,
			},
		},
		{
			name: "empty values",
			args: map[string]string{},
			want: config.EnvsConfig{
				ServerAddress:        "",
				ServerAddressIsValue: false,
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

			assert.Equal(t, internalTest.want.ServerAddress, config.ServerAddress)
			assert.Equal(t, internalTest.want.ServerAddressIsValue, config.ServerAddressIsValue)
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
				"-" + config.FlagServerAddress, "example.com:8080",
			},
			want: config.FlagsConfig{
				ServerAddress:        "example.com:8080",
				ServerAddressIsValue: true,
			},
		},
		{
			name: "partial values",
			args: []string{
				"-" + config.FlagServerAddress, "example.com:8080",
			},
			want: config.FlagsConfig{
				ServerAddress:        "example.com:8080",
				ServerAddressIsValue: true,
			},
		},
		{
			name: "empty values",
			args: []string{},
			want: config.FlagsConfig{
				ServerAddress:        "",
				ServerAddressIsValue: false,
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

			assert.Equal(t, internalTest.want.ServerAddress, config.ServerAddress)
			assert.Equal(t, internalTest.want.ServerAddressIsValue, config.ServerAddressIsValue)
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
}

/*
	===== OverrideConfigFromEnvs =====
*/

func TestOverrideConfigFromEnvs(t *testing.T) {
	t.Parallel()

	configEnvs := &config.EnvsConfig{
		ServerAddress:        "",
		ServerAddressIsValue: false,
	}

	defaultConfig := config.CreateConfigDefault().
		OverrideConfigFromEnvs(configEnvs)

	assert.Equal(t, config.DefaultServerAddress, defaultConfig.ServerAddress)
}

/*
	===== OverrideConfigFromFlags =====
*/

func TestOverrideConfigFromFlags(t *testing.T) {
	t.Parallel()

	configFlags := &config.FlagsConfig{
		ServerAddress:        "",
		ServerAddressIsValue: false,
	}

	defaultConfig := config.CreateConfigDefault().
		OverrideConfigFromFlags(configFlags)

	assert.Equal(t, config.DefaultServerAddress, defaultConfig.ServerAddress)
}

/*
	===== ValidateConfig =====
*/

func TestValidateConfig(t *testing.T) {
	t.Parallel()

	defaultConfig := &config.Config{
		ServerAddress: "localhost:8080",
	}

	defaultConfig = defaultConfig.ValidateConfig()

	assert.Equal(t, "http://localhost:8080", defaultConfig.ServerAddress)
}
