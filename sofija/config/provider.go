package config

import (
	"fmt"

	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// Config provider which loads vars from the real config file
var FileProviderSet = wire.NewSet(
	DefaultFileProvider,
	wire.Bind(new(Provider), new(*FileProvider)),
)

// Config provider which loads vars from the test config file (obviously, used for tests)
var TestFileProviderSet = wire.NewSet(
	TestFileProvider,
	wire.Bind(new(Provider), new(*FileProvider)),
)

type Provider interface {
	Load(cfg interface{}) (string, error)
}

// Implementation of the provider interface, which uses viper module to load config vars
type FileProvider struct {
	v viper.Viper
}

func DefaultFileProvider() *FileProvider {
	return NewFileProvider("config")
}

func TestFileProvider() *FileProvider {
	return NewFileProvider("config_test")
}

// Creates a new file provider which reads vars from the given config file
func NewFileProvider(configFile string) *FileProvider {
	v := viper.New()
	v.SetConfigName(configFile)
	v.SetConfigType("yaml")
	// when running from project root
	v.AddConfigPath("secrets/")
	// when running from subdirs (e.g. tests)
	v.AddConfigPath("../../secrets/")
	v.AddConfigPath("../../../secrets/")

	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("error reading viper config %w", err))
	}

	fp := &FileProvider{
		*v,
	}

	return fp
}

// Loads the config file into the provided interface
func (fp FileProvider) Load(cfg interface{}) (string, error) {
	err := fp.v.Unmarshal(cfg)
	if err != nil {
		return "", errors.Wrap(err, "error unmarshalling config")
	}
	return fp.v.ConfigFileUsed(), nil
}
