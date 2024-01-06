package config

import (
	"os"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

const (
	// DefaultGracefulShutdownTimeout is the default graceful shutdown timeout.
	DefaultGracefulShutdownTimeout = time.Second * 30

	// MaxGracefulShutdownTimeout is the maximum graceful shutdown timeout.
	MaxGracefulShutdownTimeout = time.Minute * 5
)

// Config contains broker server application configuration.
type Config struct {
	SubscriberPort          int           `yaml:"subscriberPort"`
	PublisherPort           int           `yaml:"publisherPort"`
	GracefulShutdownTimeout time.Duration `yaml:"gracefulShutdownTimeout"`
}

// LoadConfig loads the configuration from the given path.
func LoadConfig(path string) (Config, error) {
	reader := &reader{
		osReadFile:    os.ReadFile,
		yamlUnmarshal: yaml.Unmarshal,
	}
	return reader.readConfig(path)
}

type reader struct {
	osReadFile    func(string) ([]byte, error)
	yamlUnmarshal func([]byte, interface{}) error
}

func (r *reader) readConfig(path string) (Config, error) {
	data, err := r.osReadFile(path)
	if err != nil {
		return Config{}, errors.Wrap(err, "read file")
	}

	var config Config
	if err := r.yamlUnmarshal(data, &config); err != nil {
		return Config{}, errors.Wrap(err, "unmarshal yaml")
	}

	if config.GracefulShutdownTimeout == 0 {
		config.GracefulShutdownTimeout = DefaultGracefulShutdownTimeout
	} else if config.GracefulShutdownTimeout > MaxGracefulShutdownTimeout {
		config.GracefulShutdownTimeout = MaxGracefulShutdownTimeout
	}

	return config, nil
}
